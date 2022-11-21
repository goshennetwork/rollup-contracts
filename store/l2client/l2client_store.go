package l2client

import (
	"encoding/binary"

	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type Store struct {
	store schema.KeyValueDB
}

func NewStore(db schema.KeyValueDB) *Store {
	return &Store{
		store: db,
	}
}

func NewMemStore() *Store {
	return &Store{
		store: overlaydb.NewOverlayDB(storage.NewFakeDB()),
	}
}

func (self *Store) GetHeadExecutedQueueBlock() *schema.ChainedEnqueueBlockInfo {
	v, err := self.store.Get(schema.CurrentQueueBlockKey)
	utils.Ensure(err)
	if len(v) == 0 {
		return &schema.ChainedEnqueueBlockInfo{TotalEnqueuedTx: 0, PrevEnqueueBlock: 0}
	}
	b := &schema.ChainedEnqueueBlockInfo{}
	err = b.Deserialization(codec.NewZeroCopySource(v))
	utils.Ensure(err)
	return b
}

func (self *Store) StoreHeadExecutedQueueBlock(headQueue *schema.ChainedEnqueueBlockInfo) {
	self.store.Put(schema.CurrentQueueBlockKey, codec.SerializeToBytes(headQueue))
}

func (self *Store) StoreTotalCheckedBatchNum(batchNum uint64) {
	self.store.Put(schema.L2ClientCheckBatchNumKey, codec.NewZeroCopySink(nil).WriteUint64(batchNum).Bytes())
}

func (self *Store) GetTotalCheckedBatchNum() uint64 {
	v, err := self.store.Get(schema.L2ClientCheckBatchNumKey)
	utils.Ensure(err)
	if len(v) == 0 {
		return 0
	}
	d, err := codec.NewZeroCopySource(v).ReadUint64()
	utils.Ensure(err)
	return d
}

func (self *Store) StoreCheckedBlockNum(batchIndex, blockNum uint64) {
	self.store.Put(genBatchIndexKey(batchIndex), codec.NewZeroCopySink(nil).WriteUint64(blockNum).Bytes())
}

func (self *Store) GetTotalCheckedBlockNum(batchIndex uint64) uint64 {
	v, err := self.store.Get(genBatchIndexKey(batchIndex))
	utils.Ensure(err)
	if len(v) == 0 { //genesis block do not need to check
		return 1
	}
	d, err := codec.NewZeroCopySource(v).ReadUint64()
	utils.Ensure(err)
	return d
}

func (self *Store) StoreReadStorageProof(batchIndex uint64, proofs [][]byte) {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteUint64(uint64(len(proofs)))
	for _, proof := range proofs {
		sink.WriteVarBytes(proof)
	}
	self.store.Put(genProofKey(batchIndex), sink.Bytes())
}

func (self *Store) GetReadStorageProof(batchIndex uint64) [][]byte {
	v, err := self.store.Get(genProofKey(batchIndex))
	utils.Ensure(err)
	if len(v) == 0 {
		return nil
	}
	source := codec.NewZeroCopySource(v)
	num, eof := source.NextUint64()
	if eof {
		return nil
	}
	proofs := make([][]byte, 0)
	for i := uint64(0); i < num; i++ {
		data, _, ill, eof := source.NextVarBytes()
		if ill || eof {
			return nil
		}
		proofs = append(proofs, data)
	}
	return proofs
}

func genBatchIndexKey(batchIndex uint64) []byte {
	var b [9]byte
	b[0] = schema.L2ClientCheckBlockNumPrefix
	binary.BigEndian.PutUint64(b[1:], batchIndex)
	return b[:]
}

func genProofKey(batchIndex uint64) []byte {
	var b [9]byte
	b[0] = schema.L2ClientProofPrefix
	binary.BigEndian.PutUint64(b[1:], batchIndex)
	return b[:]
}
