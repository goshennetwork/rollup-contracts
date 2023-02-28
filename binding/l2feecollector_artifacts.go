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

var binL2FeeCollectorStr = "0x608060405234801561001057600080fd5b5061001a3361001f565b61006f565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b610a5f8061007e6000396000f3fe6080604052600436106100745760003560e01c8063a1db97821161004e578063a1db9782146100f0578063c311d04914610110578063c33af25014610130578063f2fde38b1461015057600080fd5b8063715018a61461008057806374f823ec146100975780638da5cb5b146100b757600080fd5b3661007b57005b600080fd5b34801561008c57600080fd5b50610095610170565b005b3480156100a357600080fd5b506100956100b23660046108e9565b610184565b3480156100c357600080fd5b506000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156100fc57600080fd5b5061009561010b36600461092a565b61019c565b34801561011c57600080fd5b5061009561012b366004610956565b6101d1565b34801561013c57600080fd5b5061009561014b36600461092a565b610204565b34801561015c57600080fd5b5061009561016b36600461096f565b610216565b6101786102cf565b6101826000610350565b565b61018c6102cf565b6101978383836103c5565b505050565b6101a46102cf565b6101cd826101c760005473ffffffffffffffffffffffffffffffffffffffff1690565b836103c5565b5050565b6101d96102cf565b6102016101fb60005473ffffffffffffffffffffffffffffffffffffffff1690565b82610452565b50565b61020c6102cf565b6101cd8282610452565b61021e6102cf565b73ffffffffffffffffffffffffffffffffffffffff81166102c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b61020181610350565b60005473ffffffffffffffffffffffffffffffffffffffff163314610182576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102bd565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001790526101979084906105ac565b804710156104bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a20696e73756666696369656e742062616c616e636500000060448201526064016102bd565b60008273ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d8060008114610516576040519150601f19603f3d011682016040523d82523d6000602084013e61051b565b606091505b5050905080610197576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603a60248201527f416464726573733a20756e61626c6520746f2073656e642076616c75652c207260448201527f6563697069656e74206d6179206861766520726576657274656400000000000060648201526084016102bd565b600061060e826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166106b89092919063ffffffff16565b805190915015610197578080602001905181019061062c9190610993565b610197576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016102bd565b60606106c784846000856106cf565b949350505050565b606082471015610761576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016102bd565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161078a91906109e5565b60006040518083038185875af1925050503d80600081146107c7576040519150601f19603f3d011682016040523d82523d6000602084013e6107cc565b606091505b50915091506107dd878383876107e8565b979650505050505050565b6060831561087e5782516000036108775773ffffffffffffffffffffffffffffffffffffffff85163b610877576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016102bd565b50816106c7565b6106c783838151156108935781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102bd9190610a01565b73ffffffffffffffffffffffffffffffffffffffff8116811461020157600080fd5b6000806000606084860312156108fe57600080fd5b8335610909816108c7565b92506020840135610919816108c7565b929592945050506040919091013590565b6000806040838503121561093d57600080fd5b8235610948816108c7565b946020939093013593505050565b60006020828403121561096857600080fd5b5035919050565b60006020828403121561098157600080fd5b813561098c816108c7565b9392505050565b6000602082840312156109a557600080fd5b8151801515811461098c57600080fd5b60005b838110156109d05781810151838201526020016109b8565b838111156109df576000848401525b50505050565b600082516109f78184602087016109b5565b9190910192915050565b6020815260008251806020840152610a208160408501602087016109b5565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016919091016040019291505056fea164736f6c634300080d000a"

var binRuntimeL2FeeCollectorStr = "0x6080604052600436106100745760003560e01c8063a1db97821161004e578063a1db9782146100f0578063c311d04914610110578063c33af25014610130578063f2fde38b1461015057600080fd5b8063715018a61461008057806374f823ec146100975780638da5cb5b146100b757600080fd5b3661007b57005b600080fd5b34801561008c57600080fd5b50610095610170565b005b3480156100a357600080fd5b506100956100b23660046108e9565b610184565b3480156100c357600080fd5b506000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156100fc57600080fd5b5061009561010b36600461092a565b61019c565b34801561011c57600080fd5b5061009561012b366004610956565b6101d1565b34801561013c57600080fd5b5061009561014b36600461092a565b610204565b34801561015c57600080fd5b5061009561016b36600461096f565b610216565b6101786102cf565b6101826000610350565b565b61018c6102cf565b6101978383836103c5565b505050565b6101a46102cf565b6101cd826101c760005473ffffffffffffffffffffffffffffffffffffffff1690565b836103c5565b5050565b6101d96102cf565b6102016101fb60005473ffffffffffffffffffffffffffffffffffffffff1690565b82610452565b50565b61020c6102cf565b6101cd8282610452565b61021e6102cf565b73ffffffffffffffffffffffffffffffffffffffff81166102c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b61020181610350565b60005473ffffffffffffffffffffffffffffffffffffffff163314610182576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102bd565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001790526101979084906105ac565b804710156104bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a20696e73756666696369656e742062616c616e636500000060448201526064016102bd565b60008273ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d8060008114610516576040519150601f19603f3d011682016040523d82523d6000602084013e61051b565b606091505b5050905080610197576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603a60248201527f416464726573733a20756e61626c6520746f2073656e642076616c75652c207260448201527f6563697069656e74206d6179206861766520726576657274656400000000000060648201526084016102bd565b600061060e826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166106b89092919063ffffffff16565b805190915015610197578080602001905181019061062c9190610993565b610197576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016102bd565b60606106c784846000856106cf565b949350505050565b606082471015610761576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016102bd565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161078a91906109e5565b60006040518083038185875af1925050503d80600081146107c7576040519150601f19603f3d011682016040523d82523d6000602084013e6107cc565b606091505b50915091506107dd878383876107e8565b979650505050505050565b6060831561087e5782516000036108775773ffffffffffffffffffffffffffffffffffffffff85163b610877576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016102bd565b50816106c7565b6106c783838151156108935781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102bd9190610a01565b73ffffffffffffffffffffffffffffffffffffffff8116811461020157600080fd5b6000806000606084860312156108fe57600080fd5b8335610909816108c7565b92506020840135610919816108c7565b929592945050506040919091013590565b6000806040838503121561093d57600080fd5b8235610948816108c7565b946020939093013593505050565b60006020828403121561096857600080fd5b5035919050565b60006020828403121561098157600080fd5b813561098c816108c7565b9392505050565b6000602082840312156109a557600080fd5b8151801515811461098c57600080fd5b60005b838110156109d05781810151838201526020016109b8565b838111156109df576000848401525b50505050565b600082516109f78184602087016109b5565b9190910192915050565b6020815260008251806020840152610a208160408501602087016109b5565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016919091016040019291505056fea164736f6c634300080d000a"

var abiL2FeeCollectorStr = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract IERC20","name":"token","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract IERC20","name":"token","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawERC20To","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawEth","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address payable","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"withdrawEthTo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`
