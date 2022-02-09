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
    function memset(uint256 dest, bytes32 src, uint256 len) internal pure {
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
}
