// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IMemory.sol";
import "./MemoryLayout.sol";
import "../interfaces/IStateTransition.sol";
import "../interfaces/IAddressResolver.sol";

contract StateTransition is IStateTransition {
    using MemoryLayout for IMemory;
    bytes32 public imageStateRoot;
    IMemory public riscvMem;
    uint256 public upgradeHeight;
    bytes32 public pendingImageStateRoot;
    IAddressResolver public resolver;

    event UpgradeToNewRoot(uint256 blockNumber, bytes32 newImageStateRoot);

    constructor(bytes32 _imageStateRoot, IAddressResolver _resolver) {
        imageStateRoot = _imageStateRoot;
        resolver = _resolver;
    }

    function upgradeToNewRoot(uint256 blockNumber, bytes32 newImageStateRoot) public {
        require(msg.sender == resolver.dao(), "only dao");
        require(upgradeHeight == 0, "upgrading");
        require(blockNumber > resolver.rollupStateChainContainer().chainSize(), "illegal height");
        require(newImageStateRoot != bytes32(0), "illegal new root");
        upgradeHeight = blockNumber;
        pendingImageStateRoot = newImageStateRoot;

        emit UpgradeToNewRoot(blockNumber, newImageStateRoot);
    }

    function generateStartState(
        uint256 blockNumber,
        bytes32 parentHash,
        bytes32 txhash,
        bytes32 coinbase,
        uint256 gasLimit,
        uint256 timestemp
    ) external returns (bytes32) {
        require(msg.sender == address(resolver.challengeFactory()), "only challenge factory");
        bytes32 inputHash = keccak256(abi.encodePacked(blockNumber, parentHash, txhash, coinbase, gasLimit, timestemp));
        if (upgradeHeight > 0 && blockNumber >= upgradeHeight) {
            imageStateRoot = pendingImageStateRoot;
            upgradeHeight = 0;
            pendingImageStateRoot = bytes32(0);
        }
        return riscvMem.writeInput(imageStateRoot, inputHash);
    }

    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external view {
        require(riscvMem.isHalt(finalState) == true, "not halted");
        require(riscvMem.mustReadOutput(finalState) == outputRoot, "mismatch root");
    }

    // TODO: only challenge contract
    function executeNextStep(bytes32 stateHash) external pure returns (bytes32 nextStateHash) {
        //fix warning
        stateHash;
        nextStateHash;
        revert("todo");
    }
}
