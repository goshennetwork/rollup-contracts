// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IChainStorageContainer.sol";
import "../interfaces/IAddressResolver.sol";

contract ChainStorageContainer is IChainStorageContainer {
    IAddressResolver addressResolver;
    bytes32[] chain;
    // the last chain element time stamp, it is simply set with largest timestamp in current tx batch.
    uint64 public override lastTimestamp;

    //the total num of elements in chain, we cut the chain simply change this num
    uint64 public override chainSize;

    string owner;

    constructor(string memory _owner, address _addressResolver) {
        owner = _owner;
        addressResolver = IAddressResolver(_addressResolver);
    }

    modifier onlyOwner() {
        require(
            msg.sender == addressResolver.resolve(owner),
            "ChainStorageContainer: Function can only be called by the owner."
        );
        _;
    }

    function append(bytes32 _element) public onlyOwner {
        if (chainSize < chain.length) {
            //has some unused storage, reuse it
            chain[chainSize] = _element;
        } else {
            //append new chain storage
            chain.push(_element);
        }
        chainSize++;
    }

    function resize(uint64 _newSize) public onlyOwner {
        require(_newSize <= chain.length, "can't resize beyond chain length");
        chainSize = _newSize;
    }

    function setLastTimestamp(uint64 _timestamp) public onlyOwner {
        lastTimestamp = _timestamp;
    }

    function get(uint64 _index) public view returns (bytes32) {
        require(_index < chainSize, "beyond chain size");
        return chain[_index];
    }
}
