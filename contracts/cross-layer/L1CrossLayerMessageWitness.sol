// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;
import { MerkleMountainRange, CompactMerkleTree } from "../libraries/MerkleMountainRange.sol";
import "../interfaces/IL1CrossLayerMessageWitness.sol";
import "../interfaces/IAddressResolver.sol";
import "../libraries/Types.sol";

contract L1CrossLayerMessageWitness is IL1CrossLayerMessageWitness {
    using Types for Types.Block;
    using MerkleMountainRange for CompactMerkleTree;
    IAddressResolver addressResolver;

    CompactMerkleTree compactMerkleTree;
    mapping(bytes32 => bool) public successRelayedMessages;
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
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, _sender, _message, _messageIndex);
        bytes32 _hash = keccak256(_crossLayerCalldata);

        require(addressResolver.rollupStateChain().verifyStateInfo(_stateInfo), "wrong state info");
        require(addressResolver.rollupStateChain().isStateConfirmed(_stateInfo), "state not confirmed yet");
        require(_block.hash() == _stateInfo.blockHash, "wrong block provide");
        MerkleMountainRange.verifyLeafHashInclusion(_hash, _messageIndex, _proof, _block.mmrRoot, _block.mmrSize);
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

    function sendMessage(
        address _target,
        bytes calldata _message,
        uint64 _gasLimit
    ) public {
        //should buy gas
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(
            _target,
            msg.sender,
            _message,
            compactMerkleTree.treeSize
        );
        compactMerkleTree.appendLeafHash(keccak256(_crossLayerCalldata));
        addressResolver.rollupInputChain().enqueue(
            address(addressResolver.l2CrossLayerMessageWitness()),
            _gasLimit,
            _crossLayerCalldata
        );
        //todo: emit
        //do not need event,already emit in rollupInputChain
    }

    function replayMessage(
        bytes memory _crossLayerCalldata,
        uint64 _queueIndex,
        uint64 _oldGasLimit,
        uint64 _newGasLimit
    ) public {
        (bytes32 _infoHash, ) = addressResolver.rollupInputChain().getQueueTxInfo(_queueIndex);
        //same as rollupInputChain
        bytes32 _txHash = keccak256(
            abi.encode(
                address(this),
                address(addressResolver.l2CrossLayerMessageWitness()),
                _oldGasLimit,
                _crossLayerCalldata
            )
        );
        require(_txHash == _infoHash, "Provided message has not been enqueued");
        addressResolver.rollupInputChain().enqueue(
            address(addressResolver.l2CrossLayerMessageWitness()),
            _newGasLimit,
            _crossLayerCalldata
        );
    }

    function mmrRoot() public view returns (bytes32) {
        return compactMerkleTree.rootHash;
    }

    function totalSize() public view returns (uint64) {
        return compactMerkleTree.treeSize;
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
