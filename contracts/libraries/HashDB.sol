// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./BytesSlice.sol";

library HashDB {
    struct Preimage {
        uint32 length;
        mapping(uint32 => bytes) partials;
    }
    //length stored in -1 index
    using HashDB for mapping(bytes32 => Preimage);
    using BytesSlice for Slice;
    uint256 constant PartialSize = 1024;
    bytes32 constant EMPTY_HASH = keccak256("");

    function preimageAtIndex(
        mapping(bytes32 => Preimage) storage partialImage,
        bytes32 _hash,
        uint32 _index
    ) internal view returns (bytes memory _ret, uint32 _len) {
        if (_hash == EMPTY_HASH) {
            return ("", 0);
        }
        Preimage storage preimage = partialImage[_hash];
        bytes memory _ret = preimage.partials[_index];
        require(_ret.length > 0, "not exist");
        return (_ret, uint32(_ret.length));
    }

    function insertPartialImage(
        mapping(bytes32 => Preimage) storage partialImage,
        bytes memory _node,
        uint32 _index
    ) internal {
        uint256 _length = _node.length;
        require(_index * PartialSize < _length, "wrong index");
        bytes32 _hash = keccak256(_node);
        Preimage storage preimage = partialImage[_hash];
        if (preimage.length == 0) {
            //not set yet
            require(_length > 0, "empty image exist");
            preimage.length = uint32(_length);
        }
        uint256 _left = (1 + _index) * PartialSize <= _length ? PartialSize : _length % PartialSize;
        preimage.partials[_index] = BytesSlice.slice(_node, _index * PartialSize, _left).toBytes();
    }

    function insertPreimage(mapping(bytes32 => Preimage) storage partialImage, bytes memory _node) internal {
        for (uint256 i = 0; i < (_node.length + PartialSize - 1) / PartialSize; i++) {
            partialImage.insertPartialImage(_node, uint32(i));
        }
    }

    function preimage(mapping(bytes32 => Preimage) storage partialImage, bytes32 _hash)
        internal
        view
        returns (bytes memory)
    {
        if (_hash == EMPTY_HASH) {
            return "";
        }
        Preimage storage preimage = partialImage[_hash];
        require(preimage.length > 0, "no node");
        uint32 _num = uint32((preimage.length + PartialSize - 1) / PartialSize);
        bytes[] memory _partials = new bytes[](_num);
        uint32 _len;
        for (uint32 i = 0; i < _num; i++) {
            //do not use function, for saving gas
            _partials[i] = preimage.partials[i];
        }
        bytes memory _data = BytesSlice.concat("", _partials);
        require(_data.length == preimage.length, "no complete");
        return _data;
    }
}
