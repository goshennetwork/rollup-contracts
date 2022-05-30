package deploy

import (
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontio/ontology/common/log"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
	utils2 "github.com/ontology-layer-2/rollup-contracts/utils"
	"github.com/urfave/cli/v2"
)

func DeployCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "deploy",
		Usage: "deploy all contract to l1",
		Action: func(ctx *cli.Context) error {
			cfgFile := ctx.String(ConfigFlag.Name)
			verbose := !ctx.Bool(QuietFlag.Name)
			submit := ctx.Bool(SubmitFlag.Name)
			return deployCmd(cfgFile, verbose, submit)
		},
		Flags: []cli.Flag{
			ConfigFlag,
			QuietFlag,
			SubmitFlag,
		},
	}

	return cmd
}

var SubmitFlag = &cli.BoolFlag{
	Name:  "submit",
	Usage: "submit transaction to remote chain",
}

var ConfigFlag = &cli.StringFlag{
	Name:  "cfg",
	Usage: "specify config file",
	Value: "deployment-config.json",
}

var QuietFlag = &cli.BoolFlag{
	Name:    "quiet",
	Aliases: []string{"q"},
	Usage:   "disable print deploy log",
}

func deployCmd(cfgFile string, verbose, submit bool) error {
	conf := &deploy.L1ChainEnv{}
	err := utils.LoadJsonFile(cfgFile, conf)
	utils.Ensure(err)
	l1Client, err := jsonrpc.NewClient(conf.RpcUrl)
	utils.Ensure(err)
	chainId1, err := l1Client.Eth().ChainID()
	utils.Ensure(err)
	if conf.ChainId != chainId1.Uint64() {
		log.Errorf("chain id mismatched, config:%d, remote: %d", conf.ChainId, chainId1)
		return nil
	}
	signer := contract.NewSigner(conf.PrivKey, l1Client, chainId1.Uint64())
	signer.Submit = submit
	results := deploy.DeployL1Contract(signer, conf.L1ChainConfig)

	utils2.AtomicWriteFile("address.json", utils.JsonString(results.Addresses()))

	return nil
}
