// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";
import "./DisputeTree.sol";

contract TestDisputeTree is TestBase {
    mapping(uint256 => DisputeTree.DisputeNode) testTree;
    uint128 N_SECTION = 7;

    //test middle special case
    function testMiddle() public pure {
        //common case
        uint128 return1 = DisputeTree.middle(1, 3);
        require(return1 == 2, "middle compute error");

        //when (_upper - _lower) / 2 == 0
        uint128 return2 = DisputeTree.middle(1, 2);
        require(return2 == 1, "middle compute error");

        //when _upper == _lower
        uint128 return3 = DisputeTree.middle(2, 2);
        require(return3 == 2, "middle compute error");
    }

    function testSectionFewSteps() public {
        uint256 _key0;
        uint256 _key1;
        uint256 _key2;
        uint128 _sections;
        uint128 _start = 0;
        uint128 _end = 2;
        uint128 _stepUpper = DisputeTree.midStep(N_SECTION - 1, 0, _start, _end);
        require(_stepUpper == 0, "wrong section0");
        _stepUpper = DisputeTree.midStep(N_SECTION - 1, 6, _start, _end);
        require(_stepUpper == 2, "wrong section1");
    }

    //test middle special case
    function testFailMiddle() public {
        //when _upper < _lower , revert
        vm.expectRevert("Arithmetic over/underflow");
        DisputeTree.middle(2, 1);
    }

    //test EncodeNodeKey result
    function testEncodeAndDecode(uint128 lower, uint128 upper) public pure {
        uint256 encoded = DisputeTree.encodeNodeKey(lower, upper);
        (uint128 _lower, uint128 _upper) = DisputeTree.decodeNodeKey(encoded);
        require(_lower == lower && upper == _upper, "decode lower error");
    }

    /* test addNewChild
   1.test Fail */

    //test Fail no parent node
    //when tree[_parentKey].parent == 0 , revert (parent not exist)
    function testAddNewChildWithNoParentNode() public {
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(0, address(1), 100);
        testTree[0] = node;
        vm.expectRevert("parent not exist");
        DisputeTree.addNewChild(testTree, 2, 1, 0, 100, address(1));
    }

    //test Fail if sub state not proven
    //when sub state is proven, the node already exist and its expire time is zero
    function testAddNewChildWithSubStateNotProven() public {
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(1, address(1), 100);
        testTree[DisputeTree.encodeNodeKey(1, 10)] = node;
        vm.expectRevert("Err Node");
        DisputeTree.addNewChild(testTree, 2, 1, DisputeTree.encodeNodeKey(1, 10), 100, address(1));
    }

    //test Fail already init
    //when selected branch has already been selected revert (Err Node)
    function testAddNewChildAlreadySelected() public {
        uint256 return1 = DisputeTree.encodeNodeKey(1, 10);
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(1, address(1), 100);
        testTree[return1] = node;
        //(5,7)
        (uint128 stepLower, uint128 stepUpper) = DisputeTree.decodeNodeKey(return1);
        stepLower = DisputeTree.middle(stepLower, stepUpper);
        //(6,7)
        uint256 _childKey = DisputeTree.encodeNodeKey(stepLower, stepUpper);
        DisputeTree.DisputeNode storage childnode = testTree[_childKey];
        childnode.parent = 1;
        DisputeTree.addNewChild(testTree, 2, 1, return1, 100, address(1));
        vm.expectRevert("Err Node");
        DisputeTree.addNewChild(testTree, 2, 1, return1, 100, address(1));
    }

    /* test addNewChild
2.test Pass */
    //test return childkey
    function testAddNewChildPass() public {
        uint256 return1 = DisputeTree.encodeNodeKey(1, 10);
        //compute return _childkey
        (uint128 stepLower, uint128 stepUpper) = DisputeTree.decodeNodeKey(return1);
        stepUpper = DisputeTree.midStep(N_SECTION - 1, 0, stepLower, stepUpper);
        uint256 _childKey = DisputeTree.encodeNodeKey(stepLower, stepUpper);

        testTree[return1] = DisputeTree.DisputeNode(1, address(1), 100);
        for (uint128 i = 0; i < N_SECTION; i++) {
            stepUpper = DisputeTree.midStep(N_SECTION - 1, i, 1, 10);
            testTree[DisputeTree.encodeNodeKey(stepLower, stepUpper)] = DisputeTree.DisputeNode(return1, address(0), 0);
            stepLower = stepUpper;
        }
        uint256 returnChildkey = DisputeTree.addNewChild(testTree, N_SECTION, 0, return1, 100, address(1));
        //test returnChildkey same as _childkey
        require(returnChildkey == _childKey, "return childkey invalid");
    }

    /*test getFirstLeafNode*/
    function testGetFirstLeafNode() public {
        uint256 return1 = DisputeTree.encodeNodeKey(1, 10);
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(return1, address(1), 100);
        testTree[return1] = node;
        uint128 _stepLower = 1;
        for (uint128 i = 0; i < N_SECTION; i++) {
            uint128 _stepUpper = DisputeTree.midStep(N_SECTION - 1, i, 1, 10);
            testTree[DisputeTree.encodeNodeKey(_stepLower, _stepUpper)] = DisputeTree.DisputeNode(
                return1,
                address(0),
                0
            );
            _stepLower = _stepUpper;
        }
        //root ==> [1,10]
        //1.no child getFirstLeafNode
        (uint256 key1, uint256 depth1, bool oneBranch1) = DisputeTree.getFirstLeafNode(testTree, N_SECTION, return1);
        (uint128 stepLower1, uint128 stepUpper1) = DisputeTree.decodeNodeKey(key1);
        require(stepLower1 == 1 && stepUpper1 == 10, "no child case error");
        require(depth1 == 1 && oneBranch1 == true, "no child case depth&branch error");
        console.log("1,10");
        //2.add child [5,10]
        DisputeTree.addNewChild(testTree, N_SECTION, 1, return1, 100, address(1));
        (uint256 key2, uint256 depth2, bool oneBranch2) = DisputeTree.getFirstLeafNode(testTree, N_SECTION, return1);
        (uint128 stepLower2, uint128 stepUpper2) = DisputeTree.decodeNodeKey(key2);
        require(stepLower2 == 2 && stepUpper2 == 3, "add child case error");
        require(depth2 == 2 && oneBranch2 == true, "add child case depth&branch error");

        //3.one step return
        uint256 return2 = DisputeTree.encodeNodeKey(1, 2);
        DisputeTree.DisputeNode memory node2 = DisputeTree.DisputeNode(1, address(1), 100);
        testTree[return2] = node2;
        (uint256 key3, uint256 depth3, bool oneBranch3) = DisputeTree.getFirstLeafNode(testTree, N_SECTION, return2);
        (uint128 stepLower3, uint128 stepUpper3) = DisputeTree.decodeNodeKey(key3);
        require(stepLower3 == 1 && stepUpper3 == 2, "one step case error");
        require(depth3 == 1 && oneBranch3 == true, "one step case depth&branch error");
    }
}
