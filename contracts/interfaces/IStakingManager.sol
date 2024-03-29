// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/interfaces/IERC20.sol";
import "../libraries/Types.sol";

interface IStakingManager {
    //proposer deposit for staking.
    event Deposited(address indexed proposer, uint256 amount);
    //proposer start withdraw.
    event WithdrawStarted(address indexed proposer, uint256 needComfirmedBlock);
    //proposer finalize withdraw.
    event WithdrawFinalized(address indexed proposer, uint256 amount);
    //challenger slash the proposer
    event DepositSlashed(address indexed proposer, address indexed challenger, uint256 blockHeight, bytes32 _blockHash);
    //challenger or DAO gets deposit.
    event DepositClaimed(address indexed proposer, address indexed receiver, uint256 amount);

    enum StakingState {
        // Before depositing or after getting slashed, a user is unstaked
        UNSTAKED,
        // After depositing, a user is staking
        STAKING,
        // After a user has applied a withdrawal
        WITHDRAWING,
        // After challenge success, the staked token is waiting to be distributed to the challenger.
        SLASHING
    }

    /// record a proposer's staking info
    struct StakingInfo {
        // The user's state
        StakingState state;
        // After which confirmed state height the proposer can withdraw
        uint64 needConfirmedHeight;
        // The L1 time of the first successful challenge.
        uint64 firstSlashTime;
        // The state chain height
        uint64 earliestChallengeHeight;
        // The earliest observed block hash for this bond which has had fraud
        bytes32 earliestChallengeBlockHash;
    }

    /// The token address used for staking.
    function token() external view returns (IERC20);

    /// Proposer call this function to add collateral, then he can publish block state root.
    function deposit() external;

    /// Check whether `_who` is staked.
    function isStaking(address _who) external view returns (bool);

    /// Proposer ask for unstaking, after this call the proposer can not publish state info any more.
    function startWithdrawal() external;

    /// Withdraw to collateral When all  state the proposer published are confirmed.
    function finalizeWithdrawal(Types.StateInfo memory _stateInfo) external;

    /// Slash the proposer's collateral. can only be called by Challenge contract.
    function slash(
        uint64 _chainHeight,
        bytes32 _stateRoot,
        address _proposer
    ) external;

    /// claim slashed collateral. Can only be called by Challenge contract.
    /// @notice revert if 1. new block hash not confirmed; 2. the new confirmed block hash
    /// is the same as this proposer's.
    function claim(address _proposer, Types.StateInfo memory _stateInfo) external;

    /// Claim slashed collateral to governance. Can be called by anybody.
    /// @notice revert if 1. new block hash not confirmed; 2. the new confirmed block hash
    /// is not the same as this proposer's.
    function claimToGovernance(address _proposer, Types.StateInfo memory _stateInfo) external;
}
