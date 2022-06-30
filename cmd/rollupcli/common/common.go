package common

import (
	"encoding/json"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/registry"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/config"
)

func SetUpL1(cfgPath string) (*contract.Signer, *config.RollupCliConfig, error) {
	conf, err := LoadConf(cfgPath)
	if err != nil {
		return nil, nil, err
	}
	return setupSignerL1(conf), conf, nil
}

func SetUpL2(cfgPath string) (*contract.Signer, *config.RollupCliConfig, error) {
	conf, err := LoadConf(cfgPath)
	if err != nil {
		return nil, nil, err
	}
	return setupSignerL2(conf), conf, nil
}

func setupSignerL1(conf *config.RollupCliConfig) *contract.Signer {
	client, err := jsonrpc.NewClient(conf.L1Rpc)
	utils.Ensure(err)
	return getSigner(conf, client)
}

func setupSignerL2(conf *config.RollupCliConfig) *contract.Signer {
	client, err := jsonrpc.NewClient(conf.L2Rpc)
	utils.Ensure(err)
	client.GasLimitFactor = nil
	return getSigner(conf, client)
}

func getSigner(conf *config.RollupCliConfig, client *jsonrpc.Client) *contract.Signer {
	chainId, err := client.Eth().ChainID()
	utils.Ensure(err)
	signer := contract.NewSigner(conf.PrivKey, client, chainId.Uint64())
	return signer
}

func LoadConf(cfgPath string) (*config.RollupCliConfig, error) {
	conf := &config.RollupCliConfig{}
	err := utils.LoadJsonFile(cfgPath, conf)
	if err != nil {
		return nil, err
	}
	registerAbiAndContract(conf)
	return conf, nil
}

func registerAbiAndContract(conf *config.RollupCliConfig) {
	registry.Instance().RegisterFromAbi(binding.L1StandardBridgeAbi())
	registry.Instance().RegisterFromAbi(binding.L1CrossLayerWitnessAbi())
	registry.Instance().RegisterFromAbi(binding.RollupInputChainAbi())
	var l1AddrMap map[string]string
	err := json.Unmarshal([]byte(utils.JsonStr(conf.L1Addresses)), &l1AddrMap)
	utils.Ensure(err)
	for name, addr := range l1AddrMap {
		registry.Instance().RegisterContractAlias(web3.HexToAddress(addr), name)
	}

	var l2AddrMap map[string]string
	err = json.Unmarshal([]byte(utils.JsonStr(conf.L2Genesis.L2ContractAddressConfig)), &l2AddrMap)
	utils.Ensure(err)
	for name, addr := range l2AddrMap {
		registry.Instance().RegisterContractAlias(web3.HexToAddress(addr), name)
	}
}
