package deploy

import (
	"fmt"

	"github.com/goshennetwork/rollup-contracts/binding"
	"github.com/goshennetwork/rollup-contracts/config"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
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
	receipt := binding.DeployProxyAdmin(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	fmt.Println("deploy l2 proxy admin, address:", receipt.ContractAddress.String())
	proxyAdmin := binding.NewProxyAdmin(receipt.ContractAddress, signer.Client)
	proxyAdmin.Contract().SetFrom(signer.Address())
	return proxyAdmin
}

func DeployL2StandardTokenFactory(signer *contract.Signer, l2Bridge web3.Address) web3.Address {
	receipt := binding.DeployL2StandardTokenFactory(signer.Client, signer.Address(), l2Bridge).Sign(signer).SendTransaction(signer).EnsureNoRevert()

	return receipt.ContractAddress
}

func DeployL2TokenBridge(signer *contract.Signer, proxyAdmin web3.Address, bridgeLogicAddress *web3.Address) (
	bridge *binding.L2StandardBridge, logic *binding.L2StandardBridge) {
	bridgeReceipt := binding.DeployL2StandardBridge(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	bridgeLogic := bridgeReceipt.ContractAddress
	if bridgeLogicAddress != nil {
		bridgeLogic = *bridgeLogicAddress
	}
	proxyReceipt := binding.DeployTransparentUpgradeableProxy(signer.Client, signer.Address(),
		bridgeLogic, proxyAdmin, []byte{}).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	fmt.Println("deploy l2 token bridge, address:", proxyReceipt.ContractAddress.String())
	bridge = binding.NewL2StandardBridge(proxyReceipt.ContractAddress, signer.Client)
	bridge.Contract().SetFrom(signer.Address())
	logic = binding.NewL2StandardBridge(bridgeReceipt.ContractAddress, signer.Client)
	logic.Contract().SetFrom(signer.Address())
	return
}

func DeployL2FeeCollector(signer *contract.Signer, owner web3.Address) *binding.L2FeeCollector {
	receipt := binding.DeployL2FeeCollector(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	fmt.Println("deploy l2 fee collector, address:", receipt.ContractAddress.String())
	collector := binding.NewL2FeeCollector(receipt.ContractAddress, signer.Client)
	collector.Contract().SetFrom(signer.Address())
	if owner.IsZero() == false && owner != signer.Address() {
		collector.TransferOwnership(owner).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	}

	return collector
}

func DeployL2CrossLayerWitness(signer *contract.Signer, proxyAdmin web3.Address, witnessLogicAddress *web3.Address) (
	witness *binding.L2CrossLayerWitness, logic *binding.L2CrossLayerWitness) {
	witnessReceipt := binding.DeployL2CrossLayerWitness(signer.Client, signer.Address()).
		Sign(signer).SendTransaction(signer).EnsureNoRevert()
	witnessLogic := witnessReceipt.ContractAddress

	if witnessLogicAddress != nil {
		witnessLogic = *witnessLogicAddress
	}
	proxyReceipt := binding.DeployTransparentUpgradeableProxy(signer.Client, signer.Address(),
		witnessLogic, proxyAdmin, []byte{}).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	fmt.Println("deploy l2 cross layer witness, address:", proxyReceipt.ContractAddress.String())

	witness = binding.NewL2CrossLayerWitness(proxyReceipt.ContractAddress, signer.Client)
	witness.Contract().SetFrom(signer.Address())
	witness.Initialize().Sign(signer).SendTransaction(signer).EnsureNoRevert()

	logic = binding.NewL2CrossLayerWitness(witnessReceipt.ContractAddress, signer.Client)
	logic.Contract().SetFrom(signer.Address())
	return
}

func DeployL2Contracts(signer *contract.Signer, cfg *L2ChainDeployConfig) *L2Contracts {
	proxyAdmin := DeployProxyAdmin(signer)
	collector := DeployL2FeeCollector(signer, cfg.FeeCollectorOwner)
	witness, witnessLogic := DeployL2CrossLayerWitness(signer, proxyAdmin.Contract().Addr(), nil)
	bridge, bridgeLogic := DeployL2TokenBridge(signer, proxyAdmin.Contract().Addr(), nil)

	return &L2Contracts{
		ProxyAdmin:          proxyAdmin,
		L2CrossLayerWitness: witness,
		L2FeeCollector:      collector,
		L2StandardBridge:    bridge,

		L2CrossLayerWitnessLogic: witnessLogic,
		L2StandardBridgeLogic:    bridgeLogic,
	}
}
