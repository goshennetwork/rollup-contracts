// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./RLPReader.sol";
import "./RLPWriter.sol";
import "./BytesSlice.sol";

contract RLPWriterTest {
    function testWriteAddress(address addr) public pure {
        bytes memory result = RLPWriter.writeAddress(addr);
        address addr2 = RLPReader.readAddress(result);
        require(addr == addr2);
    }
}
