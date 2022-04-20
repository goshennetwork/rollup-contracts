// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/MerkleMountainRange.sol";
import "../libraries/Constants.sol";
import "../interfaces/IL2CrossLayerMessageWitness.sol";
import "../interfaces/IBuiltinContext.sol";
import "../libraries/Types.sol";
import "../preDeployed/PreDeployed.sol";

contract L2CrossLayerMessageWitness is IL2CrossLayerMessageWitness {
    IBuiltinContext builtinContext = IBuiltinContext(PreDeployed.BUILTIN_CONTEXT);

    MerkleMountainRange.RootNode[] trees;
    bytes32 mmrRoot;
    uint64 totalSize;
    mapping(bytes32 => bool) public successRelayedMessages;
    address private crossDomainMsgSender;

    address l1CrossLayer;
    address owner;

    constructor(address _l1CrossLayer) {
        l1CrossLayer = _l1CrossLayer;
        owner = msg.sender;
    }

    ///used for update L1 contract
    function setL1CrossLayer(address _anotherCrossLayer) public {
        require(msg.sender == owner, "only owner can set crossLayer");
        l1CrossLayer = _anotherCrossLayer;
    }

    function l1Sender() public view returns (address) {
        require(crossDomainMsgSender != address(0), "crossDomainMsgSender not set yet");
        return crossDomainMsgSender;
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
        crossDomainMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossDomainMsgSender = address(0);
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
        uint64 _messageNonce,
        Types.MMRInclusionProof memory _proof
    ) public {
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, _sender, _message, _messageNonce);
        bytes32 _hash = keccak256(_crossLayerCalldata);
        (bytes32 _l1Root, uint64 _totalSize) = builtinContext.l1MMRRoot();
        require(MerkleMountainRange.verifyTrees(_l1Root, _totalSize, _proof.trees), "wrong mmr proof");
        require(
            _hash == _proof.leaf &&
                MerkleMountainRange.verifyLeafTree(
                    _proof.trees,
                    _totalSize,
                    _proof.leafIndex,
                    _proof.siblings,
                    _proof.leaf
                ),
            "wrong inclusion proof"
        );
        require(successRelayedMessages[_hash] == false, "provided message already been relayed");
        crossDomainMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossDomainMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit MessageRelayed(_hash);
        } else {
            emit MessageRelayFailed(_hash);
        }
    }

    function sendMessage(address _target, bytes calldata _message) public {
        //should buy gas
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, msg.sender, _message, totalSize);
        MerkleMountainRange.appendLeafHash(trees, keccak256(_crossLayerCalldata));
        mmrRoot = MerkleMountainRange.genMMRRoot(trees);
        emit MessageSent(_target, msg.sender, _message, totalSize);
        totalSize++;
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
