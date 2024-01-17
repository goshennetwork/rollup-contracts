// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../interfaces/IChainStorageContainer.sol";
import "../interfaces/IAddressResolver.sol";

contract ChainStorageContainer is IChainStorageContainer, Initializable {
    IAddressResolver public resolver;
    bytes32[] private chain;

    //who can change the state of this container
    string public owner;

    function initialize(string memory _owner, address _addressResolver) public initializer {
        owner = _owner;
        resolver = IAddressResolver(_addressResolver);
    }

    modifier onlyOwner() {
        require(
            msg.sender == resolver.resolve(owner), "ChainStorageContainer: Function can only be called by the owner."
        );
        _;
    }

    function chainSize() external view returns (uint64) {
        return uint64(chain.length);
    }

    function append(bytes32 _element) public onlyOwner returns (uint64) {
        chain.push(_element);
        return uint64(chain.length);
    }

    function resize(uint64 _newSize) public onlyOwner {
        require(_newSize <= chain.length, "can't resize beyond chain length");
        assembly {
            sstore(chain.slot, _newSize)
        }
    }

    function get(uint64 _index) public view returns (bytes32) {
        require(_index < chain.length, "beyond chain size");
        return chain[_index];
    }
}
