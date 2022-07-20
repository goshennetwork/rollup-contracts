package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiUpgradeableBeacon *abi.ABI

// UpgradeableBeaconAbi returns the abi of the UpgradeableBeacon contract
func UpgradeableBeaconAbi() *abi.ABI {
	return abiUpgradeableBeacon
}

var binUpgradeableBeacon []byte

// UpgradeableBeaconBin returns the bin of the UpgradeableBeacon contract
func UpgradeableBeaconBin() []byte {
	return binUpgradeableBeacon
}

var binRuntimeUpgradeableBeacon []byte

// UpgradeableBeaconBinRuntime returns the runtime bin of the UpgradeableBeacon contract
func UpgradeableBeaconBinRuntime() []byte {
	return binRuntimeUpgradeableBeacon
}

func init() {
	var err error
	abiUpgradeableBeacon, err = abi.NewABI(abiUpgradeableBeaconStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse UpgradeableBeacon abi: %v", err))
	}
	if len(binUpgradeableBeaconStr) != 0 {
		binUpgradeableBeacon, err = hex.DecodeString(binUpgradeableBeaconStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse UpgradeableBeacon bin: %v", err))
		}
	}
	if len(binRuntimeUpgradeableBeaconStr) != 0 {
		binRuntimeUpgradeableBeacon, err = hex.DecodeString(binRuntimeUpgradeableBeaconStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse UpgradeableBeacon bin runtime: %v", err))
		}
	}
}

var binUpgradeableBeaconStr = "0x608060405234801561001057600080fd5b506040516106e33803806106e383398101604081905261002f91610151565b61003833610047565b61004181610097565b50610181565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6100aa8161014260201b61038d1760201c565b6101205760405162461bcd60e51b815260206004820152603360248201527f5570677261646561626c65426561636f6e3a20696d706c656d656e746174696f60448201527f6e206973206e6f74206120636f6e747261637400000000000000000000000000606482015260840160405180910390fd5b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b03163b151590565b60006020828403121561016357600080fd5b81516001600160a01b038116811461017a57600080fd5b9392505050565b610553806101906000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c8063715018a611610050578063715018a6146100c45780638da5cb5b146100cc578063f2fde38b146100ea57600080fd5b80633659cfe61461006c5780635c60da1b14610081575b600080fd5b61007f61007a366004610509565b6100fd565b005b60015473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b61007f6101d0565b60005473ffffffffffffffffffffffffffffffffffffffff1661009b565b61007f6100f8366004610509565b61025d565b60005473ffffffffffffffffffffffffffffffffffffffff163314610183576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b61018c816103a9565b60405173ffffffffffffffffffffffffffffffffffffffff8216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60005473ffffffffffffffffffffffffffffffffffffffff163314610251576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161017a565b61025b6000610494565b565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161017a565b73ffffffffffffffffffffffffffffffffffffffff8116610381576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f6464726573730000000000000000000000000000000000000000000000000000606482015260840161017a565b61038a81610494565b50565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b73ffffffffffffffffffffffffffffffffffffffff81163b61044d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603360248201527f5570677261646561626c65426561636f6e3a20696d706c656d656e746174696f60448201527f6e206973206e6f74206120636f6e747261637400000000000000000000000000606482015260840161017a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60006020828403121561051b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461053f57600080fd5b939250505056fea164736f6c634300080d000a"

var binRuntimeUpgradeableBeaconStr = "0x608060405234801561001057600080fd5b50600436106100675760003560e01c8063715018a611610050578063715018a6146100c45780638da5cb5b146100cc578063f2fde38b146100ea57600080fd5b80633659cfe61461006c5780635c60da1b14610081575b600080fd5b61007f61007a366004610509565b6100fd565b005b60015473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b61007f6101d0565b60005473ffffffffffffffffffffffffffffffffffffffff1661009b565b61007f6100f8366004610509565b61025d565b60005473ffffffffffffffffffffffffffffffffffffffff163314610183576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b61018c816103a9565b60405173ffffffffffffffffffffffffffffffffffffffff8216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60005473ffffffffffffffffffffffffffffffffffffffff163314610251576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161017a565b61025b6000610494565b565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161017a565b73ffffffffffffffffffffffffffffffffffffffff8116610381576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f6464726573730000000000000000000000000000000000000000000000000000606482015260840161017a565b61038a81610494565b50565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b73ffffffffffffffffffffffffffffffffffffffff81163b61044d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603360248201527f5570677261646561626c65426561636f6e3a20696d706c656d656e746174696f60448201527f6e206973206e6f74206120636f6e747261637400000000000000000000000000606482015260840161017a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60006020828403121561051b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461053f57600080fd5b939250505056fea164736f6c634300080d000a"

var abiUpgradeableBeaconStr = `[{"inputs":[{"internalType":"address","name":"implementation_","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"implementation","type":"address"}],"name":"Upgraded","type":"event"},{"inputs":[],"name":"implementation","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newImplementation","type":"address"}],"name":"upgradeTo","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
