package rv32i

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/hardhat"
	"github.com/mitchellh/mapstructure"
	"github.com/ontology-layer-2/rollup-contracts/tests"
	"github.com/pkg/errors"
)

func TestRV32I(t *testing.T) {
	tests, err := filepath.Glob("test_case/isa/rv32ui-v-*")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range tests {
		fmt.Println(f)
		image, entry, err := getProgramImage(f)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("entry: ", entry)
		fmt.Println(image[entry])
		ret, err := start(image, entry)
		if err != nil {
			t.Log(err)
			r, _ := web3.DecodeRevert(ret)
			t.Fatalf("revert: %s", r)
		}
	}
}

//func TestHello(t *testing.T) {
//	runFile(t, "riscv-ia")
//}

func runFile(t *testing.T, fileName string) {
	image, entry, err := getProgramImage(fileName)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(entry)
	ret, err := start(image, entry)
	if err != nil {
		t.Log(err)
		r, _ := web3.DecodeRevert(ret)
		t.Fatalf("revert: %s", r)
	}
}

func start(ram map[uint32]uint32, entrypoint uint32) ([]byte, error) {
	this := newCase()
	fmt.Println("copy ram...")
	if ret, err := this.copyRam(ram); err != nil {
		return ret, err
	}
	fmt.Println("init register...")
	for r := uint32(0); r < 33; r++ {
		if ret, err := this.writeRegister(r, 0); err != nil {
			return ret, err
		}
	}
	fmt.Println("start...")
	now := time.Now()
	r, err := this.start(entrypoint)
	var i interface{}
	var root common.Hash
	var num uint32
	var insn uint32
	if err == nil {
		i, _ = this.rvAbi.Methods["start"].Outputs.Decode(r)
		if err := mapstructure.Decode(i.(map[string]interface{})["0"], &root); err != nil {
			panic(err)
		}
		if err := mapstructure.Decode(i.(map[string]interface{})["1"], &num); err != nil {
			panic(err)
		}
		if err := mapstructure.Decode(i.(map[string]interface{})["2"], &insn); err != nil {
			panic(err)
		}
	}
	fmt.Printf(" consume %d time: %v, last insn: 0x%x\n", num, time.Since(now), insn)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, this.readMemory(root, 1200))
	fmt.Println(string(b))
	return r, err
}

type testCase struct {
	evm *vm.EVM
	//interpreter contract bi
	rvAbi *abi.ABI
	//interpreter contract addr
	rvAddr common.Address
	//machine state abi
	ramAbi *abi.ABI
	//machine state addr
	ramAddr common.Address
	ramTrie *trie.Trie
	sender  vm.AccountRef
	mdb     *trie.Database
}

func newCase() *testCase {
	//get contract artifact
	rvA, err := hardhat.GetArtifact("Interpretor", "out")
	if err != nil {
		panic(err)
	}

	rvAbi, err := abi.NewABI(rvA.Abi)
	if err != nil {
		panic(err)
	}

	ramA, err := hardhat.GetArtifact("MachineState", "out")
	if err != nil {
		panic(err)
	}
	ramAbi, err := abi.NewABI(ramA.Abi)
	if err != nil {
		panic(err)

	}
	ramAddr := common.BytesToAddress([]byte("MachineState"))
	vmevm := tests.NewEVMWithCode(map[common.Address][]byte{ramAddr: ramA.DeployedBytecode})
	sender := vm.AccountRef(common.BytesToAddress([]byte("test")))
	mdb := trie.NewDatabase(memorydb.New())
	trie, err := trie.New(common.Hash{}, mdb)
	if err != nil {
		panic(err)
	}
	//constructor(address state)
	type CCC struct {
		State common.Address
	}
	input, err := rvAbi.Constructor.Inputs.Encode(CCC{State: ramAddr})
	if err != nil {
		panic(err)
	}
	_, rvAddr, _, err := vmevm.Create(sender, append(rvA.Bytecode, input...), math.MaxUint64, new(big.Int))
	if err != nil {
		panic(err)
	}
	return &testCase{vmevm, rvAbi, rvAddr, ramAbi, ramAddr, trie, sender, mdb}
}

//copy ram to evm
func (this *testCase) copyRam(ram map[uint32]uint32) ([]byte, error) {
	for k, v := range ram {
		//key is big endian
		kk, vv := make([]byte, 4), make([]byte, 4)
		binary.BigEndian.PutUint32(kk[:], k)
		//v is little endian
		binary.LittleEndian.PutUint32(vv[:], v)
		this.ramTrie.Update(kk, vv)
	}
	if _, err := this.ramTrie.Commit(nil); err != nil {
		panic(err)
	}
	err := this.mdb.Commit(this.ramTrie.Hash(), false, func(hash common.Hash) {
		node, err := this.mdb.Node(hash)
		if err != nil {
			panic(err)
		}
		_, err = this.insertTrieNode(node)
		if err != nil {
			panic(err)
		}
	})
	return nil, err
}

func (this *testCase) writeMemory(k, v uint32) (ret []byte, err error) {
	defer func() { err = errors.Wrap(err, "write memory") }()
	//function writeMemory(
	//        bytes32 root,
	//        uint32 ptr,
	//        uint32 value
	//    ) public returns (bytes32)
	input := this.ramAbi.Methods["writeMemory"].MustEncodeIDAndInput(this.ramTrie.Hash(), k, v)
	kk, vv := make([]byte, 4), make([]byte, 4)
	//key is big endian
	binary.BigEndian.PutUint32(kk[:], k)
	//v is little endian
	binary.LittleEndian.PutUint32(vv[:], v)
	this.ramTrie.Update(kk, vv)
	ret, _, err = this.evm.Call(this.sender, this.ramAddr, input, math.MaxUint64, new(big.Int))
	return
}

func (this *testCase) readMemory(root common.Hash, k uint32) uint32 {
	//function readMemory(bytes32 root, uint32 ptr) public view returns (uint32)
	method := this.ramAbi.Methods["readMemory"]
	input := method.MustEncodeIDAndInput(root, k)
	ret, _, err := this.evm.Call(this.sender, this.ramAddr, input, math.MaxUint64, new(big.Int))
	if err != nil {
		panic(err)
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		panic(err)
	}
	var out uint32
	if err := mapstructure.Decode(i.(map[string]interface{})["0"], &out); err != nil {
		panic(err)
	}
	return out
}

func (this *testCase) insertTrieNode(data []byte) (ret []byte, err error) {
	//function insertTrieNode(bytes calldata _node)public
	input := this.ramAbi.Methods["insertTrieNode"].MustEncodeIDAndInput(data)
	ret, _, err = this.evm.Call(this.sender, this.ramAddr, input, math.MaxUint64, new(big.Int))
	return
}

func (this *testCase) writeRegister(k, v uint32) (ret []byte, err error) {
	if k > (1<<8)-1 {
		panic("out of uint8")
	}
	defer func() { err = errors.Wrap(err, "write register") }()
	// function writeRegister(
	//        bytes32 root,
	//        uint32 regid,
	//        uint32 value
	//    ) public returns (bytes32)
	input := this.ramAbi.Methods["writeRegister"].MustEncodeIDAndInput(this.ramTrie.Hash(), k, v)
	vv := make([]byte, 4)
	//register value is little endian
	binary.LittleEndian.PutUint32(vv[:], v)
	if k != 0 { //x0 do not write to db
		this.ramTrie.Update([]byte{uint8(k)}, vv)
	}
	ret, _, err = this.evm.Call(this.sender, this.ramAddr, input, math.MaxUint64, new(big.Int))
	return
}

func (this *testCase) start(entrypoint uint32) (ret []byte, err error) {
	//function start(bytes32 _root,uint32 _entrypoint) public
	input := this.rvAbi.Methods["start"].MustEncodeIDAndInput(this.ramTrie.Hash(), entrypoint)
	ret, _, err = this.evm.Call(this.sender, this.rvAddr, input, math.MaxUint64, new(big.Int))
	return
}
