package resolver

import (
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"gotest.tools/assert"
)

func expectErr(err error) {
	if err == nil {
		panic("expect err,but got nil")
	}
}

func TestNewAddr(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	test := web3.Address{}
	test.SetBytes([]byte("test"))
	err := addressManager.NewAddr(c.Sender, c.Vm, "test", test)
	assert.NilError(t, err)
	gotTest, err := addressManager.GetAddr(c.Sender, c.Vm, "test")
	assert.NilError(t, err)
	assert.Equal(t, test, gotTest)

	//no owner
	err = addressManager.NewAddr(evm.AccountRef{1, 1, 1, 1}, c.Vm, "test", test)
	expectErr(err)
}

func TestUpdateAddr(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	test := web3.Address{}
	test.SetBytes([]byte("test"))
	//update empty
	err := addressManager.UpdateAddr(c.Sender, c.Vm, "test", test)
	expectErr(err)

	err = addressManager.NewAddr(c.Sender, c.Vm, "test", test)
	assert.NilError(t, err)
	test.SetBytes([]byte("test1"))
	err = addressManager.UpdateAddr(c.Sender, c.Vm, "test", test)
	//no owner
	err = addressManager.UpdateAddr(evm.AccountRef{1, 1, 1, 1}, c.Vm, "test", test)
	expectErr(err)

	gotTest, err := addressManager.GetAddr(c.Sender, c.Vm, "test")
	assert.NilError(t, err)
	assert.Equal(t, test, gotTest)
}

func TestResolve(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	//empty addr
	_, err := addressManager.Resolve(c.Sender, c.Vm, "test")
	expectErr(err)
}

func TestGet(t *testing.T) {
	c := contracts.NewCase()
	var addressManager = AddressManager(contracts.NewAddressManager(c.Sender, c.Vm))
	err := addressManager.NewAddr(c.Sender, c.Vm, "L1CrossLayerWitness", web3.Address{53, 53, 53, 53})
	assert.NilError(t, err)
	_, err = addressManager.L1CrossLayerWitness(c.Sender, c.Vm)
	assert.NilError(t, err)
}
