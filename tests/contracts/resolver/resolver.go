package resolver

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/utils/common/math"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/build"
)

type AddressManager web3.Address

//function newAddr(string memory _name, address _addr) public onlyOwner noEmptyAddr(_addr)
func (self AddressManager) NewAddr(sender evm.AccountRef, vm *evm.EVM, name string, addr web3.Address) error {
	input := build.AddressManagerAbi().Methods["newAddr"].MustEncodeIDAndInput(name, addr)
	_, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	return err
}

//function updateAddr(string memory _name, address _addr) public onlyOwner noEmptyAddr(_addr)
func (self AddressManager) UpdateAddr(sender evm.AccountRef, vm *evm.EVM, name string, addr web3.Address) error {
	input := build.AddressManagerAbi().Methods["updateAddr"].MustEncodeIDAndInput(name, addr)
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return fmt.Errorf(web3.DecodeRevert(ret))
	}
	return err
}

//function getAddr(string memory _name) public view returns (address)
func (self AddressManager) GetAddr(sender evm.AccountRef, vm *evm.EVM, name string) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["getAddr"]
	input := method.MustEncodeIDAndInput(name)
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function resolve(string memory _name) public view returns (address)
func (self AddressManager) Resolve(sender evm.AccountRef, vm *evm.EVM, name string) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["resolve"]
	input := method.MustEncodeIDAndInput(name)
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function dao() public view returns (address)
func (self AddressManager) Dao(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["dao"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

// function rollupInputChain() public view returns (IRollupInputChain)
func (self AddressManager) RollupInputChain(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["rollupInputChain"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function rollupInputChainContainer() public view returns (IChainStorageContainer)
func (self AddressManager) RollupInputChainContainer(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["rollupInputChainContainer"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function rollupStateChain() public view returns (IRollupStateChain)
func (self AddressManager) RollupStateChain(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["rollupStateChain"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function rollupStateChainContainer() public view returns (IChainStorageContainer)
func (self AddressManager) RollupStateChainContainer(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["rollupStateChainContainer"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function stakingManager() public view returns (IStakingManager)
func (self AddressManager) StakingManager(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["stakingManager"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function challengeFactory() public view returns (IChallengeFactory)
func (self AddressManager) ChallengeFactory(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["challengeFactory"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function l1CrossLayerWitness() public view returns (IL1CrossLayerWitness)
func (self AddressManager) L1CrossLayerWitness(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["l1CrossLayerWitness"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}

//function l2CrossLayerWitness() public view returns (IL2CrossLayerWitness)
func (self AddressManager) L2CrossLayerWitness(sender evm.AccountRef, vm *evm.EVM) (out web3.Address, err error) {
	method := build.AddressManagerAbi().Methods["l2CrossLayerWitness"]
	input := method.ID()
	ret, _, err := vm.Call(sender, web3.Address(self), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return web3.Address{}, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return web3.Address{}, err
	}
	if err := contracts.Decode(i, &out); err != nil {
		return web3.Address{}, err
	}
	return
}
