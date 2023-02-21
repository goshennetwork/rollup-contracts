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

// AddressManager is a solidity contract
type AddressManager struct {
	c *contract.Contract
}

// DeployAddressManager deploys a new AddressManager contract
func DeployAddressManager(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiAddressManager, binAddressManager)
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

// ChallengeBeacon calls the challengeBeacon method in the solidity contract
func (_a *AddressManager) ChallengeBeacon(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("challengeBeacon", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

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

// FeeToken calls the feeToken method in the solidity contract
func (_a *AddressManager) FeeToken(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("feeToken", web3.EncodeBlock(block...))
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

// GetAddrByHash calls the getAddrByHash method in the solidity contract
func (_a *AddressManager) GetAddrByHash(val0 [32]byte, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("getAddrByHash", web3.EncodeBlock(block...), val0)
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

// L1StandardBridge calls the l1StandardBridge method in the solidity contract
func (_a *AddressManager) L1StandardBridge(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("l1StandardBridge", web3.EncodeBlock(block...))
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

// MachineState calls the machineState method in the solidity contract
func (_a *AddressManager) MachineState(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("machineState", web3.EncodeBlock(block...))
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

// StateTransition calls the stateTransition method in the solidity contract
func (_a *AddressManager) StateTransition(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("stateTransition", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Whitelist calls the whitelist method in the solidity contract
func (_a *AddressManager) Whitelist(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("whitelist", web3.EncodeBlock(block...))
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

// Initialize sends a initialize transaction in the solidity contract
func (_a *AddressManager) Initialize() *contract.Txn {
	return _a.c.Txn("initialize")
}

// RenounceOwnership sends a renounceOwnership transaction in the solidity contract
func (_a *AddressManager) RenounceOwnership() *contract.Txn {
	return _a.c.Txn("renounceOwnership")
}

// SetAddress sends a setAddress transaction in the solidity contract
func (_a *AddressManager) SetAddress(name string, addr web3.Address) *contract.Txn {
	return _a.c.Txn("setAddress", name, addr)
}

// SetAddressBatch sends a setAddressBatch transaction in the solidity contract
func (_a *AddressManager) SetAddressBatch(names []string, addrs []web3.Address) *contract.Txn {
	return _a.c.Txn("setAddressBatch", names, addrs)
}

// TransferOwnership sends a transferOwnership transaction in the solidity contract
func (_a *AddressManager) TransferOwnership(newOwner web3.Address) *contract.Txn {
	return _a.c.Txn("transferOwnership", newOwner)
}

// events

func (_a *AddressManager) AddressSetTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{AddressSetEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *AddressManager) FilterAddressSetEvent(startBlock uint64, endBlock ...uint64) ([]*AddressSetEvent, error) {
	topic := _a.AddressSetTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*AddressSetEvent, 0)
	evts := _a.c.Abi.Events["AddressSet"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem AddressSetEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *AddressManager) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *AddressManager) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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
