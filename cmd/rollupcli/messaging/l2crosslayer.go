package messaging

import (
	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/common/hexutil"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	cli "github.com/urfave/cli/v2"
)

func L2CrossLayerWitnessCommand() *cli.Command {
	return &cli.Command{
		Name: "l2",
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
		Subcommands: []*cli.Command{
			{
				Name:   "send",
				Action: l2SendMessageCmd,
				Usage:  "send message to l1",
				Flags: []cli.Flag{
					flags.TargetFlag,
					flags.MessageFlag,
					flags.SubmitFlag,
				},
			},
		},
	}
}

func l2SendMessageCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL2(path)
	if err != nil {
		return err
	}
	target := ctx.String(flags.TargetFlag.Name)
	message := ctx.String(flags.MessageFlag.Name)
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)

	witness := binding.NewL2CrossLayerWitness(conf.L2Genesis.L2CrossLayerWitness, signer.Client)
	witness.Contract().SetFrom(signer.Address())
	receipt := witness.SendMessage(web3.HexToAddress(target), hexutil.MustDecode(message)).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	log.Infof("Sends a cross layer message, receipt:%s", utils.JsonString(receipt.Thin()))
	return nil
}
