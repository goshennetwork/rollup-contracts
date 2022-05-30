package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiL2FeeCollector *abi.ABI

// L2FeeCollectorAbi returns the abi of the L2FeeCollector contract
func L2FeeCollectorAbi() *abi.ABI {
	return abiL2FeeCollector
}

var binL2FeeCollector []byte

// L2FeeCollectorBin returns the bin of the L2FeeCollector contract
func L2FeeCollectorBin() []byte {
	return binL2FeeCollector
}

var binRuntimeL2FeeCollector []byte

// L2FeeCollectorBinRuntime returns the runtime bin of the L2FeeCollector contract
func L2FeeCollectorBinRuntime() []byte {
	return binRuntimeL2FeeCollector
}

func init() {
	var err error
	abiL2FeeCollector, err = abi.NewABI(abiL2FeeCollectorStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse L2FeeCollector abi: %v", err))
	}
	if len(binL2FeeCollectorStr) != 0 {
		binL2FeeCollector, err = hex.DecodeString(binL2FeeCollectorStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse L2FeeCollector bin: %v", err))
		}
	}
	if len(binRuntimeL2FeeCollectorStr) != 0 {
		binRuntimeL2FeeCollector, err = hex.DecodeString(binRuntimeL2FeeCollectorStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse L2FeeCollector bin runtime: %v", err))
		}
	}
}

var binL2FeeCollectorStr = "0x608060405234801561001057600080fd5b5061001a3361001f565b61006f565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b61090a8061007e6000396000f3fe6080604052600436106100745760003560e01c8063a1db97821161004e578063a1db9782146100e3578063c311d04914610103578063c33af25014610123578063f2fde38b1461014357600080fd5b8063715018a61461008057806374f823ec146100975780638da5cb5b146100b757600080fd5b3661007b57005b600080fd5b34801561008c57600080fd5b50610095610163565b005b3480156100a357600080fd5b506100956100b236600461075b565b6101a2565b3480156100c357600080fd5b50600054604080516001600160a01b039092168252519081900360200190f35b3480156100ef57600080fd5b506100956100fe36600461079c565b6101dc565b34801561010f57600080fd5b5061009561011e3660046107c8565b610226565b34801561012f57600080fd5b5061009561013e36600461079c565b61026e565b34801561014f57600080fd5b5061009561015e3660046107e1565b6102a2565b6000546001600160a01b031633146101965760405162461bcd60e51b815260040161018d906107fe565b60405180910390fd5b6101a06000610336565b565b6000546001600160a01b031633146101cc5760405162461bcd60e51b815260040161018d906107fe565b6101d7838383610386565b505050565b6000546001600160a01b031633146102065760405162461bcd60e51b815260040161018d906107fe565b6102228261021c6000546001600160a01b031690565b83610386565b5050565b6000546001600160a01b031633146102505760405162461bcd60e51b815260040161018d906107fe565b61026b6102656000546001600160a01b031690565b826103d8565b50565b6000546001600160a01b031633146102985760405162461bcd60e51b815260040161018d906107fe565b61022282826103d8565b6000546001600160a01b031633146102cc5760405162461bcd60e51b815260040161018d906107fe565b6001600160a01b0381166103315760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161018d565b61026b815b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b031663a9059cbb60e01b1790526101d79084906104f1565b804710156104285760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a20696e73756666696369656e742062616c616e6365000000604482015260640161018d565b6000826001600160a01b03168260405160006040518083038185875af1925050503d8060008114610475576040519150601f19603f3d011682016040523d82523d6000602084013e61047a565b606091505b50509050806101d75760405162461bcd60e51b815260206004820152603a60248201527f416464726573733a20756e61626c6520746f2073656e642076616c75652c207260448201527f6563697069656e74206d61792068617665207265766572746564000000000000606482015260840161018d565b6000610546826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166105c39092919063ffffffff16565b8051909150156101d757808060200190518101906105649190610833565b6101d75760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840161018d565b60606105d284846000856105dc565b90505b9392505050565b60608247101561063d5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840161018d565b6001600160a01b0385163b6106945760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161018d565b600080866001600160a01b031685876040516106b09190610885565b60006040518083038185875af1925050503d80600081146106ed576040519150601f19603f3d011682016040523d82523d6000602084013e6106f2565b606091505b509150915061070282828661070d565b979650505050505050565b6060831561071c5750816105d5565b82511561072c5782518084602001fd5b8160405162461bcd60e51b815260040161018d91906108a1565b6001600160a01b038116811461026b57600080fd5b60008060006060848603121561077057600080fd5b833561077b81610746565b9250602084013561078b81610746565b929592945050506040919091013590565b600080604083850312156107af57600080fd5b82356107ba81610746565b946020939093013593505050565b6000602082840312156107da57600080fd5b5035919050565b6000602082840312156107f357600080fd5b81356105d581610746565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b60006020828403121561084557600080fd5b815180151581146105d557600080fd5b60005b83811015610870578181015183820152602001610858565b8381111561087f576000848401525b50505050565b60008251610897818460208701610855565b9190910192915050565b60208152600082518060208401526108c0816040850160208701610855565b601f01601f1916919091016040019291505056fea2646970667358221220afdbee5f3523b5dfa5413fba1510d5eed88089691cdc53c889dde11669cb949464736f6c634300080e0033"

var binRuntimeL2FeeCollectorStr = "0x6080604052600436106100745760003560e01c8063a1db97821161004e578063a1db9782146100e3578063c311d04914610103578063c33af25014610123578063f2fde38b1461014357600080fd5b8063715018a61461008057806374f823ec146100975780638da5cb5b146100b757600080fd5b3661007b57005b600080fd5b34801561008c57600080fd5b50610095610163565b005b3480156100a357600080fd5b506100956100b236600461075b565b6101a2565b3480156100c357600080fd5b50600054604080516001600160a01b039092168252519081900360200190f35b3480156100ef57600080fd5b506100956100fe36600461079c565b6101dc565b34801561010f57600080fd5b5061009561011e3660046107c8565b610226565b34801561012f57600080fd5b5061009561013e36600461079c565b61026e565b34801561014f57600080fd5b5061009561015e3660046107e1565b6102a2565b6000546001600160a01b031633146101965760405162461bcd60e51b815260040161018d906107fe565b60405180910390fd5b6101a06000610336565b565b6000546001600160a01b031633146101cc5760405162461bcd60e51b815260040161018d906107fe565b6101d7838383610386565b505050565b6000546001600160a01b031633146102065760405162461bcd60e51b815260040161018d906107fe565b6102228261021c6000546001600160a01b031690565b83610386565b5050565b6000546001600160a01b031633146102505760405162461bcd60e51b815260040161018d906107fe565b61026b6102656000546001600160a01b031690565b826103d8565b50565b6000546001600160a01b031633146102985760405162461bcd60e51b815260040161018d906107fe565b61022282826103d8565b6000546001600160a01b031633146102cc5760405162461bcd60e51b815260040161018d906107fe565b6001600160a01b0381166103315760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161018d565b61026b815b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b031663a9059cbb60e01b1790526101d79084906104f1565b804710156104285760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a20696e73756666696369656e742062616c616e6365000000604482015260640161018d565b6000826001600160a01b03168260405160006040518083038185875af1925050503d8060008114610475576040519150601f19603f3d011682016040523d82523d6000602084013e61047a565b606091505b50509050806101d75760405162461bcd60e51b815260206004820152603a60248201527f416464726573733a20756e61626c6520746f2073656e642076616c75652c207260448201527f6563697069656e74206d61792068617665207265766572746564000000000000606482015260840161018d565b6000610546826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166105c39092919063ffffffff16565b8051909150156101d757808060200190518101906105649190610833565b6101d75760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840161018d565b60606105d284846000856105dc565b90505b9392505050565b60608247101561063d5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840161018d565b6001600160a01b0385163b6106945760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161018d565b600080866001600160a01b031685876040516106b09190610885565b60006040518083038185875af1925050503d80600081146106ed576040519150601f19603f3d011682016040523d82523d6000602084013e6106f2565b606091505b509150915061070282828661070d565b979650505050505050565b6060831561071c5750816105d5565b82511561072c5782518084602001fd5b8160405162461bcd60e51b815260040161018d91906108a1565b6001600160a01b038116811461026b57600080fd5b60008060006060848603121561077057600080fd5b833561077b81610746565b9250602084013561078b81610746565b929592945050506040919091013590565b600080604083850312156107af57600080fd5b82356107ba81610746565b946020939093013593505050565b6000602082840312156107da57600080fd5b5035919050565b6000602082840312156107f357600080fd5b81356105d581610746565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b60006020828403121561084557600080fd5b815180151581146105d557600080fd5b60005b83811015610870578181015183820152602001610858565b8381111561087f576000848401525b50505050565b60008251610897818460208701610855565b9190910192915050565b60208152600082518060208401526108c0816040850160208701610855565b601f01601f1916919091016040019291505056fea2646970667358221220afdbee5f3523b5dfa5413fba1510d5eed88089691cdc53c889dde11669cb949464736f6c634300080e0033"

var abiL2FeeCollectorStr = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract IERC20","name":"token","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract IERC20","name":"token","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawERC20To","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawEth","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address payable","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawEthTo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`
