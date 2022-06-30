// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./RollupInputChain.sol";
import "../test-helper/TestBase.sol";
import "../libraries/console.sol";

contract TestRollupInputChain is TestBase, RollupInputChain {
    address sender = address(0x7777);
    uint64 constant txBaseSize = 213;
    uint64 constant txNumPerBlock = 1;
    uint64 constant BaseTxNum = 50;

    function setUp() public {
        vm.startPrank(sender);
        super.initialize();
        dao.setProposerWhitelist(sender, true);
        dao.setSequencerWhitelist(sender, true);
        feeToken.approve(address(stakingManager), stakingManager.price());
        stakingManager.deposit();
        vm.stopPrank();
    }

    function testAppend1Transfer() public {
        uint64 txNum = 1;
        uint64 batchNum = txNum / txNumPerBlock;
        fakeAppendBatch(batchNum, generateFakeTransferTx(txNum));
    }

    function testAppend10Transfer() public {
        uint64 txNum = 10;
        uint64 batchNum = txNum / txNumPerBlock;
        fakeAppendBatch(batchNum, generateFakeTransferTx(txNum));
    }

    function testAppend50Transfer() public {
        uint64 txNum = 50;
        uint64 batchNum = txNum / txNumPerBlock;
        fakeAppendBatch(batchNum, generateFakeTransferTx(txNum));
    }

    function testAppend100Transfer() public {
        uint64 txNum = 100;
        uint64 batchNum = txNum / txNumPerBlock;
        fakeAppendBatch(batchNum, generateFakeTransferTx(txNum));
    }

    function testAppend200Transfer() public {
        uint64 txNum = 200;
        uint64 batchNum = txNum / txNumPerBlock;
        fakeAppendBatch(batchNum, generateFakeTransferTx(txNum));
    }

    function fakeTx() internal returns (bytes memory) {
        bytes[] memory _rlpList = getRlpList(
            0xffff,
            0xffff,
            0xffff,
            address(0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf),
            ""
        );
        _rlpList[6] = RLPWriter.writeUint(
            115792089237316195423570985008687907852837564279074904382605163141518161494337
        );
        _rlpList[7] = RLPWriter.writeUint(
            115792089237316195423570985008687907852837564279074904382605163141518161494337
        );
        _rlpList[8] = RLPWriter.writeUint(0xffff);
        return RLPWriter.writeList(_rlpList);
    }

    function generateFakeTransferTx(uint64 txNum) internal returns (bytes memory) {
        bytes memory tx = fakeTx();
        uint64 txLen = uint64(tx.length);
        uint64 l = txLen * txNum;
        bytes memory b = new bytes(l);
        uint64 ptr;
        assembly {
            ptr := add(b, 32)
        }
        for (uint64 i = 0; i < txNum; i++) {
            uint64 txPtr;
            assembly {
                txPtr := add(tx, 32)
            }
            for (uint64 offset = 0; offset < txLen; offset += 32) {
                assembly {
                    mstore(add(ptr, offset), mload(add(txPtr, offset)))
                }
            }
            ptr += txLen;
        }
        return b;
    }

    function fakeAppendBatch(uint64 batchNum, bytes memory data) internal {
        // now support at least one sub batch
        require(batchNum > 0, "no sub batch");
        uint64 batchIndex = 0;
        uint64 queueNum = 0;
        uint64 pendingQueueIndex = 0;
        uint64 subBatchNum = batchNum;
        uint64 time0Start = 0;
        uint32[] memory timeDiff = new uint32[](batchNum - 1);
        bytes memory _info = abi.encodePacked(
            batchIndex,
            queueNum,
            pendingQueueIndex,
            subBatchNum,
            time0Start,
            timeDiff,
            data
        );
        vm.startPrank(sender, sender);
        (bool success, ) = address(rollupInputChain).call(
            abi.encodePacked(abi.encodeWithSignature("appendBatch()"), _info)
        );
        require(success, "failed");
        vm.stopPrank();
    }
}
