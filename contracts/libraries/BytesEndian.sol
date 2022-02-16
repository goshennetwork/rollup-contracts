// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

library BytesEndian {
    function revertEndian(uint32 val) internal pure returns (uint32) {
        uint256 v1 = (val >> 0) & 0xff;
        uint256 v2 = (val >> 8) & 0xff;
        uint256 v3 = (val >> 16) & 0xff;
        uint256 v4 = (val >> 24) & 0xff;
        uint256 le = (v1 << 24) + (v2 << 16) + (v3 << 8) + v4;
        return uint32(le);
    }

    function revertEndianUint16(uint16 val) internal pure returns (uint16) {
        uint256 v1 = (val >> 0) & 0xff;
        uint256 v2 = (val >> 8) & 0xff;
        uint256 le = (v1 << 8) + v2;
        return uint16(le);
    }

    function bytes2ToUint16(bytes2 val) internal pure returns (uint16) {
        return revertEndianUint16(uint16(val));
    }

    // little endian bytes to uint32
    function bytes4ToUint32(bytes4 _val) internal pure returns (uint32) {
        return revertEndian(uint32(_val));
    }

    function uint32ToLEBytes(uint32 val) internal pure returns (bytes4) {
        return bytes4(revertEndian(val));
    }
}
