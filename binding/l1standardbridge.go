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

// L1StandardBridge is a solidity contract
type L1StandardBridge struct {
	c *contract.Contract
}

// DeployL1StandardBridge deploys a new L1StandardBridge contract
func DeployL1StandardBridge(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiL1StandardBridge, binL1StandardBridge)
}

// NewL1StandardBridge creates a new instance of the contract at a specific address
func NewL1StandardBridge(addr web3.Address, provider *jsonrpc.Client) *L1StandardBridge {
	return &L1StandardBridge{c: contract.NewContract(addr, abiL1StandardBridge, provider)}
}

// Contract returns the contract object
func (_a *L1StandardBridge) Contract() *contract.Contract {
	return _a.c
}

// calls

// CrossLayerWitness calls the crossLayerWitness method in the solidity contract
func (_a *L1StandardBridge) CrossLayerWitness(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("crossLayerWitness", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Deposits calls the deposits method in the solidity contract
func (_a *L1StandardBridge) Deposits(val0 web3.Address, val1 web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("deposits", web3.EncodeBlock(block...), val0, val1)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// L2TokenBridge calls the l2TokenBridge method in the solidity contract
func (_a *L1StandardBridge) L2TokenBridge(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("l2TokenBridge", web3.EncodeBlock(block...))
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

// DepositERC20 sends a depositERC20 transaction in the solidity contract
func (_a *L1StandardBridge) DepositERC20(l1Token web3.Address, l2Token web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("depositERC20", l1Token, l2Token, amount, data)
}

// DepositERC20To sends a depositERC20To transaction in the solidity contract
func (_a *L1StandardBridge) DepositERC20To(l1Token web3.Address, l2Token web3.Address, to web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("depositERC20To", l1Token, l2Token, to, amount, data)
}

// DepositETH sends a depositETH transaction in the solidity contract
func (_a *L1StandardBridge) DepositETH(data []byte) *contract.Txn {
	return _a.c.Txn("depositETH", data)
}

// DepositETHTo sends a depositETHTo transaction in the solidity contract
func (_a *L1StandardBridge) DepositETHTo(to web3.Address, data []byte) *contract.Txn {
	return _a.c.Txn("depositETHTo", to, data)
}

// DonateETH sends a donateETH transaction in the solidity contract
func (_a *L1StandardBridge) DonateETH() *contract.Txn {
	return _a.c.Txn("donateETH")
}

// FinalizeERC20Withdrawal sends a finalizeERC20Withdrawal transaction in the solidity contract
func (_a *L1StandardBridge) FinalizeERC20Withdrawal(l1Token web3.Address, l2Token web3.Address, from web3.Address, to web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("finalizeERC20Withdrawal", l1Token, l2Token, from, to, amount, data)
}

// FinalizeETHWithdrawal sends a finalizeETHWithdrawal transaction in the solidity contract
func (_a *L1StandardBridge) FinalizeETHWithdrawal(from web3.Address, to web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("finalizeETHWithdrawal", from, to, amount, data)
}

// Initialize sends a initialize transaction in the solidity contract
func (_a *L1StandardBridge) Initialize(l1witness web3.Address, l2TokenBridge web3.Address) *contract.Txn {
	return _a.c.Txn("initialize", l1witness, l2TokenBridge)
}

// events

func (_a *L1StandardBridge) ERC20DepositInitiatedTopicFilter(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address) [][]web3.Hash {

	var l1TokenRule []interface{}
	for _, _l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, _l1TokenItem)
	}

	var l2TokenRule []interface{}
	for _, _l2TokenItem := range l2Token {
		l2TokenRule = append(l2TokenRule, _l2TokenItem)
	}

	var fromRule []interface{}
	for _, _fromItem := range from {
		fromRule = append(fromRule, _fromItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{ERC20DepositInitiatedEventID}, l1TokenRule, l2TokenRule, fromRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1StandardBridge) FilterERC20DepositInitiatedEvent(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address, startBlock uint64, endBlock ...uint64) ([]*ERC20DepositInitiatedEvent, error) {
	topic := _a.ERC20DepositInitiatedTopicFilter(l1Token, l2Token, from)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ERC20DepositInitiatedEvent, 0)
	evts := _a.c.Abi.Events["ERC20DepositInitiated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ERC20DepositInitiatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L1StandardBridge) ERC20WithdrawalFinalizedTopicFilter(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address) [][]web3.Hash {

	var l1TokenRule []interface{}
	for _, _l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, _l1TokenItem)
	}

	var l2TokenRule []interface{}
	for _, _l2TokenItem := range l2Token {
		l2TokenRule = append(l2TokenRule, _l2TokenItem)
	}

	var fromRule []interface{}
	for _, _fromItem := range from {
		fromRule = append(fromRule, _fromItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{ERC20WithdrawalFinalizedEventID}, l1TokenRule, l2TokenRule, fromRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1StandardBridge) FilterERC20WithdrawalFinalizedEvent(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address, startBlock uint64, endBlock ...uint64) ([]*ERC20WithdrawalFinalizedEvent, error) {
	topic := _a.ERC20WithdrawalFinalizedTopicFilter(l1Token, l2Token, from)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ERC20WithdrawalFinalizedEvent, 0)
	evts := _a.c.Abi.Events["ERC20WithdrawalFinalized"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ERC20WithdrawalFinalizedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L1StandardBridge) ETHDepositInitiatedTopicFilter(from []web3.Address, to []web3.Address) [][]web3.Hash {

	var fromRule []interface{}
	for _, _fromItem := range from {
		fromRule = append(fromRule, _fromItem)
	}

	var toRule []interface{}
	for _, _toItem := range to {
		toRule = append(toRule, _toItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{ETHDepositInitiatedEventID}, fromRule, toRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1StandardBridge) FilterETHDepositInitiatedEvent(from []web3.Address, to []web3.Address, startBlock uint64, endBlock ...uint64) ([]*ETHDepositInitiatedEvent, error) {
	topic := _a.ETHDepositInitiatedTopicFilter(from, to)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ETHDepositInitiatedEvent, 0)
	evts := _a.c.Abi.Events["ETHDepositInitiated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ETHDepositInitiatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L1StandardBridge) ETHWithdrawalFinalizedTopicFilter(from []web3.Address, to []web3.Address) [][]web3.Hash {

	var fromRule []interface{}
	for _, _fromItem := range from {
		fromRule = append(fromRule, _fromItem)
	}

	var toRule []interface{}
	for _, _toItem := range to {
		toRule = append(toRule, _toItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{ETHWithdrawalFinalizedEventID}, fromRule, toRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1StandardBridge) FilterETHWithdrawalFinalizedEvent(from []web3.Address, to []web3.Address, startBlock uint64, endBlock ...uint64) ([]*ETHWithdrawalFinalizedEvent, error) {
	topic := _a.ETHWithdrawalFinalizedTopicFilter(from, to)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ETHWithdrawalFinalizedEvent, 0)
	evts := _a.c.Abi.Events["ETHWithdrawalFinalized"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ETHWithdrawalFinalizedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L1StandardBridge) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1StandardBridge) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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
