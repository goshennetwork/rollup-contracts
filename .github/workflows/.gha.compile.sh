#!/bin/bash
mkdir rv
git clone --recursive https://github.com/r1cs/riscv-tests.git@compiled
cp -r ./riscv-tests/isa ./tests/rv32i/test_case