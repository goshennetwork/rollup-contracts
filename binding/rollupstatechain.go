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

// RollupStateChain is a solidity contract
type RollupStateChain struct {
	c *contract.Contract
}

// DeployRollupStateChain deploys a new RollupStateChain contract
func DeployRollupStateChain(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiRollupStateChain, binRollupStateChain)
}

// NewRollupStateChain creates a new instance of the contract at a specific address
func NewRollupStateChain(addr web3.Address, provider *jsonrpc.Client) *RollupStateChain {
	return &RollupStateChain{c: contract.NewContract(addr, abiRollupStateChain, provider)}
}

// Contract returns the contract object
func (_a *RollupStateChain) Contract() *contract.Contract {
	return _a.c
}

// calls

// FraudProofWindow calls the fraudProofWindow method in the solidity contract
func (_a *RollupStateChain) FraudProofWindow(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("fraudProofWindow", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// IsStateConfirmed calls the isStateConfirmed method in the solidity contract
func (_a *RollupStateChain) IsStateConfirmed(stateInfo StateInfo, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("isStateConfirmed", web3.EncodeBlock(block...), stateInfo)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["_confirmed"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// TotalSubmittedState calls the totalSubmittedState method in the solidity contract
func (_a *RollupStateChain) TotalSubmittedState(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("totalSubmittedState", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// VerifyStateInfo calls the verifyStateInfo method in the solidity contract
func (_a *RollupStateChain) VerifyStateInfo(stateInfo StateInfo, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("verifyStateInfo", web3.EncodeBlock(block...), stateInfo)
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

// AppendStateBatch sends a appendStateBatch transaction in the solidity contract
func (_a *RollupStateChain) AppendStateBatch(blockHashes [][32]byte, startAt uint64) *contract.Txn {
	return _a.c.Txn("appendStateBatch", blockHashes, startAt)
}

// Initialize sends a initialize transaction in the solidity contract
func (_a *RollupStateChain) Initialize(addressResolver web3.Address, fraudProofWindow *big.Int) *contract.Txn {
	return _a.c.Txn("initialize", addressResolver, fraudProofWindow)
}

// RollbackStateBefore sends a rollbackStateBefore transaction in the solidity contract
func (_a *RollupStateChain) RollbackStateBefore(stateInfo StateInfo) *contract.Txn {
	return _a.c.Txn("rollbackStateBefore", stateInfo)
}

// events

func (_a *RollupStateChain) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *RollupStateChain) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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

func (_a *RollupStateChain) StateBatchAppendedTopicFilter(startIndex []uint64, proposer []web3.Address) [][]web3.Hash {

	var startIndexRule []interface{}
	for _, _startIndexItem := range startIndex {
		startIndexRule = append(startIndexRule, _startIndexItem)
	}

	var proposerRule []interface{}
	for _, _proposerItem := range proposer {
		proposerRule = append(proposerRule, _proposerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{StateBatchAppendedEventID}, startIndexRule, proposerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *RollupStateChain) FilterStateBatchAppendedEvent(startIndex []uint64, proposer []web3.Address, startBlock uint64, endBlock ...uint64) ([]*StateBatchAppendedEvent, error) {
	topic := _a.StateBatchAppendedTopicFilter(startIndex, proposer)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*StateBatchAppendedEvent, 0)
	evts := _a.c.Abi.Events["StateBatchAppended"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem StateBatchAppendedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *RollupStateChain) StateRollbackedTopicFilter(stateIndex []uint64, blockHash [][32]byte) [][]web3.Hash {

	var stateIndexRule []interface{}
	for _, _stateIndexItem := range stateIndex {
		stateIndexRule = append(stateIndexRule, _stateIndexItem)
	}

	var blockHashRule []interface{}
	for _, _blockHashItem := range blockHash {
		blockHashRule = append(blockHashRule, _blockHashItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{StateRollbackedEventID}, stateIndexRule, blockHashRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *RollupStateChain) FilterStateRollbackedEvent(stateIndex []uint64, blockHash [][32]byte, startBlock uint64, endBlock ...uint64) ([]*StateRollbackedEvent, error) {
	topic := _a.StateRollbackedTopicFilter(stateIndex, blockHash)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*StateRollbackedEvent, 0)
	evts := _a.c.Abi.Events["StateRollbacked"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem StateRollbackedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
