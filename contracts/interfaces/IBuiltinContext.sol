// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IBuiltinContext {
    ///@dev Builtin contract that use mmrOracle to get l1 mmr state
    function l1MMRRoot() external view returns (bytes32, uint64);

    ///@return get mmr inclusion proof
    function getMMRInclusionProof(
        bytes32 root,
        uint64 treeSize,
        uint64 leafIndex
    ) external view returns (bytes32[] memory);
}
