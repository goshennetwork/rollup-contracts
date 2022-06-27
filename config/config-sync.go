package config

import (
	"math/big"

	"github.com/laizy/web3"
)

const (
	ADDRESS_MANAGER = "AddressManager"

	///DAO
	DAO = "DAO"
	///RollupInputChain
	ROLLUP_INPUT_CHAIN = "RollupInputChain"
	///ChainStorageContainer of RollupInputChain
	ROLLUP_INPUT_CHAIN_CONTAINER = "RollupInputChainContainer"
	///RollupStateChain
	ROLLUP_STATE_CHAIN = "RollupStateChain"
	///ChainStorageContainer of RollupStateChain
	ROLLUP_STATE_CHAIN_CONTAINER = "RollupStateChainContainer"
	///StakingManager
	STAKING_MANAGER = "StakingManager"
	///ChallengeFactory
	CHALLENGE_FACTORY = "ChallengeFactory"
	///L1CrossLayerWitness
	L1_CROSS_LAYER_WITNESS = "L1CrossLayerWitness"
	///L2CrossLayerWitness
	L2_CROSS_LAYER_WITNESS = "L2CrossLayerWitness"
)

const (
	DefaultDeployConfigName = "deploy-l1-config.json"
	DefaultSyncConfigName   = "contracts-sync-config.json"
	DefaultSyncDbName       = "sync-db"
	DefaultL1MMRFile        = "l1tree.db"
	DefaultL2MMRFile        = "l2tree.db"
	DefaultContractName     = "contracts.json"
)

type DeployConfig struct {
	L1Client   string
	PrivateKey string
	DAO        web3.Address
	*RollupInputChainConfig
	*RollupStateChainConfig
	*StakingManagerConfig
	*Erc20Config
}

type Erc20Config struct {
	Name   string
	Symbol string
}
type RollupInputChainConfig struct {
	MaxTxGasLimit           uint64
	MaxCrossLayerTxGasLimit uint64
	L2ChainId               uint64
}

type RollupStateChainConfig struct {
	FraudProofWindow *big.Int
}

type StakingManagerConfig struct {
	Price *big.Int
}

type Contracts struct {
	RollupInputChain    web3.Address
	RollupStateChain    web3.Address
	StakingManager      web3.Address
	AddressManager      web3.Address
	L1CrossLayerWitness web3.Address
	Dao                 web3.Address
}

type RollupConfig struct {
	SyncConfig SyncConfig
	Contracts  Contracts
}
