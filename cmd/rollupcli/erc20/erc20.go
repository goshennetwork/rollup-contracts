package erc20

import (
	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract/builtin/erc20"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/urfave/cli/v2"
)

func ERC20Cmd() *cli.Command {
	return &cli.Command{
		Name:        "erc20",
		Subcommands: SubCommand(),
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
	}
}

func SubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "transfer",
			Usage: "transfer erc20  from privKey's account  to addr2",
			Flags: []cli.Flag{
				flags.ToFlag,
				flags.AmountFlag,
				flags.SubmitFlag,
			},
			Action: transferErc20,
		},
	}
}

func transferErc20(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)
	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	to := ctx.String(flags.ToFlag.Name)
	amount := ctx.Float64(flags.AmountFlag.Name)
	submit := ctx.Bool(flags.SubmitFlag.Name)
	signer.Submit = submit
	erc20 := erc20.NewERC20(conf.L1Addresses.FeeToken, signer.Client)
	erc20.Contract().SetFrom(signer.Address())
	depositAmt := erc20.AmountFloatWithDecimals(amount)
	receipt := erc20.Transfer(web3.HexToAddress(to), depositAmt).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	log.Infof("transfer erc20 to %s: %s", to, utils.JsonString(receipt.Thin()))
	balance, err := erc20.BalanceOf(web3.HexToAddress(to), web3.Latest)
	utils.Ensure(err)
	log.Infof("target balance : %s", balance.String())

	return nil
}
