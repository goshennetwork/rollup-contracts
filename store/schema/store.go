package schema

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

var ErrNotFound = errors.New("not found")

type StoreIterator interface {
	Next() bool //Next item. If item available return true, otherwise return false
	//Prev() bool           //previous item. If item available return true, otherwise return false
	First() bool //First item. If item available return true, otherwise return false
	//Last() bool           //Last item. If item available return true, otherwise return false
	//Seek(key []byte) bool //Seek key. If item available return true, otherwise return false
	Key() []byte   //Return the current item key
	Value() []byte //Return the current item value
	Release()      //Close iterator
	Error() error  // Error returns any accumulated error.
}

type PersistStore interface {
	Put(key []byte, value []byte) error      //Put the key-value pair to store
	Get(key []byte) ([]byte, error)          //Get the value if key in store
	Has(key []byte) (bool, error)            //Whether the key is exist in store
	Delete(key []byte) error                 //Delete the key in store
	BatchCommit(batch *leveldb.Batch) error  //Commit batch to store
	Close() error                            //Close store
	NewIterator(prefix []byte) StoreIterator //Return the iterator of store
}

type KeyValueDB interface {
	KeyValueReader
	KeyValueWriter
}

type KeyValueReader interface {
	Get(key []byte) (value []byte, err error) // returns nil, nil if not found
}

type KeyValueWriter interface {
	Put(key []byte, value []byte)
	Delete(key []byte)
}
