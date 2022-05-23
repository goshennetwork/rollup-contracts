package contracts

import (
	"math/big"
	"strconv"
	"time"

	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/evm/params"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/executor"
	"github.com/laizy/web3/utils/common/math"
	"github.com/mitchellh/mapstructure"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/build"
)

//Decode args must be pointer
func Decode(outputMap interface{}, args ...interface{}) error {
	for i := 0; i < len(args); i++ {
		if err := mapstructure.Decode(outputMap.(map[string]interface{})[strconv.Itoa(i)], args[i]); err != nil {
			return err
		}
	}
	return nil
}

func NewEVM() *evm.EVM {
	var hashFn evm.GetHashFunc = func(u uint64) web3.Hash {
		var h web3.Hash
		h.SetBytes(crypto.Keccak256(new(big.Int).SetUint64(u).Bytes()))
		return h
	}
	caccheDB := storage.NewCacheDB(overlaydb.NewOverlayDB(storage.NewFakeDB()))
	statedb := storage.NewStateDB(caccheDB)
	ctx := executor.NewEVMBlockContext(0, uint64(time.Now().Unix()), hashFn)
	vmenv := evm.NewEVM(ctx, evm.TxContext{}, statedb, params.MainnetChainConfig, evm.Config{})
	return vmenv
}

//constructor()
func NewAddressManager(sender evm.AccountRef, vm *evm.EVM) web3.Address {
	_, addr, _, err := vm.Create(sender, build.AddressManagerBin(), math.MaxUint64, new(big.Int))
	if err != nil {
		panic(err)
	}
	return addr
}

//constructor(string memory _owner, address _addressResolver)
func NewChainStorageContainer(sender evm.AccountRef, vm *evm.EVM, owner string, resolver web3.Address) web3.Address {
	param, err := abi.Encode([]interface{}{owner, resolver}, build.ChainStorageContainerAbi().Constructor.Inputs)
	if err != nil {
		panic(err)
	}
	inputs := append(build.ChainStorageContainerBin(), param...)
	_, addr, _, err := vm.Create(sender, inputs, math.MaxUint64, new(big.Int))
	if err != nil {
		panic(err)
	}
	return addr
}

//  constructor(
//        address _addressResolver,
//        uint256 _maxTxGasLimit,
//        uint256 _maxCrossLayerTxGasLimit
//    )
func NewRollupInputChain(sender evm.AccountRef, vm *evm.EVM, addressResolver web3.Address, maxTxGasLimit *big.Int, maxCrossLayerTxGasLimit *big.Int) web3.Address {

	param, err := abi.Encode([]interface{}{addressResolver, maxTxGasLimit, maxCrossLayerTxGasLimit}, build.RollupInputChainAbi().Constructor.Inputs)
	if err != nil {
		panic(err)
	}
	inputs := append(build.RollupInputChainBin(), param...)
	_, addr, _, err := vm.Create(sender, inputs, math.MaxUint64, new(big.Int))
	if err != nil {
		panic(err)
	}
	return addr
}

//    constructor(
//        address _DAOAddress,
//        address _challengeFactory,
//        address _rollupStateChain,
//        address _erc20,
//        uint256 _price
//    )
func NewStakingManager(sender evm.AccountRef, vm *evm.EVM, daoAddress, challengeFactory, rollupStateChain, erc20 web3.Address, price *big.Int) web3.Address {
	param, err := abi.Encode([]interface{}{daoAddress, challengeFactory, rollupStateChain, erc20, price}, build.StakingManagerAbi().Constructor.Inputs)
	if err != nil {
		panic(err)
	}
	inputs := append(build.StakingManagerBin(), param...)
	_, addr, _, err := vm.Create(sender, inputs, math.MaxUint64, new(big.Int))
	if err != nil {
		panic(err)
	}
	return addr
}

// constructor(string memory name_, string memory symbol_) {
//        _name = name_;
//        _symbol = symbol_;
//    }
func NewERC20(sender evm.AccountRef, vm *evm.EVM, name, symbol string) web3.Address {
	param, err := abi.Encode([]interface{}{name, symbol}, build.ERC20Abi().Constructor.Inputs)
	if err != nil {
		panic(err)
	}
	input := append(build.ERC20Bin(), param...)
	_, addr, _, err := vm.Create(sender, input, math.MaxUint64, new(big.Int))
	if err != nil {
		panic(err)
	}
	return addr
}

type TestCase struct {
	Vm     *evm.EVM
	Sender evm.AccountRef
}

func NewCase() *TestCase {
	vm := NewEVM()
	c := &TestCase{
		Vm:     vm,
		Sender: evm.AccountRef{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
	}
	c.Vm.Origin = web3.Address(c.Sender)
	return c
}
