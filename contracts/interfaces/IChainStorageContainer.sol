// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

interface IChainStorageContainer {
    function append(bytes32 _element) external;

    ////@@return cut chain size, so newSize must smaller than size
    function resize(uint64 _newSize) external;

    ///@return chain size
    function chainSize() external view returns (uint64);

    function setLastTimestamp(uint64 _timestamp) external;

    function lastTimestamp() external view returns (uint64);

    function get(uint64 _index) external view returns (bytes32);
}
