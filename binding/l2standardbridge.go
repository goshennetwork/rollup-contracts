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

// L2StandardBridge is a solidity contract
type L2StandardBridge struct {
	c *contract.Contract
}

// DeployL2StandardBridge deploys a new L2StandardBridge contract
func DeployL2StandardBridge(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiL2StandardBridge, binL2StandardBridge)
}

// NewL2StandardBridge creates a new instance of the contract at a specific address
func NewL2StandardBridge(addr web3.Address, provider *jsonrpc.Client) *L2StandardBridge {
	return &L2StandardBridge{c: contract.NewContract(addr, abiL2StandardBridge, provider)}
}

// Contract returns the contract object
func (_a *L2StandardBridge) Contract() *contract.Contract {
	return _a.c
}

// calls

// CrossLayerWitness calls the crossLayerWitness method in the solidity contract
func (_a *L2StandardBridge) CrossLayerWitness(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// L1TokenBridge calls the l1TokenBridge method in the solidity contract
func (_a *L2StandardBridge) L1TokenBridge(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("l1TokenBridge", web3.EncodeBlock(block...))
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

// FinalizeERC20Deposit sends a finalizeERC20Deposit transaction in the solidity contract
func (_a *L2StandardBridge) FinalizeERC20Deposit(l1Token web3.Address, l2Token web3.Address, from web3.Address, to web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("finalizeERC20Deposit", l1Token, l2Token, from, to, amount, data)
}

// FinalizeETHDeposit sends a finalizeETHDeposit transaction in the solidity contract
func (_a *L2StandardBridge) FinalizeETHDeposit(from web3.Address, to web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("finalizeETHDeposit", from, to, amount, data)
}

// Initialize sends a initialize transaction in the solidity contract
func (_a *L2StandardBridge) Initialize(witness web3.Address) *contract.Txn {
	return _a.c.Txn("initialize", witness)
}

// Withdraw sends a withdraw transaction in the solidity contract
func (_a *L2StandardBridge) Withdraw(l2Token web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("withdraw", l2Token, amount, data)
}

// WithdrawETH sends a withdrawETH transaction in the solidity contract
func (_a *L2StandardBridge) WithdrawETH(data []byte) *contract.Txn {
	return _a.c.Txn("withdrawETH", data)
}

// WithdrawETHTo sends a withdrawETHTo transaction in the solidity contract
func (_a *L2StandardBridge) WithdrawETHTo(to web3.Address, data []byte) *contract.Txn {
	return _a.c.Txn("withdrawETHTo", to, data)
}

// WithdrawTo sends a withdrawTo transaction in the solidity contract
func (_a *L2StandardBridge) WithdrawTo(l2Token web3.Address, to web3.Address, amount *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("withdrawTo", l2Token, to, amount, data)
}

// events

func (_a *L2StandardBridge) DepositFailedTopicFilter(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address) [][]web3.Hash {

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
	query = append(query, []interface{}{DepositFailedEventID}, l1TokenRule, l2TokenRule, fromRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2StandardBridge) FilterDepositFailedEvent(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address, startBlock uint64, endBlock ...uint64) ([]*DepositFailedEvent, error) {
	topic := _a.DepositFailedTopicFilter(l1Token, l2Token, from)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*DepositFailedEvent, 0)
	evts := _a.c.Abi.Events["DepositFailed"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem DepositFailedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L2StandardBridge) DepositFinalizedTopicFilter(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address) [][]web3.Hash {

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
	query = append(query, []interface{}{DepositFinalizedEventID}, l1TokenRule, l2TokenRule, fromRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2StandardBridge) FilterDepositFinalizedEvent(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address, startBlock uint64, endBlock ...uint64) ([]*DepositFinalizedEvent, error) {
	topic := _a.DepositFinalizedTopicFilter(l1Token, l2Token, from)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*DepositFinalizedEvent, 0)
	evts := _a.c.Abi.Events["DepositFinalized"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem DepositFinalizedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L2StandardBridge) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2StandardBridge) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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

func (_a *L2StandardBridge) WithdrawalInitiatedTopicFilter(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address) [][]web3.Hash {

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
	query = append(query, []interface{}{WithdrawalInitiatedEventID}, l1TokenRule, l2TokenRule, fromRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2StandardBridge) FilterWithdrawalInitiatedEvent(l1Token []web3.Address, l2Token []web3.Address, from []web3.Address, startBlock uint64, endBlock ...uint64) ([]*WithdrawalInitiatedEvent, error) {
	topic := _a.WithdrawalInitiatedTopicFilter(l1Token, l2Token, from)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*WithdrawalInitiatedEvent, 0)
	evts := _a.c.Abi.Events["WithdrawalInitiated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem WithdrawalInitiatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
