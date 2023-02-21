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

// L2CrossLayerWitness is a solidity contract
type L2CrossLayerWitness struct {
	c *contract.Contract
}

// DeployL2CrossLayerWitness deploys a new L2CrossLayerWitness contract
func DeployL2CrossLayerWitness(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiL2CrossLayerWitness, binL2CrossLayerWitness)
}

// NewL2CrossLayerWitness creates a new instance of the contract at a specific address
func NewL2CrossLayerWitness(addr web3.Address, provider *jsonrpc.Client) *L2CrossLayerWitness {
	return &L2CrossLayerWitness{c: contract.NewContract(addr, abiL2CrossLayerWitness, provider)}
}

// Contract returns the contract object
func (_a *L2CrossLayerWitness) Contract() *contract.Contract {
	return _a.c
}

// calls

// CrossLayerSender calls the crossLayerSender method in the solidity contract
func (_a *L2CrossLayerWitness) CrossLayerSender(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("crossLayerSender", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// MmrRoots calls the mmrRoots method in the solidity contract
func (_a *L2CrossLayerWitness) MmrRoots(arg0 uint64, block ...web3.BlockNumber) (retval0 [32]byte, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("mmrRoots", web3.EncodeBlock(block...), arg0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// SuccessRelayedMessages calls the successRelayedMessages method in the solidity contract
func (_a *L2CrossLayerWitness) SuccessRelayedMessages(arg0 [32]byte, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("successRelayedMessages", web3.EncodeBlock(block...), arg0)
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
func (_a *L2CrossLayerWitness) Initialize() *contract.Txn {
	return _a.c.Txn("initialize")
}

// RelayMessage sends a relayMessage transaction in the solidity contract
func (_a *L2CrossLayerWitness) RelayMessage(target web3.Address, sender web3.Address, message []byte, messageIndex uint64, mmrRoot [32]byte, mmrSize uint64) *contract.Txn {
	return _a.c.Txn("relayMessage", target, sender, message, messageIndex, mmrRoot, mmrSize)
}

// ReplayMessage sends a replayMessage transaction in the solidity contract
func (_a *L2CrossLayerWitness) ReplayMessage(target web3.Address, sender web3.Address, message []byte, messageIndex uint64, proof [][32]byte, mmrSize uint64) *contract.Txn {
	return _a.c.Txn("replayMessage", target, sender, message, messageIndex, proof, mmrSize)
}

// SendMessage sends a sendMessage transaction in the solidity contract
func (_a *L2CrossLayerWitness) SendMessage(target web3.Address, message []byte) *contract.Txn {
	return _a.c.Txn("sendMessage", target, message)
}

// events

func (_a *L2CrossLayerWitness) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2CrossLayerWitness) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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

func (_a *L2CrossLayerWitness) MessageRelayFailedTopicFilter(messageIndex []uint64, msgHash [][32]byte) [][]web3.Hash {

	var messageIndexRule []interface{}
	for _, _messageIndexItem := range messageIndex {
		messageIndexRule = append(messageIndexRule, _messageIndexItem)
	}

	var msgHashRule []interface{}
	for _, _msgHashItem := range msgHash {
		msgHashRule = append(msgHashRule, _msgHashItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{MessageRelayFailedEventID}, messageIndexRule, msgHashRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2CrossLayerWitness) FilterMessageRelayFailedEvent(messageIndex []uint64, msgHash [][32]byte, startBlock uint64, endBlock ...uint64) ([]*MessageRelayFailedEvent, error) {
	topic := _a.MessageRelayFailedTopicFilter(messageIndex, msgHash)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*MessageRelayFailedEvent, 0)
	evts := _a.c.Abi.Events["MessageRelayFailed"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem MessageRelayFailedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L2CrossLayerWitness) MessageRelayedTopicFilter(messageIndex []uint64, msgHash [][32]byte) [][]web3.Hash {

	var messageIndexRule []interface{}
	for _, _messageIndexItem := range messageIndex {
		messageIndexRule = append(messageIndexRule, _messageIndexItem)
	}

	var msgHashRule []interface{}
	for _, _msgHashItem := range msgHash {
		msgHashRule = append(msgHashRule, _msgHashItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{MessageRelayedEventID}, messageIndexRule, msgHashRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2CrossLayerWitness) FilterMessageRelayedEvent(messageIndex []uint64, msgHash [][32]byte, startBlock uint64, endBlock ...uint64) ([]*MessageRelayedEvent, error) {
	topic := _a.MessageRelayedTopicFilter(messageIndex, msgHash)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*MessageRelayedEvent, 0)
	evts := _a.c.Abi.Events["MessageRelayed"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem MessageRelayedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L2CrossLayerWitness) MessageSentTopicFilter(messageIndex []uint64, target []web3.Address, sender []web3.Address) [][]web3.Hash {

	var messageIndexRule []interface{}
	for _, _messageIndexItem := range messageIndex {
		messageIndexRule = append(messageIndexRule, _messageIndexItem)
	}

	var targetRule []interface{}
	for _, _targetItem := range target {
		targetRule = append(targetRule, _targetItem)
	}

	var senderRule []interface{}
	for _, _senderItem := range sender {
		senderRule = append(senderRule, _senderItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{MessageSentEventID}, messageIndexRule, targetRule, senderRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L2CrossLayerWitness) FilterMessageSentEvent(messageIndex []uint64, target []web3.Address, sender []web3.Address, startBlock uint64, endBlock ...uint64) ([]*MessageSentEvent, error) {
	topic := _a.MessageSentTopicFilter(messageIndex, target, sender)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*MessageSentEvent, 0)
	evts := _a.c.Abi.Events["MessageSent"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem MessageSentEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
