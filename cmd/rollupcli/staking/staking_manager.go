package staking

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract/builtin/erc20"
	"github.com/laizy/web3/utils"
	"github.com/ontio/ontology/common/log"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/urfave/cli/v2"
)

func StakingCommand() *cli.Command {
	return &cli.Command{
		Name: "staking",
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
		Subcommands: []*cli.Command{
			{
				Name:   "deposit",
				Action: depositCmd,
				Flags: []cli.Flag{
					flags.SubmitFlag,
				},
			},
			{
				Name:   "startWithdrawal",
				Action: startWithdrawalCmd,
				Flags: []cli.Flag{
					flags.SubmitFlag,
				},
				Description: "Starts the withdrawal for a publisher",
			},
			/*
				{
					Name:   "finalizeWithdrawal",
					Action: finalizeWithdrawalCmd,
					Flags: []cli.Flag{
						flags.SubmitFlag,
					},
					Description: "Finalizes a pending withdrawal from a publisher",
				},
			*/
		},
	}
}

func startWithdrawalCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)

	staking := binding.NewStakingManager(conf.L1Addresses.StakingManager, signer.Client)
	staking.Contract().SetFrom(signer.Address())
	receipt := staking.StartWithdrawal().Sign(signer).SendTransaction(signer)
	log.Infof("start withdrawal receipt:%s", utils.JsonString(receipt.Thin()))
	return nil
}

/*
type StakingState uint8

const UNSTAKED = 0
const STAKING = 1
const WITHDRAWING = 2
const SLASHING = 3
func finalizeWithdrawalCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	submit := ctx.Bool(flags.SubmitFlag.Name)
	if submit {
		signer.Submit = true
	}
	staking := binding.NewStakingManager(conf.L1Addresses.StakingManager, signer.Client)
	staking.Contract().SetFrom(signer.Address())
	status, needComfirmHeight, _, _, _, err := staking.GetStakingInfo(signer.Address())
	if status != WITHDRAWING {
		log.Warn("not in withdrawing state")
		return nil
	}

	// todo: get state info
	receipt := staking.FinalizeWithdrawal().Sign(signer).SendTransaction(signer)
	log.Infof("finailize withdraw receipt:%s", utils.JsonString(receipt.Thin()))
	return nil
}
*/

func depositCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)
	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	submit := ctx.Bool(flags.SubmitFlag.Name)
	if submit {
		signer.Submit = true
	}
	staking := binding.NewStakingManager(conf.L1Addresses.StakingManager, signer.Client)
	staking.Contract().SetFrom(signer.Address())
	tokenAddr, err := staking.Token(web3.Latest)
	if err != nil {
		return err
	}
	log.Infof("staking token address:%s", tokenAddr.String())
	token := erc20.NewERC20(tokenAddr, signer.Client)
	allowance, err := token.Allowance(signer.Address(), conf.L1Addresses.StakingManager)
	if err != nil {
		return err
	}
	if allowance.Cmp(web3.Ether(1)) < 0 {
		log.Info("allowance not enough, do approve first...")
		price, err := staking.Price()
		utils.Ensure(err)
		receipt := token.Approve(conf.L1Addresses.StakingManager, price).Sign(signer).SendTransaction(signer)
		log.Infof("approve receipt:%s", utils.JsonString(receipt.Thin()))
	}
	receipt := staking.Deposit().Sign(signer).SendTransaction(signer)
	log.Infof("staking manager deposit, receipt:%s", utils.JsonString(receipt.Thin()))
	return nil
}
