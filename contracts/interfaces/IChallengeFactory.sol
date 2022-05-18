// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IStakingManager.sol";
import "../interfaces/IStateTransition.sol";
import "../interfaces/IRollupStateChain.sol";

interface IChallengeFactory {
    event ChallengeStarted(
        uint256 indexed _l2BlockN,
        address indexed _proposer,
        bytes32 _startSystemState,
        uint256 expireAfterBlock
    );

    /**
     * @dev start a challenge game, challenger need to deposit first
     * @param _challengedStateInfo Challenged state
     * @param _parentStateInfo The parent state of challenged state
     * @notice revert if:
     * 1.There exist challenge with challenged state
     * 2.Provide wrong state info
     * 3.token transfer failed
     * @return true if create challenge success
     */
    function newChallange(Types.StateInfo memory _challengedStateInfo, Types.StateInfo memory _parentStateInfo)
        external
        returns (bool);

    ///@return Challenge contract address
    ///@notice revert if not exist challenge to given stateIndex
    function getChallengedContract(uint64 _stateIndex) external view returns (address);

    ///@return StakingManager
    function stakingManager() external view returns (IStakingManager);

    ///@return StateTranstion
    function executor() external view returns (IStateTransition);

    ///@return RollupStateChain
    function rollupStateChain() external view returns (IRollupStateChain);

    ///@return DAO
    function dao() external view returns (address);

    ///@return true if given addr is challenge contract
    function isChallengeContract(address _addr) external view returns (bool);
}
