package rollup

import (
	"math/rand"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/common"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"gotest.tools/assert"
)

func expectErr(err error) {
	if err == nil {
		panic("expect err,but got nil")
	}
}

func TestChainSize(t *testing.T) {
	chainEnv := contracts.LocalChainEnv
	signer := contracts.SetupLocalSigner(chainEnv)
	l1Chain := contracts.DeployL1Contract(signer, chainEnv.L1ChainConfig)

	size, err := l1Chain.InputChainStorage.ChainSize()
	utils.Ensure(err)
	assert.Equal(t, size, uint64(0))
}

func TestAppend(t *testing.T) {
	chainEnv := contracts.LocalChainEnv
	signer := contracts.SetupLocalSigner(chainEnv)
	l1Chain := contracts.DeployL1Contract(signer, chainEnv.L1ChainConfig)

	// not owner
	element := common.BytesToHash([]byte("element"))
	receipt := l1Chain.InputChainStorage.Append(element).SetGasLimit(5000000).SetGasPrice(2000).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 0)

	// change storage owner
	l1Chain.AddressManager.SetAddress("RollupInputChain", signer.Address()).Sign(signer).SendTransaction(signer)

	receipt = l1Chain.InputChainStorage.Append(element).SetGasLimit(5000000).SetGasPrice(2000).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
}

func TestSetTimestamp(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	var chainStorageContainer = ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "deployer", web3.Address(addressManager)))
	err := addressManager.NewAddr(c.Sender, c.Vm, "deployer", web3.Address(c.Sender))
	assert.NilError(t, err)
	gotTimestamp, err := chainStorageContainer.LastTimestamp(c.Sender, c.Vm)
	assert.NilError(t, err)
	assert.Equal(t, gotTimestamp, uint64(0))
	timestamp := uint64(time.Now().Unix())
	err = chainStorageContainer.SetLastTimestamp(c.Sender, c.Vm, timestamp)
	assert.NilError(t, err)
	gotTimestamp, err = chainStorageContainer.LastTimestamp(c.Sender, c.Vm)
	assert.NilError(t, err)
	assert.Equal(t, timestamp, gotTimestamp)

	//not deployer
	err = chainStorageContainer.SetLastTimestamp(evm.AccountRef{1, 1, 1, 1}, c.Vm, timestamp)
	expectErr(err)
}

func TestRandom(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randSlice := make([]byte, r.Intn(10000))
	r.Read(randSlice)
	c := contracts.NewCase()
	var addressManager = resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	var chainStorageContainer = ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "deployer", web3.Address(addressManager)))
	err := addressManager.NewAddr(c.Sender, c.Vm, "deployer", web3.Address(c.Sender))
	assert.NilError(t, err)
	calcSize := uint64(0)
	//append
	for _, v := range randSlice {
		err = chainStorageContainer.Append(c.Sender, c.Vm, common.Hash{v})
		assert.NilError(t, err)
		calcSize++
		if v&1 == 0 { //resize
			err = chainStorageContainer.Resize(c.Sender, c.Vm, calcSize-1)
			assert.NilError(t, err)
			calcSize--
		}
		size, err := chainStorageContainer.ChainSize(c.Sender, c.Vm)
		assert.NilError(t, err)
		assert.Equal(t, calcSize, size)
	}

	//get
	index := uint64(0)
	for _, v := range randSlice {
		if v&1 == 1 { //not resized
			out, err := chainStorageContainer.Get(c.Sender, c.Vm, index)
			assert.NilError(t, err)
			index++
			assert.Equal(t, out, common.Hash{v})
		}
	}

}
