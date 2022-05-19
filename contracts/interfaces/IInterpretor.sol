// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface IInterpretor {
    function step(bytes32 root) external returns (bytes32, bool);
}
