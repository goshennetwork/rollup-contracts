package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiRollupStateChain *abi.ABI

// RollupStateChainAbi returns the abi of the RollupStateChain contract
func RollupStateChainAbi() *abi.ABI {
	return abiRollupStateChain
}

var binRollupStateChain []byte

// RollupStateChainBin returns the bin of the RollupStateChain contract
func RollupStateChainBin() []byte {
	return binRollupStateChain
}

var binRuntimeRollupStateChain []byte

// RollupStateChainBinRuntime returns the runtime bin of the RollupStateChain contract
func RollupStateChainBinRuntime() []byte {
	return binRuntimeRollupStateChain
}

func init() {
	var err error
	abiRollupStateChain, err = abi.NewABI(abiRollupStateChainStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse RollupStateChain abi: %v", err))
	}
	if len(binRollupStateChainStr) != 0 {
		binRollupStateChain, err = hex.DecodeString(binRollupStateChainStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse RollupStateChain bin: %v", err))
		}
	}
	if len(binRuntimeRollupStateChainStr) != 0 {
		binRuntimeRollupStateChain, err = hex.DecodeString(binRuntimeRollupStateChainStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse RollupStateChain bin runtime: %v", err))
		}
	}
}

var binRollupStateChainStr = "0x608060405234801561001057600080fd5b506116f3806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063935b0d781161005b578063935b0d78146100d2578063cd6dc687146100f3578063e46020a114610106578063e9c706df1461011d57600080fd5b8063325aeae21461008257806376ef0aaa146100aa57806392927f11146100bf575b600080fd5b6100956100903660046113e4565b610130565b60405190151581526020015b60405180910390f35b6100bd6100b83660046113e4565b610157565b005b6100956100cd3660046113e4565b610549565b6100da61071f565b60405167ffffffffffffffff90911681526020016100a1565b6100bd610101366004611460565b610824565b61010f60015481565b6040519081526020016100a1565b6100bd61012b36600461148c565b6109fa565b600042600154836040015167ffffffffffffffff1661014f9190611573565b111592915050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635dbaf68b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156101c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e8919061158b565b6040517fb363ff8500000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff919091169063b363ff8590602401602060405180830381865afa158015610254573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061027891906115a8565b610308576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f6e6c79207065726d6974746564206279206368616c6c656e676520636f6e7460448201527f726163740000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b61031181610549565b610377576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420737461746520696e666f000000000000000000000000000060448201526064016102ff565b61038081610130565b156103e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f737461746520636f6e6669726d6564000000000000000000000000000000000060448201526064016102ff565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610454573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610478919061158b565b60208201516040517f5682afa900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff9190911690635682afa990602401600060405180830381600087803b1580156104f157600080fd5b505af1158015610505573d6000803e3d6000fd5b50508251602084015160405191935067ffffffffffffffff1691507f911ace459082010270e47b0b03415673320c53da9e7918fc7d0b0c379f80514590600090a350565b600080600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105dd919061158b565b90508073ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa15801561062a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061064e91906115ca565b67ffffffffffffffff16836020015167ffffffffffffffff1610801561071857506106788361125d565b60208401516040517fada8679800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff83169063ada8679890602401602060405180830381865afa1580156106f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061071691906115e7565b145b9392505050565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561078d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107b1919061158b565b73ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061081f91906115ca565b905090565b600054610100900460ff16158080156108445750600054600160ff909116105b8061085e5750303b15801561085e575060005460ff166001145b6108ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016102ff565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561094857600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000073ffffffffffffffffffffffffffffffffffffffff861602179055600182905580156109f557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166393e59dc16040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a67573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a8b919061158b565b6040517f42b4632e00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff91909116906342b4632e90602401602060405180830381865afa158015610af7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b1b91906115a8565b610b81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6f6e6c792070726f706f7365720000000000000000000000000000000000000060448201526064016102ff565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610bef573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c13919061158b565b90508073ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c8491906115ca565b67ffffffffffffffff168267ffffffffffffffff1614610d00576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f737461727420706f73206d69736d61746368000000000000000000000000000060448201526064016102ff565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166322828cc26040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d6d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d91919061158b565b6040517f6f49712b00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9190911690636f49712b90602401602060405180830381865afa158015610dfd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e2191906115a8565b610e87576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600860248201527f756e7374616b656400000000000000000000000000000000000000000000000060448201526064016102ff565b6000835111610ef2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f6e6f20626c6f636b20686173686573000000000000000000000000000000000060448201526064016102ff565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166374aee6c96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f83919061158b565b73ffffffffffffffffffffffffffffffffffffffff1663761a26616040518163ffffffff1660e01b8152600401602060405180830381865afa158015610fcd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ff191906115ca565b67ffffffffffffffff1683518273ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015611048573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061106c91906115ca565b67ffffffffffffffff166110809190611573565b11156110e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f65786365656420696e70757420636861696e206865696768740000000000000060448201526064016102ff565b604080516080810182526000808252602082018190524267ffffffffffffffff81169383019390935233606083015284905b86518110156111f95786818151811061113557611135611600565b602090810291909101810151845267ffffffffffffffff83169084015273ffffffffffffffffffffffffffffffffffffffff8516636483ec256111778561125d565b6040518263ffffffff1660e01b815260040161119591815260200190565b6020604051808303816000875af11580156111b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111d891906115ca565b50816111e38161162f565b92505080806111f190611656565b91505061111a565b508467ffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167ffd1ab91e7c217cde3474f0c085a92f117c977c8a9c04b903d549129f00de539a858960405161124d92919061168e565b60405180910390a3505050505050565b600061126882611276565b805190602001209050919050565b60608160000151826020015183604001518460600151604051602001611305949392919093845260c092831b7fffffffffffffffff00000000000000000000000000000000000000000000000090811660208601529190921b16602883015260601b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016603082015260440190565b6040516020818303038152906040529050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156113915761139161131b565b604052919050565b67ffffffffffffffff811681146113af57600080fd5b50565b80356113bd81611399565b919050565b73ffffffffffffffffffffffffffffffffffffffff811681146113af57600080fd5b6000608082840312156113f657600080fd5b6040516080810181811067ffffffffffffffff821117156114195761141961131b565b60405282358152602083013561142e81611399565b6020820152604083013561144181611399565b60408201526060830135611454816113c2565b60608201529392505050565b6000806040838503121561147357600080fd5b823561147e816113c2565b946020939093013593505050565b6000806040838503121561149f57600080fd5b823567ffffffffffffffff808211156114b757600080fd5b818501915085601f8301126114cb57600080fd5b81356020828211156114df576114df61131b565b8160051b92506114f081840161134a565b828152928401810192818101908985111561150a57600080fd5b948201945b848610156115285785358252948201949082019061150f565b965061153790508782016113b2565b9450505050509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561158657611586611544565b500190565b60006020828403121561159d57600080fd5b8151610718816113c2565b6000602082840312156115ba57600080fd5b8151801515811461071857600080fd5b6000602082840312156115dc57600080fd5b815161071881611399565b6000602082840312156115f957600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600067ffffffffffffffff80831681810361164c5761164c611544565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361168757611687611544565b5060010190565b60006040820167ffffffffffffffff851683526020604081850152818551808452606086019150828701935060005b818110156116d9578451835293830193918301916001016116bd565b509097965050505050505056fea164736f6c634300080d000a"

var binRuntimeRollupStateChainStr = "0x608060405234801561001057600080fd5b506004361061007d5760003560e01c8063935b0d781161005b578063935b0d78146100d2578063cd6dc687146100f3578063e46020a114610106578063e9c706df1461011d57600080fd5b8063325aeae21461008257806376ef0aaa146100aa57806392927f11146100bf575b600080fd5b6100956100903660046113e4565b610130565b60405190151581526020015b60405180910390f35b6100bd6100b83660046113e4565b610157565b005b6100956100cd3660046113e4565b610549565b6100da61071f565b60405167ffffffffffffffff90911681526020016100a1565b6100bd610101366004611460565b610824565b61010f60015481565b6040519081526020016100a1565b6100bd61012b36600461148c565b6109fa565b600042600154836040015167ffffffffffffffff1661014f9190611573565b111592915050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635dbaf68b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156101c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e8919061158b565b6040517fb363ff8500000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff919091169063b363ff8590602401602060405180830381865afa158015610254573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061027891906115a8565b610308576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f6e6c79207065726d6974746564206279206368616c6c656e676520636f6e7460448201527f726163740000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b61031181610549565b610377576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420737461746520696e666f000000000000000000000000000060448201526064016102ff565b61038081610130565b156103e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f737461746520636f6e6669726d6564000000000000000000000000000000000060448201526064016102ff565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610454573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610478919061158b565b60208201516040517f5682afa900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff9190911690635682afa990602401600060405180830381600087803b1580156104f157600080fd5b505af1158015610505573d6000803e3d6000fd5b50508251602084015160405191935067ffffffffffffffff1691507f911ace459082010270e47b0b03415673320c53da9e7918fc7d0b0c379f80514590600090a350565b600080600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105dd919061158b565b90508073ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa15801561062a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061064e91906115ca565b67ffffffffffffffff16836020015167ffffffffffffffff1610801561071857506106788361125d565b60208401516040517fada8679800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff83169063ada8679890602401602060405180830381865afa1580156106f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061071691906115e7565b145b9392505050565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561078d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107b1919061158b565b73ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061081f91906115ca565b905090565b600054610100900460ff16158080156108445750600054600160ff909116105b8061085e5750303b15801561085e575060005460ff166001145b6108ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016102ff565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561094857600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000073ffffffffffffffffffffffffffffffffffffffff861602179055600182905580156109f557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166393e59dc16040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a67573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a8b919061158b565b6040517f42b4632e00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff91909116906342b4632e90602401602060405180830381865afa158015610af7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b1b91906115a8565b610b81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6f6e6c792070726f706f7365720000000000000000000000000000000000000060448201526064016102ff565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663388f2a0a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610bef573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c13919061158b565b90508073ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c8491906115ca565b67ffffffffffffffff168267ffffffffffffffff1614610d00576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f737461727420706f73206d69736d61746368000000000000000000000000000060448201526064016102ff565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166322828cc26040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d6d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d91919061158b565b6040517f6f49712b00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9190911690636f49712b90602401602060405180830381865afa158015610dfd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e2191906115a8565b610e87576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600860248201527f756e7374616b656400000000000000000000000000000000000000000000000060448201526064016102ff565b6000835111610ef2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f6e6f20626c6f636b20686173686573000000000000000000000000000000000060448201526064016102ff565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166374aee6c96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f83919061158b565b73ffffffffffffffffffffffffffffffffffffffff1663761a26616040518163ffffffff1660e01b8152600401602060405180830381865afa158015610fcd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ff191906115ca565b67ffffffffffffffff1683518273ffffffffffffffffffffffffffffffffffffffff166331fe09496040518163ffffffff1660e01b8152600401602060405180830381865afa158015611048573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061106c91906115ca565b67ffffffffffffffff166110809190611573565b11156110e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f65786365656420696e70757420636861696e206865696768740000000000000060448201526064016102ff565b604080516080810182526000808252602082018190524267ffffffffffffffff81169383019390935233606083015284905b86518110156111f95786818151811061113557611135611600565b602090810291909101810151845267ffffffffffffffff83169084015273ffffffffffffffffffffffffffffffffffffffff8516636483ec256111778561125d565b6040518263ffffffff1660e01b815260040161119591815260200190565b6020604051808303816000875af11580156111b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111d891906115ca565b50816111e38161162f565b92505080806111f190611656565b91505061111a565b508467ffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167ffd1ab91e7c217cde3474f0c085a92f117c977c8a9c04b903d549129f00de539a858960405161124d92919061168e565b60405180910390a3505050505050565b600061126882611276565b805190602001209050919050565b60608160000151826020015183604001518460600151604051602001611305949392919093845260c092831b7fffffffffffffffff00000000000000000000000000000000000000000000000090811660208601529190921b16602883015260601b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016603082015260440190565b6040516020818303038152906040529050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156113915761139161131b565b604052919050565b67ffffffffffffffff811681146113af57600080fd5b50565b80356113bd81611399565b919050565b73ffffffffffffffffffffffffffffffffffffffff811681146113af57600080fd5b6000608082840312156113f657600080fd5b6040516080810181811067ffffffffffffffff821117156114195761141961131b565b60405282358152602083013561142e81611399565b6020820152604083013561144181611399565b60408201526060830135611454816113c2565b60608201529392505050565b6000806040838503121561147357600080fd5b823561147e816113c2565b946020939093013593505050565b6000806040838503121561149f57600080fd5b823567ffffffffffffffff808211156114b757600080fd5b818501915085601f8301126114cb57600080fd5b81356020828211156114df576114df61131b565b8160051b92506114f081840161134a565b828152928401810192818101908985111561150a57600080fd5b948201945b848610156115285785358252948201949082019061150f565b965061153790508782016113b2565b9450505050509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561158657611586611544565b500190565b60006020828403121561159d57600080fd5b8151610718816113c2565b6000602082840312156115ba57600080fd5b8151801515811461071857600080fd5b6000602082840312156115dc57600080fd5b815161071881611399565b6000602082840312156115f957600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600067ffffffffffffffff80831681810361164c5761164c611544565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361168757611687611544565b5060010190565b60006040820167ffffffffffffffff851683526020604081850152818551808452606086019150828701935060005b818110156116d9578451835293830193918301916001016116bd565b509097965050505050505056fea164736f6c634300080d000a"

var abiRollupStateChainStr = `[{"inputs":[{"internalType":"bytes32[]","name":"_blockHashes","type":"bytes32[]"},{"internalType":"uint64","name":"_startAt","type":"uint64"}],"name":"appendStateBatch","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"fraudProofWindow","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_addressResolver","type":"address"},{"internalType":"uint256","name":"_fraudProofWindow","type":"uint256"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"uint64","name":"timestamp","type":"uint64"},{"internalType":"address","name":"proposer","type":"address"}],"internalType":"struct Types.StateInfo","name":"_stateInfo","type":"tuple"}],"name":"isStateConfirmed","outputs":[{"internalType":"bool","name":"_confirmed","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"uint64","name":"timestamp","type":"uint64"},{"internalType":"address","name":"proposer","type":"address"}],"internalType":"struct Types.StateInfo","name":"_stateInfo","type":"tuple"}],"name":"rollbackStateBefore","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"totalSubmittedState","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"uint64","name":"timestamp","type":"uint64"},{"internalType":"address","name":"proposer","type":"address"}],"internalType":"struct Types.StateInfo","name":"_stateInfo","type":"tuple"}],"name":"verifyStateInfo","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_proposer","type":"address"},{"indexed":true,"internalType":"uint64","name":"_startIndex","type":"uint64"},{"indexed":false,"internalType":"uint64","name":"_timestamp","type":"uint64"},{"indexed":false,"internalType":"bytes32[]","name":"_blockHash","type":"bytes32[]"}],"name":"StateBatchAppended","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint64","name":"_stateIndex","type":"uint64"},{"indexed":true,"internalType":"bytes32","name":"_blockHash","type":"bytes32"}],"name":"StateRollbacked","type":"event"}]`
