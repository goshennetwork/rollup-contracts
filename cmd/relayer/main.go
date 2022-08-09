package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	utils2 "github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/ontology-layer-2/rollup-contracts/config"
	"github.com/ontology-layer-2/rollup-contracts/store"
	"github.com/ontology-layer-2/rollup-contracts/store/leveldbstore"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
	sync_service "github.com/ontology-layer-2/rollup-contracts/sync-service"
	"github.com/ontology-layer-2/rollup-contracts/utils"
	cli "github.com/urfave/cli/v2"
)

func main() {
	utils.InitLog("./relayer.log")
	app := &cli.App{
		Name:   "relayer",
		Usage:  "relay service for relay L2 -> L1 message",
		Action: run,
		Flags: []cli.Flag{
			flags.SubmitFlag,
			flags.ConfigFlag,
			flags.AccountsFlag,
			flags.GasLimit,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func run(ctx *cli.Context) error {
	relayerDBpath := "relayer.db"
	db, err := leveldbstore.NewLevelDBStore(relayerDBpath)
	utils2.Ensure(err)
	path := ctx.String(flags.ConfigFlag.Name)
	l1signer, cfg, err := common.SetUpL1(path)
	utils2.Ensure(err)
	l2signer, _, err := common.SetUpL2(path)
	utils2.Ensure(err)
	l1signer.Submit = ctx.Bool(flags.SubmitFlag.Name)
	gaslimit := ctx.Uint64(flags.GasLimit.Name)
	accountsStr := ctx.String(flags.AccountsFlag.Name)
	var accounts []web3.Address
	if len(accountsStr) > 0 {
		accountsStr = strings.Trim(accountsStr, " ")
		fmt.Println(accountsStr)
		for _, s := range strings.Split(accountsStr, ",") {
			accounts = append(accounts, web3.HexToAddress(s))
		}
	}
	_ = gaslimit
	relayService := NewRelayService(db, l1signer, l2signer.Client, cfg, accounts, gaslimit)
	relayService.Start()
	syncService := sync_service.NewSyncService(db, l1signer.Client, l2signer.Client, cfg)
	syncService.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	relayService.Stop()
	return syncService.Stop()
}

type RelayService struct {
	store      *store.Storage
	l1Signer   *contract.Signer
	l2Client   *jsonrpc.Client
	Whitelists []web3.Address
	cfg        *config.RollupCliConfig
	gasLimit   uint64
	quit       chan struct{}
}

func NewRelayService(db schema.PersistStore, l1signer *contract.Signer, l2client *jsonrpc.Client, cfg *config.RollupCliConfig, whitelist []web3.Address, gaslimit uint64) *RelayService {
	return &RelayService{
		store.NewStorage(db),
		l1signer,
		l2client,
		whitelist,
		cfg,
		gaslimit,
		make(chan struct{}),
	}
}

func (s *RelayService) Start() error {
	go s.start()
	return nil
}

func (s *RelayService) start() {
	ticker := time.NewTicker(500)
	pendingIndex := s.store.Relayer().GetPendingL1MsgIndex()
	var err error
	l1height := uint64(0)
	for {
		select {
		case <-s.quit:
			return
		case <-ticker.C:
			l1height, err = s.l1Signer.Client.Eth().BlockNumber()
			if err != nil {
				log.Errorf("get l2 blockNumber", "err", err)
				continue
			}
			syncedHeight := s.store.GetLastSyncedL1Height()
			if syncedHeight+6 < l1height {
				log.Warn("waiting sync finished", "l1Height", l1height, "syncedHeight", syncedHeight)
				continue
			}
			pendingIndex, err = s.TryRelay(pendingIndex)
			if err != nil {
				log.Error("try relay", "err", err)
			}
		}

	}
}

func (s *RelayService) TryRelay(msgIndex uint64) (uint64, error) {
	info, err := s.store.L2CrossLayerWitness().GetSentMessage(msgIndex)
	if err != nil {
		//have no L2 -> l1 message to relay
		return msgIndex, fmt.Errorf("no l2 -> l1 message to relay, msgIndex %d", msgIndex)
	}
	if s.ShouldRelay(info.Sender) {
		ErrAlreadyRelayed := fmt.Errorf("already relayed, msg index %d", msgIndex)
		_, err = s.store.L1CrossLayerWitness().GetRelayFailedMessage(msgIndex)
		if err != nil { //not failed, now check whether relayed
			_, err = s.store.L1CrossLayerWitness().GetRelayedMessage(msgIndex)
			if err != nil {
				//not relayed, try to relay
				if err := RelayMessage(s.l2Client, s.l1Signer, s.cfg.L1Addresses.L1CrossLayerWitness, msgIndex, s.gasLimit); err != nil {
					return msgIndex, err
				}
			} else { //already relayed once,
				err = ErrAlreadyRelayed
			}
		} else {
			//already relayed once,
			err = ErrAlreadyRelayed
		}
	} else { //no need to relay, just ignore
		log.Info("ignore relay message", "msgIndex", msgIndex)
	}

	pending := msgIndex + 1
	writer := s.store.Writer()
	writer.Relayer().StorePendingL1MsgIndex(pending)
	writer.Commit()
	return pending, err
}

// if whitelist is empty,default relay all message once
func (s *RelayService) ShouldRelay(sender web3.Address) bool {
	for _, addr := range s.Whitelists {
		if addr != sender {
			return false
		}
	}
	return true
}

func (s *RelayService) Stop() error {
	close(s.quit)
	return nil
}
func RelayMessage(l2client *jsonrpc.Client, l1signer *contract.Signer, l1Witness web3.Address, msgIndex, gasLimit uint64) error {
	params, err := l2client.L2().GetL1RelayMsgParams(msgIndex)
	if err != nil {
		return err
	}
	l1CrossLayerWitness := binding.NewL1CrossLayerWitness(l1Witness, l1signer.Client)
	l1CrossLayerWitness.Contract().SetFrom(l1signer.Address())
	StateInfo := binding.StateInfo{
		params.StateInfo.BlockHash,
		uint64(params.StateInfo.Index),
		uint64(params.StateInfo.Timestamp),
		params.StateInfo.Proposer,
	}
	proofs := make([][32]byte, len(params.Proof))
	for i, v := range params.Proof {
		proofs[i] = v
	}
	signedTx := l1CrossLayerWitness.RelayMessage(params.Target, params.Sender, params.Message, uint64(params.MessageIndex), params.RLPHeader, StateInfo, proofs).Sign(l1signer)
	if signedTx.Gas > gasLimit {
		return fmt.Errorf("out of gasLimit")
	}
	receipt := signedTx.SendTransaction(l1signer)
	if receipt.Status == 0 {
		log.Error("relay failed", "msgIndex", msgIndex)
		return nil
	}
	log.Info("relay message", "receipt", utils2.JsonStr(receipt))
	return nil
}
