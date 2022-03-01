#!/bin/bash
mkdir rv
docker run --rm -it -v=$PWD/rv:/build rpirea/riscv32-unknown-elf:pulp bash -c 'apt-get update && apt-get install -y autoconf automake autotools-dev curl python3 libmpc-dev libmpfr-dev libgmp-dev gawk build-essential bison flex texinfo gperf libtool patchutils bc zlib1g-dev libexpat-dev && git clone --recursive https://github.com/r1cs/riscv-tests.git && \
cd riscv-tests && autoconf && ./configure --with-xlen=32 --prefix=/build && make && make install'