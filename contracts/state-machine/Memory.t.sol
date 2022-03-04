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

        ptr += 4;
        for (uint32 i = 0; i < 4; i++) {
            root = mem.writeMemoryByte(root, ptr + i, bytes1(uint8(i + 1)));
        }
        require(mem.readMemoryBytes4(root, ptr) == hex"01020304");
        require(mem.readMemory(root, ptr) == value);

        ptr += 4;
        root = mem.writeMemoryBytes2(root, ptr, bytes2(hex"0102"));
        root = mem.writeMemoryBytes2(root, ptr + 2, bytes2(hex"0304"));
        require(mem.readMemoryBytes4(root, ptr) == hex"01020304");

        ptr = 0;
        root = mem.writeMemoryBytes32(root, ptr, bytes32(hex"dead1234567879"));
        require(mem.readMemoryBytes32(root, ptr) == bytes32(hex"dead1234567879"));
    }

    function testMemoryString() public {
        uint32 ptr = 0;
        root = mem.writeMemoryByte(root, ptr, bytes1("a"));
        require(keccak256(bytes(mem.readMemoryString(root, ptr, 1))) == keccak256("a"));

        ptr = 4;
        root = mem.writeMemoryBytes2(root, ptr, bytes2("hi"));
        require(keccak256(bytes(mem.readMemoryString(root, ptr, 2))) == keccak256("hi"));

        ptr = 8;
        root = mem.writeMemoryBytes4(root, ptr, bytes4("hii"));
        require(keccak256(bytes(mem.readMemoryString(root, ptr, 3))) == keccak256("hii"));

        ptr = 12;
        root = mem.writeMemoryBytes4(root, ptr, bytes4("hell"));
        require(keccak256(bytes(mem.readMemoryString(root, ptr, 4))) == keccak256("hell"));

        ptr = 16;
        root = mem.writeMemoryBytes4(root, ptr, bytes4(""));
        require(keccak256(bytes(mem.readMemoryString(root, ptr, 0))) == keccak256(""));

        ptr = 20;
        for (uint32 offset = 0; offset < 32; offset += 4) {
            root = mem.writeMemoryBytes4(root, ptr + offset, bytes4("stri"));
        }
        root = mem.writeMemoryBytes2(root, ptr + 32, bytes2("st"));

        root = mem.writeMemoryBytes4(root, ptr + 28, bytes4("stri"));
        require(
            keccak256(bytes(mem.readMemoryString(root, ptr, 34))) == keccak256("stristristristristristristristrist")
        );
    }
}
