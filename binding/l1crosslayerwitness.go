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

// L1CrossLayerWitness is a solidity contract
type L1CrossLayerWitness struct {
	c *contract.Contract
}

// DeployL1CrossLayerWitness deploys a new L1CrossLayerWitness contract
func DeployL1CrossLayerWitness(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiL1CrossLayerWitness, binL1CrossLayerWitness)
}

// NewL1CrossLayerWitness creates a new instance of the contract at a specific address
func NewL1CrossLayerWitness(addr web3.Address, provider *jsonrpc.Client) *L1CrossLayerWitness {
	return &L1CrossLayerWitness{c: contract.NewContract(addr, abiL1CrossLayerWitness, provider)}
}

// Contract returns the contract object
func (_a *L1CrossLayerWitness) Contract() *contract.Contract {
	return _a.c
}

// calls

// BlockedMessages calls the blockedMessages method in the solidity contract
func (_a *L1CrossLayerWitness) BlockedMessages(val0 [32]byte, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("blockedMessages", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// CrossLayerSender calls the crossLayerSender method in the solidity contract
func (_a *L1CrossLayerWitness) CrossLayerSender(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// IsMessageSucceed calls the isMessageSucceed method in the solidity contract
func (_a *L1CrossLayerWitness) IsMessageSucceed(messageHash [32]byte, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("isMessageSucceed", web3.EncodeBlock(block...), messageHash)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// MmrRoot calls the mmrRoot method in the solidity contract
func (_a *L1CrossLayerWitness) MmrRoot(block ...web3.BlockNumber) (retval0 [32]byte, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("mmrRoot", web3.EncodeBlock(block...))
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
func (_a *L1CrossLayerWitness) SuccessRelayedMessages(val0 [32]byte, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("successRelayedMessages", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// TotalSize calls the totalSize method in the solidity contract
func (_a *L1CrossLayerWitness) TotalSize(block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("totalSize", web3.EncodeBlock(block...))
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

// AllowMessage sends a allowMessage transaction in the solidity contract
func (_a *L1CrossLayerWitness) AllowMessage(messageHashes [][32]byte) *contract.Txn {
	return _a.c.Txn("allowMessage", messageHashes)
}

// BlockMessage sends a blockMessage transaction in the solidity contract
func (_a *L1CrossLayerWitness) BlockMessage(messageHashes [][32]byte) *contract.Txn {
	return _a.c.Txn("blockMessage", messageHashes)
}

// Initialize sends a initialize transaction in the solidity contract
func (_a *L1CrossLayerWitness) Initialize(addressResolver web3.Address) *contract.Txn {
	return _a.c.Txn("initialize", addressResolver)
}

// RelayMessage sends a relayMessage transaction in the solidity contract
func (_a *L1CrossLayerWitness) RelayMessage(target web3.Address, sender web3.Address, message []byte, messageIndex uint64, rlpHeader []byte, stateInfo StateInfo, proof [][32]byte) *contract.Txn {
	return _a.c.Txn("relayMessage", target, sender, message, messageIndex, rlpHeader, stateInfo, proof)
}

// SendMessage sends a sendMessage transaction in the solidity contract
func (_a *L1CrossLayerWitness) SendMessage(target web3.Address, message []byte) *contract.Txn {
	return _a.c.Txn("sendMessage", target, message)
}

// events

func (_a *L1CrossLayerWitness) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1CrossLayerWitness) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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

func (_a *L1CrossLayerWitness) MessageAllowedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{MessageAllowedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1CrossLayerWitness) FilterMessageAllowedEvent(startBlock uint64, endBlock ...uint64) ([]*MessageAllowedEvent, error) {
	topic := _a.MessageAllowedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*MessageAllowedEvent, 0)
	evts := _a.c.Abi.Events["MessageAllowed"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem MessageAllowedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L1CrossLayerWitness) MessageBlockedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{MessageBlockedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1CrossLayerWitness) FilterMessageBlockedEvent(startBlock uint64, endBlock ...uint64) ([]*MessageBlockedEvent, error) {
	topic := _a.MessageBlockedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*MessageBlockedEvent, 0)
	evts := _a.c.Abi.Events["MessageBlocked"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem MessageBlockedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *L1CrossLayerWitness) MessageRelayFailedTopicFilter(msgHash [][32]byte) [][]web3.Hash {

	var msgHashRule []interface{}
	for _, _msgHashItem := range msgHash {
		msgHashRule = append(msgHashRule, _msgHashItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{MessageRelayFailedEventID}, msgHashRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *L1CrossLayerWitness) FilterMessageRelayFailedEvent(msgHash [][32]byte, startBlock uint64, endBlock ...uint64) ([]*MessageRelayFailedEvent, error) {
	topic := _a.MessageRelayFailedTopicFilter(msgHash)

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

func (_a *L1CrossLayerWitness) MessageRelayedTopicFilter(messageIndex []uint64, msgHash [][32]byte) [][]web3.Hash {

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

func (_a *L1CrossLayerWitness) FilterMessageRelayedEvent(messageIndex []uint64, msgHash [][32]byte, startBlock uint64, endBlock ...uint64) ([]*MessageRelayedEvent, error) {
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

func (_a *L1CrossLayerWitness) MessageSentTopicFilter(messageIndex []uint64, target []web3.Address, sender []web3.Address) [][]web3.Hash {

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

func (_a *L1CrossLayerWitness) FilterMessageSentEvent(messageIndex []uint64, target []web3.Address, sender []web3.Address, startBlock uint64, endBlock ...uint64) ([]*MessageSentEvent, error) {
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
