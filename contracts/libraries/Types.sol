// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./MerkleMountainRange.sol";
import "./RLPReader.sol";
import "./BytesSlice.sol";

library Types {
    using RLPReader for Slice;
    ///block info,need fix
    struct Block {
        ///....more
        bytes32 mmrRoot;
        uint64 mmrSize;
    }

    // mmr info is stored in header's seal field
    function decodeMMRFromRlpHeader(bytes memory header) internal pure returns (bytes32 mmrRoot, uint64 mmrSize) {
        Slice memory rawRlp = BytesSlice.fromBytes(header);
        Slice[] memory fields = rawRlp.readList();
        require(fields.length >= 15);
        bytes memory mmrRootField = fields[13].readBytes();
        require(mmrRootField.length == 32, "ill mmrRoot length");
        mmrRoot = bytes32(mmrRootField);
        bytes memory mmrSizeField = fields[14].readBytes();
        require(mmrSizeField.length <= 8, "ill mmrSize length");
        mmrSize = uint64(bytes8(mmrSizeField));
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
