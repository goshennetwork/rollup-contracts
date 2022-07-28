package rollup

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/stretchr/testify/assert"
)

func TestL1Store(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	l1BridgeStore := newL1BridgeMemStore()
	testL1Deposit(t, 100, l1BridgeStore.StoreDeposit, l1BridgeStore.GetDeposit)
	testL1Withdraw(t, 100, l1BridgeStore.StoreWithdrawal, l1BridgeStore.GetWithdrawal)
}

func TestL2Store(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	l2BridgeStore := newL2BridgeMemStore()
	testL2DepositFinalized(t, 100, l2BridgeStore.StoreDepositFinalized, l2BridgeStore.GetDepositFinalized)
	testL2Withdraw(t, 100, l2BridgeStore.StoreWithdrawal, l2BridgeStore.GetWithdrawal)
	testL2DepositFailed(t, 100, l2BridgeStore.StoreDepositFailed, l2BridgeStore.GetDepositFailed)
}

func testL1Deposit(t *testing.T, len int, store func(events []*binding.DepositInitiatedEvent), get func(hash web3.Hash) (binding.CrossLayerInfos, error)) {
	evts := genRandomL1DepositEvts(len)
	store(evts)

	for i := 0; i < len; i++ {
		txHash := evts[i].Raw.TransactionHash
		info, err := get(txHash)
		assert.Nil(t, err)
		ori := make(binding.CrossLayerInfos, 0)
		for _, e := range evts {
			if e.Raw.TransactionHash == txHash {
				ori = append(ori, e.GetTokenCrossInfo())
			}
		}
		assert.True(t, reflect.DeepEqual(ori, info))
	}
}

func testL1Withdraw(t *testing.T, len int, store func(events []*binding.WithdrawalFinalizedEvent), get func(hash web3.Hash) (binding.CrossLayerInfos, error)) {
	evts := genRandomL1WithdrawalEvts(len)
	store(evts)

	for i := 0; i < len; i++ {
		txHash := evts[i].Raw.TransactionHash
		info, err := get(txHash)
		assert.Nil(t, err)
		ori := make(binding.CrossLayerInfos, 0)
		for _, e := range evts {
			if e.Raw.TransactionHash == txHash {
				ori = append(ori, e.GetTokenCrossInfo())
			}
		}
		assert.True(t, reflect.DeepEqual(ori, info))

	}
}

func testL2DepositFinalized(t *testing.T, len int, store func(events []*binding.DepositFinalizedEvent), get func(hash web3.Hash) (binding.CrossLayerInfos, error)) {
	evts := genRandomL2Deposits(len)
	store(evts)

	for i := 0; i < len; i++ {
		txHash := evts[i].Raw.TransactionHash
		info, err := get(txHash)
		assert.Nil(t, err)
		ori := make(binding.CrossLayerInfos, 0)
		for _, e := range evts {
			if e.Raw.TransactionHash == txHash {
				ori = append(ori, e.GetTokenCrossInfo())
			}
		}
		assert.True(t, reflect.DeepEqual(ori, info))
	}
}

func testL2Withdraw(t *testing.T, len int, store func(events []*binding.WithdrawalInitiatedEvent), get func(hash web3.Hash) (binding.CrossLayerInfos, error)) {
	evts := genRandomL2Withdrawals(len)
	store(evts)

	for i := 0; i < len; i++ {
		txHash := evts[i].Raw.TransactionHash
		info, err := get(txHash)
		assert.Nil(t, err)
		ori := make(binding.CrossLayerInfos, 0)
		for _, e := range evts {
			if e.Raw.TransactionHash == txHash {
				ori = append(ori, e.GetTokenCrossInfo())
			}
		}
		assert.True(t, reflect.DeepEqual(ori, info))

	}
}

func testL2DepositFailed(t *testing.T, len int, store func(events []*binding.DepositFailedEvent), get func(hash web3.Hash) (binding.CrossLayerInfos, error)) {
	evts := genRandomL2DepositFailed(len)
	store(evts)

	for i := 0; i < len; i++ {
		txHash := evts[i].Raw.TransactionHash
		info, err := get(txHash)
		assert.Nil(t, err)
		ori := make(binding.CrossLayerInfos, 0)
		for _, e := range evts {
			if e.Raw.TransactionHash == txHash {
				ori = append(ori, e.GetTokenCrossInfo())
			}
		}
		assert.True(t, reflect.DeepEqual(ori, info))

	}
}

func genRandomL1DepositEvts(length int) []*binding.DepositInitiatedEvent {
	result := make([]*binding.DepositInitiatedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.DepositInitiatedEvent{
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

func genRandomL1WithdrawalEvts(length int) []*binding.WithdrawalFinalizedEvent {
	result := make([]*binding.WithdrawalFinalizedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.WithdrawalFinalizedEvent{
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
		_, _ = rand.Read(evt.L1Token[:])
		_, _ = rand.Read(evt.L2Token[:])
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
