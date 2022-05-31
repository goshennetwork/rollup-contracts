package deploy

import (
	"fmt"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/config"
)

type L2ChainEnv struct {
	ChainId     uint64
	RpcUrl      string
	PrivKey     string
	ChainConfig *L2ChainDeployConfig
}

type L2ChainDeployConfig struct {
	FeeCollectorOwner web3.Address
	L1TokenBridge     web3.Address
}

type L2Contracts struct {
	L2CrossLayerWitness *binding.L2CrossLayerWitness
	L2FeeCollector      *binding.L2FeeCollector
	L2StandardBridge    *binding.L2StandardBridge
}

func (self *L2Contracts) Addresses() *config.L2ContractAddressConfig {
	return &config.L2ContractAddressConfig{
		L2CrossLayerWitness: self.L2CrossLayerWitness.Contract().Addr(),
		L2FeeCollector:      self.L2FeeCollector.Contract().Addr(),
		L2StandardBridge:    self.L2StandardBridge.Contract().Addr(),
	}
}

func DeployL2TokenBridge(signer *contract.Signer, l2Witness, l1Bridge web3.Address) *binding.L2StandardBridge {
	receipt := binding.DeployL2StandardBridge(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l2 token bridge, address:", receipt.ContractAddress.String())
	bridge := binding.NewL2StandardBridge(receipt.ContractAddress, signer.Client)
	bridge.Contract().SetFrom(signer.Address())
	if l2Witness.IsZero() == false && l1Bridge.IsZero() == false {
		bridge.Initialize(l2Witness, l1Bridge).Sign(signer).SendTransaction(signer)
	}

	return bridge
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
	witness := DeployL2CrossLayerWitness(signer)
	bridge := DeployL2TokenBridge(signer, witness.Contract().Addr(), cfg.L1TokenBridge)

	return &L2Contracts{
		L2CrossLayerWitness: witness,
		L2FeeCollector:      collector,
		L2StandardBridge:    bridge,
	}
}

// func BuildL2GenesisData(signer *contract.Signer)
