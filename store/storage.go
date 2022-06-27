package store

import (
	"path/filepath"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/config"
	"github.com/ontology-layer-2/rollup-contracts/merkle"
	"github.com/ontology-layer-2/rollup-contracts/store/l2client"
	"github.com/ontology-layer-2/rollup-contracts/store/overlaydb"
	"github.com/ontology-layer-2/rollup-contracts/store/resolver"
	"github.com/ontology-layer-2/rollup-contracts/store/rollup"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type Storage struct {
	diskdb         schema.PersistStore
	*StorageWriter // this is acturaly read only
}

func NewStorage(diskdb schema.PersistStore, dbPath string) *Storage {
	overlay := &ReadOnlyDB{overlaydb.NewOverlayDB(diskdb)}
	writer := &StorageWriter{
		overlay: overlay,
	}
	l1TreeSize, l1Hashes, err := writer.GetL1CompactMerkleTree()
	utils.Ensure(err)
	l2TreeSize, l2Hashes, err := writer.GetL1CompactMerkleTree()
	utils.Ensure(err)
	l1FileHashStore, err := merkle.NewFileHashStore(dbPath+string(filepath.Separator)+config.DefaultL1MMRFile, l1TreeSize)
	utils.Ensure(err)
	l2FileHashStore, err := merkle.NewFileHashStore(dbPath+string(filepath.Separator)+config.DefaultL2MMRFile, l2TreeSize)
	utils.Ensure(err)
	l1CompactTree := merkle.NewTree(l1TreeSize, l1Hashes, l1FileHashStore)
	l2CompactTree := merkle.NewTree(l2TreeSize, l2Hashes, l2FileHashStore)
	writer.l1CompactMerkleTree = l1CompactTree
	writer.l2CompactMerkleTree = l2CompactTree
	return &Storage{
		diskdb:        diskdb,
		StorageWriter: writer,
	}
}

type StorageWriter struct {
	l1CompactMerkleTree *merkle.CompactMerkleTree
	l2CompactMerkleTree *merkle.CompactMerkleTree
	overlay             KeyValueDBWithCommit
}

func (self *Storage) Writer() *StorageWriter {
	return &StorageWriter{overlay: overlaydb.NewOverlayDB(self.diskdb), l1CompactMerkleTree: self.l1CompactMerkleTree, l2CompactMerkleTree: self.l2CompactMerkleTree}
}

func (self *StorageWriter) InputChain() *rollup.InputChain {
	return rollup.NewInputStore(self.overlay)
}

func (self *StorageWriter) AddressManager() *resolver.AddressManager {
	return resolver.NewStore(self.overlay)
}

func (self *StorageWriter) StateChain() *rollup.StateChain {
	return rollup.NewStateStore(self.overlay)
}

func (self *StorageWriter) L1TokenBridge() *rollup.L1BridgeStore {
	return rollup.NewL1BridgeStore(self.overlay)
}

func (self *StorageWriter) L1CrossLayerWitness() *rollup.L1WitnessStore {
	return rollup.NewL1WitnessStore(self.overlay, self.l1CompactMerkleTree)
}

func (self *StorageWriter) L2TokenBridge() *rollup.L2BridgeStore {
	return rollup.NewL2BridgeStore(self.overlay)
}

func (self *StorageWriter) L2CrossLayerWitness() *rollup.L2WitnessStore {
	return rollup.NewL2WitnessStore(self.overlay, self.l2CompactMerkleTree)
}

func (self *StorageWriter) L2Client() *l2client.Store {
	return l2client.NewStore(self.overlay)
}

func (self *StorageWriter) SetLastSyncedL1Height(lastEndHeight uint64) {
	self.overlay.Put(schema.LastSyncedL1HeightKey, codec.NewZeroCopySink(nil).WriteUint64(lastEndHeight).Bytes())
}

func (self *StorageWriter) SetLastSyncedL1Timestamp(lastTimestamp uint64) {
	self.overlay.Put(schema.LastSyncedL1TimestampKey, codec.NewZeroCopySink(nil).WriteUint64(lastTimestamp).Bytes())
}

// GetLastSyncedL1Timestamp get last synced l1 timestamp, if not exist, return nil
func (self *StorageWriter) GetLastSyncedL1Timestamp() *uint64 {
	v, err := self.overlay.Get(schema.LastSyncedL1TimestampKey)
	utils.Ensure(err)
	if len(v) == 0 {
		return nil
	}
	timestamp, err := codec.NewZeroCopySource(v).ReadUint64()
	utils.Ensure(err)

	return &timestamp
}

func (self *StorageWriter) Commit() {
	self.overlay.CommitTo()
}

func (self *StorageWriter) GetLastSyncedL1Height() uint64 {
	v, err := self.overlay.Get(schema.LastSyncedL1HeightKey)
	utils.Ensure(err)
	if len(v) == 0 {
		return 0
	}
	height, err := codec.NewZeroCopySource(v).ReadUint64()
	utils.Ensure(err)

	return height
}

func (self *StorageWriter) GetL1CompactMerkleTree() (uint64, []web3.Hash, error) {
	v, err := self.overlay.Get(schema.L1CompactMerkleTreeKey)
	if err != nil {
		return 0, []web3.Hash{}, err
	}
	if len(v) == 0 {
		return 0, []web3.Hash{}, nil
	}
	return schema.DeserializeCompactMerkleTree(v)
}

func (self *StorageWriter) GetL2CompactMerkleTree() (uint64, []web3.Hash, error) {
	v, err := self.overlay.Get(schema.L2CompactMerkleTreeKey)
	if err != nil {
		return 0, []web3.Hash{}, err
	}
	if len(v) == 0 {
		return 0, []web3.Hash{}, nil
	}
	return schema.DeserializeCompactMerkleTree(v)
}

func (self *StorageWriter) StoreL1CompactMerkleTree() {
	v := schema.SerializeCompactMerkleTree(self.l1CompactMerkleTree)
	self.overlay.Put(schema.L1CompactMerkleTreeKey, v)
}

func (self *StorageWriter) StoreL2CompactMerkleTree() {
	v := schema.SerializeCompactMerkleTree(self.l2CompactMerkleTree)
	self.overlay.Put(schema.L2CompactMerkleTreeKey, v)
}

type ReadOnlyDB struct {
	schema.KeyValueDB
}

func (self *ReadOnlyDB) Put([]byte, []byte) {
	panic("read only")
}

func (self *ReadOnlyDB) Delete([]byte) {
	panic("read only")
}

func (self *ReadOnlyDB) CommitTo() {
	panic("read only")
}

type KeyValueDBWithCommit interface {
	schema.KeyValueDB
	CommitTo()
}
