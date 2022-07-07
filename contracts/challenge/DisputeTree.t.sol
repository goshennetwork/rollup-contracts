// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";
import "./DisputeTree.sol";

contract TestDisputeTree is TestBase {
    mapping(uint256 => DisputeTree.DisputeNode) testTree;

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

    //test Search Node not find
    function testSearchNodeWithMidStepNotFind() public {
        vm.expectRevert("not found");
        DisputeTree.searchNodeWithMidStep(1, 10, 10);
    }

    //test Search Node pass
    function testSearchNode2() public pure {
        for (uint256 i = 2; i < 10; i++) {
            uint256 return1 = DisputeTree.searchNodeWithMidStep(1, 10, i);
            (uint256 lower, uint256 upper) = DisputeTree.decodeNodeKey(return1);
            require((lower + upper) / 2 == i);
        }
    }

    /* test addNewChild
   1.test Fail */

    //test Fail no parent node
    //when tree[_parentKey].parent == 0 , revert (parent not exist)
    function testAddNewChildWithNoParentNode() public {
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(0, address(1), 100, bytes32("0x0"));
        testTree[0] = node;
        vm.expectRevert("parent not exist");
        DisputeTree.addNewChild(testTree, 0, false, 100, address(1));
    }

    //test Fail mid state not proven
    //when tree[_parentKey].midStateRoot == 0 , revert (parent mid state not proven)
    function testAddNewChildWithMidStateRootNotProven() public {
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(1, address(1), 100, bytes32(0));
        testTree[0] = node;
        vm.expectRevert("parent mid state not proven");
        DisputeTree.addNewChild(testTree, 0, false, 100, address(1));
    }

    //test Fail already init
    //when node.parent != 0 , revert (already init)
    function testAddNewChildAlreadyInit() public {
        uint256 return1 = DisputeTree.searchNodeWithMidStep(1, 10, 6);
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(1, address(1), 100, bytes32("0x0"));
        testTree[return1] = node;
        //(5,7)
        (uint128 stepLower, uint128 stepUpper) = DisputeTree.decodeNodeKey(return1);
        stepLower = DisputeTree.middle(stepLower, stepUpper);
        //(6,7)
        uint256 _childKey = DisputeTree.encodeNodeKey(stepLower, stepUpper);
        DisputeTree.DisputeNode storage childnode = testTree[_childKey];
        childnode.parent = 1;
        vm.expectRevert("already init");
        DisputeTree.addNewChild(testTree, return1, false, 100, address(1));
    }

    /* test addNewChild
2.test Pass */
    //test return childkey
    function testAddNewChildPass() public {
        uint256 return1 = DisputeTree.encodeNodeKey(1, 10);
        //compute return _childkey
        (uint128 stepLower, uint128 stepUpper) = DisputeTree.decodeNodeKey(return1);
        stepLower = DisputeTree.middle(stepLower, stepUpper);
        uint256 _childKey = DisputeTree.encodeNodeKey(stepLower, stepUpper);

        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(1, address(1), 100, bytes32("0x0"));
        testTree[return1] = node;
        uint256 returnChildkey = DisputeTree.addNewChild(testTree, return1, false, 100, address(1));
        //test returnChildkey same as _childkey
        require(returnChildkey == _childKey, "return childkey invalid");
    }

    /*test getFirstLeafNode*/
    function testGetFirstLeafNode() public {
        uint256 return1 = DisputeTree.encodeNodeKey(1, 10);
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(1, address(1), 100, bytes32("0x0"));
        testTree[return1] = node;
        //root ==> [1,10]
        //1.no child getFirstLeafNode
        (uint256 key1, uint256 depth1, bool oneBranch1) = DisputeTree.getFirstLeafNode(testTree, return1);
        (uint128 stepLower1, uint128 stepUpper1) = DisputeTree.decodeNodeKey(key1);
        require(stepLower1 == 1 && stepUpper1 == 10, "no child case error");
        require(depth1 == 1 && oneBranch1 == true, "no child case depth&branch error");

        //2.add child [5,10]
        DisputeTree.addNewChild(testTree, return1, false, 100, address(1));
        (uint256 key2, uint256 depth2, bool oneBranch2) = DisputeTree.getFirstLeafNode(testTree, return1);
        (uint128 stepLower2, uint128 stepUpper2) = DisputeTree.decodeNodeKey(key2);
        require(stepLower2 == 5 && stepUpper2 == 10, "add child case error");
        require(depth2 == 2 && oneBranch2 == true, "add child case depth&branch error");

        //3.one step return
        uint256 return2 = DisputeTree.encodeNodeKey(1, 2);
        DisputeTree.DisputeNode memory node2 = DisputeTree.DisputeNode(1, address(1), 100, bytes32("0x0"));
        testTree[return2] = node2;
        (uint256 key3, uint256 depth3, bool oneBranch3) = DisputeTree.getFirstLeafNode(testTree, return2);
        (uint128 stepLower3, uint128 stepUpper3) = DisputeTree.decodeNodeKey(key3);
        require(stepLower3 == 1 && stepUpper3 == 2, "one step case error");
        require(depth3 == 1 && oneBranch3 == true, "one step case depth&branch error");
    }

    /*test removeSelfBranch*/

    //test Fail when tree[_leafKey].parent == 0
    function testFailRemoveSelfBranch() public {
        uint256 leafkey1 = DisputeTree.encodeNodeKey(1, 10);
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(0, address(1), 100, bytes32("0x0"));
        testTree[leafkey1] = node;
        DisputeTree.removeSelfBranch(testTree, leafkey1);
    }

    //test-remove
    function testRemoveSelfBranch() public {
        //1.root node condition: _parentKey == _leafKey
        uint256 leafkey1 = DisputeTree.encodeNodeKey(1, 10);
        DisputeTree.DisputeNode memory node = DisputeTree.DisputeNode(leafkey1, address(1), 100, bytes32("0x0"));
        testTree[leafkey1] = node;
        DisputeTree.removeSelfBranch(testTree, leafkey1);
        require(testTree[leafkey1].parent == 0, "root node remove error");

        //2.remove one branch
        uint256 leafkey2 = DisputeTree.encodeNodeKey(1, 10);
        DisputeTree.DisputeNode memory node2 = DisputeTree.DisputeNode(leafkey2, address(1), 100, bytes32("0x0"));
        testTree[leafkey2] = node2;
        // add two child , get two branches
        DisputeTree.addNewChild(testTree, leafkey2, false, 100, address(1));
        DisputeTree.addNewChild(testTree, leafkey2, true, 100, address(1));

        //before remove , testTree[1,10] first LeafNode is [1,5]
        (uint256 key1, , ) = DisputeTree.getFirstLeafNode(testTree, leafkey2);
        (uint128 stepLower1, uint128 stepUpper1) = DisputeTree.decodeNodeKey(key1);
        require(stepLower1 == 1 && stepUpper1 == 5, "before remove error");

        //remove
        uint256 leafkey3 = DisputeTree.encodeNodeKey(1, 5);
        DisputeTree.removeSelfBranch(testTree, leafkey3);

        //after remove , testTree[1,10] first LeafNode is [5,10]
        (uint256 key2, , ) = DisputeTree.getFirstLeafNode(testTree, leafkey2);
        (uint128 stepLower2, uint128 stepUpper2) = DisputeTree.decodeNodeKey(key2);
        require(stepLower2 == 5 && stepUpper2 == 10, "remove one branch error");
    }
}
