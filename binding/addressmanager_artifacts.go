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

var binAddressManagerStr = "0x608060405234801561001057600080fd5b50610f36806100206000396000f3fe608060405234801561001057600080fd5b50600436106101215760003560e01c80638129fc1c116100ad578063d502db9711610071578063d502db97146101de578063d9f68c64146101f1578063e0b45f00146101f9578063f10073df14610222578063f2fde38b1461022a57600080fd5b80638129fc1c146101a25780638669d0ab146101aa5780638da5cb5b146101b25780639b2ea4bd146101c3578063c64c6601146101d657600080fd5b80634a7955e2116100f45780634a7955e21461016d5780635dbaf68b14610175578063715018a61461017d57806374aee6c9146101875780637f14099a1461018f57600080fd5b806322828cc214610126578063388f2a0a1461014a5780634162169f14610152578063461a44781461015a575b600080fd5b61012e61023d565b6040516001600160a01b03909116815260200160405180910390f35b61012e610273565b61012e6102b3565b61012e610168366004610af3565b6102d5565b61012e610391565b61012e6103c4565b6101856103f7565b005b61012e61042d565b61018561019d366004610b7c565b610460565b6101856105db565b61012e610650565b6033546001600160a01b031661012e565b6101856101d1366004610bfd565b610690565b61012e610716565b61012e6101ec366004610af3565b610748565b61012e610777565b61012e610207366004610c4f565b6065602052600090815260409020546001600160a01b031681565b61012e6107ad565b610185610238366004610c68565b6107e3565b600061026e6040518060400160405280600e81526020016d29ba30b5b4b733a6b0b730b3b2b960911b8152506102d5565b905090565b600061026e6040518060400160405280601981526020017f526f6c6c75705374617465436861696e436f6e7461696e6572000000000000008152506102d5565b600061026e6040518060400160405280600381526020016244414f60e81b8152505b60405163d502db9760e01b81526000908190309063d502db97906102fd908690600401610ce8565b602060405180830381865afa15801561031a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061033e9190610cfb565b90506001600160a01b03811661038b5760405162461bcd60e51b815260206004820152600d60248201526c1b9bc81b985b59481cd85d9959609a1b60448201526064015b60405180910390fd5b92915050565b600061026e6040518060400160405280601081526020016f2937b6363ab829ba30ba32a1b430b4b760811b8152506102d5565b600061026e6040518060400160405280601081526020016f4368616c6c656e6765466163746f727960801b8152506102d5565b6033546001600160a01b031633146104215760405162461bcd60e51b815260040161038290610d18565b61042b6000610877565b565b600061026e6040518060400160405280601081526020016f2937b6363ab824b7383aba21b430b4b760811b8152506102d5565b6033546001600160a01b0316331461048a5760405162461bcd60e51b815260040161038290610d18565b828181146104cc5760405162461bcd60e51b815260206004820152600f60248201526e0d8cadccee8d040dad2e6dac2e8c6d608b1b6044820152606401610382565b60005b818110156105d3573660008787848181106104ec576104ec610d4d565b90506020028101906104fe9190610d63565b91509150600086868581811061051657610516610d4d565b905060200201602081019061052b9190610c68565b9050600061056e84848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506108c992505050565b9050600061057c82846108f9565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c858583866040516105b39493929190610daa565b60405180910390a1505050505080806105cb90610df1565b9150506104cf565b505050505050565b60006105e7600161096e565b905080156105ff576000805461ff0019166101001790555b6106076109f6565b801561064d576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b600061026e6040518060400160405280601981526020017f526f6c6c7570496e707574436861696e436f6e7461696e6572000000000000008152506102d5565b6033546001600160a01b031633146106ba5760405162461bcd60e51b815260040161038290610d18565b60006106c5836108c9565b905060006106d382846108f9565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c84828560405161070893929190610e18565b60405180910390a150505050565b600061026e6040518060400160405280600f81526020016e29ba30ba32aa3930b739b4ba34b7b760891b8152506102d5565b600060656000610757846108c9565b81526020810191909152604001600020546001600160a01b031692915050565b600061026e604051806040016040528060138152602001724c3243726f73734c617965725769746e65737360681b8152506102d5565b600061026e604051806040016040528060138152602001724c3143726f73734c617965725769746e65737360681b8152506102d5565b6033546001600160a01b0316331461080d5760405162461bcd60e51b815260040161038290610d18565b6001600160a01b0381166108725760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610382565b61064d815b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6000816040516020016108dc9190610e4b565b604051602081830303815290604052805190602001209050919050565b60006001600160a01b03821661093e5760405162461bcd60e51b815260206004820152600a60248201526932b6b83a3c9030b2323960b11b6044820152606401610382565b5060009182526065602052604090912080546001600160a01b031981166001600160a01b03938416179091551690565b60008054610100900460ff16156109b5578160ff1660011480156109915750303b155b6109ad5760405162461bcd60e51b815260040161038290610e67565b506000919050565b60005460ff8084169116106109dc5760405162461bcd60e51b815260040161038290610e67565b506000805460ff191660ff92909216919091179055600190565b600054610100900460ff16610a1d5760405162461bcd60e51b815260040161038290610eb5565b61042b600054610100900460ff16610a475760405162461bcd60e51b815260040161038290610eb5565b61042b33610877565b634e487b7160e01b600052604160045260246000fd5b600082601f830112610a7757600080fd5b813567ffffffffffffffff80821115610a9257610a92610a50565b604051601f8301601f19908116603f01168101908282118183101715610aba57610aba610a50565b81604052838152866020858801011115610ad357600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060208284031215610b0557600080fd5b813567ffffffffffffffff811115610b1c57600080fd5b610b2884828501610a66565b949350505050565b60008083601f840112610b4257600080fd5b50813567ffffffffffffffff811115610b5a57600080fd5b6020830191508360208260051b8501011115610b7557600080fd5b9250929050565b60008060008060408587031215610b9257600080fd5b843567ffffffffffffffff80821115610baa57600080fd5b610bb688838901610b30565b90965094506020870135915080821115610bcf57600080fd5b50610bdc87828801610b30565b95989497509550505050565b6001600160a01b038116811461064d57600080fd5b60008060408385031215610c1057600080fd5b823567ffffffffffffffff811115610c2757600080fd5b610c3385828601610a66565b9250506020830135610c4481610be8565b809150509250929050565b600060208284031215610c6157600080fd5b5035919050565b600060208284031215610c7a57600080fd5b8135610c8581610be8565b9392505050565b60005b83811015610ca7578181015183820152602001610c8f565b83811115610cb6576000848401525b50505050565b60008151808452610cd4816020860160208601610c8c565b601f01601f19169290920160200192915050565b602081526000610c856020830184610cbc565b600060208284031215610d0d57600080fd5b8151610c8581610be8565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b634e487b7160e01b600052603260045260246000fd5b6000808335601e19843603018112610d7a57600080fd5b83018035915067ffffffffffffffff821115610d9557600080fd5b602001915036819003821315610b7557600080fd5b6060815283606082015283856080830137600060808583018101919091526001600160a01b039384166020830152919092166040830152601f909201601f19160101919050565b600060018201610e1157634e487b7160e01b600052601160045260246000fd5b5060010190565b606081526000610e2b6060830186610cbc565b6001600160a01b0394851660208401529290931660409091015292915050565b60008251610e5d818460208701610c8c565b9190910192915050565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b60608201526080019056fea2646970667358221220ec0c6358a940350c4ef94e1cfeeba46442ad008f0b948371644b89de7542771364736f6c634300080e0033"

var binRuntimeAddressManagerStr = "0x608060405234801561001057600080fd5b50600436106101215760003560e01c80638129fc1c116100ad578063d502db9711610071578063d502db97146101de578063d9f68c64146101f1578063e0b45f00146101f9578063f10073df14610222578063f2fde38b1461022a57600080fd5b80638129fc1c146101a25780638669d0ab146101aa5780638da5cb5b146101b25780639b2ea4bd146101c3578063c64c6601146101d657600080fd5b80634a7955e2116100f45780634a7955e21461016d5780635dbaf68b14610175578063715018a61461017d57806374aee6c9146101875780637f14099a1461018f57600080fd5b806322828cc214610126578063388f2a0a1461014a5780634162169f14610152578063461a44781461015a575b600080fd5b61012e61023d565b6040516001600160a01b03909116815260200160405180910390f35b61012e610273565b61012e6102b3565b61012e610168366004610af3565b6102d5565b61012e610391565b61012e6103c4565b6101856103f7565b005b61012e61042d565b61018561019d366004610b7c565b610460565b6101856105db565b61012e610650565b6033546001600160a01b031661012e565b6101856101d1366004610bfd565b610690565b61012e610716565b61012e6101ec366004610af3565b610748565b61012e610777565b61012e610207366004610c4f565b6065602052600090815260409020546001600160a01b031681565b61012e6107ad565b610185610238366004610c68565b6107e3565b600061026e6040518060400160405280600e81526020016d29ba30b5b4b733a6b0b730b3b2b960911b8152506102d5565b905090565b600061026e6040518060400160405280601981526020017f526f6c6c75705374617465436861696e436f6e7461696e6572000000000000008152506102d5565b600061026e6040518060400160405280600381526020016244414f60e81b8152505b60405163d502db9760e01b81526000908190309063d502db97906102fd908690600401610ce8565b602060405180830381865afa15801561031a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061033e9190610cfb565b90506001600160a01b03811661038b5760405162461bcd60e51b815260206004820152600d60248201526c1b9bc81b985b59481cd85d9959609a1b60448201526064015b60405180910390fd5b92915050565b600061026e6040518060400160405280601081526020016f2937b6363ab829ba30ba32a1b430b4b760811b8152506102d5565b600061026e6040518060400160405280601081526020016f4368616c6c656e6765466163746f727960801b8152506102d5565b6033546001600160a01b031633146104215760405162461bcd60e51b815260040161038290610d18565b61042b6000610877565b565b600061026e6040518060400160405280601081526020016f2937b6363ab824b7383aba21b430b4b760811b8152506102d5565b6033546001600160a01b0316331461048a5760405162461bcd60e51b815260040161038290610d18565b828181146104cc5760405162461bcd60e51b815260206004820152600f60248201526e0d8cadccee8d040dad2e6dac2e8c6d608b1b6044820152606401610382565b60005b818110156105d3573660008787848181106104ec576104ec610d4d565b90506020028101906104fe9190610d63565b91509150600086868581811061051657610516610d4d565b905060200201602081019061052b9190610c68565b9050600061056e84848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506108c992505050565b9050600061057c82846108f9565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c858583866040516105b39493929190610daa565b60405180910390a1505050505080806105cb90610df1565b9150506104cf565b505050505050565b60006105e7600161096e565b905080156105ff576000805461ff0019166101001790555b6106076109f6565b801561064d576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b600061026e6040518060400160405280601981526020017f526f6c6c7570496e707574436861696e436f6e7461696e6572000000000000008152506102d5565b6033546001600160a01b031633146106ba5760405162461bcd60e51b815260040161038290610d18565b60006106c5836108c9565b905060006106d382846108f9565b90507f9416a153a346f93d95f94b064ae3f148b6460473c6e82b3f9fc2521b873fcd6c84828560405161070893929190610e18565b60405180910390a150505050565b600061026e6040518060400160405280600f81526020016e29ba30ba32aa3930b739b4ba34b7b760891b8152506102d5565b600060656000610757846108c9565b81526020810191909152604001600020546001600160a01b031692915050565b600061026e604051806040016040528060138152602001724c3243726f73734c617965725769746e65737360681b8152506102d5565b600061026e604051806040016040528060138152602001724c3143726f73734c617965725769746e65737360681b8152506102d5565b6033546001600160a01b0316331461080d5760405162461bcd60e51b815260040161038290610d18565b6001600160a01b0381166108725760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610382565b61064d815b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6000816040516020016108dc9190610e4b565b604051602081830303815290604052805190602001209050919050565b60006001600160a01b03821661093e5760405162461bcd60e51b815260206004820152600a60248201526932b6b83a3c9030b2323960b11b6044820152606401610382565b5060009182526065602052604090912080546001600160a01b031981166001600160a01b03938416179091551690565b60008054610100900460ff16156109b5578160ff1660011480156109915750303b155b6109ad5760405162461bcd60e51b815260040161038290610e67565b506000919050565b60005460ff8084169116106109dc5760405162461bcd60e51b815260040161038290610e67565b506000805460ff191660ff92909216919091179055600190565b600054610100900460ff16610a1d5760405162461bcd60e51b815260040161038290610eb5565b61042b600054610100900460ff16610a475760405162461bcd60e51b815260040161038290610eb5565b61042b33610877565b634e487b7160e01b600052604160045260246000fd5b600082601f830112610a7757600080fd5b813567ffffffffffffffff80821115610a9257610a92610a50565b604051601f8301601f19908116603f01168101908282118183101715610aba57610aba610a50565b81604052838152866020858801011115610ad357600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060208284031215610b0557600080fd5b813567ffffffffffffffff811115610b1c57600080fd5b610b2884828501610a66565b949350505050565b60008083601f840112610b4257600080fd5b50813567ffffffffffffffff811115610b5a57600080fd5b6020830191508360208260051b8501011115610b7557600080fd5b9250929050565b60008060008060408587031215610b9257600080fd5b843567ffffffffffffffff80821115610baa57600080fd5b610bb688838901610b30565b90965094506020870135915080821115610bcf57600080fd5b50610bdc87828801610b30565b95989497509550505050565b6001600160a01b038116811461064d57600080fd5b60008060408385031215610c1057600080fd5b823567ffffffffffffffff811115610c2757600080fd5b610c3385828601610a66565b9250506020830135610c4481610be8565b809150509250929050565b600060208284031215610c6157600080fd5b5035919050565b600060208284031215610c7a57600080fd5b8135610c8581610be8565b9392505050565b60005b83811015610ca7578181015183820152602001610c8f565b83811115610cb6576000848401525b50505050565b60008151808452610cd4816020860160208601610c8c565b601f01601f19169290920160200192915050565b602081526000610c856020830184610cbc565b600060208284031215610d0d57600080fd5b8151610c8581610be8565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b634e487b7160e01b600052603260045260246000fd5b6000808335601e19843603018112610d7a57600080fd5b83018035915067ffffffffffffffff821115610d9557600080fd5b602001915036819003821315610b7557600080fd5b6060815283606082015283856080830137600060808583018101919091526001600160a01b039384166020830152919092166040830152601f909201601f19160101919050565b600060018201610e1157634e487b7160e01b600052601160045260246000fd5b5060010190565b606081526000610e2b6060830186610cbc565b6001600160a01b0394851660208401529290931660409091015292915050565b60008251610e5d818460208701610c8c565b9190910192915050565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b60608201526080019056fea2646970667358221220ec0c6358a940350c4ef94e1cfeeba46442ad008f0b948371644b89de7542771364736f6c634300080e0033"

var abiAddressManagerStr = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"_name","type":"string"},{"indexed":false,"internalType":"address","name":"_old","type":"address"},{"indexed":false,"internalType":"address","name":"_new","type":"address"}],"name":"AddressSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"challengeFactory","outputs":[{"internalType":"contract IChallengeFactory","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"dao","outputs":[{"internalType":"contract IDAO","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"name":"getAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"getAddrByHash","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"l1CrossLayerWitness","outputs":[{"internalType":"contract IL1CrossLayerWitness","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"l2CrossLayerWitness","outputs":[{"internalType":"contract IL2CrossLayerWitness","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"name":"resolve","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupInputChain","outputs":[{"internalType":"contract IRollupInputChain","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupInputChainContainer","outputs":[{"internalType":"contract IChainStorageContainer","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupStateChain","outputs":[{"internalType":"contract IRollupStateChain","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rollupStateChainContainer","outputs":[{"internalType":"contract IChainStorageContainer","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"},{"internalType":"address","name":"_addr","type":"address"}],"name":"setAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string[]","name":"_names","type":"string[]"},{"internalType":"address[]","name":"_addrs","type":"address[]"}],"name":"setAddressBatch","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"stakingManager","outputs":[{"internalType":"contract IStakingManager","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"stateTransition","outputs":[{"internalType":"contract IStateTransition","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
