// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../MemoryLayout.sol";
import "../../libraries/MerkleTrie.sol";

library Syscall {
    uint8 constant INPUT = 0xfe; // input key flag
    uint8 constant OUTPUT = 0xff; // output key flag

    function writeOutput(mapping(bytes32 => HashDB.Preimage) storage hashdb, bytes32 root, bytes32 value)
        internal
        returns (bytes32)
    {
        return MerkleTrie.update(hashdb, genKey(OUTPUT), bytes.concat(value), root);
    }

    function readOutput(mapping(bytes32 => HashDB.Preimage) storage hashdb, bytes32 root)
        internal
        view
        returns (bytes32)
    {
        (bool exist, bytes memory value) = MerkleTrie.get(hashdb, genKey(OUTPUT), root);
        return exist ? bytes32(value) : bytes32(0);
    }

    function writeInput(mapping(bytes32 => HashDB.Preimage) storage hashdb, bytes32 root, bytes32 value)
        internal
        returns (bytes32)
    {
        return MerkleTrie.update(hashdb, genKey(INPUT), bytes.concat(value), root);
    }

    function readInput(mapping(bytes32 => HashDB.Preimage) storage hashdb, bytes32 root)
        internal
        view
        returns (bytes32)
    {
        (bool exist, bytes memory value) = MerkleTrie.get(hashdb, genKey(INPUT), root);
        return exist ? bytes32(value) : bytes32(0);
    }

    function genKey(uint8 flag) internal pure returns (bytes memory) {
        return bytes.concat(bytes1(flag));
    }
}
