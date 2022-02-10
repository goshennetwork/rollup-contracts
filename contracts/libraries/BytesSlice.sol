// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

struct Slice {
    uint256 len;
    uint256 ptr;
}

library BytesSlice {
    function memcpy(
        uint256 dest,
        uint256 src,
        uint256 len
    ) internal pure {
        // Copy word-length chunks while possible
        for (; len >= 32; len -= 32) {
            assembly {
                mstore(dest, mload(src))
            }
            dest += 32;
            src += 32;
        }

        // Copy remaining bytes
        if (len == 0) {
            //have no remaining
            return;
        }
        uint256 mask = 256**(32 - len) - 1;
        assembly {
            let srcpart := and(mload(src), not(mask))
            let destpart := and(mload(dest), mask)
            mstore(dest, or(destpart, srcpart))
        }
    }

    /**
     * Set memory at dest, the equivalent golang code is `copy(mem[dest:], src[:len])`
     */
    function memset(
        uint256 dest,
        bytes32 src,
        uint256 len
    ) internal pure {
        if (len == 0) {
            return;
        }
        uint256 mask = 256**(32 - len) - 1;
        uint256 srcval = uint256(src);
        assembly {
            let srcpart := and(srcval, not(mask))
            let destpart := and(mload(dest), mask)
            mstore(dest, or(destpart, srcpart))
        }
    }

    /**
     * @dev Returns a slice containing the entire string.
     * @param self The string to make a slice from.
     * @return A newly allocated slice containing the entire string.
     */
    function fromString(string memory self) internal pure returns (Slice memory) {
        return fromBytes(bytes(self));
    }

    /**
     * @dev Returns a slice containing the entire string.
     * @param self The bytes to make a slice from.
     * @return A newly allocated slice containing the entire string.
     */
    function fromBytes(bytes memory self) internal pure returns (Slice memory) {
        uint256 ptr;
        assembly {
            ptr := add(self, 0x20)
        }
        return Slice(self.length, ptr);
    }

    /*
     * @dev Copies a slice to a new string.
     * @param self The slice to copy.
     * @return A newly allocated string containing the slice's text.
     */
    function toString(Slice memory self) internal pure returns (string memory) {
        return string(toBytes(self));
    }

    /**
     * @dev Copies a slice to a new bytes.
     * @param self The slice to copy.
     * @return A newly allocated bytes containing the slice's text.
     */
    function toBytes(Slice memory self) internal pure returns (bytes memory) {
        bytes memory ret = new bytes(self.len);
        uint256 retptr;
        assembly {
            retptr := add(ret, 32)
        }

        memcpy(retptr, self.ptr, self.len);
        return ret;
    }

    function keccak(Slice memory self) internal pure returns (bytes32) {
        uint256 ptr = self.ptr;
        uint256 len = self.len;
        bytes32 result;
        assembly {
            result := keccak256(ptr, len)
        }
        return result;
    }

    /**
     * Concatenate a variable number of bytes. note: you can also use built-in function `bytes.concat`
     * if the list size is static.
     */
    function concat(bytes[] memory _list) internal pure returns (bytes memory) {
        if (_list.length == 0) {
            return new bytes(0);
        }

        uint256 len;
        uint256 i = 0;
        for (; i < _list.length; i++) {
            len += _list[i].length;
        }

        bytes memory flattened = new bytes(len);
        uint256 flattenedPtr;
        assembly {
            flattenedPtr := add(flattened, 0x20)
        }

        for (i = 0; i < _list.length; i++) {
            bytes memory item = _list[i];

            uint256 listPtr;
            assembly {
                listPtr := add(item, 0x20)
            }

            BytesSlice.memcpy(flattenedPtr, listPtr, item.length);
            flattenedPtr += _list[i].length;
        }

        return flattened;
    }

    function equal(bytes memory left, bytes memory right) internal pure returns (bool) {
        return keccak(fromBytes(left)) == keccak(fromBytes(right));
    }

    function equal(string memory left, string memory right) internal pure returns (bool) {
        return equal(bytes(left), bytes(right));
    }

    function slice(
        bytes memory buff,
        uint256 start,
        uint256 length
    ) internal pure returns (Slice memory) {
        return slice(fromBytes(buff), start, length);
    }

    function slice(
        Slice memory buff,
        uint256 start,
        uint256 length
    ) internal pure returns (Slice memory) {
        require(buff.len >= start + length, "oob");
        return Slice({ ptr: buff.ptr + start, len: length });
    }

    function slice(bytes memory buff, uint256 start) internal pure returns (Slice memory) {
        return slice(buff, start, buff.length - start);
    }

    function toNibbles(bytes memory _bytes) internal pure returns (bytes memory) {
        bytes memory nibbles = new bytes(_bytes.length * 2);

        for (uint256 i = 0; i < _bytes.length; i++) {
            nibbles[i * 2] = _bytes[i] >> 4;
            nibbles[i * 2 + 1] = bytes1(uint8(_bytes[i]) % 16);
        }

        return nibbles;
    }

    function fromNibbles(bytes memory _bytes) internal pure returns (bytes memory) {
        bytes memory ret = new bytes(_bytes.length / 2);

        for (uint256 i = 0; i < ret.length; i++) {
            ret[i] = (_bytes[i * 2] << 4) | (_bytes[i * 2 + 1]);
        }

        return ret;
    }

    function toBytes32PadLeft(bytes memory _bytes) internal pure returns (bytes32) {
        bytes32 ret;
        uint256 len = _bytes.length <= 32 ? _bytes.length : 32;
        assembly {
            ret := shr(mul(sub(32, len), 8), mload(add(_bytes, 32)))
        }
        return ret;
    }

    function toBytes32(bytes memory _bytes) internal pure returns (bytes32) {
        bytes32 ret;
        assembly {
            ret := mload(add(_bytes, 32))
        }
        if (_bytes.length < 32) {
            uint256 mask ;
            unchecked {
                mask = 256**(32 - _bytes.length) - 1;
            }
            assembly {
                ret := and(ret, not(mask))
            }
        }
        return ret;
    }

    function toUint256(bytes memory _bytes) internal pure returns (uint256) {
        return uint256(toBytes32(_bytes));
    }

    function toUint24(bytes memory _bytes, uint256 _start) internal pure returns (uint24) {
        require(_bytes.length >= _start + 3, "oob");
        uint24 tempUint;

        assembly {
            tempUint := mload(add(add(_bytes, 0x3), _start))
        }

        return tempUint;
    }

    function toUint8(bytes memory _bytes, uint256 _start) internal pure returns (uint8) {
        require(_bytes.length >= _start + 1, "oob");
        uint8 tempUint;

        assembly {
            tempUint := mload(add(add(_bytes, 0x1), _start))
        }

        return tempUint;
    }

    function toAddress(bytes memory _bytes, uint256 _start) internal pure returns (address) {
        require(_bytes.length >= _start + 20, "oob");
        address tempAddress;

        assembly {
            tempAddress := div(mload(add(add(_bytes, 0x20), _start)), 0x1000000000000000000000000)
        }

        return tempAddress;
    }

    function genRevertHex(bytes memory _reason) internal pure returns (bytes memory) {
        bytes memory reason = toNibbles(_reason);
        for (uint256 i = 0; i < reason.length; i++) {
            if (reason[i] < bytes1(uint8(10))) {
                reason[i] = bytes1(uint8(reason[i]) + uint8(0x30));
            } else {
                reason[i] = bytes1(uint8(reason[i]) + uint8(0x61 - 10));
            }
        }

        // func id of Error(string)
        return abi.encodeWithSelector(0x08c379a0, string(reason));
    }
}
