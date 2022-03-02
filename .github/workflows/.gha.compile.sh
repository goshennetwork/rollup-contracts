#!/bin/bash
mkdir rv
docker run --rm -v=$PWD/rv:/build maxxing/riscv32-toolchain:latest bash -c 'apt-get --allow-releaseinfo-change update && apt-get install -y git && apt-get install -y autoconf automake autotools-dev curl python3 libmpc-dev libmpfr-dev libgmp-dev gawk build-essential bison flex texinfo gperf libtool patchutils bc zlib1g-dev libexpat-dev && cd /build  && git clone --recursive https://github.com/r1cs/riscv-tests.git && \
       cd riscv-tests   && autoconf && ./configure  --prefix=/build --with-xlen=32 && make && make install'
cp -r ./rv/share/riscv-tests/isa ./tests/rv32i/test_case