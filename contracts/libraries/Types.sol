// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./MerkleMountainRange.sol";

library Types {
    ///block info,need fix
    struct Block {
        ///....more
        bytes32 mmrRoot;
        uint64 mmrSize;
    }

    function encode(Block memory _block) internal pure returns (bytes memory) {
        return abi.encodePacked(_block.mmrRoot, _block.mmrSize);
    }

    function hash(Block memory _block) internal pure returns (bytes32) {
        return keccak256(encode(_block));
    }

    struct StateInfo {
        bytes32 blockHash;
        uint64 index;
        uint64 timestamp;
        address proposer;
    }

    using Types for StateInfo;

    function encode(StateInfo memory _stateInfo) internal pure returns (bytes memory) {
        return abi.encodePacked(_stateInfo.blockHash, _stateInfo.index, _stateInfo.timestamp, _stateInfo.proposer);
    }

    ///hash state info
    function hash(StateInfo memory _stateInfo) internal pure returns (bytes32) {
        return keccak256(_stateInfo.encode());
    }
}
