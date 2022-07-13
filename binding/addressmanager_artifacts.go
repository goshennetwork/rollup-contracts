package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiAddressManager *abi.ABI

// AddressManagerAbi returns the abi of the AddressManager contract
func AddressManagerAbi() *abi.ABI {
	return abiAddressManager
}

var binAddressManager []byte

// AddressManagerBin returns the bin of the AddressManager contract
func AddressManagerBin() []byte {
	return binAddressManager
}

var binRuntimeAddressManager []byte

// AddressManagerBinRuntime returns the runtime bin of the AddressManager contract
func AddressManagerBinRuntime() []byte {
	return binRuntimeAddressManager
}

func init() {
	var err error
	abiAddressManager, err = abi.NewABI(abiAddressManagerStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse AddressManager abi: %v", err))
	}
	if len(binAddressManagerStr) != 0 {
		binAddressManager, err = hex.DecodeString(binAddressManagerStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse AddressManager bin: %v", err))
		}
	}
	if len(binRuntimeAddressManagerStr) != 0 {
		binRuntimeAddressManager, err = hex.DecodeString(binRuntimeAddressManagerStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse AddressManager bin runtime: %v", err))
		}
	}
}

var binAddressManagerStr = "0x608060405234801561001057600080fd5b506112ca806100206000396000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c80638129fc1c116100cd578063d502db9711610081578063e0b45f0011610066578063e0b45f0014610243578063f10073df14610279578063f2fde38b1461028157600080fd5b8063d502db9714610228578063d9f68c641461023b57600080fd5b80638da5cb5b116100b25780638da5cb5b146101ef5780639b2ea4bd1461020d578063c64c66011461022057600080fd5b80638129fc1c146101df5780638669d0ab146101e757600080fd5b80634a7955e211610124578063715018a611610109578063715018a6146101ba57806374aee6c9146101c45780637f14099a146101cc57600080fd5b80634a7955e2146101aa5780635dbaf68b146101b257600080fd5b806322828cc214610156578063388f2a0a146101875780634162169f1461018f578063461a447814610197575b600080fd5b61015e610294565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b61015e6102d9565b61015e610319565b61015e6101a5366004610eac565b610355565b61015e610461565b61015e6104a1565b6101c26104e1565b005b61015e6104f5565b6101c26101da366004610f35565b610535565b6101c26106b6565b61015e610848565b60335473ffffffffffffffffffffffffffffffffffffffff1661015e565b6101c261021b366004610fc3565b610888565b61015e6108ec565b61015e610236366004610eac565b61092c565b61015e610968565b61015e610251366004611015565b60656020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b61015e6109a8565b6101c261028f36600461102e565b6109e8565b60006102d46040518060400160405280600e81526020017f5374616b696e674d616e61676572000000000000000000000000000000000000815250610355565b905090565b60006102d46040518060400160405280601981526020017f526f6c6c75705374617465436861696e436f6e7461696e657200000000000000815250610355565b60006102d46040518060400160405280600381526020017f44414f00000000000000000000000000000000000000000000000000000000008152505b6040517fd502db970000000000000000000000000000000000000000000000000000000081526000908190309063d502db97906103969086906004016110cc565b602060405180830381865afa1580156103b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103d791906110df565b905073ffffffffffffffffffffffffffffffffffffffff811661045b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6e6f206e616d652073617665640000000000000000000000000000000000000060448201526064015b60405180910390fd5b92915050565b60006102d46040518060400160405280601081526020017f526f6c6c75705374617465436861696e00000000000000000000000000000000815250610355565b60006102d46040518060400160405280601081526020017f4368616c6c656e6765466163746f727900000000000000000000000000000000815250610355565b6104e9610a9c565b6104f36000610b1d565b565b60006102d46040518060400160405280601081526020017f526f6c6c7570496e707574436861696e00000000000000000000000000000000815250610355565b61053d610a9c565b828181146105a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f6c656e677468206d69736d6174636800000000000000000000000000000000006044820152606401610452565b60005b818110156106ae573660008787848181106105c7576105c76110fc565b90506020028101906105d9919061112b565b9150915060008686858181106105f1576105f16110fc565b9050602002016020810190610606919061102e565b9050600061064984848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610b9492505050565b905060006106578284610bc4565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c8585838660405161068e9493929190611190565b60405180910390a1505050505080806106a690611202565b9150506105aa565b505050505050565b600054610100900460ff16158080156106d65750600054600160ff909116105b806106f05750303b1580156106f0575060005460ff166001145b61077c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610452565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156107da57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6107e2610c98565b801561084557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b60006102d46040518060400160405280601981526020017f526f6c6c7570496e707574436861696e436f6e7461696e657200000000000000815250610355565b610890610a9c565b600061089b83610b94565b905060006108a98284610bc4565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c8482856040516108de93929190611261565b60405180910390a150505050565b60006102d46040518060400160405280600f81526020017f53746174655472616e736974696f6e0000000000000000000000000000000000815250610355565b60006065600061093b84610b94565b815260208101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff1692915050565b60006102d46040518060400160405280601381526020017f4c3243726f73734c617965725769746e65737300000000000000000000000000815250610355565b60006102d46040518060400160405280601381526020017f4c3143726f73734c617965725769746e65737300000000000000000000000000815250610355565b6109f0610a9c565b73ffffffffffffffffffffffffffffffffffffffff8116610a93576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610452565b61084581610b1d565b60335473ffffffffffffffffffffffffffffffffffffffff1633146104f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610452565b6033805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600081604051602001610ba791906112a1565b604051602081830303815290604052805190602001209050919050565b600073ffffffffffffffffffffffffffffffffffffffff8216610c43576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600a60248201527f656d7074792061646472000000000000000000000000000000000000000000006044820152606401610452565b5060009182526065602052604090912080547fffffffffffffffffffffffff0000000000000000000000000000000000000000811673ffffffffffffffffffffffffffffffffffffffff938416179091551690565b600054610100900460ff16610d2f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610452565b6104f3600054610100900460ff16610dc9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610452565b6104f333610b1d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f830112610e1257600080fd5b813567ffffffffffffffff80821115610e2d57610e2d610dd2565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715610e7357610e73610dd2565b81604052838152866020858801011115610e8c57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060208284031215610ebe57600080fd5b813567ffffffffffffffff811115610ed557600080fd5b610ee184828501610e01565b949350505050565b60008083601f840112610efb57600080fd5b50813567ffffffffffffffff811115610f1357600080fd5b6020830191508360208260051b8501011115610f2e57600080fd5b9250929050565b60008060008060408587031215610f4b57600080fd5b843567ffffffffffffffff80821115610f6357600080fd5b610f6f88838901610ee9565b90965094506020870135915080821115610f8857600080fd5b50610f9587828801610ee9565b95989497509550505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461084557600080fd5b60008060408385031215610fd657600080fd5b823567ffffffffffffffff811115610fed57600080fd5b610ff985828601610e01565b925050602083013561100a81610fa1565b809150509250929050565b60006020828403121561102757600080fd5b5035919050565b60006020828403121561104057600080fd5b813561104b81610fa1565b9392505050565b60005b8381101561106d578181015183820152602001611055565b8381111561107c576000848401525b50505050565b6000815180845261109a816020860160208601611052565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061104b6020830184611082565b6000602082840312156110f157600080fd5b815161104b81610fa1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261116057600080fd5b83018035915067ffffffffffffffff82111561117b57600080fd5b602001915036819003821315610f2e57600080fd5b60608152836060820152838560808301376000608085830181019190915273ffffffffffffffffffffffffffffffffffffffff9384166020830152919092166040830152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361125a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b6060815260006112746060830186611082565b73ffffffffffffffffffffffffffffffffffffffff94851660208401529290931660409091015292915050565b600082516112b3818460208701611052565b919091019291505056fea164736f6c634300080d000a"

var binRuntimeAddressManagerStr = "0x608060405234801561001057600080fd5b50600436106101515760003560e01c80638129fc1c116100cd578063d502db9711610081578063e0b45f0011610066578063e0b45f0014610243578063f10073df14610279578063f2fde38b1461028157600080fd5b8063d502db9714610228578063d9f68c641461023b57600080fd5b80638da5cb5b116100b25780638da5cb5b146101ef5780639b2ea4bd1461020d578063c64c66011461022057600080fd5b80638129fc1c146101df5780638669d0ab146101e757600080fd5b80634a7955e211610124578063715018a611610109578063715018a6146101ba57806374aee6c9146101c45780637f14099a146101cc57600080fd5b80634a7955e2146101aa5780635dbaf68b146101b257600080fd5b806322828cc214610156578063388f2a0a146101875780634162169f1461018f578063461a447814610197575b600080fd5b61015e610294565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b61015e6102d9565b61015e610319565b61015e6101a5366004610eac565b610355565b61015e610461565b61015e6104a1565b6101c26104e1565b005b61015e6104f5565b6101c26101da366004610f35565b610535565b6101c26106b6565b61015e610848565b60335473ffffffffffffffffffffffffffffffffffffffff1661015e565b6101c261021b366004610fc3565b610888565b61015e6108ec565b61015e610236366004610eac565b61092c565b61015e610968565b61015e610251366004611015565b60656020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b61015e6109a8565b6101c261028f36600461102e565b6109e8565b60006102d46040518060400160405280600e81526020017f5374616b696e674d616e61676572000000000000000000000000000000000000815250610355565b905090565b60006102d46040518060400160405280601981526020017f526f6c6c75705374617465436861696e436f6e7461696e657200000000000000815250610355565b60006102d46040518060400160405280600381526020017f44414f00000000000000000000000000000000000000000000000000000000008152505b6040517fd502db970000000000000000000000000000000000000000000000000000000081526000908190309063d502db97906103969086906004016110cc565b602060405180830381865afa1580156103b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103d791906110df565b905073ffffffffffffffffffffffffffffffffffffffff811661045b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6e6f206e616d652073617665640000000000000000000000000000000000000060448201526064015b60405180910390fd5b92915050565b60006102d46040518060400160405280601081526020017f526f6c6c75705374617465436861696e00000000000000000000000000000000815250610355565b60006102d46040518060400160405280601081526020017f4368616c6c656e6765466163746f727900000000000000000000000000000000815250610355565b6104e9610a9c565b6104f36000610b1d565b565b60006102d46040518060400160405280601081526020017f526f6c6c7570496e707574436861696e00000000000000000000000000000000815250610355565b61053d610a9c565b828181146105a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f6c656e677468206d69736d6174636800000000000000000000000000000000006044820152606401610452565b60005b818110156106ae573660008787848181106105c7576105c76110fc565b90506020028101906105d9919061112b565b9150915060008686858181106105f1576105f16110fc565b9050602002016020810190610606919061102e565b9050600061064984848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610b9492505050565b905060006106578284610bc4565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c8585838660405161068e9493929190611190565b60405180910390a1505050505080806106a690611202565b9150506105aa565b505050505050565b600054610100900460ff16158080156106d65750600054600160ff909116105b806106f05750303b1580156106f0575060005460ff166001145b61077c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610452565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156107da57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6107e2610c98565b801561084557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b60006102d46040518060400160405280601981526020017f526f6c6c7570496e707574436861696e436f6e7461696e657200000000000000815250610355565b610890610a9c565b600061089b83610b94565b905060006108a98284610bc4565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c8482856040516108de93929190611261565b60405180910390a150505050565b60006102d46040518060400160405280600f81526020017f53746174655472616e736974696f6e0000000000000000000000000000000000815250610355565b60006065600061093b84610b94565b815260208101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff1692915050565b60006102d46040518060400160405280601381526020017f4c3243726f73734c617965725769746e65737300000000000000000000000000815250610355565b60006102d46040518060400160405280601381526020017f4c3143726f73734c617965725769746e65737300000000000000000000000000815250610355565b6109f0610a9c565b73ffffffffffffffffffffffffffffffffffffffff8116610a93576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610452565b61084581610b1d565b60335473ffffffffffffffffffffffffffffffffffffffff1633146104f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610452565b6033805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600081604051602001610ba791906112a1565b604051602081830303815290604052805190602001209050919050565b600073ffffffffffffffffffffffffffffffffffffffff8216610c43576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600a60248201527f656d7074792061646472000000000000000000000000000000000000000000006044820152606401610452565b5060009182526065602052604090912080547fffffffffffffffffffffffff0000000000000000000000000000000000000000811673ffffffffffffffffffffffffffffffffffffffff938416179091551690565b600054610100900460ff16610d2f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610452565b6104f3600054610100900460ff16610dc9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610452565b6104f333610b1d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f830112610e1257600080fd5b813567ffffffffffffffff80821115610e2d57610e2d610dd2565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715610e7357610e73610dd2565b81604052838152866020858801011115610e8c57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060208284031215610ebe57600080fd5b813567ffffffffffffffff811115610ed557600080fd5b610ee184828501610e01565b949350505050565b60008083601f840112610efb57600080fd5b50813567ffffffffffffffff811115610f1357600080fd5b6020830191508360208260051b8501011115610f2e57600080fd5b9250929050565b60008060008060408587031215610f4b57600080fd5b843567ffffffffffffffff80821115610f6357600080fd5b610f6f88838901610ee9565b90965094506020870135915080821115610f8857600080fd5b50610f9587828801610ee9565b95989497509550505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461084557600080fd5b60008060408385031215610fd657600080fd5b823567ffffffffffffffff811115610fed57600080fd5b610ff985828601610e01565b925050602083013561100a81610fa1565b809150509250929050565b60006020828403121561102757600080fd5b5035919050565b60006020828403121561104057600080fd5b813561104b81610fa1565b9392505050565b60005b8381101561106d578181015183820152602001611055565b8381111561107c576000848401525b50505050565b6000815180845261109a816020860160208601611052565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061104b6020830184611082565b6000602082840312156110f157600080fd5b815161104b81610fa1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261116057600080fd5b83018035915067ffffffffffffffff82111561117b57600080fd5b602001915036819003821315610f2e57600080fd5b60608152836060820152838560808301376000608085830181019190915273ffffffffffffffffffffffffffffffffffffffff9384166020830152919092166040830152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361125a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b6060815260006112746060830186611082565b73ffffffffffffffffffffffffffffffffffffffff94851660208401529290931660409091015292915050565b600082516112b3818460208701611052565b919091019291505056fea164736f6c634300080d000a"

var abiAddressManagerStr = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"_name","type":"string"},{"indexed":false,"internalType":"address","name":"_old","type":"address"},{"indexed":false,"internalType":"address","name":"_new","type":"address"}],"name":"AddressSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"challengeFactory","outputs":[{"internalType":"contract IChallengeFactory","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"dao","outputs":[{"internalType":"contract IDAO","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"name":"getAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"getAddrByHash","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"l1CrossLayerWitness","outputs":[{"internalType":"contract IL1CrossLayerWitness","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"l2CrossLayerWitness","outputs":[{"internalType":"contract IL2CrossLayerWitness","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"name":"resolve","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupInputChain","outputs":[{"internalType":"contract IRollupInputChain","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupInputChainContainer","outputs":[{"internalType":"contract IChainStorageContainer","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupStateChain","outputs":[{"internalType":"contract IRollupStateChain","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupStateChainContainer","outputs":[{"internalType":"contract IChainStorageContainer","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"},{"internalType":"address","name":"_addr","type":"address"}],"name":"setAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string[]","name":"_names","type":"string[]"},{"internalType":"address[]","name":"_addrs","type":"address[]"}],"name":"setAddressBatch","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"stakingManager","outputs":[{"internalType":"contract IStakingManager","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"stateTransition","outputs":[{"internalType":"contract IStateTransition","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
