package merkle_mountain_tree

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/hardhat"
	"github.com/ontology-layer-2/rollup-contracts/tests"
	"github.com/pkg/errors"

	"math/big"
	"testing"
)

var emptyRoot = common.Hash{}

type testCase struct {
	cAbi     *abi.ABI
	vm       *vm.EVM
	trie     *trie.Trie
	db       *trie.Database
	contract common.Address
	sender   vm.AccountRef
}

func newCase() *testCase {
	//get contract artifact
	ars, err := hardhat.GetArtifact("MerkleMountainRange", "out")
	if err != nil {
		panic(err)
	}

	abi1, err := abi.NewABI(ars.Abi)
	if err != nil {
		panic(err)
	}
	//setup evm
	contractAddr := common.BytesToAddress([]byte("MerkleMountainRange"))
	vmenv := tests.NewEVMWithCode(map[common.Address][]byte{contractAddr: ars.DeployedBytecode})
	sender := vm.AccountRef(common.BytesToAddress([]byte("test")))
	db := trie.NewDatabase(memorydb.New())
	emptyTrie, err := trie.New(common.Hash{}, db)
	if err != nil {
		panic(err)
	}
	return &testCase{abi1, vmenv, emptyTrie, db, contractAddr, sender}
}

func (this *testCase) call(input []byte) ([]byte, error) {
	ret, _, err := this.vm.Call(this.sender, this.contract, input, math.MaxUint64, new(big.Int))
	return ret, err
}

func (this *testCase) append(hash common.Hash) error {
	//function append(bytes32 _leafHash)public
	input := this.cAbi.Methods["append"].MustEncodeIDAndInput(hash)
	ret, err := this.call(input)
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}
	return nil

}

//todo: check mmr
func TestVerify(t *testing.T) {
	newCase()
}
