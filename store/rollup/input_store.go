package rollup

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type InputChain struct {
	store schema.KeyValueDB
}

func NewInputStore(db schema.KeyValueDB) *InputChain {
	return &InputChain{
		store: db,
	}
}

func NewInputMemStore() *InputChain {
	return &InputChain{
		store: overlaydb.NewOverlayDB(storage.NewFakeDB()),
	}
}

func (self *InputChain) putInfo(info *schema.InputChainInfo) {
	self.store.Put(schema.CurrentRollupInputChainInfoKey, codec.SerializeToBytes(info))
}

func (self *InputChain) StoreQueueSize(size uint64) {
	self.store.Put(schema.CurrentQueueSizeKey, codec.NewZeroCopySink(nil).WriteUint64BE(size).Bytes())
}

func (self *InputChain) QueueSize() uint64 {
	data, err := self.store.Get(schema.CurrentQueueSizeKey)
	utils.Ensure(err)
	if len(data) == 0 { // not exist
		return 0
	}
	return binary.BigEndian.Uint64(data)
}

func (self *InputChain) GetInfo() *schema.InputChainInfo {
	data, err := self.store.Get(schema.CurrentRollupInputChainInfoKey)
	utils.Ensure(err)
	if len(data) == 0 { // not exist
		return &schema.InputChainInfo{0, 0}
	}
	source := codec.NewZeroCopySource(data)
	bed := &schema.InputChainInfo{}
	err = bed.DeSerialization(source)
	utils.Ensure(err)
	return bed
}

func (self *InputChain) StoreEnqueuedTransaction(queues ...*binding.TransactionEnqueuedEvent) error {
	size := self.QueueSize()
	for _, queue := range queues {
		if queue.QueueIndex > size { // check consistent, wired situation will happen when l1 roll back some block,but old queue is permitted.
			return fmt.Errorf("wrong queue index, expect: %d, found: %d", size, queue.QueueIndex)
		}
		txn := schema.EnqueuedTransactionFromEvent(queue)

		if size != queue.QueueIndex { //old queue, check whether is the same, otherwise rollback happended
			oldTxn, err := self.GetEnqueuedTransaction(queue.QueueIndex)
			if err != nil {
				return err
			}
			if bytes.Equal(codec.SerializeToBytes(oldTxn), codec.SerializeToBytes(txn)) {
				continue
			} else { // find out inconsistent
				return fmt.Errorf("enqueue tx inconsistant, queueIndex: %d")
			}
		}
		self.putEnqueuedTransaction(txn)
		size += 1
	}
	self.StoreQueueSize(size)
	return nil
}

func (self *InputChain) GetAppendedTransaction(index uint64) (*schema.AppendedTransaction, error) {
	key := genRollupInputBatchKey(index)
	data, err := self.store.Get(key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 { // not found
		return nil, schema.ErrNotFound
	}
	source := codec.NewZeroCopySource(data)
	txn := &schema.AppendedTransaction{}
	err = txn.Deserialization(source)
	if err != nil {
		return nil, err
	}
	return txn, nil
}

func (self *InputChain) GetEnqueuedTransaction(queueIndex uint64) (*schema.EnqueuedTransaction, error) {
	data, err := self.store.Get(genQueueElementKey(queueIndex))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, schema.ErrNotFound
	}
	source := codec.NewZeroCopySource(data)
	enqueued := &schema.EnqueuedTransaction{}
	err = enqueued.Deserialization(source)
	if err != nil {
		return nil, err
	}
	return enqueued, nil
}

func (self *InputChain) GetEnqueuedTransactions(startIndex uint64, num uint64) ([]*schema.EnqueuedTransaction, error) {
	queues := make([]*schema.EnqueuedTransaction, num)
	var err error
	for i := uint64(0); i < num; i++ {
		queues[i], err = self.GetEnqueuedTransaction(startIndex + i)
		if err != nil {
			return nil, err
		}
	}
	return queues, err
}

//will update info in memory
func (self *InputChain) StoreSequencerBatches(queueSize uint64, batches ...*binding.InputBatchAppendedEvent) error {
	info := self.GetInfo()
	for _, batch := range batches {
		if batch.Index > info.TotalBatches { //happen when roll back, old is permitted
			return fmt.Errorf("wrong batch index, expect: %d, found: %d", info.TotalBatches, batch.Index)
		}

		// check the queue info
		if batch.StartQueueIndex > info.PendingQueueIndex { // wired batch
			return fmt.Errorf("wrong start queue index, expect:%d, found:%d", info.PendingQueueIndex, batch.StartQueueIndex)
		}
		if batch.StartQueueIndex+batch.QueueNum > queueSize { // reach unlocated queue, wired
			return fmt.Errorf("wired batch or queue found, local queue size: %d, batch queue num: %d, pending queue index: %d", queueSize, batch.QueueNum, info.PendingQueueIndex)
		}
		if batch.Index == info.TotalBatches {
			//now should consistent
			if batch.StartQueueIndex != info.PendingQueueIndex {
				//wired, mayble should never happen?
				panic(1)
			}
		}
		txn := &schema.AppendedTransaction{
			Proposer:        batch.Proposer,
			Index:           batch.Index,
			StartQueueIndex: batch.StartQueueIndex,
			QueueNum:        batch.QueueNum,
			InputHash:       batch.InputHash,
		}
		if batch.Index < info.TotalBatches { // old batch when roll back happen, check the info
			oldtx, err := self.GetAppendedTransaction(batch.Index)
			if err != nil {
				return err
			}
			if bytes.Equal(codec.SerializeToBytes(oldtx), codec.SerializeToBytes(txn)) {
				continue
			} else { // inconsistent just rollback
				return fmt.Errorf("find inconsistent input batch, index: %d", batch.Index)
			}
		}
		self.store.Put(genRollupInputBatchKey(batch.Index), codec.SerializeToBytes(txn))
		info.TotalBatches = batch.Index + 1
		info.PendingQueueIndex += batch.QueueNum
	}
	self.putInfo(info)
	return nil
}

//returned data already trim function selector in calldata
func (self *InputChain) GetSequencerBatchData(index uint64) ([]byte, error) {
	v, err := self.store.Get(genRollupInputBatchDataKey(index))
	utils.Ensure(err)
	if len(v) == 0 {
		return nil, schema.ErrNotFound
	}
	return v, nil
}

func (self *InputChain) StoreSequencerBatchData(txs []*web3.Transaction, indexes []uint64) {
	if len(txs) != len(indexes) {
		panic(fmt.Errorf("wrong num of batch data and indexes, %d vs %d", len(txs), len(indexes)))
	}
	for i, tx := range txs {
		self.store.Put(genRollupInputBatchDataKey(indexes[i]), tx.Input[4:])
	}
}

//write enqueue element to db.
func (self *InputChain) putEnqueuedTransaction(txn *schema.EnqueuedTransaction) {
	self.store.Put(genQueueElementKey(txn.QueueIndex), codec.SerializeToBytes(txn))
}

func (self *InputChain) GetNumPendingQueueElements() (uint64, error) {
	info := self.GetInfo()

	return info.TotalBatches - info.PendingQueueIndex, nil
}

func genRollupInputBatchKey(batchIndex uint64) []byte {
	key := make([]byte, 9, 9)
	key[0] = schema.RollupInputBatchKey
	binary.BigEndian.PutUint64(key[1:], batchIndex)
	return key
}

func genQueueElementKey(queueIndex uint64) []byte {
	key := make([]byte, 9, 9)
	key[0] = schema.SequencerQueuePrefix
	binary.BigEndian.PutUint64(key[1:], queueIndex)
	return key
}

func genRollupInputBatchDataKey(batchIndex uint64) []byte {
	key := make([]byte, 9, 9)
	key[0] = schema.RollupInputBatchDataKey
	binary.BigEndian.PutUint64(key[1:], batchIndex)
	return key
}
