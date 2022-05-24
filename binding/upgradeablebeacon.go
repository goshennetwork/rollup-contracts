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

// UpgradeableBeacon is a solidity contract
type UpgradeableBeacon struct {
	c *contract.Contract
}

// DeployUpgradeableBeacon deploys a new UpgradeableBeacon contract
func DeployUpgradeableBeacon(provider *jsonrpc.Client, from web3.Address, implementation web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiUpgradeableBeacon, binUpgradeableBeacon, implementation)
}

// NewUpgradeableBeacon creates a new instance of the contract at a specific address
func NewUpgradeableBeacon(addr web3.Address, provider *jsonrpc.Client) *UpgradeableBeacon {
	return &UpgradeableBeacon{c: contract.NewContract(addr, abiUpgradeableBeacon, provider)}
}

// Contract returns the contract object
func (_a *UpgradeableBeacon) Contract() *contract.Contract {
	return _a.c
}

// calls

// Implementation calls the implementation method in the solidity contract
func (_a *UpgradeableBeacon) Implementation(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("implementation", web3.EncodeBlock(block...))
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
func (_a *UpgradeableBeacon) Owner(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// RenounceOwnership sends a renounceOwnership transaction in the solidity contract
func (_a *UpgradeableBeacon) RenounceOwnership() *contract.Txn {
	return _a.c.Txn("renounceOwnership")
}

// TransferOwnership sends a transferOwnership transaction in the solidity contract
func (_a *UpgradeableBeacon) TransferOwnership(newOwner web3.Address) *contract.Txn {
	return _a.c.Txn("transferOwnership", newOwner)
}

// UpgradeTo sends a upgradeTo transaction in the solidity contract
func (_a *UpgradeableBeacon) UpgradeTo(newImplementation web3.Address) *contract.Txn {
	return _a.c.Txn("upgradeTo", newImplementation)
}

// events

func (_a *UpgradeableBeacon) OwnershipTransferredTopicFilter(previousOwner []web3.Address, newOwner []web3.Address) [][]web3.Hash {

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

func (_a *UpgradeableBeacon) FilterOwnershipTransferredEvent(previousOwner []web3.Address, newOwner []web3.Address, startBlock uint64, endBlock ...uint64) ([]*OwnershipTransferredEvent, error) {
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

func (_a *UpgradeableBeacon) UpgradedTopicFilter(implementation []web3.Address) [][]web3.Hash {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{UpgradedEventID}, implementationRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *UpgradeableBeacon) FilterUpgradedEvent(implementation []web3.Address, startBlock uint64, endBlock ...uint64) ([]*UpgradedEvent, error) {
	topic := _a.UpgradedTopicFilter(implementation)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*UpgradedEvent, 0)
	evts := _a.c.Abi.Events["Upgraded"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem UpgradedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
