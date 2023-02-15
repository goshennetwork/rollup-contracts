module github.com/ontology-layer-2/rollup-contracts

go 1.16

require (
	github.com/andybalholm/brotli v1.0.4
	github.com/ethereum/go-ethereum v1.10.3
	github.com/laizy/log v0.1.0
	github.com/laizy/web3 v0.1.14-0.20230206030724-97fda71ec1d5
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/protolambda/go-kzg v0.0.0-20221224134646-c91cee5e954e
	github.com/stretchr/testify v1.7.2
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/urfave/cli/v2 v2.10.2
)

replace github.com/ethereum/go-ethereum v1.10.3 => github.com/protolambda/go-ethereum v1.7.4-0.20220917163714-e091e9a7d5a1

replace github.com/laizy/web3 => ../web3
