// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./UnsafeSign.sol";
import "./Constants.sol";

contract ConstantsTest {
    function testWitnessAddress() public pure {
        require(Constants.L1_CROSS_LAYER_WITNESS == UnsafeSign.GADDR);
    }
}
