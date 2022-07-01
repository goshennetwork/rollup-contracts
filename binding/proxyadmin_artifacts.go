package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiProxyAdmin *abi.ABI

// ProxyAdminAbi returns the abi of the ProxyAdmin contract
func ProxyAdminAbi() *abi.ABI {
	return abiProxyAdmin
}

var binProxyAdmin []byte

// ProxyAdminBin returns the bin of the ProxyAdmin contract
func ProxyAdminBin() []byte {
	return binProxyAdmin
}

var binRuntimeProxyAdmin []byte

// ProxyAdminBinRuntime returns the runtime bin of the ProxyAdmin contract
func ProxyAdminBinRuntime() []byte {
	return binRuntimeProxyAdmin
}

func init() {
	var err error
	abiProxyAdmin, err = abi.NewABI(abiProxyAdminStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse ProxyAdmin abi: %v", err))
	}
	if len(binProxyAdminStr) != 0 {
		binProxyAdmin, err = hex.DecodeString(binProxyAdminStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse ProxyAdmin bin: %v", err))
		}
	}
	if len(binRuntimeProxyAdminStr) != 0 {
		binRuntimeProxyAdmin, err = hex.DecodeString(binRuntimeProxyAdminStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse ProxyAdmin bin runtime: %v", err))
		}
	}
}

var binProxyAdminStr = "0x608060405234801561001057600080fd5b5061001a3361001f565b61006f565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b610a488061007e6000396000f3fe60806040526004361061007b5760003560e01c80639623609d1161004e5780639623609d1461012b57806399a88ec41461013e578063f2fde38b1461015e578063f3b7dead1461017e57600080fd5b8063204e1c7a14610080578063715018a6146100c95780637eff275e146100e05780638da5cb5b14610100575b600080fd5b34801561008c57600080fd5b506100a061009b3660046107e4565b61019e565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b3480156100d557600080fd5b506100de610255565b005b3480156100ec57600080fd5b506100de6100fb366004610808565b6102e7565b34801561010c57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166100a0565b6100de610139366004610870565b6103ee565b34801561014a57600080fd5b506100de610159366004610808565b6104fc565b34801561016a57600080fd5b506100de6101793660046107e4565b6105d1565b34801561018a57600080fd5b506100a06101993660046107e4565b610701565b60008060008373ffffffffffffffffffffffffffffffffffffffff166040516101ea907f5c60da1b00000000000000000000000000000000000000000000000000000000815260040190565b600060405180830381855afa9150503d8060008114610225576040519150601f19603f3d011682016040523d82523d6000602084013e61022a565b606091505b50915091508161023957600080fd5b8080602001905181019061024d9190610964565b949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102db576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b6102e5600061074d565b565b60005473ffffffffffffffffffffffffffffffffffffffff163314610368576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b6040517f8f28397000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8281166004830152831690638f283970906024015b600060405180830381600087803b1580156103d257600080fd5b505af11580156103e6573d6000803e3d6000fd5b505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461046f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b6040517f4f1ef28600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841690634f1ef2869034906104c59086908690600401610981565b6000604051808303818588803b1580156104de57600080fd5b505af11580156104f2573d6000803e3d6000fd5b5050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461057d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b6040517f3659cfe600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8281166004830152831690633659cfe6906024016103b8565b60005473ffffffffffffffffffffffffffffffffffffffff163314610652576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b73ffffffffffffffffffffffffffffffffffffffff81166106f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016102d2565b6106fe8161074d565b50565b60008060008373ffffffffffffffffffffffffffffffffffffffff166040516101ea907ff851a44000000000000000000000000000000000000000000000000000000000815260040190565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b73ffffffffffffffffffffffffffffffffffffffff811681146106fe57600080fd5b6000602082840312156107f657600080fd5b8135610801816107c2565b9392505050565b6000806040838503121561081b57600080fd5b8235610826816107c2565b91506020830135610836816107c2565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008060006060848603121561088557600080fd5b8335610890816107c2565b925060208401356108a0816107c2565b9150604084013567ffffffffffffffff808211156108bd57600080fd5b818601915086601f8301126108d157600080fd5b8135818111156108e3576108e3610841565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561092957610929610841565b8160405282815289602084870101111561094257600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b60006020828403121561097657600080fd5b8151610801816107c2565b73ffffffffffffffffffffffffffffffffffffffff8316815260006020604081840152835180604085015260005b818110156109cb578581018301518582016060015282016109af565b818111156109dd576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160600194935050505056fea264697066735822122073ea7bea1ba9ae78220cbb25e63e8e8c4b4763447ea4bba4d3e3424bf9d7d39864736f6c634300080d0033"

var binRuntimeProxyAdminStr = "0x60806040526004361061007b5760003560e01c80639623609d1161004e5780639623609d1461012b57806399a88ec41461013e578063f2fde38b1461015e578063f3b7dead1461017e57600080fd5b8063204e1c7a14610080578063715018a6146100c95780637eff275e146100e05780638da5cb5b14610100575b600080fd5b34801561008c57600080fd5b506100a061009b3660046107e4565b61019e565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b3480156100d557600080fd5b506100de610255565b005b3480156100ec57600080fd5b506100de6100fb366004610808565b6102e7565b34801561010c57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166100a0565b6100de610139366004610870565b6103ee565b34801561014a57600080fd5b506100de610159366004610808565b6104fc565b34801561016a57600080fd5b506100de6101793660046107e4565b6105d1565b34801561018a57600080fd5b506100a06101993660046107e4565b610701565b60008060008373ffffffffffffffffffffffffffffffffffffffff166040516101ea907f5c60da1b00000000000000000000000000000000000000000000000000000000815260040190565b600060405180830381855afa9150503d8060008114610225576040519150601f19603f3d011682016040523d82523d6000602084013e61022a565b606091505b50915091508161023957600080fd5b8080602001905181019061024d9190610964565b949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102db576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b6102e5600061074d565b565b60005473ffffffffffffffffffffffffffffffffffffffff163314610368576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b6040517f8f28397000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8281166004830152831690638f283970906024015b600060405180830381600087803b1580156103d257600080fd5b505af11580156103e6573d6000803e3d6000fd5b505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461046f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b6040517f4f1ef28600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841690634f1ef2869034906104c59086908690600401610981565b6000604051808303818588803b1580156104de57600080fd5b505af11580156104f2573d6000803e3d6000fd5b5050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461057d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b6040517f3659cfe600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8281166004830152831690633659cfe6906024016103b8565b60005473ffffffffffffffffffffffffffffffffffffffff163314610652576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102d2565b73ffffffffffffffffffffffffffffffffffffffff81166106f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016102d2565b6106fe8161074d565b50565b60008060008373ffffffffffffffffffffffffffffffffffffffff166040516101ea907ff851a44000000000000000000000000000000000000000000000000000000000815260040190565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b73ffffffffffffffffffffffffffffffffffffffff811681146106fe57600080fd5b6000602082840312156107f657600080fd5b8135610801816107c2565b9392505050565b6000806040838503121561081b57600080fd5b8235610826816107c2565b91506020830135610836816107c2565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008060006060848603121561088557600080fd5b8335610890816107c2565b925060208401356108a0816107c2565b9150604084013567ffffffffffffffff808211156108bd57600080fd5b818601915086601f8301126108d157600080fd5b8135818111156108e3576108e3610841565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561092957610929610841565b8160405282815289602084870101111561094257600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b60006020828403121561097657600080fd5b8151610801816107c2565b73ffffffffffffffffffffffffffffffffffffffff8316815260006020604081840152835180604085015260005b818110156109cb578581018301518582016060015282016109af565b818111156109dd576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160600194935050505056fea264697066735822122073ea7bea1ba9ae78220cbb25e63e8e8c4b4763447ea4bba4d3e3424bf9d7d39864736f6c634300080d0033"

var abiProxyAdminStr = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[{"internalType":"contract TransparentUpgradeableProxy","name":"proxy","type":"address"},{"internalType":"address","name":"newAdmin","type":"address"}],"name":"changeProxyAdmin","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract TransparentUpgradeableProxy","name":"proxy","type":"address"}],"name":"getProxyAdmin","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract TransparentUpgradeableProxy","name":"proxy","type":"address"}],"name":"getProxyImplementation","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract TransparentUpgradeableProxy","name":"proxy","type":"address"},{"internalType":"address","name":"implementation","type":"address"}],"name":"upgrade","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract TransparentUpgradeableProxy","name":"proxy","type":"address"},{"internalType":"address","name":"implementation","type":"address"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"upgradeAndCall","outputs":[],"stateMutability":"payable","type":"function"}]`
