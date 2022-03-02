#!/usr/bin/env bash
set -ex

VERSION=$(git describe --always --tags --long)
bash ./.github/workflows/.gha.gofmt.sh
bash ./.github/workflows/.gha.compile.sh
bash ./.github/workflows/.gha.gotest.sh
