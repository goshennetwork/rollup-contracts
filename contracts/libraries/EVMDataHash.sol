// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library EVMDataHash {
    /// @dev "hello, world"
    bytes32 constant FAKE_VERSION_HASH_LIST_0 =
        bytes32(uint256(0x0194afc31faadfe83ee2a8a35cd92ec08b864e3be052ddfa82c44cf12cca9738));

    //    bytes32 constant FAKE_VERSION_HASH_LIST_1 = bytes32(uint256(0xffff));

    function datahash(uint256 _index) internal returns (bytes32) {
        if (_index == 0) {
            return FAKE_VERSION_HASH_LIST_0;
        }
        //        if (_index == 1) {
        //            return FAKE_VERSION_HASH_LIST_1;
        //        }
        return bytes32(0);
    }
}
