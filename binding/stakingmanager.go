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

// StakingManager is a solidity contract
type StakingManager struct {
	c *contract.Contract
}

// DeployStakingManager deploys a new StakingManager contract
func DeployStakingManager(provider *jsonrpc.Client, from web3.Address) *contract.Txn {
	return contract.DeployContract(provider, from, abiStakingManager, binStakingManager)
}

// NewStakingManager creates a new instance of the contract at a specific address
func NewStakingManager(addr web3.Address, provider *jsonrpc.Client) *StakingManager {
	return &StakingManager{c: contract.NewContract(addr, abiStakingManager, provider)}
}

// Contract returns the contract object
func (_a *StakingManager) Contract() *contract.Contract {
	return _a.c
}

// calls

// IsStaking calls the isStaking method in the solidity contract
func (_a *StakingManager) IsStaking(who web3.Address, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("isStaking", web3.EncodeBlock(block...), who)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Price calls the price method in the solidity contract
func (_a *StakingManager) Price(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("price", web3.EncodeBlock(block...))
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
func (_a *StakingManager) RollupStateChain(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// Token calls the token method in the solidity contract
func (_a *StakingManager) Token(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("token", web3.EncodeBlock(block...))
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

// Claim sends a claim transaction in the solidity contract
func (_a *StakingManager) Claim(proposer web3.Address, stateInfo StateInfo) *contract.Txn {
	return _a.c.Txn("claim", proposer, stateInfo)
}

// ClaimToGovernance sends a claimToGovernance transaction in the solidity contract
func (_a *StakingManager) ClaimToGovernance(proposer web3.Address, stateInfo StateInfo) *contract.Txn {
	return _a.c.Txn("claimToGovernance", proposer, stateInfo)
}

// Deposit sends a deposit transaction in the solidity contract
func (_a *StakingManager) Deposit() *contract.Txn {
	return _a.c.Txn("deposit")
}

// FinalizeWithdrawal sends a finalizeWithdrawal transaction in the solidity contract
func (_a *StakingManager) FinalizeWithdrawal(stateInfo StateInfo) *contract.Txn {
	return _a.c.Txn("finalizeWithdrawal", stateInfo)
}

// Initialize sends a initialize transaction in the solidity contract
func (_a *StakingManager) Initialize(DAOAddress web3.Address, challengeFactory web3.Address, rollupStateChain web3.Address, erc20 web3.Address, price *big.Int) *contract.Txn {
	return _a.c.Txn("initialize", DAOAddress, challengeFactory, rollupStateChain, erc20, price)
}

// Slash sends a slash transaction in the solidity contract
func (_a *StakingManager) Slash(chainHeight uint64, blockHash [32]byte, proposer web3.Address) *contract.Txn {
	return _a.c.Txn("slash", chainHeight, blockHash, proposer)
}

// StartWithdrawal sends a startWithdrawal transaction in the solidity contract
func (_a *StakingManager) StartWithdrawal() *contract.Txn {
	return _a.c.Txn("startWithdrawal")
}

// events

func (_a *StakingManager) DepositClaimedTopicFilter(proposer []web3.Address, receiver []web3.Address) [][]web3.Hash {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{DepositClaimedEventID}, proposerRule, receiverRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *StakingManager) FilterDepositClaimedEvent(proposer []web3.Address, receiver []web3.Address, startBlock uint64, endBlock ...uint64) ([]*DepositClaimedEvent, error) {
	topic := _a.DepositClaimedTopicFilter(proposer, receiver)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*DepositClaimedEvent, 0)
	evts := _a.c.Abi.Events["DepositClaimed"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem DepositClaimedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *StakingManager) DepositSlashedTopicFilter(proposer []web3.Address, challenger []web3.Address) [][]web3.Hash {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{DepositSlashedEventID}, proposerRule, challengerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *StakingManager) FilterDepositSlashedEvent(proposer []web3.Address, challenger []web3.Address, startBlock uint64, endBlock ...uint64) ([]*DepositSlashedEvent, error) {
	topic := _a.DepositSlashedTopicFilter(proposer, challenger)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*DepositSlashedEvent, 0)
	evts := _a.c.Abi.Events["DepositSlashed"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem DepositSlashedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *StakingManager) DepositedTopicFilter(proposer []web3.Address) [][]web3.Hash {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{DepositedEventID}, proposerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *StakingManager) FilterDepositedEvent(proposer []web3.Address, startBlock uint64, endBlock ...uint64) ([]*DepositedEvent, error) {
	topic := _a.DepositedTopicFilter(proposer)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*DepositedEvent, 0)
	evts := _a.c.Abi.Events["Deposited"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem DepositedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *StakingManager) InitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{InitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *StakingManager) FilterInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*InitializedEvent, error) {
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

func (_a *StakingManager) WithdrawFinalizedTopicFilter(proposer []web3.Address) [][]web3.Hash {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{WithdrawFinalizedEventID}, proposerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *StakingManager) FilterWithdrawFinalizedEvent(proposer []web3.Address, startBlock uint64, endBlock ...uint64) ([]*WithdrawFinalizedEvent, error) {
	topic := _a.WithdrawFinalizedTopicFilter(proposer)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*WithdrawFinalizedEvent, 0)
	evts := _a.c.Abi.Events["WithdrawFinalized"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem WithdrawFinalizedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

func (_a *StakingManager) WithdrawStartedTopicFilter(proposer []web3.Address) [][]web3.Hash {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{WithdrawStartedEventID}, proposerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *StakingManager) FilterWithdrawStartedEvent(proposer []web3.Address, startBlock uint64, endBlock ...uint64) ([]*WithdrawStartedEvent, error) {
	topic := _a.WithdrawStartedTopicFilter(proposer)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*WithdrawStartedEvent, 0)
	evts := _a.c.Abi.Events["WithdrawStarted"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem WithdrawStartedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
