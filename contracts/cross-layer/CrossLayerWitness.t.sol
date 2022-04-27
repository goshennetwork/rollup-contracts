// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../rollup/RollupStateChain.sol";
import "../rollup//RollupInputChain.sol";
import "../rollup/ChainStorageContainer.sol";
import "./L1CrossLayerWitness.sol";
import "./L2CrossLayerWitness.sol";
import "../libraries/Types.sol";
import { CompactMerkleTree, MerkleMountainRange } from "../libraries/MerkleMountainRange.sol";
import "../libraries/RLPWriter.sol";

interface VM {
    function prank(address sender) external;

    function warp(uint256 x) external;

    function startPrank(address sender) external;

    function stopPrank() external;
}

contract TestCrossLayerWitness {
    using MerkleMountainRange for CompactMerkleTree;
    CompactMerkleTree _trees;
    AddressManager addressManager;
    IRollupStateChain rollupStateChain;
    IRollupInputChain rollupInputChain;
    IL1CrossLayerWitness l1CrossLayerWitness;
    IL2CrossLayerWitness l2CrossLayerWitness;
    VM vm = VM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    address sender = address(0x7777);
    address dao = address(0x6666);
    address challengerFactory = address(0x5555);

    function setUp() public {
        vm.startPrank(sender);
        addressManager = new AddressManager();
        l1CrossLayerWitness = new L1CrossLayerWitness(address(addressManager));
        l2CrossLayerWitness = new L2CrossLayerWitness();
        ERC20 erc20 = new ERC20("test", "test");
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
        addressManager.newAddr(AddressName.L1_CROSS_LAYER_WITNESS, address(l1CrossLayerWitness));
        addressManager.newAddr(AddressName.L2_CROSS_LAYER_WITNESS, address(l2CrossLayerWitness));
        addressManager.newAddr(AddressName.DAO, sender);
        vm.stopPrank();
    }

    function testL1SendMsg() public {
        vm.startPrank(sender);
        l1CrossLayerWitness.sendMessage(address(0xdead), "test");
        require(l1CrossLayerWitness.totalSize() == 1, "size 1");
    }

    function testL2SendMsg() public {
        vm.startPrank(sender);
        l2CrossLayerWitness.sendMessage(address(0xdead), "test");
    }

    function testL1RelayMessage() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        MerkleMountainRange.appendLeafHash(_trees, _hash);
        bytes32[] memory _proof;

        bytes[] memory list = new bytes[](14);
        list[13] = abi.encodePacked(_trees.rootHash, _trees.treeSize);
        bytes[] memory encodedList = new bytes[](14);
        for (uint256 i = 0; i < list.length; i++) {
            encodedList[i] = RLPWriter.writeBytes(list[i]);
        }
        bytes memory rlpData = RLPWriter.writeList(encodedList);
        Types.StateInfo memory stateInfo;
        stateInfo.blockHash = keccak256(rlpData);
        vm.startPrank(address(rollupStateChain));
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.warp(3);
        vm.stopPrank();
        vm.startPrank(sender);
        bool success = l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
        require(success, "failed");
    }

    function testL2RelayMessage() public {
        vm.startPrank(Constants.L1_CROSS_LAYER_WITNESS);
        bool success = l2CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            bytes32(0),
            0
        );
        require(success, "failed");
    }

    function testL2ReplayMessage() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        MerkleMountainRange.appendLeafHash(_trees, _hash);
        bytes32[] memory _proof;

        vm.startPrank(Constants.L1_CROSS_LAYER_WITNESS);
        bool success = l2CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("fake()"),
            0,
            _trees.rootHash,
            _trees.treeSize
        );
        require(!success, "success");
        vm.stopPrank();
        vm.startPrank(sender);
        success = l2CrossLayerWitness.replayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            _proof,
            _trees.treeSize
        );
        require(success, "failed");
    }
}
