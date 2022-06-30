package schema

const (
	StateBatchPrefix        = 0x00
	RollupInputBatchKey     = 0x01 // batchIndex -> TransactionBatch
	SequencerQueuePrefix    = 0x02 // queueIndex -> QueueElement
	RollupInputBatchDataKey = 0x03 // batchIndex -> TransactionBatchData

	L1TokenBridgeETHDepositKey      = 0x04
	L1TokenBridgeERC20DepositKey    = 0x05
	L1TokenBridgeETHWithdrawalKey   = 0x06
	L1TokenBridgeERC20WithdrawalKey = 0x07

	L2TokenBridgeWithdrawalKey       = 0x08
	L2TokenBridgeDepositFinalizedKey = 0x09
	L2TokenBridgeDepositFailedKey    = 0x0A

	L1WitnessSentMessageKey = 0x0C // maybe duplicated with TransactionEnqueued
	L2WitnessSentMessageKey = 0x0D

	L2ClientExecutedQueuePrefix = 0x10 //l2 block number -> executed queueNum
	L2ClientCheckBlockNumPrefix = 0x11 //batch num -> checked l2 block num

	AddressNamePrefix = 0x20 // name -> address
)

var (
	LastSyncedL1HeightKey              = []byte{0x10} //last sync rollupInputContract's L1 block height
	LastSyncedL1TimestampKey           = []byte{0x11} //last sync l1 block timestamp
	CurrentRollupInputChainInfo        = []byte{0x12} // -> current rollupInputChain info
	RollupStateLastL1BlockHeightKey    = []byte{0x13} // last sync rollupStateContract's L1 block height
	CurrentRollupStateChainInfo        = []byte{0x14} // -> current rollupStateChain info
	AddressManagerLastL1BlockHeightKey = []byte{0x15}
	L1CompactMerkleTreeKey             = []byte{0x16}
	L2CompactMerkleTreeKey             = []byte{0x17}

	TotalUploadedBlock       = []byte{0x20} //->total uploaded l2 block size
	L2ClientCheckBatchNumKey = []byte{0x21} //-> checked batch num
	CurrentQueueBlockKey     = []byte{0x22} //-> head queue block
)