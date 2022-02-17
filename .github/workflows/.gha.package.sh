#!/usr/bin/env bash
set -ex

VERSION=$(git describe --always --tags --long)
PLATFORM=""

if [[ ${RUNNER_OS} == 'Linux' ]]; then
  PLATFORM="linux"
elif [[ ${RUNNER_OS} == 'macOS' ]]; then
  PLATFORM="darwin"
else
  PLATFORM="windows"
  exit 1
fi



env GO111MODULE=on ./build.sh
cd build/bin &&  zip  -q -r ontology-rollup-${PLATFORM}.zip * && cd -
mv build/bin/ontology-rollup-${PLATFORM}.zip .

set +x
echo "ontology-rollup-${PLATFORM}.zip |" $(md5sum ontology-rollup-${PLATFORM}.zip | cut -d ' ' -f1)
