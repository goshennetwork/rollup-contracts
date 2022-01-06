// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface IMemory {
    function read(bytes32 stateHash, uint32 addr) external view returns (uint32);

    function readBytes32(bytes32 stateHash, uint32 addr) external view returns (bytes32);

    function write(
        bytes32 stateHash,
        uint32 addr,
        uint32 val
    ) external returns (bytes32);

    function writeBytes32(
        bytes32 stateHash,
        uint32 addr,
        bytes32 val
    ) external returns (bytes32);
}
