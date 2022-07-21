package deploy

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/config"
)

type L1ChainEnv struct {
	ChainId     uint64
	RpcUrl      string
	PrivKey     string
	ChainConfig *config.L1ChainDeployConfig
}

type L1Contracts struct {
	AddressManager      *binding.AddressManager
	InputChainStorage   *binding.ChainStorageContainer
	StateChainStorage   *binding.ChainStorageContainer
	RollupInputChain    *binding.RollupInputChain
	RollupStateChain    *binding.RollupStateChain
	L1CrossLayerWitness *binding.L1CrossLayerWitness
	L1StandardBridge    *binding.L1StandardBridge
	StakingManager      *binding.StakingManager
	ChallengeBeacon     *binding.UpgradeableBeacon
	ChallengeLogic      *binding.Challenge
	ChallengeFactory    *binding.ChallengeFactory
	FeeToken            *binding.ERC20
	DAO                 *binding.DAO
}

func (self *L1Contracts) Addresses() *config.L1ContractAddressConfig {
	return &config.L1ContractAddressConfig{
		AddressManager:      self.AddressManager.Contract().Addr(),
		InputChainStorage:   self.InputChainStorage.Contract().Addr(),
		StateChainStorage:   self.StateChainStorage.Contract().Addr(),
		RollupInputChain:    self.RollupInputChain.Contract().Addr(),
		RollupStateChain:    self.RollupStateChain.Contract().Addr(),
		L1CrossLayerWitness: self.L1CrossLayerWitness.Contract().Addr(),
		L1StandardBridge:    self.L1StandardBridge.Contract().Addr(),
		StakingManager:      self.StakingManager.Contract().Addr(),
		ChallengeBeacon:     self.ChallengeBeacon.Contract().Addr(),
		ChallengeLogic:      self.ChallengeLogic.Contract().Addr(),
		ChallengeFactory:    self.ChallengeFactory.Contract().Addr(),
		FeeToken:            self.FeeToken.Contract().Addr(),
		DAO:                 self.DAO.Contract().Addr(),
	}

}

func DeployChallengeLogic(signer *contract.Signer) *binding.Challenge {
	receipt := binding.DeployChallenge(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)

	challenge := binding.NewChallenge(receipt.ContractAddress, signer.Client)
	challenge.Contract().SetFrom(signer.Address())

	return challenge
}

func DeployBeacon(signer *contract.Signer, impl web3.Address) *binding.UpgradeableBeacon {
	receipt := binding.DeployUpgradeableBeacon(signer.Client, signer.Address(), impl).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)

	beacon := binding.NewUpgradeableBeacon(receipt.ContractAddress, signer.Client)
	beacon.Contract().SetFrom(signer.Address())

	return beacon
}

func DeployTestFeeToken(signer *contract.Signer) *binding.ERC20 {
	receipt := binding.DeployTestERC20(signer.Client, signer.Address(), "TestFeeToken", "TFT").Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)

	feeToken := binding.NewERC20(receipt.ContractAddress, signer.Client)
	feeToken.Contract().SetFrom(signer.Address())

	return feeToken
}

func DeployDAOLogic(signer *contract.Signer) *binding.DAO {
	receipt := binding.DeployDAO(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	dao := binding.NewDAO(receipt.ContractAddress, signer.Client)
	dao.Contract().SetFrom(signer.Address())
	return dao
}

func DeployChallengeFactoryLogic(signer *contract.Signer) *binding.ChallengeFactory {
	receipt := binding.DeployChallengeFactory(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	factory := binding.NewChallengeFactory(receipt.ContractAddress, signer.Client)
	factory.Contract().SetFrom(signer.Address())
	return factory
}

func DeployStakingManagerLogic(signer *contract.Signer) *binding.StakingManager {
	receipt := binding.DeployStakingManager(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	staking := binding.NewStakingManager(receipt.ContractAddress, signer.Client)
	staking.Contract().SetFrom(signer.Address())
	return staking
}

func DeployRollupInputChainLogic(signer *contract.Signer) *binding.RollupInputChain {
	receipt := binding.DeployRollupInputChain(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	rollupInputChain := binding.NewRollupInputChain(receipt.ContractAddress, signer.Client)
	rollupInputChain.Contract().SetFrom(signer.Address())
	return rollupInputChain
}

func DeployRollupStateChainLogic(signer *contract.Signer) *binding.RollupStateChain {
	receipt := binding.DeployRollupStateChain(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	rollupStateChain := binding.NewRollupStateChain(receipt.ContractAddress, signer.Client)
	rollupStateChain.Contract().SetFrom(signer.Address())
	return rollupStateChain
}

func DeployChainStorageLogic(signer *contract.Signer) *binding.ChainStorageContainer {
	receipt := binding.DeployChainStorageContainer(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	chainStorage := binding.NewChainStorageContainer(receipt.ContractAddress, signer.Client)
	chainStorage.Contract().SetFrom(signer.Address())
	return chainStorage
}

func DeployL1CrossLayerWitnessLogic(signer *contract.Signer) *binding.L1CrossLayerWitness {
	receipt := binding.DeployL1CrossLayerWitness(signer.Client, signer.Address()).
		Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l1 cross layer witness, address:", receipt.ContractAddress.String())
	l1CrossLayerWitness := binding.NewL1CrossLayerWitness(receipt.ContractAddress, signer.Client)
	l1CrossLayerWitness.Contract().SetFrom(signer.Address())
	return l1CrossLayerWitness
}

func DeployUpgradeProxy(signer *contract.Signer, logic web3.Address, admin web3.Address, data []byte) *binding.TransparentUpgradeableProxy {
	receipt := binding.DeployTransparentUpgradeableProxy(signer.Client, signer.Address(), logic, admin, data).Sign(signer).SendTransaction(signer).EnsureNoRevert()
	p := binding.NewTransparentUpgradeableProxy(receipt.ContractAddress, signer.Client)
	p.Contract().SetFrom(signer.Address())
	return p
}

func DeployAddressManagerLogic(signer *contract.Signer) *binding.AddressManager {
	receipt := binding.DeployAddressManager(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy address manager, address:", receipt.ContractAddress.String())
	addrMan := binding.NewAddressManager(receipt.ContractAddress, signer.Client)
	addrMan.Contract().SetFrom(signer.Address())
	return addrMan
}

func DeployL1StandardBridgeLogic(signer *contract.Signer) *binding.L1StandardBridge {
	receipt := binding.DeployL1StandardBridge(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l1 standard bridge, address:", receipt.ContractAddress.String())

	bridge := binding.NewL1StandardBridge(receipt.ContractAddress, signer.Client)
	bridge.Contract().SetFrom(signer.Address())
	return bridge
}

// TODO: using proxy
func DeployL1Contracts(signer *contract.Signer, cfg *config.L1ChainDeployConfig) *L1Contracts {
	// deploy address manager
	addrLogic := DeployAddressManagerLogic(signer)
	addrProxy := DeployUpgradeProxy(signer, addrLogic.Contract().Addr(), cfg.Admin, addrLogic.Initialize().Data)
	resolver := addrProxy.Contract().Addr()
	l1CrossLayerWitnessLogic := DeployL1CrossLayerWitnessLogic(signer)
	l1CrossLayerWitnessProxy := DeployUpgradeProxy(signer, l1CrossLayerWitnessLogic.Contract().Addr(), cfg.Admin, l1CrossLayerWitnessLogic.Initialize(resolver).Data)
	chainContainerLogic := DeployChainStorageLogic(signer)
	inputChainContainerProxy := DeployUpgradeProxy(signer, chainContainerLogic.Contract().Addr(), cfg.Admin, chainContainerLogic.Initialize("RollupInputChain", resolver).Data)
	stateChainContainerProxy := DeployUpgradeProxy(signer, chainContainerLogic.Contract().Addr(), cfg.Admin, chainContainerLogic.Initialize("RollupStateChain", resolver).Data)

	rollupInputChainLogic := DeployRollupInputChainLogic(signer)
	rollupInputChainProxy := DeployUpgradeProxy(signer, rollupInputChainLogic.Contract().Addr(), cfg.Admin, rollupInputChainLogic.Initialize(resolver, cfg.MaxEnqueueTxGasLimit, cfg.MaxWitnessTxExecGasLimit, cfg.L2ChainId).Data)
	rollupStateChainLogic := DeployRollupStateChainLogic(signer)
	rollupStateChainProxy := DeployUpgradeProxy(signer, rollupStateChainLogic.Contract().Addr(), cfg.Admin, rollupStateChainLogic.Initialize(resolver, big.NewInt(0).SetUint64(cfg.FraudProofWindow)).Data)

	var feeToken *binding.ERC20
	if cfg.FeeToken.IsZero() {
		feeToken = DeployTestFeeToken(signer)
	} else {
		feeToken = binding.NewERC20(cfg.FeeToken, signer.Client)
		feeToken.Contract().SetFrom(signer.Address())
	}

	daoLogic := DeployDAOLogic(signer)
	daoProxy := DeployUpgradeProxy(signer, daoLogic.Contract().Addr(), cfg.Admin, daoLogic.Initialize().Data)
	challenge := DeployChallengeLogic(signer)
	beacon := DeployBeacon(signer, challenge.Contract().Addr())
	factoryLogic := DeployChallengeFactoryLogic(signer)
	factoryProxy := DeployUpgradeProxy(signer, factoryLogic.Contract().Addr(), cfg.Admin, factoryLogic.Initialize(resolver, beacon.Contract().Addr(), new(big.Int).SetUint64(cfg.BlockLimitPerRound), cfg.ChallengerDeposit).Data)
	stakingLogic := DeployStakingManagerLogic(signer)
	stakingProxy := DeployUpgradeProxy(signer, stakingLogic.Contract().Addr(), cfg.Admin, stakingLogic.Initialize(daoProxy.Contract().Addr(), factoryProxy.Contract().Addr(),
		rollupStateChainProxy.Contract().Addr(), feeToken.Contract().Addr(), cfg.StakingAmount).Data)

	bridgeLogic := DeployL1StandardBridgeLogic(signer)
	bridgeProxy := DeployUpgradeProxy(signer, bridgeLogic.Contract().Addr(), cfg.Admin, bridgeLogic.Initialize(l1CrossLayerWitnessProxy.Contract().Addr(), cfg.L2StandardBridge).Data)

	names := []string{
		"L1CrossLayerWitness",
		"RollupInputChainContainer",
		"RollupStateChainContainer",
		"RollupInputChain",
		"RollupStateChain",
		"DAO",
		"StakingManager",
		"StakingManager",
		"ChallengeFactory",
		"L2CrossLayerWitness",
	}
	addrs := []web3.Address{
		l1CrossLayerWitnessProxy.Contract().Addr(),
		inputChainContainerProxy.Contract().Addr(),
		stateChainContainerProxy.Contract().Addr(),
		rollupInputChainProxy.Contract().Addr(),
		rollupStateChainProxy.Contract().Addr(),
		daoProxy.Contract().Addr(),
		stakingProxy.Contract().Addr(),
		stakingProxy.Contract().Addr(),
		factoryProxy.Contract().Addr(),
		cfg.L2CrossLayerWitness,
	}
	manager := binding.NewAddressManager(resolver, signer.Client)
	manager.Contract().SetFrom(signer.Address())
	manager.SetAddressBatch(names, addrs).Sign(signer).SendTransaction(signer).EnsureNoRevert()

	l1 := &L1Contracts{
		AddressManager:      manager,
		InputChainStorage:   binding.NewChainStorageContainer(inputChainContainerProxy.Contract().Addr(), signer.Client),
		StateChainStorage:   binding.NewChainStorageContainer(stateChainContainerProxy.Contract().Addr(), signer.Client),
		RollupInputChain:    binding.NewRollupInputChain(rollupInputChainProxy.Contract().Addr(), signer.Client),
		RollupStateChain:    binding.NewRollupStateChain(rollupStateChainProxy.Contract().Addr(), signer.Client),
		L1CrossLayerWitness: binding.NewL1CrossLayerWitness(l1CrossLayerWitnessProxy.Contract().Addr(), signer.Client),
		L1StandardBridge:    binding.NewL1StandardBridge(bridgeProxy.Contract().Addr(), signer.Client),
		FeeToken:            feeToken,
		ChallengeLogic:      challenge,
		ChallengeBeacon:     beacon,
		ChallengeFactory:    binding.NewChallengeFactory(factoryProxy.Contract().Addr(), signer.Client),
		StakingManager:      binding.NewStakingManager(stakingProxy.Contract().Addr(), signer.Client),
		DAO:                 binding.NewDAO(daoProxy.Contract().Addr(), signer.Client),
	}
	v := reflect.ValueOf(l1).Elem()
	for i := 0; i < v.NumField(); i++ {
		v.Field(i).Interface().(Bind).Contract().SetFrom(signer.Address())
	}
	return l1
}

type Bind interface {
	Contract() *contract.Contract
}
