// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

interface IExecutor {
    function executeNextStep(bytes32 stateHash) external returns (bytes32 nextStateHash);

    //generate initial state, by set value of specific memory to inputHash
    function generateInitialState(bytes32 inputHash) returns (bytes32);
}
