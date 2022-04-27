package build

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

type OwnershipTransferredEvent struct {
	PreviousOwner web3.Address
	NewOwner      web3.Address

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
