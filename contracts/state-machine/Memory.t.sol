// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

//import "../console.sol";
import "./MachineState.sol";

contract MemoryTest {
    MachineState mem;
    bytes32 root;

    function setUp() public {
        mem = new MachineState();
        root = MerkleTrie.KECCAK256_RLP_NULL_BYTES;
    }

    function testMemUint32() public {
        uint32 ptr = 0;
        uint32 value = 0x04030201;
        root = mem.writeMemory(root, ptr, value);
        require(mem.readMemoryBytes4(root, ptr) == hex"01020304");
        require(mem.readMemory(root, ptr) == value);
        for (uint32 i = 0; i < 4; i++) {
            require(mem.readMemoryByte(root, ptr + i) == bytes1(uint8(i + 1)));
        }
        require(mem.readMemoryBytes2(root, ptr) == bytes2(hex"0102"));
        require(mem.readMemoryBytes2(root, ptr + 2) == bytes2(hex"0304"));
        require(mem.readMemoryBytes32(root, ptr) == bytes32(hex"01020304"));
    }
}
