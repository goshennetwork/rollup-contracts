package main

import (
	"os"

	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/genesis"

	"github.com/laizy/log"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/deploy"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/gateway"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/messaging"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/staking"
	"github.com/ontology-layer-2/rollup-contracts/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	utils.InitLog("./log/rollup.log")
	app := &cli.App{
		Name:  "rullup",
		Usage: "rullup cli tool",
		Commands: []*cli.Command{
			gateway.GatewayCommand(),
			deploy.DeployCmd(),
			messaging.L1CrossLayerWitnessCommand(),
			staking.StakingCommand(),
			genesis.GenesisCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
