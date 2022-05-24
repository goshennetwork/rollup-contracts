package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiRollupInputChain *abi.ABI

// RollupInputChainAbi returns the abi of the RollupInputChain contract
func RollupInputChainAbi() *abi.ABI {
	return abiRollupInputChain
}

var binRollupInputChain []byte

// RollupInputChainBin returns the bin of the RollupInputChain contract
func RollupInputChainBin() []byte {
	return binRollupInputChain
}

var binRuntimeRollupInputChain []byte

// RollupInputChainBinRuntime returns the runtime bin of the RollupInputChain contract
func RollupInputChainBinRuntime() []byte {
	return binRuntimeRollupInputChain
}

func init() {
	var err error
	abiRollupInputChain, err = abi.NewABI(abiRollupInputChainStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse RollupInputChain abi: %v", err))
	}
	if len(binRollupInputChainStr) != 0 {
		binRollupInputChain, err = hex.DecodeString(binRollupInputChainStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse RollupInputChain bin: %v", err))
		}
	}
	if len(binRuntimeRollupInputChainStr) != 0 {
		binRuntimeRollupInputChain, err = hex.DecodeString(binRuntimeRollupInputChainStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse RollupInputChain bin runtime: %v", err))
		}
	}
}

var binRollupInputChainStr = "0x608060405234801561001057600080fd5b5061171a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806384b83dbe1161008c578063876ed5cb11610066578063876ed5cb146101d757806390257565146101e0578063989a8366146101f3578063a85006ca1461020657600080fd5b806384b83dbe14610174578063866328781461018e5780638745e9df146101be57600080fd5b8063282bb0d3116100c8578063282bb0d3146101435780632f8f10f81461015a578063761a26611461016257806378f4b2f21461016a57600080fd5b806303d24d43146100ef5780630f1d1767146100f957806319d8ac611461010c575b600080fd5b6100f7610219565b005b6100f76101073660046112b1565b6109ff565b60005461012690600160901b90046001600160401b031681565b6040516001600160401b0390911681526020015b60405180910390f35b61014c61271081565b60405190815260200161013a565b600254610126565b610126610d8c565b61014c620186a081565b60005461012690600160501b90046001600160401b031681565b6101a161019c366004611386565b610e60565b604080519283526001600160401b0390911660208301520161013a565b600054610126906201000090046001600160401b031681565b61014c61c35081565b61014c6101ee366004611386565b610f0a565b6100f76102013660046113aa565b610ff1565b600354610126906001600160401b031681565b600160009054906101000a90046001600160a01b03166001600160a01b0316634162169f6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561026c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061029091906113f5565b6040516397230e8760e01b81523360048201526001600160a01b0391909116906397230e8790602401602060405180830381865afa1580156102d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102fa9190611412565b61033c5760405162461bcd60e51b815260206004820152600e60248201526d37b7363c9039b2b8bab2b731b2b960911b60448201526064015b60405180910390fd5b600160009054906101000a90046001600160a01b03166001600160a01b03166322828cc26040518163ffffffff1660e01b8152600401602060405180830381865afa15801561038f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b391906113f5565b604051636f49712b60e01b81523360048201526001600160a01b039190911690636f49712b90602401602060405180830381865afa1580156103f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041d9190611412565b6104695760405162461bcd60e51b815260206004820152601b60248201527f53657175656e6365722073686f756c64206265207374616b696e6700000000006044820152606401610333565b60015460408051638669d0ab60e01b815290516000926001600160a01b031691638669d0ab9160048083019260209291908290030181865afa1580156104b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d791906113f5565b60035490915060043560c090811c91600c3590911c906001600160401b031681146105445760405162461bcd60e51b815260206004820152601d60248201527f696e636f72726563742070656e64696e6720717565756520696e6465780000006044820152606401610333565b6000610550838361144a565b6002549091506001600160401b03821611156105ba5760405162461bcd60e51b815260206004820152602360248201527f617474656d707420746f20617070656e6420756e617661696c61626c6520717560448201526265756560e81b6064820152608401610333565b60006105c683856110c9565b6003805467ffffffffffffffff19166001600160401b03851617905590506014803560c01c806106235760405162461bcd60e51b81526020600482015260086024820152670dcde40c4c2e8c6d60c31b6044820152606401610333565b61062e60088361144a565b600054909250823560c01c90600160901b90046001600160401b031681118015610660575042816001600160401b0316105b6106a45760405162461bcd60e51b8152602060048201526015602482015274077726f6e672062617463682074696d657374616d7605c1b6044820152606401610333565b6106af60088461144a565b925060015b826001600160401b0316816001600160401b0316101561070057833560e01c6106dd818461144a565b92506106ea60048661144a565b94505080806106f890611475565b9150506106b4565b506001600160401b03851615610768576000600261071f60018861149c565b6001600160401b031681548110610738576107386114c4565b60009182526020909120600160029092020101546001600160401b0390811691508216811115610766578091505b505b60025442906001600160401b03871610156107b9576002866001600160401b031681548110610799576107996114c4565b60009182526020909120600160029092020101546001600160401b031690505b806001600160401b0316826001600160401b03161061081a5760405162461bcd60e51b815260206004820152601d60248201527f6c6173742062617463682074696d657374616d7020746f6f20686967680000006044820152606401610333565b3661082685602061144a565b6001600160401b0316111561086c5760405162461bcd60e51b815260206004820152600c60248201526b0eee4dedcce40d8cadccee8d60a31b6044820152606401610333565b600061087b36600481846114da565b604051610889929190611504565b604080519182900382206020830152810187905260600160408051808303601f19018152908290528051602090910120636483ec2560e01b82526004820181905291506001600160a01b038b1690636483ec2590602401600060405180830381600087803b1580156108fa57600080fd5b505af115801561090e573d6000803e3d6000fd5b5050505082600060126101000a8154816001600160401b0302191690836001600160401b0316021790555060018a6001600160a01b03166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015610979573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061099d9190611514565b6109a7919061149c565b604080516001600160401b038c8116825260208201859052928316928b169133917fc50b6c6e635d01801616d102e4dba8b6e98a856be9044844453c4884572a9451910160405180910390a450505050505050505050565b600033321415610a5d5733905061c35082511115610a585760405162461bcd60e51b8152602060048201526016602482015275746f6f206c6172676520547820646174612073697a6560501b6044820152606401610333565b610bae565b600160009054906101000a90046001600160a01b03166001600160a01b031663f10073df6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ab0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ad491906113f5565b6001600160a01b0316336001600160a01b031614610b345760405162461bcd60e51b815260206004820152601e60248201527f636f6e74726163742063616e206e6f7420656e7175657565204c3220547800006044820152606401610333565b61271082511115610b925760405162461bcd60e51b815260206004820152602260248201527f746f6f206c617267652063726f7373206c6179657220547820646174612073696044820152617a6560f01b6064820152608401610333565b50600054600160501b90046001600160401b03169150600b609b1b5b6000546001600160401b036201000090910481169084161115610c0b5760405162461bcd60e51b81526020600482015260156024820152741d1bdbc81a1a59da08151e0819d85cc81b1a5b5a5d605a1b6044820152606401610333565b620186a0836001600160401b03161015610c5e5760405162461bcd60e51b81526020600482015260146024820152731d1bdbc81b1bddc8151e0819d85cc81b1a5b5a5d60621b6044820152606401610333565b600081858585604051602001610c77949392919061155d565b60408051808303601f190181528282528051602091820120838301909252818352426001600160401b038181169285019283526002805460018082018355600083905296517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace9183029182015593517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf909401805467ffffffffffffffff1916949092169390931790559054919350916001600160a01b038881169290861691610d40916115b9565b6001600160401b03167f6c5721ff22ed986e360747385d8b4dcf65d3c3c1722b8085174cb91ab0980531888886604051610d7c939291906115d0565b60405180910390a4505050505050565b60015460408051638669d0ab60e01b815290516000926001600160a01b031691638669d0ab9160048083019260209291908290030181865afa158015610dd6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dfa91906113f5565b6001600160a01b03166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e37573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e5b9190611514565b905090565b60025460009081906001600160401b03841610610ebf5760405162461bcd60e51b815260206004820152601960248201527f717565756520696e646578206f766572206361706163697479000000000000006044820152606401610333565b60006002846001600160401b031681548110610edd57610edd6114c4565b60009182526020909120600290910201805460019091015490956001600160401b03909116945092505050565b60015460408051638669d0ab60e01b815290516000926001600160a01b031691638669d0ab9160048083019260209291908290030181865afa158015610f54573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f7891906113f5565b6040516315b50cf360e31b81526001600160401b03841660048201526001600160a01b03919091169063ada8679890602401602060405180830381865afa158015610fc7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610feb919061161b565b92915050565b6000610ffd60016111e6565b90508015611015576000805461ff0019166101001790555b600180546001600160a01b0319166001600160a01b0386161790556000805471ffffffffffffffffffffffffffffffff00001916620100006001600160401b038681169190910267ffffffffffffffff60501b191691909117600160501b9185169190910217905580156110c3576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b6000806110d7836028611634565b6001600160401b031690506000816001600160401b038111156110fc576110fc61129b565b6040519080825280601f01601f191660200182016040528015611126576020820181803683370190505b5090506000805b856001600160401b03168110156111dc5760006002611155836001600160401b038b16611663565b81548110611165576111656114c4565b6000918252602091829020604080518082019091526002929092020180548083526001909101546001600160401b0381168385015286880182815260c09190911b6001600160c01b0319169301839052909250906111c460288661144a565b945050505080806111d49061167b565b91505061112d565b5050209392505050565b60008054610100900460ff161561122d578160ff1660011480156112095750303b155b6112255760405162461bcd60e51b815260040161033390611696565b506000919050565b60005460ff8084169116106112545760405162461bcd60e51b815260040161033390611696565b506000805460ff191660ff92909216919091179055600190565b6001600160a01b038116811461128357600080fd5b50565b6001600160401b038116811461128357600080fd5b634e487b7160e01b600052604160045260246000fd5b6000806000606084860312156112c657600080fd5b83356112d18161126e565b925060208401356112e181611286565b915060408401356001600160401b03808211156112fd57600080fd5b818601915086601f83011261131157600080fd5b8135818111156113235761132361129b565b604051601f8201601f19908116603f0116810190838211818310171561134b5761134b61129b565b8160405282815289602084870101111561136457600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b60006020828403121561139857600080fd5b81356113a381611286565b9392505050565b6000806000606084860312156113bf57600080fd5b83356113ca8161126e565b925060208401356113da81611286565b915060408401356113ea81611286565b809150509250925092565b60006020828403121561140757600080fd5b81516113a38161126e565b60006020828403121561142457600080fd5b815180151581146113a357600080fd5b634e487b7160e01b600052601160045260246000fd5b60006001600160401b0380831681851680830382111561146c5761146c611434565b01949350505050565b60006001600160401b038083168181141561149257611492611434565b6001019392505050565b60006001600160401b03838116908316818110156114bc576114bc611434565b039392505050565b634e487b7160e01b600052603260045260246000fd5b600080858511156114ea57600080fd5b838611156114f757600080fd5b5050820193919092039150565b8183823760009101908152919050565b60006020828403121561152657600080fd5b81516113a381611286565b60005b8381101561154c578181015183820152602001611534565b838111156110c35750506000910152565b60006bffffffffffffffffffffffff19808760601b168352808660601b166014840152506001600160401b0360c01b8460c01b16602883015282516115a9816030850160208701611531565b9190910160300195945050505050565b6000828210156115cb576115cb611434565b500390565b60006001600160401b0380861683526060602084015284518060608501526115ff816080860160208901611531565b9316604083015250601f91909101601f19160160800192915050565b60006020828403121561162d57600080fd5b5051919050565b60006001600160401b038083168185168183048111821515161561165a5761165a611434565b02949350505050565b6000821982111561167657611676611434565b500190565b600060001982141561168f5761168f611434565b5060010190565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b60608201526080019056fea2646970667358221220a6d8e75973453c066e8bac1c8d43e33ce919fd0fed915deba8e5bbc5c7fa70dd64736f6c634300080b0033"

var binRuntimeRollupInputChainStr = "0x608060405234801561001057600080fd5b50600436106100ea5760003560e01c806384b83dbe1161008c578063876ed5cb11610066578063876ed5cb146101d757806390257565146101e0578063989a8366146101f3578063a85006ca1461020657600080fd5b806384b83dbe14610174578063866328781461018e5780638745e9df146101be57600080fd5b8063282bb0d3116100c8578063282bb0d3146101435780632f8f10f81461015a578063761a26611461016257806378f4b2f21461016a57600080fd5b806303d24d43146100ef5780630f1d1767146100f957806319d8ac611461010c575b600080fd5b6100f7610219565b005b6100f76101073660046112b1565b6109ff565b60005461012690600160901b90046001600160401b031681565b6040516001600160401b0390911681526020015b60405180910390f35b61014c61271081565b60405190815260200161013a565b600254610126565b610126610d8c565b61014c620186a081565b60005461012690600160501b90046001600160401b031681565b6101a161019c366004611386565b610e60565b604080519283526001600160401b0390911660208301520161013a565b600054610126906201000090046001600160401b031681565b61014c61c35081565b61014c6101ee366004611386565b610f0a565b6100f76102013660046113aa565b610ff1565b600354610126906001600160401b031681565b600160009054906101000a90046001600160a01b03166001600160a01b0316634162169f6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561026c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061029091906113f5565b6040516397230e8760e01b81523360048201526001600160a01b0391909116906397230e8790602401602060405180830381865afa1580156102d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102fa9190611412565b61033c5760405162461bcd60e51b815260206004820152600e60248201526d37b7363c9039b2b8bab2b731b2b960911b60448201526064015b60405180910390fd5b600160009054906101000a90046001600160a01b03166001600160a01b03166322828cc26040518163ffffffff1660e01b8152600401602060405180830381865afa15801561038f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b391906113f5565b604051636f49712b60e01b81523360048201526001600160a01b039190911690636f49712b90602401602060405180830381865afa1580156103f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041d9190611412565b6104695760405162461bcd60e51b815260206004820152601b60248201527f53657175656e6365722073686f756c64206265207374616b696e6700000000006044820152606401610333565b60015460408051638669d0ab60e01b815290516000926001600160a01b031691638669d0ab9160048083019260209291908290030181865afa1580156104b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d791906113f5565b60035490915060043560c090811c91600c3590911c906001600160401b031681146105445760405162461bcd60e51b815260206004820152601d60248201527f696e636f72726563742070656e64696e6720717565756520696e6465780000006044820152606401610333565b6000610550838361144a565b6002549091506001600160401b03821611156105ba5760405162461bcd60e51b815260206004820152602360248201527f617474656d707420746f20617070656e6420756e617661696c61626c6520717560448201526265756560e81b6064820152608401610333565b60006105c683856110c9565b6003805467ffffffffffffffff19166001600160401b03851617905590506014803560c01c806106235760405162461bcd60e51b81526020600482015260086024820152670dcde40c4c2e8c6d60c31b6044820152606401610333565b61062e60088361144a565b600054909250823560c01c90600160901b90046001600160401b031681118015610660575042816001600160401b0316105b6106a45760405162461bcd60e51b8152602060048201526015602482015274077726f6e672062617463682074696d657374616d7605c1b6044820152606401610333565b6106af60088461144a565b925060015b826001600160401b0316816001600160401b0316101561070057833560e01c6106dd818461144a565b92506106ea60048661144a565b94505080806106f890611475565b9150506106b4565b506001600160401b03851615610768576000600261071f60018861149c565b6001600160401b031681548110610738576107386114c4565b60009182526020909120600160029092020101546001600160401b0390811691508216811115610766578091505b505b60025442906001600160401b03871610156107b9576002866001600160401b031681548110610799576107996114c4565b60009182526020909120600160029092020101546001600160401b031690505b806001600160401b0316826001600160401b03161061081a5760405162461bcd60e51b815260206004820152601d60248201527f6c6173742062617463682074696d657374616d7020746f6f20686967680000006044820152606401610333565b3661082685602061144a565b6001600160401b0316111561086c5760405162461bcd60e51b815260206004820152600c60248201526b0eee4dedcce40d8cadccee8d60a31b6044820152606401610333565b600061087b36600481846114da565b604051610889929190611504565b604080519182900382206020830152810187905260600160408051808303601f19018152908290528051602090910120636483ec2560e01b82526004820181905291506001600160a01b038b1690636483ec2590602401600060405180830381600087803b1580156108fa57600080fd5b505af115801561090e573d6000803e3d6000fd5b5050505082600060126101000a8154816001600160401b0302191690836001600160401b0316021790555060018a6001600160a01b03166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015610979573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061099d9190611514565b6109a7919061149c565b604080516001600160401b038c8116825260208201859052928316928b169133917fc50b6c6e635d01801616d102e4dba8b6e98a856be9044844453c4884572a9451910160405180910390a450505050505050505050565b600033321415610a5d5733905061c35082511115610a585760405162461bcd60e51b8152602060048201526016602482015275746f6f206c6172676520547820646174612073697a6560501b6044820152606401610333565b610bae565b600160009054906101000a90046001600160a01b03166001600160a01b031663f10073df6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ab0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ad491906113f5565b6001600160a01b0316336001600160a01b031614610b345760405162461bcd60e51b815260206004820152601e60248201527f636f6e74726163742063616e206e6f7420656e7175657565204c3220547800006044820152606401610333565b61271082511115610b925760405162461bcd60e51b815260206004820152602260248201527f746f6f206c617267652063726f7373206c6179657220547820646174612073696044820152617a6560f01b6064820152608401610333565b50600054600160501b90046001600160401b03169150600b609b1b5b6000546001600160401b036201000090910481169084161115610c0b5760405162461bcd60e51b81526020600482015260156024820152741d1bdbc81a1a59da08151e0819d85cc81b1a5b5a5d605a1b6044820152606401610333565b620186a0836001600160401b03161015610c5e5760405162461bcd60e51b81526020600482015260146024820152731d1bdbc81b1bddc8151e0819d85cc81b1a5b5a5d60621b6044820152606401610333565b600081858585604051602001610c77949392919061155d565b60408051808303601f190181528282528051602091820120838301909252818352426001600160401b038181169285019283526002805460018082018355600083905296517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace9183029182015593517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf909401805467ffffffffffffffff1916949092169390931790559054919350916001600160a01b038881169290861691610d40916115b9565b6001600160401b03167f6c5721ff22ed986e360747385d8b4dcf65d3c3c1722b8085174cb91ab0980531888886604051610d7c939291906115d0565b60405180910390a4505050505050565b60015460408051638669d0ab60e01b815290516000926001600160a01b031691638669d0ab9160048083019260209291908290030181865afa158015610dd6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dfa91906113f5565b6001600160a01b03166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e37573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e5b9190611514565b905090565b60025460009081906001600160401b03841610610ebf5760405162461bcd60e51b815260206004820152601960248201527f717565756520696e646578206f766572206361706163697479000000000000006044820152606401610333565b60006002846001600160401b031681548110610edd57610edd6114c4565b60009182526020909120600290910201805460019091015490956001600160401b03909116945092505050565b60015460408051638669d0ab60e01b815290516000926001600160a01b031691638669d0ab9160048083019260209291908290030181865afa158015610f54573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f7891906113f5565b6040516315b50cf360e31b81526001600160401b03841660048201526001600160a01b03919091169063ada8679890602401602060405180830381865afa158015610fc7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610feb919061161b565b92915050565b6000610ffd60016111e6565b90508015611015576000805461ff0019166101001790555b600180546001600160a01b0319166001600160a01b0386161790556000805471ffffffffffffffffffffffffffffffff00001916620100006001600160401b038681169190910267ffffffffffffffff60501b191691909117600160501b9185169190910217905580156110c3576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b6000806110d7836028611634565b6001600160401b031690506000816001600160401b038111156110fc576110fc61129b565b6040519080825280601f01601f191660200182016040528015611126576020820181803683370190505b5090506000805b856001600160401b03168110156111dc5760006002611155836001600160401b038b16611663565b81548110611165576111656114c4565b6000918252602091829020604080518082019091526002929092020180548083526001909101546001600160401b0381168385015286880182815260c09190911b6001600160c01b0319169301839052909250906111c460288661144a565b945050505080806111d49061167b565b91505061112d565b5050209392505050565b60008054610100900460ff161561122d578160ff1660011480156112095750303b155b6112255760405162461bcd60e51b815260040161033390611696565b506000919050565b60005460ff8084169116106112545760405162461bcd60e51b815260040161033390611696565b506000805460ff191660ff92909216919091179055600190565b6001600160a01b038116811461128357600080fd5b50565b6001600160401b038116811461128357600080fd5b634e487b7160e01b600052604160045260246000fd5b6000806000606084860312156112c657600080fd5b83356112d18161126e565b925060208401356112e181611286565b915060408401356001600160401b03808211156112fd57600080fd5b818601915086601f83011261131157600080fd5b8135818111156113235761132361129b565b604051601f8201601f19908116603f0116810190838211818310171561134b5761134b61129b565b8160405282815289602084870101111561136457600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b60006020828403121561139857600080fd5b81356113a381611286565b9392505050565b6000806000606084860312156113bf57600080fd5b83356113ca8161126e565b925060208401356113da81611286565b915060408401356113ea81611286565b809150509250925092565b60006020828403121561140757600080fd5b81516113a38161126e565b60006020828403121561142457600080fd5b815180151581146113a357600080fd5b634e487b7160e01b600052601160045260246000fd5b60006001600160401b0380831681851680830382111561146c5761146c611434565b01949350505050565b60006001600160401b038083168181141561149257611492611434565b6001019392505050565b60006001600160401b03838116908316818110156114bc576114bc611434565b039392505050565b634e487b7160e01b600052603260045260246000fd5b600080858511156114ea57600080fd5b838611156114f757600080fd5b5050820193919092039150565b8183823760009101908152919050565b60006020828403121561152657600080fd5b81516113a381611286565b60005b8381101561154c578181015183820152602001611534565b838111156110c35750506000910152565b60006bffffffffffffffffffffffff19808760601b168352808660601b166014840152506001600160401b0360c01b8460c01b16602883015282516115a9816030850160208701611531565b9190910160300195945050505050565b6000828210156115cb576115cb611434565b500390565b60006001600160401b0380861683526060602084015284518060608501526115ff816080860160208901611531565b9316604083015250601f91909101601f19160160800192915050565b60006020828403121561162d57600080fd5b5051919050565b60006001600160401b038083168185168183048111821515161561165a5761165a611434565b02949350505050565b6000821982111561167657611676611434565b500190565b600060001982141561168f5761168f611434565b5060010190565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b60608201526080019056fea2646970667358221220a6d8e75973453c066e8bac1c8d43e33ce919fd0fed915deba8e5bbc5c7fa70dd64736f6c634300080b0033"

var abiRollupInputChainStr = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"proposer","type":"address"},{"indexed":true,"internalType":"uint256","name":"startQueueIndex","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"queueNum","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"chainHeight","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"inputHash","type":"bytes32"}],"name":"TransactionAppended","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint64","name":"queueIndex","type":"uint64"},{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"gaslimit","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"},{"indexed":false,"internalType":"uint64","name":"timestamp","type":"uint64"}],"name":"TransactionEnqueued","type":"event"},{"inputs":[],"name":"MAX_CROSS_LAYER_TX_SIZE","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_ROLLUP_TX_SIZE","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_ROLLUP_TX_GAS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"appendBatch","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"chainHeight","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_target","type":"address"},{"internalType":"uint64","name":"_gasLimit","type":"uint64"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"enqueue","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"_inputIndex","type":"uint64"}],"name":"getInputHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"_queueIndex","type":"uint64"}],"name":"getQueueTxInfo","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"},{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_addressResolver","type":"address"},{"internalType":"uint64","name":"_maxTxGasLimit","type":"uint64"},{"internalType":"uint64","name":"_maxCrossLayerTxGasLimit","type":"uint64"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"lastTimestamp","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxCrossLayerTxGasLimit","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxEnqueueTxGasLimit","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"pendingQueueIndex","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalQueue","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"}]`
