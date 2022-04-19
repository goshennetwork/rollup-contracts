// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./MerkleMountainRange.sol";
import "./console.sol";
import { CompactMerkleTree, MerkleMountainRange } from "./MerkleMountainRange.sol";

contract MMRTest {
    using MerkleMountainRange for CompactMerkleTree;
    CompactMerkleTree _trees;

    function getTreeSize() public returns (uint64) {
        return _trees.treeSize;
    }

    function getRootHash() public returns (bytes32) {
        return _trees.rootHash;
    }

    function append(bytes32 _leafHash) public {
        _trees.appendLeafHash(_leafHash);
    }

    function verifyProof(
        bytes32 _leafHash,
        uint64 _leafIndex,
        bytes32[] memory _proof,
        bytes32 _rootHash,
        uint64 _treeSize
    ) public {
        MerkleMountainRange.verifyLeafHashInclusion(_leafHash, _leafIndex, _proof, _rootHash, _treeSize);
    }

    function testAppend() public {
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 1, "0");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 1, "1");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 2, "2");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 1, "3");

        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 2, "4");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 2, "5");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 3, "6");
    }

    function append7Leaf() internal {
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));

        MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(1)));

        MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(2)));

        MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(3)));

        MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(4)));

        MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(5)));

        MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(6)));
    }

    //todo: verify
    function testVerify() public {
        bytes32[] memory _proof;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        MerkleMountainRange.verifyLeafHashInclusion(bytes32(0), uint64(0), _proof, _trees.rootHash, _trees.treeSize);
    }
}
