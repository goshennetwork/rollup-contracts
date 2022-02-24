package rv32i

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io/ioutil"
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

func getProgramImage(s string) (map[uint32]uint32, uint32, error) {
	fileimage := make(map[uint32]uint32)
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
func writeToMap(d []byte, offset uint32, m map[uint32]uint32) {
	if len(d)&3 != 0 {
		panic("not align 4")
	}
	i := 0
	for i < len(d)-4 {
		m[offset] = binary.LittleEndian.Uint32(d[i : i+4])
		i += 4
		offset += 4
	} //last 4 byte
	m[offset] = binary.LittleEndian.Uint32(d[i:])
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
