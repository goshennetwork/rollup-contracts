// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./MerkleMountainRange.sol";
import "./console.sol";

contract MMRTest {
    function testAppend() public {
        MerkleMountainRange.RootNode[] memory _trees;
        uint256 _total;
        _total++;
        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///1
        _total++;
        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///1
        _total++;
        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///2
        _total++;
        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///1
        _total++;
        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///2
        _total++;
        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///2
        _total++;
        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///3\

        uint64 indexInTrees;
        uint64 preTotalNum;
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, 0);
        require(indexInTrees == 0 && preTotalNum == 0, "0");
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, 1);
        require(indexInTrees == 0 && preTotalNum == 0, "1");
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, 2);
        require(indexInTrees == 0 && preTotalNum == 0, "2");
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, 3);
        require(indexInTrees == 0 && preTotalNum == 0, "3");
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, 4);
        require(indexInTrees == 1 && preTotalNum == 4, "4");
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, 5);
        require(indexInTrees == 1 && preTotalNum == 4, "4");
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, 6);
        require(indexInTrees == 2 && preTotalNum == 6, "6");
    }

    function append7Leaf() internal returns (MerkleMountainRange.RootNode[] memory) {
        MerkleMountainRange.RootNode[] memory _trees;

        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        console.logUint(_trees.length); ///1

        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(1)));
        console.logUint(_trees.length); ///1

        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(2)));
        console.logUint(_trees.length); ///2

        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(3)));
        console.logUint(_trees.length); ///1

        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(4)));
        console.logUint(_trees.length); ///2

        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(5)));
        console.logUint(_trees.length); ///2

        _trees = MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(6)));
        console.logUint(_trees.length); ///3\
        return _trees;
    }

    function testGenRoot() public {
        MerkleMountainRange.RootNode[] memory _trees = append7Leaf();
        bytes32 _root = MerkleMountainRange.genMMRRoot(_trees);
        require(MerkleMountainRange.verifyTrees(_root, _trees), "verify");
    }

    function testverifyLeafTree() public {
        MerkleMountainRange.RootNode[] memory _trees = append7Leaf();

        ///index 2
        uint64 _indexInAll = 2;
        uint64 indexInTrees;
        uint64 preTotalNum;
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, _indexInAll);

        bytes32[] memory _siblings = new bytes32[](2);
        _siblings[0] = bytes32(uint256(3));
        _siblings[1] = keccak256(abi.encodePacked(bytes32(0), bytes32(uint256(1))));
        bool _correct = MerkleMountainRange.verifyLeafTree(
            _siblings,
            _trees[indexInTrees].hash,
            _indexInAll - preTotalNum,
            uint64(1) << _trees[indexInTrees].level,
            bytes32(uint256(2))
        );
        require(_correct, "index 2");

        //index 6
        _indexInAll = 6;
        (indexInTrees, preTotalNum) = MerkleMountainRange.chooseTreeNode(_trees, _indexInAll);
        assembly {
            //do not need siblings
            mstore(_siblings, 0)
        }
        _correct = MerkleMountainRange.verifyLeafTree(
            _siblings,
            _trees[indexInTrees].hash,
            _indexInAll - preTotalNum,
            uint64(1) << _trees[indexInTrees].level,
            bytes32(uint256(6))
        );
        require(_correct, "index 6");
    }
}
