// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library MerkleMountainRange {
    // CompactMerkleTree calculate merkle tree with compact hash store in HashStore
    struct CompactMerkleTree {
        bytes32 root;
        bytes32[] hashes;
        uint64 treeSize;
    }

    function appendLeafHash(CompactMerkleTree storage tree, bytes32 leaf) {
        bytes32[] storage hashes = tree.hashes;
        uint64 size = tree.hashes.length;
        for (uint s = tree.treeSize; s%2 == 1; s = s>>1) {
            leaf = keccak256(abi.encodePacked(hashes[size - 1], leaf));
            size -= 1;
        }
        self.treeSize += 1;
        // resize hashes
        assembly {
            sstore(hashes.slot, size)
        }
        hashes.push(leaf);
        self.rootHash = bytes32(0);
    }

    struct RootNode {
        uint64 level;
        bytes32 hash;
    }

    /**
     * @dev Append a leaf to a list of perfect binary tree
     * @param _trees A list of trees' root, ranging from highest level to lowest level, including root level,root hash
     * @param _hash New leaf(level 0)
     */
    function appendLeafHash(RootNode[] memory _trees, bytes32 _hash) internal pure returns (RootNode[] memory) {
        RootNode[] memory _arrs = new RootNode[](_trees.length + 1);
        for (uint256 i = 0; i < _trees.length; i++) {
            _arrs[i] = _trees[i];
        }
        _arrs[_trees.length] = RootNode({ level: 0, hash: _hash });
        uint256 _len = _arrs.length;
        uint256 i = _len - 1;
        for (; i > 0; i--) {
            if (_arrs[i].level == _arrs[i - 1].level) {
                //merge
                _trees[i - 1].level += 1;
                _trees[i - 1].hash = keccak256(abi.encodePacked(_arrs[i - 1].hash, _arrs[i].hash));
            } else {
                break;
            }
        }
        _len = i + 1;
        assembly {
            //resize
            mstore(_arrs, _len)
        }
        return _arrs;
    }

    /**
     * @dev Verify a list of ranged trees is right///maybe do not need it, locally store a set of trees
     */
    function verifyTrees(bytes32 _root, RootNode[] memory _trees) internal pure returns (bool) {
        return genMMRRoot(_trees) == _root;
    }

    /**
     * @dev Chose a tree from a list of ranged trees by leaf index
     * @param _trees A list of ranged trees
     * @param _index Leaf index
     * @return _treeIndex Tree index in list of range trees, which contains wanted leaf
     * @return _preTotal Previous trees have stored total number of leaves
     * @notice revert if find noting (wanted leaf beyond local record)
     */
    function chooseTreeNode(RootNode[] memory _trees, uint64 _index)
        internal
        pure
        returns (uint64 _treeIndex, uint64 _preTotal)
    {
        uint64 _num;
        for (uint64 i = 0; i < _trees.length; i++) {
            uint64 _pre = _num;
            _num += uint64(1) << _trees[i].level;
            if (_num > _index) {
                return (i, _pre);
            }
        }
        revert("find nothing");
    }

    /**
     * @dev Verify leaf in a perfect binary tree
     * @param _siblings Leaf node's sibling to generate merkle tree root,ranging from low level to highest level(except root)
     * @param _treeRoot Root of this tree
     * @param _index Leaf index in this tree
     * @param _totalLeaf Total leaf in this tree,because this tree is perfect binary tree, _totalLeaf must be 2^n
     * @param _leaf Leaf node
     * @notice revert if:
     * - have no leaf
     * - leaf node's index beyond tree capcity
     * - provided sibling num not equal to (log2(total leaf num in tree))
     */
    function verifyLeafTree(
        bytes32[] memory _siblings,
        bytes32 _treeRoot,
        uint64 _index,
        uint64 _totalLeaf,
        bytes32 _leaf
    ) internal pure returns (bool) {
        require(_totalLeaf > 0, "Total leaves must be greater than zero.");
        require(_index < _totalLeaf, "Index out of bounds.");
        bytes32 _computedRoot = _leaf;
        uint64 _calculatedNum = 1;
        for (uint256 i = 0; i < _siblings.length; i++) {
            if ((_index & 1) == 1) {
                _computedRoot = keccak256(abi.encodePacked(_siblings[i], _computedRoot));
            } else {
                _computedRoot = keccak256(abi.encodePacked(_computedRoot, _siblings[i]));
            }
            _calculatedNum <<= 1;
            _index >>= 1;
        }
        require(_calculatedNum == _totalLeaf, "calculated num should equal to total leaf");
        return _computedRoot == _treeRoot;
    }

    /**
     * @dev Gen merkle moutain tree roots.It simply hash one by one
     * @param _trees A list of ranged tree root node
     */
    function genMMRRoot(RootNode[] memory _trees) internal pure returns (bytes32) {
        require(_trees.length > 0, "no tree");
        bytes32 _left = _trees[0].hash;
        for (uint256 i = 1; i < _trees.length; i++) {
            _left = keccak256(abi.encodePacked(_left, _trees[i].hash));
        }
        return _left;
    }
}
