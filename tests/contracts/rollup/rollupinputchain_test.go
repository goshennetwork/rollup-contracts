package rollup

import (
	"bytes"
	"fmt"
	"github.com/ontology-layer-2/rollup-contracts/deploy"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/common"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
)

func EnqueueTransactionHash(sender, target web3.Address, gasLimit uint64, data []byte, nonce uint64) web3.Hash {
	key := contracts.LocalL1ChainEnv.PrivKey
	if sender == L1CrossLayerFakeSender {
		key = L1CrossLayerFakeKey
		fmt.Println(nonce)
	}
	txdata := CompleteTxData(target, gasLimit, data, nonce)
	r, s, v := Sign(target, gasLimit, data, nonce, key)
	txdata.V = v
	txdata.R = r
	txdata.S = s
	_s, err := types.NewEIP155Signer(new(big.Int).SetUint64(contracts.LocalL1ChainEnv.L1ChainConfig.L2ChainId)).Sender(types.NewTx(txdata))
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(sender.Bytes(), _s.Bytes()) {
		panic("wrong sig")
	}
	return web3.Hash(types.NewTx(txdata).Hash())
}

func TestEnqueue(t *testing.T) {
	chainEnv := contracts.LocalL1ChainEnv
	signer := contracts.SetupLocalSigner(chainEnv.ChainId, chainEnv.PrivKey)
	l1Chain := deploy.DeployL1Contract(signer, chainEnv.L1ChainConfig)

	target, gasLimit, data, nonce := web3.Address{1, 1}, uint64(900_000), []byte("test"), uint64(0)
	r, s, v := Sign(target, gasLimit, data, nonce, contracts.LocalL1ChainEnv.PrivKey)
	receipt := l1Chain.RollupInputChain.Enqueue(target, gasLimit, data, nonce, r, s, v.Uint64()).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)

	txHash, _, err := l1Chain.RollupInputChain.GetQueueTxInfo(0)
	utils.Ensure(err)
	utils.EnsureTrue(txHash == EnqueueTransactionHash(signer.Address(), target, gasLimit, data, nonce))

	receipt = l1Chain.L1CrossLayerWitness.SendMessage(target, data).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(strings.Contains(utils.JsonStr(receipt), "MessageSent"))

	msgHash := contracts.CrossLayerMessageHash(target, signer.Address(), 0, data)
	mmrRoot, err := l1Chain.L1CrossLayerWitness.MmrRoot()
	utils.Ensure(err)
	utils.EnsureTrue(mmrRoot == msgHash)
	crossLayerMsg := contracts.EncodeL1ToL2CallData(target, signer.Address(), data, 0, mmrRoot, 1)

	txHash, _, err = l1Chain.RollupInputChain.GetQueueTxInfo(1)
	utils.Ensure(err)
	size, err := l1Chain.L1CrossLayerWitness.TotalSize()
	utils.Ensure(err)
	utils.EnsureTrue(txHash == EnqueueTransactionHash(L1CrossLayerFakeSender, chainEnv.L1ChainConfig.L2CrossLayerWitness, chainEnv.L1ChainConfig.MaxCrossLayerTxGasLimit, crossLayerMsg, size-1))
}

func TestAppendBatches(t *testing.T) {
	chainEnv := contracts.LocalL1ChainEnv
	signer := contracts.SetupLocalSigner(chainEnv.ChainId, chainEnv.PrivKey)
	l1Chain := deploy.DeployL1Contract(signer, chainEnv.L1ChainConfig)

	l1Chain.FeeToken.Approve(l1Chain.StakingManager.Contract().Addr(), chainEnv.L1ChainConfig.StakingAmount).Sign(signer).SendTransaction(signer)
	l1Chain.StakingManager.Deposit().Sign(signer).SendTransaction(signer)
	l1Chain.DAO.SetSequencerWhitelist(signer.Address(), true).Sign(signer).SendTransaction(signer)

	batches := &binding.RollupInputBatches{
		QueueNum:    0,
		QueueStart:  0,
		BatchNum:    1,
		Batch0Time:  uint64(1),
		BatchesData: common.Hash{}.Bytes(),
	}
	receipt := l1Chain.RollupInputChain.AppendInputBatches(batches).Sign(signer).SendTransaction(signer)

	fmt.Println(utils.JsonString(receipt))

	//utils.EnsureTrue(strings.Contains(utils.JsonStr(receipt), "TransactionAppended"))
	height, err := l1Chain.RollupInputChain.ChainHeight()
	utils.Ensure(err)
	utils.EnsureTrue(height == uint64(1))
}
