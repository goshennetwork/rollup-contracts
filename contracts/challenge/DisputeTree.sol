// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

library DisputeTree {
    struct DisputeNode {
        uint256 parent;
        address challenger;
        uint256 expireAfterBlock;
        bytes32 midStateRoot;
    }

    function middle(uint128 _lower, uint128 _upper) internal pure returns (uint128) {
        return _lower + (_upper - _lower) / 2;
    }

    function encodeNodeKey(uint128 _stepLower, uint128 _stepUpper) internal pure returns (uint256) {
        return uint256(_stepLower) + (uint256(_stepUpper) << 128);
    }

    function decodeNodeKey(uint256 nodeKey) internal pure returns (uint128 stepLower, uint128 stepUpper) {
        stepLower = uint128(nodeKey);
        stepUpper = uint128(nodeKey >> 128);
    }

    function searchNodeWithMidStep(
        uint128 _stepLower,
        uint128 _stepUpper,
        uint256 _midStep
    ) internal view returns (uint256) {
        while (_stepUpper - _stepLower > 1) {
            uint128 _stateStep = middle(_stepLower, _stepUpper);
            if (_midStep < _stateStep) {
                //so wanted is in left child.
                _stepUpper = _stateStep;
            } else if (_midStep > _stateStep) {
                //so wanted is in right.
                _stepLower = _stateStep;
            } else {
                //find out.
                return encodeNodeKey(_stepLower, _stepUpper);
            }
        }
        revert("not found");
    }

    function addNewChild(
        mapping(uint256 => DisputeNode) storage tree,
        uint256 _parentKey,
        bool _isLeftChild,
        uint256 _expireAfterBlock,
        address _challenger
    ) internal returns (uint256) {
        DisputeNode storage parent = tree[_parentKey];
        require(parent.parent != 0, "parent not exist");
        require(parent.midStateRoot != 0, "parent mid state not proven");
        (uint128 stepLower, uint128 stepUpper) = decodeNodeKey(_parentKey);
        require(stepUpper > stepLower + 1, "one step have no child");
        if (_isLeftChild) {
            stepUpper = middle(stepLower, stepUpper);
        } else {
            stepLower = middle(stepLower, stepUpper);
        }
        uint256 _childKey = encodeNodeKey(stepLower, stepUpper);
        DisputeNode storage node = tree[_childKey];
        require(node.parent == 0, "already init");
        node.parent = _parentKey;
        node.challenger = _challenger;
        node.expireAfterBlock = _expireAfterBlock;
        return _childKey;
    }

    function isChildNode(uint256 _parentKey, uint256 _childKey) internal pure returns (bool) {
        (uint128 parentLower, uint128 parentUpper) = decodeNodeKey(_parentKey);
        (uint128 childLower, uint128 childUpper) = decodeNodeKey(_childKey);
        return _parentKey != _childKey && childLower >= parentLower && childUpper <= parentUpper;
    }

    /**
     * @dev Get the lowest branch in the disputeNode tree
     */
    function getFirstLeafNode(mapping(uint256 => DisputeNode) storage tree, uint256 _rootKey)
        internal
        view
        returns (
            uint256,
            uint64,
            bool
        )
    {
        uint64 _depth;
        bool _oneBranch = true;
        (uint128 _stepLower, uint128 _stepUpper) = decodeNodeKey(_rootKey);
        while (_stepUpper - _stepLower > 1) {
            _depth++;
            uint128 _stepMid = middle(_stepLower, _stepUpper);
            //now check branch.
            bool _leftChildExist = tree[encodeNodeKey(_stepLower, _stepMid)].parent != 0;
            bool _rightChildExist = tree[encodeNodeKey(_stepMid, _stepUpper)].parent != 0;
            if (_leftChildExist && _rightChildExist) {
                _oneBranch = false;
            }

            if (_leftChildExist) {
                //find left child,
                _stepUpper = _stepMid;
                continue;
            }
            //not left,maybe right
            if (_rightChildExist) {
                //find right child,
                _stepLower = _stepMid;
                continue;
            }

            // no child
            return (encodeNodeKey(_stepLower, _stepUpper), _depth, _oneBranch);
        }
        _depth++;
        //find one step, one step is surely leaf.
        return (encodeNodeKey(_stepLower, _stepUpper), _depth, _oneBranch);
    }

    function removeSelfBranch(mapping(uint256 => DisputeNode) storage tree, uint256 _leafKey) internal {
        DisputeNode storage node = tree[_leafKey];
        require(node.parent > 0);
        while (true) {
            DisputeNode storage childNode = tree[_leafKey];

            uint256 _parentKey = childNode.parent;
            // remove
            childNode.parent = 0;
            if (_parentKey == 0 || _parentKey == _leafKey) {
                // root node condition: _parentKey == _leafKey
                return;
            }
            DisputeNode storage parentNode = tree[_parentKey];
            assert(parentNode.parent > 0);

            (uint128 stepLower, uint128 stepUpper) = decodeNodeKey(_parentKey);
            (uint128 childStepLower, ) = decodeNodeKey(_leafKey);
            uint256 _siblingKey;
            if (stepLower == childStepLower) {
                _siblingKey = encodeNodeKey(middle(stepLower, stepUpper), stepUpper);
            } else {
                _siblingKey = encodeNodeKey(stepLower, middle(stepLower, stepUpper));
            }

            if (tree[_siblingKey].parent == 0) {
                //parent have no branch,just remove it and go on.
                _leafKey = _parentKey;
            } else {
                //have branch,terminate
                return;
            }
        }
        revert("should never happen");
    }
}
