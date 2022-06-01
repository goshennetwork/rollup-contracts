package contracts

import (
	"fmt"
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
)

var DevL2ChainEnv = &deploy.L2ChainEnv{
	ChainId: 0x539,
	//	21772,
	RpcUrl:      "http://172.24.78.20:12345/",
	PrivKey:     "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	ChainConfig: &deploy.L2ChainDeployConfig{},
}

func TestTransfer(t *testing.T) {
	chainEnv := DevL2ChainEnv
	client, err := jsonrpc.NewClient(chainEnv.RpcUrl)
	utils.Ensure(err)
	signer := contract.NewSigner(chainEnv.PrivKey, client, chainEnv.ChainId)
	signer.Submit = true

	fmt.Println(client.Eth().GetBalance(signer.Address(), 0))

	signed := signer.TransferEther(web3.HexToAddress("0xd7804Ab82801DF2a4499C3157d907CE0C11e1bac"), web3.Ether(1), "hello")
	signer.SendTransaction(signed)
}
