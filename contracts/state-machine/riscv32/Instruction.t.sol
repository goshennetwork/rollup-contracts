// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./Instruction.sol";
import "../../libraries/console.sol";

contract InstructionTest {
    function testHaha() public {
        bytes4 mask = bytes4(hex"ff");
        console.logBytes4(mask);
        console.logBytes4(mask >> 8);
        console.logBytes4(~(mask >> 8));
        console.logBytes4((mask >> 8) & mask);
    }
}
