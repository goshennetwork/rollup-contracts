// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Constants.sol";
import { MerkleMountainRange, CompactMerkleTree } from "../libraries/MerkleMountainRange.sol";
import "../interfaces/IL2CrossLayerMessageWitness.sol";
import "../interfaces/IBuiltinContext.sol";
import "../libraries/Types.sol";
import "../predeployed/PreDeployed.sol";

contract L2CrossLayerMessageWitness is IL2CrossLayerMessageWitness {
    using MerkleMountainRange for CompactMerkleTree;
    IBuiltinContext builtinContext = IBuiltinContext(PreDeployed.BUILTIN_CONTEXT);

    CompactMerkleTree compactMerkleTree;
    mapping(bytes32 => bool) public successRelayedMessages;
    address private crossLayerMsgSender;

    function l1Sender() public view returns (address) {
        require(crossLayerMsgSender != address(0), "crossDomainMsgSender not set yet");
        return crossLayerMsgSender;
    }

    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex
    ) public {
        require(msg.sender == Constants.L1_CROSS_LAYER_WITNESS, "wrong sender");
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, _sender, _message, _messageIndex);
        bytes32 _hash = keccak256(_crossLayerCalldata);
        require(successRelayedMessages[_hash] == false, "already relayed");
        crossLayerMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossLayerMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit MessageRelayed(_hash);
        } else {
            emit MessageRelayFailed(_hash);
        }
    }

    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes32[] memory _proof
    ) public {
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, _sender, _message, _messageIndex);
        bytes32 _hash = keccak256(_crossLayerCalldata);
        (bytes32 _l1Root, uint64 _totalSize) = builtinContext.l1MMRRoot();
        MerkleMountainRange.verifyLeafHashInclusion(_hash, _messageIndex, _proof, _l1Root, _totalSize);
        require(successRelayedMessages[_hash] == false, "provided message already been relayed");
        crossLayerMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossLayerMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit MessageRelayed(_hash);
        } else {
            emit MessageRelayFailed(_hash);
        }
    }

    function sendMessage(address _target, bytes calldata _message) public {
        uint64 _messageIndex = compactMerkleTree.treeSize;
        //should buy gas
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, msg.sender, _message, _messageIndex);
        compactMerkleTree.appendLeafHash(keccak256(_crossLayerCalldata));
        emit MessageSent(_target, msg.sender, _message, _messageIndex);
    }

    function _encodeCrossLayerCallData(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "relayMessage(address,address,bytes,uint64)",
                _target,
                _sender,
                _message,
                _messageIndex
            );
    }
}
