// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../rollup/RollupStateChain.sol";
import "../rollup/RollupInputChain.sol";
import "../rollup/ChainStorageContainer.sol";

interface VM {
    function prank(address sender) external;

    function warp(uint256 x) external;

    function startPrank(address sender) external;

    function stopPrank() external;
}

contract TestStakingManager {
    AddressManager addressManager;
    RollupStateChain rollupStateChain;
    IRollupInputChain rollupInputChain;
    IStakingManager stakingManager;
    VM vm = VM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    address sender = address(0x7777);
    address dao = address(0x6666);
    address challengerFactory;

    function setUp() public {
        vm.startPrank(sender);
        challengerFactory = address(new MockChallengeFactory());
        ERC20 erc20 = new ERC20("test", "test");
        addressManager = new AddressManager();
        rollupStateChain = new RollupStateChain(address(addressManager), 3);
        stakingManager = new StakingManager(dao, challengerFactory, address(rollupStateChain), address(erc20), 0);
        rollupInputChain = new RollupInputChain(address(addressManager), 2_000_000, 1_000_000);
        address stateStorage = address(
            new ChainStorageContainer(AddressName.ROLLUP_STATE_CHAIN, address(addressManager))
        );
        address inputStorage = address(
            new ChainStorageContainer(AddressName.ROLLUP_INPUT_CHAIN, address(addressManager))
        );
        addressManager.newAddr(AddressName.ROLLUP_INPUT_CHAIN, address(rollupInputChain));
        addressManager.newAddr(AddressName.STAKING_MANAGER, address(stakingManager));
        addressManager.newAddr(AddressName.ROLLUP_STATE_CHAIN_CONTAINER, stateStorage);
        addressManager.newAddr(AddressName.ROLLUP_INPUT_CHAIN_CONTAINER, inputStorage);
        addressManager.newAddr(AddressName.ROLLUP_STATE_CHAIN, address(rollupStateChain));
        addressManager.newAddr(AddressName.CHALLENGE_FACTORY, challengerFactory);
        vm.stopPrank();
    }

    function testDeposit() public {
        vm.startPrank(sender);
        stakingManager.deposit();
        require(stakingManager.isStaking(sender), "not staking");
    }

    function testWithdraw() public {
        vm.startPrank(sender);
        stakingManager.deposit();
        stakingManager.startWithdrawal();
        vm.stopPrank();
        vm.startPrank(address(rollupStateChain));

        Types.StateInfo memory stateInfo;
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.warp(3);
        vm.stopPrank();
        vm.startPrank(sender);
        stakingManager.finalizeWithdrawal(stateInfo);
    }

    function testSlash() public {
        vm.startPrank(sender);
        stakingManager.deposit();
        vm.stopPrank();
        vm.startPrank(address(rollupStateChain));

        Types.StateInfo memory stateInfo;
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.stopPrank();
        vm.startPrank(challengerFactory);
        stakingManager.slash(0, Types.hash(stateInfo), sender);
        vm.warp(3);
        stakingManager.claim(sender, stateInfo);
    }
}

contract MockChallengeFactory {
    function isChallengeContract(address _addr) external view returns (bool) {
        return _addr == address(this);
    }
}
