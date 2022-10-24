package schema

const (
	StateBatchPrefix        = 0x00
	RollupInputBatchKey     = 0x01 // batchIndex -> TransactionBatch
	SequencerQueuePrefix    = 0x02 // queueIndex -> QueueElement
	RollupInputBatchDataKey = 0x03 // batchIndex -> TransactionBatchData

	L1TokenBridgeDepositKey    = 0x04
	L1TokenBridgeWithdrawalKey = 0x05

	L2TokenBridgeWithdrawalKey       = 0x08
	L2TokenBridgeDepositFinalizedKey = 0x09
	L2TokenBridgeDepositFailedKey    = 0x0A

	L1WitnessSentMessageKey = 0x0C // maybe duplicated with TransactionEnqueued
	L2WitnessSentMessageKey = 0x0D

	L2ClientCheckBlockNumPrefix = 0x10 //batch index -> checked l2 block num
L2ClientProofPrefix         = 0x11 //batch index -> read-storage-proof	L1MMRDataPrefix             = 0x16
	L2MMRDataPrefix             = 0x17

	AddressNamePrefix = 0x20 // name -> address
)

var (
	LastSyncedL1HeightKey              = []byte{0x10} //last sync rollupInputContract's L1 block height
	LastSyncedL1TimestampKey           = []byte{0x11} //last sync l1 block timestamp
	CurrentRollupInputChainInfoKey     = []byte{0x12} // -> current rollupInputChain info
	RollupStateLastL1BlockHeightKey    = []byte{0x13} // last sync rollupStateContract's L1 block height
	CurrentRollupStateChainInfoKey     = []byte{0x14} // -> current rollupStateChain info
	AddressManagerLastL1BlockHeightKey = []byte{0x15}
	L1CompactMerkleTreeKey             = []byte{0x16}
	L2CompactMerkleTreeKey             = []byte{0x17}
	LastSyncedL2HeightKey              = []byte{0x18}
	LastSyncedL1Hash                   = []byte{0x19}
	HighestL1CheckPointInfoKey         = []byte{0x1a} //rollback to this l1 height point, this point is simply last LastSyncedL1Height
	PendingL1CheckPointInfoKey         = []byte{0x1b}
	L1DbVersionKey                     = []byte{0x1c}

	L2ClientCheckBatchNumKey = []byte{0x20} //-> checked batch num
	CurrentQueueBlockKey     = []byte{0x21} //-> head queue block
	L2ClientVersion          = []byte{0x22} //-> l2 client version
	L2ClientConfirmPoint     = []byte{0x23} // -> l2 client confirm point info
	L2ClientPendingPoint     = []byte{0x24} // -> l2 client pending point info
)
