package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/laizy/log"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/config"
	"github.com/ontology-layer-2/rollup-contracts/store/leveldbstore"
	sync_service "github.com/ontology-layer-2/rollup-contracts/sync-service"
	utils2 "github.com/ontology-layer-2/rollup-contracts/utils"
)

func main() {
	utils2.InitLog("./rollup-sync.log")
	var dbDir = flag.String("dbDir", config.DefaultSyncDbName, "set sync db dir")
	var deployOnL1Height = flag.Uint64("deployOnL1Height", 0, "set l1 rollup contracts deploy on height for speed up")
	var minConfirmNum = flag.Uint64("minConfirmNum", 6, "set min l1 confirm block num to avoid l1 chain reorg")
	flag.Parse()
	var cfg config.RollupCliConfig
	utils.Ensure(utils.LoadJsonFile(config.DefaultRollupConfigName, &cfg))
	db, err := leveldbstore.NewLevelDBStore(*dbDir)
	utils.Ensure(err)
	client, err := jsonrpc.NewClient(cfg.L1Rpc)
	utils.Ensure(err)
	syncService := sync_service.NewSyncService(db, client, &cfg, *dbDir, *deployOnL1Height, *minConfirmNum)
	syncService.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	log.Info("shuting down!!!")
	syncService.Stop()
}
