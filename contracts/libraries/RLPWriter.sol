// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./BytesSlice.sol";

// TODO: these code have performance bottleneck

library RLPWriter {
    function writeBytes(bytes memory _in) internal pure returns (bytes memory) {
        bytes memory encoded;

        if (_in.length == 0) {
            encoded = new bytes(1);
            encoded[0] = 0x80;
        } else if (_in.length == 1 && uint8(_in[0]) < 128) {
            encoded = _in;
        } else {
            encoded = abi.encodePacked(_writeLength(_in.length, 128), _in);
        }

        return encoded;
    }

    function writeList(bytes[] memory _list) internal pure returns (bytes memory) {
        uint256 len = 0;
        for (uint256 i = 0; i < _list.length; i++) {
            len += _list[i].length;
        }
        bytes memory prefix = _writeLength(len, 192);

        return BytesSlice.concat(prefix, _list, len);
    }

    function writeString(string memory _in) internal pure returns (bytes memory) {
        return writeBytes(bytes(_in));
    }

    /**
     * RLP encodes an address.
     * @param _in The address to encode.
     * @return The RLP encoded address in bytes.
     */
    function writeAddress(address _in) internal pure returns (bytes memory) {
        return abi.encodePacked(bytes1(uint8(20 + 128)), _in);
    }

    /**
     * RLP encodes a bytes32 value.
     * @param _in The bytes32 to encode.
     * @return _out The RLP encoded bytes32 in bytes.
     */
    function writeBytes32(bytes32 _in) internal pure returns (bytes memory _out) {
        return abi.encodePacked(bytes1(uint8(32 + 128)), _in);
    }

    /**
     * RLP encodes a uint.
     * @param _in The uint256 to encode.
     * @return The RLP encoded uint256 in bytes.
     */
    function writeUint(uint256 _in) internal pure returns (bytes memory) {
        return writeBytes(_toBinary(_in));
    }

    /**
     * RLP encodes a bool.
     * @param _in The bool to encode.
     * @return The RLP encoded bool in bytes.
     */
    function writeBool(bool _in) internal pure returns (bytes memory) {
        bytes memory encoded = new bytes(1);
        encoded[0] = (_in ? bytes1(0x01) : bytes1(0x80));
        return encoded;
    }

    /*********************
     * Private Functions *
     *********************/

    /**
     * Encode the first byte, followed by the `len` in binary form if `length` is more than 55.
     * @param _len The length of the string or the payload.
     * @param _offset 128 if item is string, 192 if item is list.
     * @return RLP encoded bytes.
     */
    function _writeLength(uint256 _len, uint256 _offset) private pure returns (bytes memory) {
        bytes memory encoded;

        if (_len < 56) {
            encoded = new bytes(1);
            encoded[0] = bytes1(uint8(_len + _offset));
        } else {
            uint256 lenLen;
            uint256 i = 1;
            while (_len / i != 0) {
                lenLen++;
                i *= 256;
            }

            encoded = new bytes(lenLen + 1);
            encoded[0] = bytes1(uint8(lenLen) + uint8(_offset) + 55);
            for (i = 1; i <= lenLen; i++) {
                encoded[i] = bytes1(uint8((_len / (256**(lenLen - i))) % 256));
            }
        }

        return encoded;
    }

    /**
     * Encode integer in big endian binary form with no leading zeroes.
     * @param _x The integer to encode.
     * @return RLP encoded bytes.
     */
    function _toBinary(uint256 _x) internal pure returns (bytes memory) {
        uint256 len = lenBytes(_x);
        bytes memory res = new bytes(32);
        _x = _x << ((32 - len) * 8);
        assembly {
            mstore(add(res, 32), _x)
            mstore(res, len)
        }
        return res;
    }

    // returns the minimum number of bytes required to represent x; the result is 0 for x == 0.
    function lenBytes(uint256 _x) private pure returns (uint256) {
        uint256 n = 0;
        if (_x >= 1 << 128) {
            n += 16;
            _x >>= 128;
        }
        if (_x >= 1 << 64) {
            n += 8;
            _x >>= 64;
        }
        if (_x >= 1 << 32) {
            n += 4;
            _x >>= 32;
        }
        if (_x >= 1 << 16) {
            n += 2;
            _x >>= 16;
        }
        if (_x >= 1 << 8) {
            n += 1;
            _x >>= 8;
        }
        if (_x > 0) {
            n += 1;
        }
        return n;
    }
}
