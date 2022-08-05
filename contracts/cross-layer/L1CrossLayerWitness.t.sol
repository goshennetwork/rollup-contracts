pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";

contract TestL1CrossLayerWitness is TestBase, L1CrossLayerWitness {
    using MerkleMountainRange for CompactMerkleTree;

    CompactMerkleTree _compactMerkleTree;
    address sender = address(0x7878);

    function setUp() public {
        _initialize(sender);
    }

    function testCrossLayerSender() public {
        vm.expectRevert(bytes("no cross layer sender"));
        l1CrossLayerWitness.crossLayerSender();
    }

    function testRelayMessageCallL1System() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(rollupInputChain),
            sender,
            0,
            abi.encodeWithSignature("chainHeight()")
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
        vm.startPrank(address(addressManager));
        vm.expectRevert(bytes("can't relay message to l1 system"));
        l1CrossLayerWitness.relayMessage(
            address(rollupInputChain),
            sender,
            abi.encodeWithSignature("chainHeight()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testRelayMessageWithWrongStateInfo() public {
        // set wrong sender
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
        Types.StateInfo memory failStateInfo;
        Types.StateInfo memory stateInfo;
        stateInfo.blockHash = keccak256(rlpData);
        vm.startPrank(address(rollupStateChain));
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.warp(3);
        vm.stopPrank();
        vm.startPrank(address(l1CrossLayerWitness));
        vm.expectRevert(bytes("wrong state info"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            failStateInfo,
            _proof
        );
    }

    function testRelayMessageWithStateNotConfirmedYed() public {
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
        vm.warp(0);
        vm.stopPrank();
        vm.startPrank(address(l1CrossLayerWitness));
        vm.expectRevert(bytes("state not confirmed yet"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testRelayMessageWithWrongBlockProvide() public {
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
        // set wrong stateInfo.blockHash
        stateInfo.blockHash = keccak256(encodedList[1]);
        vm.startPrank(address(rollupStateChain));
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.warp(3);
        vm.stopPrank();
        vm.startPrank(address(l1CrossLayerWitness));
        vm.expectRevert(bytes("wrong block provide"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testRelayMessageWithWrongData() public {
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
        sender = address(0x9);
        vm.expectRevert(bytes("mmr root differ"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testRelayMessageWithWrongTreeSize() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            1,
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
        vm.expectRevert(bytes("leaf index out of bounds"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            1,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testRelayMessageWithMessageAlreadyExists() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        l1CrossLayerWitness.mockSetSuccessRelayedMessages(_hash);
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
        vm.expectRevert(bytes("provided message already been relayed"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testRelayMessageWithMessageBlocked() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(addressManager),
            sender,
            0,
            abi.encodeWithSignature("dao()")
        );
        l1CrossLayerWitness.mockSetBlockedMessages(_hash);
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
        vm.expectRevert(bytes("message blocked"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testRelayMessageWithWrongTargetAddr() public {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            address(0x9),
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
        (bytes32 _mmrRoot, uint64 _mmrSize) = Types.decodeMMRFromRlpHeader(rlpData);
        vm.expectEmit(true, true, true, true, address(l1CrossLayerWitness));
        emit MessageRelayFailed(0, _hash, _mmrSize, _mmrRoot);
        bool success = l1CrossLayerWitness.relayMessage(
            address(0x9),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
        require(!success, "failed");
        require(!l1CrossLayerWitness.successRelayedMessages(_hash), "failed");
    }

    function testRelayMessage() public {
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
        vm.expectEmit(true, true, true, true, address(l1CrossLayerWitness));
        emit MessageRelayed(0, _hash);
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
        require(l1CrossLayerWitness.successRelayedMessages(_hash), "failed");
    }

    function testRelayMessageWithPaused() public {
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
        vm.startPrank(address(addressManager.dao()));
        l1CrossLayerWitness.pause();
        vm.stopPrank();
        vm.startPrank(sender);
        vm.expectRevert(bytes("Pausable: paused"));
        l1CrossLayerWitness.relayMessage(
            address(addressManager),
            sender,
            abi.encodeWithSignature("dao()"),
            0,
            rlpData,
            stateInfo,
            _proof
        );
    }

    function testSendMessage() public {
        vm.startPrank(sender);
        address MockL2Target = address(0x9899);
        bytes memory _message = bytes("0x01");
        uint64 treeSize = _compactMerkleTree.treeSize;
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(MockL2Target, sender, treeSize, _message);
        bytes32 _mmrRoot = _compactMerkleTree.appendLeafHash(_hash);
        vm.expectEmit(true, true, true, true, address(l1CrossLayerWitness));
        emit MessageSent(treeSize, MockL2Target, sender, _mmrRoot, _message);
        l1CrossLayerWitness.sendMessage(MockL2Target, _message);
        require(l1CrossLayerWitness.totalSize() == 1, "size not 1");
    }

    function testSendMessageWithCallerIsThis() public {
        vm.startPrank(address(l1CrossLayerWitness));
        address MockL2Target = address(0x9899);
        bytes memory _message = bytes("0x01");
        vm.expectRevert("wired situation");
        l1CrossLayerWitness.sendMessage(MockL2Target, _message);
    }

    function testSendMessageWithPaused() public {
        vm.startPrank(address(addressManager.dao()));
        l1CrossLayerWitness.pause();
        vm.stopPrank();
        vm.startPrank(sender);
        address MockL2Target = address(0x9899);
        bytes memory _message = bytes("0x01");
        vm.expectRevert("Pausable: paused");
        l1CrossLayerWitness.sendMessage(MockL2Target, _message);
    }

    function testBlockMessage() public {
        vm.startPrank(address(addressManager.dao()));
        bytes32[] memory messageHashes = new bytes32[](2);
        messageHashes[0] = bytes32(uint256(0x0));
        messageHashes[1] = bytes32(uint256(0x1));
        vm.expectEmit(true, true, true, true, address(l1CrossLayerWitness));
        emit MessageBlocked(messageHashes);
        l1CrossLayerWitness.blockMessage(messageHashes);
        require(l1CrossLayerWitness.blockedMessages(bytes32(uint256(0x0))) == true, "failed");
        require(l1CrossLayerWitness.blockedMessages(bytes32(uint256(0x1))) == true, "failed");
    }

    function testBlockMessageWithWrongSender() public {
        bytes32[] memory messageHashes = new bytes32[](1);
        messageHashes[0] = bytes32(uint256(0x0));
        vm.expectRevert("only dao allowed");
        l1CrossLayerWitness.blockMessage(messageHashes);
    }

    function testAllowMessage() public {
        vm.startPrank(address(addressManager.dao()));
        bytes32[] memory messageHashes = new bytes32[](1);
        messageHashes[0] = bytes32(uint256(0x1));
        l1CrossLayerWitness.blockMessage(messageHashes);
        require(l1CrossLayerWitness.blockedMessages(bytes32(uint256(0x1))) == true, "failed");
        vm.stopPrank();
        vm.startPrank(address(addressManager.dao()));
        vm.expectEmit(true, true, true, true, address(l1CrossLayerWitness));
        emit MessageAllowed(messageHashes);
        l1CrossLayerWitness.allowMessage(messageHashes);
        require(l1CrossLayerWitness.blockedMessages(bytes32(uint256(0x0))) == false, "failed");
    }

    function testAllowMessageWithWrongSender() public {
        bytes32[] memory messageHashes = new bytes32[](1);
        messageHashes[0] = bytes32(uint256(0x0));
        vm.expectRevert("only dao allowed");
        l1CrossLayerWitness.allowMessage(messageHashes);
    }
}
