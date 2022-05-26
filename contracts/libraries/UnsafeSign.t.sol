// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./UnsafeSign.sol";

contract TestUnsafeSign {
    function testSign() public {
        for (uint256 i = 0; i < 10000; i++) {
            bytes32 signedHash = keccak256(abi.encode(i, "test"));
            (uint256 r, uint256 s, uint64 v) = UnsafeSign.GetRSV(signedHash, 0);
            address sender = ecrecover(signedHash, uint8(v - 35), bytes32(r), bytes32(s));
            require(sender == address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266), "wrong sender");
        }
    }
}
