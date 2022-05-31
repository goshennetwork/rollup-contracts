package main

import (
	"os"

	"github.com/ontio/ontology/common/log"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/deploy"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/gateway"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/messaging"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/staking"
	"github.com/urfave/cli/v2"
)

func main() {
	log.Init(os.Stdout, "./Log/")
	app := &cli.App{
		Name:  "rullup",
		Usage: "rullup cli tool",
		Commands: []*cli.Command{
			gateway.GatewayCommand(),
			deploy.DeployCmd(),
			messaging.L1CrossLayerWitnessCommand(),
			staking.StakingCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
