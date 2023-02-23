// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./EVMPreCompiled.sol";

library BlobDB {
    uint32 constant FIELD_ELEMENTS_PERBLOB = 4096;

    using BlobDB for mapping(bytes32 => uint256[4096]);
    uint256 constant VERSION_MASK = ((1 << 8) - 1) << 248;

    uint256 constant ELEMENT_MASK = (1 << 8) - 1;

    uint256 constant EXIST_FLAG = 1;

    function insertBlobAt(
        mapping(bytes32 => uint256[4096]) storage db,
        bytes32 versionHash,
        uint256 x,
        uint256 y,
        bytes1[48] memory commitment,
        bytes1[48] memory proof
    ) internal {
        EVMPreCompiled.point_evaluation_precompile(abi.encodePacked(versionHash, x, y, commitment, proof));
        if (y & ELEMENT_MASK != 0) {
            /// @dev right blob should never happen
            revert("invalid y");
        }
        db[versionHash][x] = y + EXIST_FLAG;
    }

    function readBlobAt(
        mapping(bytes32 => uint256[4096]) storage db,
        bytes32 versionHash,
        uint32 index
    ) internal view returns (bytes32) {
        /// @notice should check index overhead?
        uint256 element = db[versionHash][index];
        if (element & ELEMENT_MASK == 0) {
            /// @dev not exist, if an element is inserted into blob db, the last byte will set to 1.
            revert("no element");
        }
        return bytes32(element - EXIST_FLAG);
    }
}
