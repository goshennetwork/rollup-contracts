// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

interface ICrossLayerWitness {
    event MessageRelayFailed(bytes32 indexed _msgHash, uint64 _mmrSize, bytes32 _mmrRoot);
    event MessageRelayed(uint64 indexed _messageIndex, bytes32 indexed _msgHash);
    event MessageSent(
        uint64 indexed _messageIndex,
        address indexed _target,
        address indexed _sender,
        bytes32 _mmrRoot,
        bytes _message
    );

    /**
     * @dev Send L1->L2 tx to l2,record tx in local mmr
     * @param _target EVM call target
     * @param _message EVM call data
     */
    function sendMessage(address _target, bytes calldata _message) external;

    function crossLayerSender() external view returns (address);
}
