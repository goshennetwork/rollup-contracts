package store

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
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

func NewStorage(diskdb schema.PersistStore) *Storage {
	overlay := &ReadOnlyDB{overlaydb.NewOverlayDB(diskdb)}
	writer := &StorageWriter{
		overlay: overlay,
	}
	return &Storage{
		diskdb:        diskdb,
		StorageWriter: writer,
	}
}

type StorageWriter struct {
	overlay KeyValueDBWithCommit
}

func (self *Storage) Writer() *StorageWriter {
	return &StorageWriter{overlay: overlaydb.NewOverlayDB(self.diskdb)}
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
	return rollup.NewL1WitnessStore(self.overlay)
}

func (self *StorageWriter) L2TokenBridge() *rollup.L2BridgeStore {
	return rollup.NewL2BridgeStore(self.overlay)
}

func (self *StorageWriter) L2CrossLayerWitness() *rollup.L2WitnessStore {
	return rollup.NewL2WitnessStore(self.overlay)
}

func (self *StorageWriter) L2Client() *l2client.Store {
	return l2client.NewStore(self.overlay)
}

func (self *StorageWriter) L1MMR() *rollup.MMR {
	return rollup.NewL1MMR(self.overlay)
}

func (self *StorageWriter) L2MMR() *rollup.MMR {
	return rollup.NewL2MMR(self.overlay)
}

func (self *StorageWriter) SetLastSyncedL1Height(lastEndHeight uint64) {
	self.overlay.Put(schema.LastSyncedL1HeightKey, codec.NewZeroCopySink(nil).WriteUint64(lastEndHeight).Bytes())
}

func (self *StorageWriter) SetLastSyncedL1Timestamp(lastTimestamp uint64) {
	self.overlay.Put(schema.LastSyncedL1TimestampKey, codec.NewZeroCopySink(nil).WriteUint64(lastTimestamp).Bytes())
}

func (self *StorageWriter) SetLastSyncedL1Hash(hash web3.Hash) {
	self.overlay.Put(schema.LastSyncedL1Hash, codec.NewZeroCopySink(nil).WriteHash(hash).Bytes())
}

func (self *StorageWriter) SetHighestL1CheckPointInfo(info *schema.L1CheckPointInfo) {
	self.overlay.Put(schema.HighestL1CheckPointInfoKey, codec.SerializeToBytes(info))
}

func (self *StorageWriter) SetPendingL1CheckPointInfo(info *schema.L1CheckPointInfo) {
	self.overlay.Put(schema.PendingL1CheckPointInfoKey, codec.SerializeToBytes(info))
}

func (self *StorageWriter) SetL1DbVersion(version uint64) {
	self.overlay.Put(schema.L1DbVersionKey, codec.NewZeroCopySink(nil).WriteUint64BE(version).Bytes())
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

func (self *StorageWriter) GetL1DbVersion() uint64 {
	v, err := self.overlay.Get(schema.L1DbVersionKey)
	utils.Ensure(err)
	if len(v) == 0 {
		return 0
	}
	return codec.NewZeroCopyReader(v).ReadUint64BE()
}

func (self *StorageWriter) GetHighestL1CheckPointInfo() *schema.L1CheckPointInfo {
	v, err := self.overlay.Get(schema.HighestL1CheckPointInfoKey)
	utils.Ensure(err)
	if len(v) == 0 {
		return nil
	}
	ret := new(schema.L1CheckPointInfo)
	err = ret.Deserialization(codec.NewZeroCopySource(v))
	utils.Ensure(err)
	return ret
}

func (self *StorageWriter) GetPendingL1CheckPointInfo() *schema.L1CheckPointInfo {
	v, err := self.overlay.Get(schema.PendingL1CheckPointInfoKey)
	utils.Ensure(err)
	if len(v) == 0 {
		return nil
	}
	ret := new(schema.L1CheckPointInfo)
	err = ret.Deserialization(codec.NewZeroCopySource(v))
	utils.Ensure(err)
	return ret
}

func (self *StorageWriter) GetLastSyncedL1Hash() web3.Hash {
	v, err := self.overlay.Get(schema.LastSyncedL1Hash)
	utils.Ensure(err)
	if len(v) == 0 {
		return web3.Hash{}
	}
	hash, err := codec.NewZeroCopySource(v).ReadHash()
	utils.Ensure(err)

	return hash
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

func (self *StorageWriter) SetLastSyncedL2Height(height uint64) {
	self.overlay.Put(schema.LastSyncedL2HeightKey, codec.NewZeroCopySink(nil).WriteUint64(height).Bytes())
}

func (self *StorageWriter) GetLastSyncedL2Height() uint64 {
	v, err := self.overlay.Get(schema.LastSyncedL2HeightKey)
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

// DirtyKeys suppose there is only write operation, just return all write key
func (self *StorageWriter) Dirty() ([][]byte, [][]byte) {
	b := make([][]byte, 0, 1024)
	v := make([][]byte, 0, 1024)
	self.overlay.(*overlaydb.OverlayDB).GetWriteSet().ForEach(
		func(key, val []byte) {
			if len(val) == 0 {
				//suppose only support write operation just panic
				panic(1)
			}
			_copy := make([]byte, len(key))
			copy(_copy, key)
			b = append(b, _copy)

			//copy old value
			old, _ := self.overlay.(*overlaydb.OverlayDB).Store.Get(key)
			_copy = make([]byte, len(old))
			copy(_copy, old)
			v = append(v, _copy)
		})
	return b, v
}

func (self *StorageWriter) Cover(key []byte, value []byte) {
	if len(value) == 0 {
		self.overlay.Delete(key)
	} else {
		self.overlay.Put(key, value)
	}
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
