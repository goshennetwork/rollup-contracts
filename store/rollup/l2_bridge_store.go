package rollup

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type L2BridgeStore struct {
	store schema.KeyValueDB
}

func NewL2BridgeStore(db schema.KeyValueDB) *L2BridgeStore {
	return &L2BridgeStore{
		store: db,
	}
}

// private method for test
func newL2BridgeMemStore() *L2BridgeStore {
	return &L2BridgeStore{
		store: overlaydb.NewOverlayDB(storage.NewFakeDB()),
	}
}

func (self *L2BridgeStore) StoreWithdrawal(events []*binding.WithdrawalInitiatedEvent) {
	cached := make(map[web3.Hash]binding.CrossLayerInfos, 0)
	for _, evt := range events {
		data, ok := cached[evt.Raw.TransactionHash]
		if !ok {
			data = make([]*binding.CrossLayerInfo, 0)
		}
		data = append(data, evt.GetTokenCrossInfo())
		cached[evt.Raw.TransactionHash] = data
	}
	for txHash, evts := range cached {
		self.store.Put(genL2WithdrawalInitKey(txHash), codec.SerializeToBytes(evts))
	}
}

func (self *L2BridgeStore) StoreDepositFinalized(events []*binding.DepositFinalizedEvent) {
	cached := make(map[web3.Hash]binding.CrossLayerInfos, 0)
	for _, evt := range events {
		data, ok := cached[evt.Raw.TransactionHash]
		if !ok {
			data = make([]*binding.CrossLayerInfo, 0)
		}
		data = append(data, evt.GetTokenCrossInfo())
		cached[evt.Raw.TransactionHash] = data
	}
	for txHash, evts := range cached {
		self.store.Put(genDepositFinalizedKey(txHash), codec.SerializeToBytes(evts))
	}
}

func (self *L2BridgeStore) StoreDepositFailed(events []*binding.DepositFailedEvent) {
	cached := make(map[web3.Hash]binding.CrossLayerInfos, 0)
	for _, evt := range events {
		data, ok := cached[evt.Raw.TransactionHash]
		if !ok {
			data = make([]*binding.CrossLayerInfo, 0)
		}
		data = append(data, evt.GetTokenCrossInfo())
		cached[evt.Raw.TransactionHash] = data
	}
	for txHash, evts := range cached {
		self.store.Put(genDepositFailedKey(txHash), codec.SerializeToBytes(evts))
	}
}
