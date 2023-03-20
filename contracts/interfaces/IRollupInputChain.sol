// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IRollupInputChain {
    ///EVENT
    event TransactionEnqueued(
        uint64 indexed queueIndex,
        address indexed from,
        address indexed to,
        bytes rlpTx,
        uint64 timestamp
    );

    event InputBatchAppended(
        address indexed proposer,
        uint64 indexed index,
        uint64 startQueueIndex,
        uint64 queueNum,
        bytes32 inputHash
    );

    /**
     * @dev Adds a transaction to the queue.This function do not need to check tx or pay tx's gas fee,it's paid in L2.Normal EOA just need
     to send a L2 tx.However, L1CrossLayerWitness do not need to sign L2 tx, it's signed by this function
     * @param _target Target contract to send the transaction to.
     * @param _gasLimit Gas limit for the given transaction.
     * @param _data Transaction data.
     * @param _nonce sender's nonce in L2, start from 1<<63 now
     * @param r,s,v tx signature,some tx's param is set on contract: gasPrice(1 GWEI), value(0), chainId
     * @notice Revert if :
     * - contract caller isn't l1CrossLayerWitness contract(make sure L1 contract can't act as L2 EOA)
     * - call data size overhead
     * - nonce not consistent with recorded nonce in local contract
     *
     * - Anyone who tries to use UnsafeSigner's private key to enqueue
     */
    function enqueue(
        address _target,
        uint64 _gasLimit,
        bytes memory _data,
        uint64 _nonce,
        uint256 r,
        uint256 s,
        uint64 v
    ) external;

    /**
     * append a batches of sequenced tx to input chain.Only staking sender permitted
     * @dev The info is in calldata,format as: // format: batchIndex(uint64) + batchIndex(uint64)+ queueNum(uint64) + queueStartIndex(uint64)  + subBatchNum(uint64) + subBatch0Time(uint64) +
    // subBatchLeftTimeDiff([]uint32) + batchesData
    // batchesData: version(0) + rlp([][]transaction)
     *
     * @notice Revert if:
     * - sender isn't EOA
     * - sender isn't staking
     * - batchIndex not equal to pending batch index
     * - queue_start_index not equal to pending queue index
     * - pending queue length beyond queue length locally(make sure can't attempt to append nonexistent queue)
     * - first sequenced tx's timestamp smaller than  lastTimeStamp or block.timestamp(make sure block.timestamp >= sequenced_tx_timestamp >= lastTimestamp)
     * - txs' largest timestamp larger than next pending queue timestamp(block.timestamp queued_tx_timestamp >= lastTimestamp )
     * - next lastTimestamp larger than next pending queue timestamp, which making sequencer stuck(make sure next_lastTimestamp <=pending_queue_timestamp)
     */
    function appendInputBatch() external;

    /**
     * @dev set force delayed seconds, when a queued tx expired after the force delayed seconds, everyone can force append it to input.
     * @param _forceDelayedSeconds  with which, the enqueued tx can't be force push
     * @notice only permitted by dao
     */
    function setForceDelayedSeconds(uint64 _forceDelayedSeconds) external;

    /**
     * @dev force flush queue if tx queued is already expire forceDelayedSeconds
     * @param _queueStartIndex start queue index
     * @param _queueNum total queue num to force flush
     * @notice required:
     * - only EOA
     * - queue start index is right
     * - queue num should larger than 0
     * - total queue not overhead
     * - the newest tx has expired forceDelayedSeconds
     */
    function forceFlushQueue(uint64 _queueStartIndex, uint64 _queueNum) external;

    ///@return total sequenced input num
    function chainHeight() external view returns (uint64);

    ///@return next pending queue index
    function pendingQueueIndex() external view returns (uint64);

    ///@return total queueNum
    function totalQueue() external view returns (uint64);

    ///@return lastTimestamp of RollupInputChain
    function lastTimestamp() external view returns (uint64);

    ///@return input hash related input index in rollup input chain.
    function getInputHash(uint64 _inputIndex) external view returns (bytes32);

    function getQueueTxInfo(uint64 _queueIndex) external view returns (bytes32, uint64);

    /// @return sender's nonce
    function getNonceByAddress(address _sender) external view returns (uint64);

    /// @return the force second delayed time for a queue to batch
    function forceDelayedSeconds() external view returns (uint64);
}
