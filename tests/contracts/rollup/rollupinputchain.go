package rollup

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/laizy/web3/utils/codec"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/utils/common"
	"github.com/laizy/web3/utils/common/math"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/build"
)

// format: queueNum(uint64) + queueStart(uint64) + batchNum(uint64) + batch0Time(uint64) +
// batchLeftTimeDiff([]uint32) + batchesData
type RollupInputBatches struct {
	QueueNum          uint64
	QueueStart        uint64
	BatchNum          uint64
	Batch0Time        uint64
	BatchLeftTimeDiff []uint32
	BatchesData       []byte
}

func (self *RollupInputBatches) Calldata() []byte {
	//function appendBatch() public
	funcSelecter := build.RollupInputChainAbi().Methods["appendBatch"].ID()
	return append(funcSelecter, self.Encode()...)
}

func (self *RollupInputBatches) Encode() []byte {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteBytes(uint64ToBytes(self.QueueNum))
	sink.WriteBytes(uint64ToBytes(self.QueueStart))
	sink.WriteBytes(uint64ToBytes(self.BatchNum))
	sink.WriteBytes(uint64ToBytes(self.Batch0Time))
	for _, diff := range self.BatchLeftTimeDiff {
		sink.WriteBytes(uint32ToBytes(diff))
	}
	sink.WriteBytes(self.BatchesData)
	return sink.Bytes()
}

func uint64ToBytes(i uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], i)
	return b[:]
}

func uint32ToBytes(i uint32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], i)
	return b[:]
}

type RollupInputChain web3.Address

//function enqueue(
//        address _target,
//        uint64 _gasLimit,
//        bytes memory _data
//    ) public
func (self RollupInputChain) Enqueue(sender evm.AccountRef, vm *evm.EVM, target web3.Address, gasLimit uint64, data []byte) error {
	input := build.RollupInputChainAbi().Methods["enqueue"].MustEncodeIDAndInput(target, gasLimit, data)
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return fmt.Errorf(web3.DecodeRevert(ret))
	}
	return err
}

//function appendBatch() public
func (self RollupInputChain) AppendBatch(sender evm.AccountRef, vm *evm.EVM, batches *RollupInputBatches) error {
	input := batches.Calldata()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf(web3.DecodeRevert(ret))
	}
	return err
}

//function chainHeight() public view returns (uint64)
func (self RollupInputChain) ChainHeight(sender evm.AccountRef, vm *evm.EVM) (out uint64, err error) {
	method := build.RollupInputChainAbi().Methods["chainHeight"]
	ret, _, err := vm.Call(sender, web3.Address(self), method.ID(), math.MaxUint64, new(big.Int))
	if err != nil {
		return 0, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return 0, err
	}
	err = contracts.Decode(i, &out)
	return
}

//function pendingQueueIndex() external view returns (uint64)
func (self RollupInputChain) PendingQueueIndex(sender evm.AccountRef, vm *evm.EVM) (out uint64, err error) {
	method := build.RollupInputChainAbi().Methods["pendingQueueIndex"]
	ret, _, err := vm.Call(sender, web3.Address(self), method.ID(), math.MaxUint64, new(big.Int))
	if err != nil {
		return 0, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return 0, err
	}
	err = contracts.Decode(i, &out)
	return
}

//function lastTimestamp() public view returns (uint64)
func (self RollupInputChain) LastTimestamp(sender evm.AccountRef, vm *evm.EVM) (out uint64, err error) {
	method := build.RollupInputChainAbi().Methods["lastTimestamp"]
	ret, _, err := vm.Call(sender, web3.Address(self), method.ID(), math.MaxUint64, new(big.Int))
	if err != nil {
		return 0, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return 0, err
	}
	err = contracts.Decode(i, &out)
	return
}

//function getQueueTxInfo(uint64 _queueIndex) public view returns (bytes32, uint64)
func (self RollupInputChain) GetQueueTxInfo(sender evm.AccountRef, vm *evm.EVM, queueIndex uint64) (txHash common.Hash, timestamp uint64, err error) {
	method := build.RollupInputChainAbi().Methods["getQueueTxInfo"]
	input := method.MustEncodeIDAndInput(queueIndex)
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return [32]byte{}, 0, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return [32]byte{}, 0, err
	}
	err = contracts.Decode(i, &txHash, &timestamp)
	return
}
