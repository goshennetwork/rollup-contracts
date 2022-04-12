// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IAddressManager.sol";
import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";

contract AddressManager is IAddressManager, Ownable {
    mapping(bytes32 => address) private addrs;

    ///cant set empty address
    modifier noEmptyAddr(address _addr) {
        require(_addr != address(0), "set empty addr not allowed");
        _;
    }

    function newAddr(string memory _name, address _addr) public onlyOwner noEmptyAddr(_addr) {
        bytes32 _hash = hash(_name);
        require(addrs[_hash] == address(0), "address already exist");
        addrs[_hash] = _addr;
    }

    function updateAddr(string memory _name, address _addr) public onlyOwner noEmptyAddr(_addr) {
        bytes32 _hash = hash(_name);
        require(addrs[_hash] != address(0), "can't update empty addr, use newAddr instead");
        addrs[_hash] = _addr;
    }

    function getAddr(string memory _name) public view returns (address) {
        return addrs[hash(_name)];
    }

    function hash(string memory _name) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(_name));
    }
}
