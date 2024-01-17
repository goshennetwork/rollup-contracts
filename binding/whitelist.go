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

// Whitelist is a solidity contract
type Whitelist struct {
	c *contract.Contract
}

// DeployWhitelist deploys a new Whitelist contract
func DeployWhitelist(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiWhitelist, binWhitelist)
}

// NewWhitelist creates a new instance of the contract at a specific address
func NewWhitelist(addr web3.Address, provider *jsonrpc.Client) *Whitelist {
	return &Whitelist{c: contract.NewContract(addr, abiWhitelist, provider)}
}

// Contract returns the contract object
func (_a *Whitelist) Contract() *contract.Contract {
	return _a.c
}

// calls

// CanChallenge calls the canChallenge method in the solidity contract
func (_a *Whitelist) CanChallenge(val0 web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("canChallenge", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// CanPropose calls the canPropose method in the solidity contract
func (_a *Whitelist) CanPropose(val0 web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("canPropose", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// CanSequence calls the canSequence method in the solidity contract
func (_a *Whitelist) CanSequence(val0 web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("canSequence", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Resolver calls the resolver method in the solidity contract
func (_a *Whitelist) Resolver(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("resolver", web3.EncodeBlock(block...))
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
func (_a *Whitelist) Initialize(resolver web3.Address) *contract.Txn {
	return _a.c.Txn("initialize", resolver)
}

// SetChallenger sends a setChallenger transaction in the solidity contract
func (_a *Whitelist) SetChallenger(challenger web3.Address, enabled bool) *contract.Txn {
	return _a.c.Txn("setChallenger", challenger, enabled)
}

// SetProposer sends a setProposer transaction in the solidity contract
func (_a *Whitelist) SetProposer(proposer web3.Address, enabled bool) *contract.Txn {
	return _a.c.Txn("setProposer", proposer, enabled)
}

// SetSequencer sends a setSequencer transaction in the solidity contract
func (_a *Whitelist) SetSequencer(sequencer web3.Address, enabled bool) *contract.Txn {
	return _a.c.Txn("setSequencer", sequencer, enabled)
}

// events

func (_a *Whitelist) ChallengerUpdatedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{ChallengerUpdatedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Whitelist) FilterChallengerUpdatedEvent(startBlock uint64, endBlock ...uint64) ([]*ChallengerUpdatedEvent, error) {
	topic := _a.ChallengerUpdatedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ChallengerUpdatedEvent, 0)
	evts := _a.c.Abi.Events["ChallengerUpdated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ChallengerUpdatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *Whitelist) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Whitelist) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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

func (_a *Whitelist) ProposerUpdatedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{ProposerUpdatedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Whitelist) FilterProposerUpdatedEvent(startBlock uint64, endBlock ...uint64) ([]*ProposerUpdatedEvent, error) {
	topic := _a.ProposerUpdatedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ProposerUpdatedEvent, 0)
	evts := _a.c.Abi.Events["ProposerUpdated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ProposerUpdatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *Whitelist) SequencerUpdatedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{SequencerUpdatedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Whitelist) FilterSequencerUpdatedEvent(startBlock uint64, endBlock ...uint64) ([]*SequencerUpdatedEvent, error) {
	topic := _a.SequencerUpdatedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*SequencerUpdatedEvent, 0)
	evts := _a.c.Abi.Events["SequencerUpdated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem SequencerUpdatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
