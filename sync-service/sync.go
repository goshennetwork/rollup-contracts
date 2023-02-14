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

var (
	ErrNoBlock    = errors.New("no block")
	ErrNoTx       = errors.New("no transaction")
	ErrShortQueue = errors.New("shortage of queue")
	ErrInputHash  = errors.New("inconsistent input hash")
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
func RollBack(writer *store.StorageWriter, f byte) uint64 {
	switch f {
	case 1:
		info1 := writer.GetHighestL1CheckPointInfo1()
		if info1 == nil { //wired, panic
			panic(1)
		}
		for i := 0; i < len(info1.DirtyKey); i++ {
			reverse := len(info1.DirtyKey) - 1 - i
			writer.Cover(info1.DirtyKey[reverse], info1.DirtyValue[reverse])
		}
		return info1.StartPoint
	case 2:
		info2 := writer.GetHighestL1CheckPointInfo2()
		if info2 == nil {
			panic(1)
		}
		for i := 0; i < len(info2.DirtyKey); i++ {
			reverse := len(info2.DirtyKey) - 1 - i
			writer.Cover(info2.DirtyKey[reverse], info2.DirtyValue[reverse])
		}

		return info2.StartPoint

	case 3:
		info3 := writer.GetHighestL1CheckPointInfo3()

		if info3 == nil {
			panic(1)
		}
		for i := 0; i < len(info3.DirtyKey); i++ {
			reverse := len(info3.DirtyKey) - 1 - i
			writer.Cover(info3.DirtyKey[reverse], info3.DirtyValue[reverse])
		}
		return info3.StartPoint
	default:
		panic(1)
	}
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
	timer := time.NewTimer(0)
	defer timer.Stop()
	isSetup := lastHeight == 0
	round := 0
	startHeight := lastHeight + 1
	errSpan := 10 * time.Second
	for {
		select {
		case <-self.quit:
			return nil
		case <-timer.C:
		}
		if startHeight < self.conf.DeployOnL1Height { //speedup
			startHeight = self.conf.DeployOnL1Height
		}
		l1Height, err := self.l1client.Eth().BlockNumber()
		if err != nil {
			log.Warnf("l1 get block number error: %s", err)
			timer.Reset(errSpan)
			continue
		}
		if isSetup && startHeight+self.conf.MinConfirmBlockNum+2 > l1Height { //only setup period make sure first 2 block must confirmed
			log.Warn("l1 block too low,waiting..")
			timer.Reset(errSpan)
			continue
		}
		endHeight, err := CalcEndBlock(startHeight, l1Height)
		if err != nil {
			log.Warnf("l1 sync service: %s", err)
			timer.Reset(errSpan)
			continue
		}
		//be sure setup first 2 round will not roll back.
		if isSetup && round < 2 { //ez first 2 block
			round++
			endHeight = startHeight
		}
		rollback := func(flag byte) {
			if flag == 0 { //no need to rollback
				return
			}
			//which will make change history record
			self.dirtyLock.Lock()

			writer := self.db.Writer()
			startHeight = RollBack(writer, flag)
			lastEnd := startHeight - 1
			b, err := self.l1client.Eth().GetBlockByNumber(web3.BlockNumber(lastEnd), false)
			if err != nil || b == nil {
				if err == nil {
					err = ErrNoBlock
				}
				//unlock
				self.dirtyLock.Unlock()

				log.Warnf("l1 network err: %s", err)
				return
			}
			writer.SetLastSyncedL1Height(lastEnd)
			writer.SetLastSyncedL1Timestamp(b.Timestamp)
			writer.SetLastSyncedL1Hash(b.Hash)
			writer.SetL1DbVersion(writer.GetL1DbVersion() + 1)
			writer.Commit()
			self.dirtyLock.Unlock()
			log.Info("roll back")
			return
		}
		confirmedEndHeight := l1Height
		if confirmedEndHeight > self.conf.MinConfirmBlockNum {
			confirmedEndHeight -= self.conf.MinConfirmBlockNum
		}
		confirmedLastHeight := self.db.GetConfirmedLastSyncedL1Height()
		confirmedStartHeight := confirmedLastHeight + 1
		if confirmedStartHeight < self.conf.DeployOnL1Height {
			confirmedStartHeight = self.conf.DeployOnL1Height
		}
		confirmedEndHeight, err = CalcEndBlock(confirmedStartHeight, confirmedEndHeight)
		if err == nil { //ignore confirmed calc block
			if err := self.syncConfirmSyncL1Contracts(confirmedStartHeight, confirmedEndHeight); err != nil { //only network error
				log.Warnf("sync confirmed l1 contract failed: %s", err) //confirmed sync do not effect unsafe sync
			}
		}

		if startHeight+self.conf.MinConfirmBlockNum > endHeight {
			startHeight = endHeight - self.conf.MinConfirmBlockNum
		}
		if flag, err := self.syncL1Contracts(startHeight, endHeight); err != nil {
			//wired situation happened ,try to rollback
			log.Warnf("l1 sync error: %s,trying to rollback", err)
			rollback(flag)
			timer.Reset(errSpan)
			continue
		}
		startHeight = endHeight + 1
		isSetup = false
		log.Debugf("l1 sync to :%d", endHeight)
		timer.Reset(0)
	}
}

//duplicated sync is fine

func (self *SyncService) syncConfirmSyncL1Contracts(startHeight, endHeight uint64) error {
	overlay := self.db.Writer()
	if err := self.syncRollupStateChain(overlay, startHeight, endHeight); err != nil {
		return fmt.Errorf("sync rollup state chain: %w", err)
	}
	if err := self.syncL1Bridge(overlay, startHeight, endHeight); err != nil {
		return fmt.Errorf("sync l1 bridge: %w", err)
	}
	overlay.StoreConfirmedLastSyncedL1Height(endHeight)
	overlay.Commit()
	return nil
}

func (self *SyncService) syncL1Contracts(startHeight, endHeight uint64) (byte, error) {
	block, err := self.l1client.Eth().GetBlockByNumber(web3.BlockNumber(endHeight), false) // get block first
	if err != nil || block == nil {
		if err == nil {
			err = ErrNoBlock
		}
		return 0, err
	}
	queueStore, inputBatchStore, crossLayerStore := self.db.Writer(), self.db.Writer(), self.db.Writer()
	err, dirtyQueue := self.SyncRollupInputQueues(queueStore, startHeight, endHeight)
	if err != nil {
		return 1, fmt.Errorf("sync rollup input queue: %w", err)
	}
	err, dirtyInputBatch := self.syncRollupInputChainBatches(inputBatchStore, queueStore, startHeight, endHeight)
	if err != nil {
		if errors.Is(err, ErrShortQueue) || errors.Is(err, ErrInputHash) {
			return 1, fmt.Errorf("sync rollup input chain: %w", err)
		}
		return 2, fmt.Errorf("sync rollup input chain: %w", err)
	}
	err, dirtyCrossLayer := self.syncL1Witness(crossLayerStore, startHeight, endHeight)
	if err != nil {
		return 3, fmt.Errorf("sync l1 witness: %w", err)
	}
	if dirtyQueue {
		queueStore.StoreHighestL1CheckPointInfo1(startHeight)
	}
	if dirtyInputBatch {
		inputBatchStore.StoreHighestL1CheckPointInfo2(startHeight)
	}
	if dirtyCrossLayer {
		crossLayerStore.StoreHighestL1CheckPointInfo3(startHeight)
	}
	queueStore.Commit()
	crossLayerStore.Commit()
	inputBatchStore.SetLastSyncedL1Timestamp(block.Timestamp)
	inputBatchStore.SetLastSyncedL1Height(endHeight)
	inputBatchStore.SetLastSyncedL1Hash(block.Hash)
	inputBatchStore.Commit()
	return 0, nil
}

func (self *SyncService) SyncRollupInputQueues(kvdb *store.StorageWriter, startHeight, endHeight uint64) (error, bool) {
	rollupInputContract := binding.NewRollupInputChain(self.conf.L1Addresses.RollupInputChain, self.l1client)
	queues, err := rollupInputContract.FilterTransactionEnqueuedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return err, false
	}
	if len(queues) > 0 {
		//r1cs debug
		lastQueue := queues[len(queues)-1]
		log.Info("r1cs debug transaction enqueued event", "startQueueIndex", queues[0].QueueIndex, "lastQueueIndex", lastQueue.QueueIndex, "lastBlockHash", lastQueue.Raw.BlockHash, "lastBlockNumber", lastQueue.Raw.BlockNumber)
	}
	inputStore := kvdb.InputChain()
	if err := inputStore.StoreEnqueuedTransaction(queues...); err != nil {
		return err, false
	}
	return nil, len(queues) > 0
}

func (self *SyncService) syncRollupInputChainBatches(kvdb *store.StorageWriter, queueStore *store.StorageWriter, startHeight, endHeight uint64) (error, bool) {
	rollupInputContract := binding.NewRollupInputChain(self.conf.L1Addresses.RollupInputChain, self.l1client)

	batches, err := rollupInputContract.FilterInputBatchAppendedEvent(nil, nil, startHeight, endHeight)
	if err != nil {
		log.Errorf("sync fetch sequenced batch err: %s", err)
		return err, false
	}
	inputStore := kvdb.InputChain()
	if len(batches) > 0 { //r1cs debug
		last := batches[len(batches)-1]
		log.Info("r1cs debug input batches event", "firstBatchIndex", batches[0].Index, "lastBatchIndex", last.Index, "lastBlockHash", last.Raw.BlockHash, "lastBlockNumber", last.Raw.BlockNumber, "startBlockNumber", startHeight)
	}
	txs := make([]*web3.Transaction, 0)
	txBatchIndexes := make([]uint64, 0)
	for _, batch := range batches {
		// get transaction
		tx, err := self.l1client.Eth().GetTransactionByHash(batch.Raw.TransactionHash)
		if err != nil || tx == nil {
			if err == nil {
				err = ErrNoTx
			}
			return err, false
		}
		txs = append(txs, tx)
		txBatchIndexes = append(txBatchIndexes, batch.Index)
	}
	queueSize := queueStore.InputChain().QueueSize()
	for _, batch := range batches {
		if batch.StartQueueIndex+batch.QueueNum > queueSize {
			return ErrShortQueue, false
		}
	}
	if err := inputStore.StoreSequencerBatches(queueSize, batches...); err != nil {
		return err, false
	}
	inputStore.StoreSequencerBatchData(txs, txBatchIndexes)
	info := inputStore.GetInfo()
	log.Infof("queueTotalSize: %d, inputChain totalSize: %d", queueSize, info.TotalBatches)
	//now check
	for _, batch := range batches {
		batchData, err := inputStore.GetSequencerBatchData(batch.Index)
		if err != nil {
			return err, false
		}
		b := &binding.RollupInputBatches{}
		if err := b.Decode(batchData); err != nil {
			return fmt.Errorf("decode input batches failed, err: %s", err), false
		}
		queueHash := schema.CalcQueueHash(nil)
		if b.QueueNum > 0 {
			queues, err := queueStore.InputChain().GetEnqueuedTransactions(b.QueueStart, b.QueueNum)
			if err != nil {
				return err, false
			}
			queueHash = schema.CalcQueueHash(queues)
		}
		h := b.InputHash(queueHash)
		if h != batch.InputHash {
			log.Errorf("get wrong input, expected hash:%x, but %s, batchInfo: %s", batch.InputHash, h.String(), utils.JsonString(batch))
			return ErrInputHash, false
		}
	}

	return nil, len(batches) > 0
}

func (self *SyncService) syncL1Witness(kvdb *store.StorageWriter, startHeight, endHeight uint64) (error, bool) {
	l1Witness := binding.NewL1CrossLayerWitness(self.conf.L1Addresses.L1CrossLayerWitness, self.l1client)
	l1SentMsgs, err := l1Witness.FilterMessageSentEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("syncL1Witness: filter sent message, %s", err), false
	}
	l1BridgeStore := kvdb.L1CrossLayerWitness()
	beforeNum := l1BridgeStore.TotalMessage()
	if err := l1BridgeStore.StoreSentMessage(l1SentMsgs); err != nil {
		return fmt.Errorf("store sent message: %w", err), false
	}
	log.Infof("syncL1Witness: from %d to %d", startHeight, endHeight)
	return nil, l1BridgeStore.TotalMessage()-beforeNum > 0
}

func (self *SyncService) syncL1Bridge(kvdb *store.StorageWriter, startHeight, endHeight uint64) error {
	l1TokenBridge := binding.NewL1StandardBridge(self.conf.L1Addresses.L1StandardBridge, self.l1client)
	depositEvts, err := l1TokenBridge.FilterDepositInitiatedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("filter eth deposit, %s", err)
	}
	withdrawalEvts, err := l1TokenBridge.FilterWithdrawalFinalizedEvent(nil, nil, nil, startHeight, endHeight)
	if err != nil {
		return fmt.Errorf("filter eth withdrawal, %s", err)
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
		return fmt.Errorf("filter state batch appended event: %w", err)
	}
	if len(statesBatches) > 0 { // r1cs debug
		last := statesBatches[len(statesBatches)-1]
		log.Info("r1cs debug state batches", "lastBatchIndex", last.StartIndex, "lastBlockHash", last.Raw.BlockHash, "lastBlockNumber", last.Raw.BlockNumber)
	}
	stateStore := kvdb.StateChain()
	if err := stateStore.StoreBatchInfo(statesBatches...); err != nil {
		return fmt.Errorf("store batch info: %w", err)
	}
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
