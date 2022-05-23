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

// RollupInputChain is a solidity contract
type RollupInputChain struct {
	c *contract.Contract
}

// DeployRollupInputChain deploys a new RollupInputChain contract
func DeployRollupInputChain(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiRollupInputChain, binRollupInputChain, args...)
}

// NewRollupInputChain creates a new instance of the contract at a specific address
func NewRollupInputChain(addr web3.Address, provider *jsonrpc.Client) *RollupInputChain {
	return &RollupInputChain{c: contract.NewContract(addr, abiRollupInputChain, provider)}
}

// Contract returns the contract object
func (_a *RollupInputChain) Contract() *contract.Contract {
	return _a.c
}

// calls

// MAXCROSSLAYERTXSIZE calls the MAX_CROSS_LAYER_TX_SIZE method in the solidity contract
func (_a *RollupInputChain) MAXCROSSLAYERTXSIZE(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("MAX_CROSS_LAYER_TX_SIZE", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// MAXROLLUPTXSIZE calls the MAX_ROLLUP_TX_SIZE method in the solidity contract
func (_a *RollupInputChain) MAXROLLUPTXSIZE(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("MAX_ROLLUP_TX_SIZE", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// MINROLLUPTXGAS calls the MIN_ROLLUP_TX_GAS method in the solidity contract
func (_a *RollupInputChain) MINROLLUPTXGAS(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("MIN_ROLLUP_TX_GAS", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// ChainHeight calls the chainHeight method in the solidity contract
func (_a *RollupInputChain) ChainHeight(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("chainHeight", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// GetQueueTxInfo calls the getQueueTxInfo method in the solidity contract
func (_a *RollupInputChain) GetQueueTxInfo(queueIndex uint64, block ...web3.BlockNumber) (retval0 [32]byte, retval1 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("getQueueTxInfo", web3.EncodeBlock(block...), queueIndex)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}
	if err = mapstructure.Decode(out["1"], &retval1); err != nil {
		err = fmt.Errorf("failed to encode output at index 1")
	}

	return
}

// LastTimestamp calls the lastTimestamp method in the solidity contract
func (_a *RollupInputChain) LastTimestamp(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("lastTimestamp", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// MaxCrossLayerTxGasLimit calls the maxCrossLayerTxGasLimit method in the solidity contract
func (_a *RollupInputChain) MaxCrossLayerTxGasLimit(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("maxCrossLayerTxGasLimit", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// MaxEnqueueTxGasLimit calls the maxEnqueueTxGasLimit method in the solidity contract
func (_a *RollupInputChain) MaxEnqueueTxGasLimit(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("maxEnqueueTxGasLimit", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// PendingQueueIndex calls the pendingQueueIndex method in the solidity contract
func (_a *RollupInputChain) PendingQueueIndex(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("pendingQueueIndex", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// TotalQueue calls the totalQueue method in the solidity contract
func (_a *RollupInputChain) TotalQueue(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("totalQueue", web3.EncodeBlock(block...))
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

// AppendBatch sends a appendBatch transaction in the solidity contract
func (_a *RollupInputChain) AppendBatch() *contract.Txn {
	return _a.c.Txn("appendBatch")
}

// Enqueue sends a enqueue transaction in the solidity contract
func (_a *RollupInputChain) Enqueue(target web3.Address, gasLimit uint64, data []byte) *contract.Txn {
	return _a.c.Txn("enqueue", target, gasLimit, data)
}

// events

var TransactionAppendedEventID = crypto.Keccak256Hash([]byte("TransactionAppended(address,uint256,uint256,uint256,bytes32)"))

func (_a *RollupInputChain) TransactionAppendedTopicFilter(proposer []web3.Address, startQueueIndex []*big.Int, chainHeight []*big.Int) [][]web3.Hash {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	var startQueueIndexRule []interface{}
	for _, startQueueIndexItem := range startQueueIndex {
		startQueueIndexRule = append(startQueueIndexRule, startQueueIndexItem)
	}

	var chainHeightRule []interface{}
	for _, chainHeightItem := range chainHeight {
		chainHeightRule = append(chainHeightRule, chainHeightItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{TransactionAppendedEventID}, proposerRule, startQueueIndexRule, chainHeightRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *RollupInputChain) FilterTransactionAppendedEvent(proposer []web3.Address, startQueueIndex []*big.Int, chainHeight []*big.Int, startBlock uint64, endBlock ...uint64) ([]*TransactionAppendedEvent, error) {
	topic := _a.TransactionAppendedTopicFilter(proposer, startQueueIndex, chainHeight)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*TransactionAppendedEvent, 0)
	evts := _a.c.Abi.Events["TransactionAppended"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem TransactionAppendedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var TransactionEnqueuedEventID = crypto.Keccak256Hash([]byte("TransactionEnqueued(uint64,address,address,uint256,bytes,uint64)"))

func (_a *RollupInputChain) TransactionEnqueuedTopicFilter(queueIndex []uint64, from []web3.Address, to []web3.Address) [][]web3.Hash {

	var queueIndexRule []interface{}
	for _, queueIndexItem := range queueIndex {
		queueIndexRule = append(queueIndexRule, queueIndexItem)
	}

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{TransactionEnqueuedEventID}, queueIndexRule, fromRule, toRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *RollupInputChain) FilterTransactionEnqueuedEvent(queueIndex []uint64, from []web3.Address, to []web3.Address, startBlock uint64, endBlock ...uint64) ([]*TransactionEnqueuedEvent, error) {
	topic := _a.TransactionEnqueuedTopicFilter(queueIndex, from, to)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*TransactionEnqueuedEvent, 0)
	evts := _a.c.Abi.Events["TransactionEnqueued"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem TransactionEnqueuedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
