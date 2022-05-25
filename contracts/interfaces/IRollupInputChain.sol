// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IRollupInputChain {
    ///EVENT
    event TransactionEnqueued(
        uint64 indexed queueIndex,
        address indexed from,
        address indexed to,
        uint256 gaslimit,
        bytes data,
        uint64 nonce,
        uint256 r,
        uint256 s,
        uint64 v,
        uint64 timestamp
    );

    /**
     * @dev Adds a transaction to the queue.This function do not need to check tx or pay tx's gas fee,it's paid in L2.Normal EOA just need
     to send a L2 tx.However, L1CrossLayerWitness do not need to sign L2 tx, it's signed by l2 system
     * @param _target Target contract to send the transaction to.
     * @param _gasLimit Gas limit for the given transaction.
     * @param _data Transaction data.
     * @param _nonce sender's nonce in L2
     * @param r,s,v tx signature,some tx's param is set on contract: gasPrice(1 GWEI), value(0), chainId
     * @notice Revert if contract caller isn't l1CrossLayerWitness contract(make sure L1 contract can't act as L2 EOA)
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

    event TransactionAppended(
        address indexed proposer,
        uint256 indexed startQueueIndex,
        uint256 queueNum,
        uint256 indexed chainHeight,
        bytes32 inputHash
    );

    /**
     * append a batches of sequenced tx to input chain.Only staking sender permitted
     * @dev The info is in calldata,format as:
     *  uint64 (num_queue) || uint64 (queue_start_index)||uint64 (num_sequenced) || [uint64,uint64...] (timestamp)) || uint64 (batch_version) [batch_sequenced,batch...]
     * @notice Revert if:
     * - sender isn't staking
     * - queue_start_index not equal to pending queue index
     * - pending queue length beyond queue length locally(make sure can't attempt to append nonexistent queue)
     * - first sequenced tx's timestamp smaller than or equal to  lastTimeStamp(make sure next sequenced tx timestamp larger than lastTimestamp)
     * - last sequenced tx's timestamp smaller than last appended queue(to make sure last sequenced tx timestamp is largest)
     * - last sequenced tx's timestamp larger than or equal to next pending queue timestamp(make sure next pending queue timestamp larger than lastTimestamp )
     * - sequenced tx n timestamp not larger than sequenced tx n-1 timestamp(make sure all sequenced tx timestamp larger than lastTimestamp)
     */
    function appendBatch() external;

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

    //    function getQueueTxInfo(uint64 _queueIndex) external view returns (bytes32, uint64);
}
