// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface IDAO {
    /// @notice who can submit input to L1 contract
    function sequencerWhitelist(address) external view returns (bool);

    /// @notice who can submit output to L1 contract
    function proposerWhitelist(address) external view returns (bool);

    /// @notice who can challenge
    function challengerWhitelist(address) external view returns (bool);
}
