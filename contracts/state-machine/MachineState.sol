// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./Memory.sol";
import "./riscv32/Register.sol";
import "../libraries/BytesSlice.sol";
import "./riscv32/Syscall.sol";
import "../interfaces/IMachineState.sol";
import "../libraries/HashDB.sol";
import "../libraries/BlobDB.sol";

contract MachineState is IMachineState {
    using HashDB for mapping(bytes32 => HashDB.Preimage);
    using Memory for mapping(bytes32 => HashDB.Preimage);
    using Register for mapping(bytes32 => HashDB.Preimage);
    using Syscall for mapping(bytes32 => HashDB.Preimage);
    mapping(bytes32 => HashDB.Preimage) hashdb;

    using BlobDB for mapping(bytes32 => uint256[]);
    mapping(bytes32 => uint256[]) blobdb;

    function insertBlobAt(
        bytes32 _versionHash,
        uint64 _index,
        uint256 _y,
        bytes1[48] memory _commitment,
        bytes1[48] memory _proof
    ) public {
        blobdb.insertBlobAt(_versionHash, _index, _y, _commitment, _proof);
    }

    function readBlobAt(bytes32 _versionHash, uint32 _index) public view returns (bytes32) {
        return blobdb.readBlobAt(_versionHash, _index);
    }

    function insertPreimage(bytes calldata _node) public {
        hashdb.insertPreimage(_node);
    }

    function insertPartialImage(bytes calldata _node, uint32 _index) public {
        hashdb.insertPartialImage(_node, _index);
    }

    function preimage(bytes32 _hash) public view returns (bytes memory _ret) {
        return hashdb.preimage(_hash);
    }

    function preimageAt(bytes32 _hash, uint32 pos) public view returns (uint32) {
        uint32 _index = uint32(pos / HashDB.PartialSize);
        uint256 offset = uint256(pos) % HashDB.PartialSize; // convert to local offset
        bytes memory _ret = hashdb.preimageAtIndex(_hash, _index);
        uint256 len = _ret.length - offset;
        if (len > 4) {
            len = 4;
        }
        bytes memory _data = BytesSlice.toBytes(BytesSlice.slice(_ret, offset, len));
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

    function writeMemoryAddr(
        bytes32 root,
        uint32 ptr,
        address val
    ) public returns (bytes32) {
        bytes20 data = bytes20(val);
        for (uint32 i; i < 20; i += 4) {
            root = hashdb.writeMemoryBytes4(root, ptr + i, bytes4(data));
            data <<= 32;
        }
        return root;
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

    function reserve(bytes32 root, uint32 addr) public returns (bytes32) {
        return hashdb.writeRegister(root, Register.REG_RESV, addr);
    }

    function unReserve(bytes32 root) public returns (bytes32) {
        return hashdb.writeRegister(root, Register.REG_RESV, 0);
    }

    function isReserved(bytes32 root, uint32 addr) public view returns (bool) {
        return hashdb.readRegister(root, Register.REG_RESV) == addr;
    }
}
