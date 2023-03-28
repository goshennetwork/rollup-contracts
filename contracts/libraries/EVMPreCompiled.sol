// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library EVMPreCompiled {
    /// @dev rawtx: 0xf86e808405f5e100830ece6494dad685e17c1e5208ffa3bd852c20a502774cd134880de0b6b3a76400008082aa3ba0b79fa6d01c478fd2b2217b9a80f6015604359cd1ee2825e73157eed4e9b8d34ca001183e35973bdbf83adf470f035f77074f3fa85d15c8bbe7f721bc24144ed4e5
    /// @dev encoded: "0x8b3980f872f870f86e808405f5e100830ece6494dad685e17c1e5208ffa3bd852c20a502774cd134880de0b6b3a76400008082aa3ba0b79fa6d01c478fd2b2217b9a80f6015604359cd1ee2825e73157eed4e9b8d34ca001183e35973bdbf83adf470f035f77074f3fa85d15c8bbe7f721bc24144ed4e503"
    bytes32 constant FAKE_VERSION_HASH_LIST_0 =
        bytes32(uint256(0x0199976106a5466c745aa976a2d02e257eba22e1f0a56da55151769d9f22b99f));

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
        (bool ok, ) = POINT_EVALUATION_PRECOMPILE_ADDRESS.staticcall(d);
        require(ok, "point_evaluation failed");
    }
}
