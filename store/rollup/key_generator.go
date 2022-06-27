package rollup

import (
	"github.com/laizy/web3"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

func genL1ETHDepositKey(l1TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L1TokenBridgeETHDepositKey, l1TxHash)
}

func genL1ERC20DepositInitKey(l1TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L1TokenBridgeERC20DepositKey, l1TxHash)
}

func genL1ETHWithdrawalKey(l1TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L1TokenBridgeETHWithdrawalKey, l1TxHash)
}

func genL1ERC20WithdrawalFinalizedKey(l1TxHash web3.Hash) []byte {
	return genKeyByTxHash(schema.L1TokenBridgeERC20WithdrawalKey, l1TxHash)
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
