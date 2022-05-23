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

// AddressManager is a solidity contract
type AddressManager struct {
	c *contract.Contract
}

// DeployAddressManager deploys a new AddressManager contract
func DeployAddressManager(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiAddressManager, binAddressManager, args...)
}

// NewAddressManager creates a new instance of the contract at a specific address
func NewAddressManager(addr web3.Address, provider *jsonrpc.Client) *AddressManager {
	return &AddressManager{c: contract.NewContract(addr, abiAddressManager, provider)}
}

// Contract returns the contract object
func (_a *AddressManager) Contract() *contract.Contract {
	return _a.c
}

// calls

// ChallengeFactory calls the challengeFactory method in the solidity contract
func (_a *AddressManager) ChallengeFactory(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("challengeFactory", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Dao calls the dao method in the solidity contract
func (_a *AddressManager) Dao(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("dao", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// GetAddr calls the getAddr method in the solidity contract
func (_a *AddressManager) GetAddr(name string, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("getAddr", web3.EncodeBlock(block...), name)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// L1CrossLayerWitness calls the l1CrossLayerWitness method in the solidity contract
func (_a *AddressManager) L1CrossLayerWitness(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("l1CrossLayerWitness", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// L2CrossLayerWitness calls the l2CrossLayerWitness method in the solidity contract
func (_a *AddressManager) L2CrossLayerWitness(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("l2CrossLayerWitness", web3.EncodeBlock(block...))
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
func (_a *AddressManager) Owner(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// Resolve calls the resolve method in the solidity contract
func (_a *AddressManager) Resolve(name string, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("resolve", web3.EncodeBlock(block...), name)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// RollupInputChain calls the rollupInputChain method in the solidity contract
func (_a *AddressManager) RollupInputChain(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("rollupInputChain", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// RollupInputChainContainer calls the rollupInputChainContainer method in the solidity contract
func (_a *AddressManager) RollupInputChainContainer(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("rollupInputChainContainer", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// RollupStateChain calls the rollupStateChain method in the solidity contract
func (_a *AddressManager) RollupStateChain(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("rollupStateChain", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// RollupStateChainContainer calls the rollupStateChainContainer method in the solidity contract
func (_a *AddressManager) RollupStateChainContainer(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("rollupStateChainContainer", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// StakingManager calls the stakingManager method in the solidity contract
func (_a *AddressManager) StakingManager(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("stakingManager", web3.EncodeBlock(block...))
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

// NewAddr sends a newAddr transaction in the solidity contract
func (_a *AddressManager) NewAddr(name string, addr web3.Address) *contract.Txn {
	return _a.c.Txn("newAddr", name, addr)
}

// RenounceOwnership sends a renounceOwnership transaction in the solidity contract
func (_a *AddressManager) RenounceOwnership() *contract.Txn {
	return _a.c.Txn("renounceOwnership")
}

// TransferOwnership sends a transferOwnership transaction in the solidity contract
func (_a *AddressManager) TransferOwnership(newOwner web3.Address) *contract.Txn {
	return _a.c.Txn("transferOwnership", newOwner)
}

// UpdateAddr sends a updateAddr transaction in the solidity contract
func (_a *AddressManager) UpdateAddr(name string, addr web3.Address) *contract.Txn {
	return _a.c.Txn("updateAddr", name, addr)
}

// events

var AddressUpdatedEventID = crypto.Keccak256Hash([]byte("AddressUpdated(string,address,address)"))

func (_a *AddressManager) AddressUpdatedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{AddressUpdatedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *AddressManager) FilterAddressUpdatedEvent(startBlock uint64, endBlock ...uint64) ([]*AddressUpdatedEvent, error) {
	topic := _a.AddressUpdatedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*AddressUpdatedEvent, 0)
	evts := _a.c.Abi.Events["AddressUpdated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem AddressUpdatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var OwnershipTransferredEventID = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))

func (_a *AddressManager) OwnershipTransferredTopicFilter(previousOwner []web3.Address, newOwner []web3.Address) [][]web3.Hash {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{OwnershipTransferredEventID}, previousOwnerRule, newOwnerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *AddressManager) FilterOwnershipTransferredEvent(previousOwner []web3.Address, newOwner []web3.Address, startBlock uint64, endBlock ...uint64) ([]*OwnershipTransferredEvent, error) {
	topic := _a.OwnershipTransferredTopicFilter(previousOwner, newOwner)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*OwnershipTransferredEvent, 0)
	evts := _a.c.Abi.Events["OwnershipTransferred"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem OwnershipTransferredEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
