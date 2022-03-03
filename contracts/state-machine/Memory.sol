// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/MerkleTrie.sol";
import "../libraries/BytesSlice.sol";
import "../libraries/BytesEndian.sol";

library Memory {
    function writeMemory(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr,
        uint32 value
    ) internal returns (bytes32) {
        return writeMemoryBytes4(hashdb, root, ptr, BytesEndian.uint32ToLEBytes(value));
    }

    function writeMemoryBytes4(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr,
        bytes4 value
    ) internal returns (bytes32) {
        require(ptr & 3 == 0, "write non-aligned mem ptr");
        return MerkleTrie.update(hashdb, uint32ToBytes(ptr), BytesSlice.bytes4ToBytes(value), root);
    }

    function writeMemoryByte(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr,
        bytes1 value
    ) internal returns (bytes32) {
        uint32 offset = (ptr & 3);
        ptr = ptr - offset;
        bytes4 data = readMemoryBytes4(hashdb, root, ptr);
        uint32 shift = 8 * offset;
        bytes4 mask = ~(bytes4(hex"ff") >> shift);
        data = data & mask;
        data = data | (bytes4(value) >> shift);
        return writeMemoryBytes4(hashdb, root, ptr, data);
    }

    function writeMemoryBytes2(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr,
        bytes2 value
    ) internal returns (bytes32) {
        uint32 offset = (ptr & 3);
        require(offset != 3, " write data cross 4byte boundry");
        ptr = ptr - offset;
        bytes4 data = readMemoryBytes4(hashdb, root, ptr);
        uint32 shift = 8 * offset;
        bytes4 mask = ~(bytes4(hex"ffff") >> shift);
        data = data & mask;
        data = data | (bytes4(value) >> shift);
        return writeMemoryBytes4(hashdb, root, ptr, data);
    }

    function writeMemoryBytes32(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr,
        bytes32 val
    ) internal returns (bytes32) {
        for (uint32 i = 0; i < 32; i += 4) {
            root = writeMemoryBytes4(hashdb, root, ptr + i, bytes4(val));
            val <<= 32;
        }
        return root;
    }

    function readMemoryBytes2(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr
    ) internal view returns (bytes2) {
        uint32 offset = (ptr & 3);
        require(offset != 3, "read data cross 4byte boundry");
        bytes4 data = readMemoryBytes4(hashdb, root, ptr - offset);
        return bytes2(data << (offset * 8));
    }

    function readMemoryByte(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr
    ) internal view returns (bytes1) {
        uint32 offset = (ptr & 3);
        bytes4 data = readMemoryBytes4(hashdb, root, ptr - offset);
        return bytes1(data << (offset * 8));
    }

    function readMemoryBytes4(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr
    ) internal view returns (bytes4) {
        require(ptr & 3 == 0, "non-aligned mem ptr");
        (bool exists, bytes memory value) = MerkleTrie.get(hashdb, uint32ToBytes(ptr), root);
        return exists ? BytesSlice.bytesToBytes4(value) : bytes4(0);
    }

    function readMemory(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr
    ) internal view returns (uint32) {
        bytes4 result = readMemoryBytes4(hashdb, root, ptr);
        return BytesEndian.bytes4ToUint32(result);
    }

    function readMemoryBytes32(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 ptr
    ) internal view returns (bytes32) {
        bytes32 ret = 0;
        ptr += 32;
        for (uint32 i = 0; i < 8; i += 1) {
            ptr -= 4;
            ret >>= 32;
            ret |= bytes32(readMemoryBytes4(hashdb, root, ptr));
        }
        return ret;
    }

    function readString(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        uint32 addr,
        uint32 len
    ) internal view returns (string memory) {
        if (len == 0) {
            //maybe should panic?
            return "";
        }
        bytes memory msg;
        for (uint32 offset = 0; offset < len; offset += 4) {
            bytes4 piece = readMemoryBytes4(hashdb, root, addr + offset);
            msg = abi.encodePacked(msg, piece);
        }
        return string(BytesSlice.toBytes(BytesSlice.slice(msg, 0, len)));
    }

    // note: we use big endian encoding to store memory in trie.
    function uint32ToBytes(uint32 data) internal pure returns (bytes memory) {
        return bytes.concat(bytes4(data));
    }
}
