package cfg

import (
	"os"

	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/flags"
	"github.com/ontology-layer-2/rollup-contracts/config"
	"github.com/urfave/cli/v2"
)

func CfgCommand() *cli.Command {
	return &cli.Command{
		Name:   "cfg",
		Action: genCfg,
		Flags: []cli.Flag{
			flags.L1RpcFlag,
			flags.L2RpcFlag,
			flags.NameFlag,
			flags.ResolverFlag,
			flags.PrivateFlag,
		},
		Description: "generate config json file from address manager",
	}
}

func genCfg(ctx *cli.Context) error {
	l1, l2, name, resolverAddr, key := ctx.String(flags.L1RpcFlag.Name), ctx.String(flags.L2RpcFlag.Name), ctx.String(flags.NameFlag.Name), web3.HexToAddress(ctx.String(flags.ResolverFlag.Name)), ctx.String(flags.PrivateFlag.Name)
	var cfg config.RollupCliConfig
	cfg.L2Genesis = &config.DefaultL2GenesisConfig
	cfg.L1Rpc, cfg.L2Rpc, cfg.PrivKey, cfg.MinConfirmBlockNum = l1, l2, key, 6

	l1client, err := jsonrpc.NewClient(l1)
	utils.Ensure(err)
	resolver := binding.NewAddressManager(resolverAddr, l1client)
	challengeBeacon := getAddr(resolver.ChallengeBeacon())
	beacon := binding.NewUpgradeableBeacon(challengeBeacon, l1client)
	challengeLogic := getAddr(beacon.Implementation())
	l1contracts := &config.L1ContractAddressConfig{
		resolverAddr,
		getAddr(resolver.RollupInputChainContainer()),
		getAddr(resolver.RollupStateChainContainer()),
		getAddr(resolver.RollupInputChain()),
		getAddr(resolver.RollupStateChain()),
		getAddr(resolver.L1CrossLayerWitness()),
		getAddr(resolver.L1StandardBridge()),
		getAddr(resolver.StakingManager()),
		challengeBeacon,
		challengeLogic,
		getAddr(resolver.ChallengeFactory()),
		getAddr(resolver.FeeToken()),
		getAddr(resolver.Dao()),
	}

	cfg.L1Addresses = l1contracts
	//todo: get debloy height
	cfg.DeployOnL1Height = 0
	return os.WriteFile(name, []byte(utils.JsonString(cfg)), 0x0666)
}

func getAddr(addr web3.Address, err error) web3.Address {
	utils.Ensure(err)
	return addr
}
