// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../libraries/MerkleMountainRange.sol";
import "../interfaces/IL1CrossLayerWitness.sol";
import "../interfaces/IAddressResolver.sol";
import "../libraries/Types.sol";
import "./CrossLayerCodec.sol";

contract L1CrossLayerWitness is IL1CrossLayerWitness, Initializable {
    using Types for Types.Block;
    using MerkleMountainRange for CompactMerkleTree;
    IAddressResolver addressResolver;

    CompactMerkleTree compactMerkleTree;
    mapping(bytes32 => bool) public successRelayedMessages;
    mapping(bytes32 => bool) public blockedMessages;
    address private crossLayerMsgSender;

    function initialize(address _addressResolver) public initializer {
        addressResolver = IAddressResolver(_addressResolver);
    }

    function crossLayerSender() external view returns (address) {
        require(crossLayerMsgSender != address(0), "no cross layer sender");
        return crossLayerMsgSender;
    }

    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes memory _rlpHeader,
        Types.StateInfo memory _stateInfo,
        bytes32[] memory _proof
    ) public {
        require(crossLayerMsgSender == address(0), "reentrancy");
        require(_target != address(addressResolver.rollupInputChain()), "can't relay message to l1 system");
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(_target, _sender, _messageIndex, _message);
        require(addressResolver.rollupStateChain().verifyStateInfo(_stateInfo), "wrong state info");
        require(addressResolver.rollupStateChain().isStateConfirmed(_stateInfo), "state not confirmed yet");
        require(keccak256(_rlpHeader) == _stateInfo.blockHash, "wrong block provide");
        (bytes32 _mmrRoot, uint64 _mmrSize) = Types.decodeMMRFromRlpHeader(_rlpHeader);
        MerkleMountainRange.verifyLeafHashInclusion(_hash, _messageIndex, _proof, _mmrRoot, _mmrSize);
        require(successRelayedMessages[_hash] == false, "provided message already been relayed");
        require(blockedMessages[_hash] == false, "message blocked");
        crossLayerMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossLayerMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit MessageRelayed(_messageIndex, _hash);
        } else {
            emit MessageRelayFailed(_hash, _mmrSize, _mmrRoot);
        }
    }

    function sendMessage(address _target, bytes calldata _message) public {
        require(msg.sender != address(this), "wired situation");
        uint64 treeSize = compactMerkleTree.treeSize;
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(_target, msg.sender, treeSize, _message);
        compactMerkleTree.appendLeafHash(_hash);
        bytes memory _crossLayerCalldata = CrossLayerCodec.encodeL1ToL2CallData(
            _target,
            msg.sender,
            _message,
            treeSize,
            compactMerkleTree.rootHash,
            treeSize + 1
        );
        addressResolver.rollupInputChain().enqueue(
            address(addressResolver.l2CrossLayerWitness()),
            0,
            _crossLayerCalldata
        );
        emit MessageSent(treeSize, _target, msg.sender, _message);
    }

    function blockMessage(bytes32[] memory _messageHashes) public {
        require(msg.sender == addressResolver.dao(), "only dao allowed");
        for (uint256 i = 0; i < _messageHashes.length; i++) {
            blockedMessages[_messageHashes[i]] = true;
        }
        emit MessageBlocked(_messageHashes);
    }

    function allowMessage(bytes32[] memory _messageHashes) public {
        require(msg.sender == addressResolver.dao(), "only dao allowed");
        for (uint256 i = 0; i < _messageHashes.length; i++) {
            blockedMessages[_messageHashes[i]] = false;
        }
        emit MessageAllowed(_messageHashes);
    }

    function mmrRoot() public view returns (bytes32) {
        return compactMerkleTree.rootHash;
    }

    function totalSize() public view returns (uint64) {
        return compactMerkleTree.treeSize;
    }
}
