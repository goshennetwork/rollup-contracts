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

// ChallengeFactory is a solidity contract
type ChallengeFactory struct {
	c *contract.Contract
}

// DeployChallengeFactory deploys a new ChallengeFactory contract
func DeployChallengeFactory(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiChallengeFactory, binChallengeFactory)
}

// NewChallengeFactory creates a new instance of the contract at a specific address
func NewChallengeFactory(addr web3.Address, provider *jsonrpc.Client) *ChallengeFactory {
	return &ChallengeFactory{c: contract.NewContract(addr, abiChallengeFactory, provider)}
}

// Contract returns the contract object
func (_a *ChallengeFactory) Contract() *contract.Contract {
	return _a.c
}

// calls

// BlockLimitPerRound calls the blockLimitPerRound method in the solidity contract
func (_a *ChallengeFactory) BlockLimitPerRound(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("blockLimitPerRound", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// ChallengeBeacon calls the challengeBeacon method in the solidity contract
func (_a *ChallengeFactory) ChallengeBeacon(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("challengeBeacon", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// ChallengerDeposit calls the challengerDeposit method in the solidity contract
func (_a *ChallengeFactory) ChallengerDeposit(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("challengerDeposit", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Dao calls the dao method in the solidity contract
func (_a *ChallengeFactory) Dao(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("dao", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Executor calls the executor method in the solidity contract
func (_a *ChallengeFactory) Executor(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("executor", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// GetChallengedContract calls the getChallengedContract method in the solidity contract
func (_a *ChallengeFactory) GetChallengedContract(arg0 [32]byte, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("getChallengedContract", web3.EncodeBlock(block...), arg0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// IsChallengeContract calls the isChallengeContract method in the solidity contract
func (_a *ChallengeFactory) IsChallengeContract(arg0 web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("isChallengeContract", web3.EncodeBlock(block...), arg0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// RollupStateChain calls the rollupStateChain method in the solidity contract
func (_a *ChallengeFactory) RollupStateChain(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("rollupStateChain", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// StakingManager calls the stakingManager method in the solidity contract
func (_a *ChallengeFactory) StakingManager(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("stakingManager", web3.EncodeBlock(block...))
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
func (_a *ChallengeFactory) Initialize(resolver web3.Address, beacon web3.Address, blockLimitPerRound *big.Int, challengerDeposit *big.Int) *contract.Txn {
	return _a.c.Txn("initialize", resolver, beacon, blockLimitPerRound, challengerDeposit)
}

// NewChallenge sends a newChallenge transaction in the solidity contract
func (_a *ChallengeFactory) NewChallenge(challengedStateInfo StateInfo, parentStateInfo StateInfo) *contract.Txn {
	return _a.c.Txn("newChallenge", challengedStateInfo, parentStateInfo)
}

// events

func (_a *ChallengeFactory) ChallengeStartedTopicFilter(l2BlockN []*big.Int, proposer []web3.Address) [][]web3.Hash {

	var l2BlockNRule []interface{}
	for _, _l2BlockNItem := range l2BlockN {
		l2BlockNRule = append(l2BlockNRule, _l2BlockNItem)
	}

	var proposerRule []interface{}
	for _, _proposerItem := range proposer {
		proposerRule = append(proposerRule, _proposerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{ChallengeStartedEventID}, l2BlockNRule, proposerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *ChallengeFactory) FilterChallengeStartedEvent(l2BlockN []*big.Int, proposer []web3.Address, startBlock uint64, endBlock ...uint64) ([]*ChallengeStartedEvent, error) {
	topic := _a.ChallengeStartedTopicFilter(l2BlockN, proposer)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ChallengeStartedEvent, 0)
	evts := _a.c.Abi.Events["ChallengeStarted"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ChallengeStartedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *ChallengeFactory) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *ChallengeFactory) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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
