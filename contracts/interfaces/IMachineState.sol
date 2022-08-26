// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface IMachineState {
    function insertPreimage(bytes calldata _node) external;

    function preimage(bytes32 _hash) external view returns (bytes memory);

    function preimageAt(bytes32 _hash, uint32 pos) external view returns (uint32);

    function writeMemory(
        bytes32 root,
        uint32 ptr,
        uint32 value
    ) external returns (bytes32);

    function writeMemoryBytes4(
        bytes32 root,
        uint32 ptr,
        bytes4 value
    ) external returns (bytes32);

    function writeMemoryByte(
        bytes32 root,
        uint32 ptr,
        bytes1 value
    ) external returns (bytes32);

    function writeMemoryBytes2(
        bytes32 root,
        uint32 ptr,
        bytes2 value
    ) external returns (bytes32);

    function writeMemoryBytes32(
        bytes32 root,
        uint32 ptr,
        bytes32 val
    ) external returns (bytes32);

    function writeMemoryAddr(
        bytes32 root,
        uint32 ptr,
        address val
    ) external returns (bytes32);

    function readMemoryBytes2(bytes32 root, uint32 ptr) external view returns (bytes2);

    function readMemoryByte(bytes32 root, uint32 ptr) external view returns (bytes1);

    function readMemoryBytes4(bytes32 root, uint32 ptr) external view returns (bytes4);

    function readMemory(bytes32 root, uint32 ptr) external view returns (uint32);

    function readMemoryBytes32(bytes32 root, uint32 ptr) external view returns (bytes32);

    function writeRegisterBytes4(
        bytes32 root,
        uint32 regid,
        bytes4 value
    ) external returns (bytes32);

    function writeRegister(
        bytes32 root,
        uint32 regid,
        uint32 value
    ) external returns (bytes32);

    function readRegisterBytes4(bytes32 root, uint32 regid) external view returns (bytes4);

    function readRegister(bytes32 root, uint32 regid) external view returns (uint32);

    function readMemoryString(
        bytes32 _root,
        uint32 addr,
        uint32 len
    ) external view returns (string memory);

    function writeOutput(bytes32 root, bytes32 hash) external returns (bytes32);

    function readOutput(bytes32 root) external view returns (bytes32);

    function writeInput(bytes32 root, bytes32 hash) external returns (bytes32);

    function readInput(bytes32 root) external view returns (bytes32);

    function genReservedKey(uint32 addr) external pure returns (bytes memory);

    function lr(bytes32 root, uint32 addr) external returns (bytes32);

    function sc(bytes32 root, uint32 addr) external returns (bytes32);

    function isReserved(bytes32 root, uint32 addr) external view returns (bool);
}
