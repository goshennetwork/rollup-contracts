package contracts

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
)

type ChainEnv struct {
	ChainId           uint64
	RpcUrl            string
	PrivKey           string
	L1ChainConfig *L1ChainConfig
}

type L1ChainConfig struct {
	FeeToken web3.Address
	FraudProofWindow uint64 // block number
	MaxEnqueueTxGasLimit uint64
	MaxCrossLayerTxGasLimit uint64
}

type L1Contracts struct {
	AddressManager *binding.AddressManager
	InputChainStorage *binding.ChainStorageContainer
	StateChainStorage *binding.ChainStorageContainer
	RollupInputChain *binding.RollupInputChain
	RollupStateChain *binding.RollupStateChain
	L1CrossLayerWitness *binding.L1CrossLayerWitness
	StakingManager *binding.StakingManager
	ChallengeBeacon *binding.
	FeeToken *binding.ERC20
}

func DeployTestFeeToken(signer *contract.Signer) *binding.ERC20 {
	receipt := binding.DeployTestERC20(signer.Client, signer.Address(), "TestFeeToken", "TFT").Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)

	feeToken := binding.NewERC20(receipt.ContractAddress, signer.Client)
	feeToken.Contract().SetFrom(signer.Address())

	return feeToken
}

func DeployStakingManager(signer *contract.Signer) *binding.StakingManager {
	receipt := binding.DeployStakingManager(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	staking := binding.NewStakingManager(receipt.ContractAddress, signer.Client)
	staking.Contract().SetFrom(signer.Address())
	staking.Initialize()
}

func DeployRollupInputChain(signer *contract.Signer, addrMan web3.Address, maxEnqueueTxGasLimit,
	maxCrossLayerTxGasLimit uint64 )*binding.RollupInputChain {
	receipt := binding.DeployRollupInputChain(signer.Client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	rollupInputChain := binding.NewRollupInputChain(receipt.ContractAddress, signer.Client)
	rollupInputChain.Contract().SetFrom(signer.Address())
	rollupInputChain.Initialize(addrMan, maxEnqueueTxGasLimit, maxCrossLayerTxGasLimit).Sign(signer).SendTransaction(signer)

	return rollupInputChain
}

func DeployRollupStateChain(signer *contract.Signer, addrMan web3.Address, fraudProofWindow uint64) *binding.RollupStateChain {
	receipt := binding.DeployRollupStateChain(signer.Client, signer.Address(), ).Sign(signer).SendTransaction(signer)
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

// TODO: using proxy
func DeployL1Contract(signer *contract.Signer, cfg *L1ChainConfig) *L1Contracts {
	client := signer.Client
	// deploy address manager
	receipt := binding.DeployAddressManager(client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy address manager, address:", receipt.ContractAddress.String())

	addrMan := binding.NewAddressManager(receipt.ContractAddress, signer.Client)
	addrMan.Contract().SetFrom(signer.Address())
	utils.EnsureTrue(1 ==addrMan.Initialize().Sign(signer).SendTransaction(signer).Status)

	l1CrossLayerWitness := DeployL1CrossLayerWitness(signer, addrMan.Contract().Addr())
	addrMan.SetAddress("L1CrossLayerWitness", l1CrossLayerWitness.Contract().Addr())

	inputChainContainer := DeployChainStorage(signer, addrMan.Contract().Addr(), "RollupInputChain")
	stateChainContainer := DeployChainStorage(signer, addrMan.Contract().Addr(), "RollupStateChain")
	addrMan.SetAddress("RollupInputChainContainer", inputChainContainer.Contract().Addr())
	addrMan.SetAddress("RollupStateChainContainer", stateChainContainer.Contract().Addr())

	rollupInputChain := DeployRollupInputChain(signer, addrMan.Contract().Addr(), cfg.MaxEnqueueTxGasLimit,
		cfg.MaxCrossLayerTxGasLimit)
	rollupStateChain := DeployRollupStateChain(signer, addrMan.Contract().Addr(), cfg.FraudProofWindow)
	addrMan.SetAddress("RollupInputChain", rollupInputChain.Contract().Addr())
	addrMan.SetAddress("RollupStateChain", rollupStateChain.Contract().Addr())

	var feeToken *binding.ERC20
	if cfg.FeeToken.IsZero() {
		feeToken = DeployTestFeeToken(signer)
	} else {
		feeToken = binding.NewERC20(cfg.FeeToken, signer.Client)
		feeToken.Contract().SetFrom(signer.Address())
	}

	return &L1Contracts{
		AddressManager: addrMan,
		InputChainStorage: inputChainContainer,
		StateChainStorage: stateChainContainer,
		RollupInputChain: rollupInputChain,
		L1CrossLayerWitness: l1CrossLayerWitness,
		FeeToken: feeToken,
	}
}