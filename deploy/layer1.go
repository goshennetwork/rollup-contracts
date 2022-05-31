package deploy

import (
	"fmt"
	"math/big"

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

func DeployDAO(signer *contract.Signer) *binding.DAO {
	receipt := binding.DeployDAO(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	dao := binding.NewDAO(receipt.ContractAddress, signer.Client)
	dao.Contract().SetFrom(signer.Address())
	dao.Initialize().Sign(signer).SendTransaction(signer)

	return dao
}

func DeployChallengeFactory(signer *contract.Signer, addrMan, beacon web3.Address, blockLimitPerRound uint64, challengerDeposit *big.Int) *binding.ChallengeFactory {
	receipt := binding.DeployChallengeFactory(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	factory := binding.NewChallengeFactory(receipt.ContractAddress, signer.Client)
	factory.Contract().SetFrom(signer.Address())
	factory.Initialize(addrMan, beacon, big.NewInt(0).SetUint64(blockLimitPerRound), challengerDeposit).Sign(signer).SendTransaction(signer)

	return factory
}

func DeployStakingManager(signer *contract.Signer, dao, challengeFactory, rollupStateChain,
	feeToken web3.Address, price *big.Int) *binding.StakingManager {
	receipt := binding.DeployStakingManager(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	staking := binding.NewStakingManager(receipt.ContractAddress, signer.Client)
	staking.Contract().SetFrom(signer.Address())
	staking.Initialize(dao, challengeFactory, rollupStateChain, feeToken, price).Sign(signer).SendTransaction(signer)

	return staking
}

func DeployRollupInputChain(signer *contract.Signer, addrMan web3.Address, maxEnqueueTxGasLimit,
	maxCrossLayerTxGasLimit uint64, l2ChainID uint64) *binding.RollupInputChain {
	receipt := binding.DeployRollupInputChain(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	rollupInputChain := binding.NewRollupInputChain(receipt.ContractAddress, signer.Client)
	rollupInputChain.Contract().SetFrom(signer.Address())
	rollupInputChain.Initialize(addrMan, maxEnqueueTxGasLimit, maxCrossLayerTxGasLimit, l2ChainID).Sign(signer).SendTransaction(signer)

	return rollupInputChain
}

func DeployRollupStateChain(signer *contract.Signer, addrMan web3.Address, fraudProofWindow uint64) *binding.RollupStateChain {
	receipt := binding.DeployRollupStateChain(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	rollupStateChain := binding.NewRollupStateChain(receipt.ContractAddress, signer.Client)
	rollupStateChain.Contract().SetFrom(signer.Address())
	rollupStateChain.Initialize(addrMan, big.NewInt(0).SetUint64(fraudProofWindow)).Sign(signer).SendTransaction(signer)

	return rollupStateChain
}

func DeployChainStorage(signer *contract.Signer, addrMan web3.Address, owner string) *binding.ChainStorageContainer {
	receipt := binding.DeployChainStorageContainer(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Printf("deploy chain storage, owner: %s, address:%s\n", owner, receipt.ContractAddress.String())
	chainStorage := binding.NewChainStorageContainer(receipt.ContractAddress, signer.Client)
	chainStorage.Contract().SetFrom(signer.Address())
	chainStorage.Initialize(owner, addrMan).Sign(signer).SendTransaction(signer)
	fmt.Println("initialized chain storage")

	return chainStorage
}

func DeployL1CrossLayerWitness(signer *contract.Signer, addrMan web3.Address) *binding.L1CrossLayerWitness {
	receipt := binding.DeployL1CrossLayerWitness(signer.Client, signer.Address()).
		Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l1 cross layer witness, address:", receipt.ContractAddress.String())
	l1CrossLayerWitness := binding.NewL1CrossLayerWitness(receipt.ContractAddress, signer.Client)
	l1CrossLayerWitness.Contract().SetFrom(signer.Address())
	l1CrossLayerWitness.Initialize(addrMan).Sign(signer).SendTransaction(signer)

	return l1CrossLayerWitness
}

func DeployAddressManager(signer *contract.Signer) *binding.AddressManager {
	receipt := binding.DeployAddressManager(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy address manager, address:", receipt.ContractAddress.String())

	addrMan := binding.NewAddressManager(receipt.ContractAddress, signer.Client)
	addrMan.Contract().SetFrom(signer.Address())
	utils.EnsureTrue(1 == addrMan.Initialize().Sign(signer).SendTransaction(signer).Status)

	return addrMan
}

func DeployL1StandardBridge(signer *contract.Signer, l1witness, l2bridge web3.Address) *binding.L1StandardBridge {
	receipt := binding.DeployL1StandardBridge(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy l1 standard bridge, address:", receipt.ContractAddress.String())

	bridge := binding.NewL1StandardBridge(receipt.ContractAddress, signer.Client)
	bridge.Contract().SetFrom(signer.Address())

	bridge.Initialize(l1witness, l2bridge).Sign(signer).SendTransaction(signer)

	return bridge
}

// TODO: using proxy
func DeployL1Contracts(signer *contract.Signer, cfg *config.L1ChainDeployConfig) *L1Contracts {
	// deploy address manager
	addrMan := DeployAddressManager(signer)
	l1CrossLayerWitness := DeployL1CrossLayerWitness(signer, addrMan.Contract().Addr())
	inputChainContainer := DeployChainStorage(signer, addrMan.Contract().Addr(), "RollupInputChain")
	stateChainContainer := DeployChainStorage(signer, addrMan.Contract().Addr(), "RollupStateChain")

	rollupInputChain := DeployRollupInputChain(signer, addrMan.Contract().Addr(), cfg.MaxEnqueueTxGasLimit, cfg.MaxCrossLayerTxGasLimit, cfg.L2ChainId)
	rollupStateChain := DeployRollupStateChain(signer, addrMan.Contract().Addr(), cfg.FraudProofWindow)

	var feeToken *binding.ERC20
	if cfg.FeeToken.IsZero() {
		feeToken = DeployTestFeeToken(signer)
	} else {
		feeToken = binding.NewERC20(cfg.FeeToken, signer.Client)
		feeToken.Contract().SetFrom(signer.Address())
	}

	dao := DeployDAO(signer)
	challenge := DeployChallengeLogic(signer)
	beacon := DeployBeacon(signer, challenge.Contract().Addr())
	factory := DeployChallengeFactory(signer, addrMan.Contract().Addr(), beacon.Contract().Addr(), cfg.BlockLimitPerRound, cfg.ChallengerDeposit)

	staking := DeployStakingManager(signer, dao.Contract().Addr(), factory.Contract().Addr(),
		rollupStateChain.Contract().Addr(), feeToken.Contract().Addr(), cfg.StakingAmount)

	bridge := DeployL1StandardBridge(signer, l1CrossLayerWitness.Contract().Addr(), cfg.L2StandardBridge)

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
		l1CrossLayerWitness.Contract().Addr(),
		inputChainContainer.Contract().Addr(),
		stateChainContainer.Contract().Addr(),
		rollupInputChain.Contract().Addr(),
		rollupStateChain.Contract().Addr(),
		dao.Contract().Addr(),
		staking.Contract().Addr(),
		staking.Contract().Addr(),
		factory.Contract().Addr(),
		cfg.L2CrossLayerWitness,
	}
	addrMan.SetAddressBatch(names, addrs).Sign(signer).SendTransaction(signer)

	return &L1Contracts{
		AddressManager:      addrMan,
		InputChainStorage:   inputChainContainer,
		StateChainStorage:   stateChainContainer,
		RollupInputChain:    rollupInputChain,
		RollupStateChain:    rollupStateChain,
		L1CrossLayerWitness: l1CrossLayerWitness,
		L1StandardBridge:    bridge,
		FeeToken:            feeToken,
		ChallengeLogic:      challenge,
		ChallengeBeacon:     beacon,
		ChallengeFactory:    factory,
		StakingManager:      staking,
		DAO:                 dao,
	}
}
