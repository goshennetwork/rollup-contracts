// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./MerkleTrie.sol";
import "./console.sol";

contract MockMerkleTrie {
    using MerkleTrie for mapping(bytes32 => bytes);
    mapping(bytes32 => bytes) _hashdb;
    mapping(bytes32 => bytes) _rawdb;
    bytes32 public root;

    constructor() {
        root = MerkleTrie.KECCAK256_RLP_NULL_BYTES;
    }

    function update(bytes memory _key, bytes memory _value) external returns (bytes32) {
        _rawdb[keccak256(_key)] = _value;
        root = _hashdb.update(_key, _value, root);
        return root;
    }

    function get(bytes memory _key) external view returns (bool, bytes memory) {
        return _hashdb.get(_key, root);
    }

    function getRaw(bytes memory _key) external view returns (bytes memory) {
        return _rawdb[keccak256(_key)];
    }
}

contract MerkleTrieTest {
    MockMerkleTrie trie;

    function setUp() public {
        trie = new MockMerkleTrie();
    }

    function testGetSet() public {
        string[2][4] memory kvs = [[hex"00", ""], [hex"01", ""], ["02", "cccccc"], ["00000000000003", "d"]];
        for (uint256 i = 0; i < kvs.length; i++) {
            bytes memory k = bytes(kvs[i][0]);
            bytes memory v = bytes(kvs[i][1]);

            trie.update(k, v);
            (bool exist, bytes memory value) = trie.get(k);
            require(exist);
            require(BytesSlice.equal(value, v));
        }
        for (uint256 i = 0; i < kvs.length; i++) {
            bytes memory k = bytes(kvs[i][0]);
            bytes memory v = bytes(kvs[i][1]);
            (bool exist, bytes memory value) = trie.get(k);
            require(exist);
            require(BytesSlice.equal(value, v));
        }
    }

    /**
     *  TODO: this fuzz use too much time to execute
    function testGetFuzz(bytes[2][] memory kvs) public {
        for (uint256 i = 0; i < kvs.length; i++) {
            bytes32 root = trie.update(kvs[i][0], kvs[i][1]);
            console.logBytes32(root);
        }

        for (uint256 i = 0; i < kvs.length; i++) {
            (bool exist, bytes memory value) = trie.get(kvs[i][0]);
            console.logBool(exist);
            console.logBytes(value);
            require(exist);
            require(BytesSlice.equal(value, trie.getRaw(kvs[i][0])));
        }
    }
    */
}