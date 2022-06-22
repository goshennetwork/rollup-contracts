package main

import (
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
	flag.Parse()

	var cfg config.RollupCliConfig
	utils.Ensure(utils.LoadJsonFile(*cfgName, &cfg))
	client, err := jsonrpc.NewClient(cfg.L1Rpc)
	if err != nil {
		panic(err)
	}
	chainId, err := client.Eth().ChainID()
	if err != nil {
		panic(err)
	}
	signer := contract.NewSigner(cfg.PrivKey, client, chainId.Uint64())
	signer.Submit = true

	stakingManager := binding.NewStakingManager(cfg.L1Addresses.StakingManager, client)
	stakingManager.Contract().SetFrom(signer.Address())
	dao := binding.NewDAO(cfg.L1Addresses.DAO, client)
	dao.Contract().SetFrom(signer.Address())
	err = checkPermission(stakingManager, dao, signer.Address())
	if err != nil {
		log.Error(err.Error())
		return
	}
	stateChain := binding.NewRollupStateChain(cfg.L1Addresses.RollupStateChain, client)
	inputChain := binding.NewRollupInputChain(cfg.L1Addresses.RollupInputChain, client)
	stateChain.Contract().SetFrom(signer.Address())
	inputChain.Contract().SetFrom(signer.Address())
	l2Client, err := jsonrpc.NewClient(cfg.L2Rpc)
	utils.Ensure(err)
	uploader := NewUploadService(l2Client, signer, stateChain, inputChain)
	uploader.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	uploader.Stop()
}

type UploadBackend struct {
	l2client   *jsonrpc.Client
	signer     *contract.Signer
	stateChain *binding.RollupStateChain
	inputChain *binding.RollupInputChain
	quit       chan struct{}
}

func NewUploadService(l2client *jsonrpc.Client, signer *contract.Signer, stateChain *binding.RollupStateChain, inputChain *binding.RollupInputChain) *UploadBackend {
	return &UploadBackend{l2client, signer, stateChain, inputChain, make(chan struct{})}
}

func (self *UploadBackend) Start() error {
	go self.runTxTask()
	return nil
}

//fixme: now only support one sequencer.
func (self *UploadBackend) runTxTask() {
	interval := 10 * time.Second
	ticker := time.NewTimer(0)
	defer ticker.Stop()
	for range ticker.C {
		ticker.Reset(interval)
		select {
		case <-self.quit:
			return
		default:
		}
		batchCode, err := self.l2client.L2().GetPendingTxBatches()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		var batch binding.RollupInputBatches
		err = batch.Decode(batchCode)
		utils.Ensure(err)
		log.Infof("start upload input batch: %s", utils.JsonString(batch))
		receipt := self.inputChain.AppendInputBatches(batch.Encode()).Sign(self.signer).SendTransaction(self.signer)
		if receipt.IsReverted() {
			log.Errorf("append input batch failed: %s", utils.JsonString(receipt))
		}
	}
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
