package deploy

import (
	"fmt"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
)

type L2ChainEnv struct {
	ChainId       uint64
	RpcUrl        string
	PrivKey       string
	L2ChainConfig *L2ChainDeployConfig
}

type L2ChainDeployConfig struct {
	FeeCollectorOwner       web3.Address
}

type L2Contracts struct {
	L2CrossLayerWitness *binding.L2CrossLayerWitness
	L2FeeCollector *binding.L2FeeCollector
}

func DeployL2FeeCollector(signer *contract.Signer, owner web3.Address) *binding.L2FeeCollector {
	receipt := binding.DeployL2FeeCollector(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l2 fee collector, address:", receipt.ContractAddress.String())
	collector := binding.NewL2FeeCollector(receipt.ContractAddress, signer.Client)
	collector.Contract().SetFrom(signer.Address())
	if owner.IsZero() == false && owner != signer.Address() {
		collector.TransferOwnership(owner).Sign(signer).SendTransaction(signer)
	}

	return collector
}

func DeployL2CrossLayerWitness(signer *contract.Signer) *binding.L2CrossLayerWitness {
	receipt := binding.DeployL2CrossLayerWitness(signer.Client, signer.Address()).
		Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l2 cross layer witness, address:", receipt.ContractAddress.String())
	l2CrossLayerWitness := binding.NewL2CrossLayerWitness(receipt.ContractAddress, signer.Client)
	l2CrossLayerWitness.Contract().SetFrom(signer.Address())
	// todo: deploy after proxy
	l2CrossLayerWitness.Initialize().Sign(signer).SendTransaction(signer)

	return l2CrossLayerWitness
}

func DeployL2Contracts(signer *contract.Signer, cfg *L2ChainDeployConfig) *L2Contracts {
	collector := DeployL2FeeCollector(signer, cfg.FeeCollectorOwner)
	witness := 	DeployL2CrossLayerWitness(signer)

	return &L2Contracts{
		L2CrossLayerWitness: witness,
		L2FeeCollector:      collector,
	}
}

// func BuildL2GenesisData(signer *contract.Signer)