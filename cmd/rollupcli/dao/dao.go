package dao

import (
	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/urfave/cli/v2"
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "dao",
		Subcommands: SubCommand(),
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
	}
}

func SubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "sequencerWhitelist",
			Usage: "transfer erc20  from privKey's account  to addr2",
			Flags: []cli.Flag{
				flags.AccountFlag,
				flags.EnabledFlag,
				flags.SubmitFlag,
			},
			Action: sequencerWhitelist,
		},
	}
}

func sequencerWhitelist(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)
	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	acc := ctx.String(flags.AccountFlag.Name)
	enabled := ctx.Bool(flags.EnabledFlag.Name)
	submit := ctx.Bool(flags.SubmitFlag.Name)
	signer.Submit = submit
	setSequencerWhitelist(signer, conf.L1Addresses.DAO, web3.HexToAddress(acc), enabled)
	return nil
}

func setSequencerWhitelist(signer *contract.Signer, dao web3.Address, sequencer web3.Address, enabled bool) {
	c := binding.NewDAO(dao, signer.Client)
	c.Contract().SetFrom(signer.Address())
	receipt := c.SetSequencerWhitelist(sequencer, enabled).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	log.Info("set sequencer whitelist", "dao", dao, "sequencer", sequencer, "enabled", enabled, "receipt", utils.JsonStr(receipt))
}
