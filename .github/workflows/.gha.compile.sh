#!/bin/bash
mkdir rv
git clone --branch compiled --recursive https://github.com/r1cs/riscv-tests.git
cp -r ./riscv-tests/isa ./tests/rv32i