package main

import (
	"os"

	"github.com/laizy/log"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/cfg"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/dao"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/deploy"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/erc20"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/gateway"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/genesis"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/messaging"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/staking"
	"github.com/ontology-layer-2/rollup-contracts/utils"
	cli "github.com/urfave/cli/v2"
)

func main() {
	utils.InitLog("./rollup.log")
	app := &cli.App{
		Name:  "rullup",
		Usage: "rullup cli tool",
		Commands: []*cli.Command{
			gateway.GatewayCommand(),
			deploy.DeployCmd(),
			messaging.CrossLayerWitnessCommand(),
			staking.StakingCommand(),
			genesis.GenesisCommand(),
			erc20.ERC20Cmd(),
			dao.Cmd(),
			cfg.CfgCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
