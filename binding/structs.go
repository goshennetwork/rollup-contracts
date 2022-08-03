package binding

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = web3.HexToAddress
	_ = crypto.Keccak256Hash
)

var AddressSetEventID = crypto.Keccak256Hash([]byte("AddressSet(string,address,address)"))

type AddressSetEvent struct {
	Name string
	Old  web3.Address
	New  web3.Address

	Raw *web3.Log
}

var AdminChangedEventID = crypto.Keccak256Hash([]byte("AdminChanged(address,address)"))

type AdminChangedEvent struct {
	PreviousAdmin web3.Address
	NewAdmin      web3.Address

	Raw *web3.Log
}

var ApprovalEventID = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))

type ApprovalEvent struct {
	Owner   web3.Address
	Spender web3.Address
	Value   *big.Int

	Raw *web3.Log
}

var BeaconUpgradedEventID = crypto.Keccak256Hash([]byte("BeaconUpgraded(address)"))

type BeaconUpgradedEvent struct {
	Beacon web3.Address

	Raw *web3.Log
}

var ChallengeInitializedEventID = crypto.Keccak256Hash([]byte("ChallengeInitialized(uint128,bytes32)"))

type ChallengeInitializedEvent struct {
	SystemEndStep  *big.Int
	MidSystemState [32]byte

	Raw *web3.Log
}

var ChallengeStartedEventID = crypto.Keccak256Hash([]byte("ChallengeStarted(uint256,address,bytes32,uint256,address)"))

type ChallengeStartedEvent struct {
	L2BlockN         *big.Int
	Proposer         web3.Address
	StartSystemState [32]byte
	ExpireAfterBlock *big.Int
	Contract         web3.Address

	Raw *web3.Log
}

var ChallengerUpdatedEventID = crypto.Keccak256Hash([]byte("ChallengerUpdated(address,bool)"))

type ChallengerUpdatedEvent struct {
	Challenger web3.Address
	Enabled    bool

	Raw *web3.Log
}

var ChallengerWhitelistUpdatedEventID = crypto.Keccak256Hash([]byte("ChallengerWhitelistUpdated(address,bool)"))

type ChallengerWhitelistUpdatedEvent struct {
	Challenger web3.Address
	Enabled    bool

	Raw *web3.Log
}

var DepositClaimedEventID = crypto.Keccak256Hash([]byte("DepositClaimed(address,address,uint256)"))

type DepositClaimedEvent struct {
	Proposer web3.Address
	Receiver web3.Address
	Amount   *big.Int

	Raw *web3.Log
}

var DepositFailedEventID = crypto.Keccak256Hash([]byte("DepositFailed(address,address,address,address,uint256,bytes)"))

type DepositFailedEvent struct {
	L1Token web3.Address
	L2Token web3.Address
	From    web3.Address
	To      web3.Address
	Amount  *big.Int
	Data    []byte

	Raw *web3.Log
}

var DepositFinalizedEventID = crypto.Keccak256Hash([]byte("DepositFinalized(address,address,address,address,uint256,bytes)"))

type DepositFinalizedEvent struct {
	L1Token web3.Address
	L2Token web3.Address
	From    web3.Address
	To      web3.Address
	Amount  *big.Int
	Data    []byte

	Raw *web3.Log
}

var DepositInitiatedEventID = crypto.Keccak256Hash([]byte("DepositInitiated(address,address,address,address,uint256,bytes)"))

type DepositInitiatedEvent struct {
	L1Token web3.Address
	L2Token web3.Address
	From    web3.Address
	To      web3.Address
	Amount  *big.Int
	Data    []byte

	Raw *web3.Log
}

var DepositSlashedEventID = crypto.Keccak256Hash([]byte("DepositSlashed(address,address,uint256,bytes32)"))

type DepositSlashedEvent struct {
	Proposer    web3.Address
	Challenger  web3.Address
	BlockHeight *big.Int
	BlockHash   [32]byte

	Raw *web3.Log
}

var DepositedEventID = crypto.Keccak256Hash([]byte("Deposited(address,uint256)"))

type DepositedEvent struct {
	Proposer web3.Address
	Amount   *big.Int

	Raw *web3.Log
}

var DisputeBranchSelectedEventID = crypto.Keccak256Hash([]byte("DisputeBranchSelected(address,uint256[],uint256)"))

type DisputeBranchSelectedEvent struct {
	Challenger       web3.Address
	NodeKey          []*big.Int
	ExpireAfterBlock *big.Int

	Raw *web3.Log
}

var InitializedEventID = crypto.Keccak256Hash([]byte("Initialized(uint8)"))

type InitializedEvent struct {
	Version uint8

	Raw *web3.Log
}

var InputBatchAppendedEventID = crypto.Keccak256Hash([]byte("InputBatchAppended(address,uint64,uint64,uint64,bytes32)"))

type InputBatchAppendedEvent struct {
	Proposer        web3.Address
	Index           uint64
	StartQueueIndex uint64
	QueueNum        uint64
	InputHash       [32]byte

	Raw *web3.Log
}

var MessageAllowedEventID = crypto.Keccak256Hash([]byte("MessageAllowed(bytes32[])"))

type MessageAllowedEvent struct {
	MessageHashes [][32]byte

	Raw *web3.Log
}

var MessageBlockedEventID = crypto.Keccak256Hash([]byte("MessageBlocked(bytes32[])"))

type MessageBlockedEvent struct {
	MessageHashes [][32]byte

	Raw *web3.Log
}

var MessageRelayFailedEventID = crypto.Keccak256Hash([]byte("MessageRelayFailed(uint64,bytes32,uint64,bytes32)"))

type MessageRelayFailedEvent struct {
	MessageIndex uint64
	MsgHash      [32]byte
	MmrSize      uint64
	MmrRoot      [32]byte

	Raw *web3.Log
}

var MessageRelayedEventID = crypto.Keccak256Hash([]byte("MessageRelayed(uint64,bytes32)"))

type MessageRelayedEvent struct {
	MessageIndex uint64
	MsgHash      [32]byte

	Raw *web3.Log
}

var MessageSentEventID = crypto.Keccak256Hash([]byte("MessageSent(uint64,address,address,bytes32,bytes)"))

type MessageSentEvent struct {
	MessageIndex uint64
	Target       web3.Address
	Sender       web3.Address
	MmrRoot      [32]byte
	Message      []byte

	Raw *web3.Log
}

var MidStateRevealedEventID = crypto.Keccak256Hash([]byte("MidStateRevealed(uint256[],bytes32[])"))

type MidStateRevealedEvent struct {
	NodeKeys   []*big.Int
	StateRoots [][32]byte

	Raw *web3.Log
}

var OneStepTransitionEventID = crypto.Keccak256Hash([]byte("OneStepTransition(uint256,bytes32,bytes32)"))

type OneStepTransitionEvent struct {
	StartStep    *big.Int
	RevealedRoot [32]byte
	ExecutedRoot [32]byte

	Raw *web3.Log
}

var OwnershipTransferredEventID = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))

type OwnershipTransferredEvent struct {
	PreviousOwner web3.Address
	NewOwner      web3.Address

	Raw *web3.Log
}

var PausedEventID = crypto.Keccak256Hash([]byte("Paused(address)"))

type PausedEvent struct {
	Account web3.Address

	Raw *web3.Log
}

var ProposerTimeoutEventID = crypto.Keccak256Hash([]byte("ProposerTimeout(uint256)"))

type ProposerTimeoutEvent struct {
	NodeKey *big.Int

	Raw *web3.Log
}

var ProposerUpdatedEventID = crypto.Keccak256Hash([]byte("ProposerUpdated(address,bool)"))

type ProposerUpdatedEvent struct {
	Proposer web3.Address
	Enabled  bool

	Raw *web3.Log
}

var ProposerWhitelistUpdatedEventID = crypto.Keccak256Hash([]byte("ProposerWhitelistUpdated(address,bool)"))

type ProposerWhitelistUpdatedEvent struct {
	Proposer web3.Address
	Enabled  bool

	Raw *web3.Log
}

var ProposerWinEventID = crypto.Keccak256Hash([]byte("ProposerWin(address,uint256)"))

type ProposerWinEvent struct {
	Winner web3.Address
	Amount *big.Int

	Raw *web3.Log
}

var SequencerUpdatedEventID = crypto.Keccak256Hash([]byte("SequencerUpdated(address,bool)"))

type SequencerUpdatedEvent struct {
	Submitter web3.Address
	Enabled   bool

	Raw *web3.Log
}

var SequencerWhitelistUpdatedEventID = crypto.Keccak256Hash([]byte("SequencerWhitelistUpdated(address,bool)"))

type SequencerWhitelistUpdatedEvent struct {
	Submitter web3.Address
	Enabled   bool

	Raw *web3.Log
}

var StateBatchAppendedEventID = crypto.Keccak256Hash([]byte("StateBatchAppended(address,uint64,uint64,bytes32[])"))

type StateBatchAppendedEvent struct {
	Proposer   web3.Address
	StartIndex uint64
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

var StateRollbackedEventID = crypto.Keccak256Hash([]byte("StateRollbacked(uint64,bytes32)"))

type StateRollbackedEvent struct {
	StateIndex uint64
	BlockHash  [32]byte

	Raw *web3.Log
}

var TransactionEnqueuedEventID = crypto.Keccak256Hash([]byte("TransactionEnqueued(uint64,address,address,bytes,uint64)"))

type TransactionEnqueuedEvent struct {
	QueueIndex uint64
	From       web3.Address
	To         web3.Address
	RlpTx      []byte
	Timestamp  uint64

	Raw *web3.Log
}

var TransferEventID = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

type TransferEvent struct {
	From  web3.Address
	To    web3.Address
	Value *big.Int

	Raw *web3.Log
}

var UnpausedEventID = crypto.Keccak256Hash([]byte("Unpaused(address)"))

type UnpausedEvent struct {
	Account web3.Address

	Raw *web3.Log
}

var UpgradedEventID = crypto.Keccak256Hash([]byte("Upgraded(address)"))

type UpgradedEvent struct {
	Implementation web3.Address

	Raw *web3.Log
}

var WithdrawFinalizedEventID = crypto.Keccak256Hash([]byte("WithdrawFinalized(address,uint256)"))

type WithdrawFinalizedEvent struct {
	Proposer web3.Address
	Amount   *big.Int

	Raw *web3.Log
}

var WithdrawStartedEventID = crypto.Keccak256Hash([]byte("WithdrawStarted(address,uint256)"))

type WithdrawStartedEvent struct {
	Proposer           web3.Address
	NeedComfirmedBlock *big.Int

	Raw *web3.Log
}

var WithdrawalFinalizedEventID = crypto.Keccak256Hash([]byte("WithdrawalFinalized(address,address,address,address,uint256,bytes)"))

type WithdrawalFinalizedEvent struct {
	L1Token web3.Address
	L2Token web3.Address
	From    web3.Address
	To      web3.Address
	Amount  *big.Int
	Data    []byte

	Raw *web3.Log
}

var WithdrawalInitiatedEventID = crypto.Keccak256Hash([]byte("WithdrawalInitiated(address,address,address,address,uint256,bytes)"))

type WithdrawalInitiatedEvent struct {
	L1Token web3.Address
	L2Token web3.Address
	From    web3.Address
	To      web3.Address
	Amount  *big.Int
	Data    []byte

	Raw *web3.Log
}
