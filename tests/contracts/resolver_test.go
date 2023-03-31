package contracts

import (
	"testing"

	"github.com/goshennetwork/rollup-contracts/deploy"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils"
	"github.com/stretchr/testify/assert"
)

func TestResolver(t *testing.T) {
	signer := SetupLocalSigner(LocalL1ChainEnv.ChainId, LocalL1ChainEnv.PrivKey)
	addrMan := deploy.DeployL1Contracts(signer, LocalL1ChainEnv.ChainConfig).AddressManager

	receipt := addrMan.SetAddress("signer", signer.Address()).Sign(signer).SendTransaction(signer)
	assert.Equal(t, receipt.Status, uint64(1))
	signerAddr, err := addrMan.GetAddr("signer")
	utils.Ensure(err)
	assert.Equal(t, signerAddr, signer.Address())
	updateAddr := crypto.CreateAddress(signer.Address(), 10)
	receipt = addrMan.SetAddress("signer", updateAddr).Sign(signer).SendTransaction(signer)
	assert.Equal(t, receipt.Status, uint64(1))
	signerAddr, err = addrMan.GetAddr("signer")
	utils.Ensure(err)
	assert.Equal(t, signerAddr, updateAddr)
}
