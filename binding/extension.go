package binding

import (
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils/codec"
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
	funcSelecter := RollupInputChainAbi().Methods["appendBatch"].ID()
	return append(funcSelecter, self.Encode()...)
}

func (self *RollupInputBatches) Encode() []byte {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteUint64BE(self.QueueNum)
	sink.WriteUint64BE(self.QueueStart)
	sink.WriteUint64BE(self.BatchNum)
	sink.WriteUint64BE(self.Batch0Time)
	for _, diff := range self.BatchLeftTimeDiff {
		sink.WriteUint32BE(diff)
	}
	sink.WriteBytes(self.BatchesData)
	return sink.Bytes()
}

// AppendBatch sends a appendBatch transaction in the solidity contract
func (_a *RollupInputChain) AppendInputBatches(batches *RollupInputBatches) *contract.Txn {
	txn :=  _a.c.Txn("appendBatch")
	txn.Data = batches.Calldata()

	return txn
}
