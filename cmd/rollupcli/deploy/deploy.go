package deploy

import (
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontio/ontology/common/log"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
	utils2 "github.com/ontology-layer-2/rollup-contracts/utils"
	"github.com/urfave/cli/v2"
)

type DeployConfig struct {
	L1 *deploy.L1ChainEnv
	L2 *deploy.L2ChainEnv
}

func LoadDeployConfig(file string) *DeployConfig {
	cfg := &DeployConfig{}
	err := utils.LoadJsonFile(file, cfg)
	utils.Ensure(err)
	cfg.L1.ChainConfig.L2ChainId = cfg.L2.ChainId
	return cfg
}

func DeployCmd() *cli.Command {
	cmd := &cli.Command{
		Name:        "deploy",
		Usage:       "deploy and initialize contracts",
		Subcommands: deploySubCommands(),
	}

	return cmd
}

func deploySubCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "l1",
			Usage: "deploy all l1 contracts",
			Action: func(ctx *cli.Context) error {
				cfgFile := ctx.String(ConfigFlag.Name)
				verbose := !ctx.Bool(QuietFlag.Name)
				submit := ctx.Bool(flags.SubmitFlag.Name)
				return deployL1Cmd(cfgFile, verbose, submit)
			},
			Flags: []cli.Flag{
				ConfigFlag,
				QuietFlag,
				flags.SubmitFlag,
			},
		},
		{
			Name:  "l2",
			Usage: "deploy all l2 contracts",
			Action: func(ctx *cli.Context) error {
				cfgFile := ctx.String(ConfigFlag.Name)
				verbose := !ctx.Bool(QuietFlag.Name)
				submit := ctx.Bool(flags.SubmitFlag.Name)
				return deployL2Cmd(cfgFile, verbose, submit)
			},
			Flags: []cli.Flag{
				ConfigFlag,
				QuietFlag,
				flags.SubmitFlag,
			},
		},
		{
			Name:  "l2init",
			Usage: "initialize l2 token bridge",
			Action: func(ctx *cli.Context) error {
				cfgFile := ctx.String(flags.ConfigFlag.Name)
				verbose := !ctx.Bool(QuietFlag.Name)
				submit := ctx.Bool(flags.SubmitFlag.Name)
				return initL2BridgeCmd(cfgFile, verbose, submit)
			},
			Flags: []cli.Flag{
				flags.ConfigFlag,
				QuietFlag,
				flags.SubmitFlag,
			},
		},
	}
}

var ConfigFlag = &cli.StringFlag{
	Name:  "cfg",
	Usage: "specify config file",
	Value: "deploy-config.json",
}

var QuietFlag = &cli.BoolFlag{
	Name:    "quiet",
	Aliases: []string{"q"},
	Usage:   "disable print deploy log",
}

func deployL1Cmd(cfgFile string, verbose, submit bool) error {
	cfg := LoadDeployConfig(cfgFile)
	conf := cfg.L1
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
	results := deploy.DeployL1Contracts(signer, conf.ChainConfig)

	utils2.AtomicWriteFile("addressl1.json", utils.JsonString(results.Addresses()))

	return nil
}
