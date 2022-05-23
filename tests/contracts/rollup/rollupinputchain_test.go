package rollup

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils/common"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/resolver"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/staking"
	"gotest.tools/assert"
)

func TestEnqueue(t *testing.T) {
	c := contracts.NewCase()
	addrManager := resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	container := ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "RollupInputChain", web3.Address(addrManager)))
	rollupInputChain := RollupInputChain(contracts.NewRollupInputChain(c.Sender, c.Vm, web3.Address(addrManager), new(big.Int).SetUint64(2_000_000), new(big.Int).SetUint64(1_000_000)))
	_ = addrManager.NewAddr(c.Sender, c.Vm, "RollupInputChainContainer", web3.Address(container))
	_ = addrManager.NewAddr(c.Sender, c.Vm, "RollupInputChain", web3.Address(rollupInputChain))
	_ = addrManager.NewAddr(c.Sender, c.Vm, "L1CrossLayerWitness", web3.Address{53, 53, 53, 53})
	target := web3.Address{1, 1}
	err := rollupInputChain.Enqueue(c.Sender, c.Vm, target, 900_000, []byte("test"))
	assert.NilError(t, err)
	_, _, err = rollupInputChain.GetQueueTxInfo(c.Sender, c.Vm, 0)
	assert.NilError(t, err)

	c.Vm.Origin = web3.Address{2, 2, 2, 2}
	_ = addrManager.UpdateAddr(c.Sender, c.Vm, "L1CrossLayerWitness", c.Sender.Address())
	err = rollupInputChain.Enqueue(c.Sender, c.Vm, target, 900_000, []byte("test"))
	assert.NilError(t, err)
}

func TestAppendBatches(t *testing.T) {
	c := contracts.NewCase()
	addrManager := resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	container := ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "RollupInputChain", web3.Address(addrManager)))
	rollupInputChain := RollupInputChain(contracts.NewRollupInputChain(c.Sender, c.Vm, web3.Address(addrManager), new(big.Int).SetUint64(2_000_000), new(big.Int).SetUint64(1_000_000)))
	erc20 := contracts.NewERC20(c.Sender, c.Vm, "test", "v0")
	stakingManager := staking.StakingManager(contracts.NewStakingManager(c.Sender, c.Vm, c.Sender.Address(), web3.Address{}, web3.Address(rollupInputChain), erc20, new(big.Int)))
	err := stakingManager.Deposit(c.Sender, c.Vm)
	assert.NilError(t, err)
	_ = addrManager.NewAddr(c.Sender, c.Vm, "RollupInputChainContainer", web3.Address(container))
	_ = addrManager.NewAddr(c.Sender, c.Vm, "RollupInputChain", web3.Address(rollupInputChain))
	_ = addrManager.NewAddr(c.Sender, c.Vm, "L1CrossLayerWitness", web3.Address{53, 53, 53, 53})
	_ = addrManager.NewAddr(c.Sender, c.Vm, "StakingManager", web3.Address(stakingManager))

	batches := &RollupInputBatches{
		QueueNum:    0,
		QueueStart:  0,
		BatchNum:    1,
		Batch0Time:  uint64(time.Now().Unix()),
		BatchesData: common.Hash{}.Bytes(),
	}
	c.Vm.Context.Time = new(big.Int).SetInt64(time.Now().Add(1 * time.Hour).Unix())
	err = rollupInputChain.AppendBatch(c.Sender, c.Vm, batches)
	assert.NilError(t, err)
	height, err := rollupInputChain.ChainHeight(c.Sender, c.Vm)
	assert.NilError(t, err)
	assert.Equal(t, height, uint64(1))
	timestamp, err := rollupInputChain.LastTimestamp(c.Sender, c.Vm)
	assert.NilError(t, err)
	assert.Equal(t, timestamp, batches.Batch0Time)
}

func TestRandomBatches(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randSlice := make([]byte, r.Intn(100)+32)
	r.Read(randSlice)

	c := contracts.NewCase()
	addrManager := resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	container := ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "RollupInputChain", web3.Address(addrManager)))
	rollupInputChain := RollupInputChain(contracts.NewRollupInputChain(c.Sender, c.Vm, web3.Address(addrManager), new(big.Int).SetUint64(2_000_000), new(big.Int).SetUint64(1_000_000)))
	erc20 := contracts.NewERC20(c.Sender, c.Vm, "test", "v0")
	stakingManager := staking.StakingManager(contracts.NewStakingManager(c.Sender, c.Vm, c.Sender.Address(), web3.Address{}, web3.Address(rollupInputChain), erc20, new(big.Int)))
	err := stakingManager.Deposit(c.Sender, c.Vm)
	assert.NilError(t, err)
	_ = addrManager.NewAddr(c.Sender, c.Vm, "RollupInputChainContainer", web3.Address(container))
	_ = addrManager.NewAddr(c.Sender, c.Vm, "RollupInputChain", web3.Address(rollupInputChain))
	_ = addrManager.NewAddr(c.Sender, c.Vm, "L1CrossLayerWitness", web3.Address{53, 53, 53, 53})
	_ = addrManager.NewAddr(c.Sender, c.Vm, "StakingManager", web3.Address(stakingManager))
	totalQueue := uint64(0)
	timestamp := uint64(time.Now().Unix())
	var lastTimestamp uint64
	for _, v := range randSlice {
		c.Vm.Context.Time.SetUint64(timestamp)
		err = rollupInputChain.Enqueue(c.Sender, c.Vm, c.Sender.Address(), 900_000, nil)
		assert.NilError(t, err)
		totalQueue++
		if v&1 == 0 {
			batchNum := r.Uint64()%100 + 1
			timeDiff := make([]uint32, batchNum)
			pendginQueue, err := rollupInputChain.PendingQueueIndex(c.Sender, c.Vm)
			assert.NilError(t, err)
			batches := &RollupInputBatches{
				QueueNum:          totalQueue - pendginQueue,
				QueueStart:        pendginQueue,
				BatchNum:          batchNum,
				Batch0Time:        timestamp + 1,
				BatchLeftTimeDiff: timeDiff,
				BatchesData:       common.CopyBytes(randSlice),
			}
			c.Vm.Context.Time.SetUint64(timestamp + 2)
			err = rollupInputChain.AppendBatch(c.Sender, c.Vm, batches)
			assert.NilError(t, err)
			timestamp += 10
		}

		gotTimestamp, err := rollupInputChain.LastTimestamp(c.Sender, c.Vm)
		assert.NilError(t, err)
		if gotTimestamp < lastTimestamp {
			t.Fatal("wrong  timestamp")
		}
	}
}
