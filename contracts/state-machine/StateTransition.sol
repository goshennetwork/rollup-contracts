// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IMemory.sol";
import "./MemoryLayout.sol";
import "../interfaces/IStateTransition.sol";

contract StateTransition is IStateTransition {
    using MemoryLayout for IMemory;
    bytes32 public imageStateRoot;
    IMemory public riscvMem;

    constructor(bytes32 _imageStateRoot) {
        imageStateRoot = _imageStateRoot;
    }

    function generateStartState(
        uint256 blockNumber,
        bytes32 parentHash,
        bytes32 txhash,
        bytes32 coinbase,
        uint256 gasLimit,
        uint256 timestemp
    ) external returns (bytes32) {
        bytes32 inputHash = keccak256(abi.encodePacked(blockNumber, parentHash, txhash, coinbase, gasLimit, timestemp));
        return riscvMem.writeInput(imageStateRoot, inputHash);
    }

    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external {
        require(riscvMem.isHalt(finalState) == true, "not halted");
        require(riscvMem.mustReadOutput(finalState) == outputRoot, "mismatch root");
    }

    function executeNextStep(bytes32 stateHash) external returns (bytes32 nextStateHash) {
        revert("todo");
    }
}
