package genesis

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/u256"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
	utils2 "github.com/ontology-layer-2/rollup-contracts/utils"
	"github.com/urfave/cli/v2"
)

var FeeCollectorOwner = &cli.StringFlag{
	Name:     "feeOwner",
	Usage:    "fee collector contract owner address",
	Required: true,
}

var FeeCollector = &cli.StringFlag{
	Name:  "fee",
	Usage: "fee collector contract address",
	Value: "0xfee0000000000000000000000000000000000fee",
}

var L2Witness = &cli.StringFlag{
	Name:  "witness",
	Usage: "l2 witness contract address",
	Value: "0x2210000000000000000000000000000000000221",
}

var L2WitnessBalance = &cli.Float64Flag{
	Name:  "balance",
	Usage: "balance of l2 witness contract",
	Value: 1000000000,
}

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
			FeeCollectorOwner,
			FeeCollector,
			L2Witness,
			L2WitnessBalance,
			OutFile,
		},
	}
}

func GenesisCmd(ctx *cli.Context) error {
	owner := ctx.String(FeeCollectorOwner.Name)
	fee := ctx.String(FeeCollector.Name)
	witness := ctx.String(L2Witness.Name)
	amount := ctx.Float64(L2WitnessBalance.Name)
	outfile := ctx.String(OutFile.Name)
	balance := u256.New(uint64(amount * 1e9)).Mul(web3.Ether(1)).Div(uint64(1e9))

	conf := &deploy.GenesisConfig{
		FeeCollectorOwner:   web3.HexToAddress(owner),
		FeeCollector:        web3.HexToAddress(fee),
		L2CrossLayerWitness: web3.HexToAddress(witness),
		WitnessBalance:      balance.ToBigInt(),
	}

	result := deploy.BuildL2GenesisData(conf)
	utils2.AtomicWriteFile(outfile, utils.JsonString(result))

	return nil
}
