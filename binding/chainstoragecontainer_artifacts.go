package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiChainStorageContainer *abi.ABI

// ChainStorageContainerAbi returns the abi of the ChainStorageContainer contract
func ChainStorageContainerAbi() *abi.ABI {
	return abiChainStorageContainer
}

var binChainStorageContainer []byte

// ChainStorageContainerBin returns the bin of the ChainStorageContainer contract
func ChainStorageContainerBin() []byte {
	return binChainStorageContainer
}

var binRuntimeChainStorageContainer []byte

// ChainStorageContainerBinRuntime returns the runtime bin of the ChainStorageContainer contract
func ChainStorageContainerBinRuntime() []byte {
	return binRuntimeChainStorageContainer
}

func init() {
	var err error
	abiChainStorageContainer, err = abi.NewABI(abiChainStorageContainerStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse ChainStorageContainer abi: %v", err))
	}
	if len(binChainStorageContainerStr) != 0 {
		binChainStorageContainer, err = hex.DecodeString(binChainStorageContainerStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse ChainStorageContainer bin: %v", err))
		}
	}
	if len(binRuntimeChainStorageContainerStr) != 0 {
		binRuntimeChainStorageContainer, err = hex.DecodeString(binRuntimeChainStorageContainerStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse ChainStorageContainer bin runtime: %v", err))
		}
	}
}

var binChainStorageContainerStr = "0x608060405234801561001057600080fd5b50610c07806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80636483ec251161005b5780636483ec25146101045780637ab4339d146101175780638da5cb5b1461012a578063ada867981461013f57600080fd5b806304f3bcec1461008257806331fe0949146100d25780635682afa9146100ef575b600080fd5b6000546100a89062010000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6001545b60405167ffffffffffffffff90911681526020016100c9565b6101026100fd366004610881565b610160565b005b6100d66101123660046108b2565b610339565b61010261012536600461092a565b6104ce565b6101326106b3565b6040516100c99190610a0a565b61015261014d366004610881565b610741565b6040519081526020016100c9565b6000546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a4478906101bc90600290600401610ad0565b602060405180830381865afa1580156101d9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101fd9190610bae565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102be57604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e60648201526084015b60405180910390fd5b60015467ffffffffffffffff82161115610334576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f63616e277420726573697a65206265796f6e6420636861696e206c656e67746860448201526064016102b5565b600155565b600080546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a44789061039690600290600401610ad0565b602060405180830381865afa1580156103b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103d79190610bae565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461049357604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e60648201526084016102b5565b5060018054808201825560008290527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601829055545b919050565b600054610100900460ff16158080156104ee5750600054600160ff909116105b806105085750303b158015610508575060005460ff166001145b610594576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016102b5565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156105f257600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b82516106059060029060208601906107e8565b50600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000073ffffffffffffffffffffffffffffffffffffffff85160217905580156106ae57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b600280546106c090610a7d565b80601f01602080910402602001604051908101604052809291908181526020018280546106ec90610a7d565b80156107395780601f1061070e57610100808354040283529160200191610739565b820191906000526020600020905b81548152906001019060200180831161071c57829003601f168201915b505050505081565b60015460009067ffffffffffffffff8316106107b9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6265796f6e6420636861696e2073697a6500000000000000000000000000000060448201526064016102b5565b60018267ffffffffffffffff16815481106107d6576107d6610bcb565b90600052602060002001549050919050565b8280546107f490610a7d565b90600052602060002090601f016020900481019282610816576000855561085c565b82601f1061082f57805160ff191683800117855561085c565b8280016001018555821561085c579182015b8281111561085c578251825591602001919060010190610841565b5061086892915061086c565b5090565b5b80821115610868576000815560010161086d565b60006020828403121561089357600080fd5b813567ffffffffffffffff811681146108ab57600080fd5b9392505050565b6000602082840312156108c457600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461091c57600080fd5b50565b80356104c9816108fa565b6000806040838503121561093d57600080fd5b823567ffffffffffffffff8082111561095557600080fd5b818501915085601f83011261096957600080fd5b81358181111561097b5761097b6108cb565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156109c1576109c16108cb565b816040528281528860208487010111156109da57600080fd5b826020860160208301376000602084830101528096505050505050610a016020840161091f565b90509250929050565b600060208083528351808285015260005b81811015610a3757858101830151858201604001528201610a1b565b81811115610a49576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b600181811c90821680610a9157607f821691505b602082108103610aca577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600060208083526000845481600182811c915080831680610af257607f831692505b8583108103610b28577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b878601838152602001818015610b455760018114610b7457610b9f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528782019650610b9f565b60008b81526020902060005b86811015610b9957815484820152908501908901610b80565b83019750505b50949998505050505050505050565b600060208284031215610bc057600080fd5b81516108ab816108fa565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c634300080d000a"

var binRuntimeChainStorageContainerStr = "0x608060405234801561001057600080fd5b506004361061007d5760003560e01c80636483ec251161005b5780636483ec25146101045780637ab4339d146101175780638da5cb5b1461012a578063ada867981461013f57600080fd5b806304f3bcec1461008257806331fe0949146100d25780635682afa9146100ef575b600080fd5b6000546100a89062010000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6001545b60405167ffffffffffffffff90911681526020016100c9565b6101026100fd366004610881565b610160565b005b6100d66101123660046108b2565b610339565b61010261012536600461092a565b6104ce565b6101326106b3565b6040516100c99190610a0a565b61015261014d366004610881565b610741565b6040519081526020016100c9565b6000546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a4478906101bc90600290600401610ad0565b602060405180830381865afa1580156101d9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101fd9190610bae565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102be57604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e60648201526084015b60405180910390fd5b60015467ffffffffffffffff82161115610334576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f63616e277420726573697a65206265796f6e6420636861696e206c656e67746860448201526064016102b5565b600155565b600080546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a44789061039690600290600401610ad0565b602060405180830381865afa1580156103b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103d79190610bae565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461049357604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e60648201526084016102b5565b5060018054808201825560008290527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601829055545b919050565b600054610100900460ff16158080156104ee5750600054600160ff909116105b806105085750303b158015610508575060005460ff166001145b610594576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016102b5565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156105f257600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b82516106059060029060208601906107e8565b50600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000073ffffffffffffffffffffffffffffffffffffffff85160217905580156106ae57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b600280546106c090610a7d565b80601f01602080910402602001604051908101604052809291908181526020018280546106ec90610a7d565b80156107395780601f1061070e57610100808354040283529160200191610739565b820191906000526020600020905b81548152906001019060200180831161071c57829003601f168201915b505050505081565b60015460009067ffffffffffffffff8316106107b9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6265796f6e6420636861696e2073697a6500000000000000000000000000000060448201526064016102b5565b60018267ffffffffffffffff16815481106107d6576107d6610bcb565b90600052602060002001549050919050565b8280546107f490610a7d565b90600052602060002090601f016020900481019282610816576000855561085c565b82601f1061082f57805160ff191683800117855561085c565b8280016001018555821561085c579182015b8281111561085c578251825591602001919060010190610841565b5061086892915061086c565b5090565b5b80821115610868576000815560010161086d565b60006020828403121561089357600080fd5b813567ffffffffffffffff811681146108ab57600080fd5b9392505050565b6000602082840312156108c457600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461091c57600080fd5b50565b80356104c9816108fa565b6000806040838503121561093d57600080fd5b823567ffffffffffffffff8082111561095557600080fd5b818501915085601f83011261096957600080fd5b81358181111561097b5761097b6108cb565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156109c1576109c16108cb565b816040528281528860208487010111156109da57600080fd5b826020860160208301376000602084830101528096505050505050610a016020840161091f565b90509250929050565b600060208083528351808285015260005b81811015610a3757858101830151858201604001528201610a1b565b81811115610a49576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b600181811c90821680610a9157607f821691505b602082108103610aca577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600060208083526000845481600182811c915080831680610af257607f831692505b8583108103610b28577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b878601838152602001818015610b455760018114610b7457610b9f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528782019650610b9f565b60008b81526020902060005b86811015610b9957815484820152908501908901610b80565b83019750505b50949998505050505050505050565b600060208284031215610bc057600080fd5b81516108ab816108fa565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c634300080d000a"

var abiChainStorageContainerStr = `[{"inputs":[{"internalType":"bytes32","name":"_element","type":"bytes32"}],"name":"append","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"chainSize","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"_index","type":"uint64"}],"name":"get","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_owner","type":"string"},{"internalType":"address","name":"_addressResolver","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"_newSize","type":"uint64"}],"name":"resize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"resolver","outputs":[{"internalType":"contract IAddressResolver","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"}]`
