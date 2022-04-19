// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./MerkleMountainRange.sol";
import "./console.sol";
import "./MerkleMountainRange.sol";

contract MMRTest {
    using MerkleMountainRange for MerkleMountainRange.CompactMerkleTree;
    MerkleMountainRange.CompactMerkleTree _trees;

    function getTreeSize() public returns (uint64) {
        return _trees.treeSize;
    }

    function getRootHash() public returns (bytes32) {
        return _trees.merkleRoot();
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
        uint256 _total;
        _total++;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.hashes.length); ///1
        _total++;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.hashes.length); ///1
        _total++;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.hashes.length); ///2
        _total++;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.hashes.length); //1
        _total++;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.hashes.length); ///2
        _total++;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.hashes.length); ///2
        _total++;
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.hashes.length); ///3\
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

    //    function testGenRoot() public {
    //        MerkleMountainRange.RootNode[] memory _trees = append7Leaf();
    //        bytes32 _root = MerkleMountainRange.genMMRRoot(_trees);
    //        require(MerkleMountainRange.verifyTrees(_root, _trees), "verify");
    //    }
    //
    //    function testverifyLeafTree() public {
    //        MerkleMountainRange.RootNode[] memory _trees = append7Leaf();
    //
    //        ///index 2
    //        uint64 _indexInAll = 2;
    //        uint64 indexInTrees;
    //        uint64 preTotalNum;
    //        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, _indexInAll);
    //
    //        bytes32[] memory _siblings = new bytes32[](2);
    //        _siblings[0] = bytes32(uint256(3));
    //        _siblings[1] = keccak256(abi.encodePacked(bytes32(0), bytes32(uint256(1))));
    //        bool _correct = MerkleMountainRange.verifyLeafTree(
    //            _siblings,
    //            _trees[indexInTrees].hash,
    //            _indexInAll - preTotalNum,
    //            uint64(1) << _trees[indexInTrees].level,
    //            bytes32(uint256(2))
    //        );
    //        require(_correct, "index 2");
    //
    //        //index 6
    //        _indexInAll = 6;
    //        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, _indexInAll);
    //        assembly {
    //            //do not need siblings
    //            mstore(_siblings, 0)
    //        }
    //        _correct = MerkleMountainRange.verifyLeafTree(
    //            _siblings,
    //            _trees[indexInTrees].hash,
    //            _indexInAll - preTotalNum,
    //            uint64(1) << _trees[indexInTrees].level,
    //            bytes32(uint256(6))
    //        );
    //        require(_correct, "index 6");
    //    }
}
