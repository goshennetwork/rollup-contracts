// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

interface IStateTransition {
    function upgradeToNewRoot(uint256 blockNumber, bytes32 newImageStateRoot) external;

    /**
    * @dev Only Challenge factory permitted, because it acts like a button to switch different version of system
    * @param rollupInputHash RollupInput hash in RollupInputChain
    * @param blockNumber state's block number in RollupStateChain(same as RollupInputChain)
    * @param parentBlockHash Parent block's hash

    */
    function generateStartState(
        bytes32 rollupInputHash,
        uint64 blockNumber,
        bytes32 parentBlockHash
    ) external returns (bytes32);

    ///@dev validate final state, revert if final state is not halt or output inconsistent root
    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external;

    ///@dev Exec one step transition
    function executeNextStep(bytes32 stateHash) external returns (bytes32 nextStateHash, bool halt);
}
