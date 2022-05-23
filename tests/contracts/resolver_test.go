package contracts

import (
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/jsonrpc/transport"
	"github.com/laizy/web3/utils"
	"gotest.tools/assert"
)

var LocalChainEnv = &ChainEnv{
	ChainId: 1234,
	RpcUrl:  "local",
	PrivKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
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

func TestResolver(t *testing.T) {
	signer := SetupLocalSigner(LocalChainEnv)

	addrMan:= DeployL1Contract(signer).AddressManager

	receipt := addrMan.NewAddr("signer", signer.Address()).Sign(signer).SendTransaction(signer)
	assert.Equal(t, receipt.Status, uint64(1))
	signerAddr, err := addrMan.GetAddr("signer")
	utils.Ensure(err)
	assert.Equal(t, signerAddr, signer.Address())
	updateAddr := crypto.CreateAddress(signer.Address(), 10)
	receipt = addrMan.UpdateAddr("signer", updateAddr).Sign(signer).SendTransaction(signer)
	assert.Equal(t, receipt.Status, uint64(1))
	signerAddr, err = addrMan.GetAddr("signer")
	utils.Ensure(err)
	assert.Equal(t, signerAddr, updateAddr)
}