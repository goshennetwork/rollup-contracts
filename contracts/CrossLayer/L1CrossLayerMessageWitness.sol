// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;
import "../libraries/MerkleMountainRange.sol";
import "../interfaces/IL1CrossLayerMessageWitness.sol";
import "../interfaces/IAddressResolver.sol";
import "../libraries/Types.sol";

contract L1CrossLayerMessageWitness is IL1CrossLayerMessageWitness {
    using Types for Types.Block;
    IAddressResolver addressResolver;

    MerkleMountainRange.RootNode[] trees;
    bytes32 public override mmrRoot;
    uint64 public override totalSize;
    mapping(bytes32 => bool) public successRelayedMessages;
    address private crossDomainMsgSender;

    constructor(address _addressResolver) {
        addressResolver = IAddressResolver(_addressResolver);
    }

    function l2Sender() public view returns (address) {
        require(crossDomainMsgSender != address(0), "crossDomainMsgSender not set yet");
        return crossDomainMsgSender;
    }

    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        Types.Block memory _block,
        Types.StateInfo memory _stateInfo,
        Types.MMRInclusionProof memory _proof
    ) public {
        require(crossDomainMsgSender == address(0), "reentrancy");
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, _sender, _message, _messageIndex);
        bytes32 _hash = keccak256(_crossLayerCalldata);

        require(addressResolver.scc().verifyStateInfo(_stateInfo), "wrong state info");
        require(addressResolver.scc().isStateConfirmed(_stateInfo), "state not confirmed yet");
        require(_block.hash() == _stateInfo.blockHash, "wrong block provide");
        require(MerkleMountainRange.verifyTrees(_block.mmrRoot, _block.mmrSize, _proof.trees), "wrong mmr proof");
        require(
            _hash == _proof.leaf &&
                MerkleMountainRange.verifyLeafTree(
                    _proof.trees,
                    _block.mmrSize,
                    _proof.leafIndex,
                    _proof.siblings,
                    _proof.leaf
                ),
            "wrong inclustion proof"
        );
        require(successRelayedMessages[_hash] == false, "provided message already been relayed");
        crossDomainMsgSender = _sender;
        (bool success, ) = _target.call(_message);
        crossDomainMsgSender = address(0);
        if (success) {
            successRelayedMessages[_hash] = true;
            emit RelayedMessage(_hash);
        } else {
            emit FailedRelayedMessage(_hash);
        }
    }

    function sendMessage(
        address _target,
        bytes calldata _message,
        uint64 _gasLimit
    ) public {
        //should buy gas
        bytes memory _crossLayerCalldata = _encodeCrossLayerCallData(_target, msg.sender, _message, totalSize);
        MerkleMountainRange.appendLeafHash(trees, keccak256(_crossLayerCalldata));
        mmrRoot = MerkleMountainRange.genMMRRoot(trees);
        addressResolver.ctc().enqueue(
            address(addressResolver.l2CrossDomainMessageWitness()),
            _gasLimit,
            _crossLayerCalldata
        );
        totalSize++;
        //do not need event,already emit in ctc
    }

    function replayMessage(
        bytes memory _crossLayerCalldata,
        uint64 _queueIndex,
        uint64 _oldGasLimit,
        uint64 _newGasLimit
    ) public {
        Types.QueueElement memory _element = addressResolver.ctc().getQueueElement(_queueIndex);
        //same as ctc
        bytes32 _txHash = keccak256(
            abi.encode(
                address(this),
                address(addressResolver.l2CrossDomainMessageWitness()),
                _oldGasLimit,
                _crossLayerCalldata
            )
        );
        require(_txHash == _element.transactionHash, "Provided message has not been enqueued");
        addressResolver.ctc().enqueue(
            address(addressResolver.l2CrossDomainMessageWitness()),
            _newGasLimit,
            _crossLayerCalldata
        );
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
