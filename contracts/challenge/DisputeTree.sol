// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

library DisputeTree {
    struct DisputeNode {
        uint256 parent;
        address challenger;
        uint256 expireAfterBlock;
    }

    function middle(uint128 _lower, uint128 _upper) internal pure returns (uint128) {
        return _lower + (_upper - _lower) / 2;
    }

    /// @dev mid step return the upper step of a provided steps in an interval
    /// @param midSteps the num of mid step of an interval, bisection only have one mid step
    /// @param _i the i'th piece of the interval, should <= midSteps
    /// @param _lower the start num of interval
    /// @param _upper the end num of interval
    function midStep(
        uint128 midSteps,
        uint128 _i,
        uint128 _lower,
        uint128 _upper
    ) internal pure returns (uint128) {
        return ((_upper - _lower) * (_i + 1)) / (midSteps + 1) + _lower;
    }

    function encodeNodeKey(uint128 _stepLower, uint128 _stepUpper) internal pure returns (uint256) {
        return uint256(_stepLower) + (uint256(_stepUpper) << 128);
    }

    function decodeNodeKey(uint256 nodeKey) internal pure returns (uint128 stepLower, uint128 stepUpper) {
        stepLower = uint128(nodeKey);
        stepUpper = uint128(nodeKey >> 128);
    }

    //    function searchNodeWithEndStep(
    //        uint128 _stepLower,
    //        uint128 _stepUpper,
    //        uint256 _EndStep
    //    ) internal pure returns (uint256) {
    //        while (_stepUpper - _stepLower > 1) {
    //            uint128 _stateStep = middle(_stepLower, _stepUpper);
    //            if (_midStep < _stateStep) {
    //                //so wanted is in left child.
    //                _stepUpper = _stateStep;
    //            } else if (_midStep > _stateStep) {
    //                //so wanted is in right.
    //                _stepLower = _stateStep;
    //            } else {
    //                //find out.
    //                return encodeNodeKey(_stepLower, _stepUpper);
    //            }
    //        }
    //        revert("not found");
    //    }

    function addNewChild(
        mapping(uint256 => DisputeNode) storage tree,
        uint128 _NSection,
        uint128 _Nth,
        uint256 _parentKey,
        uint256 _expireAfterBlock,
        address _challenger
    ) internal returns (uint256) {
        DisputeNode storage parent = tree[_parentKey];
        require(parent.parent != 0, "parent not exist");
        require(_Nth < _NSection, "Err Nth");
        (uint128 stepLower, uint128 stepUpper) = decodeNodeKey(_parentKey);
        require(stepUpper > stepLower + 1, "one step have no child");
        uint128 _childStepUpper = midStep(_NSection - 1, _Nth, stepLower, stepUpper);
        uint128 _childStepLower = stepLower;
        if (_Nth > 0) {
            _childStepLower = midStep(_NSection - 1, _Nth - 1, stepLower, stepUpper);
        }
        uint256 _childKey = encodeNodeKey(_childStepLower, _childStepUpper);
        DisputeNode storage node = tree[_childKey];
        require(node.parent != 0 && node.expireAfterBlock == 0, "Err Node");
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
    function getFirstLeafNode(
        mapping(uint256 => DisputeNode) storage tree,
        uint128 _nSection,
        uint256 _rootKey
    )
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
        (uint128 _stepStart, uint128 _stepEnd) = decodeNodeKey(_rootKey);
        while (_stepEnd - _stepStart > 1) {
            uint128 _stepLower = _stepStart;
            _depth++;
            uint256 _tempNextNodeKey;
            /// @dev maybe remained step num less than n section.
            for (uint128 i = 0; i < _nSection; i++) {
                uint128 _stepUpper = midStep(_nSection - 1, i, _stepStart, _stepEnd);
                uint256 _nodeKey = encodeNodeKey(_stepLower, _stepUpper);
                if (tree[_nodeKey].parent != 0 && tree[_nodeKey].expireAfterBlock != 0) {
                    if (_tempNextNodeKey != 0) {
                        /// @notice duplicated, so there is more than one branch todo: maybe just return, no need to go on?
                        _oneBranch = false;
                        return (encodeNodeKey(_stepStart, _stepEnd), _depth, _oneBranch);
                    }
                    _tempNextNodeKey = _nodeKey;
                }
                _stepLower = _stepUpper;
            }
            if (_tempNextNodeKey == 0) {
                /// @dev no child
                return (encodeNodeKey(_stepStart, _stepEnd), _depth, _oneBranch);
            }
            /// exist next node
            (_stepStart, _stepEnd) = decodeNodeKey(_tempNextNodeKey);
        }
        _depth++;
        //find one step, one step is surely leaf.
        return (encodeNodeKey(_stepStart, _stepEnd), _depth, _oneBranch);
    }

    //    function removeSelfBranch(mapping(uint256 => DisputeNode) storage tree, uint256 _leafKey) internal {
    //        DisputeNode storage node = tree[_leafKey];
    //        require(node.parent > 0);
    //        while (true) {
    //            DisputeNode storage childNode = tree[_leafKey];
    //
    //            uint256 _parentKey = childNode.parent;
    //            // remove
    //            childNode.parent = 0;
    //            if (_parentKey == 0 || _parentKey == _leafKey) {
    //                // root node condition: _parentKey == _leafKey
    //                return;
    //            }
    //            DisputeNode storage parentNode = tree[_parentKey];
    //            assert(parentNode.parent > 0);
    //
    //            (uint128 stepLower, uint128 stepUpper) = decodeNodeKey(_parentKey);
    //            (uint128 childStepLower, ) = decodeNodeKey(_leafKey);
    //            uint256 _siblingKey;
    //            if (stepLower == childStepLower) {
    //                _siblingKey = encodeNodeKey(middle(stepLower, stepUpper), stepUpper);
    //            } else {
    //                _siblingKey = encodeNodeKey(stepLower, middle(stepLower, stepUpper));
    //            }
    //
    //            if (tree[_siblingKey].parent == 0) {
    //                //parent have no branch,just remove it and go on.
    //                _leafKey = _parentKey;
    //            } else {
    //                //have branch,terminate
    //                return;
    //            }
    //        }
    //        revert("should never happen");
    //    }
}
