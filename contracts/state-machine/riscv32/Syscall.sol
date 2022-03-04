// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../MemoryLayout.sol";
import "../../libraries/MerkleTrie.sol";

library Syscall {
    uint32 constant OUTPUT = 0; //output key flag
    uint32 constant INPUT = 1; //input key flag

    function writeOutput(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        bytes32 value
    ) internal returns (bytes32) {
        return MerkleTrie.update(hashdb, genKey(OUTPUT), bytes.concat(value), root);
    }

    function readOutput(mapping(bytes32 => bytes) storage hashdb, bytes32 root) internal view returns (bytes32) {
        (bool exist, bytes memory value) = MerkleTrie.get(hashdb, genKey(OUTPUT), root);
        return exist ? bytes32(value) : bytes32(0);
    }

    function writeInput(
        mapping(bytes32 => bytes) storage hashdb,
        bytes32 root,
        bytes32 value
    ) internal returns (bytes32) {
        return MerkleTrie.update(hashdb, genKey(INPUT), bytes.concat(value), root);
    }

    function readInput(mapping(bytes32 => bytes) storage hashdb, bytes32 root) internal view returns (bytes32) {
        (bool exist, bytes memory value) = MerkleTrie.get(hashdb, genKey(INPUT), root);
        return exist ? bytes32(value) : bytes32(0);
    }

    function genKey(uint32 flag) internal view returns (bytes memory) {
        return bytes.concat(bytes2(uint16(flag)));
    }
}
