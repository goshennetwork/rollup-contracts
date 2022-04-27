package tests

import (
	"math/big"
	"time"

	"github.com/laizy/web3/executor/fakedb"

	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/evm/params"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/executor"
)

func NewEVMWithCode(contracts map[web3.Address][]byte) *evm.EVM {
	var hashFn evm.GetHashFunc = func(u uint64) web3.Hash {
		var h web3.Hash
		h.SetBytes(crypto.Keccak256(new(big.Int).SetUint64(u).Bytes()))
		return h
	}
	caccheDB := storage.NewCacheDB(overlaydb.NewOverlayDB(&fakedb.FakeDB{}))
	statedb := storage.NewStateDB(caccheDB, web3.Hash{}, web3.Hash{})
	ctx := executor.NewEVMBlockContext(0, uint64(time.Now().Unix()), hashFn)
	vmenv := evm.NewEVM(ctx, evm.TxContext{}, statedb, params.MainnetChainConfig, evm.Config{})
	for addr, code := range contracts {
		vmenv.StateDB.CreateAccount(addr)
		// set the receiver's (the executing contract) code for execution.
		vmenv.StateDB.SetCode(addr, code)
	}
	return vmenv
}
