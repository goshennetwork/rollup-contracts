// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./HashDB.sol";
import "./console.sol";

contract TestPartialHashDB {
    using HashDB for mapping(bytes32 => HashDB.Preimage);

    mapping(bytes32 => HashDB.Preimage) partialImage;

    function setUp() public {}

    function testShortInsert() public {
        bytes memory data = bytes("hello world");
        partialImage.insertPartialImage(data, 0);
        bytes32 _hash = keccak256(data);
        bytes memory got = partialImage.preimageAtIndex(_hash, 0);
        require(keccak256(got) == _hash, "not equal");
    }

    function testLongInsert() public {
        bytes memory data = new bytes(10000);
        assembly {
            mstore(add(data, 0x20), 1000)
        }
        bytes32 _hash = keccak256(data);
        partialImage.insertPreimage(data);
        bytes memory got = partialImage.preimage(_hash);
        require(keccak256(got) == _hash, "not equal");
    }

    function test1024Insert() public {
        bytes memory data = new bytes(1024);
        assembly {
            mstore(add(data, 0x20), 1000)
        }
        bytes32 _hash = keccak256(data);
        partialImage.insertPreimage(data);
        bytes memory got = partialImage.preimage(_hash);
        require(keccak256(got) == _hash, "not equal");
    }
}
