package rollup

import (
	"math/big"

	"github.com/laizy/web3/utils/common"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/utils/common/math"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/build"
)

type ChainStorageContainer web3.Address

//function chainSize() external view returns (uint64)
func (self ChainStorageContainer) ChainSize(sender evm.AccountRef, vm *evm.EVM) (out uint64, err error) {
	method := build.ChainStorageContainerAbi().Methods["chainSize"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return 0, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return 0, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return 0, err
	}
	return
}

//function append(bytes32 _element) public onlyOwner
func (self ChainStorageContainer) Append(sender evm.AccountRef, vm *evm.EVM, element common.Hash) error {
	input := build.ChainStorageContainerAbi().Methods["append"].MustEncodeIDAndInput(element)
	_, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return err
	}
	return nil
}

//function resize(uint64 _newSize) public onlyOwner
func (self ChainStorageContainer) Resize(sender evm.AccountRef, vm *evm.EVM, newSize uint64) error {
	input := build.ChainStorageContainerAbi().Methods["resize"].MustEncodeIDAndInput(newSize)
	_, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return err
	}
	return nil
}

//function setLastTimestamp(uint64 _timestamp) public onlyOwner
func (self ChainStorageContainer) SetLastTimestamp(sender evm.AccountRef, vm *evm.EVM, timestamp uint64) error {
	input := build.ChainStorageContainerAbi().Methods["setLastTimestamp"].MustEncodeIDAndInput(timestamp)
	_, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return err
	}
	return nil
}

//function lastTimestamp() external view returns (uint64)
func (self ChainStorageContainer) LastTimestamp(sender evm.AccountRef, vm *evm.EVM) (out uint64, err error) {
	method := build.ChainStorageContainerAbi().Methods["lastTimestamp"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return 0, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return 0, err
	}
	err = contracts.Decode(i, &out)
	return
}

//function get(uint64 _index) public view returns (bytes32)
func (self ChainStorageContainer) Get(sender evm.AccountRef, vm *evm.EVM, index uint64) (out common.Hash, err error) {
	method := build.ChainStorageContainerAbi().Methods["get"]
	input := method.MustEncodeIDAndInput(index)
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return common.Hash{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return common.Hash{}, err
	}
	err = contracts.Decode(i, &out)
	return
}
