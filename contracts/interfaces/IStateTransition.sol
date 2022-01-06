// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

interface IStateTransition {
    function executeNextStep(bytes32 stateHash) external returns (bytes32 nextStateHash);

    function generateStartState(
        bytes32 parentHash,
        bytes32 txhash,
        bytes32 coinbase,
        uint256 gasLimit,
        uint256 timestemp
    ) external returns (bytes32);

    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external;
}
