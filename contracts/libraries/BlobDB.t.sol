// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./EVMPreCompiled.sol";
import "./BlobDB.sol";

contract TestBlobDB {
    function testWn() public {
        require(BlobDB.calcWn(0) == 1, "w0");
        require(
            BlobDB.calcWn(1) == 39033254847818212395286706435128746857159659164139250548781411570340225835782,
            "w1"
        );
        require(
            BlobDB.calcWn(2) == 49307615728544765012166121802278658070711169839041683575071795236746050763237,
            "w2"
        );
        require(
            BlobDB.calcWn(3) == 24708315984211871914193122998736790630152527847838377463928930981829811447635,
            "w3"
        );
        require(BlobDB.calcWn(4096) == 1, "w4096");
    }
}
