package sync_service

import (
	"fmt"
	"time"

	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/config"
	"github.com/ontology-layer-2/rollup-contracts/store"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type SyncService struct {
	conf   *config.SyncConfig
	client *jsonrpc.Client
	db     *store.Storage
	quit   chan struct{}
}

func NewSyncService(diskdb schema.PersistStore,
	client *jsonrpc.Client, cfg *config.SyncConfig) *SyncService {
	return &SyncService{
		db:     store.NewStorage(diskdb, cfg.DbDir),
		conf:   cfg,
		client: client,
		quit:   make(chan struct{}),
	}
}

func (self *SyncService) Start() error {
	go self.start()
	return nil
}

func (self *SyncService) start() error {
	lastHeight := self.db.GetLastSyncedL1Height()
	startHeight := lastHeight + 1
	for {
		select {
		case <-self.quit:
			return nil
		default:

		}
		if startHeight < self.conf.StartSyncHeight { //speedup
			startHeight = self.conf.StartSyncHeight
		}
		l1Height, err := self.client.Eth().BlockNumber()
		if err != nil {
			log.Warnf("get block number error: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		if l1Height > self.conf.MinConfirmBlockNum {
			l1Height -= self.conf.MinConfirmBlockNum
		}
		endHeight, err := CalcEndBlock(startHeight, l1Height)
		if err != nil {
			log.Warnf("Input sync service: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		err = self.syncL1Contracts(startHeight, endHeight)
		if err != nil {
			log.Warnf("sync error: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		startHeight = endHeight + 1
		log.Debugf("input sync to :%d", endHeight)
	}
}

func (self *SyncService) syncL1Contracts(startHeight, endHeight uint64) error {
	overlay := self.db.Writer()
	err := self.syncAddrManager(overlay, startHeight, endHeight)
	if err != nil {
		return err
	}
	err = self.syncRollupInputChain(overlay, startHeight, endHeight)
	if err != nil {
		return err
	}
	err = self.syncRollupStateChain(overlay, startHeight, endHeight)
	if err != nil {
		return err
	}
	err = self.syncL1Witness(overlay, startHeight, endHeight)
	if err != nil {
		return err
	}
	err = self.syncL1Bridge(overlay, startHeight, endHeight)
	if err != nil {
		return err
	}
	block, err := self.client.Eth().GetBlockByNumber(web3.BlockNumber(endHeight), false)
	if err != nil {
		return err
	}
	overlay.SetLastSyncedL1Timestamp(block.Timestamp)
	overlay.SetLastSyncedL1Height(endHeight)
	overlay.Commit()
	return nil
}

func (self *SyncService) syncRollupInputChain(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	rollupInputContract := binding.NewRollupInputChain(self.conf.RollupInputChain, self.client)
	queues, err := rollupInputContract.FilterTransactionEnqueuedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return err
	}
	batches, err := rollupInputContract.FilterTransactionAppendedEvent(nil, nil, startHeight, endHeight)
	if err != nil {
		log.Errorf("sync fetch sequenced batch err:%s", err)
		return err
	}
	txs := make([]*web3.Transaction, 0)
	txBatchIndexes := make([]uint64, 0)
	for _, batch := range batches {
		// get transaction
		tx, err := self.client.Eth().GetTransactionByHash(batch.Raw.TransactionHash)
		if err != nil {
			log.Errorf("sync fetch sequenced batch tx, %s", err)
			return err
		}
		txs = append(txs, tx)
		txBatchIndexes = append(txBatchIndexes, batch.Index)
	}
	inputStore := kvdb.InputChain()
	inputStore.StoreEnqueuedTransaction(queues...)
	inputStore.StoreSequencerBatches(batches...)
	inputStore.StoreSequencerBatchData(txs, txBatchIndexes)
	info := inputStore.GetInfo()
	log.Infof("queueTotalSize: %d, inputChain totalSize: %d", info.QueueSize, info.TotalBatches)
	//now check
	for _, batch := range batches {
		batchData, err := inputStore.GetSequencerBatchData(batch.Index)
		utils.Ensure(err)
		b := &binding.RollupInputBatches{}
		if err := b.Decode(batchData); err != nil {
			log.Errorf("decode input batches failed, err: %s", err)
			return err
		}
		queueHash := schema.CalcQueueHash(nil)
		if b.QueueNum > 0 {
			queues, err := inputStore.GetEnqueuedTransactions(b.QueueStart, b.QueueNum)
			if err != nil {
				return err
			}
			queueHash = schema.CalcQueueHash(queues)
		}
		h := b.InputHash(queueHash)
		if h != batch.InputHash {
			return fmt.Errorf("get wrong input, expected hash:%x, but %x", batch.InputHash, h)
		}
	}
	return nil
}

func (self *SyncService) syncL1Witness(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l1Witness := binding.NewL1CrossLayerWitness(self.conf.L1CrossLayerWitness, self.client)
	l1SentMsgs, err := l1Witness.FilterMessageSentEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Witness: filter sent message, %s", err)
	}
	l1BridgeStore := kvdb.L1CrossLayerWitness()
	l1BridgeStore.StoreSentMessage(l1SentMsgs)
	kvdb.StoreL1CompactMerkleTree()
	log.Infof("syncL1Witness: from %d to %d", startHeight, endHeight)
	return nil
}

func (self *SyncService) syncL1Bridge(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l1TokenBridge := binding.NewL1StandardBridge(self.conf.L1TokenBridge, self.client)
	ethDepositEvts, err := l1TokenBridge.FilterETHDepositInitiatedEvent(nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Bridge: filter eth deposit, %s", err)
	}
	ethWithdrawalEvts, err := l1TokenBridge.FilterETHWithdrawalFinalizedEvent(nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Bridge: filter eth withdrawal, %s", err)
	}
	erc20DepositEvts, err := l1TokenBridge.FilterERC20DepositInitiatedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Bridge: filter erc20 deposit, %s", err)
	}
	erc20WithdrawalEvts, err := l1TokenBridge.FilterERC20WithdrawalFinalizedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Bridge: filter erc20 withdrawal, %s", err)
	}
	l1BridgeStore := kvdb.L1TokenBridge()
	l1BridgeStore.StoreETHDeposit(ethDepositEvts)
	l1BridgeStore.StoreETHWithdrawal(ethWithdrawalEvts)
	l1BridgeStore.StoreERC20Deposit(erc20DepositEvts)
	l1BridgeStore.StoreERC20Withdrawal(erc20WithdrawalEvts)
	log.Infof("syncL1Bridge: from %d to %d", startHeight, endHeight)
	return nil
}

func (self *SyncService) Stop() error {
	self.quit <- struct{}{}
	return nil
}

func (self *SyncService) syncRollupStateChain(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	rollupStateContract := binding.NewRollupStateChain(self.conf.RollupStateChain, self.client)
	statesBatches, err := rollupStateContract.FilterStateBatchAppendedEvent(nil, nil, startHeight, endHeight)
	if err != nil {
		return err
	}
	stateStore := kvdb.StateChain()
	stateStore.StoreBatchInfo(statesBatches...)
	info := stateStore.GetInfo()
	log.Infof("total state chain size: %d", info.TotalSize)
	return nil
}

func (self *SyncService) syncAddrManager(writer *store.StorageWriter, startHeight, endHeight uint64) error {
	addrMan := binding.NewAddressManager(self.conf.AddressManager, self.client)
	updated, err := addrMan.FilterAddressSetEvent(startHeight, endHeight)
	if err != nil {
		return err
	}
	addrStore := writer.AddressManager()
	for _, v := range updated {
		addrStore.SetAddress(v.Name, v.New)
	}
	return nil
}

func (self *SyncService) syncL2Witness(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l2Witness := binding.NewL2CrossLayerWitness(self.conf.L2CrossLayerWitness, self.client)
	l2SentMsgs, err := l2Witness.FilterMessageSentEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL2Witness: filter sent message, %s", err)
	}
	l2WitnessStore := kvdb.L2CrossLayerWitness()
	l2WitnessStore.StoreSentMessage(l2SentMsgs)
	kvdb.StoreL2CompactMerkleTree()
	log.Infof("syncL2Witness: from %d to %d", startHeight, endHeight)
	return nil
}

func (self *SyncService) syncL2Bridge(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l2TokenBridge := binding.NewL2StandardBridge(self.conf.L2TokenBridge, self.client)
	tokenWithdrawalEvts, err := l2TokenBridge.FilterWithdrawalInitiatedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL2Bridge: filter eth withdrawal, %s", err)
	}
	tokenDepositEvts, err := l2TokenBridge.FilterDepositFinalizedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL2Bridge: filter erc20 deposit, %s", err)
	}
	tokenDepositFailedEvts, err := l2TokenBridge.FilterDepositFailedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL2Bridge: filter erc20 withdrawal, %s", err)
	}
	l2BridgeStore := kvdb.L2TokenBridge()
	l2BridgeStore.StoreWithdrawal(tokenWithdrawalEvts)
	l2BridgeStore.StoreDepositFinalized(tokenDepositEvts)
	l2BridgeStore.StoreDepositFailed(tokenDepositFailedEvts)
	log.Infof("syncL1Bridge: from %d to %d", startHeight, endHeight)
	return nil
}

func CalcEndBlock(start, largest uint64) (uint64, error) {
	if largest < start {
		return 0, fmt.Errorf("beyond: start %d, largest %d", start, largest)
	}
	calc := start + 1024
	if (calc) < largest {
		return calc, nil
	} else {
		return largest, nil
	}
}