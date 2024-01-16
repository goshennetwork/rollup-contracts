// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./RollupStateChain.sol";
import "./RollupInputChain.sol";
import "./ChainStorageContainer.sol";
import "../test-helper/TestBase.sol";

contract TestRollupStateChain is TestBase {
    address sender = address(0x7777); //admin
    address testAddress = address(0x8888);
    event StateBatchAppended(
        address indexed _proposer,
        uint64 indexed _startIndex,
        uint64 _timestamp,
        bytes32[] _blockHash
    );
    event StateRollbacked(uint64 indexed _stateIndex, bytes32 indexed _blockHash);

    function setUp() public {
        vm.startPrank(sender);
        super._initialize(sender);
        whitelist.setProposer(sender, true);
        whitelist.setProposer(sender, true);
        stakeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        vm.stopPrank();
    }

    /* test appendStateBatch
   1.test Fail
*/
    // test address not proposer
    function testAppendNotSequencer() public {
        bytes32[] memory states = new bytes32[](1);
        vm.startPrank(testAddress, testAddress);
        vm.expectRevert("only proposer");
        rollupStateChain.appendStateBatch(states, 0);
        vm.stopPrank();
    }

    // test in case of duplicated
    function testAppendDup() public {
        bytes32[] memory states = new bytes32[](1);
        vm.startPrank(sender);
        vm.expectRevert("start pos mismatch");
        rollupStateChain.appendStateBatch(states, 1);
        vm.stopPrank();
    }

    //test Fail proposer unstaked
    function testAppendBatchSequencerNoStaking() public {
        vm.startPrank(sender);
        whitelist.setProposer(testAddress, true);
        vm.stopPrank();
        bytes32[] memory states = new bytes32[](1);
        vm.startPrank(testAddress);
        vm.expectRevert("unstaked");
        rollupStateChain.appendStateBatch(states, 0);
        vm.stopPrank();
    }

    //append empty state
    function testAppend0() public {
        bytes32[] memory states;
        vm.startPrank(sender);
        vm.expectRevert("no block hashes");
        rollupStateChain.appendStateBatch(states, 0);
        vm.stopPrank();
    }

    // test stateChian.length + _blockHashes.length >  inputChain.length
    function testAppendExceedInputChain() public {
        vm.startPrank(address(rollupInputChain));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        vm.stopPrank();

        bytes32[] memory states = new bytes32[](5);
        vm.startPrank(sender);
        vm.expectRevert("exceed input chain height");
        rollupStateChain.appendStateBatch(states, 0);
        vm.stopPrank();
    }

    /* test appendStateBatch
   2.test pass
*/

    // test Append 2*batch (1+3 == 4) + test event
    function testAppendStateBatch() public {
        vm.startPrank(address(rollupInputChain));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        vm.stopPrank();

        vm.startPrank(sender);
        bytes32[] memory states = new bytes32[](1);
        //test eventEmit
        vm.expectEmit(true, true, false, true);
        emit StateBatchAppended(sender, 0, uint64(block.timestamp), states);
        rollupStateChain.appendStateBatch(states, 0);

        states = new bytes32[](3);
        rollupStateChain.appendStateBatch(states, 1);
        require(rollupStateChain.totalSubmittedState() == 4, "should 4");
        vm.stopPrank();
    }

    /* test rollbackStateBefore
   1.test Fail
*/
    // test address not challenge contract
    function testRollbackNotchallenge() public {
        vm.startPrank(testAddress, testAddress);
        Types.StateInfo memory stateInfo;
        stateInfo.timestamp = uint64(block.timestamp);
        stateInfo.proposer = sender;
        stateInfo.index = 3;

        vm.expectRevert("only permitted by challenge contract");
        rollupStateChain.rollbackStateBefore(stateInfo);
        vm.stopPrank();
    }

    // test verify state info
    function testRollbackStateInfoInvalid() public {
        vm.startPrank(challengerFactory, challengerFactory);
        Types.StateInfo memory stateInfo;
        stateInfo.timestamp = uint64(block.timestamp);
        stateInfo.proposer = sender;
        stateInfo.index = 3;
        vm.expectRevert("invalid state info");
        rollupStateChain.rollbackStateBefore(stateInfo);
        vm.stopPrank();
    }

    // test state info confirmed
    function testRollbackStateInfoConfirmed() public {
        vm.startPrank(address(rollupInputChain));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        vm.stopPrank();
        Types.StateInfo memory stateInfo;
        vm.startPrank(sender);
        bytes32[] memory states = new bytes32[](4);
        rollupStateChain.appendStateBatch(states, 0);
        stateInfo.timestamp = uint64(block.timestamp);
        stateInfo.proposer = sender;
        stateInfo.index = 3;
        vm.stopPrank();

        vm.startPrank(challengerFactory, challengerFactory);
        vm.warp(1000000000000); //it will always work
        vm.expectRevert("state confirmed");
        rollupStateChain.rollbackStateBefore(stateInfo);
        vm.stopPrank();
    }

    /* test rollbackStateBefore
   2.test Pass
*/
    // test rollback*3 & test event
    function testRollbackBefore() public {
        vm.startPrank(address(rollupInputChain));
        Types.StateInfo memory stateInfo;
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        vm.stopPrank();
        vm.startPrank(sender);
        bytes32[] memory states = new bytes32[](4);
        rollupStateChain.appendStateBatch(states, 0);
        require(rollupStateChain.totalSubmittedState() == 4, "4");
        stateInfo.timestamp = uint64(block.timestamp);
        stateInfo.proposer = sender;
        stateInfo.index = 3;
        vm.stopPrank();
        vm.startPrank(challengerFactory);

        //test eventEmit
        vm.expectEmit(true, true, false, true);
        emit StateRollbacked(stateInfo.index, stateInfo.blockHash);
        rollupStateChain.rollbackStateBefore(stateInfo);
        require(rollupStateChain.totalSubmittedState() == 3, "3");
        stateInfo.index = 1;
        rollupStateChain.rollbackStateBefore(stateInfo);
        require(rollupStateChain.totalSubmittedState() == 1, "1");
        stateInfo.index = 0;
        rollupStateChain.rollbackStateBefore(stateInfo);
        require(rollupStateChain.totalSubmittedState() == 0, "0");
    }
}
