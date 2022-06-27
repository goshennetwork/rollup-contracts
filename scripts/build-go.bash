set -ex

go build -o build/rollup ./cmd/rollupcli
go build -o build/uploader ./cmd/uploader
go build -o build/sync-service ./cmd/sync-service
