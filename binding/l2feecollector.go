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

// L2FeeCollector is a solidity contract
type L2FeeCollector struct {
	c *contract.Contract
}

// DeployL2FeeCollector deploys a new L2FeeCollector contract
func DeployL2FeeCollector(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiL2FeeCollector, binL2FeeCollector)
}

// NewL2FeeCollector creates a new instance of the contract at a specific address
func NewL2FeeCollector(addr web3.Address, provider *jsonrpc.Client) *L2FeeCollector {
	return &L2FeeCollector{c: contract.NewContract(addr, abiL2FeeCollector, provider)}
}

// Contract returns the contract object
func (_a *L2FeeCollector) Contract() *contract.Contract {
	return _a.c
}

// calls

// Owner calls the owner method in the solidity contract
func (_a *L2FeeCollector) Owner(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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
func (_a *L2FeeCollector) RenounceOwnership() *contract.Txn {
	return _a.c.Txn("renounceOwnership")
}

// TransferOwnership sends a transferOwnership transaction in the solidity contract
func (_a *L2FeeCollector) TransferOwnership(newOwner web3.Address) *contract.Txn {
	return _a.c.Txn("transferOwnership", newOwner)
}

// WithdrawERC20 sends a withdrawERC20 transaction in the solidity contract
func (_a *L2FeeCollector) WithdrawERC20(token web3.Address) *contract.Txn {
	return _a.c.Txn("withdrawERC20", token)
}

// WithdrawEth sends a withdrawEth transaction in the solidity contract
func (_a *L2FeeCollector) WithdrawEth() *contract.Txn {
	return _a.c.Txn("withdrawEth")
}

// events

func (_a *L2FeeCollector) OwnershipTransferredTopicFilter(previousOwner []web3.Address, newOwner []web3.Address) [][]web3.Hash {

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

func (_a *L2FeeCollector) FilterOwnershipTransferredEvent(previousOwner []web3.Address, newOwner []web3.Address, startBlock uint64, endBlock ...uint64) ([]*OwnershipTransferredEvent, error) {
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
