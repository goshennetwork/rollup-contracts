// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

interface IChainStorageContainer {
    ///@dev Append element to chain
    function append(bytes32 _element) external;

    ////@dev cut chain size
    ///@notice Revert if _newSize larger than chain size.
    function resize(uint64 _newSize) external;

    ///@return chain size
    function chainSize() external view returns (uint64);

    ///@dev set lastTimeStamp,only use in RollupInputChain
    function setLastTimestamp(uint64 _timestamp) external;

    ///@return lastTimeStamp, only used in RollupInputChain
    function lastTimestamp() external view returns (uint64);

    ///@dev Get element from chain by specific index
    ///@notice Revert if element can't get from chain(index>=chain size)
    function get(uint64 _index) external view returns (bytes32);
}
