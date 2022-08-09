package relayer

import (
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type RelayerStore struct {
	store schema.KeyValueDB
}

func NewStore(db schema.KeyValueDB) *RelayerStore {
	return &RelayerStore{
		store: db,
	}
}

func NewMemStore() *RelayerStore {
	return &RelayerStore{
		store: overlaydb.NewOverlayDB(storage.NewFakeDB()),
	}
}

func (self *RelayerStore) StorePendingL1MsgIndex(pendingIndex uint64) {
	self.store.Put(schema.L1RelayerPendingMsgIndex, codec.NewZeroCopySink(nil).WriteUint64(pendingIndex).Bytes())
}

func (self *RelayerStore) GetPendingL1MsgIndex() uint64 {
	v, err := self.store.Get(schema.L1RelayerPendingMsgIndex)
	utils.Ensure(err)
	if len(v) == 0 {
		return 0
	}
	data, err := codec.NewZeroCopySource(v).ReadUint64()
	utils.Ensure(err)
	return data
}
