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

var binChainStorageContainerStr = "0x608060405234801561001057600080fd5b50610955806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806331fe0949146100675780635682afa9146100895780636483ec251461009e5780637ab4339d146100b15780638da5cb5b146100c4578063ada86798146100d9575b600080fd5b6001545b60405167ffffffffffffffff90911681526020015b60405180910390f35b61009c6100973660046105c3565b6100fa565b005b61006b6100ac3660046105f4565b61020b565b61009c6100bf366004610646565b6102ee565b6100cc610393565b6040516100809190610708565b6100ec6100e73660046105c3565b610421565b604051908152602001610080565b6000546040516308c3488f60e31b8152620100009091046001600160a01b03169063461a44789061013090600290600401610798565b602060405180830381865afa15801561014d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101719190610840565b6001600160a01b0316336001600160a01b0316146101aa5760405162461bcd60e51b81526004016101a19061085d565b60405180910390fd5b60015467ffffffffffffffff821611156102065760405162461bcd60e51b815260206004820181905260248201527f63616e277420726573697a65206265796f6e6420636861696e206c656e67746860448201526064016101a1565b600155565b600080546040516308c3488f60e31b8152620100009091046001600160a01b03169063461a44789061024290600290600401610798565b602060405180830381865afa15801561025f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102839190610840565b6001600160a01b0316336001600160a01b0316146102b35760405162461bcd60e51b81526004016101a19061085d565b5060018054808201825560008290527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601829055545b919050565b60006102fa60016104a2565b90508015610312576000805461ff0019166101001790555b825161032590600290602086019061052a565b506000805462010000600160b01b031916620100006001600160a01b03851602179055801561038e576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b600280546103a09061075d565b80601f01602080910402602001604051908101604052809291908181526020018280546103cc9061075d565b80156104195780601f106103ee57610100808354040283529160200191610419565b820191906000526020600020905b8154815290600101906020018083116103fc57829003601f168201915b505050505081565b60015460009067ffffffffffffffff8316106104735760405162461bcd60e51b81526020600482015260116024820152706265796f6e6420636861696e2073697a6560781b60448201526064016101a1565b60018267ffffffffffffffff1681548110610490576104906108bb565b90600052602060002001549050919050565b60008054610100900460ff16156104e9578160ff1660011480156104c55750303b155b6104e15760405162461bcd60e51b81526004016101a1906108d1565b506000919050565b60005460ff8084169116106105105760405162461bcd60e51b81526004016101a1906108d1565b506000805460ff191660ff92909216919091179055600190565b8280546105369061075d565b90600052602060002090601f016020900481019282610558576000855561059e565b82601f1061057157805160ff191683800117855561059e565b8280016001018555821561059e579182015b8281111561059e578251825591602001919060010190610583565b506105aa9291506105ae565b5090565b5b808211156105aa57600081556001016105af565b6000602082840312156105d557600080fd5b813567ffffffffffffffff811681146105ed57600080fd5b9392505050565b60006020828403121561060657600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b6001600160a01b038116811461063857600080fd5b50565b80356102e981610623565b6000806040838503121561065957600080fd5b823567ffffffffffffffff8082111561067157600080fd5b818501915085601f83011261068557600080fd5b8135818111156106975761069761060d565b604051601f8201601f19908116603f011681019083821181831017156106bf576106bf61060d565b816040528281528860208487010111156106d857600080fd5b8260208601602083013760006020848301015280965050505050506106ff6020840161063b565b90509250929050565b600060208083528351808285015260005b8181101561073557858101830151858201604001528201610719565b81811115610747576000604083870101525b50601f01601f1916929092016040019392505050565b600181811c9082168061077157607f821691505b6020821081141561079257634e487b7160e01b600052602260045260246000fd5b50919050565b600060208083526000845481600182811c9150808316806107ba57607f831692505b8583108114156107d857634e487b7160e01b85526022600452602485fd5b8786018381526020018180156107f5576001811461080657610831565b60ff19861682528782019650610831565b60008b81526020902060005b8681101561082b57815484820152908501908901610812565b83019750505b50949998505050505050505050565b60006020828403121561085257600080fd5b81516105ed81610623565b602080825260409082018190527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e20908201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e606082015260800190565b634e487b7160e01b600052603260045260246000fd5b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b60608201526080019056fea264697066735822122071d4d55aed36b79131eee43572e6bdac7b418577044f6b85db167fbac787863964736f6c634300080b0033"

var binRuntimeChainStorageContainerStr = "0x608060405234801561001057600080fd5b50600436106100625760003560e01c806331fe0949146100675780635682afa9146100895780636483ec251461009e5780637ab4339d146100b15780638da5cb5b146100c4578063ada86798146100d9575b600080fd5b6001545b60405167ffffffffffffffff90911681526020015b60405180910390f35b61009c6100973660046105c3565b6100fa565b005b61006b6100ac3660046105f4565b61020b565b61009c6100bf366004610646565b6102ee565b6100cc610393565b6040516100809190610708565b6100ec6100e73660046105c3565b610421565b604051908152602001610080565b6000546040516308c3488f60e31b8152620100009091046001600160a01b03169063461a44789061013090600290600401610798565b602060405180830381865afa15801561014d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101719190610840565b6001600160a01b0316336001600160a01b0316146101aa5760405162461bcd60e51b81526004016101a19061085d565b60405180910390fd5b60015467ffffffffffffffff821611156102065760405162461bcd60e51b815260206004820181905260248201527f63616e277420726573697a65206265796f6e6420636861696e206c656e67746860448201526064016101a1565b600155565b600080546040516308c3488f60e31b8152620100009091046001600160a01b03169063461a44789061024290600290600401610798565b602060405180830381865afa15801561025f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102839190610840565b6001600160a01b0316336001600160a01b0316146102b35760405162461bcd60e51b81526004016101a19061085d565b5060018054808201825560008290527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601829055545b919050565b60006102fa60016104a2565b90508015610312576000805461ff0019166101001790555b825161032590600290602086019061052a565b506000805462010000600160b01b031916620100006001600160a01b03851602179055801561038e576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b600280546103a09061075d565b80601f01602080910402602001604051908101604052809291908181526020018280546103cc9061075d565b80156104195780601f106103ee57610100808354040283529160200191610419565b820191906000526020600020905b8154815290600101906020018083116103fc57829003601f168201915b505050505081565b60015460009067ffffffffffffffff8316106104735760405162461bcd60e51b81526020600482015260116024820152706265796f6e6420636861696e2073697a6560781b60448201526064016101a1565b60018267ffffffffffffffff1681548110610490576104906108bb565b90600052602060002001549050919050565b60008054610100900460ff16156104e9578160ff1660011480156104c55750303b155b6104e15760405162461bcd60e51b81526004016101a1906108d1565b506000919050565b60005460ff8084169116106105105760405162461bcd60e51b81526004016101a1906108d1565b506000805460ff191660ff92909216919091179055600190565b8280546105369061075d565b90600052602060002090601f016020900481019282610558576000855561059e565b82601f1061057157805160ff191683800117855561059e565b8280016001018555821561059e579182015b8281111561059e578251825591602001919060010190610583565b506105aa9291506105ae565b5090565b5b808211156105aa57600081556001016105af565b6000602082840312156105d557600080fd5b813567ffffffffffffffff811681146105ed57600080fd5b9392505050565b60006020828403121561060657600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b6001600160a01b038116811461063857600080fd5b50565b80356102e981610623565b6000806040838503121561065957600080fd5b823567ffffffffffffffff8082111561067157600080fd5b818501915085601f83011261068557600080fd5b8135818111156106975761069761060d565b604051601f8201601f19908116603f011681019083821181831017156106bf576106bf61060d565b816040528281528860208487010111156106d857600080fd5b8260208601602083013760006020848301015280965050505050506106ff6020840161063b565b90509250929050565b600060208083528351808285015260005b8181101561073557858101830151858201604001528201610719565b81811115610747576000604083870101525b50601f01601f1916929092016040019392505050565b600181811c9082168061077157607f821691505b6020821081141561079257634e487b7160e01b600052602260045260246000fd5b50919050565b600060208083526000845481600182811c9150808316806107ba57607f831692505b8583108114156107d857634e487b7160e01b85526022600452602485fd5b8786018381526020018180156107f5576001811461080657610831565b60ff19861682528782019650610831565b60008b81526020902060005b8681101561082b57815484820152908501908901610812565b83019750505b50949998505050505050505050565b60006020828403121561085257600080fd5b81516105ed81610623565b602080825260409082018190527f436861696e53746f72616765436f6e7461696e65723a2046756e6374696f6e20908201527f63616e206f6e6c792062652063616c6c656420627920746865206f776e65722e606082015260800190565b634e487b7160e01b600052603260045260246000fd5b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b60608201526080019056fea264697066735822122071d4d55aed36b79131eee43572e6bdac7b418577044f6b85db167fbac787863964736f6c634300080b0033"

var abiChainStorageContainerStr = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"inputs":[{"internalType":"bytes32","name":"_element","type":"bytes32"}],"name":"append","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"chainSize","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"_index","type":"uint64"}],"name":"get","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_owner","type":"string"},{"internalType":"address","name":"_addressResolver","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"_newSize","type":"uint64"}],"name":"resize","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
