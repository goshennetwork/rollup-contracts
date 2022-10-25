package sync_service

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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
	conf      *config.RollupCliConfig
	l1client  *jsonrpc.Client
	l2client  *jsonrpc.Client
	db        *store.Storage
	quit      chan struct{}
	wg        sync.WaitGroup
	dirtyLock *sync.Mutex
	running   uint32
}

func NewSyncService(diskdb schema.PersistStore, dirtyLock *sync.Mutex,
	l1client *jsonrpc.Client, l2client *jsonrpc.Client, cfg *config.RollupCliConfig) *SyncService {
	return &SyncService{
		db:        store.NewStorage(diskdb),
		conf:      cfg,
		l1client:  l1client,
		l2client:  l2client,
		quit:      make(chan struct{}),
		dirtyLock: dirtyLock,
	}
}

func (self *SyncService) isRunning() bool {
	return atomic.LoadUint32(&self.running) == 1
}

func (self *SyncService) Start() error {
	if !atomic.CompareAndSwapUint32(&self.running, 0, 1) {
		return errors.New("already running")
	}
	self.wg.Add(2)
	go func() {
		defer self.wg.Done()
		self.startL1Sync()
	}()
	go func() {
		defer self.wg.Done()
		self.startL2Sync()
	}()
	return nil
}

// RollPending  try rollback pending and highest info, and return the start height
func RollBack(writer *store.StorageWriter) uint64 {
	pendingInfo := writer.GetPendingL1CheckPointInfo()
	if pendingInfo == nil { //if there is no pending info, wired, just panic
		panic(1)
	}
	highestInfo := writer.GetHighestL1CheckPointInfo()
	if highestInfo == nil { //if there is no highest info, wired just panic
		panic(1)
	}
	for i := 0; i < len(pendingInfo.DirtyKey); i++ {
		reverse := len(pendingInfo.DirtyKey) - 1 - i
		writer.Cover(pendingInfo.DirtyKey[reverse], pendingInfo.DirtyValue[reverse])
	}
	for i := 0; i < len(highestInfo.DirtyKey); i++ {
		reverse := len(highestInfo.DirtyKey) - 1 - i
		writer.Cover(highestInfo.DirtyKey[reverse], highestInfo.DirtyValue[reverse])
	}
	//now clean confirm and pending info
	writer.SetHighestL1CheckPointInfo(&schema.L1CheckPointInfo{highestInfo.StartPoint, highestInfo.StartPoint, nil, nil})
	writer.SetPendingL1CheckPointInfo(&schema.L1CheckPointInfo{highestInfo.StartPoint, highestInfo.StartPoint, nil, nil})
	return highestInfo.StartPoint
}

func (self *SyncService) startL2Sync() error {
	lastHeight := self.db.GetLastSyncedL2Height()
	startHeight := lastHeight + 1
	for {
		select {
		case <-self.quit:
			return nil
		default:
		}
		l2Info, err := self.l2client.L2().GlobalInfo()
		if err != nil {
			log.Warnf("l2 get global info error: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		endHeight, err := CalcEndBlock(startHeight, uint64(l2Info.L2CheckedBlockNum)-1)
		if err != nil {
			log.Warnf("l2 sync service: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		err = self.syncL2Contracts(startHeight, endHeight)
		if err != nil {
			log.Warnf("l2 sync error: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		startHeight = endHeight + 1
		log.Debugf("l2 sync to :%d", endHeight)
	}
}

func (self *SyncService) startL1Sync() error {
	lastHeight := self.db.GetLastSyncedL1Height()
	isSetup := lastHeight == 0
	round := 0
	startHeight := lastHeight + 1
	for {
		select {
		case <-self.quit:
			return nil
		default:
		}
		if startHeight < self.conf.DeployOnL1Height { //speedup
			startHeight = self.conf.DeployOnL1Height
		}
		l1Height, err := self.l1client.Eth().BlockNumber()
		if err != nil {
			log.Warnf("l1 get block number error: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		if isSetup && startHeight+self.conf.MinConfirmBlockNum+2 > l1Height { //only setup period make sure first 2 block must confirmed
			log.Warn("l1 block too low,waiting..")
			continue
		}
		endHeight, err := CalcEndBlock(startHeight, l1Height)
		if err != nil {
			log.Warnf("l1 sync service: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		//be sure setup first 2 round will not roll back.
		if isSetup && round < 2 { //ez first 2 block
			round++
			endHeight = startHeight
		}
		//now check whether reorg first
		lastHash := self.db.GetLastSyncedL1Hash()
		b, err := self.l1client.Eth().GetBlockByNumber(web3.BlockNumber(startHeight-1), false)
		if err != nil {
			log.Warnf("l1 network err: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		if lastHash != (web3.Hash{}) && b.Hash != lastHash { //reorg happen, just rollback, make sure grap the del lock,
			//which will make change history record
			self.dirtyLock.Lock()

			writer := self.db.Writer()
			startHeight = RollBack(writer)
			lastEnd := startHeight - 1
			b, err := self.l1client.Eth().GetBlockByNumber(web3.BlockNumber(lastEnd), false)
			if err != nil {
				//unlock
				self.dirtyLock.Unlock()

				log.Warnf("l1 network err: %s", err)
				time.Sleep(15 * time.Second)
				continue
			}
			writer.SetLastSyncedL1Height(lastEnd)
			writer.SetLastSyncedL1Timestamp(b.Timestamp)
			writer.SetLastSyncedL1Hash(b.Hash)
			writer.SetL1DbVersion(writer.GetL1DbVersion() + 1)
			writer.Commit()
			self.dirtyLock.Unlock()

			log.Info("roll back")
			continue
		}

		err = self.syncL1Contracts(startHeight, endHeight)
		if err != nil {
			log.Warnf("l1 sync error: %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		startHeight = endHeight + 1
		isSetup = false
		log.Debugf("l1 sync to :%d", endHeight)
	}
}

func (self *SyncService) syncL1Contracts(startHeight, endHeight uint64) error {
	overlay := self.db.Writer()
	err := self.syncRollupInputChain(overlay, startHeight, endHeight)
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
	block, err := self.l1client.Eth().GetBlockByNumber(web3.BlockNumber(endHeight), false)
	if err != nil {
		return err
	}
	//now check point
	highestCheckpointInfo := overlay.GetHighestL1CheckPointInfo()
	dirtyk, dirtyv := overlay.Dirty()
	if highestCheckpointInfo == nil { //first, just record as highest check point info
		overlay.SetHighestL1CheckPointInfo(&schema.L1CheckPointInfo{startHeight, endHeight, dirtyk, dirtyv})
	} else {
		pendingCheckpoint := overlay.GetPendingL1CheckPointInfo()
		if pendingCheckpoint == nil {
			//open a new pend
			pendingCheckpoint = &schema.L1CheckPointInfo{startHeight, endHeight + 1, nil, nil}
		} else {
			//check consistence of pending key
			if pendingCheckpoint.EndPoint != startHeight { //wired should never happen
				//not consistence
				panic(1)
			}
		}
		pendingCheckpoint.DirtyKey = append(pendingCheckpoint.DirtyKey, dirtyk...)
		pendingCheckpoint.DirtyValue = append(pendingCheckpoint.DirtyValue, dirtyv...)
		pendingCheckpoint.EndPoint = endHeight + 1
		if pendingCheckpoint.OldEnough() { //reached, just make pending to highest
			overlay.SetHighestL1CheckPointInfo(pendingCheckpoint)
			//remove pending info, next pending start is end +1
			overlay.SetPendingL1CheckPointInfo(&schema.L1CheckPointInfo{endHeight + 1, endHeight + 1, nil, nil})
		} else { //not reach height, just add to pending
			overlay.SetPendingL1CheckPointInfo(pendingCheckpoint)
		}
	}
	overlay.SetLastSyncedL1Timestamp(block.Timestamp)
	overlay.SetLastSyncedL1Height(endHeight)
	overlay.SetLastSyncedL1Hash(block.Hash)
	overlay.Commit()
	return nil
}

func (self *SyncService) syncRollupInputChain(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	rollupInputContract := binding.NewRollupInputChain(self.conf.L1Addresses.RollupInputChain, self.l1client)
	queues, err := rollupInputContract.FilterTransactionEnqueuedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return err
	}
	batches, err := rollupInputContract.FilterInputBatchAppendedEvent(nil, nil, startHeight, endHeight)
	if err != nil {
		log.Errorf("sync fetch sequenced batch err:%s", err)
		return err
	}
	txs := make([]*web3.Transaction, 0)
	txBatchIndexes := make([]uint64, 0)
	for _, batch := range batches {
		// get transaction
		tx, err := self.l1client.Eth().GetTransactionByHash(batch.Raw.TransactionHash)
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
	l1Witness := binding.NewL1CrossLayerWitness(self.conf.L1Addresses.L1CrossLayerWitness, self.l1client)
	l1SentMsgs, err := l1Witness.FilterMessageSentEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Witness: filter sent message, %s", err)
	}
	l1BridgeStore := kvdb.L1CrossLayerWitness()
	l1BridgeStore.StoreSentMessage(l1SentMsgs)
	log.Infof("syncL1Witness: from %d to %d", startHeight, endHeight)
	return nil
}

func (self *SyncService) syncL1Bridge(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l1TokenBridge := binding.NewL1StandardBridge(self.conf.L1Addresses.L1StandardBridge, self.l1client)
	depositEvts, err := l1TokenBridge.FilterDepositInitiatedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Bridge: filter eth deposit, %s", err)
	}
	withdrawalEvts, err := l1TokenBridge.FilterWithdrawalFinalizedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Bridge: filter eth withdrawal, %s", err)
	}

	l1BridgeStore := kvdb.L1TokenBridge()
	l1BridgeStore.StoreDeposit(depositEvts)
	l1BridgeStore.StoreWithdrawal(withdrawalEvts)
	log.Infof("syncL1Bridge: from %d to %d", startHeight, endHeight)
	return nil
}

func (self *SyncService) Stop() error {
	close(self.quit)
	self.wg.Wait()
	return nil
}

func (self *SyncService) syncRollupStateChain(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	rollupStateContract := binding.NewRollupStateChain(self.conf.L1Addresses.RollupStateChain, self.l1client)
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
	addrMan := binding.NewAddressManager(self.conf.L1Addresses.AddressManager, self.l1client)
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

func (self *SyncService) syncL2Contracts(startHeight, endHeight uint64) error {
	writer := self.db.Writer()
	if err := self.syncL2Witness(writer, startHeight, endHeight); err != nil {
		return err
	}
	if err := self.syncL2Bridge(writer, startHeight, endHeight); err != nil {
		return err
	}
	writer.SetLastSyncedL2Height(endHeight)
	writer.Commit()
	return nil
}

func (self *SyncService) syncL2Witness(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l2Witness := binding.NewL2CrossLayerWitness(self.conf.L2Genesis.L2CrossLayerWitness, self.l2client)
	l2SentMsgs, err := l2Witness.FilterMessageSentEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL2Witness: filter sent message, %s", err)
	}
	l2WitnessStore := kvdb.L2CrossLayerWitness()
	l2WitnessStore.StoreSentMessage(l2SentMsgs)
	log.Infof("syncL2Witness: from %d to %d", startHeight, endHeight)
	return nil
}

func (self *SyncService) syncL2Bridge(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l2TokenBridge := binding.NewL2StandardBridge(self.conf.L1Addresses.L1StandardBridge, self.l2client)
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
	log.Infof("syncL2Bridge: from %d to %d", startHeight, endHeight)
	return nil
}

func CalcEndBlock(start, largest uint64) (uint64, error) {
	if largest < start {
		return 0, fmt.Errorf("beyond: start %d, largest %d", start, largest)
	}
	calc := start + 1024 // every 1024 range
	if (calc) < largest {
		return calc, nil
	} else {
		return largest, nil
	}
}
