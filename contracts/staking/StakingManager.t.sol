// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../rollup/RollupStateChain.sol";
import "../rollup/RollupInputChain.sol";
import "../rollup/ChainStorageContainer.sol";
import "../test-helper/TestBase.sol";

contract TestStakingManager is TestBase {
    address sender = address(0x7777);

    function setUp() public {
        vm.startPrank(sender);
        super.initialize();
        vm.stopPrank();
    }

    function testDeposit() public {
        vm.startPrank(sender);
        feeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        require(stakingManager.isStaking(sender), "not staking");
    }

    function testWithdraw() public {
        vm.startPrank(sender);
        feeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        stakingManager.startWithdrawal();
        vm.stopPrank();
        vm.startPrank(address(rollupStateChain));

        Types.StateInfo memory stateInfo;
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.warp(fraudProofWindow);
        vm.stopPrank();
        vm.startPrank(sender);
        stakingManager.finalizeWithdrawal(stateInfo);
    }

    function testSlash() public {
        vm.startPrank(sender);
        feeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        vm.stopPrank();
        vm.startPrank(address(rollupStateChain));

        Types.StateInfo memory stateInfo;
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.stopPrank();
        vm.startPrank(challengerFactory);
        stakingManager.slash(0, Types.hash(stateInfo), sender);
        vm.warp(fraudProofWindow);
        stakingManager.claim(sender, stateInfo);
    }
}
