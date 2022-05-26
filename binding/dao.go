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

// DAO is a solidity contract
type DAO struct {
	c *contract.Contract
}

// DeployDAO deploys a new DAO contract
func DeployDAO(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiDAO, binDAO)
}

// NewDAO creates a new instance of the contract at a specific address
func NewDAO(addr web3.Address, provider *jsonrpc.Client) *DAO {
	return &DAO{c: contract.NewContract(addr, abiDAO, provider)}
}

// Contract returns the contract object
func (_a *DAO) Contract() *contract.Contract {
	return _a.c
}

// calls

// ChallengerWhitelist calls the challengerWhitelist method in the solidity contract
func (_a *DAO) ChallengerWhitelist(val0 web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("challengerWhitelist", web3.EncodeBlock(block...), val0)
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
func (_a *DAO) Owner(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// ProposerWhitelist calls the proposerWhitelist method in the solidity contract
func (_a *DAO) ProposerWhitelist(val0 web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("proposerWhitelist", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// SequencerWhitelist calls the sequencerWhitelist method in the solidity contract
func (_a *DAO) SequencerWhitelist(val0 web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("sequencerWhitelist", web3.EncodeBlock(block...), val0)
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
func (_a *DAO) Initialize() *contract.Txn {
	return _a.c.Txn("initialize")
}

// RenounceOwnership sends a renounceOwnership transaction in the solidity contract
func (_a *DAO) RenounceOwnership() *contract.Txn {
	return _a.c.Txn("renounceOwnership")
}

// SetChallengerWhitelist sends a setChallengerWhitelist transaction in the solidity contract
func (_a *DAO) SetChallengerWhitelist(challenger web3.Address, enabled bool) *contract.Txn {
	return _a.c.Txn("setChallengerWhitelist", challenger, enabled)
}

// SetProposerWhitelist sends a setProposerWhitelist transaction in the solidity contract
func (_a *DAO) SetProposerWhitelist(proposer web3.Address, enabled bool) *contract.Txn {
	return _a.c.Txn("setProposerWhitelist", proposer, enabled)
}

// SetSequencerWhitelist sends a setSequencerWhitelist transaction in the solidity contract
func (_a *DAO) SetSequencerWhitelist(sequencer web3.Address, enabled bool) *contract.Txn {
	return _a.c.Txn("setSequencerWhitelist", sequencer, enabled)
}

// TransferERC20 sends a transferERC20 transaction in the solidity contract
func (_a *DAO) TransferERC20(token web3.Address, to web3.Address, amount *big.Int) *contract.Txn {
	return _a.c.Txn("transferERC20", token, to, amount)
}

// TransferOwnership sends a transferOwnership transaction in the solidity contract
func (_a *DAO) TransferOwnership(newOwner web3.Address) *contract.Txn {
	return _a.c.Txn("transferOwnership", newOwner)
}

// events

func (_a *DAO) ChallengerWhitelistUpdatedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{ChallengerWhitelistUpdatedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *DAO) FilterChallengerWhitelistUpdatedEvent(startBlock uint64, endBlock ...uint64) ([]*ChallengerWhitelistUpdatedEvent, error) {
	topic := _a.ChallengerWhitelistUpdatedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ChallengerWhitelistUpdatedEvent, 0)
	evts := _a.c.Abi.Events["ChallengerWhitelistUpdated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ChallengerWhitelistUpdatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *DAO) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *DAO) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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

func (_a *DAO) OwnershipTransferredTopicFilter(previousOwner []web3.Address, newOwner []web3.Address) [][]web3.Hash {

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

func (_a *DAO) FilterOwnershipTransferredEvent(previousOwner []web3.Address, newOwner []web3.Address, startBlock uint64, endBlock ...uint64) ([]*OwnershipTransferredEvent, error) {
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

func (_a *DAO) ProposerWhitelistUpdatedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{ProposerWhitelistUpdatedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *DAO) FilterProposerWhitelistUpdatedEvent(startBlock uint64, endBlock ...uint64) ([]*ProposerWhitelistUpdatedEvent, error) {
	topic := _a.ProposerWhitelistUpdatedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ProposerWhitelistUpdatedEvent, 0)
	evts := _a.c.Abi.Events["ProposerWhitelistUpdated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ProposerWhitelistUpdatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *DAO) SequencerWhitelistUpdatedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{SequencerWhitelistUpdatedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *DAO) FilterSequencerWhitelistUpdatedEvent(startBlock uint64, endBlock ...uint64) ([]*SequencerWhitelistUpdatedEvent, error) {
	topic := _a.SequencerWhitelistUpdatedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*SequencerWhitelistUpdatedEvent, 0)
	evts := _a.c.Abi.Events["SequencerWhitelistUpdated"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem SequencerWhitelistUpdatedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
