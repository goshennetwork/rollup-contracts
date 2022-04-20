// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IChainStorageContainer.sol";
import "../interfaces/IAddressResolver.sol";

contract ChainStorageContainer is IChainStorageContainer {
    IAddressResolver addressResolver;
    bytes32[] chain;
    // the last chain element time stamp, it is simply set with largest timestamp in current input batch.
    uint64 public override lastTimestamp;

    //who can change the state of this container
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

    function chainSize() external view returns (uint64) {
        return uint64(chain.length);
    }

    function append(bytes32 _element) public onlyOwner {
        chain.push(_element);
    }

    function resize(uint64 _newSize) public onlyOwner {
        require(_newSize <= chain.length, "can't resize beyond chain length");
        assembly {
            sstore(chain.slot, _newSize)
        }
    }

    function setLastTimestamp(uint64 _timestamp) public onlyOwner {
        lastTimestamp = _timestamp;
    }

    function get(uint64 _index) public view returns (bytes32) {
        require(_index < chain.length, "beyond chain size");
        return chain[_index];
    }
}
