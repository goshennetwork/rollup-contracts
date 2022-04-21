// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;
import "../libraries/MerkleMountainRange.sol";
import "../interfaces/IL1CrossLayerMessageWitness.sol";
import "../interfaces/IAddressResolver.sol";
import "../libraries/Types.sol";
import "./CrossLayerCodec.sol";

contract L1CrossLayerMessageWitness is IL1CrossLayerMessageWitness {
    using Types for Types.Block;
    using MerkleMountainRange for CompactMerkleTree;
    IAddressResolver addressResolver;

    CompactMerkleTree compactMerkleTree;
    mapping(bytes32 => bool) public successRelayedMessages;
    mapping(bytes32 => bool) public blockedMessages;
    address private crossLayerMsgSender;

    constructor(address _addressResolver) {
        addressResolver = IAddressResolver(_addressResolver);
    }

    function l2Sender() public view returns (address) {
        require(crossLayerMsgSender != address(0), "crossLayerMsgSender not set");
        return crossLayerMsgSender;
    }

    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        Types.Block memory _block,
        Types.StateInfo memory _stateInfo,
        bytes32[] memory _proof
    ) public {
        require(crossLayerMsgSender == address(0), "reentrancy");
        require(_target != address(addressResolver.rollupInputChain()), "can't relay message to l1 system");
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(_target, _sender, _messageIndex, _message);
        require(addressResolver.rollupStateChain().verifyStateInfo(_stateInfo), "wrong state info");
        require(addressResolver.rollupStateChain().isStateConfirmed(_stateInfo), "state not confirmed yet");
        require(_block.hash() == _stateInfo.blockHash, "wrong block provide");
        MerkleMountainRange.verifyLeafHashInclusion(_hash, _messageIndex, _proof, _block.mmrRoot, _block.mmrSize);
        require(successRelayedMessages[_hash] == false, "provided message already been relayed");
        require(blockedMessages[_hash] == false, "message blocked");
        crossLayerMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossLayerMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit MessageRelayed(_messageIndex, _hash);
        } else {
            emit MessageRelayFailed(_hash, _block.mmrSize, _block.mmrRoot);
        }
    }

    function sendMessage(
        address _target,
        bytes calldata _message,
        uint64 _gasLimit
    ) public {
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
            address(addressResolver.l2CrossLayerMessageWitness()),
            _gasLimit,
            _crossLayerCalldata
        );
    }

    function replayMessage(
        bytes memory _crossLayerCalldata,
        uint64 _queueIndex,
        uint64 _oldGasLimit,
        uint64 _newGasLimit
    ) public {
        (bytes32 _infoHash, ) = addressResolver.rollupInputChain().getQueueTxInfo(_queueIndex);
        // same as rollupInputChain
        bytes32 _txHash = keccak256(
            abi.encode(
                address(this),
                address(addressResolver.l2CrossLayerMessageWitness()),
                _oldGasLimit,
                _crossLayerCalldata
            )
        );
        require(_txHash == _infoHash, "message not in queue");
        addressResolver.rollupInputChain().enqueue(
            address(addressResolver.l2CrossLayerMessageWitness()),
            _newGasLimit,
            _crossLayerCalldata
        );
    }

    function blockMessage(bytes32[] memory _messageHashes) public {
        require(msg.sender == addressResolver.dao(), "only dao allowed");
        for (uint256 i = 0; i < _messageHashes.length; i++) {
            blockedMessages[_messageHashes[i]] = true;
        }
    }

    function allowMessage(bytes32[] memory _messageHashes) public {
        require(msg.sender == addressResolver.dao(), "only dao allowed");
        for (uint256 i = 0; i < _messageHashes.length; i++) {
            blockedMessages[_messageHashes[i]] = false;
        }
    }

    function mmrRoot() public view returns (bytes32) {
        return compactMerkleTree.rootHash;
    }

    function totalSize() public view returns (uint64) {
        return compactMerkleTree.treeSize;
    }
}
