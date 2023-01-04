// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

library UnsafeMath {
    function unsafeIncrement(uint256 i) internal pure returns (uint256) {
        unchecked {
            return i + 1;
        }
    }
}
