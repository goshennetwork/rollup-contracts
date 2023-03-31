package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/goshennetwork/rollup-contracts/config"
	"github.com/goshennetwork/rollup-contracts/store/leveldbstore"
	sync_service "github.com/goshennetwork/rollup-contracts/sync-service"
	utils2 "github.com/goshennetwork/rollup-contracts/utils"
	"github.com/laizy/log"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
)

func main() {
	utils2.InitLog("./rollup-sync.log")
	var dbDir = flag.String("dbDir", config.DefaultSyncDbName, "set sync db name")
	flag.Parse()
	var cfg config.RollupCliConfig
	utils.Ensure(utils.LoadJsonFile(config.DefaultRollupConfigName, &cfg))
	db, err := leveldbstore.NewLevelDBStore(*dbDir)
	utils.Ensure(err)
	l1client, err := jsonrpc.NewClient(cfg.L1Rpc)
	utils.Ensure(err)
	l2client, err := jsonrpc.NewClient(cfg.L1Rpc)
	utils.Ensure(err)
	utils.Ensure(err)
	syncService := sync_service.NewSyncService(db, l1client, l2client, &cfg)
	syncService.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	log.Info("shuting down!!!")
	syncService.Stop()
}
