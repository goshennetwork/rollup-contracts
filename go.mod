module github.com/ontology-layer-2/rollup-contracts

go 1.16

require (
	github.com/andybalholm/brotli v1.0.4
	github.com/ethereum/go-ethereum v1.10.3
	github.com/laizy/log v0.1.0
	github.com/laizy/web3 v0.1.14-0.20230215075148-fe3cc512a961
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/protolambda/go-kzg v0.0.0-20221224134646-c91cee5e954e
	github.com/stretchr/testify v1.7.2
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/umbracle/fastrlp v0.1.0 // indirect
	github.com/urfave/cli/v2 v2.10.2
	github.com/valyala/fastjson v1.6.4 // indirect
	golang.org/x/crypto v0.6.0 // indirect
)

replace github.com/ethereum/go-ethereum v1.10.3 => github.com/protolambda/go-ethereum v1.7.4-0.20220917163714-e091e9a7d5a1
