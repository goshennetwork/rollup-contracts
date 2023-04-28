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

// L2StandardTokenFactory is a solidity contract
type L2StandardTokenFactory struct {
	c *contract.Contract
}

// DeployL2StandardTokenFactory deploys a new L2StandardTokenFactory contract
func DeployL2StandardTokenFactory(provider *jsonrpc.Client, from web3.Address, l2StandardBridge web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiL2StandardTokenFactory, binL2StandardTokenFactory, l2StandardBridge)
}

// NewL2StandardTokenFactory creates a new instance of the contract at a specific address
func NewL2StandardTokenFactory(addr web3.Address, provider *jsonrpc.Client) *L2StandardTokenFactory {
	return &L2StandardTokenFactory{c: contract.NewContract(addr, abiL2StandardTokenFactory, provider)}
}

// Contract returns the contract object
func (_a *L2StandardTokenFactory) Contract() *contract.Contract {
	return _a.c
}

// calls

// txns

// CreateStandardL2Token sends a createStandardL2Token transaction in the solidity contract
func (_a *L2StandardTokenFactory) CreateStandardL2Token(l1Token web3.Address, name string, symbol string) *contract.Txn {
	return _a.c.Txn("createStandardL2Token", l1Token, name, symbol)
}

// events

func (_a *L2StandardTokenFactory) StandardL2TokenCreatedTopicFilter(l1Token []web3.Address, l2Token []web3.Address) [][]web3.Hash {

	var l1TokenRule []interface{}
	for _, _l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, _l1TokenItem)
	}

	var l2TokenRule []interface{}
	for _, _l2TokenItem := range l2Token {
		l2TokenRule = append(l2TokenRule, _l2TokenItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{StandardL2TokenCreatedEventID}, l1TokenRule, l2TokenRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2StandardTokenFactory) FilterStandardL2TokenCreatedEvent(l1Token []web3.Address, l2Token []web3.Address, startBlock uint64, endBlock ...uint64) ([]*StandardL2TokenCreatedEvent, error) {
	topic := _a.StandardL2TokenCreatedTopicFilter(l1Token, l2Token)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*StandardL2TokenCreatedEvent, 0)
	evts := _a.c.Abi.Events["StandardL2TokenCreated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem StandardL2TokenCreatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
