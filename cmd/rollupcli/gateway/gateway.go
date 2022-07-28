package gateway

import (
	"fmt"
	"math/big"

	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/contract/builtin/erc20"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/u256"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	cli "github.com/urfave/cli/v2"
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
		{
			Name:   "withdrawEth",
			Action: WithdrawEthCmd,
			Flags: []cli.Flag{
				flags.AmountFlag,
				flags.ToFlag,
				flags.SubmitFlag,
			},
		},
		{
			Name:   "relayMsgToL1",
			Action: FinalizeWithdraw,
			Flags: []cli.Flag{
				flags.MsgIndexFlag,
				flags.SubmitFlag,
			},
			Usage: "relay l2->l1 msg to l1",
		},
	}
}

func FinalizeWithdraw(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)

	signer, conf, err := common.SetUpL1(path)
	if err != nil {
		return err
	}
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)
	msgIndex := ctx.Uint64(flags.MsgIndexFlag.Name)
	return RelayMsg(signer, conf.L2Rpc, conf.L1Addresses.L1CrossLayerWitness, msgIndex)
}

func RelayMsg(signer *contract.Signer, l2url string, L1CrossLayerWitness web3.Address, msgIndex uint64) error {
	l2client, err := jsonrpc.NewClient(l2url)
	utils.Ensure(err)
	params, err := l2client.L2().GetL1RelayMsgParams(msgIndex)
	if err != nil {
		return err
	}
	log.Info("r1cs debug", "params", utils.JsonStr(params))
	if params == nil {
		return fmt.Errorf("no params found, msgIndex: %d", msgIndex)
	}

	l1Cross := binding.NewL1CrossLayerWitness(L1CrossLayerWitness, signer.Client)
	l1Cross.Contract().SetFrom(signer.Address())
	var proof [][32]byte
	for _, h := range params.Proof {
		proof = append(proof, h)
	}
	stateInfo := binding.FromRPCStateInfo(params.StateInfo)
	log.Info("r1cs debug", "stateHash", stateInfo.Hash(), "stateInfo", utils.JsonStr(stateInfo))
	txn := l1Cross.RelayMessage(params.Target, params.Sender, params.Message, uint64(params.MessageIndex), params.RLPHeader, *stateInfo, proof).Sign(signer)
	log.Info("r1cs debug", "data", txn.Input)
	r := txn.SendTransaction(signer).EnsureNoRevert()
	log.Info("relay message succeed", "msgIndex", msgIndex, "txHash", r.TransactionHash, "log", utils.JsonStr(r))
	return nil
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
	depositAmt := u256.New(l1Tok.AmountFloatWithDecimals(amount))
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
	withDrawAmt := u256.New(l2Tok.AmountFloatWithDecimals(amount))
	balance, err := signer.Eth().GetBalance(signer.Address(), web3.Latest)
	utils.Ensure(err)
	log.Infof("balance of %s is %s ether", signer.Address().String(), u256.New(balance).ToFixNum(18))
	WithdrawToERC20ToL1(signer, conf.L2Genesis.L2StandardBridge, web3.HexToAddress(to), web3.HexToAddress(l2Token), withDrawAmt)
	return nil
}

func WithdrawEthCmd(ctx *cli.Context) error {
	path := ctx.String(flags.ConfigFlag.Name)
	signer, conf, err := common.SetUpL2(path)
	if err != nil {
		return err
	}
	amount := ctx.Float64(flags.AmountFlag.Name)
	to := ctx.String(flags.ToFlag.Name)
	if to == "" {
		to = signer.Address().String()
	}
	signer.Submit = ctx.Bool(flags.SubmitFlag.Name)
	withdrawAmt := u256.New(uint64(amount * 1e9)).Mul(web3.Ether(1)).Div(uint64(1e9))
	WithdrawEthTo(signer, conf.L2Genesis.L2StandardBridge, web3.HexToAddress(to), withdrawAmt.ToBigInt())
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

func WithdrawEthTo(signer *contract.Signer, l2Bridge web3.Address, target web3.Address, amount *big.Int) {
	gateway := binding.NewL2StandardBridge(l2Bridge, signer.Client)
	gateway.Contract().SetFrom(signer.Address())
	r := gateway.WithdrawETHTo(target, nil).SetValue(amount).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	log.Infof("withdrawal eth to l1 :%s", utils.JsonString(r.Thin()))
}
