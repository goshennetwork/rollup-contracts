// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../interfaces/IMachineState.sol";
import "./MemoryLayout.sol";
import "../interfaces/IStateTransition.sol";
import "../interfaces/IAddressResolver.sol";
import "./MachineState.sol";
import "./riscv32/Interpretor.sol";

contract StateTransition is IStateTransition, Initializable {
    using MemoryLayout for IMachineState;
    IMachineState public mstate;
    Interpretor interpretor;
    bytes32 public imageStateRoot;
    uint256 public upgradeHeight;
    bytes32 public pendingImageStateRoot;
    IAddressResolver public resolver;

    event UpgradeToNewRoot(uint256 blockNumber, bytes32 newImageStateRoot);

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

    function upgradeToNewRoot(uint256 blockNumber, bytes32 newImageStateRoot) public {
        require(msg.sender == resolver.dao(), "only dao");
        require(blockNumber > resolver.rollupStateChainContainer().chainSize(), "illegal height");
        require(newImageStateRoot != bytes32(0), "illegal new root");
        upgradeHeight = blockNumber;
        pendingImageStateRoot = newImageStateRoot;

        emit UpgradeToNewRoot(blockNumber, newImageStateRoot);
    }

    function generateStartState(
        bytes32 rollupInputHash,
        uint64 blockNumber,
        bytes32 parentBlockHash
    ) external returns (bytes32) {
        require(msg.sender == address(resolver.challengeFactory()), "only challenge factory");
        bytes32 inputHash = keccak256(abi.encodePacked(rollupInputHash, parentBlockHash));
        if (upgradeHeight > 0 && blockNumber >= upgradeHeight) {
            imageStateRoot = pendingImageStateRoot;
            upgradeHeight = 0;
            pendingImageStateRoot = bytes32(0);
        }
        return mstate.writeInput(imageStateRoot, inputHash);
    }

    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external view {
        require(mstate.isHalt(finalState) == true, "not halted");
        require(mstate.mustReadOutput(finalState) == outputRoot, "mismatch root");
    }

    function executeNextStep(bytes32 stateHash) external returns (bytes32 nextStateHash) {
        (nextStateHash,) = interpretor.step(stateHash);
        return nextStateHash;
    }
}
