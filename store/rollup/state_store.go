package rollup

import (
	"encoding/binary"
	"errors"

	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type StateChain struct {
	store schema.KeyValueDB
}

func NewStateStore(db schema.KeyValueDB) *StateChain {
	return &StateChain{
		store: db,
	}
}

func NewStateMemStore() *StateChain {
	return &StateChain{
		store: overlaydb.NewOverlayDB(storage.NewFakeDB()),
	}
}

// update info in memory
func (self *StateChain) StoreBatchInfo(states ...*binding.StateBatchAppendedEvent) error {
	info := self.GetInfo()
	for _, state := range states {
		blockNumber := state.Raw.BlockNumber
		if !(info.LastEventBlock < blockNumber) && !(info.LastEventBlock == blockNumber && info.LastEventIndex < state.Raw.LogIndex) {
			return errors.New("older state event found") //may happen when roll back happen, just return, re sync to check l1 block number
		}

		// rollback happend if info.TotalSize > state.StartIndex
		if info.TotalSize < state.StartIndex { // when a gap appear, maybe one mid block is rolled back
			return errors.New("wired state found")
		}
		self.putStateBatchInfo(state)
		info.TotalSize = state.StartIndex + uint64(len(state.BlockHash))
		info.LastEventBlock = state.Raw.BlockNumber
		info.LastEventIndex = state.Raw.LogIndex
	}

	self.StoreInfo(info)
	return nil
}

func (self *StateChain) putStateBatchInfo(states *binding.StateBatchAppendedEvent) {
	for i, v := range states.BlockHash {
		index := states.StartIndex + uint64(i)
		self.store.Put(genStateBatchKey(index), codec.SerializeToBytes(&schema.RollupStateBatchInfo{Index: index, Proposer: states.Proposer, Timestamp: states.Timestamp, BlockHash: v}))
	}
}

func (self *StateChain) GetState(index uint64) (*schema.RollupStateBatchInfo, error) {
	v, err := self.store.Get(genStateBatchKey(index))
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, schema.ErrNotFound
	}
	codecs := &schema.RollupStateBatchInfo{}
	if err := codecs.Deserialization(codec.NewZeroCopySource(v)); err != nil {
		return nil, err
	}
	return codecs, nil
}

func (self *StateChain) GetLastL1BlockHeight() (uint64, error) {
	v, err := self.store.Get(schema.RollupStateLastL1BlockHeightKey)
	if err != nil {
		return 0, err
	}
	if len(v) == 0 {
		return 0, nil
	}
	return codec.NewZeroCopySource(v).ReadUint64()
}

func (self *StateChain) StoreLastL1BlockHeight(lastEndHeight uint64) {
	self.store.Put(schema.RollupStateLastL1BlockHeightKey, codec.NewZeroCopySink(nil).WriteUint64(lastEndHeight).Bytes())
}

func (self *StateChain) StoreInfo(info *schema.StateChainInfo) {
	self.store.Put(schema.CurrentRollupStateChainInfoKey, codec.SerializeToBytes(info))
}

func (self *StateChain) GetInfo() *schema.StateChainInfo {
	v, err := self.store.Get(schema.CurrentRollupStateChainInfoKey)
	utils.Ensure(err)
	if len(v) == 0 { // not exist
		return &schema.StateChainInfo{TotalSize: 0}
	}
	info := new(schema.StateChainInfo)
	err = info.Deserialization(codec.NewZeroCopySource(v))
	utils.Ensure(err)
	return info
}

func genStateBatchKey(batchIndex uint64) []byte {
	var b [9]byte
	b[0] = schema.StateBatchPrefix
	binary.BigEndian.PutUint64(b[1:], batchIndex)
	return b[:]
}
