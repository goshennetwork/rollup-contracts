#!/usr/bin/env bash
set -ex

VERSION=$(git describe --always --tags --long)
GOPRIVATE=github.com/ontology-layer-2 go get github.com/ethereum/go-ethereum@v1.10.3
bash ./.github/workflows/.gha.gofmt.sh
bash ./.github/workflows/.gha.gotest.sh
