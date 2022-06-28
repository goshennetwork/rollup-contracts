package deploy

import (
	"github.com/laizy/log"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/cmd/rollupcli/common"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
	utils2 "github.com/ontology-layer-2/rollup-contracts/utils"
)

func deployL2Cmd(cfgFile string, verbose, submit bool) error {
	cfg := LoadDeployConfig(cfgFile)
	conf := cfg.L2
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
	results := deploy.DeployL2Contracts(signer, conf.ChainConfig)

	utils2.AtomicWriteFile("addressl2.json", utils.JsonString(results.Addresses()))

	return nil
}

func initL2BridgeCmd(cfgFile string, verbose, submit bool) error {
	signer, conf, err := common.SetUpL2(cfgFile)
	if err != nil {
		return err
	}
	signer.Submit = submit
	l2Witness := conf.L2Genesis.L2CrossLayerWitness
	l1bridge := conf.L1Addresses.L1StandardBridge
	if l2Witness.IsZero() || l1bridge.IsZero() {
		log.Error("need set l2 witness or l1 bridge")
		return nil
	}
	bridge := binding.NewL2StandardBridge(conf.L2Genesis.L2StandardBridge, signer.Client)
	bridge.Contract().SetFrom(signer.Address())
	bridge.Initialize(conf.L2Genesis.L2CrossLayerWitness, conf.L1Addresses.L1StandardBridge).Sign(signer).SendTransaction(signer)

	return nil
}
