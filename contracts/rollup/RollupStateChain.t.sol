// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./RollupStateChain.sol";
import "./RollupInputChain.sol";
import "./ChainStorageContainer.sol";

interface VM {
    function prank(address sender) external;

    function startPrank(address sender) external;

    function stopPrank() external;
}

contract TestRollupStateChain {
    AddressManager addressManager;
    RollupStateChain rollupStateChain;
    IRollupInputChain rollupInputChain;
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
        IStakingManager stakingManager = new StakingManager(
            dao,
            challengerFactory,
            address(rollupStateChain),
            address(erc20),
            0
        );
        stakingManager.deposit();
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

contract MockChallengeFactory {
    function isChallengeContract(address _addr) external view returns (bool) {
        return _addr == address(this);
    }
}
