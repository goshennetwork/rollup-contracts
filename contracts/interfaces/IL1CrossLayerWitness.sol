// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IL1CrossLayerWitness {
    event MessageRelayFailed(bytes32 indexed _msgHash, uint64 _mmrSize, bytes32 _mmrRoot);
    event MessageRelayed(uint64 indexed _messageIndex, bytes32 indexed _msgHash);

    /**
     * @dev Relay L2 -> L1 message that in L2CrossLayerWitness contract.
     * @param _target EVM call Target
     * @param _sender EVM call sender
     * @param _message EVM call data
     * @param _messageIndex index in l2 merkle mountain range's leaf
     * @param _rlpHeader L2 block header contains l2 mmr root and size
     * @param _stateInfo L2 stateInfo contains block hash
     * @param _proof MMR proof that used to proof provided info surly exists in l2 block mmr
     * @notice Revert if:
     * - reentrancy
     * - target is l1 system contract.(In this case, anyone can't send any calldata to L2 relay contract)
     * - provide wrong state info(not exist in StateCommitChain)
     * - provided state info not confirmed.(only confirmed state is right)
     * - provided block is not consistent with state recorded
     * - provided _proof can't proof message indeed exist in l2 block
     * - message already relayed
     * - message blocked
     */
    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes memory _rlpHeader,
        Types.StateInfo memory _stateInfo,
        bytes32[] memory _proof
    ) external;

    /**
     * @dev Send L1->L2 tx to l2,record tx in local mmr
     * @param _target EVM call target
     * @param _message EVM call data
     * @param _gasLimit EVM call gasLimit
     */
    function sendMessage(
        address _target,
        bytes calldata _message,
        uint64 _gasLimit
    ) external;

    /**
     * @dev Replay failed L2->L1 message.We only assume that poor gasLimit is the only failed reason.So this update old gasLimit
     * @param _crossLayerCalldata ols EVM call data to L2CrossLayerWitness contract
     * @param _queueIndex Replayed tx in queue index
     * @param _oldGasLimit Old gasLimit
     * @param _newGasLimit New gasLimit
     * @notice Revert if:
     * - Provided message not enqueued
     */
    function replayMessage(
        bytes memory _crossLayerCalldata,
        uint64 _queueIndex,
        uint64 _oldGasLimit,
        uint64 _newGasLimit
    ) external;

    event MessageBlocked(bytes32[] _messageHashes);

    /**
     * @dev Block a list of l2 -> l1 message.Only allowed by DAO
     * @param _messageHashes A list of blocked message hash
     */
    function blockMessage(bytes32[] memory _messageHashes) external;

    event MessageAllowed(bytes32[] _messageHashes);

    ///@dev allow a list of L2->L1 message.Only allowed by DAO
    function allowMessage(bytes32[] memory _messageHashes) external;

    ///@return merkle mountain root used to proof l1 -> l2 tx existence
    function mmrRoot() external view returns (bytes32);

    ///@return merkle mountain tx total num(leaf num)
    function totalSize() external view returns (uint64);
}
