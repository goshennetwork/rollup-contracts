package leveldbstore

import (
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

//LevelDB store
type LevelDBStore struct {
	db *leveldb.DB // LevelDB instance
}

// used to compute the size of bloom filter bits array .
// too small will lead to high false positive rate.
const BITSPERKEY = 10

//NewLevelDBStore return LevelDBStore instance
func NewLevelDBStore(file string) (*LevelDBStore, error) {
	openFileCache := opt.DefaultOpenFilesCacheCapacity
	maxOpenFiles, err := fdlimit.Current()
	if err == nil && maxOpenFiles < openFileCache*5 {
		openFileCache = maxOpenFiles / 5
	}

	if openFileCache < 16 {
		openFileCache = 16
	}

	// default Options
	o := opt.Options{
		NoSync:                 false,
		OpenFilesCacheCapacity: openFileCache,
		Filter:                 filter.NewBloomFilter(BITSPERKEY),
	}

	db, err := leveldb.OpenFile(file, &o)

	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}

	if err != nil {
		return nil, err
	}

	return &LevelDBStore{
		db: db,
	}, nil
}

func NewMemLevelDBStore() *LevelDBStore {
	store := storage.NewMemStorage()
	// default Options
	o := opt.Options{
		NoSync: false,
		Filter: filter.NewBloomFilter(BITSPERKEY),
	}
	db, err := leveldb.Open(store, &o)
	if err != nil {
		panic(err)
	}

	return &LevelDBStore{
		db: db,
	}
}

//Put a key-value pair to leveldb
func (self *LevelDBStore) Put(key []byte, value []byte) error {
	return self.db.Put(key, value, nil)
}

//Get the value of a key from leveldb
func (self *LevelDBStore) Get(key []byte) ([]byte, error) {
	dat, err := self.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, schema.ErrNotFound
		}
		return nil, err
	}
	return dat, nil
}

//Has return whether the key is exist in leveldb
func (self *LevelDBStore) Has(key []byte) (bool, error) {
	return self.db.Has(key, nil)
}

//Delete the the in leveldb
func (self *LevelDBStore) Delete(key []byte) error {
	return self.db.Delete(key, nil)
}

//NewBatch start commit batch
func NewBatch() *leveldb.Batch {
	return new(leveldb.Batch)
}

//BatchCommit commit batch to leveldb
func (self *LevelDBStore) BatchCommit(batch *leveldb.Batch) error {
	err := self.db.Write(batch, nil)
	if err != nil {
		return err
	}
	return nil
}

//Close leveldb
func (self *LevelDBStore) Close() error {
	err := self.db.Close()
	return err
}

//NewIterator return a iterator of leveldb with the key prefix
func (self *LevelDBStore) NewIterator(prefix []byte) schema.StoreIterator {

	iter := self.db.NewIterator(util.BytesPrefix(prefix), nil)

	return iter
}
