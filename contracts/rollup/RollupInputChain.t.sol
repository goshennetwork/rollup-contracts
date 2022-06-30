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


contract TestRollupInputChain is TestBase,RollupInputChain {
    address testAddress = address(0x8888);
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
    function testFail_if_L1Witness_Enqueue() public {
        // address l1witness = addressManager.resolve(AddressName.L1_CROSS_LAYER_WITNESS);
        address l1witness2 = 0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf ;
        vm.startPrank(l1witness2,l1witness2);
        rollupInputChain.enqueue(address(1),10000000,bytes("0x100"),100,1,1,1);
        vm.stopPrank();
    }


    //test if{}  when (_nonce != pendingNonce) revert("wrong nonce")
    function testFail_if_pendingNonce_Enqueue() public {
        vm.startPrank(testAddress,testAddress);
        uint64 pendingNonce = rollupInputChain.getNonceByAddress(testAddress) + 1;
        rollupInputChain.enqueue(address(1),10000000,bytes("0x0"),pendingNonce,1,1,3000);
        vm.stopPrank();
    }

    //test if-else： else{}: when (msg.sender != l1CrossLayerWitness) 
    //                      revert("contract can not enqueue L2 Tx")
    function testFail_else_Nol1CrossLayerWitness_Enqueue() public {        
        vm.startPrank(testAddress);
        rollupInputChain.enqueue(address(1),10000000,bytes("0x0"),100,1,1,3000);
        vm.stopPrank();
    }

    //test gasLimit toohigh
    // when (_gasLimit > maxEnqueueTxGasLimit)  revert("too high Tx gas limit")         
    function testFail_gasLimitTooHigh_Enqueue() public {        
        uint64 highGasLimit = rollupInputChain.maxEnqueueTxGasLimit() + uint64(100) ;
        uint64 allowedNonce = getNonceByAddress(testAddress);
        vm.startPrank(testAddress,testAddress);
        rollupInputChain.enqueue(address(1),highGasLimit,bytes("0x0"),allowedNonce,1,1,3000);
        vm.stopPrank();
    }

    //test gasLimit toolow
    // when (_gasLimit < MIN_ENQUEUE_TX_GAS )  revert("too low Tx gas limit")         
    function testFail_gasLimitToolow_Enqueue() public {        
        uint64 mingasLimit = uint64(MIN_ENQUEUE_TX_GAS) - 10 ;
        uint64 pendingNonce = getNonceByAddress(testAddress);
        vm.startPrank(testAddress,testAddress);
        rollupInputChain.enqueue(address(1),mingasLimit,bytes("0x0"),pendingNonce,1,1,3000);
        vm.stopPrank();
    }

    //test v of (r,s,v)
    // when (_pureV > 28)  revert ("invalid v")
    // invalidV must (> 36 + 2*l2ChainID )
    function testFail_invalidV_Enqueue() public {        
        uint64 gasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10 ;
        uint64 pendingNonce = getNonceByAddress(testAddress);
        uint64 invalidV = 37 + 2 *(rollupInputChain.l2ChainID());
        vm.startPrank(testAddress,testAddress);
        rollupInputChain.enqueue(address(1),gasLimit,bytes("0x0"),pendingNonce,1,1,invalidV);
        vm.stopPrank();
    }

    //test wrong sign
    // when  sender != ecrecover(_signTxHash, uint8(_pureV), bytes32(r), bytes32(s)
    // revert ("wrong sign")
    function testFail_invalidSign_Enqueue() public {        
        uint64 gasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10 ;
        uint64 pendingNonce = getNonceByAddress(testAddress);
        uint64 allowedV = 36 + 2 *(rollupInputChain.l2ChainID());
        vm.startPrank(testAddress,testAddress);
        rollupInputChain.enqueue(address(1),gasLimit,bytes("0x0"),pendingNonce,1,1,allowedV);
        vm.stopPrank();
    }

    //Test Fail() if{msg.sender == tx.origin}  
    // when _rlpTx.length > _maxTxSize  ; revert ("too large tx data size")
    function testFail_dataTooLarge_Enqueue() public {
        uint64 pendingNonce = 1 << 63;//when initialize ; pendingNonce always 2^63
        uint64 gasLimit = uint64(MIN_ENQUEUE_TX_GAS) + 10 ;
        uint64 allowedV = 36 + 2 *(rollupInputChain.l2ChainID());

        bytes memory data = new bytes(50002);
        address a = 0x576Dacb2e7Cb8DADbd9665CA9e62107AdD049EB0 ;
        uint r = 1;
        uint s = 1;
        vm.startPrank(a,a);
        rollupInputChain.enqueue(address(1),gasLimit,data,pendingNonce,r,s,allowedV);
        vm.stopPrank();
    }

// Test Enqueue() 
/*2.Test pass*/
    
    // if{msg.sender == tx.origin}  enqueue tx*2  
    // test queue.length  &&  queue.context{timestamp + transactionHash} && emit event
    function test_if_Enqueue() public {
        address SENDER = UnsafeSign.G2ADDR ;
        bytes32  Rlptx0 ;
        bytes32  Rlptx1 ;
        //enqueue
        vm.startPrank(SENDER,SENDER);
        Rlptx0 = enqueue1(bytes("0x0"));

        //test queue length
        require(rollupInputChain.totalQueue() == 1 , "enqueue failed");
        Rlptx1 = enqueue1(bytes("0x1"));
        require(rollupInputChain.totalQueue() == 2 , "enqueue failed");
        //test queue context
        uint64 b ;
        bytes32 aa ;
        (aa , b) = rollupInputChain.getQueueTxInfo(0);
        require(b == uint64(block.timestamp) , "queue[0] storage timestamp different");
        require(Rlptx0 == aa , "queue[0] storage transactionHash different");
        (aa , b) = rollupInputChain.getQueueTxInfo(1);
        require(b == uint64(block.timestamp) , "queue[1] storage timestamp different");
        require(Rlptx1 == aa , "queue[1] storage transactionHash different");
        vm.stopPrank();
    }
    //helper function1  + test emit event
    // when if-else: if{}: {msg.sender == tx.origin}
    // input data , enqueue rollupInputChian
    function enqueue1(bytes memory _data) public returns(bytes32  ) {
        address SENDER = UnsafeSign.G2ADDR ;
        uint64 pendingNonce  = rollupInputChain.getNonceByAddress(SENDER);
        uint64 gasLimit = rollupInputChain.maxEnqueueTxGasLimit()/3;
        // create allowed  sign enqueue
        //data0
        bytes[] memory allowedList0 = rollupInputChain.getRlpList(pendingNonce, gasLimit, GAS_PRICE, address(1), _data);
        bytes32 signTxHash0 = keccak256(RLPWriter.writeList(allowedList0));
        uint r ;
        uint s ;
        uint64 v ;
        ( r ,  s ,  v) = UnsafeSign.Sign2(signTxHash0, rollupInputChain.l2ChainID());
        // uint pureV = v - 2 * rollupInputChain.l2ChainID() - 8;
        // address signAddress = ecrecover(signTxHash0, uint8(pureV), bytes32(r), bytes32(s));
        allowedList0[6] = RLPWriter.writeUint(v);
        allowedList0[7] = RLPWriter.writeUint(r);
        allowedList0[8] = RLPWriter.writeUint(s);
        bytes memory rlpTx = RLPWriter.writeList(allowedList0);

        //test eventEmit
        vm.expectEmit(true , true , true , true);
        emit TransactionEnqueued(uint64(rollupInputChain.totalQueue()),SENDER,address(1),rlpTx,uint64(block.timestamp));
        //enqueue rollupInputChian
        rollupInputChain.enqueue(address(1),gasLimit,_data,pendingNonce,r,s,v);
        return keccak256(rlpTx) ;
    }

    // Test Enqueue() 
    // if-else: else {}:  {msg.sender == l1CrossLayerWitness()}  enqueue tx*2  
    // test queue.length  &&  queue.context{timestamp + transactionHash} && emit event
    function test_else_Enqueue() public {
        bytes32  Rlptx0 ;
        bytes32  Rlptx1 ;
        vm.startPrank(address(l1CrossLayerWitness) ,testAddress);
        //enqueue
        Rlptx0 = enqueue2(bytes("0x0"));
        //test queue length
        require(rollupInputChain.totalQueue() == 1 , "enqueue failed");
        Rlptx1 = enqueue2(bytes("0x1"));
        require(rollupInputChain.totalQueue() == 2 , "enqueue failed");
        //test queue context
        uint64 b ;
        bytes32 aa ;
        (aa , b) = rollupInputChain.getQueueTxInfo(0);
        require(b == uint64(block.timestamp) , "queue[0] storage timestamp different");
        require(Rlptx0 == aa , "queue[0] storage transactionHash different");
        (aa , b) = rollupInputChain.getQueueTxInfo(1);
        require(b == uint64(block.timestamp) , "queue[1] storage timestamp different");
        require(Rlptx1 == aa , "queue[1] storage transactionHash different");
        vm.stopPrank();
    }

    //helper function2  + test emit event
    //  if-else: else{} : {msg.sender == l1CrossLayerWitness()}
    // input data , enqueue rollupInputChian
    function enqueue2(bytes memory _data) public returns(bytes32) {
        address sender = Constants.L1_CROSS_LAYER_WITNESS;
        // console.log("sender");
        // console.log(sender);
        uint64 pendingNonce  = rollupInputChain.getNonceByAddress(sender);
        uint64 GasLimit = rollupInputChain.maxWitnessTxExecGasLimit() + uint64(16 * INTRINSIC_GAS_FACTOR * _data.length);
        uint256 Gasprice = 0 ;
        // create allowed  sign enqueue
        //data0
        bytes[] memory allowedList0 = rollupInputChain.getRlpList(pendingNonce, GasLimit, Gasprice, address(1), _data);
        bytes32 signTxHash0 = keccak256(RLPWriter.writeList(allowedList0));
        uint r ;
        uint s ;
        uint64 v ;
        (r, s, v) = UnsafeSign.Sign(signTxHash0, rollupInputChain.l2ChainID());
        allowedList0[6] = RLPWriter.writeUint(v);
        allowedList0[7] = RLPWriter.writeUint(r);
        allowedList0[8] = RLPWriter.writeUint(s);
        bytes memory rlpTx = RLPWriter.writeList(allowedList0);
        //test eventEmit
        vm.expectEmit(true , true , true , true);
        emit TransactionEnqueued(uint64(rollupInputChain.totalQueue()),sender,address(1),rlpTx,uint64(block.timestamp));
        //enqueue rollupInputChian
        rollupInputChain.enqueue(address(1),GasLimit,_data,pendingNonce,r,s,v);
        return keccak256(rlpTx) ;
    }

//Test appendBatch
/**1.Test Fail**/

    //test msg.sender
    // when  sender != sequencerWhitelist
    // revert ("only sequencer")
    function testFail_NotSequencer_appendBatch() public {        
        vm.startPrank(testAddress2,testAddress2);
        rollupInputChain.appendBatch();
        vm.stopPrank();
    }

    //test sequencer staking
    // when  isStaking(msg.sender)  ==  false
    // revert ("Sequencer should be staking")
    function testFail_NoStaking_appendBatch() public {        
        vm.startPrank(testAddress);
        dao.setSequencerWhitelist(testAddress2, true);
        vm.stopPrank();
        vm.startPrank(testAddress2);
        rollupInputChain.appendBatch();
        vm.stopPrank();
    }


    //test msg.data.length
    // when  msg.data.length < 36
    // revert ("wrong len")
    function testFail_dataLength_appendBatch() public {        
        vm.startPrank(testAddress);
        helpCall(address(rollupInputChain));
        vm.stopPrank();
    }
    function helpCall(address _rollupInputChain) public {
        (bool success, )= _rollupInputChain.call(
            abi.encodePacked(abi.encodeWithSignature("appendBatch()"),bytes("0x0"))
        ) ;
        require(success , "call failed") ;
    }
    

    //test _batchIndex
    // when  _batchIndex != chainHeight()
    // revert ("wrong batch index")
    function testFail_batchIndex_appendBatch() public {        
        vm.startPrank(testAddress);
        helpCall2(address(rollupInputChain));
        vm.stopPrank();
    }
    function helpCall2(address _rollupInputChain) public {
        (bool success, )= _rollupInputChain.call(
            abi.encodePacked(abi.encodeWithSignature("appendBatch()"),bytes("0x000000010000000000000000000000000000"))
        ) ;
        require(success , "call failed") ;
    }




}
