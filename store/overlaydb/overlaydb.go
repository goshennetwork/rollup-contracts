/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package overlaydb

import (
	"github.com/goshennetwork/rollup-contracts/store/leveldbstore"
	common "github.com/goshennetwork/rollup-contracts/store/schema"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type OverlayDB struct {
	store common.PersistStore
	memdb *MemDB
	dbErr error
}

const initCap = 4 * 1024
const initkvNum = 128

func NewOverlayDB(store common.PersistStore) *OverlayDB {
	return &OverlayDB{
		store: store,
		memdb: NewMemDB(initCap, initkvNum),
	}
}

func (self *OverlayDB) Reset() {
	self.memdb.Reset()
}

func (self *OverlayDB) Error() error {
	return self.dbErr
}

func (self *OverlayDB) SetError(err error) {
	self.dbErr = err
}

// Get Warning: if key is not found, value&err is both nil!
func (self *OverlayDB) Get(key []byte) (value []byte, err error) {
	var unknown bool
	value, unknown = self.memdb.Get(key)
	if !unknown {
		return value, nil
	}

	value, err = self.store.Get(key)
	if err != nil {
		if err == common.ErrNotFound {
			return nil, nil
		}
		self.dbErr = err
		return nil, err
	}

	return
}

func (self *OverlayDB) Put(key []byte, value []byte) {
	self.memdb.Put(key, value)
}

func (self *OverlayDB) Delete(key []byte) {
	self.memdb.Delete(key)
}

//CommitTo write memdb in this OverlayDB to levelDB
func (self *OverlayDB) CommitTo() {
	batch := leveldbstore.NewBatch()
	self.memdb.ForEach(func(key, val []byte) {
		if len(val) == 0 {
			batch.Delete(key)
		} else {
			batch.Put(key, val)
		}
	})
	if err := self.store.BatchCommit(batch); err != nil {
		panic(err)
	}
}

func (self *OverlayDB) GetWriteSet() *MemDB {
	return self.memdb
}

// param key is referenced by iterator
func (self *OverlayDB) NewIterator(key []byte) common.StoreIterator {
	prefixRange := util.BytesPrefix(key)
	backIter := self.store.NewIterator(key)
	memIter := self.memdb.NewIterator(prefixRange)

	return NewJoinIter(memIter, backIter)
}
