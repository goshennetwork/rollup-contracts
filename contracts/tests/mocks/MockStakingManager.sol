pragma solidity ^0.8.0;
import "../../interfaces/IStakingManager.sol";
import "@openzeppelin/interfaces/IERC20.sol";

contract MockStakingManager is IStakingManager {
    IERC20 public override token;

    constructor(address _erc20) {
        token = IERC20(_erc20);
    }

    /// Proposer call this function to add collateral, then he can publish block state root.
    function deposit() external override {}

    /// Check whether `_who` is staked.
    function isStaking(address _who) external view override returns (bool) {
        return true;
    }

    /// Proposer ask for unstaking, after this call the proposer can not publish state root any more.
    function startWithdrawal() external override {}

    /// Withdraw to collateral When all block states the proposer published are comfirmed.
    function finalizeWithdrawal() external override {}

    /// Slash the proposer's collateral. can only be called by Challenge contract.
    function slash(
        uint256 _blockHeight,
        bytes32 _stateRoot,
        address _proposer
    ) external override {}

    /// claim slashed collateral. Can only be called by Challenge contract.
    /// @notice revert if 1. new block root not confirmed; 2. the new comfirmed block root
    /// is the same as this proposer's.
    function claim(address _proposer) external override {
        token.transfer(msg.sender, 0.1 ether);
    }

    /// Claim slashed collateral to governance. Can be called by anybody.
    /// @notice revert if 1. new block root not confirmed; 2. the new comfirmed block root
    /// is not the same as this proposer's.
    function claimToGovernance(address _proposer) external override {}
}
