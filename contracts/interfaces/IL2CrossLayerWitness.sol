// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IL2CrossLayerWitness {
    /**
     * @dev Relay L1 -> L2 message that in L1CrossLayerWitness contract.
     * @param _target EVM call Target
     * @param _sender EVM call sender
     * @param _message EVM call data
     * @param _messageIndex index in l1 merkle mountain range's leaf
     * @param _messageIndex l1 merkle mountain range root
     * @param _mmrSize l1 merkle mountain range tree size
     * @notice Revert if:
     * - sender isn't L1CrossLayerWitness.
     * - message already relayed
     */
    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes32 _mmrRoot,
        uint64 _mmrSize
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
     * @param _mmrSize L1 merkle mountain range tree size
     * @notice Revert if:
     * - Provided mmrSize have no related mmrRoot.(which means first relay message didn't successful finish or relay succeed)
     * - Provided _proof cant proof message indeed exist in l1 mmr root got by local recorded
     * - Provided message already relayed
     */
    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes32[] memory _proof,
        uint64 _mmrSize
    ) external;

    event MessageSent(uint64 indexed _messageIndex, address indexed _target, address indexed _sender, bytes _message);

    //Send L1 -> L2 message
    /**
     * @dev Send message to L1CrossLayerWitness
     * @param _target EVM call target
     * @param _message EVM call data
     * @notice Revert if sender is self
     */
    function sendMessage(address _target, bytes calldata _message) external;
}
