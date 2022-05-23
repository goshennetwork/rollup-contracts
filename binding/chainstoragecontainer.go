package binding

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/mitchellh/mapstructure"
)

var (
	_ = json.Unmarshal
	_ = big.NewInt
	_ = fmt.Printf
	_ = utils.JsonStr
	_ = mapstructure.Decode
	_ = crypto.Keccak256Hash
)

// ChainStorageContainer is a solidity contract
type ChainStorageContainer struct {
	c *contract.Contract
}

// DeployChainStorageContainer deploys a new ChainStorageContainer contract
func DeployChainStorageContainer(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiChainStorageContainer, binChainStorageContainer, args...)
}

// NewChainStorageContainer creates a new instance of the contract at a specific address
func NewChainStorageContainer(addr web3.Address, provider *jsonrpc.Client) *ChainStorageContainer {
	return &ChainStorageContainer{c: contract.NewContract(addr, abiChainStorageContainer, provider)}
}

// Contract returns the contract object
func (_a *ChainStorageContainer) Contract() *contract.Contract {
	return _a.c
}

// calls

// ChainSize calls the chainSize method in the solidity contract
func (_a *ChainStorageContainer) ChainSize(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("chainSize", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Get calls the get method in the solidity contract
func (_a *ChainStorageContainer) Get(index uint64, block ...web3.BlockNumber) (retval0 [32]byte, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("get", web3.EncodeBlock(block...), index)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// LastTimestamp calls the lastTimestamp method in the solidity contract
func (_a *ChainStorageContainer) LastTimestamp(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("lastTimestamp", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// txns

// Append sends a append transaction in the solidity contract
func (_a *ChainStorageContainer) Append(element [32]byte) *contract.Txn {
	return _a.c.Txn("append", element)
}

// Resize sends a resize transaction in the solidity contract
func (_a *ChainStorageContainer) Resize(newSize uint64) *contract.Txn {
	return _a.c.Txn("resize", newSize)
}

// SetLastTimestamp sends a setLastTimestamp transaction in the solidity contract
func (_a *ChainStorageContainer) SetLastTimestamp(timestamp uint64) *contract.Txn {
	return _a.c.Txn("setLastTimestamp", timestamp)
}

// events
