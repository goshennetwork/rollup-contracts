// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import { Types } from "../libraries/Types.sol";
import "../interfaces/IStakingManager.sol";
import { ICanonicalTransactionChain } from "../interfaces/ICanonicalTransactionChain.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IChainStorageContainer.sol";
import "../libraries/Constants.sol";

contract CanonicalTransactionChain is ICanonicalTransactionChain {
    uint256 public constant MIN_ROLLUP_TX_GAS = 100000;
    uint256 public constant MAX_ROLLUP_TX_SIZE = 50000;

    uint256 public maxEnqueueTxGasLimit;
    uint256 public maxCrossLayerTxGasLimit;

    using Types for Types.QueueElement;
    IAddressResolver addressResolver;

    //store L1 -> L2 tx
    Types.QueueElement[] queueElements;
    // index of the first queue element not yet included
    uint64 public override pendingQueueIndex;

    constructor(address _addressResolver, uint256 _maxTxGasLimit, uint256 _maxCrossLayerTxGasLimit) {
        addressResolver = IAddressResolver(_addressResolver);
        maxEnqueueTxGasLimit = _maxTxGasLimit;
        maxCrossLayerTxGasLimit = _maxCrossLayerTxGasLimit;
    }

    function enqueue(address _target, uint256 _gasLimit, bytes memory _data) public {
        require(_data.length <= MAX_ROLLUP_TX_SIZE, "too large Tx data size");
        require(_gasLimit <= maxEnqueueTxGasLimit, "too high Tx gas limit");
        require(_gasLimit >= MIN_ROLLUP_TX_GAS, "too low Tx gas limit");
        // L1 EOA is equal to L2 EOA, but L1 contract is not except L1CrossLayerMessageWitness
        address sender;
        if (msg.sender == tx.origin) {
            sender = msg.sender;
        } else {
            require(msg.sender == addressResolver.l1CrossLayerMessageWitness(), "contract can not enqueue L2 Tx");
            require(_gasLimit <= maxCrossLayerTxGasLimit, "too high cross layer Tx gas limit");
            sender = Constants.L1_CROSS_LAYER_MESSAGE_WITNESS;
        }
        // todo: maybe need more tx params, such as tip fee, value
        bytes32 transactionHash = keccak256(abi.encode(sender, _target, _gasLimit, _data));
        uint64 _now = uint64(block.timestamp);
        queueElements.push(Types.QueueElement({ transactionHash: transactionHash, timestamp: _now }));
        emit Enqueued(uint64(queueElements.length - 1), sender, _target, _gasLimit, _data, _now);
    }

    function calculateQueueTxHash(uint64 _queueStartIndex, uint64 _queueNum) internal returns (bytes32) {
        bytes memory _queueHash = new bytes(32 * _queueNum);
        uint256 ptr;
        assembly {
            ptr := add(_queueHash, 32)
        }
        uint64 _offset;
        for (uint256 i = 0; i < _queueNum; i++) {
            bytes32 _h = (queueElements[_queueStartIndex + i].hash());
            assembly {
                mstore(add(ptr, _offset), _h)
            }
        }
        return keccak256(_queueHash);
    }

    // format: queueNum(uint64) + queueStart(uint64) + batchNum(uint64) + batchTime([]uint64) + batchesData
    function appendBatch() public {
        require(addressResolver.stakingManager().isStaking(msg.sender), "Sequencer should be staking");
        IChainStorageContainer _chain = addressResolver.ctcContainer();
        uint64 _queueNum;
        uint64 _queueStartIndex;
        assembly {
            _queueNum := shr(192, calldataload(4))
            _queueStartIndex := shr(192, calldataload(12))
        }
        require(_queueStartIndex == pendingQueueIndex, "incorrect pending queue index");
        uint64 _nextPendingQueueIndex = _queueStartIndex + _queueNum;
        require(_nextPendingQueueIndex <= queueElements.length, "attempt to append unavailable queue");
        bytes32 _queueHashes = calculateQueueTxHash(_queueStartIndex, _queueNum);
        uint64 _batchDataPos = 4 + 8 + 8; //4byte function selector, 2 uint64
        pendingQueueIndex = _nextPendingQueueIndex;
        //check sequencer timestamp
        uint64 _batchNum;
        assembly {
            _batchNum := shr(192, _batchDataPos)
        }
        uint64 _timestamp;
        uint64 _lastTimestamp;
        //clear
        _offset = 0;
        for (uint64 i = 0; i < _batchNum; i++) {
            _offset = _batchDataPos + 8 + 8 * i;
            assembly {
                _timestamp := shr(192, calldataload(_offset))
            }
            if (i == 0) {
                //first
                require(_timestamp > _chain.lastTimestamp(), "start timestamp should be larger than previous timestamp");
            }
            if (i == _batchNum - 1) {
                //last
                if (pendingQueueIndex > 0) {
                    //make sure lastBatchTimestamp is the largest
                    require(
                        _timestamp >= queueElements[pendingQueueIndex - 1].timestamp,
                        "last sequenced tx timestamp should larger than appended queue timestamp"
                    );
                }
                if (pendingQueueIndex < queueElements.length) {
                    //make sure lastBatchTimestamp smaller than pending queue.
                    require(
                        _timestamp < queueElements[pendingQueueIndex].timestamp,
                        "last batch muse less than pending queue timestamp"
                    );
                }
            }
            require(_timestamp >= _lastTimestamp, "sequenced batch timestamp should be continuous");
            _lastTimestamp = _timestamp;
        }
        //record batches info
        _chain.append(keccak256(abi.encodePacked(keccak256(msg.data), _queueHashes)));
        _chain.setLastTimestamp(_lastTimestamp);
        emit Appended(msg.sender, _queueStartIndex, _queueNum, _chain.chainSize() - 1);
    }

    function chainHeight() public view returns (uint64) {
        return addressResolver.ctcContainer().chainSize();
    }

    function lastTimestamp() public view returns (uint64) {
        return addressResolver.ctcContainer().lastTimestamp();
    }
}
