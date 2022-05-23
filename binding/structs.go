package binding

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = web3.HexToAddress
)

type AddressUpdatedEvent struct {
	Name string
	Old  web3.Address
	New  web3.Address

	Raw *web3.Log
}

type ApprovalEvent struct {
	Owner   web3.Address
	Spender web3.Address
	Value   *big.Int

	Raw *web3.Log
}

type ChallengeInitializedEvent struct {
	SystemEndStep  *big.Int
	MidSystemState [32]byte

	Raw *web3.Log
}

type ChallengeStartedEvent struct {
	L2BlockN         *big.Int
	Proposer         web3.Address
	StartSystemState [32]byte
	EndSystemState   [32]byte
	ExpireAfterBlock *big.Int

	Raw *web3.Log
}

type DepositClaimedEvent struct {
	Proposer web3.Address
	Receiver web3.Address
	Amount   *big.Int

	Raw *web3.Log
}

type DepositSlashedEvent struct {
	Proposer    web3.Address
	Challenger  web3.Address
	BlockHeight *big.Int
	BlockHash   [32]byte

	Raw *web3.Log
}

type DepositedEvent struct {
	Proposer web3.Address
	Amount   *big.Int

	Raw *web3.Log
}

type DisputeBranchSelectedEvent struct {
	Challenger       web3.Address
	NodeKey          []*big.Int
	ExpireAfterBlock *big.Int

	Raw *web3.Log
}

type MessageAllowedEvent struct {
	MessageHashes [][32]byte

	Raw *web3.Log
}

type MessageBlockedEvent struct {
	MessageHashes [][32]byte

	Raw *web3.Log
}

type MessageRelayFailedEvent struct {
	MsgHash [32]byte
	MmrSize uint64
	MmrRoot [32]byte

	Raw *web3.Log
}

type MessageRelayedEvent struct {
	MessageIndex uint64
	MsgHash      [32]byte

	Raw *web3.Log
}

type MessageSentEvent struct {
	MessageIndex uint64
	Target       web3.Address
	Sender       web3.Address
	Message      []byte

	Raw *web3.Log
}

type MidStateRevealedEvent struct {
	NodeKeys   []*big.Int
	StateRoots [][32]byte

	Raw *web3.Log
}

type OneStepTransitionEvent struct {
	StartStep    *big.Int
	RevealedRoot [32]byte
	ExecutedRoot [32]byte

	Raw *web3.Log
}

type OwnershipTransferredEvent struct {
	PreviousOwner web3.Address
	NewOwner      web3.Address

	Raw *web3.Log
}

type ProposerTimeoutEvent struct {
	NodeKey *big.Int

	Raw *web3.Log
}

type ProposerWinEvent struct {
	Winner web3.Address
	Amount *big.Int

	Raw *web3.Log
}

type StateBatchAppendedEvent struct {
	StartIndex uint64
	Proposer   web3.Address
	Timestamp  uint64
	BlockHash  [][32]byte

	Raw *web3.Log
}

type StateInfo struct {
	BlockHash [32]byte
	Index     uint64
	Timestamp uint64
	Proposer  web3.Address
}

type StateRollbackedEvent struct {
	StateIndex uint64
	BlockHash  [32]byte

	Raw *web3.Log
}

type TransactionAppendedEvent struct {
	Proposer        web3.Address
	StartQueueIndex *big.Int
	QueueNum        *big.Int
	ChainHeight     *big.Int
	InputHash       [32]byte

	Raw *web3.Log
}

type TransactionEnqueuedEvent struct {
	QueueIndex uint64
	From       web3.Address
	To         web3.Address
	Gaslimit   *big.Int
	Data       []byte
	Timestamp  uint64

	Raw *web3.Log
}

type TransferEvent struct {
	From  web3.Address
	To    web3.Address
	Value *big.Int

	Raw *web3.Log
}

type WithdrawFinalizedEvent struct {
	Proposer web3.Address
	Amount   *big.Int

	Raw *web3.Log
}

type WithdrawStartedEvent struct {
	Proposer           web3.Address
	NeedComfirmedBlock *big.Int

	Raw *web3.Log
}
