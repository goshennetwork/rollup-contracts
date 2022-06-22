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
    address toAddr = address(0x8989);
    uint256 amount = 0.5 ether;

    //    proposer deposit for staking.
    event Deposited(address indexed proposer, uint256 amount);
    //proposer start withdraw.
    event WithdrawStarted(address indexed proposer, uint256 needComfirmedBlock);
    //proposer finalize withdraw.
    event WithdrawFinalized(address indexed proposer, uint256 amount);
    //challenger slash the proposer
    event DepositSlashed(address indexed proposer, address indexed challenger, uint256 blockHeight, bytes32 _blockHash);
    //challenger or DAO gets deposit.
    event DepositClaimed(address indexed proposer, address indexed receiver, uint256 amount);

    function setUp() public {
        vm.startPrank(sender);
        super._initialize();
        vm.stopPrank();
    }

    function testDeposit() public {
        vm.startPrank(sender);
        require(!stakingManager.isStaking(sender), "not staking");
        feeToken.approve(address(stakingManager), stakingManager.price());
        vm.expectEmit(true, true, true, true, address(stakingManager));
        emit Deposited(sender, stakingManager.price());
        stakingManager.deposit();
        require(stakingManager.isStaking(sender), "not staking");
        vm.expectRevert("only unstacked user can deposit");
        stakingManager.deposit();
        feeToken.transfer(toAddr, amount);
        require(feeToken.balanceOf(toAddr) == amount);
        vm.stopPrank();
        vm.startPrank(toAddr);
        feeToken.approve(address(stakingManager), stakingManager.price());
        vm.expectRevert("ERC20: transfer amount exceeds balance");
        stakingManager.deposit();
    }

    function testWithdraw() public {
        vm.startPrank(sender);
        feeToken.approve(address(stakingManager), stakingManager.price());
        vm.expectRevert("not in staking");
        stakingManager.startWithdrawal();
        stakingManager.deposit();
        Types.StateInfo memory stateInfo;
        vm.expectRevert("not in withdrawing");
        stakingManager.finalizeWithdrawal(stateInfo);
        uint256 needConfirmedHeight = stakingManager.rollupStateChain().totalSubmittedState();
        vm.expectEmit(true, true, true, true, address(stakingManager));
        emit WithdrawStarted(sender, needConfirmedHeight);
        stakingManager.startWithdrawal();
        vm.expectRevert("not in staking");
        stakingManager.startWithdrawal();
        vm.expectRevert("incorrect state info");
        stakingManager.finalizeWithdrawal(stateInfo);
        vm.stopPrank();
        vm.startPrank(address(rollupStateChain));
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.stopPrank();
        vm.startPrank(sender);
        vm.expectRevert("provided state not confirmed");
        stakingManager.finalizeWithdrawal(stateInfo);
        vm.stopPrank();
        vm.startPrank(address(rollupStateChain));
        stateInfo.index = 1;
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.stopPrank();
        vm.startPrank(sender);
        vm.warp(fraudProofWindow);
        vm.expectRevert("should provide wanted state info");
        stakingManager.finalizeWithdrawal(stateInfo);
        vm.stopPrank();
        vm.startPrank(address(rollupStateChain));
        stateInfo.index = 0;
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.stopPrank();
        vm.startPrank(sender);
        vm.warp(fraudProofWindow);
        vm.expectEmit(true, true, true, true, address(stakingManager));
        emit WithdrawFinalized(sender, stakingManager.price());
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
