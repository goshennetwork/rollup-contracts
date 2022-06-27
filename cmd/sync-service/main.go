package main

import (
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

	var cfg config.SyncConfig
	utils.Ensure(utils.LoadJsonFile(config.DefaultSyncConfigName, &cfg))
	db, err := leveldbstore.NewLevelDBStore(cfg.DbDir)
	utils.Ensure(err)
	client, err := jsonrpc.NewClient(cfg.L1RpcUrl)
	utils.Ensure(err)
	syncService := sync_service.NewSyncService(db, client, &cfg)
	syncService.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	log.Info("shuting down!!!")
	syncService.Stop()
}
