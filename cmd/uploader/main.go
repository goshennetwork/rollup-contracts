package main

import (
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
	var l2 = flag.String("l2", "http://localhost:23333", "set l2 client url")
	cfgName := flag.String("conf", "./rollup-config.json", "rollup config file name")
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
	signer.Submit = true

	stakingManager := binding.NewStakingManager(cfg.L1Addresses.StakingManager, l1client)
	stakingManager.Contract().SetFrom(signer.Address())
	dao := binding.NewDAO(cfg.L1Addresses.DAO, l1client)
	dao.Contract().SetFrom(signer.Address())
	err = checkPermission(stakingManager, dao, signer.Address())
	if err != nil {
		log.Error(err.Error())
		return
	}
	stateChain := binding.NewRollupStateChain(cfg.L1Addresses.RollupStateChain, l1client)
	inputChain := binding.NewRollupInputChain(cfg.L1Addresses.RollupInputChain, l1client)
	stateChain.Contract().SetFrom(signer.Address())
	inputChain.Contract().SetFrom(signer.Address())
	l2Client, err := jsonrpc.NewClient(*l2)
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

func (self *UploadBackend) Start() error {
	go self.runTxTask()
	return nil
}

//fixme: now only support one sequencer.
func (self *UploadBackend) runTxTask() {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for range timer.C {
		select {
		case <-self.quit:
			return
		default:
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
	if info.L1InputInfo.TotalBatches < totalBatches {
		log.Warn("total batches not equal, waiting...", "l1 total input batch num", totalBatches, "l2 synced total batch num", info.L1InputInfo.TotalBatches)
		diff := int64(totalBatches - info.L1InputInfo.TotalBatches)
		utils.EnsureTrue(diff > 0)
		//every batch wait a block interval
		interval *= diff
		return
	}
	if info.L2CheckedBatchNum < totalBatches {
		log.Warn("l2 client have not checked all batches", "checkedBatchNum", info.L2CheckedBlockNum, "l1 total batch", info.L1InputInfo.TotalBatches)
		diff := int64(totalBatches - info.L2CheckedBatchNum)
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
	if info.L1InputInfo.PendingQueueIndex != pendingQueueIndex {
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
		receipt := self.inputChain.AppendInputBatches(b).Sign(self.signer).SendTransaction(self.signer)
		if receipt.IsReverted() {
			log.Errorf("append input batch failed: %s", utils.JsonString(receipt))
			return
		}
	}
	//no err wait for block seal
	return
}

func (self *UploadBackend) Stop() error {
	self.quit <- struct{}{}
	return nil
}

func checkPermission(stakingManager *binding.StakingManager, dao *binding.DAO, addr web3.Address) error {
	allowed, err := dao.SequencerWhitelist(addr)
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
