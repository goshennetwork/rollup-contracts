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

var binChainStorageContainerStr = "0x608060405234801561001057600080fd5b50610bb1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c80637ab4339d116100505780637ab4339d146100c15780638da5cb5b146100d4578063ada86798146100e957600080fd5b806331fe0949146100775780635682afa9146100995780636483ec25146100ae575b600080fd5b6001545b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100ac6100a736600461082b565b61010a565b005b61007b6100bc36600461085c565b6102e3565b6100ac6100cf3660046108d4565b610478565b6100dc61065d565b60405161009091906109b4565b6100fc6100f736600461082b565b6106eb565b604051908152602001610090565b6000546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a44789061016690600290600401610a7a565b602060405180830381865afa158015610183573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101a79190610b58565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461026857604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e60648201526084015b60405180910390fd5b60015467ffffffffffffffff821611156102de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f63616e277420726573697a65206265796f6e6420636861696e206c656e677468604482015260640161025f565b600155565b600080546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a44789061034090600290600401610a7a565b602060405180830381865afa15801561035d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103819190610b58565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461043d57604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e606482015260840161025f565b5060018054808201825560008290527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601829055545b919050565b600054610100900460ff16158080156104985750600054600160ff909116105b806104b25750303b1580156104b2575060005460ff166001145b61053e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161025f565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561059c57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b82516105af906002906020860190610792565b50600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000073ffffffffffffffffffffffffffffffffffffffff851602179055801561065857600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b6002805461066a90610a27565b80601f016020809104026020016040519081016040528092919081815260200182805461069690610a27565b80156106e35780601f106106b8576101008083540402835291602001916106e3565b820191906000526020600020905b8154815290600101906020018083116106c657829003601f168201915b505050505081565b60015460009067ffffffffffffffff831610610763576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6265796f6e6420636861696e2073697a65000000000000000000000000000000604482015260640161025f565b60018267ffffffffffffffff168154811061078057610780610b75565b90600052602060002001549050919050565b82805461079e90610a27565b90600052602060002090601f0160209004810192826107c05760008555610806565b82601f106107d957805160ff1916838001178555610806565b82800160010185558215610806579182015b828111156108065782518255916020019190600101906107eb565b50610812929150610816565b5090565b5b808211156108125760008155600101610817565b60006020828403121561083d57600080fd5b813567ffffffffffffffff8116811461085557600080fd5b9392505050565b60006020828403121561086e57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff811681146108c657600080fd5b50565b8035610473816108a4565b600080604083850312156108e757600080fd5b823567ffffffffffffffff808211156108ff57600080fd5b818501915085601f83011261091357600080fd5b81358181111561092557610925610875565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561096b5761096b610875565b8160405282815288602084870101111561098457600080fd5b8260208601602083013760006020848301015280965050505050506109ab602084016108c9565b90509250929050565b600060208083528351808285015260005b818110156109e1578581018301518582016040015282016109c5565b818111156109f3576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b600181811c90821680610a3b57607f821691505b602082108103610a74577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600060208083526000845481600182811c915080831680610a9c57607f831692505b8583108103610ad2577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b878601838152602001818015610aef5760018114610b1e57610b49565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528782019650610b49565b60008b81526020902060005b86811015610b4357815484820152908501908901610b2a565b83019750505b50949998505050505050505050565b600060208284031215610b6a57600080fd5b8151610855816108a4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c634300080d000a"

var binRuntimeChainStorageContainerStr = "0x608060405234801561001057600080fd5b50600436106100725760003560e01c80637ab4339d116100505780637ab4339d146100c15780638da5cb5b146100d4578063ada86798146100e957600080fd5b806331fe0949146100775780635682afa9146100995780636483ec25146100ae575b600080fd5b6001545b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100ac6100a736600461082b565b61010a565b005b61007b6100bc36600461085c565b6102e3565b6100ac6100cf3660046108d4565b610478565b6100dc61065d565b60405161009091906109b4565b6100fc6100f736600461082b565b6106eb565b604051908152602001610090565b6000546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a44789061016690600290600401610a7a565b602060405180830381865afa158015610183573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101a79190610b58565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461026857604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e60648201526084015b60405180910390fd5b60015467ffffffffffffffff821611156102de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f63616e277420726573697a65206265796f6e6420636861696e206c656e677468604482015260640161025f565b600155565b600080546040517f461a44780000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063461a44789061034090600290600401610a7a565b602060405180830381865afa15801561035d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103819190610b58565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461043d57604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e2060448201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e606482015260840161025f565b5060018054808201825560008290527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601829055545b919050565b600054610100900460ff16158080156104985750600054600160ff909116105b806104b25750303b1580156104b2575060005460ff166001145b61053e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161025f565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561059c57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b82516105af906002906020860190610792565b50600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000073ffffffffffffffffffffffffffffffffffffffff851602179055801561065857600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b6002805461066a90610a27565b80601f016020809104026020016040519081016040528092919081815260200182805461069690610a27565b80156106e35780601f106106b8576101008083540402835291602001916106e3565b820191906000526020600020905b8154815290600101906020018083116106c657829003601f168201915b505050505081565b60015460009067ffffffffffffffff831610610763576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6265796f6e6420636861696e2073697a65000000000000000000000000000000604482015260640161025f565b60018267ffffffffffffffff168154811061078057610780610b75565b90600052602060002001549050919050565b82805461079e90610a27565b90600052602060002090601f0160209004810192826107c05760008555610806565b82601f106107d957805160ff1916838001178555610806565b82800160010185558215610806579182015b828111156108065782518255916020019190600101906107eb565b50610812929150610816565b5090565b5b808211156108125760008155600101610817565b60006020828403121561083d57600080fd5b813567ffffffffffffffff8116811461085557600080fd5b9392505050565b60006020828403121561086e57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff811681146108c657600080fd5b50565b8035610473816108a4565b600080604083850312156108e757600080fd5b823567ffffffffffffffff808211156108ff57600080fd5b818501915085601f83011261091357600080fd5b81358181111561092557610925610875565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561096b5761096b610875565b8160405282815288602084870101111561098457600080fd5b8260208601602083013760006020848301015280965050505050506109ab602084016108c9565b90509250929050565b600060208083528351808285015260005b818110156109e1578581018301518582016040015282016109c5565b818111156109f3576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b600181811c90821680610a3b57607f821691505b602082108103610a74577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600060208083526000845481600182811c915080831680610a9c57607f831692505b8583108103610ad2577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b878601838152602001818015610aef5760018114610b1e57610b49565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528782019650610b49565b60008b81526020902060005b86811015610b4357815484820152908501908901610b2a565b83019750505b50949998505050505050505050565b600060208284031215610b6a57600080fd5b8151610855816108a4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c634300080d000a"

var abiChainStorageContainerStr = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"inputs":[{"internalType":"bytes32","name":"_element","type":"bytes32"}],"name":"append","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"chainSize","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"_index","type":"uint64"}],"name":"get","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_owner","type":"string"},{"internalType":"address","name":"_addressResolver","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"_newSize","type":"uint64"}],"name":"resize","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
