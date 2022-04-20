// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IL2CrossLayerMessageWitness {
    /**
     * @dev Relay L1 -> L2 message that in L1CrossLayerMessageWitness contract.
     * @param _target EVM call Target
     * @param _sender EVM call sender
     * @param _message EVM call data
     * @param _messageIndex index in l1 merkle mountain range's leaf
     * @notice Revert if:
     * - sender not L1CrossLayerMessageWitness and can't proof message indeed in l1 mmr
     * - message already relayed
     */
    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex
    ) external;

    event MessageRelayFailed(bytes32 indexed _msgHash, uint64 _mmrSize, bytes32 _mmrRoot);
    event MessageRelayed(uint64 indexed _messageIndex, bytes32 indexed _msgHash);

    /**
     * @dev Relay L1 -> L2 message when obvious relayed false
     * @param _target EVM call Target
     * @param _sender EVM call sender
     * @param _message EVM call data
     * @param _messageIndex index in l1 merkle mountain range's leaf
     * @param _proof Merkle mountain range inclusion proof
     * @notice this function get l1 mmr root and size by builtin contract.and the mmr root only after l1->l2 tx failed.
     * Revert if:
     * - Provided tree in proof not consistent with l1 mmr root got by builtinContext
     * - Provided _proof cant proof message indeed exist in l1 mmr root
     * - Provided message already relayed
     */
    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes32[] memory _proof
    ) external;

    event MessageSent(uint64 indexed _messageIndex, address indexed _target, address indexed _sender, bytes _message);

    //Send L1 -> L2 message
    /**
     * @dev Send message to L1CrossLayerMessageWitness
     * @param _target EVM call target
     * @param _message EVM call data
     */
    function sendMessage(address _target, bytes calldata _message) external;
}
