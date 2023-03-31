package rollup

import (
	"github.com/goshennetwork/rollup-contracts/store/schema"
	"github.com/laizy/web3"
)

func genL1DepositKey(l1TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L1TokenBridgeDepositKey, l1TxHash)
}

func genL1WithdrawalKey(l1TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L1TokenBridgeWithdrawalKey, l1TxHash)
}

func genL2WithdrawalInitKey(l2TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L2TokenBridgeWithdrawalKey, l2TxHash)
}

func genDepositFinalizedKey(l2TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L2TokenBridgeDepositFinalizedKey, l2TxHash)
}

func genDepositFailedKey(l2TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L2TokenBridgeDepositFailedKey, l2TxHash)
}

// use tx hash as key could prevent mismatch when block re-org
func genKeyByTxHash(prefix byte, txHash web3.Hash) []byte {
	key := make([]byte, 0)
	key = append(key, prefix)
	key = append(key, txHash.Bytes()...)
	return key
}
