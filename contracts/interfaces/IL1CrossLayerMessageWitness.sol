// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IL1CrossLayerMessageWitness {
    event MessageRelayFailed(bytes32 _calldataHash);
    event MessageRelayed(bytes32 _calldataHash);

    /**
     * @dev Relay L2 -> L1 message that in L2CrossLayerMessageWitness contract.
     * @param _target EVM call Target
     * @param _sender EVM call sender
     * @param _message EVM call data
     * @param _messageIndex index in l2 merkle mountain range's leaf
     * @param _block L2 block contains l2 mmr root and size
     * @param _stateInfo L2 stateInfo contains block hash
     * @param _proof MMR proof that used to proof provided info surly exists in l2 block mmr
     * @notice Revert if:
     * - reentrancy
     * - provide wrong state info(not exist in StateCommitChain)
     * - provided state info not confirmed.(only confirmed state is right)
     * - provided block is not consistent with state recorded
     * - provided _proof can't proof message indeed exist in l2 block
     * - message already relayed
     */
    function relayMessage(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        Types.Block memory _block,
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
     * @param _crossLayerCalldata ols EVM call data to L2CrossLayerMessageWitness contract
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

    ///@return merkle mountain root used to proof l1 -> l2 tx existence
    function mmrRoot() external view returns (bytes32);

    ///@return merkle mountain tx total num(leaf num)
    function totalSize() external view returns (uint64);
}
