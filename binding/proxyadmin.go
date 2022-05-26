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

// ProxyAdmin is a solidity contract
type ProxyAdmin struct {
	c *contract.Contract
}

// DeployProxyAdmin deploys a new ProxyAdmin contract
func DeployProxyAdmin(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiProxyAdmin, binProxyAdmin)
}

// NewProxyAdmin creates a new instance of the contract at a specific address
func NewProxyAdmin(addr web3.Address, provider *jsonrpc.Client) *ProxyAdmin {
	return &ProxyAdmin{c: contract.NewContract(addr, abiProxyAdmin, provider)}
}

// Contract returns the contract object
func (_a *ProxyAdmin) Contract() *contract.Contract {
	return _a.c
}

// calls

// GetProxyAdmin calls the getProxyAdmin method in the solidity contract
func (_a *ProxyAdmin) GetProxyAdmin(proxy web3.Address, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("getProxyAdmin", web3.EncodeBlock(block...), proxy)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// GetProxyImplementation calls the getProxyImplementation method in the solidity contract
func (_a *ProxyAdmin) GetProxyImplementation(proxy web3.Address, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("getProxyImplementation", web3.EncodeBlock(block...), proxy)
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
func (_a *ProxyAdmin) Owner(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// ChangeProxyAdmin sends a changeProxyAdmin transaction in the solidity contract
func (_a *ProxyAdmin) ChangeProxyAdmin(proxy web3.Address, newAdmin web3.Address) *contract.Txn {
	return _a.c.Txn("changeProxyAdmin", proxy, newAdmin)
}

// RenounceOwnership sends a renounceOwnership transaction in the solidity contract
func (_a *ProxyAdmin) RenounceOwnership() *contract.Txn {
	return _a.c.Txn("renounceOwnership")
}

// TransferOwnership sends a transferOwnership transaction in the solidity contract
func (_a *ProxyAdmin) TransferOwnership(newOwner web3.Address) *contract.Txn {
	return _a.c.Txn("transferOwnership", newOwner)
}

// Upgrade sends a upgrade transaction in the solidity contract
func (_a *ProxyAdmin) Upgrade(proxy web3.Address, implementation web3.Address) *contract.Txn {
	return _a.c.Txn("upgrade", proxy, implementation)
}

// UpgradeAndCall sends a upgradeAndCall transaction in the solidity contract
func (_a *ProxyAdmin) UpgradeAndCall(proxy web3.Address, implementation web3.Address, data []byte) *contract.Txn {
	return _a.c.Txn("upgradeAndCall", proxy, implementation, data)
}

// events

func (_a *ProxyAdmin) OwnershipTransferredTopicFilter(previousOwner []web3.Address, newOwner []web3.Address) [][]web3.Hash {

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

func (_a *ProxyAdmin) FilterOwnershipTransferredEvent(previousOwner []web3.Address, newOwner []web3.Address, startBlock uint64, endBlock ...uint64) ([]*OwnershipTransferredEvent, error) {
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
