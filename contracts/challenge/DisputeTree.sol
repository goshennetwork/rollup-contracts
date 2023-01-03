// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

library DisputeTree {
    struct DisputeNode {
        uint256 parent;
        address challenger;
        uint256 expireAfterBlock;
        bytes32 endStateRoot;
    }

    function middle(uint128 _lower, uint128 _upper) internal pure returns (uint128) {
        return _lower + (_upper - _lower) / 2;
    }

    /// @dev split steps to pieces section.(the max piece is _nSection, but may be lower than this
    /// if left step num less than n section)
    /// @notice can't divide one step
    function nSection(
        uint128 _nSection,
        uint128 _n,
        uint128 _lower,
        uint128 _upper
    ) internal pure returns (uint128,uint128, uint128) {
        uint128 _stepNum = _upper - _lower;
        require(_stepNum > 1, "can't divide oneStep");
        if (_stepNum < _nSection) {
            // @dev if step number smaller than n section, then n section is equal to step number
            _nSection = _stepNum;
        }
        require(_n < _nSection, "Out of N Section");
        uint128 _newLower = _lower;
        uint128 _newUpper = _upper;
        uint128 _piece = (_stepNum) / _nSection;
        _newLower = _lower + _piece * _n;
        if (_n + 1 != _nSection) {
            /// @dev not last
            _newUpper = _lower + _piece * (_n + 1);
        }
        return (_nSection,_newLower, _newUpper);
    }

    function encodeNodeKey(uint128 _stepLower, uint128 _stepUpper) internal pure returns (uint256) {
        return uint256(_stepLower) + (uint256(_stepUpper) << 128);
    }

    function decodeNodeKey(uint256 nodeKey) internal pure returns (uint128 stepLower, uint128 stepUpper) {
        stepLower = uint128(nodeKey);
        stepUpper = uint128(nodeKey >> 128);
    }

    function searchNodeWithEndStep(
        uint128 _stepLower,
        uint128 _stepUpper,
        uint256 _EndStep
    ) internal pure returns (uint256) {
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
        uint128 _NSection,
        uint128 _Nth,
        uint256 _parentKey,
        uint256 _expireAfterBlock,
        address _challenger
    ) internal returns (uint256) {
        DisputeNode storage parent = tree[_parentKey];
        require(parent.parent != 0, "parent not exist");
        (uint128 stepLower, uint128 stepUpper) = decodeNodeKey(_parentKey);
        require(stepUpper > stepLower + 1, "one step have no child");
        (,stepLower, stepUpper) = nSection(_NSection, _Nth, stepLower, stepUpper);
        uint256 _childKey = encodeNodeKey(stepLower, stepUpper);
        DisputeNode storage node = tree[_childKey];
        require(node.parent != 0 && node.expireAfterBlock==0, "already init");
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
        (uint128 _stepLower, uint128 _stepUpper) = decodeNodeKey(_rootKey);
        uint128 _tempStepLower;
        uint128 _tempStepUpper;
        while (_stepUpper - _stepLower > 1) {
            _depth++;
            uint256 _nodeKey;
            uint256 _tempNextNodeKey;
            /// @dev maybe remained step num less than n section.
            uint128 _tempNSection=_nSection;
            for (uint128 i = 0; i < _tempNSection; i++) {
                (_tempNSection,_tempStepLower, _tempStepUpper) = nSection(_nSection, i, _stepLower, _stepUpper);
                uint256 _nodeKey = encodeNodeKey(_tempStepLower, _tempStepUpper);
                if (tree[_nodeKey].parent != 0) {
                    if (_tempNextNodeKey != 0) {
                        /// @notice deplicated, so there is more than one branch todo: maybe just return, no need to go on?
                        _oneBranch = false;
                        break;
                    }
                    _tempNextNodeKey = _nodeKey;
                }
            }
            if (_tempNextNodeKey == 0) {
                /// @dev no child
                return (encodeNodeKey(_stepLower, _stepUpper), _depth, _oneBranch);
            }
            /// exist next node
            (_stepLower, _stepUpper) = decodeNodeKey(_tempNextNodeKey);
        }
        _depth++;
        //find one step, one step is surely leaf.
        return (encodeNodeKey(_stepLower, _stepUpper), _depth, _oneBranch);
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
