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

// TransparentUpgradeableProxy is a solidity contract
type TransparentUpgradeableProxy struct {
	c *contract.Contract
}

// DeployTransparentUpgradeableProxy deploys a new TransparentUpgradeableProxy contract
func DeployTransparentUpgradeableProxy(provider *jsonrpc.Client, from web3.Address, logic web3.Address, admin web3.Address, data []byte) *contract.Txn {
	return contract.DeployContract(provider, from, abiTransparentUpgradeableProxy, binTransparentUpgradeableProxy, logic, admin, data)
}

// NewTransparentUpgradeableProxy creates a new instance of the contract at a specific address
func NewTransparentUpgradeableProxy(addr web3.Address, provider *jsonrpc.Client) *TransparentUpgradeableProxy {
	return &TransparentUpgradeableProxy{c: contract.NewContract(addr, abiTransparentUpgradeableProxy, provider)}
}

// Contract returns the contract object
func (_a *TransparentUpgradeableProxy) Contract() *contract.Contract {
	return _a.c
}

// calls

// txns

// Admin sends a admin transaction in the solidity contract
func (_a *TransparentUpgradeableProxy) Admin() *contract.Txn {
	return _a.c.Txn("admin")
}

// ChangeAdmin sends a changeAdmin transaction in the solidity contract
func (_a *TransparentUpgradeableProxy) ChangeAdmin(newAdmin web3.Address) *contract.Txn {
	return _a.c.Txn("changeAdmin", newAdmin)
}

// Implementation sends a implementation transaction in the solidity contract
func (_a *TransparentUpgradeableProxy) Implementation() *contract.Txn {
	return _a.c.Txn("implementation")
}

// UpgradeTo sends a upgradeTo transaction in the solidity contract
func (_a *TransparentUpgradeableProxy) UpgradeTo(newImplementation web3.Address) *contract.Txn {
	return _a.c.Txn("upgradeTo", newImplementation)
}

// UpgradeToAndCall sends a upgradeToAndCall transaction in the solidity contract
func (_a *TransparentUpgradeableProxy) UpgradeToAndCall(newImplementation web3.Address, data []byte) *contract.Txn {
	return _a.c.Txn("upgradeToAndCall", newImplementation, data)
}

// events

func (_a *TransparentUpgradeableProxy) AdminChangedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{AdminChangedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *TransparentUpgradeableProxy) FilterAdminChangedEvent(startBlock uint64, endBlock ...uint64) ([]*AdminChangedEvent, error) {
	topic := _a.AdminChangedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*AdminChangedEvent, 0)
	evts := _a.c.Abi.Events["AdminChanged"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem AdminChangedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *TransparentUpgradeableProxy) BeaconUpgradedTopicFilter(beacon []web3.Address) [][]web3.Hash {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{BeaconUpgradedEventID}, beaconRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *TransparentUpgradeableProxy) FilterBeaconUpgradedEvent(beacon []web3.Address, startBlock uint64, endBlock ...uint64) ([]*BeaconUpgradedEvent, error) {
	topic := _a.BeaconUpgradedTopicFilter(beacon)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*BeaconUpgradedEvent, 0)
	evts := _a.c.Abi.Events["BeaconUpgraded"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem BeaconUpgradedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *TransparentUpgradeableProxy) UpgradedTopicFilter(implementation []web3.Address) [][]web3.Hash {

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

func (_a *TransparentUpgradeableProxy) FilterUpgradedEvent(implementation []web3.Address, startBlock uint64, endBlock ...uint64) ([]*UpgradedEvent, error) {
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
