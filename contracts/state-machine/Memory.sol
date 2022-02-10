// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/MerkleTrie.sol";
import "../libraries/BytesSlice.sol";

library Memory {
    function writeMemory(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr,
        uint32 value
    ) public returns (bytes32) {
        require(ptr & 3 == 0, "non-aligned mem ptr");
        return MerkleTrie.update(hashdb, uint32ToBytes(ptr), uint32ToBytes(value), root);
    }

    function writeMemoryBytes32(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr,
        bytes32 val
    ) public returns (bytes32) {
        for (uint32 i = 0; i < 32; i += 4) {
            root = writeMemory(hashdb, root, ptr + i, uint32(bytes4(val)));
            val = bytes32(uint256(val) << 32);
        }
        return root;
    }

    function readMemory(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 addr
    ) public view returns (uint32) {
        require(addr & 3 == 0, "non-aligned mem ptr");
        (bool exists, bytes memory value) = MerkleTrie.get(hashdb, uint32ToBytes(addr), root);
        return exists ? bytesToUint32(value) : 0;
    }

    function readMemoryBytes32(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 addr
    ) public view returns (bytes32) {
        uint256 ret = 0;
        for (uint32 i = 0; i < 32; i += 4) {
            ret <<= 32;
            ret |= uint256(readMemory(hashdb, root, addr + i));
        }
        return bytes32(ret);
    }

    function uint32ToBytes(uint32 data) internal pure returns (bytes memory) {
        return bytes.concat(bytes4(data));
    }

    function bytesToUint32(bytes memory dat) internal pure returns (uint32) {
        require(dat.length == 4, "wrong length value");
        bytes32 data = BytesSlice.toBytes32(dat);
        return uint32(bytes4(data));
    }
}
