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
    address sender = address(0x7777);

    function setUp() public {
        vm.startPrank(sender);
        super.initialize();
        dao.setProposerWhitelist(sender, true);
        dao.setSequencerWhitelist(sender, true);
        feeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        vm.stopPrank();
    }

    //append empty state
    function testFailAppend0() public {
        bytes32[] memory states;
        vm.startPrank(sender);
        rollupStateChain.appendStateBatch(states, 0);
    }

    function testFailAppendDup() public {
        bytes32[] memory states = new bytes32[](1);
        vm.startPrank(sender);
        rollupStateChain.appendStateBatch(states, 1);
    }

    function testAppend() public {
        vm.startPrank(address(rollupInputChain));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        vm.stopPrank();
        vm.startPrank(sender);
        bytes32[] memory states = new bytes32[](1);
        rollupStateChain.appendStateBatch(states, 0);
        states = new bytes32[](3);
        rollupStateChain.appendStateBatch(states, 1);
        require(rollupStateChain.totalSubmittedState() == 4, "should 4");
    }

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
