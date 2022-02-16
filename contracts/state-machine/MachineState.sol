// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./Memory.sol";
import "./riscv32/Register.sol";

contract MachineState {
    using Memory for mapping(bytes32 => bytes);
    using Register for mapping(bytes32 => bytes);
    mapping(bytes32 => bytes) hashdb;

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
