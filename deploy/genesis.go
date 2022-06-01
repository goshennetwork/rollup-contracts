package deploy

import (
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/jsonrpc/transport"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/common/hexutil"
)

type GenesisAccount struct {
	Code    hexutil.Bytes           `json:"code,omitempty"`
	Storage map[web3.Hash]web3.Hash `json:"storage,omitempty"`
	Balance *hexutil.Big            `json:"balance" gencodec:"required"`
	Nonce   hexutil.Uint64          `json:"nonce,omitempty"`
}

type GenesisConfig struct {
	FeeCollectorOwner   web3.Address
	FeeCollector        web3.Address
	L2CrossLayerWitness web3.Address
	WitnessBalance      *big.Int
}

func BuildL2GenesisData(cfg *GenesisConfig) map[web3.Address]*GenesisAccount {
	genesisAccts := make(map[web3.Address]*GenesisAccount)
	setBalanceForBuiltins(genesisAccts)
	privKey := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	signer, local := SetupLocalSigner(0, privKey)
	collector := DeployL2FeeCollector(signer, cfg.FeeCollectorOwner)
	witness := DeployL2CrossLayerWitness(signer)
	overlay := local.Executor.OverlayDB
	statedb := storage.NewStateDB(storage.NewCacheDB(overlay))
	genesisAccts[cfg.FeeCollector] = getContractData(statedb, collector.Contract().Addr())
	genesisAccts[cfg.L2CrossLayerWitness] = getContractData(statedb, witness.Contract().Addr())
	if cfg.WitnessBalance != nil {
		genesisAccts[cfg.L2CrossLayerWitness].Balance = (*hexutil.Big)(cfg.WitnessBalance)
	}

	return genesisAccts
}

func getContractData(statedb *storage.StateDB, address web3.Address) *GenesisAccount {
	fee := &GenesisAccount{
		Code:    statedb.GetCode(address),
		Balance: (*hexutil.Big)(statedb.GetBalance(address)),
		Storage: make(map[web3.Hash]web3.Hash),
		Nonce:   hexutil.Uint64(statedb.GetNonce(address)),
	}

	err := statedb.ForEachStorage(address, func(key, value web3.Hash) bool {
		fee.Storage[key] = value
		return true
	})
	utils.Ensure(err)

	return fee
}

func setBalanceForBuiltins(genesisAccts map[web3.Address]*GenesisAccount) {
	builtin := web3.Address{}
	for i := 0; i < 256; i++ {
		builtin[19] = byte(i)
		genesisAccts[builtin] = &GenesisAccount{
			Balance: (*hexutil.Big)(big.NewInt(1)),
		}
	}
}

func SetupLocalSigner(chainID uint64, privKey string) (*contract.Signer, *transport.Local) {
	db := storage.NewFakeDB()
	local := transport.NewLocal(db, chainID)
	client := jsonrpc.NewClientWithTransport(local)
	signer := contract.NewSigner(privKey, client, chainID)
	signer.Submit = true
	local.SetBalance(signer.Address(), web3.Ether(1000))

	return signer, local
}