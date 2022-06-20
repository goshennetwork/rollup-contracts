set -ex

go build -o build/rollup ./cmd/rollupcli
go build -o build/uploader ./cmd/uploader
