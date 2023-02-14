// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library EVMDataHash {
    bytes32 constant FAKE_VERSION_HASH_LIST_0 = bytes32(uint256(0xffff));

    function datahash(uint256 _index) internal returns (bytes32) {
        if (_index == 0) {
            return FAKE_VERSION_HASH_LIST_0;
        }
        return bytes32(0);
    }
}
