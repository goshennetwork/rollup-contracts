package rollup

import (
	"math/rand"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/utils/common"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/resolver"
	"gotest.tools/assert"
)

func expectErr(err error) {
	if err == nil {
		panic("expect err,but got nil")
	}
}

func TestChainSize(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	var chainStorageContainer = ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "deployer", web3.Address(addressManager)))
	size, err := chainStorageContainer.ChainSize(c.Sender, c.Vm)
	assert.NilError(t, err)
	assert.Equal(t, size, uint64(0))
}

func TestAppend(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	var chainStorageContainer = ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "deployer", web3.Address(addressManager)))
	element := common.BytesToHash([]byte("element"))
	//no deployer addr
	err := chainStorageContainer.Append(c.Sender, c.Vm, element)
	expectErr(err)

	err = addressManager.NewAddr(c.Sender, c.Vm, "deployer", web3.Address(c.Sender))
	assert.NilError(t, err)
	err = chainStorageContainer.Append(c.Sender, c.Vm, element)
	assert.NilError(t, err)
	size, err := chainStorageContainer.ChainSize(c.Sender, c.Vm)
	assert.NilError(t, err)
	assert.Equal(t, size, uint64(1))

	//not deployer
	err = chainStorageContainer.Append(evm.AccountRef{1, 1, 1, 1}, c.Vm, element)
	expectErr(err)
	gotElement, err := chainStorageContainer.Get(c.Sender, c.Vm, 0)
	assert.NilError(t, err)
	assert.Equal(t, element, gotElement)
}

//function resize(uint64 _newSize) public onlyOwner
func TestResize(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = resolver.AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	var chainStorageContainer = ChainStorageContainer(contracts.NewChainStorageContainer(c.Sender, c.Vm, "deployer", web3.Address(addressManager)))
	err := addressManager.NewAddr(c.Sender, c.Vm, "deployer", web3.Address(c.Sender))
	assert.NilError(t, err)
	element := common.BytesToHash([]byte("element"))
	err = chainStorageContainer.Append(c.Sender, c.Vm, element)
	assert.NilError(t, err)
	err = chainStorageContainer.Resize(c.Sender, c.Vm, 0)
	assert.NilError(t, err)
	_, err = chainStorageContainer.Get(c.Sender, c.Vm, 0)
	expectErr(err)
	err = chainStorageContainer.Append(c.Sender, c.Vm, element)
	assert.NilError(t, err)
	gotElement, err := chainStorageContainer.Get(c.Sender, c.Vm, 0)
	assert.NilError(t, err)
	assert.Equal(t, element, gotElement)
	size, err := chainStorageContainer.ChainSize(c.Sender, c.Vm)
	assert.NilError(t, err)
	assert.Equal(t, size, uint64(1))

	//not deployer
	err = chainStorageContainer.Resize(evm.AccountRef{1, 1, 1, 1}, c.Vm, 0)
	expectErr(err)
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
