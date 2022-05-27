package contracts

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/jsonrpc/transport"
)

var LocalChainEnv = &ChainEnv{
	ChainId: 1,
	RpcUrl:  "local",
	PrivKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	L1ChainConfig: &L1ChainDeployConfig{
		FraudProofWindow:        3,
		MaxEnqueueTxGasLimit:    15000000,
		MaxCrossLayerTxGasLimit: 5000000,
		StakingAmount:           web3.Ether(10),
		L2CrossLayerWitness:     web3.Address{1, 2, 3, 4, 5, 6},
		L2ChainId:               1234,
		ChallengeConfig: &ChallengeConfig{
			BlockLimitPerRound: 10,
			ChallengerDeposit:  web3.Ether(1),
		},
	},
}

func SetupLocalSigner(chainEnv *ChainEnv) *contract.Signer {
	db := storage.NewFakeDB()
	local := transport.NewLocal(db, chainEnv.ChainId)
	client := jsonrpc.NewClientWithTransport(local)
	signer := contract.NewSigner(chainEnv.PrivKey, client, chainEnv.ChainId)
	signer.Submit = true
	local.SetBalance(signer.Address(), web3.Ether(1000))

	return signer
}
