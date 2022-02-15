// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./BytesSlice.sol";

/**
 * @title RLPReader
 * @dev Adapted from "RLPReader" by Hamdi Allam (hamdi.allam97@gmail.com) and Optimism
 */
library RLPReader {
    uint256 internal constant MAX_LIST_LENGTH = 32;
    enum RLPItemType {
        DATA_ITEM,
        LIST_ITEM
    }

    /**
     * Reads an RLP list value into a list of RLP items.
     * @param rawRlp RLP list value.
     * @return Decoded RLP list items.
     */
    function readList(Slice memory rawRlp) internal pure returns (Slice[] memory) {
        (uint256 listOffset, uint256 listLength, RLPItemType itemType) = decodeKind(rawRlp);

        require(itemType == RLPItemType.LIST_ITEM, "Invalid RLP list value.");

        // Solidity in-memory arrays can't be increased in size, but *can* be decreased in size by
        // writing to the length. Since we can't know the number of RLP items without looping over
        // the entire input, we'd have to loop twice to accurately size this array. It's easier to
        // simply set a reasonable maximum list length and decrease the size before we finish.
        Slice[] memory out = new Slice[](MAX_LIST_LENGTH);

        uint256 itemCount = 0;
        uint256 offset = listOffset;
        require(listOffset + listLength == rawRlp.len, "Provided RLP List not consistent");
        while (offset < rawRlp.len) {
            require(itemCount < MAX_LIST_LENGTH, "Provided RLP list exceeds max list length");
            (uint256 itemOffset, uint256 itemLength, ) = decodeKind(
                Slice({ len: rawRlp.len - offset, ptr: rawRlp.ptr + offset })
            );

            out[itemCount] = Slice({ len: itemLength + itemOffset, ptr: rawRlp.ptr + offset });

            itemCount += 1;
            offset += itemOffset + itemLength;
        }

        // Decrease the array size to match the actual item count.
        assembly {
            mstore(out, itemCount)
        }

        return out;
    }

    /**
     * Reads an RLP list value into a list of RLP items.
     * @param rawRlp RLP list value.
     * @return Decoded RLP list items.
     */
    function readList(bytes memory rawRlp) internal pure returns (Slice[] memory) {
        return readList(BytesSlice.fromBytes(rawRlp));
    }

    function readBytes(Slice memory rawRlp) internal pure returns (bytes memory) {
        (uint256 itemOffset, uint256 itemLength, RLPItemType itemType) = decodeKind(rawRlp);

        require(itemType == RLPItemType.DATA_ITEM, "Invalid RLP bytes value.");

        return BytesSlice.toBytes(Slice({ ptr: rawRlp.ptr + itemOffset, len: itemLength }));
    }

    function readBytes(bytes memory rawRlp) internal pure returns (bytes memory) {
        return readBytes(BytesSlice.fromBytes(rawRlp));
    }

    function readString(Slice memory rawRlp) internal pure returns (string memory) {
        return string(readBytes(rawRlp));
    }

    function readString(bytes memory rawRlp) internal pure returns (string memory) {
        return readString(BytesSlice.fromBytes(rawRlp));
    }

    function readUint256(Slice memory rawRlp) internal pure returns (uint256) {
        require(rawRlp.len <= 33, "Invalid RLP bytes32 value.");

        (uint256 itemOffset, uint256 itemLength, RLPItemType itemType) = decodeKind(rawRlp);

        require(itemType == RLPItemType.DATA_ITEM, "Invalid RLP bytes32 value.");

        uint256 ptr = rawRlp.ptr + itemOffset;
        uint256 out;
        assembly {
            out := mload(ptr)

            // Shift the bytes over to match the item size.
            if lt(itemLength, 32) {
                out := div(out, exp(256, sub(32, itemLength)))
            }
        }

        return out;
    }

    function readUint256(bytes memory rawRlp) internal pure returns (uint256) {
        return readUint256(BytesSlice.fromBytes(rawRlp));
    }

    function readBytes32(bytes memory rawRlp) internal pure returns (bytes32) {
        return readBytes32(BytesSlice.fromBytes(rawRlp));
    }

    function readBytes32(Slice memory rawRlp) internal pure returns (bytes32) {
        return bytes32(readUint256(rawRlp));
    }

    function readBool(Slice memory rawRlp) internal pure returns (bool) {
        require(rawRlp.len == 1, "Invalid RLP boolean value.");

        uint256 ptr = rawRlp.ptr;
        uint256 out;
        assembly {
            out := byte(0, mload(ptr))
        }

        require(out == 0x80 || out == 1, "Invalid RLP boolean value");
        return out == 1;
    }

    function readBool(bytes memory rawRlp) internal pure returns (bool) {
        return readBool(BytesSlice.fromBytes(rawRlp));
    }

    /**
     * Like readAddress, but if the length is 1, and value is (false)"80", return empty address
     */
    function readOptionAddress(Slice memory rawRlp) internal pure returns (address) {
        if (rawRlp.len == 1 && readBool(rawRlp) == false) {
            return address(0);
        }
        require(rawRlp.len == 21, "Invalid RLP address value.");

        return address(uint160(readUint256(rawRlp)));
    }

    function readOptionAddress(bytes memory rawRlp) internal pure returns (address) {
        return readOptionAddress(BytesSlice.fromBytes(rawRlp));
    }

    function readAddress(Slice memory rawRlp) internal pure returns (address) {
        require(rawRlp.len == 21, "Invalid RLP address value.");

        return address(uint160(readUint256(rawRlp)));
    }

    function readAddress(bytes memory rawRlp) internal pure returns (address) {
        return readAddress(BytesSlice.fromBytes(rawRlp));
    }

    /**
     * Decodes the length and item type of an RLP item.
     * @param rawRlp RLP item to decode.
     * @return Start offset of the encoded data.
     * @return Length of the encoded data.
     * @return RLP item type (LIST_ITEM or DATA_ITEM).
     */
    function decodeKind(Slice memory rawRlp)
        internal
        pure
        returns (
            uint256,
            uint256,
            RLPItemType
        )
    {
        require(rawRlp.len > 0, "RLP item cannot be null.");

        uint256 ptr = rawRlp.ptr;
        uint256 prefix;
        assembly {
            prefix := byte(0, mload(ptr))
        }

        if (prefix <= 0x7f) {
            // Single byte.
            return (0, 1, RLPItemType.DATA_ITEM);
        } else if (prefix <= 0xb7) {
            // Short string.
            uint256 strLen = prefix - 0x80;
            require(rawRlp.len > strLen, "Invalid RLP short string.");
            return (1, strLen, RLPItemType.DATA_ITEM);
        } else if (prefix <= 0xbf) {
            // Long string.
            uint256 lenOfStrLen = prefix - 0xb7;
            require(rawRlp.len > lenOfStrLen, "Invalid RLP long string length.");
            uint256 strLen;
            assembly {
                // Pick out the string length. note: rlp的标准要求整数采用大端编码，且必须移除前缀0,因此这里没有做这个检查.
                strLen := div(mload(add(ptr, 1)), exp(256, sub(32, lenOfStrLen)))
            }
            require(rawRlp.len > lenOfStrLen + strLen, "Invalid RLP long string.");

            return (1 + lenOfStrLen, strLen, RLPItemType.DATA_ITEM);
        } else if (prefix <= 0xf7) {
            // Short list.
            uint256 listLen = prefix - 0xc0;
            require(rawRlp.len > listLen, "Invalid RLP short list.");

            return (1, listLen, RLPItemType.LIST_ITEM);
        } else {
            // Long list.
            uint256 lenOfListLen = prefix - 0xf7;
            require(rawRlp.len > lenOfListLen, "Invalid RLP long list length.");

            uint256 listLen;
            assembly {
                // Pick out the list length. note: 同上.
                listLen := div(mload(add(ptr, 1)), exp(256, sub(32, lenOfListLen)))
            }

            require(rawRlp.len > lenOfListLen + listLen, "Invalid RLP long list.");

            return (1 + lenOfListLen, listLen, RLPItemType.LIST_ITEM);
        }
    }
}
