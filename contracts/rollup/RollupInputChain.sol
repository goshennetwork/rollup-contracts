// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";
import "../interfaces/IStakingManager.sol";
import "../interfaces/IRollupInputChain.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IChainStorageContainer.sol";
import "../libraries/Constants.sol";

contract RollupInputChain is IRollupInputChain {
    uint256 public constant MIN_ROLLUP_TX_GAS = 100000;
    uint256 public constant MAX_ROLLUP_TX_SIZE = 50000;
    uint256 public constant MAX_CROSS_LAYER_TX_SIZE = 10000;

    uint256 public maxEnqueueTxGasLimit;
    uint256 public maxCrossLayerTxGasLimit;

    IAddressResolver addressResolver;

    //store L1 -> L2 tx
    struct QueueTxInfo {
        bytes32 transactionHash;
        uint64 timestamp;
    }

    QueueTxInfo[] queuedTxInfos;
    // index of the first queue element not yet included
    uint64 public override pendingQueueIndex;

    constructor(
        address _addressResolver,
        uint256 _maxTxGasLimit,
        uint256 _maxCrossLayerTxGasLimit
    ) {
        addressResolver = IAddressResolver(_addressResolver);
        maxEnqueueTxGasLimit = _maxTxGasLimit;
        maxCrossLayerTxGasLimit = _maxCrossLayerTxGasLimit;
    }

    function enqueue(
        address _target,
        uint64 _gasLimit,
        bytes memory _data
    ) public {
        require(_data.length <= MAX_ROLLUP_TX_SIZE, "too large Tx data size");
        require(_gasLimit <= maxEnqueueTxGasLimit, "too high Tx gas limit");
        require(_gasLimit >= MIN_ROLLUP_TX_GAS, "too low Tx gas limit");
        // L1 EOA is equal to L2 EOA, but L1 contract is not except L1CrossLayerMessageWitness
        address sender;
        if (msg.sender == tx.origin) {
            sender = msg.sender;
        } else {
            require(
                msg.sender == address(addressResolver.l1CrossLayerMessageWitness()),
                "contract can not enqueue L2 Tx"
            );
            require(_gasLimit <= maxCrossLayerTxGasLimit, "too high cross layer Tx gas limit");
            require(_data.length <= MAX_CROSS_LAYER_TX_SIZE, "too large cross layer Tx data size");
            sender = Constants.L1_CROSS_LAYER_WITNESS;
        }
        // todo: maybe need more tx params, such as tip fee, value
        bytes32 transactionHash = keccak256(abi.encode(sender, _target, _gasLimit, _data));
        uint64 _now = uint64(block.timestamp);
        queuedTxInfos.push(QueueTxInfo({ transactionHash: transactionHash, timestamp: _now }));
        emit TransactionEnqueued(uint64(queuedTxInfos.length - 1), sender, _target, _gasLimit, _data, _now);
    }

    function calculateQueueTxHash(uint64 _queueStartIndex, uint64 _queueNum) internal view returns (bytes32) {
        uint256 len = (32 + 8) * _queueNum;
        bytes memory _queueHash = new bytes(len);
        uint64 _offset = 0;
        for (uint256 i = 0; i < _queueNum; i++) {
            QueueTxInfo memory info = queuedTxInfos[_queueStartIndex + i];
            bytes32 txHash = info.transactionHash;
            bytes32 time = bytes32(uint256(info.timestamp) << 192);
            assembly {
                let ptr := add(_queueHash, _offset)
                mstore(ptr, txHash)
                ptr := add(ptr, 32)
                // @notice we reuse _queueHash's the first 32 byte length bits, so no overflow
                mstore(ptr, time)
            }
            _offset += 40;
        }

        // @notice we reuse _queueHash's length, so can not use keccak256(_queueHash)
        bytes32 result;
        assembly {
            result := keccak256(_queueHash, len)
        }
        return result;
    }

    // format: queueNum(uint64) + queueStart(uint64) + batchNum(uint64) + batch0Time(uint64) +
    // batchLeftTimeDiff([]uint32) + batchesData
    function appendBatch() public {
        require(addressResolver.stakingManager().isStaking(msg.sender), "Sequencer should be staking");
        IChainStorageContainer _chain = addressResolver.rollupInputChainContainer();
        uint64 _queueNum;
        uint64 _queueStartIndex;
        assembly {
            _queueNum := shr(192, calldataload(4))
            _queueStartIndex := shr(192, calldataload(12))
        }
        require(_queueStartIndex == pendingQueueIndex, "incorrect pending queue index");
        uint64 _nextPendingQueueIndex = _queueStartIndex + _queueNum;
        require(_nextPendingQueueIndex <= queuedTxInfos.length, "attempt to append unavailable queue");
        bytes32 _queueHashes = calculateQueueTxHash(_queueStartIndex, _queueNum);
        uint64 _batchDataPos = 4 + 8 + 8; //4byte function selector, 2 uint64
        pendingQueueIndex = _nextPendingQueueIndex;
        //check sequencer timestamp
        uint64 _batchNum;
        assembly {
            _batchNum := shr(192, _batchDataPos)
        }
        require(_batchNum > 0, "no batch");
        _batchDataPos += 8;
        uint64 _timestamp;
        assembly {
            _timestamp := shr(192, calldataload(_batchDataPos))
        }
        require(_timestamp > _chain.lastTimestamp() && _timestamp < block.timestamp, "wrong batch timestap");
        _batchDataPos += 8;
        uint64 _lastTimestamp;
        for (uint64 i = 1; i < _batchNum; i++) {
            uint32 _timediff;
            assembly {
                _timediff := shr(224, calldataload(_batchDataPos))
            }
            _timestamp += uint64(_timediff);
            _batchDataPos += 4;
        }

        if (_nextPendingQueueIndex > 0) {
            uint64 _lastIncludedQueueTime = queuedTxInfos[_nextPendingQueueIndex - 1].timestamp;
            if (_timestamp < _lastIncludedQueueTime) {
                _timestamp = _lastIncludedQueueTime;
            }
        }
        uint64 _nextTimestamp = uint64(block.timestamp);
        if (_nextPendingQueueIndex < queuedTxInfos.length) {
            _nextTimestamp = queuedTxInfos[_nextPendingQueueIndex].timestamp;
        }
        require(_timestamp < _nextTimestamp, "last batch timestamp too high");

        //input msgdata hash, queue hash, mmr root, mmr size
        _chain.append(keccak256(abi.encodePacked(keccak256(msg.data), _queueHashes)));
        _chain.setLastTimestamp(_lastTimestamp);
        emit TransactionAppended(msg.sender, _queueStartIndex, _queueNum, _chain.chainSize() - 1);
    }

    function chainHeight() public view returns (uint64) {
        return addressResolver.rollupInputChainContainer().chainSize();
    }

    function lastTimestamp() public view returns (uint64) {
        return addressResolver.rollupInputChainContainer().lastTimestamp();
    }

    function getQueueTxInfo(uint64 _queueIndex) public view returns (bytes32, uint64) {
        require(_queueIndex < queuedTxInfos.length, "queue index over capacity");
        QueueTxInfo storage info = queuedTxInfos[_queueIndex];
        return (info.transactionHash, info.timestamp);
    }
}
