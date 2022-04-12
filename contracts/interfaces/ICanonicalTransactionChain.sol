// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface ICanonicalTransactionChain {
    ///EVENT
    event Enqueued(address indexed from, address to, uint256 gaslimit, bytes data, uint64 queueIndex, uint64 timestamp);

    /**
     * Adds a transaction to the queue.This function do not need to check tx or pay tx's gas fee,it's paid in L2,so L2 treat
     * L1 -> L2 tx as origin tx.
     * @param _target Target contract to send the transaction to.
     * @param _gasLimit Gas limit for the given transaction.
     * @param _data Transaction data.
     */
    function enqueue(
        address _target,
        uint256 _gasLimit,
        bytes memory _data
    ) external;

    /**
     * append a batches of sequenced tx to tx chain.
     * @dev The info is in calldata,format as:
     *  uint64 (num_queue) || uint64 (queue_start_index)||uint64 (num_sequenced) || [uint64,uint64...] (timestamp)) || uint64 (batch_version) [batch_sequenced,batch...]
     */
    function appendBatch() external;

    ///get total sequenced tx batches num
    function chainHeight() external view returns (uint64);

    function pendingQueueIndex() external view returns (uint64);

    function lastTimestamp() external view returns (uint64);
}
