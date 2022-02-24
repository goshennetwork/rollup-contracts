package rv32i

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"path/filepath"
	"testing"

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
)

func TestRV32I(t *testing.T) {
	tests, err := filepath.Glob("test_case/isa/rv32ui-v-*")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range tests {
		fmt.Println(f)
		image, entry := GetImageWithEntrypoint(f)
		if len(image)&3 != 0 {
			panic("wrong image")
		}
		m := make(map[uint32]uint32)
		for i := uint32(0); i < uint32(len(image)); i += 4 {
			m[i] = binary.LittleEndian.Uint32(image[i : i+4])
		}
		ret, err := start(m, entry)
		if err != nil {
			t.Log(err)
			r, _ := web3.DecodeRevert(ret)
			t.Fatalf("revert: %s", r)
		}
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
	return this.start(entrypoint)
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
	trie, err := trie.New(common.Hash{}, trie.NewDatabase(memorydb.New()))
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
	return &testCase{vmevm, rvAbi, rvAddr, ramAbi, ramAddr, trie, sender}
}

func (this *testCase) newInterpretor() {
}

//copy ram to evm
func (this *testCase) copyRam(ram map[uint32]uint32) ([]byte, error) {
	for k, v := range ram {
		ret, err := this.writeMemory(k, v)
		if err != nil {
			return ret, err
		}
	}
	return nil, nil
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
