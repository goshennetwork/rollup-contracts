pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";

contract TestL2CrossLayerWitness is TestBase, L2CrossLayerWitness {
    using MerkleMountainRange for CompactMerkleTree;

    CompactMerkleTree _compactMerkleTree;
    address sender = address(0x7878);

    function setUp() public {
        _initialize();
    }

    function testCrossLayerSender() public {
        vm.expectRevert(bytes("no cross layer sender"));
        l2CrossLayerWitness.crossLayerSender();
    }

    function testRelayMessage() public {
        vm.startPrank(Constants.L1_CROSS_LAYER_WITNESS);
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        vm.expectEmit(true, true, true, true, address(l2CrossLayerWitness));
        emit MessageRelayed(0, _hash);
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

    function testRelayMessageWithWrongSender() public {
        vm.expectRevert("wrong sender");
        l2CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            bytes32(0),
            0
        );
    }

    function testRelayMessageWithAlreadyRelayed() public {
        vm.startPrank(Constants.L1_CROSS_LAYER_WITNESS);
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        bool success = l2CrossLayerWitness.mockSetSuccessRelayedMessages(_hash);
        require(success, "mockSetSuccessRelayedMessages failed");
        vm.expectRevert("already relayed");
        l2CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            bytes32(0),
            0
        );
    }

    function testRelayMessageWithCallFailed() public {
        vm.startPrank(Constants.L1_CROSS_LAYER_WITNESS);
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(0x9),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        vm.expectEmit(true, true, true, true, address(l2CrossLayerWitness));
        emit MessageRelayFailed(_hash, 0, 0);
        bool success = l2CrossLayerWitness.relayMessage(
            address(0x9),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            bytes32(0),
            0
        );
        require(!success, "failed");
    }

    function testReplayMessage() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        MerkleMountainRange.appendLeafHash(_trees, _hash);
        bytes32[] memory _proof;
        l2CrossLayerWitness.mockSetMmrRoot(_trees.treeSize, _hash);
        require(l2CrossLayerWitness.mmrRoots(_trees.treeSize) == _hash, "failed");
        vm.startPrank(sender);
        vm.expectEmit(true, true, true, true, address(l2CrossLayerWitness));
        emit MessageRelayed(0, _hash);
        bool success = l2CrossLayerWitness.replayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            _proof,
            _trees.treeSize
        );
        require(success, "failed");
    }

    function testReplayMessageWithUnknownMmrRoot() public {
        bytes32[] memory _proof;
        require(l2CrossLayerWitness.mmrRoots(_trees.treeSize) == bytes32(uint256(0)), "failed");
        vm.startPrank(sender);
        vm.expectRevert("unknown mmr root");
        l2CrossLayerWitness.replayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            _proof,
            _trees.treeSize
        );
    }

    function testReplayMessageWithMessageAlreadyRelayed() public {
        bytes32[] memory _proof;
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        MerkleMountainRange.appendLeafHash(_trees, _hash);
        l2CrossLayerWitness.mockSetSuccessRelayedMessages(_hash);
        l2CrossLayerWitness.mockSetMmrRoot(_trees.treeSize, _hash);
        require(l2CrossLayerWitness.successRelayedMessages(_hash), "failed");
        vm.startPrank(sender);
        vm.expectRevert("message already relayed");
        l2CrossLayerWitness.replayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            _proof,
            _trees.treeSize
        );
    }

    function testReplayMessageWithCallFailed() public {
        bytes32[] memory _proof;
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(0x9),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        MerkleMountainRange.appendLeafHash(_trees, _hash);
        l2CrossLayerWitness.mockSetMmrRoot(_trees.treeSize, _hash);
        vm.startPrank(sender);
        vm.expectEmit(true, true, true, true, address(l2CrossLayerWitness));
        emit MessageRelayFailed(_hash, _trees.treeSize, _hash);
        bool success = l2CrossLayerWitness.replayMessage(
            address(0x9),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            _proof,
            _trees.treeSize
        );
        require(!success, "failed");
    }

    function testSendMessage() public {
        vm.startPrank(sender);
        address MockL1Target = address(0x9899);
        bytes memory _message = bytes("0x01");
        vm.expectEmit(true, true, true, true, address(l2CrossLayerWitness));
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            MockL1Target,
            sender,
            _compactMerkleTree.treeSize,
            _message
        );
        bytes32 _mmrRoot = _compactMerkleTree.appendLeafHash(_hash);
        emit MessageSent(_compactMerkleTree.treeSize - 1, MockL1Target, sender, _mmrRoot, _message);
        l2CrossLayerWitness.sendMessage(MockL1Target, _message);
    }

    function testSendMessageWithCallerIsThis() public {
        vm.startPrank(address(l2CrossLayerWitness));
        address MockL1Target = address(0x9899);
        bytes memory _message = bytes("0x01");
        vm.expectRevert("wired situation");
        l2CrossLayerWitness.sendMessage(MockL1Target, _message);
    }
}
