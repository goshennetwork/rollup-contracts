// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";
import "./StateTransition.sol";

contract TestStateTransition is TestBase, StateTransition {
    address sender = address(0x7777); //admin

    function setUp() public {
        vm.startPrank(sender);
        super._initialize(sender);
        bytes32 root = hex"2ead174930137579e36ccc4d9b5d89c3ad532188617d6f74ed97a5e0c94f90b7";
        IMachineState _mstate = IMachineState(address(0));
        initialize(root, addressManager, _mstate);
        assert(getImageRoot(0) == root);
        assert(getImageRoot(1) == root);
        assert(getImageRoot(2) == root);
        vm.stopPrank();
    }

    function testUpgradeNewRoot() public {
        vm.startPrank(sender);
        bytes32 root = hex"2ead174930137579e36ccc4d9b5d89c3ad532188617d6f74ed97a5e0c94f90b7";
        console.log(msg.sender, address(this));
        this.upgradeToNewRoot(0, root);
        this.upgradeToNewRoot(1, root);
        // modify root
        this.upgradeToNewRoot(1, root);

        // should failed
        vm.expectRevert("duplicated upgrade");
        this.upgradeToNewRoot(0, root);

        whitelist.setProposer(sender, true);
        whitelist.setSequencer(sender, true);
        feeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        vm.stopPrank();
        vm.startPrank(address(rollupInputChain));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        addressManager.rollupInputChainContainer().append(bytes32(0));
        vm.stopPrank();
        vm.startPrank(sender);
        bytes32[] memory states = new bytes32[](3);
        rollupStateChain.appendStateBatch(states, 0);
        vm.expectRevert("ill batch index");
        this.upgradeToNewRoot(2, root);
        vm.stopPrank();
    }

    function testConsistentImageRoot() public {
        vm.startPrank(sender);
        bytes32 newRoot = hex"3ead174930137579e36ccc4d9b5d89c3ad532188617d6f74ed97a5e0c94f90b7";
        uint64 upgradeBatchIndex = 10;
        this.upgradeToNewRoot(upgradeBatchIndex, newRoot);
        for (uint64 i = 1; i < upgradeBatchIndex; i++) {
            assert(getImageRoot(i) == getImageRoot(0));
        }
        assert(getImageRoot(upgradeBatchIndex) == newRoot);
        assert(getImageRoot(upgradeBatchIndex + 1) == newRoot);

        bytes32 newRoot2 = hex"4ead174930137579e36ccc4d9b5d89c3ad532188617d6f74ed97a5e0c94f90b7";
        uint64 upgradeBatchIndex2 = 10;
        this.upgradeToNewRoot(upgradeBatchIndex2, newRoot2);
        for (uint64 i = 1; i < upgradeBatchIndex; i++) {
            assert(getImageRoot(i) == getImageRoot(0));
        }
        for (uint64 i = upgradeBatchIndex; i < upgradeBatchIndex2; i++) {
            assert(getImageRoot(i) == newRoot);
        }
        assert(getImageRoot(upgradeBatchIndex2) == newRoot2);
        assert(getImageRoot(upgradeBatchIndex2 + 1) == newRoot2);
        vm.stopPrank();
    }
}
