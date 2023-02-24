// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./EVMPreCompiled.sol";

library BlobDB {
    uint32 constant FIELD_ELEMENTS_PERBLOB = 4096;

    uint256 constant EXIST_FLAG = 1;

    using BlobDB for mapping(bytes32 => uint256[]);

    function insertBlobAt(
        mapping(bytes32 => uint256[]) storage db,
        bytes32 versionHash,
        uint256 x,
        uint256 y,
        bytes1[48] memory commitment,
        bytes1[48] memory proof
    ) internal {
        EVMPreCompiled.point_evaluation_precompile(abi.encodePacked(versionHash, x, y, commitment, proof));
        /// @notice y is checked in precompiled logic, sure never overflow
        db[versionHash][x] = y + EXIST_FLAG;
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
