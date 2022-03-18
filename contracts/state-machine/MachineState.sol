// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./Memory.sol";
import "./riscv32/Register.sol";
import "../libraries/BytesSlice.sol";
import "./riscv32/Syscall.sol";

contract MachineState {
    using Memory for mapping(bytes32 => bytes);
    using Register for mapping(bytes32 => bytes);
    using Syscall for mapping(bytes32 => bytes);
    mapping(bytes32 => bytes) hashdb;

    function insertTrieNode(bytes calldata _node) public {
        hashdb[keccak256(_node)] = _node;
    }

    function preimage(bytes32 _hash) public view returns (bytes memory _ret, uint32 _len) {
        _ret = hashdb[_hash];
        require(_ret.length > 0, "no image");
        require(_ret.length < uint32((1 << 32) - 1), "image too big");
        return (_ret, uint32(_ret.length));
    }

    function preimageAt(bytes32 _hash, uint32 pos) public view returns (uint32) {
        (bytes memory _ret, uint32 length) = preimage(_hash);
        bytes memory _data;
        uint32 len = length - pos >= 4 ? 4 : length - pos; //overflow safe
        _data = BytesSlice.toBytes(BytesSlice.slice(_ret, pos, len));
        return BytesEndian.bytes4ToUint32(bytes4(_data));
    }

    function writeMemory(
        bytes32 root,
        uint32 ptr,
        uint32 value
    ) public returns (bytes32) {
        return hashdb.writeMemory(root, ptr, value);
    }

    function writeMemoryBytes4(
        bytes32 root,
        uint32 ptr,
        bytes4 value
    ) public returns (bytes32) {
        return hashdb.writeMemoryBytes4(root, ptr, value);
    }

    function writeMemoryByte(
        bytes32 root,
        uint32 ptr,
        bytes1 value
    ) public returns (bytes32) {
        return hashdb.writeMemoryByte(root, ptr, value);
    }

    function writeMemoryBytes2(
        bytes32 root,
        uint32 ptr,
        bytes2 value
    ) public returns (bytes32) {
        return hashdb.writeMemoryBytes2(root, ptr, value);
    }

    function writeMemoryBytes32(
        bytes32 root,
        uint32 ptr,
        bytes32 val
    ) public returns (bytes32) {
        return hashdb.writeMemoryBytes32(root, ptr, val);
    }

    function readMemoryBytes2(bytes32 root, uint32 ptr) public view returns (bytes2) {
        return hashdb.readMemoryBytes2(root, ptr);
    }

    function readMemoryByte(bytes32 root, uint32 ptr) public view returns (bytes1) {
        return hashdb.readMemoryByte(root, ptr);
    }

    function readMemoryBytes4(bytes32 root, uint32 ptr) public view returns (bytes4) {
        return hashdb.readMemoryBytes4(root, ptr);
    }

    function readMemory(bytes32 root, uint32 ptr) public view returns (uint32) {
        return hashdb.readMemory(root, ptr);
    }

    function readMemoryBytes32(bytes32 root, uint32 ptr) public view returns (bytes32) {
        return hashdb.readMemoryBytes32(root, ptr);
    }

    function writeRegisterBytes4(
        bytes32 root,
        uint32 regid,
        bytes4 value
    ) public returns (bytes32) {
        return hashdb.writeRegisterBytes4(root, regid, value);
    }

    function writeRegister(
        bytes32 root,
        uint32 regid,
        uint32 value
    ) public returns (bytes32) {
        return hashdb.writeRegister(root, regid, value);
    }

    function readRegisterBytes4(bytes32 root, uint32 regid) public view returns (bytes4) {
        return hashdb.readRegisterBytes4(root, regid);
    }

    function readRegister(bytes32 root, uint32 regid) public view returns (uint32) {
        return hashdb.readRegister(root, regid);
    }

    function readMemoryString(
        bytes32 _root,
        uint32 addr,
        uint32 len
    ) public view returns (string memory) {
        return hashdb.readMemoryString(_root, addr, len);
    }

    function writeOutput(bytes32 root, bytes32 hash) public returns (bytes32) {
        return hashdb.writeOutput(root, hash);
    }

    function readOutput(bytes32 root) public view returns (bytes32) {
        return hashdb.readOutput(root);
    }

    function writeInput(bytes32 root, bytes32 hash) public returns (bytes32) {
        return hashdb.writeInput(root, hash);
    }

    function readInput(bytes32 root) public view returns (bytes32) {
        return hashdb.readInput(root);
    }

    function genReservedKey(uint32 addr) public view returns (bytes memory) {
        return bytes.concat(bytes5(bytes4(addr)));
    }

    function lr(bytes32 root, uint32 addr) public returns (bytes32) {
        return MerkleTrie.update(hashdb, genReservedKey(addr), bytes.concat(bytes1(uint8(1))), root);
    }

    function sc(bytes32 root, uint32 addr) public returns (bytes32) {
        return MerkleTrie.update(hashdb, genReservedKey(addr), bytes.concat(bytes1(uint8(0))), root);
    }

    function isReserved(bytes32 root, uint32 addr) public view returns (bool) {
        (bool exist, bytes memory value) = MerkleTrie.get(hashdb, genReservedKey(addr), root);
        if (!exist) {
            return false;
        }
        if (bytes1(value) == bytes1(0)) {
            return false;
        }
        return true;
    }
}
