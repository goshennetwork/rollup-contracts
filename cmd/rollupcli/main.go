package main

import (
	"os"

	"github.com/goshennetwork/rollup-contracts/cmd/rollupcli/deploy"
	"github.com/goshennetwork/rollup-contracts/cmd/rollupcli/erc20"
	"github.com/goshennetwork/rollup-contracts/cmd/rollupcli/gateway"
	"github.com/goshennetwork/rollup-contracts/cmd/rollupcli/genesis"
	"github.com/goshennetwork/rollup-contracts/cmd/rollupcli/messaging"
	"github.com/goshennetwork/rollup-contracts/cmd/rollupcli/staking"
	"github.com/goshennetwork/rollup-contracts/cmd/rollupcli/whitelist"
	"github.com/goshennetwork/rollup-contracts/utils"
	"github.com/laizy/log"
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
			whitelist.Cmd(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
