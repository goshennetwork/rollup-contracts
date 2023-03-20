// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./EVMPreCompiled.sol";

library BlobDB {
    uint32 constant FIELD_ELEMENTS_PERBLOB = 4096;
    uint256 constant W1 = 39033254847818212395286706435128746857159659164139250548781411570340225835782;
    uint256 constant BLS_MODULE = 52435875175126190479447740508185965837690552500527637822603658699938581184513;
    uint256 constant EXIST_FLAG = 1;
    using BlobDB for mapping(bytes32 => uint256[]);

    function calcWn(uint64 n) internal pure returns (uint256) {
        if (n == 0) {
            return 1;
        }
        uint256 ret = W1;
        for (uint64 i = 1; i < n; i++) {
            assembly {
                ret := mulmod(ret, W1, BLS_MODULE)
            }
        }
        return ret;
    }

    /// @notice index over 4096 also make point evaluation pass, but it is useless, because L2 OS will not read index that >=4096
    /// @dev the y is decided by sequencer, but the insertBlobAt method is used by challenger.So DO NOT assume the y is
    /// satisfied with application protocol that should be checked within L2 OS.So if the y is larger than modules, challenger
    /// should module it too.
    function insertBlobAt(
        mapping(bytes32 => uint256[]) storage db,
        bytes32 versionHash,
        uint32 index,
        uint256 y,
        bytes1[48] memory commitment,
        bytes1[48] memory proof
    ) internal {
        uint256 x = calcWn(index);
        EVMPreCompiled.point_evaluation_precompile(abi.encodePacked(versionHash, x, y, commitment, proof));
        /// @notice y is checked in precompiled logic, sure never overflow
        db[versionHash][index] = y + EXIST_FLAG;
    }

    function readBlobAt(
        mapping(bytes32 => uint256[]) storage db,
        bytes32 versionHash,
        uint32 index
    ) internal view returns (bytes32) {
        /// @notice should check index overhead?
        uint256 element = db[versionHash][index];
        if (element == 0) {
            /// @dev not exist, if an element is inserted into blob db, the last byte will set to 1.
            revert("no element");
        }
        return bytes32(element - EXIST_FLAG);
    }
}
