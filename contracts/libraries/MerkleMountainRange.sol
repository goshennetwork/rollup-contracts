// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// CompactMerkleTree calculate merkle tree with compact hash store in HashStore
struct CompactMerkleTree {
    bytes32 rootHash;
    bytes32[] hashes;
    uint64 treeSize;
}

library MerkleMountainRange {
    function appendLeafHash(CompactMerkleTree storage tree, bytes32 leaf) internal returns (bytes32) {
        bytes32[] storage hashes = tree.hashes;
        uint64 size = uint64(hashes.length);
        for (uint256 s = tree.treeSize; s % 2 == 1; s = s >> 1) {
            leaf = keccak256(abi.encodePacked(hashes[size - 1], leaf));
            size -= 1;
        }
        tree.treeSize += 1;
        // resize hashes
        assembly {
            sstore(hashes.slot, size)
        }
        hashes.push(leaf);
        int64 _lenHashes = int64(size + 1);
        bytes32 _accum = leaf;
        for (int64 i = _lenHashes - 2; i >= 0; i--) {
            _accum = keccak256(abi.encodePacked(tree.hashes[uint64(i)], _accum));
        }
        tree.rootHash = _accum;
        return _accum;
    }

    function verifyLeafHashInclusion(
        bytes32 _leafHash,
        uint64 _leafIndex,
        bytes32[] memory _proof,
        bytes32 _rootHash,
        uint64 _treeSize
    ) internal pure {
        require(_leafIndex < _treeSize, "leaf index out of bounds");
        require(
            calculateRootHashFromAuditPath(_leafHash, _leafIndex, _proof, _treeSize) == _rootHash, "mmr root differ"
        );
    }

    function calculateRootHashFromAuditPath(
        bytes32 _leafHash,
        uint64 _leafIndex,
        bytes32[] memory _auditPath,
        uint64 _treeSize
    ) internal pure returns (bytes32) {
        bytes32 _calculatedHash = _leafHash;
        uint64 _pos = 0;
        uint64 _pathLen = uint64(_auditPath.length);
        for (uint64 _lastNode = _treeSize - 1; _lastNode > 0; _lastNode >>= 1) {
            require(_pos < _pathLen, "proof too short");
            if (_leafIndex % 2 == 1) {
                _calculatedHash = keccak256(abi.encodePacked(_auditPath[_pos], _calculatedHash));
                _pos++;
            } else if (_leafIndex < _lastNode) {
                _calculatedHash = keccak256(abi.encodePacked(_calculatedHash, _auditPath[_pos]));
                _pos++;
            }

            _leafIndex >>= 1;
        }
        require(_pos >= _pathLen, "proof too long");
        return _calculatedHash;
    }
}
