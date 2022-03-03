// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;
import "../MemoryLayout.sol";

library Syscall {
    uint32 constant RUNTIME_INPUT = 0;
    uint32 constant RUNTIME_RETURN = 1;
    uint32 constant RUNTIME_PREIMAGE_LEN = 2;
    uint32 constant RUNTIME_PREIMAGE = 3;
    uint32 constant RUNTIME_PANIC = 4;
    uint32 constant RUNTIME_DEBUG = 5;
}
