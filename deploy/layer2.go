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
	ProxyAdmin          *binding.ProxyAdmin
	L2CrossLayerWitness *binding.L2CrossLayerWitness
	L2FeeCollector      *binding.L2FeeCollector
	L2StandardBridge    *binding.L2StandardBridge

	L2CrossLayerWitnessLogic *binding.L2CrossLayerWitness
	L2StandardBridgeLogic    *binding.L2StandardBridge
}

func (self *L2Contracts) Addresses() *config.L2ContractAddressConfig {
	return &config.L2ContractAddressConfig{
		ProxyAdmin:          self.ProxyAdmin.Contract().Addr(),
		L2CrossLayerWitness: self.L2CrossLayerWitness.Contract().Addr(),
		L2FeeCollector:      self.L2FeeCollector.Contract().Addr(),
		L2StandardBridge:    self.L2StandardBridge.Contract().Addr(),

		L2CrossLayerWitnessLogic: self.L2CrossLayerWitnessLogic.Contract().Addr(),
		L2StandardBridgeLogic:    self.L2StandardBridgeLogic.Contract().Addr(),
	}
}

func DeployProxyAdmin(signer *contract.Signer) *binding.ProxyAdmin {
	receipt := binding.DeployProxyAdmin(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l2 proxy admin, address:", receipt.ContractAddress.String())
	proxyAdmin := binding.NewProxyAdmin(receipt.ContractAddress, signer.Client)
	proxyAdmin.Contract().SetFrom(signer.Address())
	return proxyAdmin
}

func DeployL2TokenBridge(signer *contract.Signer, proxyAdmin web3.Address) (
	bridge *binding.L2StandardBridge, logic *binding.L2StandardBridge) {
	bridgeReceipt := binding.DeployL2StandardBridge(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(bridgeReceipt.Status == 1)
	// don't initialize bridge while deploy
	proxyReceipt := binding.DeployTransparentUpgradeableProxy(signer.Client, signer.Address(),
		bridgeReceipt.ContractAddress, proxyAdmin, []byte{}).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(proxyReceipt.Status == 1)
	fmt.Println("deploy l2 token bridge, address:", proxyReceipt.ContractAddress.String())
	bridge = binding.NewL2StandardBridge(proxyReceipt.ContractAddress, signer.Client)
	bridge.Contract().SetFrom(signer.Address())
	logic = binding.NewL2StandardBridge(bridgeReceipt.ContractAddress, signer.Client)
	logic.Contract().SetFrom(signer.Address())
	return
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

func DeployL2CrossLayerWitness(signer *contract.Signer, proxyAdmin web3.Address) (
	witness *binding.L2CrossLayerWitness, logic *binding.L2CrossLayerWitness) {
	witnessReceipt := binding.DeployL2CrossLayerWitness(signer.Client, signer.Address()).
		Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(witnessReceipt.Status == 1)
	proxyReceipt := binding.DeployTransparentUpgradeableProxy(signer.Client, signer.Address(),
		witnessReceipt.ContractAddress, proxyAdmin, []byte{}).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(proxyReceipt.Status == 1)
	fmt.Println("deploy l2 cross layer witness, address:", proxyReceipt.ContractAddress.String())

	witness = binding.NewL2CrossLayerWitness(proxyReceipt.ContractAddress, signer.Client)
	witness.Contract().SetFrom(signer.Address())
	r := witness.Initialize().Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(r.Status == 1)

	logic = binding.NewL2CrossLayerWitness(witnessReceipt.ContractAddress, signer.Client)
	logic.Contract().SetFrom(signer.Address())
	return
}

func DeployL2Contracts(signer *contract.Signer, cfg *L2ChainDeployConfig) *L2Contracts {
	proxyAdmin := DeployProxyAdmin(signer)
	collector := DeployL2FeeCollector(signer, cfg.FeeCollectorOwner)
	witness, witnessLogic := DeployL2CrossLayerWitness(signer, proxyAdmin.Contract().Addr())
	bridge, bridgeLogic := DeployL2TokenBridge(signer, proxyAdmin.Contract().Addr())

	return &L2Contracts{
		ProxyAdmin:          proxyAdmin,
		L2CrossLayerWitness: witness,
		L2FeeCollector:      collector,
		L2StandardBridge:    bridge,

		L2CrossLayerWitnessLogic: witnessLogic,
		L2StandardBridgeLogic:    bridgeLogic,
	}
}
