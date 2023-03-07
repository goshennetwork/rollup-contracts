#!/bin/bash
mkdir rv
git clone --branch compiled --recursive https://github.com/ontology-layer-2/riscv-tests.git
cp -r ./riscv-tests/isa ./tests/rv32i