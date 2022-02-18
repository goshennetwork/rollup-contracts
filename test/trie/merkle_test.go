package trie

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
	"github.com/laizy/web3/hardhat"
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

var contractPath = "../../contracts/libraries/MerkleTrie.t.sol"
var contractName = "../../contracts/libraries/MerkleTrie.t.sol:MockMerkleTrie"
var code, cAbi = func() ([]byte, *abi.ABI) {
	ars, err := hardhat.GetArtifact("MockMerkleTrie", "out")
	if err != nil {
		panic(err)
	}
	abi1, err := abi.NewABI(ars.Abi)
	if err != nil {
		panic(err)
	}
	return common.FromHex(ars.DeployedBytecode.String()), abi1
}()

var address = common.BytesToAddress([]byte("merkleContract"))
var sender = vm.AccountRef(common.BytesToAddress([]byte("test")))
var rawUpdateFunc = cAbi.Methods["rawUpdate"]
var rawGetFunc = cAbi.Methods["rawGet"]
var checkUpdateFunc = cAbi.Methods["checkUpdate"]
var checkGetFunc = cAbi.Methods["checkGet"]
var insertTrieNodeFunc = cAbi.Methods["insertTrieNode"]

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
	if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, true); rules.IsBerlin {
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

//copyTrie create a new EVm trie with old golang trie db
func (this *testCase) newEvmTrie() *testCase {
	n := newCase()
	n.trie = this.trie
	n.db = this.db
	return n
}

func (this *testCase) checkUpdateString(key, value string) error {
	return this.checkUpdate([]byte(key), []byte(value))
}

func (this *testCase) update(key, value []byte, root common.Hash) error {
	//function rawUpdate( bytes memory _key,bytes memory _value,bytes32 _root)external
	input := rawUpdateFunc.MustEncodeIDAndInput(key, value, root)
	ret, _, err := this.vm.Call(sender, address, input, defaultsConfig().GasLimit, new(big.Int))
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}
	return nil
}

func (this *testCase) get(key []byte, root common.Hash) error {
	//function rawGet(bytes memory _key,bytes32 _root)external returns (bytes memory)
	input := rawGetFunc.MustEncodeIDAndInput(key, root)
	ret, _, err := this.vm.Call(sender, address, input, defaultsConfig().GasLimit, new(big.Int))
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}
	return nil
}

func (this *testCase) checkUpdate(key, value []byte) error {
	/*function checkUpdate(
	      bytes memory _key,
	      bytes memory _value,
	      bytes32 _root,
	      bytes32 _expectRoot
	  ) external;
	*/
	startRoot := this.trie.Hash()
	this.trie.Update(key, value)
	this.trie.Commit(nil)
	fmt.Printf("updated: key: 0x%x, value: 0x%x, newRoot: %s\n", key, value, this.trie.Hash())
	input, err := checkUpdateFunc.EncodeIDAndInput(key, value, startRoot, this.trie.Hash())
	if err != nil {
		return errors.Wrap(err, "checkUpdate input")
	}
	ret, _, err := this.vm.Call(sender, address, input, defaultsConfig().GasLimit, new(big.Int))
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}
	return nil
}

func (this *testCase) checkGet(key []byte) error {
	//function checkGet(bytes memory _key, bytes32 _root) external override returns (bytes memory);
	input, err := checkGetFunc.EncodeIDAndInput(key, this.trie.Hash())
	if err != nil {
		return errors.Wrap(err, "checkGet input")
	}
	ret, _, err := this.vm.Call(sender, address, input, defaultsConfig().GasLimit, new(big.Int))
	if err != nil {
		s, _ := web3.DecodeRevert(ret)
		return errors.Wrap(err, s)
	}

	m, err := checkGetFunc.Outputs.Decode(ret)
	if err != nil {
		return err
	}
	out := m.(map[string]interface{})["0"].([]byte)
	if !bytes.Equal(out, this.trie.Get(key)) {
		return fmt.Errorf("not equal, want: 0x%x, but checkGet: 0x%x", this.trie.Get(key), out)
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
	err := trieCase.checkUpdate(k, v)
	ensure(t, err)
	err = trieCase.checkGet(k)
	ensure(t, err)
}

//update with missing root(i.g. not from empty)
func TestMissingRoot(t *testing.T) {
	trieCase := newCase()
	err := trieCase.update([]byte("test"), []byte("test"), common.Hash{})
	if err == nil {
		t.Fatal("update for invalid root")
	}
	err = trieCase.get([]byte(""), common.Hash{})
	if err == nil {
		t.Fatal("get for invalid root")
	}

	err = trieCase.get([]byte(""), emptyRoot)
	if err == nil {
		t.Fatal("get for invalid root")
	}

}

func TestInsert(t *testing.T) {
	trieCase := newCase()
	k, v := []byte("doe"), []byte("reindeer")
	ensure(t, trieCase.checkUpdate(k, v))
	k, v = []byte("dog"), []byte("puppy")
	ensure(t, trieCase.checkUpdate(k, v))
	k, v = []byte("dogglesworth"), []byte("cat")
	ensure(t, trieCase.checkUpdate(k, v))
}

func TestGet(t *testing.T) {
	trieCase := newCase()
	k, v := []byte("doe"), []byte("reindeer")
	ensure(t, trieCase.checkUpdate(k, v))
	k, v = []byte("dog"), []byte("puppy")
	ensure(t, trieCase.checkUpdate(k, v))
	k, v = []byte("dogglesworth"), []byte("cat")
	ensure(t, trieCase.checkUpdate(k, v))
	for i := 0; i < 2; i++ {
		ensure(t, trieCase.checkGet([]byte("dog")))
		//checkGet unknown key
		err := trieCase.checkGet([]byte("unknown"))
		if err == nil {
			t.Fatal("checkGet value for invalid key")
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
		err := trieCase.checkUpdateString(s.k, s.v)
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
		err := trieCase.checkUpdateString(s.k, s.v)
		ensure(t, err)
	}
	copied := trieCase.newEvmTrie()
	trieCase.db.Commit(trieCase.trie.Hash(), false, func(hash common.Hash) {
		value, err := trieCase.db.Node(hash)
		ensure(t, err)
		ensure(t, copied.insertTrieNode(value))
	})
	for _, s := range vals {
		err := copied.checkGet([]byte(s.k))
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
		ensure(t, trieCase.checkUpdate(s.k, s.v))
	}
	for _, s := range res {
		ensure(t, trieCase.checkGet(s.k))
	}
	copied := trieCase.newEvmTrie()
	trieCase.db.Commit(trieCase.trie.Hash(), false, func(hash common.Hash) {
		value, err := trieCase.db.Node(hash)
		ensure(t, err)
		ensure(t, copied.insertTrieNode(value))
	})
	for _, s := range res {
		ensure(t, copied.checkGet(s.k))
	}

}

// TestRandomCases tests some cases that were found via random fuzzing
func TestRandomCases(t *testing.T) {
	trieCase := newCase()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	res := Generate(r)
	for _, s := range res {
		ensure(t, trieCase.checkUpdate(s.k, s.v))
		ensure(t, trieCase.checkGet(s.k))
	}
	copied := trieCase.newEvmTrie()
	trieCase.db.Commit(trieCase.trie.Hash(), false, func(hash common.Hash) {
		value, err := trieCase.db.Node(hash)
		ensure(t, err)
		ensure(t, copied.insertTrieNode(value))
	})
	for _, s := range res { //same key may cover pre value, so checkGet it from trie to ensure correctness
		if err := copied.checkGet(s.k); err != nil {
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
