// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../libraries/Constants.sol";
import "../libraries/MerkleMountainRange.sol";
import "../interfaces/IL2CrossLayerWitness.sol";
import "../libraries/Types.sol";
import "./CrossLayerCodec.sol";

contract L2CrossLayerWitness is IL2CrossLayerWitness, Initializable {
    using MerkleMountainRange for CompactMerkleTree;
    CompactMerkleTree compactMerkleTree;
    mapping(bytes32 => bool) public successRelayedMessages;
    mapping(uint64 => bytes32) public mmrRoots;
    address private crossLayerMsgSender;

    function initialize() public initializer {}

    function crossLayerSender() external view returns (address) {
        require(crossLayerMsgSender != address(0), "no cross layer sender");
        return crossLayerMsgSender;
    }

    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes32 _mmrRoot,
        uint64 _mmrSize
    ) public returns (bool) {
        require(crossLayerMsgSender == address(0), "reentrancy");
        require(msg.sender == Constants.L1_CROSS_LAYER_WITNESS, "wrong sender");
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(_target, _sender, _messageIndex, _message);
        require(successRelayedMessages[_hash] == false, "already relayed");
        crossLayerMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossLayerMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit MessageRelayed(_messageIndex, _hash);
        } else {
            mmrRoots[_mmrSize] = _mmrRoot;
            emit MessageRelayFailed(_messageIndex, _hash, _mmrSize, _mmrRoot);
        }
        return success;
    }

    function replayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes32[] memory _proof,
        uint64 _mmrSize
    ) public returns (bool) {
        require(crossLayerMsgSender == address(0), "reentrancy");
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(_target, _sender, _messageIndex, _message);
        bytes32 _mmrRoot = mmrRoots[_mmrSize];
        require(_mmrRoot != bytes32(0), "unknown mmr root");
        MerkleMountainRange.verifyLeafHashInclusion(_hash, _messageIndex, _proof, _mmrRoot, _mmrSize);
        require(successRelayedMessages[_hash] == false, "message already relayed");
        crossLayerMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossLayerMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit MessageRelayed(_messageIndex, _hash);
        } else {
            emit MessageRelayFailed(_messageIndex, _hash, _mmrSize, _mmrRoot);
        }
        return success;
    }

    function sendMessage(address _target, bytes calldata _message) public {
        require(msg.sender != address(this), "wired situation");
        uint64 _messageIndex = compactMerkleTree.treeSize;
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(_target, msg.sender, _messageIndex, _message);
        bytes32 _mmrRoot = compactMerkleTree.appendLeafHash(_hash);
        emit MessageSent(_messageIndex, _target, msg.sender, _mmrRoot, _message);
    }
}
