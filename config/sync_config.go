package config

import (
	"github.com/laizy/web3"
)

type SyncConfig struct {
	L1RpcUrl           string
	DbDir              string
	StartSyncHeight    uint64 // system contract deployed on l1 height
	MinConfirmBlockNum uint64

	RollupInputChain web3.Address
	RollupStateChain web3.Address
	AddressManager   web3.Address
	Dao              web3.Address

	L1TokenBridge web3.Address
	L2TokenBridge web3.Address

	L1CrossLayerWitness web3.Address
	L2CrossLayerWitness web3.Address
}
