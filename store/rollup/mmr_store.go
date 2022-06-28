package rollup

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/merkle"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type MMR struct {
	dataPrefix byte
	treeKey    []byte
	store      schema.KeyValueDB
	tree       *merkle.CompactMerkleTree
}

func NewL1MMR(db schema.KeyValueDB) *MMR {
	return &MMR{
		dataPrefix: schema.L1MMRDataPrefix,
		treeKey:    schema.L1CompactMerkleTreeKey,
		store:      db,
	}
}

func NewL2MMR(db schema.KeyValueDB) *MMR {
	return &MMR{
		dataPrefix: schema.L2MMRDataPrefix,
		treeKey:    schema.L2CompactMerkleTreeKey,
		store:      db,
	}
}

func (self *MMR) GetL1CompactMerkleTree() (uint64, []web3.Hash, error) {
	v, err := self.store.Get(schema.L1CompactMerkleTreeKey)
	if err != nil {
		return 0, []web3.Hash{}, err
	}
	if len(v) == 0 {
		return 0, []web3.Hash{}, nil
	}
	return schema.DeserializeCompactMerkleTree(v)
}

// update info in memory
func (self *MMR) StoreBatchInfo(states ...*binding.StateBatchAppendedEvent) {
	info := self.GetInfo()
	for _, state := range states {
		blockNumber := state.Raw.BlockNumber
		utils.EnsureTrue(info.LastEventBlock < blockNumber || (info.LastEventBlock == blockNumber && info.LastEventIndex < state.Raw.LogIndex))

		// rollback happend if info.TotalSize > state.StartIndex
		utils.EnsureTrue(info.TotalSize >= state.StartIndex)
		self.putStateBatchInfo(state)
		info.TotalSize = state.StartIndex + uint64(len(state.BlockHash))
		info.LastEventBlock = state.Raw.BlockNumber
		info.LastEventIndex = state.Raw.LogIndex
	}

	self.StoreInfo(info)
}

func (self *MMR) putStateBatchInfo(states *binding.StateBatchAppendedEvent) {
	for i, v := range states.BlockHash {
		index := states.StartIndex + uint64(i)
		self.store.Put(genStateBatchKey(index), codec.SerializeToBytes(&schema.RollupStateBatchInfo{Index: index, Proposer: states.Proposer, Timestamp: states.Timestamp, BlockHash: v}))
	}
}

func (self *MMR) GetState(index uint64) (*schema.RollupStateBatchInfo, error) {
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

func (self *MMR) GetLastL1BlockHeight() (uint64, error) {
	v, err := self.store.Get(schema.RollupStateLastL1BlockHeightKey)
	if err != nil {
		return 0, err
	}
	if len(v) == 0 {
		return 0, nil
	}
	return codec.NewZeroCopySource(v).ReadUint64()
}

func (self *MMR) StoreLastL1BlockHeight(lastEndHeight uint64) {
	self.store.Put(schema.RollupStateLastL1BlockHeightKey, codec.NewZeroCopySink(nil).WriteUint64(lastEndHeight).Bytes())
}

func (self *MMR) StoreInfo(info *schema.MMRInfo) {
	self.store.Put(schema.CurrentRollupMMRInfoKey, codec.SerializeToBytes(info))
}

type MMRHashStore struct {
	store schema.KeyValueDB
}

func (self *MMRHashStore) Append(hash []web3.Hash) error {

}

// HashStore is an interface for persist hash
type HashStore interface {
	Append(hash []web3.Hash) error
	Flush() error
	Close()
	GetHash(pos uint64) (web3.Hash, error)
}
