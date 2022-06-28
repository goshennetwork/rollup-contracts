package gateway

import (
	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/contract/builtin/erc20"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/u256"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/urfave/cli/v2"
)

func GatewayCommand() *cli.Command {
	return &cli.Command{
		Name:        "gateway",
		Subcommands: gatewayCommands(),
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
	}
}

func gatewayCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "depositEth",
			Action: DepositEthCmd,
			Flags: []cli.Flag{
				flags.AmountFlag,
				flags.ToFlag,
				flags.SubmitFlag,
			},
		},
		{
			Name:   "depositERC20",
			Action: DepositERC20Cmd,
			Flags: []cli.Flag{
				flags.L1TokenFlag,
				flags.L2TokenFlag,
				flags.AmountFlag,
				flags.ToFlag,
				flags.SubmitFlag,
			},
		},
		{
			Name:   "withdrawToERC20",
			Action: WithdrawToERC20Cmd,
			Flags: []cli.Flag{
				flags.L2TokenFlag,
				flags.AmountFlag,
				flags.ToFlag,
				flags.SubmitFlag,
			},
		},
	}
}

func DepositERC20Cmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	l1Token := ctx.String(flags.L1TokenFlag.Name)
	l2Token := ctx.String(flags.L2TokenFlag.Name)
	amount := ctx.Float64(flags.AmountFlag.Name)
	to := ctx.String(flags.ToFlag.Name)
	if to == "" {
		to = signer.Address().String()
	}
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)

	l1Tok := erc20.NewERC20(web3.HexToAddress(l1Token), signer.Client)
	decimal, err := l1Tok.Decimals(web3.Latest)
	if err != nil {
		return err
	}

	depositAmt := u256.New(uint64(amount * 1e9)).Mul(u256.New(1).ExpUint8(decimal)).Div(uint64(1e9))

	balance, err := signer.Eth().GetBalance(signer.Address(), web3.Latest)
	utils.Ensure(err)
	log.Infof("balance of %s is %s ether", signer.Address().String(), u256.New(balance).ToFixNum(18))

	DepositERC20ToL2(signer, web3.HexToAddress(to), web3.HexToAddress(l1Token), web3.HexToAddress(l2Token), conf.L1Addresses.L1StandardBridge, depositAmt)
	return nil
}

func DepositEthCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)
	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	amount := ctx.Float64(flags.AmountFlag.Name)
	to := ctx.String(flags.ToFlag.Name)
	if to == "" {
		to = signer.Address().String()
	}
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)

	depositAmt := u256.New(uint64(amount * 1e9)).Mul(web3.Ether(1)).Div(uint64(1e9))

	balance, err := signer.Eth().GetBalance(signer.Address(), web3.Latest)
	utils.Ensure(err)
	log.Infof("balance of %s is %s ether", signer.Address().String(), u256.New(balance).ToFixNum(18))

	DepositEthToL2(signer, web3.HexToAddress(to), conf.L1Addresses.L1StandardBridge, depositAmt)
	return nil
}

func WithdrawToERC20Cmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)
	signer, conf, err := common.SetUpL2(path)
	if err != nil {
		return err
	}
	l2Token := ctx.String(flags.L2TokenFlag.Name)
	amount := ctx.Float64(flags.AmountFlag.Name)
	to := ctx.String(flags.ToFlag.Name)
	if to == "" {
		to = signer.Address().String()
	}
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)
	l2Tok := erc20.NewERC20(web3.HexToAddress(l2Token), signer.Client)
	decimal, err := l2Tok.Decimals(web3.Latest)
	if err != nil {
		return err
	}
	withDrawAmt := u256.New(uint64(amount * 1e9)).Mul(u256.New(10).ExpUint8(decimal)).Div(uint64(1e9))
	balance, err := signer.Eth().GetBalance(signer.Address(), web3.Latest)
	utils.Ensure(err)
	log.Infof("balance of %s is %s ether", signer.Address().String(), u256.New(balance).ToFixNum(18))
	WithdrawToERC20ToL1(signer, conf.L2Genesis.L2StandardBridge, web3.HexToAddress(to), web3.HexToAddress(l2Token), withDrawAmt)
	return nil
}

func DepositEthToL2(signer *contract.Signer, to, l1Bridge web3.Address, depositAmt u256.Int) {
	gateway := binding.NewL1StandardBridge(l1Bridge, signer.Client)
	gateway.Contract().SetFrom(signer.Address())
	receipt := gateway.DepositETHTo(to, nil).SetValue(depositAmt.ToBigInt()).Sign(signer).SendTransaction(signer)
	log.Infof("deposit eth to l2: %s", utils.JsonString(receipt.Thin()))
}

func DepositERC20ToL2(signer *contract.Signer, to, l1Bridge, l1Token, l2Token web3.Address, depositAmt u256.Int) {
	gateway := binding.NewL1StandardBridge(l1Bridge, signer.Client)
	gateway.Contract().SetFrom(signer.Address())
	receipt := gateway.DepositERC20To(l1Token, l2Token, to, depositAmt.ToBigInt(), nil).Sign(signer).SendTransaction(signer)
	log.Infof("deposit erc20 to l2: %s", utils.JsonString(receipt.Thin()))
}

func WithdrawToERC20ToL1(signer *contract.Signer, l2Bridge, to, l2Token web3.Address, withdrawAmt u256.Int) {
	gateway := binding.NewL2StandardBridge(l2Bridge, signer.Client)
	gateway.Contract().SetFrom(signer.Address())
	receipt := gateway.WithdrawTo(l2Token, to, withdrawAmt.ToBigInt(), nil).Sign(signer).SendTransaction(signer)
	log.Infof("withdrawal erc20 to l1: %s", utils.JsonString(receipt.Thin()))
}
