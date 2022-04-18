// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import { Types } from "../libraries/Types.sol";
import "../interfaces/IStakingManager.sol";
import { ICanonicalTransactionChain } from "../interfaces/ICanonicalTransactionChain.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IChainStorageContainer.sol";

contract CanonicalTransactionChain is ICanonicalTransactionChain {
    using Types for Types.QueueElement;
    IAddressResolver addressResolver;

    //store L1 -> L2 tx
    Types.QueueElement[] queueElements;
    // index of the first queue element not yet included
    uint64 public override pendingQueueIndex;

    constructor(address _addressResolver) {
        addressResolver = IAddressResolver(_addressResolver);
    }

    function enqueue(
        address _target,
        uint256 _gasLimit,
        bytes memory _data
    ) public {
        //We guarantee that the L2 EOA is L1 EOA, and L1 contract can't be L2 EOA except l1 crossDomainContract which is used
        //when l1 bridge try to  enqueue tx to l2
        if (msg.sender != tx.origin) {
            //the l1 bridge use cross Domain contract to enqueue tx to l2.We only allow contract as l2 EOA when sender is this contract,
            require(
                msg.sender == addressResolver.l1CrossLayerMessageWitness(),
                "contract can't act as EOA in L2 except l1 crossDomain contract"
            );
        }
        //todo: maybe need more tx params, such as tip fee,value
        bytes32 transactionHash = keccak256(abi.encode(msg.sender, _target, _gasLimit, _data));
        uint64 _now = uint64(block.timestamp);
        queueElements.push(Types.QueueElement({ transactionHash: transactionHash, timestamp: _now }));
        emit Enqueued(msg.sender, _target, _gasLimit, _data, uint64(queueElements.length - 1), _now);
    }

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
        bytes32 _queueHashes = keccak256(_queueHash);
        uint64 _sequencedIndex = 4 + 8 + 8; //4byte function selector, 2 uint64
        pendingQueueIndex = _nextPendingQueueIndex;
        //check sequencer timestamp
        uint64 _contextNum;
        assembly {
            _contextNum := shr(192, _sequencedIndex)
        }
        uint64 _timestamp;
        uint64 _lastTimestamp;
        //clear
        _offset = 0;
        for (uint64 i = 0; i < _contextNum; i++) {
            _offset = _sequencedIndex + 8 + 8 * i;
            assembly {
                _timestamp := shr(192, calldataload(_offset))
            }
            if (i == 0) {
                //first
                require(_timestamp > _chain.lastTimestamp(), "start timestamp should be larger than obvious timestamp");
            }
            if (i == _contextNum - 1) {
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
