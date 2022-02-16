et -ex

VERSION=$(git describe --always --tags --long)

if [ $RUNNER_OS == 'Linux' ]; then
  echo "linux sys"
  GOPRIVATE=github.com/ontology-layer-2 go get github.com/ontology-layer-2/go-ethereum@8283586
  bash ./.gha.gofmt.sh
  bash ./.gha.gotest.sh
elif [ $RUNNER_OS == 'osx' ]; then
  echo "osx sys"
  ./build.sh
else
  echo "win sys not supported yet"
fi
