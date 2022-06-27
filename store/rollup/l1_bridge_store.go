package rollup

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
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

func (self *L1BridgeStore) StoreETHDeposit(events []*binding.ETHDepositInitiatedEvent) {
	cached := make(map[web3.Hash]schema.L1TokenBridgeETHEvents, 0)
	for _, evt := range events {
		data, ok := cached[evt.Raw.TransactionHash]
		if !ok {
			data = make([]*schema.L1TokenBridgeETHEvent, 0)
		}
		data = append(data, &schema.L1TokenBridgeETHEvent{
			From:   evt.From,
			To:     evt.To,
			Amount: evt.Amount,
			Data:   evt.Data,
		})
		cached[evt.Raw.TransactionHash] = data
	}
	for txHash, evts := range cached {
		self.store.Put(genL1ETHDepositKey(txHash), codec.SerializeToBytes(evts))
	}
}

func (self *L1BridgeStore) StoreERC20Deposit(events []*binding.ERC20DepositInitiatedEvent) {
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
		self.store.Put(genL1ERC20DepositInitKey(txHash), codec.SerializeToBytes(evts))
	}
}

func (self *L1BridgeStore) StoreETHWithdrawal(events []*binding.ETHWithdrawalFinalizedEvent) {
	cached := make(map[web3.Hash]schema.L1TokenBridgeETHEvents, 0)
	for _, evt := range events {
		data, ok := cached[evt.Raw.TransactionHash]
		if !ok {
			data = make([]*schema.L1TokenBridgeETHEvent, 0)
		}
		data = append(data, &schema.L1TokenBridgeETHEvent{
			From:   evt.From,
			To:     evt.To,
			Amount: evt.Amount,
			Data:   evt.Data,
		})
		cached[evt.Raw.TransactionHash] = data
	}
	for txHash, evts := range cached {
		self.store.Put(genL1ETHWithdrawalKey(txHash), codec.SerializeToBytes(evts))
	}
}

func (self *L1BridgeStore) StoreERC20Withdrawal(events []*binding.ERC20WithdrawalFinalizedEvent) {
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
		self.store.Put(genL1ERC20WithdrawalFinalizedKey(txHash), codec.SerializeToBytes(evts))
	}
}
