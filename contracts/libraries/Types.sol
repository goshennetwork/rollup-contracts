// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

library Types {
    struct QueueElement {
        bytes32 transactionHash;
        uint64 timestamp;
    }

    function encode(QueueElement memory _element) internal pure returns (bytes memory) {
        return abi.encodePacked(_element.transactionHash, _element.timestamp);
    }

    function hash(QueueElement memory _element) internal pure returns (bytes32) {
        return keccak256(encode(_element));
    }

    struct StateInfo {
        bytes32 blockHash;
        uint64 index;
        uint64 timestamp;
        address proposer;
    }

    using Types for StateInfo;

    function encode(StateInfo memory _stateInfo) internal pure returns (bytes memory) {
        return abi.encodePacked(_stateInfo.blockHash, _stateInfo.timestamp, _stateInfo.proposer);
    }

    ///hash state info
    function hash(StateInfo memory _stateInfo) internal pure returns (bytes32) {
        return keccak256(_stateInfo.encode());
    }
}
