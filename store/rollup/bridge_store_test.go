package rollup

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

func TestL1Store(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	l1BridgeStore := newL1BridgeMemStore()

	l1EthDepositEvents := genRandomL1EthDepositEvts(10)
	l1EthWithdrawalEvents := genRandomL1EthWithdrawalEvts(10)
	l1Erc20DepositEvents := genRandomL1Erc20DepositEvts(10)
	l1Erc20WithdrawalEvents := genRandomL1Erc20Withdrawals(10)

	l1BridgeStore.StoreETHDeposit(l1EthDepositEvents)
	l1BridgeStore.StoreETHWithdrawal(l1EthWithdrawalEvents)
	l1BridgeStore.StoreERC20Deposit(l1Erc20DepositEvents)
	l1BridgeStore.StoreERC20Withdrawal(l1Erc20WithdrawalEvents)

	assertL1EthDepositEvtsEqual(t, l1BridgeStore.store, l1EthDepositEvents)
	assertL1EthWithdrawalEvtsEqual(t, l1BridgeStore.store, l1EthWithdrawalEvents)
	assertERC20DepositInitiatedEventEqual(t, l1BridgeStore.store, l1Erc20DepositEvents)
	assertERC20WithdrawalFinalizedEventEqual(t, l1BridgeStore.store, l1Erc20WithdrawalEvents)
}

func TestL2Store(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	l2BridgeStore := newL2BridgeMemStore()

	l2Withdrawals := genRandomL2Withdrawals(10)
	l2Deposits := genRandomL2Deposits(10)
	l2DepositFailed := genRandomL2DepositFailed(10)

	l2BridgeStore.StoreWithdrawal(l2Withdrawals)
	l2BridgeStore.StoreDepositFinalized(l2Deposits)
	l2BridgeStore.StoreDepositFailed(l2DepositFailed)

	assertWithdrawalInitiatedEventEqual(t, l2BridgeStore.store, l2Withdrawals)
	assertDepositFinalizedEventEqual(t, l2BridgeStore.store, l2Deposits)
	assertDepositFailedEventEqual(t, l2BridgeStore.store, l2DepositFailed)
}

func genRandomL1EthDepositEvts(length int) []*binding.ETHDepositInitiatedEvent {
	result := make([]*binding.ETHDepositInitiatedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.ETHDepositInitiatedEvent{
			Data: make([]byte, 20), Raw: &web3.Log{},
		}
		_, _ = rand.Read(evt.From[:])
		_, _ = rand.Read(evt.To[:])
		_, _ = rand.Read(evt.Data[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		amountData := make([]byte, 8)
		_, _ = rand.Read(amountData[:])
		evt.Amount = new(big.Int).SetBytes(amountData)
		result = append(result, evt)
	}
	return result
}

func genRandomL1EthWithdrawalEvts(length int) []*binding.ETHWithdrawalFinalizedEvent {
	result := make([]*binding.ETHWithdrawalFinalizedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.ETHWithdrawalFinalizedEvent{
			Data: make([]byte, 20), Raw: &web3.Log{},
		}
		_, _ = rand.Read(evt.From[:])
		_, _ = rand.Read(evt.To[:])
		_, _ = rand.Read(evt.Data[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		amountData := make([]byte, 8)
		_, _ = rand.Read(amountData[:])
		evt.Amount = new(big.Int).SetBytes(amountData)
		result = append(result, evt)
	}
	return result
}

func genRandomL1Erc20DepositEvts(length int) []*binding.ERC20DepositInitiatedEvent {
	result := make([]*binding.ERC20DepositInitiatedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.ERC20DepositInitiatedEvent{
			Data: make([]byte, 20), Raw: &web3.Log{},
		}
		_, _ = rand.Read(evt.From[:])
		_, _ = rand.Read(evt.To[:])
		_, _ = rand.Read(evt.Data[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		amountData := make([]byte, 8)
		_, _ = rand.Read(amountData[:])
		evt.Amount = new(big.Int).SetBytes(amountData)
		result = append(result, evt)
	}
	return result
}

func genRandomL1Erc20Withdrawals(length int) []*binding.ERC20WithdrawalFinalizedEvent {
	result := make([]*binding.ERC20WithdrawalFinalizedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.ERC20WithdrawalFinalizedEvent{
			Data: make([]byte, 20), Raw: &web3.Log{},
		}
		_, _ = rand.Read(evt.From[:])
		_, _ = rand.Read(evt.To[:])
		_, _ = rand.Read(evt.Data[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		amountData := make([]byte, 8)
		_, _ = rand.Read(amountData[:])
		evt.Amount = new(big.Int).SetBytes(amountData)
		result = append(result, evt)
	}
	return result
}

func genRandomL2Withdrawals(length int) []*binding.WithdrawalInitiatedEvent {
	result := make([]*binding.WithdrawalInitiatedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.WithdrawalInitiatedEvent{
			Data: make([]byte, 20), Raw: &web3.Log{},
		}
		_, _ = rand.Read(evt.From[:])
		_, _ = rand.Read(evt.To[:])
		_, _ = rand.Read(evt.Data[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		amountData := make([]byte, 8)
		_, _ = rand.Read(amountData[:])
		evt.Amount = new(big.Int).SetBytes(amountData)
		result = append(result, evt)
	}
	return result
}

func genRandomL2Deposits(length int) []*binding.DepositFinalizedEvent {
	result := make([]*binding.DepositFinalizedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.DepositFinalizedEvent{
			Data: make([]byte, 20), Raw: &web3.Log{},
		}
		_, _ = rand.Read(evt.From[:])
		_, _ = rand.Read(evt.To[:])
		_, _ = rand.Read(evt.Data[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		amountData := make([]byte, 8)
		_, _ = rand.Read(amountData[:])
		evt.Amount = new(big.Int).SetBytes(amountData)
		result = append(result, evt)
	}
	return result
}

func genRandomL2DepositFailed(length int) []*binding.DepositFailedEvent {
	result := make([]*binding.DepositFailedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.DepositFailedEvent{
			Data: make([]byte, 20), Raw: &web3.Log{},
		}
		_, _ = rand.Read(evt.From[:])
		_, _ = rand.Read(evt.To[:])
		_, _ = rand.Read(evt.Data[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		amountData := make([]byte, 8)
		_, _ = rand.Read(amountData[:])
		evt.Amount = new(big.Int).SetBytes(amountData)
		result = append(result, evt)
	}
	return result
}

func assertL1EthDepositEvtsEqual(t *testing.T, store schema.KeyValueDB, events []*binding.ETHDepositInitiatedEvent) {
	for _, evt := range events {
		newEvt := readL1TokenBridgeETHEvent(t, store, evt.Raw.TransactionHash, false)
		exist := false
		for _, item := range newEvt {
			if l1TokenBridgeETHDepositEqual(evt, item) {
				exist = true
				break
			}
		}
		if !exist {
			t.Fatal("failed")
		}
	}
}

func assertL1EthWithdrawalEvtsEqual(t *testing.T, store schema.KeyValueDB, events []*binding.ETHWithdrawalFinalizedEvent) {
	for _, evt := range events {
		newEvt := readL1TokenBridgeETHEvent(t, store, evt.Raw.TransactionHash, true)
		exist := false
		for _, item := range newEvt {
			if l1TokenBridgeETHWithdrawalEqual(evt, item) {
				exist = true
				break
			}
		}
		if !exist {
			t.Fatal("failed")
		}
	}
}

func assertDepositFailedEventEqual(t *testing.T, store schema.KeyValueDB, events []*binding.DepositFailedEvent) {
	for _, evt := range events {
		newEvt := readCrossLayerTokenInfo(t, store, evt.Raw.TransactionHash, genDepositFailedKey)
		exist := false
		for _, item := range newEvt {
			if crossLayerTokenInfoEqual(evt.GetTokenCrossInfo(), item) {
				exist = true
				break
			}
		}
		if !exist {
			t.Fatal("failed")
		}
	}
}

func assertDepositFinalizedEventEqual(t *testing.T, store schema.KeyValueDB, events []*binding.DepositFinalizedEvent) {
	for _, evt := range events {
		newEvt := readCrossLayerTokenInfo(t, store, evt.Raw.TransactionHash, genDepositFinalizedKey)
		exist := false
		for _, item := range newEvt {
			if crossLayerTokenInfoEqual(evt.GetTokenCrossInfo(), item) {
				exist = true
				break
			}
		}
		if !exist {
			t.Fatal("failed")
		}
	}
}

func assertERC20DepositInitiatedEventEqual(t *testing.T, store schema.KeyValueDB, events []*binding.ERC20DepositInitiatedEvent) {
	for _, evt := range events {
		newEvt := readCrossLayerTokenInfo(t, store, evt.Raw.TransactionHash, genL1ERC20DepositInitKey)
		exist := false
		for _, item := range newEvt {
			if crossLayerTokenInfoEqual(evt.GetTokenCrossInfo(), item) {
				exist = true
				break
			}
		}
		if !exist {
			t.Fatal("failed")
		}
	}
}

func assertERC20WithdrawalFinalizedEventEqual(t *testing.T, store schema.KeyValueDB, events []*binding.ERC20WithdrawalFinalizedEvent) {
	for _, evt := range events {
		newEvt := readCrossLayerTokenInfo(t, store, evt.Raw.TransactionHash, genL1ERC20WithdrawalFinalizedKey)
		exist := false
		for _, item := range newEvt {
			if crossLayerTokenInfoEqual(evt.GetTokenCrossInfo(), item) {
				exist = true
				break
			}
		}
		if !exist {
			t.Fatal("failed")
		}
	}
}

func assertWithdrawalInitiatedEventEqual(t *testing.T, store schema.KeyValueDB, events []*binding.WithdrawalInitiatedEvent) {
	for _, evt := range events {
		newEvt := readCrossLayerTokenInfo(t, store, evt.Raw.TransactionHash, genL2WithdrawalInitKey)
		exist := false
		for _, item := range newEvt {
			if crossLayerTokenInfoEqual(evt.GetTokenCrossInfo(), item) {
				exist = true
				break
			}
		}
		if !exist {
			t.Fatal("failed")
		}
	}
}

func readL1TokenBridgeETHEvent(t *testing.T, store schema.KeyValueDB, txHash web3.Hash, isWithdrawal bool) schema.L1TokenBridgeETHEvents {
	key := genL1ETHDepositKey(txHash)
	if isWithdrawal {
		key = genL1ETHWithdrawalKey(txHash)
	}
	data, err := store.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	source := codec.NewZeroCopySource(data)
	newEvt, err := schema.DeserializationL1TokenBridgeETHEvents(source)
	if err != nil {
		t.Fatal(err)
	}
	return newEvt
}

func readCrossLayerTokenInfo(t *testing.T, store schema.KeyValueDB, txHash web3.Hash,
	keyGen func(hash web3.Hash) []byte) binding.CrossLayerInfos {
	key := keyGen(txHash)
	data, err := store.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	source := codec.NewZeroCopySource(data)
	newEvt, err := binding.DeserializationCrossLayerInfos(source)
	if err != nil {
		t.Fatal(err)
	}
	return newEvt
}

func crossLayerTokenInfoEqual(item, newItem *binding.CrossLayerInfo) bool {
	return item.L1Token == newItem.L1Token || item.L2Token == newItem.L2Token || item.From == newItem.From ||
		item.To == newItem.To || item.Amount.Uint64() == newItem.Amount.Uint64() ||
		!bytes.Equal(item.Data, newItem.Data)
}

func l1TokenBridgeETHDepositEqual(info *binding.ETHDepositInitiatedEvent, newInfo *schema.L1TokenBridgeETHEvent) bool {
	return info.To == newInfo.To || info.From == newInfo.From || info.Amount.Uint64() == newInfo.Amount.Uint64() ||
		!bytes.Equal(info.Data, newInfo.Data)
}
func l1TokenBridgeETHWithdrawalEqual(info *binding.ETHWithdrawalFinalizedEvent, newInfo *schema.L1TokenBridgeETHEvent) bool {
	return info.To == newInfo.To || info.From == newInfo.From || info.Amount.Uint64() == newInfo.Amount.Uint64() ||
		!bytes.Equal(info.Data, newInfo.Data)
}
