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
    bytes32 public imageStateRoot;
    uint64 public upgradeTime;
    bytes32 public pendingImageStateRoot;
    IAddressResolver public resolver;

    event UpgradeToNewRoot(uint64 timestamp, bytes32 newImageStateRoot);

    function initialize(
        bytes32 _imageStateRoot,
        IAddressResolver _resolver,
        IMachineState _mstate
    ) public initializer {
        imageStateRoot = _imageStateRoot;
        resolver = _resolver;
        mstate = _mstate;
        interpretor = new Interpretor();
        interpretor.initialize(address(_mstate));
    }

    function upgradeToNewRoot(uint64 _upgradeTime, bytes32 newImageStateRoot) public {
        require(msg.sender == address(resolver.dao()), "only dao");
        require(_upgradeTime > block.timestamp, "illegal upgrade time");
        require(newImageStateRoot != bytes32(0), "illegal new root");
        upgradeTime = _upgradeTime;
        pendingImageStateRoot = newImageStateRoot;

        emit UpgradeToNewRoot(_upgradeTime, newImageStateRoot);
    }

    function generateStartState(
        bytes32 rollupInputHash,
        uint64 batchTimestamp,
        bytes32 parentBlockHash
    ) external returns (bytes32) {
        require(msg.sender == address(resolver.challengeFactory()), "only challenge factory");
        bytes32 inputHash = keccak256(abi.encodePacked(rollupInputHash, parentBlockHash));
        if (upgradeTime > 0 && batchTimestamp > upgradeTime) {
            // use pendingImageStateRoot
            Types.StateInfo memory _stateInfo;
            _stateInfo.timestamp = upgradeTime;
            if (resolver.rollupStateChain().isStateConfirmed(_stateInfo)) {
                imageStateRoot = pendingImageStateRoot;
                upgradeTime = 0;
                pendingImageStateRoot = bytes32(0);
            } else {
                return mstate.writeInput(pendingImageStateRoot, inputHash);
            }
        }
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
}
