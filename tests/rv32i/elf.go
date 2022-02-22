package rv32i

import (
	"debug/elf"
	"fmt"
	"io/ioutil"
	"sort"
)

func filterFileSymbol(f string, want func(string) bool) []elf.Symbol {
	file, err := elf.Open(f)
	if err != nil {
		panic(err)
	}
	symbols, err := file.Symbols()
	if err != nil {
		panic(err)
	}
	return filterSymbol(symbols, want)
}

func filterSymbol(symbols []elf.Symbol, want func(string) bool) []elf.Symbol {
	var s []elf.Symbol
	for _, symbol := range symbols {
		if want(symbol.Name) {
			s = append(s, symbol)
		}
	}
	return s
}

func GetImageWithEntrypoint(s string) ([]byte, uint32) {
	m, en, err := getProgramImage(s)
	if err != nil {
		panic(err)
	}
	//change map to slice
	addrs := make([]uint32, len(m))
	i := 0
	for addr, _ := range m {
		addrs[i] = addr
		i += 1
	}
	sort.Slice(addrs, func(i, j int) bool { //sort addr
		return addrs[i] < addrs[j]
	})
	maxaddr := addrs[len(addrs)-1]
	file := make([]byte, maxaddr+1)
	for _, addr := range addrs {
		file[addr] = m[addr]
	}
	return file, en
}

func getProgramImage(s string) (map[uint32]byte, uint32, error) {
	fileimage := make(map[uint32]byte)
	file, err := elf.Open(s)
	if err != nil {
		return nil, 0, err
	}
	if err := checkForRiscv32(file); err != nil {
		return nil, 0, err
	}

	ddd, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, 0, err
	}
	entry := file.Entry
	for _, ph := range file.Progs {
		if ph.Type == elf.PT_LOAD {
			data := make([]byte, ph.Memsz)
			copy(data, ddd[ph.Off:ph.Off+ph.Filesz])
			if err != nil {
				return nil, 0, err
			}
			writeToMap(data, uint32(ph.Vaddr), fileimage)
		}
	}
	return fileimage, uint32(entry), nil
}

//write to memory map
func writeToMap(d []byte, offset uint32, m map[uint32]byte) {
	for i := 0; i < len(d); i++ {
		m[offset] = d[i]
		offset++
	}
}
func checkForRiscv32(f *elf.File) error {
	if len(f.Progs) == 0 {
		return fmt.Errorf("no program header in elf")
	}
	if f.Class != elf.ELFCLASS32 {
		return fmt.Errorf("riscv32 target should be arch32")
	}
	if f.Data != elf.ELFDATA2LSB {
		return fmt.Errorf("riscv32 target should be LSB")
	}

	for _, ph := range f.Progs {
		if ph.Type == elf.PT_DYNAMIC {
			return fmt.Errorf("riscv32 do not support dynamic link")
		}
	}
	return nil
}
