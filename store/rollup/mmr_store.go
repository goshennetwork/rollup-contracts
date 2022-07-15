package rollup

import (
	"encoding/binary"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/merkle"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type MMR struct {
	dataPrefix byte
	treeKey    []byte
	store      schema.KeyValueDB
	tree       *merkle.CompactMerkleTree
}

func (self *MMR) genHashKey(index uint64) []byte {
	index += 1 //first index is used for hash num
	var b [9]byte
	b[0] = self.dataPrefix
	binary.BigEndian.PutUint64(b[1:], index)
	return b[:]
}

func (self *MMR) genHashSizeKey() []byte {
	var b [9]byte
	b[0] = self.dataPrefix
	binary.BigEndian.PutUint64(b[1:], 0)
	return b[:]
}

func (self *MMR) getPendingIndex() uint64 {
	size, err := self.store.Get(self.genHashSizeKey())
	utils.Ensure(err)
	if len(size) == 0 {
		return 0
	}
	r := codec.NewZeroCopyReader(size)
	i := r.ReadUint64BE()
	utils.Ensure(r.Error())
	return i
}

func (self *MMR) storePendingIndex(pendingIndex uint64) {
	self.store.Put(self.genHashSizeKey(), codec.NewZeroCopySink(nil).WriteUint64BE(pendingIndex).Bytes())
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

func (self *MMR) StoreCompactMerkleTree(tree *merkle.CompactMerkleTree) {
	self.store.Put(self.treeKey, schema.SerializeCompactMerkleTree(tree))
}

func (self *MMR) GetCompactMerkleTree() *merkle.CompactMerkleTree {
	v, err := self.store.Get(self.treeKey)
	utils.Ensure(err)
	size, hashes := uint64(0), []web3.Hash{}
	if len(v) != 0 {
		size, hashes, err = schema.DeserializeCompactMerkleTree(v)
	}
	utils.Ensure(err)
	return merkle.NewTree(size, hashes, self)
}

func (self *MMR) Append(hash []web3.Hash) error {
	pendingIndex := self.getPendingIndex()
	for _, s := range hash {
		self.store.Put(self.genHashKey(pendingIndex), codec.NewZeroCopySink(nil).WriteHash(s).Bytes())
		pendingIndex++
	}
	self.storePendingIndex(pendingIndex)
	return nil
}

func (self *MMR) GetHash(pos uint64) (web3.Hash, error) {
	v, err := self.store.Get(self.genHashKey(pos))
	utils.Ensure(err)
	if len(v) == 0 {
		return web3.Hash{}, schema.ErrNotFound
	}
	r := codec.NewZeroCopyReader(v)
	return r.ReadHash(), r.Error()
}

// HashStore is an interface for persist hash
type HashStore interface {
	Append(hash []web3.Hash) error
	GetHash(pos uint64) (web3.Hash, error)
}
