package contracts

import (
	"fmt"
	"github.com/laizy/web3/utils"

	"github.com/laizy/web3/contract"
	"github.com/ontology-layer-2/rollup-contracts/binding"
)

type ChainEnv struct {
	ChainId           uint64
	RpcUrl            string
	PrivKey           string
}

type L1Contracts struct {
	AddressManager *binding.AddressManager
	InputChainStorage *binding.ChainStorageContainer
	StateChainStorage *binding.ChainStorageContainer
	RollupInputChain *binding.ChainStorageContainer
	RollupStateChain *binding.ChainStorageContainer
}

func DeployL1Contract(signer *contract.Signer) *L1Contracts {
	client := signer.Client
	// deploy address manager
	receipt := binding.DeployAddressManager(client, signer.Address()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	fmt.Println("deploy address manager, address:", receipt.ContractAddress.String())

	addrMan := binding.NewAddressManager(receipt.ContractAddress, signer.Client)
	addrMan.Contract().SetFrom(signer.Address())

	receipt = binding.DeployChainStorageContainer(client, signer.Address(), "RollupInputChain", addrMan.Contract().Addr()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	inputChainContainer := binding.NewChainStorageContainer(receipt.ContractAddress, client)
	inputChainContainer.Contract().SetFrom(signer.Address())
	fmt.Println("deploy input chain storage, address:", receipt.ContractAddress.String())
	receipt = binding.DeployChainStorageContainer(client, signer.Address(), "RollupStateChain", addrMan.Contract().Addr()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	stateChainContainer := binding.NewChainStorageContainer(receipt.ContractAddress, client)
	stateChainContainer.Contract().SetFrom(signer.Address())
	fmt.Println("deploy state chain storage, address:", receipt.ContractAddress.String())

	addrMan.NewAddr("RollupInputChainContainer", inputChainContainer.Contract().Addr())
	addrMan.NewAddr("RollupStateChainContainer", stateChainContainer.Contract().Addr())

	receipt = binding.DeployRollupInputChain(client, signer.Address(), "RollupInputChain", addrMan.Contract().Addr()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)
	inputChainContainer := binding.NewChainStorageContainer(receipt.ContractAddress, client)
	inputChainContainer.Contract().SetFrom(signer.Address())

	return &L1Contracts{
		AddressManager: addrMan,
		InputChainStorage: inputChainContainer,
		StateChainStorage: stateChainContainer,
	}
}