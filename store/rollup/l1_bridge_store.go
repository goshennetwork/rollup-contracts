package rollup

import (
	"github.com/goshennetwork/rollup-contracts/binding"
	"github.com/goshennetwork/rollup-contracts/store/schema"
	"github.com/laizy/web3"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
)

type L1BridgeStore struct {
	store schema.KeyValueDB
}

func NewL1BridgeStore(db schema.KeyValueDB) *L1BridgeStore {
	return &L1BridgeStore{
		store: db,
	}
}

// private method for test
func newL1BridgeMemStore() *L1BridgeStore {
	return &L1BridgeStore{
		store: overlaydb.NewOverlayDB(storage.NewFakeDB()),
	}
}

func (self *L1BridgeStore) StoreDeposit(events []*binding.DepositInitiatedEvent) {
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
		self.store.Put(genL1DepositKey(txHash), codec.SerializeToBytes(evts))
	}
}

func (self *L1BridgeStore) GetDeposit(txHash web3.Hash) (binding.CrossLayerInfos, error) {
	v, err := self.store.Get(genL1DepositKey(txHash))
	utils.Ensure(err)
	if len(v) == 0 {
		return nil, schema.ErrNotFound
	}
	source := codec.NewZeroCopySource(v)
	return binding.DeserializationCrossLayerInfos(source)
}

func (self *L1BridgeStore) StoreWithdrawal(events []*binding.WithdrawalFinalizedEvent) {
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
		self.store.Put(genL1WithdrawalKey(txHash), codec.SerializeToBytes(evts))
	}
}

func (self *L1BridgeStore) GetWithdrawal(txHash web3.Hash) (binding.CrossLayerInfos, error) {
	v, err := self.store.Get(genL1WithdrawalKey(txHash))
	utils.Ensure(err)
	if len(v) == 0 {
		return nil, schema.ErrNotFound
	}
	source := codec.NewZeroCopySource(v)
	return binding.DeserializationCrossLayerInfos(source)
}
