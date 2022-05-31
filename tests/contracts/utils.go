package contracts

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/jsonrpc/transport"
	"github.com/ontology-layer-2/rollup-contracts/config"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
)

var LocalL2ChainEnv = &deploy.L2ChainEnv{
	ChainId:       1234,
	RpcUrl:        "local",
	PrivKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	L2ChainConfig: &deploy.L2ChainDeployConfig{},
}

var LocalL1ChainEnv = &deploy.L1ChainEnv{
	ChainId: 1,
	RpcUrl:  "local",
	PrivKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	L1ChainConfig: &config.L1ChainDeployConfig{
		FraudProofWindow:        3,
		MaxEnqueueTxGasLimit:    15000000,
		MaxCrossLayerTxGasLimit: 5000000,
		StakingAmount:           web3.Ether(10),
		L2CrossLayerWitness:     web3.Address{1, 2, 3, 4, 5, 6},
		L2ChainId:               LocalL2ChainEnv.ChainId,
		ChallengeConfig: &config.ChallengeConfig{
			BlockLimitPerRound: 10,
			ChallengerDeposit:  web3.Ether(1),
		},
	},
}

func SetupLocalSigner(chainID uint64, privKey string) *contract.Signer {
	db := storage.NewFakeDB()
	local := transport.NewLocal(db, chainID)
	client := jsonrpc.NewClientWithTransport(local)
	signer := contract.NewSigner(privKey, client, chainID)
	signer.Submit = true
	local.SetBalance(signer.Address(), web3.Ether(1000))

	return signer
}