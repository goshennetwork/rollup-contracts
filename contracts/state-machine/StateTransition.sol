// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../interfaces/IMachineState.sol";
import "./MemoryLayout.sol";
import "../interfaces/IStateTransition.sol";
import "../interfaces/IAddressResolver.sol";
import "./MachineState.sol";
import "./riscv32/Interpretor.sol";
import "../libraries/Types.sol";

contract StateTransition is IStateTransition, Initializable {
    using MemoryLayout for IMachineState;
    IMachineState public mstate;
    Interpretor interpretor;
    IAddressResolver public resolver;

    bytes32[] private imageStateRoots;
    uint64[] private upgradeBatchIndexes;

    event UpgradeToNewRoot(uint64 upgradeBatchIndex, bytes32 newImageStateRoot);

    function initialize(
        bytes32 _imageStateRoot,
        IAddressResolver _resolver,
        IMachineState _mstate
    ) public initializer {
        imageStateRoots.push(_imageStateRoot);
        upgradeBatchIndexes.push(0);
        resolver = _resolver;
        mstate = _mstate;
        interpretor = new Interpretor();
        interpretor.initialize(address(_mstate));
    }

    function upgradeToNewRoot(uint64 upgradeBatchIndex, bytes32 newImageStateRoot) public {
        require(msg.sender == address(resolver.dao()), "only dao");
        require(upgradeBatchIndex >= upgradeBatchIndexes[upgradeBatchIndexes.length - 1], "duplicated upgrade");
        require(upgradeBatchIndex >= resolver.rollupStateChainContainer().chainSize(), "ill batch index");
        require(newImageStateRoot != bytes32(0), "illegal new root");
        imageStateRoots.push(newImageStateRoot);
        upgradeBatchIndexes.push(upgradeBatchIndex);

        emit UpgradeToNewRoot(upgradeBatchIndex, newImageStateRoot);
    }

    function generateStartState(
        bytes32 rollupInputHash,
        uint64 batchIndex,
        bytes32 parentBlockHash
    ) external returns (bytes32) {
        require(msg.sender == address(resolver.challengeFactory()), "only challenge factory");
        bytes32 inputHash = keccak256(abi.encodePacked(rollupInputHash, parentBlockHash));
        bytes32 imageStateRoot = getImageRoot(batchIndex);
        return mstate.writeInput(imageStateRoot, inputHash);
    }

    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external view {
        require(mstate.isHalt(finalState) == true, "not halted");
        require(mstate.mustReadOutput(finalState) == outputRoot, "mismatch root");
    }

    function executeNextStep(bytes32 stateHash) external returns (bytes32 nextStateHash) {
        (nextStateHash, ) = interpretor.step(stateHash);
        return nextStateHash;
    }

    function getImageRoot(uint64 batchIndex) public view override returns (bytes32) {
        uint256 stateRootsNum = imageStateRoots.length;
        bytes32 result = imageStateRoots[0];
        for (uint256 i = 0; i < stateRootsNum; i++) {
            if (upgradeBatchIndexes[i] <= batchIndex) {
                result = imageStateRoots[i];
            } else {
                return result;
            }
        }
        return result;
    }
}
