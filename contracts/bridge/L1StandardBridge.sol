// SPDX-License-Identifier: MIT
// Copyright 2020-2021 Optimism
pragma solidity ^0.8.9;

/* Interface Imports */
import "../interfaces/IL1StandardBridge.sol";
import "../interfaces/IL1ERC20Bridge.sol";
import "../interfaces/IL2ERC20Bridge.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/* Library Imports */
import "../cross-layer/CrossLayerContext.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title L1StandardBridge
 * @dev The L1 ETH and ERC20 Bridge is a contract which stores deposited L1 funds and standard
 * tokens that are in use on L2. It synchronizes a corresponding L2 Bridge, informing it of deposits
 * and listening to it for newly finalized withdrawals.
 *
 */
contract L1StandardBridge is IL1StandardBridge, CrossLayerContext {
    using SafeERC20 for IERC20;

    address public l2TokenBridge;

    // Maps L1 token to L2 token to balance of the L1 token deposited
    mapping(address => mapping(address => uint256)) public deposits;

    // This contract lives behind a proxy, so the constructor parameters will go unused.
    constructor() CrossLayerContext(address(0)) {}

    // slither-disable-next-line external-function
    function initialize(address _l1witness, address _l2TokenBridge) public {
        require(address(crossLayerWitness) == address(0), "already initialized.");
        crossLayerWitness = ICrossLayerWitness(_l1witness);
        l2TokenBridge = _l2TokenBridge;
    }

    modifier onlyEOA() {
        require(msg.sender == tx.origin, "Account not EOA");
        _;
    }

    receive() external payable onlyEOA {
        _initiateETHDeposit(msg.sender, msg.sender, bytes(""));
    }

    function depositETH(bytes calldata _data) external payable onlyEOA {
        _initiateETHDeposit(msg.sender, msg.sender, _data);
    }

    function depositETHTo(address _to, bytes calldata _data) external payable {
        _initiateETHDeposit(msg.sender, _to, _data);
    }

    /**
     * @dev Performs the logic for deposits by storing the ETH and informing the L2 ETH Gateway of
     * the deposit.
     * @param _from Account to pull the deposit from on L1.
     * @param _to Account to give the deposit to on L2.
     * @param _data Optional data to forward to L2. This data is provided
     *        solely as a convenience for external contracts. Aside from enforcing a maximum
     *        length, these contracts provide no guarantees about its content.
     */
    function _initiateETHDeposit(
        address _from,
        address _to,
        bytes memory _data
    ) internal {
        // Construct calldata for finalizeDeposit call
        bytes memory message = abi.encodeWithSelector(
            IL2ERC20Bridge.finalizeETHDeposit.selector,
            _from,
            _to,
            msg.value,
            _data
        );

        // Send calldata into L2
        // slither-disable-next-line reentrancy-events
        sendCrossLayerMessage(l2TokenBridge, message);

        // slither-disable-next-line reentrancy-events
        emit ETHDepositInitiated(_from, _to, msg.value, _data);
    }

    function depositERC20(
        address _l1Token,
        address _l2Token,
        uint256 _amount,
        bytes calldata _data
    ) external virtual onlyEOA {
        _initiateERC20Deposit(_l1Token, _l2Token, msg.sender, msg.sender, _amount, _data);
    }

    function depositERC20To(
        address _l1Token,
        address _l2Token,
        address _to,
        uint256 _amount,
        bytes calldata _data
    ) external virtual {
        _initiateERC20Deposit(_l1Token, _l2Token, msg.sender, _to, _amount, _data);
    }

    /**
     * @dev Performs the logic for deposits by informing the L2 Deposited Token
     * contract of the deposit and calling a handler to lock the L1 funds. (e.g. transferFrom)
     *
     * @param _l1Token Address of the L1 ERC20 we are depositing
     * @param _l2Token Address of the L1 respective L2 ERC20
     * @param _from Account to pull the deposit from on L1
     * @param _to Account to give the deposit to on L2
     * @param _amount Amount of the ERC20 to deposit.
     * @param _data Optional data to forward to L2. This data is provided
     *        solely as a convenience for external contracts. Aside from enforcing a maximum
     *        length, these contracts provide no guarantees about its content.
     */
    function _initiateERC20Deposit(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _data
    ) internal {
        // When a deposit is initiated on L1, the L1 Bridge transfers the funds to itself for future
        // withdrawals. safeTransferFrom also checks if the contract has code, so this will fail if
        // _from is an EOA or address(0).
        // slither-disable-next-line reentrancy-events, reentrancy-benign
        IERC20(_l1Token).safeTransferFrom(_from, address(this), _amount);

        // Construct calldata for _l2Token.finalizeDeposit(_to, _amount)
        bytes memory message = abi.encodeWithSelector(
            IL2ERC20Bridge.finalizeERC20Deposit.selector,
            _l1Token,
            _l2Token,
            _from,
            _to,
            _amount,
            _data
        );

        // Send calldata into L2
        // slither-disable-next-line reentrancy-events, reentrancy-benign
        sendCrossLayerMessage(l2TokenBridge, message);

        // slither-disable-next-line reentrancy-benign
        deposits[_l1Token][_l2Token] = deposits[_l1Token][_l2Token] + _amount;

        // slither-disable-next-line reentrancy-events
        emit ERC20DepositInitiated(_l1Token, _l2Token, _from, _to, _amount, _data);
    }

    function finalizeETHWithdrawal(
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _data
    ) external ensureCrossLayerSender(l2TokenBridge) {
        // slither-disable-next-line reentrancy-events
        (bool success, ) = _to.call{ value: _amount }(new bytes(0));
        require(success, "ETH transfer failed");

        // slither-disable-next-line reentrancy-events
        emit ETHWithdrawalFinalized(_from, _to, _amount, _data);
    }

    /**
     * @inheritdoc IL1ERC20Bridge
     */
    function finalizeERC20Withdrawal(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _data
    ) external ensureCrossLayerSender(l2TokenBridge) {
        deposits[_l1Token][_l2Token] = deposits[_l1Token][_l2Token] - _amount;

        // When a withdrawal is finalized on L1, the L1 Bridge transfers the funds to the withdrawer
        // slither-disable-next-line reentrancy-events
        IERC20(_l1Token).safeTransfer(_to, _amount);

        // slither-disable-next-line reentrancy-events
        emit ERC20WithdrawalFinalized(_l1Token, _l2Token, _from, _to, _amount, _data);
    }

    /*****************************
     * Temporary - Migrating ETH *
     *****************************/

    /**
     * @dev Adds ETH balance to the account. This is meant to allow for ETH
     * to be migrated from an old gateway to a new gateway.
     * NOTE: This is left for one upgrade only so we are able to receive the migrated ETH from the
     * old contract
     */
    function donateETH() external payable {}
}
