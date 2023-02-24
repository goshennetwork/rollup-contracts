// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library EVMPreCompiled {
    /// @dev "hello, world"
    bytes32 constant FAKE_VERSION_HASH_LIST_0 =
        bytes32(uint256(0x0194afc31faadfe83ee2a8a35cd92ec08b864e3be052ddfa82c44cf12cca9738));

    //    bytes32 constant FAKE_VERSION_HASH_LIST_1 = bytes32(uint256(0xffff));

    address constant POINT_EVALUATION_PRECOMPILE_ADDRESS = address(0x14);

    function datahash(uint256 _index) internal returns (bytes32) {
        if (_index == 0) {
            return FAKE_VERSION_HASH_LIST_0;
        }
        //        if (_index == 1) {
        //            return FAKE_VERSION_HASH_LIST_1;
        //        }
        return bytes32(0);
    }

    // versioned hash: first 32 bytes
    // Evaluation point: next 32 bytes
    // Expected output: next 32 bytes
    // input kzg point: next 48 bytes
    // Quotient kzg: next 48 bytes
    function point_evaluation_precompile(bytes memory d) internal view {
        (bool ok, ) = POINT_EVALUATION_PRECOMPILE_ADDRESS.staticcall(abi.encodePacked(d));
        require(ok, "point_evaluation failed");
    }
}
