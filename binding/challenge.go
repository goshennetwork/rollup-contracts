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

// Challenge is a solidity contract
type Challenge struct {
	c *contract.Contract
}

// DeployChallenge deploys a new Challenge contract
func DeployChallenge(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiChallenge, binChallenge, args...)
}

// NewChallenge creates a new instance of the contract at a specific address
func NewChallenge(addr web3.Address, provider *jsonrpc.Client) *Challenge {
	return &Challenge{c: contract.NewContract(addr, abiChallenge, provider)}
}

// Contract returns the contract object
func (_a *Challenge) Contract() *contract.Contract {
	return _a.c
}

// calls

// DisputeTree calls the disputeTree method in the solidity contract
func (_a *Challenge) DisputeTree(val0 *big.Int, block ...web3.BlockNumber) (retval0 *big.Int, retval1 web3.Address, retval2 *big.Int, retval3 [32]byte, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("disputeTree", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["parent"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}
	if err = mapstructure.Decode(out["challenger"], &retval1); err != nil {
		err = fmt.Errorf("failed to encode output at index 1")
	}
	if err = mapstructure.Decode(out["expireAfterBlock"], &retval2); err != nil {
		err = fmt.Errorf("failed to encode output at index 2")
	}
	if err = mapstructure.Decode(out["midStateRoot"], &retval3); err != nil {
		err = fmt.Errorf("failed to encode output at index 3")
	}

	return
}

// Factory calls the factory method in the solidity contract
func (_a *Challenge) Factory(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("factory", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// LastSelectedNodeKey calls the lastSelectedNodeKey method in the solidity contract
func (_a *Challenge) LastSelectedNodeKey(val0 web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("lastSelectedNodeKey", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// MinChallengerDeposit calls the minChallengerDeposit method in the solidity contract
func (_a *Challenge) MinChallengerDeposit(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("minChallengerDeposit", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Stage calls the stage method in the solidity contract
func (_a *Challenge) Stage(block ...web3.BlockNumber) (retval0 uint8, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("stage", web3.EncodeBlock(block...))
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

// ClaimChallengerWin sends a claimChallengerWin transaction in the solidity contract
func (_a *Challenge) ClaimChallengerWin(challenger web3.Address, stateInfo StateInfo) *contract.Txn {
	return _a.c.Txn("claimChallengerWin", challenger, stateInfo)
}

// ClaimProposerWin sends a claimProposerWin transaction in the solidity contract
func (_a *Challenge) ClaimProposerWin() *contract.Txn {
	return _a.c.Txn("claimProposerWin")
}

// Create sends a create transaction in the solidity contract
func (_a *Challenge) Create(blockN *big.Int, systemStartState [32]byte, creator web3.Address, proposerTimeLimit *big.Int, stateInfo StateInfo) *contract.Txn {
	return _a.c.Txn("create", blockN, systemStartState, creator, proposerTimeLimit, stateInfo)
}

// ExecOneStepTransition sends a execOneStepTransition transaction in the solidity contract
func (_a *Challenge) ExecOneStepTransition(leafNodeKey *big.Int) *contract.Txn {
	return _a.c.Txn("execOneStepTransition", leafNodeKey)
}

// Initialize sends a initialize transaction in the solidity contract
func (_a *Challenge) Initialize(endStep uint64, systemEndState [32]byte, midSystemState [32]byte) *contract.Txn {
	return _a.c.Txn("initialize", endStep, systemEndState, midSystemState)
}

// ProposerTimeout sends a proposerTimeout transaction in the solidity contract
func (_a *Challenge) ProposerTimeout(nodeKey *big.Int) *contract.Txn {
	return _a.c.Txn("proposerTimeout", nodeKey)
}

// RevealMidStates sends a revealMidStates transaction in the solidity contract
func (_a *Challenge) RevealMidStates(nodeKeys []*big.Int, stateRoots [][32]byte) *contract.Txn {
	return _a.c.Txn("revealMidStates", nodeKeys, stateRoots)
}

// SelectDisputeBranch sends a selectDisputeBranch transaction in the solidity contract
func (_a *Challenge) SelectDisputeBranch(parentNodeKeys []*big.Int, isLefts []bool) *contract.Txn {
	return _a.c.Txn("selectDisputeBranch", parentNodeKeys, isLefts)
}

// events

var ChallengeInitializedEventID = crypto.Keccak256Hash([]byte("ChallengeInitialized(uint128,bytes32)"))

func (_a *Challenge) ChallengeInitializedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{ChallengeInitializedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Challenge) FilterChallengeInitializedEvent(startBlock uint64, endBlock ...uint64) ([]*ChallengeInitializedEvent, error) {
	topic := _a.ChallengeInitializedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ChallengeInitializedEvent, 0)
	evts := _a.c.Abi.Events["ChallengeInitialized"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ChallengeInitializedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var ChallengeStartedEventID = crypto.Keccak256Hash([]byte("ChallengeStarted(uint256,address,bytes32,bytes32,uint256)"))

func (_a *Challenge) ChallengeStartedTopicFilter(l2BlockN []*big.Int, proposer []web3.Address) [][]web3.Hash {

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

func (_a *Challenge) FilterChallengeStartedEvent(l2BlockN []*big.Int, proposer []web3.Address, startBlock uint64, endBlock ...uint64) ([]*ChallengeStartedEvent, error) {
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

var DisputeBranchSelectedEventID = crypto.Keccak256Hash([]byte("DisputeBranchSelected(address,uint256[],uint256)"))

func (_a *Challenge) DisputeBranchSelectedTopicFilter(challenger []web3.Address) [][]web3.Hash {

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{DisputeBranchSelectedEventID}, challengerRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Challenge) FilterDisputeBranchSelectedEvent(challenger []web3.Address, startBlock uint64, endBlock ...uint64) ([]*DisputeBranchSelectedEvent, error) {
	topic := _a.DisputeBranchSelectedTopicFilter(challenger)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*DisputeBranchSelectedEvent, 0)
	evts := _a.c.Abi.Events["DisputeBranchSelected"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem DisputeBranchSelectedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var MidStateRevealedEventID = crypto.Keccak256Hash([]byte("MidStateRevealed(uint256[],bytes32[])"))

func (_a *Challenge) MidStateRevealedTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{MidStateRevealedEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Challenge) FilterMidStateRevealedEvent(startBlock uint64, endBlock ...uint64) ([]*MidStateRevealedEvent, error) {
	topic := _a.MidStateRevealedTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*MidStateRevealedEvent, 0)
	evts := _a.c.Abi.Events["MidStateRevealed"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem MidStateRevealedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var OneStepTransitionEventID = crypto.Keccak256Hash([]byte("OneStepTransition(uint256,bytes32,bytes32)"))

func (_a *Challenge) OneStepTransitionTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{OneStepTransitionEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Challenge) FilterOneStepTransitionEvent(startBlock uint64, endBlock ...uint64) ([]*OneStepTransitionEvent, error) {
	topic := _a.OneStepTransitionTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*OneStepTransitionEvent, 0)
	evts := _a.c.Abi.Events["OneStepTransition"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem OneStepTransitionEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var ProposerTimeoutEventID = crypto.Keccak256Hash([]byte("ProposerTimeout(uint256)"))

func (_a *Challenge) ProposerTimeoutTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{ProposerTimeoutEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Challenge) FilterProposerTimeoutEvent(startBlock uint64, endBlock ...uint64) ([]*ProposerTimeoutEvent, error) {
	topic := _a.ProposerTimeoutTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ProposerTimeoutEvent, 0)
	evts := _a.c.Abi.Events["ProposerTimeout"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ProposerTimeoutEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var ProposerWinEventID = crypto.Keccak256Hash([]byte("ProposerWin(address,uint256)"))

func (_a *Challenge) ProposerWinTopicFilter() [][]web3.Hash {

	var query [][]interface{}
	query = append(query, []interface{}{ProposerWinEventID})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Challenge) FilterProposerWinEvent(startBlock uint64, endBlock ...uint64) ([]*ProposerWinEvent, error) {
	topic := _a.ProposerWinTopicFilter()

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ProposerWinEvent, 0)
	evts := _a.c.Abi.Events["ProposerWin"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ProposerWinEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
