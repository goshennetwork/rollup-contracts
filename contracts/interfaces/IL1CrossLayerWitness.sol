// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";
import "./ICrossLayerWitness.sol";

interface IL1CrossLayerWitness is ICrossLayerWitness {
    /**
     * @dev Relay L2 -> L1 message that in L2CrossLayerWitness contract.
     * @param _target EVM call Target
     * @param _sender EVM call sender
     * @param _message EVM call data
     * @param _messageIndex index in l2 merkle mountain range's leaf
     * @param _rlpHeader L2 block header contains l2 mmr root and size
     * @param _stateInfo L2 stateInfo contains block hash
     * @param _proof MMR proof that used to proof provided info surly exists in l2 block mmr
     * @return whether relay call message success
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
    ) external returns (bool);

    event MessageBlocked(bytes32[] _messageHashes);

    //check whether specific message already succeed.
    function isMessageSucceed(bytes32 _messageHash) external view returns (bool);

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
