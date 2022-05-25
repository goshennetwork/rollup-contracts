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

contract RollupInputChain is IRollupInputChain, Initializable {
    uint256 public constant MIN_ROLLUP_TX_GAS = 100000;
    uint256 public constant MAX_ROLLUP_TX_SIZE = 50000;
    uint256 public constant MAX_CROSS_LAYER_TX_SIZE = 10000;
    uint256 public constant GAS_PRICE = 1;
    uint256 public constant VALUE = 0;

    uint64 public maxEnqueueTxGasLimit;
    uint64 public maxCrossLayerTxGasLimit;
    uint64 public l2ChainId;
    address public l1CrossLayerEOA;

    uint64 public override lastTimestamp;

    IAddressResolver addressResolver;

    //store L1 -> L2 tx
    struct QueueTxInfo {
        bytes32 signHash;
        uint256 r;
        uint256 s;
        uint64 v;
        bool isL1CrossLayerWitness;
        uint64 timestamp;
    }

    QueueTxInfo[] queuedTxInfos;
    // index of the first queue element not yet included
    uint64 public override pendingQueueIndex;

    function initialize(
        address _addressResolver,
        uint64 _maxTxGasLimit,
        uint64 _maxCrossLayerTxGasLimit,
        uint64 _l2ChainId,
        address _l1CrossLayerEOA
    ) public initializer {
        addressResolver = IAddressResolver(_addressResolver);
        maxEnqueueTxGasLimit = _maxTxGasLimit;
        maxCrossLayerTxGasLimit = _maxCrossLayerTxGasLimit;
        l2ChainId = _l2ChainId;
        l1CrossLayerEOA = _l1CrossLayerEOA;
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
        bool _isL1CrossLayerWitness = false;
        bytes32 _signTxHash;
        // L1 EOA is equal to L2 EOA, but L1 contract is not except L1CrossLayerWitness
        address sender;
        if (msg.sender == tx.origin) {
            sender = msg.sender;
            require(_data.length <= MAX_ROLLUP_TX_SIZE, "too large Tx data size");
            _signTxHash = calculateSignHash(_nonce, _gasLimit, _target, _data);
        } else {
            require(msg.sender == address(addressResolver.l1CrossLayerWitness()), "contract can not enqueue L2 Tx");
            require(_data.length <= MAX_CROSS_LAYER_TX_SIZE, "too large cross layer Tx data size");
            sender = Constants.L1_CROSS_LAYER_WITNESS;
            _gasLimit = maxCrossLayerTxGasLimit;
            _signTxHash = calculateSignHash(_nonce, _gasLimit, _target, _data);
            _isL1CrossLayerWitness = true;
        }
        require(_gasLimit <= maxEnqueueTxGasLimit, "too high Tx gas limit");
        require(_gasLimit >= MIN_ROLLUP_TX_GAS, "too low Tx gas limit");

        uint64 _now = uint64(block.timestamp);
        queuedTxInfos.push(
            QueueTxInfo({
                timestamp: _now,
                signHash: _signTxHash,
                r: r,
                s: s,
                v: v,
                isL1CrossLayerWitness: _isL1CrossLayerWitness
            })
        );
        //do not need to reveal default value
        emit TransactionEnqueued(
            uint64(queuedTxInfos.length - 1),
            sender,
            _target,
            _gasLimit,
            _data,
            _nonce,
            r,
            s,
            v,
            _now
        );
    }

    //encode tx params: sender, to, gasLimit, data, nonce, r,s,v and gasPrice(1 GWEI), value(0), chainId
    //sender used to recognize tx from L1CrossLayerWitness
    function calculateSignHash(
        uint64 _nonce,
        uint64 _gasLimit,
        address _target,
        bytes memory _data
    ) public view returns (bytes32) {
        return
            keccak256(
                abi.encodePacked(
                    RLPWriter.writeUint(uint256(_nonce)),
                    RLPWriter.writeUint(GAS_PRICE),
                    RLPWriter.writeUint(uint256(_gasLimit)),
                    RLPWriter.writeAddress(_target),
                    RLPWriter.writeUint(VALUE),
                    RLPWriter.writeBytes(_data),
                    RLPWriter.writeUint(l2ChainId),
                    bytes1(0x80),
                    bytes1(0x80)
                )
            );
    }

    function completeQueueInfo(
        uint64 _queueStartIndex,
        uint64 _queueNum,
        uint64 _batchDataPos
    ) internal returns (bytes32, uint64) {
        uint64 v;
        bytes32 r;
        bytes32 s;
        uint64 _crossTxNum;
        uint256 len = (32 + 8) * _queueNum;
        bytes memory _queueHash = new bytes(len);
        uint64 _offset = 0;
        for (uint256 i = 0; i < _queueNum; i++) {
            QueueTxInfo storage info = queuedTxInfos[_queueStartIndex + i];
            if (info.isL1CrossLayerWitness) {
                _crossTxNum++;
                //check signature
                assembly {
                    r := calldataload(_batchDataPos)
                    s := calldataload(add(_batchDataPos, 32))
                    v := shr(192, calldataload(add(_batchDataPos, 64)))
                }
                uint64 _pureV = v - 32 - ((l2ChainId * 2));
                require(_pureV <= 1, "wrong v");
                address signer = ecrecover(info.signHash, uint8(_pureV), r, s);
                require(signer == l1CrossLayerEOA, "wrong sig");
                _batchDataPos += 32 + 32 + 8;
                //now complete signature
                info.r = uint256(r);
                info.s = uint256(s);
                info.v = v;
            }
            bytes32 txHash = keccak256(abi.encode(info.signHash, info.r, info.s, info.v, info.isL1CrossLayerWitness));
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
        return (result, _crossTxNum);
    }

    // format: queueNum(uint64) + queueStart(uint64) + batchNum(uint64) + batch0Time(uint64) +
    // batchLeftTimeDiff([]uint32) + crossLayerSig([]bytes72(uint256 R || uint256 S || uint64 V)) + batchesData
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
        require(_timestamp > lastTimestamp && _timestamp < block.timestamp, "wrong batch timestamp");
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
        require(_timestamp < _nextTimestamp, "last batch timestamp too high");

        //now check signature and feed info,then calc enqueue tx hash
        (bytes32 _queueHash, uint64 _crossLayerTxNum) = completeQueueInfo(_queueStartIndex, _queueNum, _batchDataPos);
        _batchDataPos += _crossLayerTxNum * (32 + 32 + 8); //every crossLayerTx have 72 byte  spece to complete signature
        require(_batchDataPos + 32 <= msg.data.length, "wrong length");
        //input msgdata hash, queue hash
        bytes32 inputHash = keccak256(abi.encodePacked(keccak256(msg.data[4:]), _queueHash));
        _chain.append(inputHash);
        lastTimestamp = _timestamp;
        emit TransactionAppended(msg.sender, _queueStartIndex, _queueNum, _chain.chainSize() - 1, inputHash);
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

    //    function getQueueTxInfo(uint64 _queueIndex) public view returns (bytes32, uint64) {
    //        require(_queueIndex < queuedTxInfos.length, "queue index over capacity");
    //        QueueTxInfo storage info = queuedTxInfos[_queueIndex];
    //        return (info.transactionHash, info.timestamp);
    //    }
}
