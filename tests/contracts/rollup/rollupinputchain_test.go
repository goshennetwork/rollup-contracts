package rollup

import (
	"fmt"
	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/laizy/web3/utils/common"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"
	"strings"
	"testing"
)

var L1CrossLayerWitnessAddr = web3.HexToAddress("0x5800000000000000000000000000000000000000")

func EnqueueTransactionHash(sender, target web3.Address, gasLimit uint64, data []byte) web3.Hash {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteAddress(sender)
	sink.WriteAddress(target)
	sink.WriteUint64BE(gasLimit)
	sink.WriteBytes(data)

	return crypto.Keccak256Hash(sink.Bytes())
}

func TestEnqueue(t *testing.T) {
	chainEnv := contracts.LocalChainEnv
	signer := contracts.SetupLocalSigner(chainEnv)
	l1Chain := contracts.DeployL1Contract(signer, chainEnv.L1ChainConfig)

	target, gasLimit, data := web3.Address{1, 1}, uint64(900_000), []byte("test")
	receipt := l1Chain.RollupInputChain.Enqueue(target, gasLimit, data).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(receipt.Status == 1)

	txHash, _, err := l1Chain.RollupInputChain.GetQueueTxInfo(0)
	utils.Ensure(err)
	utils.EnsureTrue(txHash == EnqueueTransactionHash(signer.Address(), target, gasLimit, data))

	receipt = l1Chain.L1CrossLayerWitness.SendMessage(target, data).Sign(signer).SendTransaction(signer)
	utils.EnsureTrue(strings.Contains(utils.JsonStr(receipt), "MessageSent" ))

	msgHash := contracts.CrossLayerMessageHash(target, signer.Address(), 0, data)
	mmrRoot, err := l1Chain.L1CrossLayerWitness.MmrRoot()
	utils.Ensure(err)
	utils.EnsureTrue(mmrRoot == msgHash)
	crossLayerMsg := contracts.EncodeL1ToL2CallData(target, signer.Address(), data, 0, mmrRoot, 1)

	txHash, _, err = l1Chain.RollupInputChain.GetQueueTxInfo(1)
	utils.Ensure(err)
	utils.EnsureTrue(txHash == EnqueueTransactionHash(L1CrossLayerWitnessAddr, chainEnv.L1ChainConfig.L2CrossLayerWitness, chainEnv.L1ChainConfig.MaxCrossLayerTxGasLimit , crossLayerMsg))
}

func TestAppendBatches(t *testing.T) {
	chainEnv := contracts.LocalChainEnv
	signer := contracts.SetupLocalSigner(chainEnv)
	l1Chain := contracts.DeployL1Contract(signer, chainEnv.L1ChainConfig)

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
	utils.EnsureTrue(height== uint64(1))
}
