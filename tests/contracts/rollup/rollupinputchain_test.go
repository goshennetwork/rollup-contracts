package rollup

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/goshennetwork/rollup-contracts/binding"
	"github.com/goshennetwork/rollup-contracts/deploy"
	"github.com/goshennetwork/rollup-contracts/tests/contracts"
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
)

func EnqueueTransactionHash(sender, target web3.Address, gasLimit uint64, data []byte, nonce uint64) web3.Hash {
	key := contracts.LocalL1ChainEnv.PrivKey
	gasPrice := uint64(GasPrice)
	if sender == L1CrossLayerFakeSender { //fix if just provided index in L1CrossLayer Contract
		if nonce < contracts.INIT_ENQUEUE_NONCE {
			nonce += contracts.INIT_ENQUEUE_NONCE
		}
		key = L1CrossLayerFakeKey
		fmt.Println(nonce)
		gasPrice = 0
	}
	txdata := CompleteTxData(target, gasPrice, gasLimit, data, nonce)
	r, s, v := Sign(target, gasPrice, gasLimit, data, nonce, key)
	txdata.V = v
	txdata.R = r
	txdata.S = s
	_s, err := types.NewEIP155Signer(new(big.Int).SetUint64(contracts.LocalL1ChainEnv.ChainConfig.L2ChainId)).Sender(types.NewTx(txdata))
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
	l1Chain := deploy.DeployL1Contracts(signer, chainEnv.ChainConfig)
	nonce, err := l1Chain.RollupInputChain.GetNonceByAddress(signer.Address())
	utils.Ensure(err)
	target, gasLimit, data := web3.Address{1, 1}, uint64(900_000), []byte("test")
	r, s, v := Sign(target, GasPrice, gasLimit, data, nonce, contracts.LocalL1ChainEnv.PrivKey)
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
	nonce, err = l1Chain.RollupInputChain.GetNonceByAddress(L1CrossLayerFakeSender)
	utils.Ensure(err)
	gasLimit = chainEnv.ChainConfig.MaxWitnessTxExecGasLimit
	utils.EnsureTrue(txHash == EnqueueTransactionHash(L1CrossLayerFakeSender, chainEnv.ChainConfig.L2CrossLayerWitness,
		gasLimit, crossLayerMsg, nonce-1))
}

func TestAppendBatches(t *testing.T) {
	chainEnv := contracts.LocalL1ChainEnv
	signer := contracts.SetupLocalSigner(chainEnv.ChainId, chainEnv.PrivKey)
	l1Chain := deploy.DeployL1Contracts(signer, chainEnv.ChainConfig)

	l1Chain.FeeToken.Approve(l1Chain.StakingManager.Contract().Addr(), chainEnv.ChainConfig.StakingAmount).Sign(signer).SendTransaction(signer)
	l1Chain.StakingManager.Deposit().Sign(signer).SendTransaction(signer)
	l1Chain.Whitelist.SetSequencer(signer.Address(), true).Sign(signer).SendTransaction(signer)

	batches := &binding.RollupInputBatches{
		BatchIndex: 0,
		QueueNum:   0,
		QueueStart: 0,
		SubBatches: []*binding.SubBatch{
			{
				0,
				nil,
			},
		},
	}
	receipt := l1Chain.RollupInputChain.AppendInputBatches(batches).Sign(signer).SendTransaction(signer)

	fmt.Println(utils.JsonString(receipt))

	//utils.EnsureTrue(strings.Contains(utils.JsonStr(receipt), "TransactionAppended"))
	height, err := l1Chain.RollupInputChain.ChainHeight()
	utils.Ensure(err)
	utils.EnsureTrue(height == uint64(1))
}
