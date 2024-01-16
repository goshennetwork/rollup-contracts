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
    RollupStateChain rollupStateChain;
    RollupInputChain rollupInputChain;
    L1CrossLayerWitness l1CrossLayerWitness;
    L2CrossLayerWitness l2CrossLayerWitness;
    VM vm = VM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    address sender = address(0x7777);
    address dao = address(0x6666);
    address challengerFactory = address(0x5555);

    function setUp() public {
        vm.startPrank(sender);
        addressManager = new AddressManager();
        addressManager.initialize();
        l1CrossLayerWitness = new L1CrossLayerWitness();
        l1CrossLayerWitness.initialize(address(addressManager));
        l2CrossLayerWitness = new L2CrossLayerWitness();
        ERC20 erc20 = new ERC20("test", "test");
        rollupStateChain = new RollupStateChain();
        rollupStateChain.initialize(address(addressManager), 3);
        StakingManager stakingManager = new StakingManager();
        stakingManager.initialize(address(addressManager), 0);
        rollupInputChain = new RollupInputChain();
        rollupInputChain.initialize(address(addressManager), 15000000, 3000000, 1234);
        ChainStorageContainer stateStorageContainer = new ChainStorageContainer();
        stateStorageContainer.initialize(AddressName.ROLLUP_STATE_CHAIN, address(addressManager));
        address stateStorage = address(stateStorageContainer);
        ChainStorageContainer inputStorageContainer = new ChainStorageContainer();
        inputStorageContainer.initialize(AddressName.ROLLUP_INPUT_CHAIN, address(addressManager));
        address inputStorage = address(inputStorageContainer);
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN, address(rollupInputChain));
        addressManager.setAddress(AddressName.STAKING_MANAGER, address(stakingManager));
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN_CONTAINER, stateStorage);
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN_CONTAINER, inputStorage);
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN, address(rollupStateChain));
        addressManager.setAddress(AddressName.L1_CROSS_LAYER_WITNESS, address(l1CrossLayerWitness));
        addressManager.setAddress(AddressName.L2_CROSS_LAYER_WITNESS, address(l2CrossLayerWitness));
        addressManager.setAddress(AddressName.DAO, sender);
        addressManager.setAddress(AddressName.STAKE_TOKEN, address(erc20));
        stakingManager.deposit();
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

        bytes[] memory list = new bytes[](15);
        list[13] = abi.encodePacked(_trees.rootHash);
        list[14] = abi.encodePacked(_trees.treeSize);
        bytes[] memory encodedList = new bytes[](15);
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
