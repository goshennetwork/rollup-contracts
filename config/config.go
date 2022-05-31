package config

import (
	"math/big"

	"github.com/laizy/web3"
)

type RollupCliConfig struct {
	L1Rpc       string
	L2Rpc       string
	PrivKey     string
	L1Addresses *L1ContractAddressConfig
	L2Addresses *L2ContractAddressConfig
}

type L1ContractAddressConfig struct {
	AddressManager      web3.Address
	InputChainStorage   web3.Address
	StateChainStorage   web3.Address
	RollupInputChain    web3.Address
	RollupStateChain    web3.Address
	L1CrossLayerWitness web3.Address
	L1StandardBridge    web3.Address
	StakingManager      web3.Address
	ChallengeBeacon     web3.Address
	ChallengeLogic      web3.Address
	ChallengeFactory    web3.Address
	FeeToken            web3.Address
	DAO                 web3.Address
}

type L2ContractAddressConfig struct {
	L2CrossLayerWitness web3.Address
	L2StandardBridge    web3.Address
}

type L1ChainDeployConfig struct {
	FeeToken                web3.Address
	FraudProofWindow        uint64 // block number
	MaxEnqueueTxGasLimit    uint64
	MaxCrossLayerTxGasLimit uint64
	L2CrossLayerWitness     web3.Address
	L2ChainId               uint64
	StakingAmount           *big.Int
	*ChallengeConfig
}

type ChallengeConfig struct {
	BlockLimitPerRound uint64 // proposer
	ChallengerDeposit  *big.Int
}
