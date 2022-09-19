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

	PendingState            uint64
	PendingStateCreatedTime time.Time
	PendingInput            uint64
	PendingInputCreatedTime time.Time
}

func NewUploadService(l2client *jsonrpc.Client, l1client *jsonrpc.Client, signer *contract.Signer, stateChain *binding.RollupStateChain, inputChain *binding.RollupInputChain) *UploadBackend {

	return &UploadBackend{l2client, l1client, signer, stateChain, inputChain, make(chan struct{}), 0, time.Now(), 0, time.Now()}
}

func (self *UploadBackend) AppendInputBatch(batches *binding.RollupInputBatches) error {
	replace := false
	if batches.BatchIndex == self.PendingInput { //last input not finished yet, check time passed
		interval := time.Now().Sub(self.PendingInputCreatedTime)
		if interval >= time.Minute { //now timeout 1 minute, just replace last transaction
			replace = true
		} else { //wait pending tx
			return nil
		}
	}
	txn := self.inputChain.AppendInputBatches(batches)
	if replace { //use confirmed nonce
		nonce, err := self.l1client.Eth().GetNonce(self.signer.Address(), web3.Latest)
		if err != nil { //network tolerate
			return err
		}
		txn.SetNonce(nonce)
	}
	hs, err := self.l1client.Eth().SendRawTransaction(txn.Sign(self.signer).MarshalRLP())
	if err != nil {
		return err
	}
	self.PendingInputCreatedTime = time.Now()
	self.PendingInput = batches.BatchIndex
	log.Info("sending append inputBatch tx", "hash", hs, "batchIndex", batches.BatchIndex)
	return nil
}

func (self *UploadBackend) AppendStateBatch(blockHashes [][32]byte, startAt uint64) error {
	replace := false
	if startAt == self.PendingState { //last input not finished yet, check time passed
		interval := time.Now().Sub(self.PendingStateCreatedTime)
		if interval >= time.Minute { //now timeout 1 minute, just replace last transaction
			replace = true
		} else { //wait pending tx
			return nil
		}
	}
	txn := self.stateChain.AppendStateBatch(blockHashes, startAt)
	if replace { //use confirmed nonce
		nonce, err := self.l1client.Eth().GetNonce(self.signer.Address(), web3.Latest)
		if err != nil { //network tolerate
			return err
		}
		txn.SetNonce(nonce)
	}
	hs, err := self.l1client.Eth().SendRawTransaction(txn.Sign(self.signer).MarshalRLP())
	if err != nil {
		return err
	}
	self.PendingStateCreatedTime = time.Now()
	self.PendingState = startAt
	log.Info("sending append stateBatch tx", "hash", hs, "batchIndex", startAt)
	return nil
}

func (self *UploadBackend) Start() error {
	go self.runTxTask()
	go self.runStateTask()
	return nil
}

func (self *UploadBackend) runStateTask() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-self.quit:
			return
		case <-ticker.C:
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
				log.Debug("nothing to append", "l1 state batch num", l1StateNum, "l2 checked batch num", clientInfo.L2CheckedBatchNum)
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
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-self.quit:
			return
		case <-timer.C:
			timer.Reset(time.Duration(self.handle()))
		}
	}
}

func (self *UploadBackend) handle() (interval int64) {
	interval = int64(16 * time.Second)

	info, err := self.l2client.L2().GlobalInfo()
	if err != nil {
		log.Error("get l2 client info", "err", err)
		return
	}
	sss, _ := json.MarshalIndent(info, "", " ")
	log.Debug("global info", "info", string(sss))
	totalBatches, err := self.inputChain.ChainHeight()
	if err != nil {
		log.Error("l1 get input ChainHeight", "err", err)
		return
	}
	l2clientTotalBatches := uint64(info.L1InputInfo.TotalBatches)
	if l2clientTotalBatches < totalBatches {
		log.Warn("total batches not equal, waiting...", "l1 total input batch num", totalBatches, "l2 synced total batch num", info.L1InputInfo.TotalBatches)
		diff := int64(totalBatches - l2clientTotalBatches)
		utils.EnsureTrue(diff > 0)
		//every batch wait a block interval
		interval *= diff
		return
	}
	l2checkedBatchNum := uint64(info.L2CheckedBatchNum)
	if l2checkedBatchNum < totalBatches {
		log.Warn("l2 client have not checked all batches", "checkedBatchNum", info.L2CheckedBlockNum, "l1 total batch", info.L1InputInfo.TotalBatches)
		diff := int64(totalBatches - l2checkedBatchNum)
		utils.EnsureTrue(diff > 0)
		//every batch wait a block interval
		interval *= diff
		return
	}
	pendingQueueIndex, err := self.inputChain.PendingQueueIndex()
	if err != nil {
		log.Warn("l1 get pending queue index", "err", err)
		return
	}
	if uint64(info.L1InputInfo.PendingQueueIndex) != pendingQueueIndex {
		log.Warn("pending queue index not equal, waiting...", "l1 pendingQueueIndex", pendingQueueIndex, "l2 synced pendingQueueIndex", info.L1InputInfo.PendingQueueIndex)
		return
	}
	if info.L2HeadBlockNumber < info.L2CheckedBlockNum {
		log.Warn("no block to append", "l2 checked block num", info.L2CheckedBlockNum, "head block number", info.L2HeadBlockNumber)
		return
	}

	//may happen in situation of async
	if batchCode, err := self.l2client.L2().GetPendingTxBatches(); err != nil {
		log.Error(err.Error())
		return
	} else {
		if len(batchCode) == 0 {
			log.Warn("no batch code get")
			return
		}
		//try to decode fist
		b := new(binding.RollupInputBatches)
		if err := b.Decode(batchCode); err != nil {
			log.Error("decode batchCode", "err", err)
			return
		}
		if err := self.AppendInputBatch(b); err != nil {
			log.Error("append input batch failed", "batchIndex", b.BatchIndex, "err", err)
		}
	}
	//no err wait for block seal
	return
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
