package binding

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
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
)

// ChainStorageContainer is a solidity contract
type ChainStorageContainer struct {
	c *contract.Contract
}

// DeployChainStorageContainer deploys a new ChainStorageContainer contract
func DeployChainStorageContainer(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiChainStorageContainer, binChainStorageContainer)
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

// Owner calls the owner method in the solidity contract
func (_a *ChainStorageContainer) Owner(block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("owner", web3.EncodeBlock(block...))
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

// Initialize sends a initialize transaction in the solidity contract
func (_a *ChainStorageContainer) Initialize(owner string, addressResolver web3.Address) *contract.Txn {
	return _a.c.Txn("initialize", owner, addressResolver)
}

// Resize sends a resize transaction in the solidity contract
func (_a *ChainStorageContainer) Resize(newSize uint64) *contract.Txn {
	return _a.c.Txn("resize", newSize)
}

// events

func (_a *ChainStorageContainer) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *ChainStorageContainer) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
	topic := _a.InitializedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*InitializedEvent, 0)
	evts := _a.c.Abi.Events["Initialized"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem InitializedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
