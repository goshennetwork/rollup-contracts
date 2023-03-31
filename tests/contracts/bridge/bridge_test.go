package bridge

import (
	"fmt"
	"testing"

	"github.com/goshennetwork/rollup-contracts/deploy"
	"github.com/goshennetwork/rollup-contracts/tests/contracts"
)

func TestDepositEth(t *testing.T) {

	chainEnv := contracts.LocalL1ChainEnv
	signer := contracts.SetupLocalSigner(chainEnv.ChainId, chainEnv.PrivKey)
	l1Chain := deploy.DeployL1Contracts(signer, chainEnv.ChainConfig)

	fmt.Println(l1Chain)
}
