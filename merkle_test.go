package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/compiler"
	"github.com/pkg/errors"
)

func IsSolcInstalled() bool {
	output, err := exec.Command("solc", "--version").Output()
	if err != nil {
		return false
	}

	return len(output) > 0
}

var emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

var contractPath = "./contracts/libraries/MerkleTrie.t.sol"
var contractName = "./contracts/libraries/MerkleTrie.t.sol:MockMerkleTrie"
var code, cAbi = func() ([]byte, *abi.ABI) {
	if !IsSolcInstalled() {
		panic("solc not exist")
	}
	solc := &compiler.Solidity{Path: "solc"}
	a, err := solc.Compile(contractPath)
	if err != nil {
		panic(err)
	}

	abi1, err := abi.NewABI(a[contractName].Abi)
	if err != nil {
		panic(err)
	}
	return common.FromHex(a[contractName].BinRuntime), abi1
}()

var address = common.BytesToAddress([]byte("merkleContract"))
var sender = vm.AccountRef(common.BytesToAddress([]byte("test")))
var updateFunc = cAbi.Methods["checkUpdate"]
var getFunc = cAbi.Methods["checkGet"]
var insertTrieNodeFunc = cAbi.Methods["insertTrieNode"]

type GetOutPut struct {
	Res []byte
}

//testCase is a single test object, it hold the trie info and evm storage db info.
type testCase struct {
	vm   *vm.EVM
	trie *trie.Trie
	db   *trie.Database
}

func newCase() *testCase {
	cfg := defaultsConfig()
	cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	var (
		vmenv = runtime.NewEnv(cfg)
	)
	if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber); rules.IsBerlin {
		cfg.State.PrepareAccessList(cfg.Origin, &address, vm.ActivePrecompiles(rules), nil)
	}
	vmenv.StateDB.CreateAccount(address)
	// set the receiver's (the executing contract) code for execution.
	vmenv.StateDB.SetCode(address, code)
	db := trie.NewDatabase(memorydb.New())
	emptyTrie, err := trie.New(common.Hash{}, db)
	if err != nil {
		panic(err)
	}
	return &testCase{vmenv, emptyTrie, db}
}

func (this *testCase) updateString(key, value string, root common.Hash) error {
	return this.update([]byte(key), []byte(value), root)
}

func (this *testCase) update(key, value []byte, root common.Hash) error {
	/*function update(
	      bytes memory _key,
	      bytes memory _value,
	      bytes32 _root,
	      bytes32 _expectRoot
	  ) external;
	*/
	this.trie.Update(key, value)
	this.trie.Commit(nil)
	fmt.Printf("updated: key: 0x%x, value: 0x%x, newRoot: %s\n", key, value, this.trie.Hash())
	input, err := updateFunc.EncodeIDAndInput(key, value, root, this.trie.Hash())

	if err != nil {
		return errors.Wrap(err, "update input")
	}
	ret, _, err := this.vm.Call(sender, address, input, defaultsConfig().GasLimit, new(big.Int))
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}
	return nil
}

func (this *testCase) get(key []byte, root common.Hash, want []byte) error {
	//function get(bytes memory _key, bytes32 _root) external override returns (bytes memory);
	input, err := getFunc.EncodeIDAndInput(key, root)
	if err != nil {
		return errors.Wrap(err, "get input")
	}
	ret, _, err := this.vm.Call(sender, address, input, defaultsConfig().GasLimit, new(big.Int))
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}

	m, err := getFunc.Outputs.Decode(ret)
	if err != nil {
		return err
	}
	out := m.(map[string]interface{})["0"].([]byte)
	if !bytes.Equal(out, want) {
		return fmt.Errorf("not equal, want: 0x%x, but get: 0x%x", want, out)
	}
	return nil
}

func (this *testCase) insertTrieNode(encoded []byte) error {
	//function insertTrieNode(bytes calldata anything)external;
	input, err := insertTrieNodeFunc.EncodeIDAndInput(encoded)
	if err != nil {
		return err
	}
	ret, _, err := this.vm.Call(sender, address, input, defaultsConfig().GasLimit, new(big.Int))
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}
	return nil
}

func ensure(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

//test null key
func TestNull(t *testing.T) {
	trieCase := newCase()
	k := make([]byte, 32)
	v := []byte("test")
	err := trieCase.update(k, v, emptyRoot)
	ensure(t, err)
	err = trieCase.get(k, trieCase.trie.Hash(), v)
	ensure(t, err)
}

//update with missing root(i.g. not from empty)
func TestMissingRoot(t *testing.T) {
	trieCase := newCase()
	err := trieCase.update([]byte("test"), []byte("test"), common.Hash{})
	if err == nil {
		t.Fatal("update for invalid root")
	}
}

func TestInsert(t *testing.T) {
	trieCase := newCase()
	k, v := []byte("doe"), []byte("reindeer")
	ensure(t, trieCase.update(k, v, emptyRoot))
	k, v = []byte("dog"), []byte("puppy")
	ensure(t, trieCase.update(k, v, trieCase.trie.Hash()))
	k, v = []byte("dogglesworth"), []byte("cat")
	ensure(t, trieCase.update(k, v, trieCase.trie.Hash()))
}

func TestGet(t *testing.T) {
	trieCase := newCase()
	k, v := []byte("doe"), []byte("reindeer")
	ensure(t, trieCase.update(k, v, emptyRoot))
	k, v = []byte("dog"), []byte("puppy")
	ensure(t, trieCase.update(k, v, trieCase.trie.Hash()))
	k, v = []byte("dogglesworth"), []byte("cat")
	ensure(t, trieCase.update(k, v, trieCase.trie.Hash()))
	for i := 0; i < 2; i++ {
		ensure(t, trieCase.get([]byte("dog"), trieCase.trie.Hash(), []byte("puppy")))
		//get unknown key
		err := trieCase.get([]byte("unknown"), trieCase.trie.Hash(), []byte{})
		if err == nil {
			t.Fatal("get value for invalid key")
		}
		if i == 1 {
			return
		}
	}

}

// in origin trie logic, if empty value will delete the leaf node to trim trie, but in contract, it is hard to trim triem.
// so the contract just store value to 0x80(rlp[]byte{}).
func TestEmptyValue(t *testing.T) {
	trieCase := newCase()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, s := range vals {
		err := trieCase.updateString(s.k, s.v, trieCase.trie.Hash())
		ensure(t, err)
	}
}

//rebuild a new trie from existing trie node
func TestReplication(t *testing.T) {
	trieCase := newCase()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"dog", "puppy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
	}
	for _, s := range vals {
		err := trieCase.updateString(s.k, s.v, trieCase.trie.Hash())
		ensure(t, err)
	}
	copied := newCase()
	trieCase.db.Commit(trieCase.trie.Hash(), false, func(hash common.Hash) {
		value, err := trieCase.db.Node(hash)
		ensure(t, err)
		ensure(t, copied.insertTrieNode(value))
	})
	for _, s := range vals {
		err := copied.get([]byte(s.k), trieCase.trie.Hash(), []byte(s.v))
		ensure(t, err)
	}
}

//large value test
func TestLargeValue(t *testing.T) {
	trieCase := newCase()
	res := []struct{ k, v []byte }{
		{[]byte("key1"), []byte{99, 99, 99, 99}},
		{[]byte("key2"), bytes.Repeat([]byte{1}, 32)},
	}
	for _, s := range res {
		ensure(t, trieCase.update(s.k, s.v, trieCase.trie.Hash()))
	}
	for _, s := range res {
		ensure(t, trieCase.get(s.k, trieCase.trie.Hash(), s.v))
	}
	copied := newCase()
	trieCase.db.Commit(trieCase.trie.Hash(), false, func(hash common.Hash) {
		value, err := trieCase.db.Node(hash)
		ensure(t, err)
		ensure(t, copied.insertTrieNode(value))
	})
	for _, s := range res {
		ensure(t, copied.get(s.k, trieCase.trie.Hash(), s.v))
	}

}

// TestRandomCases tests some cases that were found via random fuzzing
func TestRandomCases(t *testing.T) {
	trieCase := newCase()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	res := Generate(r)
	for _, s := range res {
		ensure(t, trieCase.update(s.k, s.v, trieCase.trie.Hash()))
		ensure(t, trieCase.get(s.k, trieCase.trie.Hash(), s.v))
	}
	copied := newCase()
	trieCase.db.Commit(trieCase.trie.Hash(), false, func(hash common.Hash) {
		value, err := trieCase.db.Node(hash)
		ensure(t, err)
		ensure(t, copied.insertTrieNode(value))
	})
	for _, s := range res { //same key may cover pre value, so get it from trie to ensure correctness
		if err := copied.get(s.k, trieCase.trie.Hash(), trieCase.trie.Get(s.k)); err != nil {
			t.Fatal(err)
		}
	}
}

//Generate random k,v for fuzzy test
func Generate(r *rand.Rand) []struct{ k, v []byte } {
	var allKeys [][]byte
	for {
		if len(allKeys) < 100 || r.Intn(100) < 60 {
			// new key
			key := make([]byte, r.Intn(50))
			r.Read(key)
			allKeys = append(allKeys, key)
		} else {
			break
		}
	}

	res := make([]struct{ k, v []byte }, len(allKeys), len(allKeys))
	for i := range allKeys {
		res[i].k = allKeys[i]
		res[i].v = make([]byte, 8)
		binary.BigEndian.PutUint64(res[i].v, uint64(i))
	}
	return res
}

func defaultsConfig() (cfg *runtime.Config) {
	cfg = new(runtime.Config)
	if cfg.ChainConfig == nil {
		cfg.ChainConfig = &params.ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      new(big.Int),
			DAOForkBlock:        new(big.Int),
			DAOForkSupport:      false,
			EIP150Block:         new(big.Int),
			EIP150Hash:          common.Hash{},
			EIP155Block:         new(big.Int),
			EIP158Block:         new(big.Int),
			ByzantiumBlock:      new(big.Int),
			ConstantinopleBlock: new(big.Int),
			PetersburgBlock:     new(big.Int),
			IstanbulBlock:       new(big.Int),
			MuirGlacierBlock:    new(big.Int),
			BerlinBlock:         new(big.Int),
			LondonBlock:         new(big.Int),
		}
	}

	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = math.MaxUint64
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
	if cfg.BaseFee == nil {
		cfg.BaseFee = big.NewInt(params.InitialBaseFee)
	}
	return
}
