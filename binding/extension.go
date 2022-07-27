package binding

import (
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
)

// format: batchIndex(uint64)+ queueNum(uint64) + queueStartIndex(uint64) + subBatchNum(uint64) + subBatch0Time(uint64) +
// subBatchLeftTimeDiff([]uint32) + batchesData
// batchesData: version(0) + rlp([][]transaction)
type RollupInputBatches struct {
	//BatchIndex ignored when calc hash, because its useless in l2 system
	BatchIndex uint64
	QueueNum   uint64
	QueueStart uint64
	SubBatches []*SubBatch
}

type SubBatch struct {
	Timestamp uint64
	Txs       []*types.Transaction
}

func (self *RollupInputBatches) Calldata() []byte {
	//function appendInputBatch() public
	funcSelecter := RollupInputChainAbi().Methods["appendInputBatch"].ID()
	return append(funcSelecter, self.Encode()...)
}

// AppendBatch sends a appendBatch transaction in the solidity contract
func (_a *RollupInputChain) AppendInputBatches(batches *RollupInputBatches) *contract.Txn {
	txn := _a.c.Txn("appendInputBatch")
	txn.Data = batches.Calldata()

	return txn
}

func (self *RollupInputBatches) EncodeWithoutIndex() []byte {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteUint64BE(self.QueueNum).WriteUint64BE(self.QueueStart)
	batchNum := uint64(len(self.SubBatches))
	if batchNum < 1 {
		return sink.WriteUint64BE(0).Bytes()
	}
	sink.WriteUint64BE(batchNum).WriteUint64BE(self.SubBatches[0].Timestamp)
	txes := [][]*types.Transaction{self.SubBatches[0].Txs}
	for i := 1; i < len(self.SubBatches); i++ {
		b := self.SubBatches[i]
		prev := self.SubBatches[i-1].Timestamp
		//equal happens when l1 block timestamp not refresh yet
		utils.EnsureTrue(b.Timestamp >= prev && prev+math.MaxUint32 > b.Timestamp)
		timeDiff := b.Timestamp - prev
		sink.WriteUint32BE(uint32(timeDiff))
		txes = append(txes, b.Txs)
	}
	rlpTx, err := rlp.EncodeToBytes(txes)
	if err != nil {
		panic(err)
	}
	sink.WriteByte(0) // version 0
	sink.WriteBytes(rlpTx)
	return sink.Bytes()
}

func (self *RollupInputBatches) Encode() []byte {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteUint64BE(self.BatchIndex)
	dataWithoutIndex := self.EncodeWithoutIndex()
	return append(sink.Bytes(), dataWithoutIndex...)
}

// InputBatchHash get input hash, ignore first 8 byte
func (self *RollupInputBatches) InputBatchHash() web3.Hash {
	return crypto.Keccak256Hash(self.EncodeWithoutIndex())
}

func (self *RollupInputBatches) InputHash(queueHash web3.Hash) web3.Hash {
	return crypto.Keccak256Hash(self.InputBatchHash().Bytes(), queueHash.Bytes())
}

func safeAdd(x, y uint64) uint64 {
	utils.EnsureTrue(y < math.MaxUint64-x)
	return x + y
}

func (self *RollupInputBatches) DecodeWithoutIndex(b []byte) error {
	reader := codec.NewZeroCopyReader(b)
	self.QueueNum = reader.ReadUint64BE()
	self.QueueStart = reader.ReadUint64BE()
	batchNum := reader.ReadUint64BE()
	if batchNum == 0 {
		//check length
		if reader.Len() != 0 {
			return fmt.Errorf("wrong b length")
		}
		return reader.Error()
	}
	batchTime := reader.ReadUint64BE()
	batchesTime := []uint64{batchTime}
	for i := uint64(0); i < batchNum-1; i++ {
		batchTime = safeAdd(batchTime, uint64(reader.ReadUint32BE()))
		if reader.Error() != nil {
			return reader.Error()
		}
		batchesTime = append(batchesTime, batchTime)
	}

	version := reader.ReadUint8()
	if version != 0 {
		return fmt.Errorf("unknown batch version: %d", version)
	}

	rawBatchesData := reader.ReadBytes(reader.Len())
	if reader.Error() != nil {
		return reader.Error()
	}

	txs := make([][]*types.Transaction, 0)
	err := rlp.DecodeBytes(rawBatchesData, &txs)
	if err != nil {
		return err
	}

	if uint64(len(txs)) != batchNum {
		return fmt.Errorf("inconsistent batch num with tx")
	}
	for i, b := range txs {
		self.SubBatches = append(self.SubBatches, &SubBatch{
			Timestamp: batchesTime[i],
			Txs:       b,
		})
	}

	return nil
}

// decode batch info and check in info correctness
func (self *RollupInputBatches) Decode(b []byte) error {
	reader := codec.NewZeroCopyReader(b[:8])
	self.BatchIndex = reader.ReadUint64BE()
	return self.DecodeWithoutIndex(b[8:])
}
