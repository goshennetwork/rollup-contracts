// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./Memory.sol";
import "./riscv32/Register.sol";
import "../libraries/BytesSlice.sol";

contract MachineState {
    using Memory for mapping(bytes32 => bytes);
    using Register for mapping(bytes32 => bytes);
    mapping(bytes32 => bytes) hashdb;

    function insertTrieNode(bytes calldata _node) public {
        hashdb[keccak256(_node)] = _node;
    }

    function preimage(bytes32 _hash) public view returns (bytes memory _ret) {
        _ret = hashdb[_hash];
        require(_ret.length > 0, "no image");
        return _ret;
    }

    function preimageLen(bytes32 _hash) public view returns (uint32) {
        uint256 _len = preimage(_hash).length;
        require(_len < uint32((1 << 32) - 1), "image too big");
        return uint32(_len);
    }

    function preimagePos(bytes32 _hash, uint32 pos) public view returns (uint32) {
        bytes memory _data = BytesSlice.toBytes(BytesSlice.slice(preimage(_hash), pos, 4));
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
}
