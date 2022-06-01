// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../libraries/Types.sol";
import "../interfaces/IStakingManager.sol";
import "../interfaces/IRollupInputChain.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IChainStorageContainer.sol";
import "../libraries/Constants.sol";
import "../libraries/RLPWriter.sol";
import "../libraries/UnsafeSign.sol";

contract RollupInputChain is IRollupInputChain, Initializable {
    uint256 public constant MIN_ROLLUP_TX_GAS = 100000;
    uint256 public constant MAX_ROLLUP_TX_SIZE = 50000;
    uint256 public constant MAX_CROSS_LAYER_TX_SIZE = 10000;
    uint256 public constant GAS_PRICE = 1_000_000_000;
    uint256 public constant VALUE = 0;

    uint64 public maxEnqueueTxGasLimit;
    uint64 public maxCrossLayerTxGasLimit;
    uint64 public l2ChainID;

    uint64 public override lastTimestamp;

    IAddressResolver addressResolver;

    //store L1 -> L2 tx
    struct QueueTxInfo {
        bytes32 transactionHash;
        uint64 timestamp;
    }

    QueueTxInfo[] queuedTxInfos;
    // index of the first queue element not yet included
    uint64 public override pendingQueueIndex;

    function initialize(
        address _addressResolver,
        uint64 _maxTxGasLimit,
        uint64 _maxCrossLayerTxGasLimit,
        uint64 _l2ChainID
    ) public initializer {
        addressResolver = IAddressResolver(_addressResolver);
        maxEnqueueTxGasLimit = _maxTxGasLimit;
        maxCrossLayerTxGasLimit = _maxCrossLayerTxGasLimit;
        l2ChainID = _l2ChainID;
    }

    function enqueue(
        address _target,
        uint64 _gasLimit,
        bytes memory _data,
        uint64 _nonce,
        uint256 r,
        uint256 s,
        uint64 v
    ) public {
        uint256 _gasPrice = GAS_PRICE;
        address sender;
        // L1 EOA is equal to L2 EOA, but L1 contract is not except L1CrossLayerWitness
        if (msg.sender == tx.origin) {
            sender = msg.sender;
            //make sure only L1CrossLayerWitness use unsafe sender
            require(sender != Constants.L1_CROSS_LAYER_WITNESS, "malicious sender");
            require(_data.length <= MAX_ROLLUP_TX_SIZE, "too large Tx data size");
        } else {
            sender = Constants.L1_CROSS_LAYER_WITNESS;
            require(msg.sender == address(addressResolver.l1CrossLayerWitness()), "contract can not enqueue L2 Tx");
            require(_data.length <= MAX_CROSS_LAYER_TX_SIZE, "too large cross layer Tx data size");
            _gasLimit = maxCrossLayerTxGasLimit;
            _gasPrice = GAS_PRICE;
        }
        require(_gasLimit <= maxEnqueueTxGasLimit, "too high Tx gas limit");
        require(_gasLimit >= MIN_ROLLUP_TX_GAS, "too low Tx gas limit");

        bytes[] memory _rlpList = getRlpList(_nonce, _gasLimit, _gasPrice, _target, _data);
        bytes32 _signTxHash = keccak256(RLPWriter.writeList(_rlpList));
        if (sender == Constants.L1_CROSS_LAYER_WITNESS) {
            // L1CrossLayer need to sign
            // help sign L1CrossLayerWitness
            (r, s, v) = UnsafeSign.Sign(_signTxHash, l2ChainID);
        }
        uint64 _pureV = v - 2 * l2ChainID - 8;
        require(_pureV <= 28, "invalid v");
        require(sender == ecrecover(_signTxHash, uint8(_pureV), bytes32(r), bytes32(s)), "wrong sign");
        //now change rsv value in tx to calc tx's hash
        _rlpList[6] = RLPWriter.writeUint(v);
        _rlpList[7] = RLPWriter.writeUint(r);
        _rlpList[8] = RLPWriter.writeUint(s);
        bytes memory _rlpTx = RLPWriter.writeList(_rlpList);
        uint64 _now = uint64(block.timestamp);
        queuedTxInfos.push(QueueTxInfo({ timestamp: _now, transactionHash: keccak256(_rlpTx) }));

        emit TransactionEnqueued(uint64(queuedTxInfos.length - 1), sender, _target, _rlpTx, _now);
    }

    //encode tx params: sender, to, gasLimit, data, nonce, r,s,v and gasPrice(1 GWEI), value(0), chainId
    //sender used to recognize tx from L1CrossLayerWitness
    function getRlpList(
        uint64 _nonce,
        uint64 _gasLimit,
        uint256 _gasPrice,
        address _target,
        bytes memory _data
    ) internal view returns (bytes[] memory) {
        bytes[] memory list = new bytes[](9);
        list[0] = RLPWriter.writeUint(uint256(_nonce));
        list[1] = RLPWriter.writeUint(_gasPrice);
        list[2] = RLPWriter.writeUint(uint256(_gasLimit));
        list[3] = RLPWriter.writeAddress(_target);
        list[4] = RLPWriter.writeUint(VALUE);
        list[5] = RLPWriter.writeBytes(_data);
        list[6] = RLPWriter.writeUint(l2ChainID);
        list[7] = abi.encodePacked(bytes1(0x80));
        list[8] = abi.encodePacked(bytes1(0x80));
        return list;
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
        require(addressResolver.dao().sequencerWhitelist(msg.sender), "only sequencer");
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
        uint64 _batchDataPos = 4 + 8 + 8;
        //4byte function selector, 2 uint64
        pendingQueueIndex = _nextPendingQueueIndex;
        //check sequencer timestamp
        uint64 _batchNum;
        assembly {
            _batchNum := shr(192, calldataload(_batchDataPos))
        }
        require(_batchNum > 0, "no batch");
        _batchDataPos += 8;
        uint64 _timestamp;
        assembly {
            _timestamp := shr(192, calldataload(_batchDataPos))
        }
        require(_timestamp >= lastTimestamp, "wrong batch timestamp");
        _batchDataPos += 8;
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
        require(_timestamp <= _nextTimestamp, "last batch timestamp too high");
        require(_batchDataPos + 32 <= msg.data.length, "wrong length");
        //input msgdata hash, queue hash
        bytes32 inputHash = keccak256(abi.encodePacked(keccak256(msg.data[4:]), _queueHashes));
        uint64 _chainSize = _chain.append(inputHash);
        lastTimestamp = _timestamp;
        emit TransactionAppended(msg.sender, _chainSize - 1, _queueStartIndex, _queueNum, inputHash);
    }

    function chainHeight() public view returns (uint64) {
        return addressResolver.rollupInputChainContainer().chainSize();
    }

    function totalQueue() public view returns (uint64) {
        return uint64(queuedTxInfos.length);
    }

    function getInputHash(uint64 _inputIndex) public view returns (bytes32) {
        return addressResolver.rollupInputChainContainer().get(_inputIndex);
    }

    function getQueueTxInfo(uint64 _queueIndex) public view returns (bytes32, uint64) {
        require(_queueIndex < queuedTxInfos.length, "queue index over capacity");
        QueueTxInfo storage info = queuedTxInfos[_queueIndex];
        return (info.transactionHash, info.timestamp);
    }
}
