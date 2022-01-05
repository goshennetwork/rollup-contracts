// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

interface IStateCommitChain {
    function getCurrentBlockHeight() external view returns (uint256);

    function isBlockComfirmed(uint256 blockHeight) external view returns (bool);

    function getBlockInfo(uint256 blockHeight)
        external
        view
        returns (
            bytes32 root,
            address proposer,
            uint256 timestamp
        );

    function rollbackBlockBefore(uint256 fraultBlock) external;
}
