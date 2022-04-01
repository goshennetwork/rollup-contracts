// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/* Library Imports */
import { OVMCodec } from "./OVMCodec.sol";
import "./console.sol";

contract OVMCodecTest {
    //fixme: why failed?
    function testHashTransaction() public {
        OVMCodec.Transaction memory _tx;
        _tx.timestamp = 121212;
        _tx.blockNumber = 10;
        _tx.l1QueueOrigin = OVMCodec.QueueOrigin.SEQUENCER_QUEUE;
        _tx.l1TxOrigin = address(bytes20(hex"1111111111111111111111111111111111111111"));
        _tx.entrypoint = address(bytes20(hex"1111111111111111111111111111111111111111"));
        _tx.gasLimit = 100;
        _tx.data = "0x1234";
        console.logBytes32(OVMCodec.hashTransaction(_tx));
        require(
            OVMCodec.hashTransaction(_tx) == bytes32(0xf07818e2db63d0140e55c9e68cfaa030f9a2d0962f671d6b339edb2207633ebd)
        );
    }
}
