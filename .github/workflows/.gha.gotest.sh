#!/bin/bash
# code from https://github.com/Seklfreak/Robyul2
unset dirs files
dirs=$(go list ./... |grep -v rollup-contracts$)
set -x -e

for d in $dirs
do
  go test -v $d
done
