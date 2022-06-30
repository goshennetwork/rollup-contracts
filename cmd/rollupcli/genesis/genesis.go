package genesis

import (
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
	utils2 "github.com/ontology-layer-2/rollup-contracts/utils"
	cli "github.com/urfave/cli/v2"
)

var OutFile = &cli.StringFlag{
	Name:  "out",
	Usage: "output file name",
	Value: "genesis-data.json",
}

func GenesisCommand() *cli.Command {
	return &cli.Command{
		Name:   "genesis",
		Usage:  "generate l2 contracts genesis data",
		Action: GenesisCmd,
		Flags: []cli.Flag{
			OutFile,
			flags.ConfigFlag,
		},
	}
}

func GenesisCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)
	rollupCfg, err := common.LoadConf(path)
	if err != nil {
		return err
	}

	outfile := ctx.String(OutFile.Name)

	result := deploy.BuildL2GenesisData(rollupCfg.L2Genesis, rollupCfg.L1Addresses.L1StandardBridge)
	utils2.AtomicWriteFile(outfile, utils.JsonString(result))

	return nil
}
