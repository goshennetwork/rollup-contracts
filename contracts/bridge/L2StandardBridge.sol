// SPDX-License-Identifier: MIT
// Copyright 2020-2021 Optimism
pragma solidity ^0.8.9;

import "../interfaces/IL1StandardBridge.sol";
import "../interfaces/IL1ERC20Bridge.sol";
import "../interfaces/IL2ERC20Bridge.sol";
import "../interfaces/IL2StandardERC20.sol";
import "@openzeppelin/contracts/utils/introspection/ERC165Checker.sol";
import "../cross-layer/CrossLayerContext.sol";

/**
 * @title L2StandardBridge
 * @dev The L2 Standard bridge is a contract which works together with the L1 Standard bridge to
 * enable ETH and ERC20 transitions between L1 and L2.
 * This contract acts as a minter for new tokens when it hears about deposits into the L1 Standard
 * bridge.
 * This contract also acts as a burner of the tokens intended for withdrawal, informing the L1
 * bridge to release L1 funds.
 */
contract L2StandardBridge is IL2ERC20Bridge, CrossLayerContextUpgradeable {
    address public l1TokenBridge;

    function initialize(address _l2witness, address _l1TokenBridge) public initializer {
        __CrossLayerContext_init(_l2witness);
        l1TokenBridge = _l1TokenBridge;
    }

    modifier onlyEOA() {
        require(msg.sender == tx.origin, "Account not EOA");
        _;
    }

    receive() external payable onlyEOA {
        _initiateETHWithdrawal(msg.sender, msg.sender, msg.value, bytes(""));
    }

    function withdrawETH(bytes calldata _data) external payable onlyEOA {
        _initiateETHWithdrawal(msg.sender, msg.sender, msg.value, _data);
    }

    function withdrawETHTo(address _to, bytes calldata _data) external payable {
        _initiateETHWithdrawal(msg.sender, _to, msg.value, _data);
    }

    function withdraw(address _l2Token, uint256 _amount, bytes calldata _data) external virtual {
        _initiateWithdrawal(_l2Token, msg.sender, msg.sender, _amount, _data);
    }

    function withdrawTo(address _l2Token, address _to, uint256 _amount, bytes calldata _data) external virtual {
        _initiateWithdrawal(_l2Token, msg.sender, _to, _amount, _data);
    }

    /**
     * @dev Performs the logic for withdrawals by burning the token and informing
     *      the L1 token Gateway of the withdrawal.
     * @param _from Account to pull the withdrawal from on L2.
     * @param _to Account to give the withdrawal to on L1.
     * @param _amount Amount of ETH to withdraw.
     * @param _data Optional data to forward to L1. This data is provided
     *        solely as a convenience for external contracts. Aside from enforcing a maximum
     *        length, these contracts provide no guarantees about its content.
     */
    function _initiateETHWithdrawal(address _from, address _to, uint256 _amount, bytes memory _data) internal {
        bytes memory message =
            abi.encodeWithSelector(IL1StandardBridge.finalizeETHWithdrawal.selector, _from, _to, _amount, _data);

        // Send message up to L1 bridge
        sendCrossLayerMessage(l1TokenBridge, message);

        // slither-disable-next-line reentrancy-events
        emit WithdrawalInitiated(address(0), address(0), _from, _to, _amount, _data);
    }

    /**
     * @dev Performs the logic for withdrawals by burning the token and informing
     *      the L1 token Gateway of the withdrawal.
     * @param _l2Token Address of L2 token where withdrawal is initiated.
     * @param _from Account to pull the withdrawal from on L2.
     * @param _to Account to give the withdrawal to on L1.
     * @param _amount Amount of the token to withdraw.
     * @param _data Optional data to forward to L1. This data is provided
     *        solely as a convenience for external contracts. Aside from enforcing a maximum
     *        length, these contracts provide no guarantees about its content.
     */
    function _initiateWithdrawal(address _l2Token, address _from, address _to, uint256 _amount, bytes calldata _data)
        internal
    {
        // slither-disable-next-line reentrancy-events
        IL2StandardERC20(_l2Token).burn(msg.sender, _amount);

        // slither-disable-next-line reentrancy-events
        address l1Token = IL2StandardERC20(_l2Token).l1Token();
        bytes memory message = abi.encodeWithSelector(
            IL1ERC20Bridge.finalizeERC20Withdrawal.selector, l1Token, _l2Token, _from, _to, _amount, _data
        );

        // Send message up to L1 bridge
        // slither-disable-next-line reentrancy-events
        sendCrossLayerMessage(l1TokenBridge, message);

        // slither-disable-next-line reentrancy-events
        emit WithdrawalInitiated(l1Token, _l2Token, msg.sender, _to, _amount, _data);
    }

    function finalizeETHDeposit(address _from, address _to, uint256 _amount, bytes calldata _data)
        external
        virtual
        ensureCrossLayerSender(l1TokenBridge)
    {
        require(address(this).balance >= _amount, "ETH not enough");
        // slither-disable-next-line reentrancy-events
        (bool success,) = _to.call{value: _amount}(new bytes(0));
        require(success, "ETH transfer failed");

        emit DepositFinalized(address(0), address(0), _from, _to, _amount, _data);
    }

    function finalizeERC20Deposit(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _data
    ) external virtual ensureCrossLayerSender(l1TokenBridge) {
        // Check the target token is compliant and
        // verify the deposited token on L1 matches the L2 deposited token representation here
        if (
            ERC165Checker
                // slither-disable-next-line reentrancy-events
                .supportsInterface(_l2Token, 0x1d1d8b63) && _l1Token == IL2StandardERC20(_l2Token).l1Token()
        ) {
            // When a deposit is finalized, we credit the account on L2 with the same amount of
            // tokens.
            // slither-disable-next-line reentrancy-events
            IL2StandardERC20(_l2Token).mint(_to, _amount);
            // slither-disable-next-line reentrancy-events
            emit DepositFinalized(_l1Token, _l2Token, _from, _to, _amount, _data);
        } else {
            // Either the L2 token which is being deposited-into disagrees about the correct address
            // of its L1 token, or does not support the correct interface.
            // This should only happen if there is a  malicious L2 token, or if a user somehow
            // specified the wrong L2 token address to deposit into.
            // In either case, we stop the process here and construct a withdrawal
            // message so that users can get their funds out in some cases.
            // There is no way to prevent malicious token contracts altogether, but this does limit
            // user error and mitigate some forms of malicious contract behavior.
            bytes memory message = abi.encodeWithSelector(
                IL1ERC20Bridge.finalizeERC20Withdrawal.selector,
                _l1Token,
                _l2Token,
                _to, // switched the _to and _from here to bounce back the deposit to the sender
                _from,
                _amount,
                _data
            );

            // Send message up to L1 bridge
            // slither-disable-next-line reentrancy-events
            sendCrossLayerMessage(l1TokenBridge, message);
            // slither-disable-next-line reentrancy-events
            emit DepositFailed(_l1Token, _l2Token, _from, _to, _amount, _data);
        }
    }
}
