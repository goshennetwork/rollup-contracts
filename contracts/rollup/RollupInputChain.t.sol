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
import "../libraries/RLPWriter.sol";
import "../libraries/UnsafeSign.sol";

contract TestRollupInputChain is TestBase, RollupInputChain {
    address testAddress = address(0x8888); //admain
    address testAddress2 = address(0x9999);

    function setUp() public {
        vm.startPrank(testAddress);
        initialize();
        dao.setProposerWhitelist(testAddress, true);
        dao.setSequencerWhitelist(testAddress, true);
        feeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        vm.stopPrank();
    }

    //Test Enqueue()
    /*1.Test Fail*/

    //test  if{}  when (msg.sender == Constants.L1_CROSS_LAYER_WITNESS) , revert:"malicious sender"
    function testEnqueueWithWitnessSender() public {
        address l1witness2 = Constants.L1_CROSS_LAYER_WITNESS;
        uint64 pendingNonce = rollupInputChain.getNonceByAddress(l1witness2);
        uint64 passV = 36 + 2 * rollupInputChain.l2ChainID();
        uint64 passGasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10;

        vm.startPrank(l1witness2, l1witness2);
        vm.expectRevert("malicious sender");
        rollupInputChain.enqueue(address(1), passGasLimit, bytes("0x100"), pendingNonce, 1, 1, passV);
        vm.stopPrank();
    }

    //test if{}  when (_nonce != pendingNonce) revert("wrong nonce")
    function testEnqueueWithWrongNonce() public {
        uint64 passV = 36 + 2 * rollupInputChain.l2ChainID();
        uint64 passGasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10;
        uint64 pendingNonce = rollupInputChain.getNonceByAddress(testAddress) + 1;

        vm.startPrank(testAddress, testAddress);
        vm.expectRevert("wrong nonce");
        rollupInputChain.enqueue(address(1), passGasLimit, bytes("0x0"), pendingNonce, 1, 1, passV);
        vm.stopPrank();
    }

    //test if-else： else{}: when (msg.sender != l1CrossLayerWitness)
    //                      revert("contract can not enqueue L2 Tx")
    function testEnqueueWithContractSender() public {
        uint64 passV = 36 + 2 * rollupInputChain.l2ChainID();
        uint64 passGasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10;
        uint64 pendingNonce = rollupInputChain.getNonceByAddress(testAddress);

        vm.startPrank(testAddress);
        vm.expectRevert("contract can not enqueue L2 Tx");
        rollupInputChain.enqueue(address(1), passGasLimit, bytes("0x0"), pendingNonce, 1, 1, passV);
        vm.stopPrank();
    }

    //test gasLimit toohigh
    // when (_gasLimit > maxEnqueueTxGasLimit)  revert("too high Tx gas limit")
    function testEnqueueWithTooHighGasLimit() public {
        uint64 passV = 36 + 2 * rollupInputChain.l2ChainID();
        uint64 allowedNonce = rollupInputChain.getNonceByAddress(testAddress);
        uint64 highGasLimit = rollupInputChain.maxEnqueueTxGasLimit() + uint64(100);

        vm.startPrank(testAddress, testAddress);
        vm.expectRevert("too high Tx gas limit");
        rollupInputChain.enqueue(address(1), highGasLimit, bytes("0x0"), allowedNonce, 1, 1, passV);
        vm.stopPrank();
    }

    //test gasLimit toolow
    // when (_gasLimit < MIN_ENQUEUE_TX_GAS )  revert("too low Tx gas limit")
    function testEnqueueWithToolowGasLimit() public {
        uint64 passV = 36 + 2 * rollupInputChain.l2ChainID();
        uint64 allowedNonce = rollupInputChain.getNonceByAddress(testAddress);
        uint64 mingasLimit = uint64(MIN_ENQUEUE_TX_GAS) - 10;

        vm.startPrank(testAddress, testAddress);
        vm.expectRevert("too low Tx gas limit");
        rollupInputChain.enqueue(address(1), mingasLimit, bytes("0x0"), allowedNonce, 1, 1, passV);
        vm.stopPrank();
    }

    //test v of (r,s,v)
    // when (_pureV > 28)  revert ("invalid v")
    // invalidV must (> 36 + 2*l2ChainID )
    function testEnqueueWithinvalidV() public {
        uint64 gasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10;
        uint64 pendingNonce = getNonceByAddress(testAddress);
        uint64 invalidV = 37 + 2 * (rollupInputChain.l2ChainID());

        vm.startPrank(testAddress, testAddress);
        vm.expectRevert("invalid v");
        rollupInputChain.enqueue(address(1), gasLimit, bytes("0x0"), pendingNonce, 1, 1, invalidV);
        vm.stopPrank();
    }

    //test wrong sign
    // when  sender != ecrecover(_signTxHash, uint8(_pureV), bytes32(r), bytes32(s)
    // revert ("wrong sign")
    function testEnqueueWithInvalidSign() public {
        uint64 gasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10;
        uint64 pendingNonce = getNonceByAddress(testAddress);
        uint64 allowedV = 36 + 2 * (rollupInputChain.l2ChainID());

        vm.startPrank(testAddress, testAddress);
        vm.expectRevert("wrong sign");
        rollupInputChain.enqueue(address(1), gasLimit, bytes("0x0"), pendingNonce, 1, 1, allowedV);
        vm.stopPrank();
    }

    //Test Fail() if{msg.sender == tx.origin}
    // when _rlpTx.length > _maxTxSize  ; revert ("too large tx data size")
    function testEnqueueWithTooLargeData() public {
        uint64 pendingNonce = 1 << 63; //when initialize ; pendingNonce always 2^63
        uint64 gasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10;
        uint64 allowedV = 36 + 2 * (rollupInputChain.l2ChainID());
        bytes memory data = new bytes(50002);
        address a = 0x576Dacb2e7Cb8DADbd9665CA9e62107AdD049EB0; //sign address
        uint256 r = 1;
        uint256 s = 1;

        vm.startPrank(a, a);
        vm.expectRevert("too large tx data size");
        rollupInputChain.enqueue(address(1), gasLimit, data, pendingNonce, r, s, allowedV);
        vm.stopPrank();
    }

    // Test Enqueue()
    /*2.Test pass*/

    // if{msg.sender == tx.origin}  enqueue tx*2
    // test queue.length  &&  queue.context{timestamp + transactionHash} && emit event
    function testSenderEnqueueTwoTx() public {
        address SENDER = UnsafeSign.G2ADDR;
        bytes32 Rlptx0;
        bytes32 Rlptx1;
        //enqueue
        vm.startPrank(SENDER, SENDER);
        Rlptx0 = enqueue1(bytes("0x0"));

        //test queue length
        require(rollupInputChain.totalQueue() == 1, "enqueue failed");
        Rlptx1 = enqueue1(bytes("0x1"));
        require(rollupInputChain.totalQueue() == 2, "enqueue failed");
        //test queue context
        uint64 b;
        bytes32 aa;
        (aa, b) = rollupInputChain.getQueueTxInfo(0);
        require(b == uint64(block.timestamp), "queue[0] storage timestamp different");
        require(Rlptx0 == aa, "queue[0] storage transactionHash different");
        (aa, b) = rollupInputChain.getQueueTxInfo(1);
        require(b == uint64(block.timestamp), "queue[1] storage timestamp different");
        require(Rlptx1 == aa, "queue[1] storage transactionHash different");
        vm.stopPrank();
    }

    //helper function1  & test emit event
    // when if-else: if{}: {msg.sender == tx.origin}
    // input data , enqueue rollupInputChian
    function enqueue1(bytes memory _data) public returns (bytes32) {
        address SENDER = UnsafeSign.G2ADDR;
        uint64 pendingNonce = rollupInputChain.getNonceByAddress(SENDER);
        uint64 gasLimit = rollupInputChain.maxEnqueueTxGasLimit() / 3;
        //state same as rollupInputChain
        l2ChainID = rollupInputChain.l2ChainID();

        // create allowed  sign enqueue
        //data0
        bytes[] memory allowedList0 = getRlpList(pendingNonce, gasLimit, GAS_PRICE, address(1), _data);
        bytes32 signTxHash0 = keccak256(RLPWriter.writeList(allowedList0));
        uint256 r;
        uint256 s;
        uint64 v;
        (r, s, v) = UnsafeSign.Sign2(signTxHash0, rollupInputChain.l2ChainID());
        // uint pureV = v - 2 * rollupInputChain.l2ChainID() - 8;
        // address signAddress = ecrecover(signTxHash0, uint8(pureV), bytes32(r), bytes32(s));
        allowedList0[6] = RLPWriter.writeUint(v);
        allowedList0[7] = RLPWriter.writeUint(r);
        allowedList0[8] = RLPWriter.writeUint(s);
        bytes memory rlpTx = RLPWriter.writeList(allowedList0);

        //test eventEmit
        vm.expectEmit(true, true, true, true);
        emit TransactionEnqueued(
            uint64(rollupInputChain.totalQueue()),
            SENDER,
            address(1),
            rlpTx,
            uint64(block.timestamp)
        );
        //enqueue rollupInputChian
        rollupInputChain.enqueue(address(1), gasLimit, _data, pendingNonce, r, s, v);
        return keccak256(rlpTx);
    }

    // Test L1witness Enqueue()
    // if-else: else {}:  {msg.sender == l1CrossLayerWitness()}  enqueue tx*2
    // test queue.length  &&  queue.context{timestamp + transactionHash} && emit event
    function testL1witnessEnqueueTwoTx() public {
        bytes32 Rlptx0;
        bytes32 Rlptx1;
        vm.startPrank(address(l1CrossLayerWitness), testAddress);
        //enqueue
        Rlptx0 = enqueue2(bytes("0x0"));
        //test queue length
        require(rollupInputChain.totalQueue() == 1, "enqueue failed");
        Rlptx1 = enqueue2(bytes("0x1"));
        require(rollupInputChain.totalQueue() == 2, "enqueue failed");
        //test queue context
        uint64 b;
        bytes32 aa;
        (aa, b) = rollupInputChain.getQueueTxInfo(0);
        require(b == uint64(block.timestamp), "queue[0] storage timestamp different");
        require(Rlptx0 == aa, "queue[0] storage transactionHash different");
        (aa, b) = rollupInputChain.getQueueTxInfo(1);
        require(b == uint64(block.timestamp), "queue[1] storage timestamp different");
        require(Rlptx1 == aa, "queue[1] storage transactionHash different");
        vm.stopPrank();
    }

    //helper function2  + test emit event
    //  if-else: else{} : {msg.sender == l1CrossLayerWitness()}
    // input data , enqueue rollupInputChian
    function enqueue2(bytes memory _data) public returns (bytes32) {
        address sender = Constants.L1_CROSS_LAYER_WITNESS;
        uint64 pendingNonce = rollupInputChain.getNonceByAddress(sender);
        uint64 GasLimit = rollupInputChain.maxWitnessTxExecGasLimit();
        uint256 Gasprice = 0;
        //state same as rollupInputChain
        l2ChainID = rollupInputChain.l2ChainID();
        // create allowed  sign enqueue
        //data0
        bytes[] memory allowedList0 = getRlpList(pendingNonce, GasLimit, Gasprice, address(1), _data);
        bytes32 signTxHash0 = keccak256(RLPWriter.writeList(allowedList0));
        uint256 r;
        uint256 s;
        uint64 v;
        (r, s, v) = UnsafeSign.Sign(signTxHash0, rollupInputChain.l2ChainID());
        allowedList0[6] = RLPWriter.writeUint(v);
        allowedList0[7] = RLPWriter.writeUint(r);
        allowedList0[8] = RLPWriter.writeUint(s);
        bytes memory rlpTx = RLPWriter.writeList(allowedList0);
        //test eventEmit
        vm.expectEmit(true, true, true, true);
        emit TransactionEnqueued(
            uint64(rollupInputChain.totalQueue()),
            sender,
            address(1),
            rlpTx,
            uint64(block.timestamp)
        );
        //enqueue rollupInputChian
        rollupInputChain.enqueue(address(1), GasLimit, _data, pendingNonce, r, s, v);
        return keccak256(rlpTx);
    }

    //Test appendBatch
    /**1.Test Fail**/

    //test Fail msg.sender
    // when  sender != sequencerWhitelist
    // revert ("only sequencer")
    function testAppendBatchNotSequencer() public {
        vm.startPrank(testAddress2, testAddress2);
        vm.expectRevert("only sequencer");
        rollupInputChain.appendBatch();
        vm.stopPrank();
    }

    //test Fail sequencer staking
    // when  isStaking(msg.sender)  ==  false
    // revert ("Sequencer should be staking")
    function testAppendBatchSequencerNoStaking() public {
        vm.startPrank(testAddress);
        dao.setSequencerWhitelist(testAddress2, true);
        vm.stopPrank();
        vm.startPrank(testAddress2);
        vm.expectRevert("Sequencer should be staking");
        rollupInputChain.appendBatch();
        vm.stopPrank();
    }

    //test Fail msg.data.length
    // when  msg.data.length < 36
    // revert ("wrong len")
    function testAppendWrongBatchDataLength() public {
        vm.startPrank(testAddress);
        vm.expectRevert("wrong len");
        helpCall(address(rollupInputChain), bytes("0x0"));
        vm.stopPrank();
    }

    function helpCall(address _rollupInputChain, bytes memory _data) public {
        (bool success, ) = _rollupInputChain.call(abi.encodePacked(abi.encodeWithSignature("appendBatch()"), _data));
        require(success, "call failed");
    }

    //test Fail _batchIndex
    // when  _batchIndex != chainHeight()
    // revert ("wrong batch index")
    function testAppendBatchWrongBatchIndex() public {
        vm.startPrank(testAddress);
        uint64 invalidIndex = rollupInputChain.chainHeight() + 10;
        vm.expectRevert("wrong batch index");
        fakeAppendBatch(invalidIndex, 0, 1, 1, 0, bytes("0x0"));
        vm.stopPrank();
    }

    //test Fail _queueStartIndex
    // when  _queueStartIndex != pendingQueueIndex
    // revert ("incorrect pending queue index")
    function testAppendBatchQueueStartIndex() public {
        vm.startPrank(testAddress);
        vm.expectRevert("incorrect pending queue index");
        fakeAppendBatch(0, 0, 1, 1, 0, bytes("0x0")); //it will always work
        vm.stopPrank();
    }

    //helper function: encode calldata to appendBatch()
    function fakeAppendBatch(
        uint64 _batchIndex,
        uint64 _queueNum,
        uint64 _queueStartIndex,
        uint64 batchNum,
        uint64 _time0Start,
        bytes memory data
    ) internal {
        // now support at least one sub batch
        uint64 batchIndex = _batchIndex;
        uint64 queueNum = _queueNum;
        uint64 pendingQueueIndex = _queueStartIndex;
        uint64 subBatchNum = batchNum; //all batch num
        uint64 time0Start = _time0Start;
        uint32[] memory timeDiff;
        bytes memory _info;
        if (batchNum == 0) {
            timeDiff = new uint32[](0);
            _info = abi.encodePacked(batchIndex, queueNum, pendingQueueIndex, subBatchNum);
        } else {
            timeDiff = new uint32[](batchNum - 1); //alwayes length = batchNum -1
            _info = abi.encodePacked(batchIndex, queueNum, pendingQueueIndex, subBatchNum, time0Start, timeDiff, data);
        }
        (bool success, ) = address(rollupInputChain).call(
            abi.encodePacked(abi.encodeWithSignature("appendBatch()"), _info)
        );
        require(success, "failed");
    }

    //test Fail _nextPendingQueueIndex
    // when _nextPendingQueueIndex(queueStart + queueNum) > queuedTxInfos.length
    // revert ("attempt to append unavailable queue")
    function testAppendBatchInvalidNextPendingQueueIndex() public {
        vm.startPrank(testAddress);
        vm.expectRevert("attempt to append unavailable queue");
        fakeAppendBatch(0, 1, 0, 1, 0, bytes("0x0"));
        vm.stopPrank();
    }

    //test Fail _queueNum
    // when _queueNum <= 0
    // revert ("nothing to append")
    function testAppendBatchInvalidQueueNum() public {
        vm.startPrank(testAddress);
        vm.expectRevert("nothing to append");
        fakeAppendBatch(0, 0, 0, 0, 0, bytes("0x0"));
        vm.stopPrank();
    }

    //test Fail msg.data.length
    // when msg.data.length != _batchDataPos
    // revert ("wrong calldata")
    function testAppendBatchInvalidMsgdataLength() public {
        //enqueue
        vm.startPrank(address(l1CrossLayerWitness), testAddress);
        bytes32 Rlptx0 = enqueue2(bytes("0x0"));
        vm.stopPrank();

        vm.startPrank(testAddress);
        bytes memory data = new bytes(502); //because _batchDataPos always change, 502 will cover almost all
        uint64 batchIndex = 0;
        uint64 queueNum = 1;
        uint64 pendingQueueIndex = 0;
        uint64 subBatchNum = 0;
        uint64 time0Start = 0;
        uint32[] memory timeDiff = new uint32[](0);
        bytes memory _info = abi.encodePacked(
            batchIndex,
            queueNum,
            pendingQueueIndex,
            subBatchNum,
            time0Start,
            timeDiff,
            data
        );
        vm.expectRevert("wrong calldata");
        (bool success, ) = address(rollupInputChain).call(
            abi.encodePacked(abi.encodeWithSignature("appendBatch()"), _info)
        );
        require(success, "failed");
        vm.stopPrank();
    }

    //test Fail _timestamp
    // when _timestamp < lastTimestamp
    // revert ("wrong batch timestamp")
    //---- we do : enqueue*2 -->  AppendBatch(queue(1)) ---> lastTimestamp = 2
    //--->AppendBatch(queue(1)) & timeStamp == 0 ---> ("wrong batch timestamp")
    function testAppendBatchInvalidTimestamp() public {
        //enqueue
        vm.startPrank(address(l1CrossLayerWitness), testAddress);
        vm.warp(2);
        bytes32 Rlptx0 = enqueue2(bytes("0x0"));
        bytes32 Rlptx2 = enqueue2(bytes("0x0"));
        vm.stopPrank();

        vm.startPrank(testAddress);
        fakeAppendBatch(0, 1, 0, 1, 0, bytes("0x0"));
        vm.expectRevert("wrong batch timestamp");
        fakeAppendBatch(1, 1, 1, 1, 0, bytes("0x0")); //This situation can work in many occasions
        vm.stopPrank();
    }

    //test Fail _timestamp
    // when _timestamp > _nextTimestamp
    // revert ("last batch timestamp too high")
    function testAppendBatchTooHighTimestamp() public {
        vm.startPrank(testAddress);
        uint64 time = uint64(block.timestamp) + 10;
        vm.expectRevert("last batch timestamp too high");
        fakeAppendBatch(0, 0, 0, 1, time, bytes("0x0"));
        vm.stopPrank();
    }

    // Test AppendBatch()
    /*2.Test pass*/

    //test event AppendBatch if batchNum == 0
    function testAppendBatchIfBatchnumEqual0() public {
        //enqueue rollupInputChain Contract
        vm.startPrank(address(l1CrossLayerWitness), testAddress);
        bytes32 Rlptx0 = enqueue2(bytes("0x0"));
        bytes32 Rlptx1 = enqueue2(bytes("0x0"));
        vm.stopPrank();

        vm.startPrank(testAddress);
        bytes32 _queueHashes = getrollupInputChainQueueHash(0, 1);
        bytes memory info = getinfo(0, 1, 0, 0, 0, "");
        bytes32 inputhash = keccak256(abi.encodePacked(keccak256(info), _queueHashes));
        //test eventEmit
        vm.expectEmit(true, true, false, true);
        emit TransactionAppended(testAddress, 0, 0, 1, inputhash);
        fakeAppendBatch(0, 1, 0, 0, 0, "");
        vm.stopPrank();
    }

    // help function : get queuehash in rollupInputChain-Contract
    function getrollupInputChainQueueHash(uint64 _queueStartIndex, uint64 _queueNum) internal view returns (bytes32) {
        uint256 len = (32 + 8) * _queueNum;
        bytes memory _queueHash = new bytes(len);
        uint64 _offset = 0;
        for (uint64 i = 0; i < _queueNum; i++) {
            (bytes32 queuehash, uint64 queuetime) = rollupInputChain.getQueueTxInfo(_queueStartIndex + i);
            QueueTxInfo memory info = QueueTxInfo(queuehash, queuetime);
            bytes32 txHash = info.transactionHash;
            bytes32 time = bytes32(uint256(info.timestamp) << 192);
            assembly {
                let ptr := add(_queueHash, _offset)
                mstore(ptr, txHash)
                ptr := add(ptr, 32)
                // @notice we reuse _queueHash's the first 32 byte length bits, so no overflow
                mstore(ptr, time)
            }
            _offset += 40;
        }
        // @notice we reuse _queueHash's length, so can not use keccak256(_queueHash)
        bytes32 result;
        assembly {
            result := keccak256(_queueHash, len)
        }
        return result;
    }

    //helper function: get Inputhash
    //Only in this case, the function name is not required
    function getinfo(
        uint64 _batchIndex,
        uint64 _queueNum,
        uint64 _queueStartIndex,
        uint64 batchNum,
        uint64 _time0Start,
        bytes memory data
    ) internal returns (bytes memory) {
        uint64 batchIndex = _batchIndex;
        uint64 queueNum = _queueNum;
        uint64 pendingQueueIndex = _queueStartIndex;
        uint64 subBatchNum = batchNum;
        uint64 time0Start = _time0Start;
        uint32[] memory timeDiff;
        bytes memory _info;
        if (batchNum == 0) {
            timeDiff = new uint32[](0);
            _info = abi.encodePacked(queueNum, pendingQueueIndex, subBatchNum);
        } else {
            timeDiff = new uint32[](batchNum - 1);
            _info = abi.encodePacked(queueNum, pendingQueueIndex, subBatchNum, time0Start, timeDiff, data);
        }
        return _info;
    }

    // test event AppendBatch if batchNum != 0
    // appendBatch*2
    function testAppendBatchIfBatchnumBiggerThan0() public {
        //enqueue
        vm.startPrank(address(l1CrossLayerWitness), testAddress);
        vm.warp(2);
        bytes32 Rlptx0 = enqueue2(bytes("0x0"));
        bytes32 Rlptx2 = enqueue2(bytes("0x0"));
        vm.stopPrank();

        vm.startPrank(testAddress);
        fakeAppendBatch(0, 1, 0, 1, 0, bytes("0x0"));

        bytes32 _queueHashes = getrollupInputChainQueueHash(1, 1);
        bytes memory info = getinfo(1, 1, 1, 1, 2, bytes("0x0"));
        bytes32 inputhash = keccak256(abi.encodePacked(keccak256(info), _queueHashes));
        vm.expectEmit(true, true, false, true);
        emit TransactionAppended(testAddress, 1, 1, 1, inputhash);
        fakeAppendBatch(1, 1, 1, 1, 2, bytes("0x0")); //it will always work
        vm.stopPrank();
    }
}
