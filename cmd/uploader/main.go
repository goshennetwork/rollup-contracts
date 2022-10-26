package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/config"
)

func main() {
	cfgName := flag.String("conf", "./rollup-config.json", "rollup config file name")
	submit := flag.Bool("submit", false, "whether submit tx to node")
	flag.Parse()
	var cfg config.RollupCliConfig
	utils.Ensure(utils.LoadJsonFile(*cfgName, &cfg))
	l1client, err := jsonrpc.NewClient(cfg.L1Rpc)
	if err != nil {
		panic(err)
	}
	chainId, err := l1client.Eth().ChainID()
	if err != nil {
		panic(err)
	}
	signer := contract.NewSigner(cfg.PrivKey, l1client, chainId.Uint64())
	signer.Submit = *submit

	stakingManager := binding.NewStakingManager(cfg.L1Addresses.StakingManager, l1client)
	stakingManager.Contract().SetFrom(signer.Address())
	whitelist := binding.NewWhitelist(cfg.L1Addresses.Whitelist, l1client)
	whitelist.Contract().SetFrom(signer.Address())
	err = checkPermission(stakingManager, whitelist, signer.Address())
	if err != nil {
		log.Error(err.Error())
		return
	}
	stateChain := binding.NewRollupStateChain(cfg.L1Addresses.RollupStateChain, l1client)
	inputChain := binding.NewRollupInputChain(cfg.L1Addresses.RollupInputChain, l1client)
	stateChain.Contract().SetFrom(signer.Address())
	inputChain.Contract().SetFrom(signer.Address())
	l2Client, err := jsonrpc.NewClient(cfg.L2Rpc)
	utils.Ensure(err)
	uploader := NewUploadService(l2Client, l1client, signer, stateChain, inputChain)
	uploader.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	uploader.Stop()
}

type UploadBackend struct {
	l2client   *jsonrpc.Client
	l1client   *jsonrpc.Client
	signer     *contract.Signer
	stateChain *binding.RollupStateChain
	inputChain *binding.RollupInputChain
	quit       chan struct{}
}

func NewUploadService(l2client *jsonrpc.Client, l1client *jsonrpc.Client, signer *contract.Signer, stateChain *binding.RollupStateChain, inputChain *binding.RollupInputChain) *UploadBackend {

	return &UploadBackend{l2client, l1client, signer, stateChain, inputChain, make(chan struct{})}
}

func (self *UploadBackend) AppendInputBatch(batches *binding.RollupInputBatches) error {
	txn := self.inputChain.AppendInputBatches(batches)
	//use confirmed nonce
	nonce, err := self.l1client.Eth().GetNonce(self.signer.Address(), web3.Latest)
	if err != nil { //network tolerate
		return err
	}
	tx := txn.SetNonce(nonce).Sign(self.signer)
	log.Infof("start sending transaction: %s, %s raw: %x\n", tx.Hash().String(), utils.JsonString(*tx.Transaction), tx.MarshalRLP())
	_, err = self.l1client.Eth().SendRawTransaction(tx.MarshalRLP())
	if err != nil {
		return err
	}
	log.Info("sending append inputBatch tx", "batchIndex", batches.BatchIndex)
	return nil
}

func (self *UploadBackend) AppendStateBatch(blockHashes [][32]byte, startAt uint64) error {
	txn := self.stateChain.AppendStateBatch(blockHashes, startAt)
	nonce, err := self.l1client.Eth().GetNonce(self.signer.Address(), web3.Latest)
	if err != nil { //network tolerate
		return err
	}
	tx := txn.SetNonce(nonce).Sign(self.signer)
	log.Infof("start sending transaction: %s, %s raw: %x\n", tx.Hash().String(), utils.JsonString(*tx.Transaction), tx.MarshalRLP())
	_, err = self.l1client.Eth().SendRawTransaction(tx.MarshalRLP())
	if err != nil {
		return err
	}
	log.Info("sending append stateBatch tx", "batchIndex", startAt)
	return nil
}

func (self *UploadBackend) Start() error {
	go self.runTxTask()
	go self.runStateTask()
	return nil
}

func (self *UploadBackend) runStateTask() {
	ticker := time.NewTicker(30 * time.Second)
	first := true
	defer ticker.Stop()

loop:
	for {
		select {
		case <-self.quit:
			return
		case <-ticker.C:
			if first {
				first = false
				ticker.Reset(time.Minute)
			}
			clientInfo, err := self.l2client.L2().GlobalInfo()
			if err != nil {
				log.Error("get global info", "err", err)
				continue
			}
			l1StateNum, err := self.stateChain.TotalSubmittedState()
			if err != nil {
				log.Error("get l1 total submitted state", "err", err)
				continue
			}
			l2checkedBatchNum := uint64(clientInfo.L2CheckedBatchNum)
			if l2checkedBatchNum <= l1StateNum {
				log.Debug("nothing to append", "l1 state batch num", l1StateNum, "l2 checked batch num", uint64(clientInfo.L2CheckedBatchNum))
				continue
			}
			num := l2checkedBatchNum - l1StateNum
			if num > 64 { // limit num
				num = 64
			}
			pendingStates := make([][32]byte, num)

			for i, _ := range pendingStates {
				index := l1StateNum + uint64(i)
				l2State, err := self.l2client.L2().GetRollupStateHash(index)
				if err != nil {
					log.Error("get state", "err", err)
					continue loop
				}
				if bytes.Equal(l2State.Bytes(), web3.Hash{}.Bytes()) {
					log.Warn("empty hash found", "batchIndex", index)
					continue loop
				}
				pendingStates[i] = l2State
			}
			log.Info("try to append state...", "start", l1StateNum, "end", l1StateNum+num)
			if err := self.AppendStateBatch(pendingStates, l1StateNum); err != nil {
				log.Error("append state batch failed", "batcIndex", l1StateNum, "err", err)
			}
		}
	}

}

//fixme: now only support one sequencer.
func (self *UploadBackend) runTxTask() {
	ticker := time.NewTicker(1)
	first := true
	defer ticker.Stop()
	for {
		select {
		case <-self.quit:
			return
		case <-ticker.C:
			if first {
				first = false
				ticker.Reset(time.Minute)
			}

			//may happen in situation of async
			if batchCode, err := self.l2client.L2().GetPendingTxBatches(); err != nil {
				log.Error(err.Error())
				continue
			} else {
				if len(batchCode) == 0 {
					log.Warn("no batch code get")
					continue
				}
				//try to decode fist
				b := new(binding.RollupInputBatches)
				if err := b.Decode(batchCode); err != nil {
					log.Error("decode batchCode", "err", err)
					continue
				}
				if err := self.AppendInputBatch(b); err != nil {
					log.Error("append input batch failed", "batchIndex", b.BatchIndex, "err", err)
				}
			}
		}
	}
}

func (self *UploadBackend) getPendingTxBatches() (*binding.RollupInputBatches, error) {
	info, err := self.l2client.L2().GlobalInfo()
	if err != nil {
		return nil, fmt.Errorf("get l2 client info: %w", err)
	}
	log.Debug("global info", "info", utils.JsonString(info))
	totalBatches, err := self.inputChain.ChainHeight()
	if err != nil {
		return nil, fmt.Errorf("l1 get input ChainHeight: %w", err)
	}
	l2clientTotalBatches := uint64(info.L1InputInfo.TotalBatches)
	if l2clientTotalBatches < totalBatches {
		return nil, fmt.Errorf("total batches not equal, waiting..., l1 total input batch num: %d, l2 synced total batch num: %d", totalBatches, uint64(info.L1InputInfo.TotalBatches))
	}
	l2checkedBatchNum := uint64(info.L2CheckedBatchNum)
	if l2checkedBatchNum < totalBatches {
		return nil, fmt.Errorf("l2 client have not checked all batches, checkedBatchNum: %d, l1 total batch: %d", uint64(info.L2CheckedBatchNum), uint64(info.L1InputInfo.TotalBatches))
	}
	pendingQueueIndex, err := self.inputChain.PendingQueueIndex()
	if err != nil {
		return nil, fmt.Errorf("l1 get pending queue index: %w", err)
	}
	if uint64(info.L1InputInfo.PendingQueueIndex) != pendingQueueIndex { //waiting l2 client catch up to newest state
		return nil, fmt.Errorf("pending queue index not equal, waiting..., l1 pendingQueueIndex: %d, l2 synced pendingQueueIndex: %d", pendingQueueIndex, uint64(info.L1InputInfo.PendingQueueIndex))
	}
	if info.L2HeadBlockNumber < info.L2CheckedBlockNum {
		return nil, fmt.Errorf("no block to append, l2 checked block num: %d, head block number: %d", uint64(info.L2CheckedBlockNum), uint64(info.L2HeadBlockNumber))
	}
	l2CheckedBlockNum := uint64(info.L2CheckedBlockNum)
	maxBlockes := uint64(info.L2HeadBlockNumber) - l2CheckedBlockNum + 1
	//todo: now simple limit upload size.should limit calldata size instead
	if maxBlockes > 512 {
		maxBlockes = 512
	}
	//batches := &binding.RollupInputBatches{
	//	QueueStart: uint64(info.L1InputInfo.PendingQueueIndex),
	//	BatchIndex: uint64(info.L2CheckedBatchNum),
	//}
	//var batchesData []byte
	//startQueueHeight, err := self.l2client.Eth().GetBlockByNumber(web3.BlockNumber(l2CheckedBlockNum-1), false)
	//for i := uint64(0); i < maxBlockes; i++ {
	//	blockNumber := i + l2CheckedBlockNum
	//	block := self.EthBackend.BlockChain().GetBlockByNumber(blockNumber)
	//	if block == nil { // should not happen except chain reorg
	//		log.Warn("nil block", "blockNumber", blockNumber)
	//		return nil
	//	}
	//	txs := block.Transactions()
	//	l2txs := FilterOrigin(txs)
	//	batches.QueueNum = block.Header().TotalExecutedQueueNum() - startQueueHeight
	//	if len(l2txs) > 0 {
	//		batches.SubBatches = append(batches.SubBatches, &binding.SubBatch{Timestamp: block.Time(), Txs: l2txs})
	//	}
	//	newBatch := batches.Encode()
	//	if len(newBatch)+4 < consts.MaxRollupInputBatchSize {
	//		batchesData = newBatch
	//	}
	//}
	//log.Info("generate batch", "index", batches.BatchIndex, "size", len(batchesData))
	//return batchesData

}

func (self *UploadBackend) Stop() error {
	close(self.quit)
	return nil
}

func checkPermission(stakingManager *binding.StakingManager, whitelist *binding.Whitelist, addr web3.Address) error {
	allowed, err := whitelist.CanSequence(addr)
	if err != nil {
		return err
	}
	if allowed == false {
		return fmt.Errorf("%s is not in sequencer whitelist", addr)
	}

	staked, err := stakingManager.IsStaking(addr, web3.Latest)
	if err != nil {
		return err
	}
	if staked == false {
		return fmt.Errorf("%s is not staked", addr)
	}

	return nil
}
