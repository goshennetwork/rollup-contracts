// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./BytesSlice.sol";

contract BytesSliceTest {
    function testToBytes32() public view {
        checkToBytes32("", 0);
        checkToBytes32(hex"01", bytes32(hex"01"));
        checkToBytes32(hex"0102", bytes32(hex"0102"));
        checkToBytes32(
            hex"010203040506070809101112131415161718192021222324252627282930",
            bytes32(hex"010203040506070809101112131415161718192021222324252627282930")
        );
        checkToBytes32(
            hex"0102030405060708091011121314151617181920212223242526272829303132",
            bytes32(hex"0102030405060708091011121314151617181920212223242526272829303132")
        );
        checkToBytes32(
            hex"010203040506070809101112131415161718192021222324252627282930313233",
            bytes32(hex"0102030405060708091011121314151617181920212223242526272829303132")
        );

        checkToBytes32(
            hex"01020304050607080910111213141516171819202122232425262728293031",
            bytes32(hex"01020304050607080910111213141516171819202122232425262728293031")
        );
    }

    function checkToBytes32(bytes memory val, bytes32 expected) private view {
        bytes32 ret = BytesSlice.toBytes32(val);
        require(ret == expected);
    }
}
