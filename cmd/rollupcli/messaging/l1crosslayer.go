package messaging

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/ontio/ontology/common/log"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/urfave/cli/v2"
)

func L1CrossLayerWitnessCommand() *cli.Command {
	return &cli.Command{
		Name: "crosslayer",
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
		Subcommands: []*cli.Command{
			{
				Name:   "pause",
				Action: pauseCmd,
				Usage:  "pause contract",
				Flags: []cli.Flag{
					flags.SubmitFlag,
				},
			},
			{
				Name:   "block",
				Action: blockMessageCmd,
				Usage:  "block l2->l1 message",
				Flags: []cli.Flag{
					flags.DataHash,
					flags.SubmitFlag,
				},
			},
			{
				Name:   "allow",
				Action: allowMessageCmd,
				Usage:  "unblock l2->l1 message",
				Flags: []cli.Flag{
					flags.DataHash,
					flags.SubmitFlag,
				},
			},
			{
				Name:   "send",
				Action: sendMessageCmd,
				Usage:  "send message to l2",
				Flags: []cli.Flag{
					flags.TargetFlag,
					flags.MessageFlag,
					flags.SubmitFlag,
				},
			},
			{
				Name:   "mmr",
				Action: mmrCmd,
				Usage:  "query mmr info",
			},
		},
	}
}

func mmrCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}

	l1Messenger := binding.NewL1CrossLayerWitness(conf.L1Addresses.L1CrossLayerWitness, signer.Client)
	root, err := l1Messenger.MmrRoot()
	utils.Ensure(err)
	size, err := l1Messenger.TotalSize()
	utils.Ensure(err)

	log.Infof("L1 -> L2 cross layer messages, size: %d, mmr root: %s", size,
		web3.Hash(root).String())
	return nil
}

func sendMessageCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	target := ctx.String(flags.TargetFlag.Name)
	mesage := ctx.String(flags.MessageFlag.Name)
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)

	l1Messenger := binding.NewL1CrossLayerWitness(conf.L1Addresses.L1CrossLayerWitness, signer.Client)
	l1Messenger.Contract().SetFrom(signer.Address())
	receipt := l1Messenger.SendMessage(web3.HexToAddress(target), web3.Hex2Bytes(mesage)).Sign(signer).SendTransaction(signer)
	log.Infof("Sends a cross layer message, receipt:%s", utils.JsonString(receipt.Thin()))
	return nil
}

func blockMessageCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	dataHash := ctx.String(flags.DataHash.Name)
	submit := ctx.Bool(flags.SubmitFlag.Name)
	if submit {
		signer.Submit = true
	}
	l1Messenger := binding.NewL1CrossLayerWitness(conf.L1Addresses.L1CrossLayerWitness, signer.Client)
	l1Messenger.Contract().SetFrom(signer.Address())
	receipt := l1Messenger.BlockMessage([][32]byte{web3.HexToHash(dataHash)}).Sign(signer).SendTransaction(signer)
	log.Infof("Block a message, messageHash:%s, receipt:%s", dataHash, utils.JsonString(receipt.Thin()))
	return nil
}

func allowMessageCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	dataHash := ctx.String(flags.DataHash.Name)
	submit := ctx.Bool(flags.SubmitFlag.Name)
	if submit {
		signer.Submit = true
	}
	l1Messenger := binding.NewL1CrossLayerWitness(conf.L1Addresses.L1CrossLayerWitness, signer.Client)
	l1Messenger.Contract().SetFrom(signer.Address())
	receipt := l1Messenger.AllowMessage([][32]byte{web3.HexToHash(dataHash)}).Sign(signer).SendTransaction(signer)
	log.Infof("Allow a message, messageHash:%s, receipt:%s", dataHash, utils.JsonString(receipt.Thin()))
	return nil
}

func pauseCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	submit := ctx.Bool(flags.SubmitFlag.Name)
	if submit {
		signer.Submit = true
	}
	l1Messenger := binding.NewL1CrossLayerWitness(conf.L1Addresses.L1CrossLayerWitness, signer.Client)
	l1Messenger.Contract().SetFrom(signer.Address())
	receipt := l1Messenger.Pause().Sign(signer).SendTransaction(signer)
	log.Infof("Pause relaying: %s", utils.JsonString(receipt.Thin()))
	return nil
}
