package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiL2StandardBridge *abi.ABI

// L2StandardBridgeAbi returns the abi of the L2StandardBridge contract
func L2StandardBridgeAbi() *abi.ABI {
	return abiL2StandardBridge
}

var binL2StandardBridge []byte

// L2StandardBridgeBin returns the bin of the L2StandardBridge contract
func L2StandardBridgeBin() []byte {
	return binL2StandardBridge
}

var binRuntimeL2StandardBridge []byte

// L2StandardBridgeBinRuntime returns the runtime bin of the L2StandardBridge contract
func L2StandardBridgeBinRuntime() []byte {
	return binRuntimeL2StandardBridge
}

func init() {
	var err error
	abiL2StandardBridge, err = abi.NewABI(abiL2StandardBridgeStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse L2StandardBridge abi: %v", err))
	}
	if len(binL2StandardBridgeStr) != 0 {
		binL2StandardBridge, err = hex.DecodeString(binL2StandardBridgeStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse L2StandardBridge bin: %v", err))
		}
	}
	if len(binRuntimeL2StandardBridgeStr) != 0 {
		binRuntimeL2StandardBridge, err = hex.DecodeString(binRuntimeL2StandardBridgeStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse L2StandardBridge bin runtime: %v", err))
		}
	}
}

var binL2StandardBridgeStr = "0x608060405234801561001057600080fd5b50611a7c806100206000396000f3fe60806040526004361061009a5760003560e01c806381de0dd511610069578063ab5c7bf11161004e578063ab5c7bf114610238578063dad7ecd91461024b578063de1b85fd1461026b57600080fd5b806381de0dd5146101f8578063920f96cc1461021857600080fd5b80630ac2b63a1461012f57806331f092651461018b57806336c717c1146101ab578063485cc955146101d857600080fd5b3661012a5733321461010d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f4163636f756e74206e6f7420454f41000000000000000000000000000000000060448201526064015b60405180910390fd5b6101283333346040518060200160405280600081525061027e565b005b600080fd5b34801561013b57600080fd5b506000546101629062010000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b34801561019757600080fd5b506101286101a63660046115d0565b6103cb565b3480156101b757600080fd5b506032546101629073ffffffffffffffffffffffffffffffffffffffff1681565b3480156101e457600080fd5b506101286101f336600461162c565b6103df565b34801561020457600080fd5b50610128610213366004611665565b6104d1565b34801561022457600080fd5b50610128610233366004611665565b610865565b6101286102463660046116d8565b61087a565b34801561025757600080fd5b5061012861026636600461171a565b610929565b6101286102793660046117b2565b610e1f565b6000631532ec3460e01b8585858560405160240161029f949392919061187d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526032549091506103429073ffffffffffffffffffffffffffffffffffffffff1682610e61565b8473ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e8787876040516103bc939291906118bc565b60405180910390a45050505050565b6103d9843333868686610ef4565b50505050565b60006103eb6001611133565b9050801561042057600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b610429836112b9565b603280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff841617905580156104cc57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b60325460005473ffffffffffffffffffffffffffffffffffffffff9182169133620100009092041614610560576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6e6f207065726d697373696f6e000000000000000000000000000000000000006044820152606401610104565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633981bc986040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105f291906118fa565b90508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610689576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f77726f6e672063726f7373206c617965722073656e64657200000000000000006044820152606401610104565b844710156106f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f455448206e6f7420656e6f7567680000000000000000000000000000000000006044820152606401610104565b6040805160008082526020820190925273ffffffffffffffffffffffffffffffffffffffff881690879060405161072a919061191e565b60006040518083038185875af1925050503d8060008114610767576040519150601f19603f3d011682016040523d82523d6000602084013e61076c565b606091505b50509050806107d7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f455448207472616e73666572206661696c6564000000000000000000000000006044820152606401610104565b8773ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd898a8a8a8a6040516108539493929190611983565b60405180910390a45050505050505050565b610873853386868686610ef4565b5050505050565b3332146108e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f4163636f756e74206e6f7420454f4100000000000000000000000000000000006044820152606401610104565b61092533333485858080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061027e92505050565b5050565b60325460005473ffffffffffffffffffffffffffffffffffffffff91821691336201000090920416146109b8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6e6f207065726d697373696f6e000000000000000000000000000000000000006044820152606401610104565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633981bc986040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a26573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a4a91906118fa565b90508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610ae1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f77726f6e672063726f7373206c617965722073656e64657200000000000000006044820152606401610104565b610b0b887f1d1d8b630000000000000000000000000000000000000000000000000000000061139d565b8015610bb257508773ffffffffffffffffffffffffffffffffffffffff1663c01e1bd66040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610b5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b8391906118fa565b73ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff16145b15610cc6576040517f40c10f1900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152602482018790528916906340c10f1990604401600060405180830381600087803b158015610c2757600080fd5b505af1158015610c3b573d6000803e3d6000fd5b505050508673ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff167fb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd8989898989604051610cb99493929190611983565b60405180910390a4610e14565b600063a9f9e67560e01b8a8a898b8a8a8a604051602401610ced97969594939291906119b9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152603254909150610d909073ffffffffffffffffffffffffffffffffffffffff1682610e61565b8773ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff168b73ffffffffffffffffffffffffffffffffffffffff167f7ea89a4591614515571c2b51f5ea06494056f261c10ab1ed8c03c7590d87bce08a8a8a8a604051610e0a9493929190611983565b60405180910390a4505b505050505050505050565b6104cc33843485858080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061027e92505050565b6000546040517fbb5ddb0f0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063bb5ddb0f90610ebe9085908590600401611a16565b600060405180830381600087803b158015610ed857600080fd5b505af1158015610eec573d6000803e3d6000fd5b505050505050565b6040517f9dc29fac0000000000000000000000000000000000000000000000000000000081523360048201526024810184905273ffffffffffffffffffffffffffffffffffffffff871690639dc29fac90604401600060405180830381600087803b158015610f6257600080fd5b505af1158015610f76573d6000803e3d6000fd5b5050505060008673ffffffffffffffffffffffffffffffffffffffff1663c01e1bd66040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610fc9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fed91906118fa565b9050600063a9f9e67560e01b8289898989898960405160240161101697969594939291906119b9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526032549091506110b99073ffffffffffffffffffffffffffffffffffffffff1682610e61565b3373ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e898989896040516108539493929190611983565b60008054610100900460ff16156111ea578160ff1660011480156111565750303b155b6111e2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610104565b506000919050565b60005460ff808416911610611281576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610104565b50600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff92909216919091179055600190565b600054610100900460ff16611350576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610104565b6000805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff909216919091179055565b60006113a8836113c2565b80156113b957506113b98383611426565b90505b92915050565b60006113ee827f01ffc9a700000000000000000000000000000000000000000000000000000000611426565b80156113bc575061141f827fffffffff00000000000000000000000000000000000000000000000000000000611426565b1592915050565b604080517fffffffff00000000000000000000000000000000000000000000000000000000831660248083019190915282518083039091018152604490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a7000000000000000000000000000000000000000000000000000000001790529051600091908290819073ffffffffffffffffffffffffffffffffffffffff871690617530906114e090869061191e565b6000604051808303818686fa925050503d806000811461151c576040519150601f19603f3d011682016040523d82523d6000602084013e611521565b606091505b509150915060208151101561153c57600093505050506113bc565b8180156115585750808060200190518101906115589190611a4d565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461158457600080fd5b50565b60008083601f84011261159957600080fd5b50813567ffffffffffffffff8111156115b157600080fd5b6020830191508360208285010111156115c957600080fd5b9250929050565b600080600080606085870312156115e657600080fd5b84356115f181611562565b935060208501359250604085013567ffffffffffffffff81111561161457600080fd5b61162087828801611587565b95989497509550505050565b6000806040838503121561163f57600080fd5b823561164a81611562565b9150602083013561165a81611562565b809150509250929050565b60008060008060006080868803121561167d57600080fd5b853561168881611562565b9450602086013561169881611562565b935060408601359250606086013567ffffffffffffffff8111156116bb57600080fd5b6116c788828901611587565b969995985093965092949392505050565b600080602083850312156116eb57600080fd5b823567ffffffffffffffff81111561170257600080fd5b61170e85828601611587565b90969095509350505050565b600080600080600080600060c0888a03121561173557600080fd5b873561174081611562565b9650602088013561175081611562565b9550604088013561176081611562565b9450606088013561177081611562565b93506080880135925060a088013567ffffffffffffffff81111561179357600080fd5b61179f8a828b01611587565b989b979a50959850939692959293505050565b6000806000604084860312156117c757600080fd5b83356117d281611562565b9250602084013567ffffffffffffffff8111156117ee57600080fd5b6117fa86828701611587565b9497909650939450505050565b60005b8381101561182257818101518382015260200161180a565b838111156103d95750506000910152565b6000815180845261184b816020860160208601611807565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600073ffffffffffffffffffffffffffffffffffffffff8087168352808616602084015250836040830152608060608301526115586080830184611833565b73ffffffffffffffffffffffffffffffffffffffff841681528260208201526060604082015260006118f16060830184611833565b95945050505050565b60006020828403121561190c57600080fd5b815161191781611562565b9392505050565b60008251611930818460208701611807565b9190910192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff8516815283602082015260606040820152600061155860608301848661193a565b600073ffffffffffffffffffffffffffffffffffffffff808a1683528089166020840152808816604084015280871660608401525084608083015260c060a0830152611a0960c08301848661193a565b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201526000611a456040830184611833565b949350505050565b600060208284031215611a5f57600080fd5b8151801515811461191757600080fdfea164736f6c634300080d000a"

var binRuntimeL2StandardBridgeStr = "0x60806040526004361061009a5760003560e01c806381de0dd511610069578063ab5c7bf11161004e578063ab5c7bf114610238578063dad7ecd91461024b578063de1b85fd1461026b57600080fd5b806381de0dd5146101f8578063920f96cc1461021857600080fd5b80630ac2b63a1461012f57806331f092651461018b57806336c717c1146101ab578063485cc955146101d857600080fd5b3661012a5733321461010d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f4163636f756e74206e6f7420454f41000000000000000000000000000000000060448201526064015b60405180910390fd5b6101283333346040518060200160405280600081525061027e565b005b600080fd5b34801561013b57600080fd5b506000546101629062010000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b34801561019757600080fd5b506101286101a63660046115d0565b6103cb565b3480156101b757600080fd5b506032546101629073ffffffffffffffffffffffffffffffffffffffff1681565b3480156101e457600080fd5b506101286101f336600461162c565b6103df565b34801561020457600080fd5b50610128610213366004611665565b6104d1565b34801561022457600080fd5b50610128610233366004611665565b610865565b6101286102463660046116d8565b61087a565b34801561025757600080fd5b5061012861026636600461171a565b610929565b6101286102793660046117b2565b610e1f565b6000631532ec3460e01b8585858560405160240161029f949392919061187d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526032549091506103429073ffffffffffffffffffffffffffffffffffffffff1682610e61565b8473ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e8787876040516103bc939291906118bc565b60405180910390a45050505050565b6103d9843333868686610ef4565b50505050565b60006103eb6001611133565b9050801561042057600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b610429836112b9565b603280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff841617905580156104cc57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b60325460005473ffffffffffffffffffffffffffffffffffffffff9182169133620100009092041614610560576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6e6f207065726d697373696f6e000000000000000000000000000000000000006044820152606401610104565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633981bc986040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105f291906118fa565b90508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610689576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f77726f6e672063726f7373206c617965722073656e64657200000000000000006044820152606401610104565b844710156106f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f455448206e6f7420656e6f7567680000000000000000000000000000000000006044820152606401610104565b6040805160008082526020820190925273ffffffffffffffffffffffffffffffffffffffff881690879060405161072a919061191e565b60006040518083038185875af1925050503d8060008114610767576040519150601f19603f3d011682016040523d82523d6000602084013e61076c565b606091505b50509050806107d7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f455448207472616e73666572206661696c6564000000000000000000000000006044820152606401610104565b8773ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd898a8a8a8a6040516108539493929190611983565b60405180910390a45050505050505050565b610873853386868686610ef4565b5050505050565b3332146108e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f4163636f756e74206e6f7420454f4100000000000000000000000000000000006044820152606401610104565b61092533333485858080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061027e92505050565b5050565b60325460005473ffffffffffffffffffffffffffffffffffffffff91821691336201000090920416146109b8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6e6f207065726d697373696f6e000000000000000000000000000000000000006044820152606401610104565b60008060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633981bc986040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a26573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a4a91906118fa565b90508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610ae1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f77726f6e672063726f7373206c617965722073656e64657200000000000000006044820152606401610104565b610b0b887f1d1d8b630000000000000000000000000000000000000000000000000000000061139d565b8015610bb257508773ffffffffffffffffffffffffffffffffffffffff1663c01e1bd66040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610b5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b8391906118fa565b73ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff16145b15610cc6576040517f40c10f1900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152602482018790528916906340c10f1990604401600060405180830381600087803b158015610c2757600080fd5b505af1158015610c3b573d6000803e3d6000fd5b505050508673ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff167fb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd8989898989604051610cb99493929190611983565b60405180910390a4610e14565b600063a9f9e67560e01b8a8a898b8a8a8a604051602401610ced97969594939291906119b9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152603254909150610d909073ffffffffffffffffffffffffffffffffffffffff1682610e61565b8773ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff168b73ffffffffffffffffffffffffffffffffffffffff167f7ea89a4591614515571c2b51f5ea06494056f261c10ab1ed8c03c7590d87bce08a8a8a8a604051610e0a9493929190611983565b60405180910390a4505b505050505050505050565b6104cc33843485858080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061027e92505050565b6000546040517fbb5ddb0f0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063bb5ddb0f90610ebe9085908590600401611a16565b600060405180830381600087803b158015610ed857600080fd5b505af1158015610eec573d6000803e3d6000fd5b505050505050565b6040517f9dc29fac0000000000000000000000000000000000000000000000000000000081523360048201526024810184905273ffffffffffffffffffffffffffffffffffffffff871690639dc29fac90604401600060405180830381600087803b158015610f6257600080fd5b505af1158015610f76573d6000803e3d6000fd5b5050505060008673ffffffffffffffffffffffffffffffffffffffff1663c01e1bd66040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610fc9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fed91906118fa565b9050600063a9f9e67560e01b8289898989898960405160240161101697969594939291906119b9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526032549091506110b99073ffffffffffffffffffffffffffffffffffffffff1682610e61565b3373ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e898989896040516108539493929190611983565b60008054610100900460ff16156111ea578160ff1660011480156111565750303b155b6111e2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610104565b506000919050565b60005460ff808416911610611281576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610104565b50600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff92909216919091179055600190565b600054610100900460ff16611350576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610104565b6000805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff909216919091179055565b60006113a8836113c2565b80156113b957506113b98383611426565b90505b92915050565b60006113ee827f01ffc9a700000000000000000000000000000000000000000000000000000000611426565b80156113bc575061141f827fffffffff00000000000000000000000000000000000000000000000000000000611426565b1592915050565b604080517fffffffff00000000000000000000000000000000000000000000000000000000831660248083019190915282518083039091018152604490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a7000000000000000000000000000000000000000000000000000000001790529051600091908290819073ffffffffffffffffffffffffffffffffffffffff871690617530906114e090869061191e565b6000604051808303818686fa925050503d806000811461151c576040519150601f19603f3d011682016040523d82523d6000602084013e611521565b606091505b509150915060208151101561153c57600093505050506113bc565b8180156115585750808060200190518101906115589190611a4d565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461158457600080fd5b50565b60008083601f84011261159957600080fd5b50813567ffffffffffffffff8111156115b157600080fd5b6020830191508360208285010111156115c957600080fd5b9250929050565b600080600080606085870312156115e657600080fd5b84356115f181611562565b935060208501359250604085013567ffffffffffffffff81111561161457600080fd5b61162087828801611587565b95989497509550505050565b6000806040838503121561163f57600080fd5b823561164a81611562565b9150602083013561165a81611562565b809150509250929050565b60008060008060006080868803121561167d57600080fd5b853561168881611562565b9450602086013561169881611562565b935060408601359250606086013567ffffffffffffffff8111156116bb57600080fd5b6116c788828901611587565b969995985093965092949392505050565b600080602083850312156116eb57600080fd5b823567ffffffffffffffff81111561170257600080fd5b61170e85828601611587565b90969095509350505050565b600080600080600080600060c0888a03121561173557600080fd5b873561174081611562565b9650602088013561175081611562565b9550604088013561176081611562565b9450606088013561177081611562565b93506080880135925060a088013567ffffffffffffffff81111561179357600080fd5b61179f8a828b01611587565b989b979a50959850939692959293505050565b6000806000604084860312156117c757600080fd5b83356117d281611562565b9250602084013567ffffffffffffffff8111156117ee57600080fd5b6117fa86828701611587565b9497909650939450505050565b60005b8381101561182257818101518382015260200161180a565b838111156103d95750506000910152565b6000815180845261184b816020860160208601611807565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600073ffffffffffffffffffffffffffffffffffffffff8087168352808616602084015250836040830152608060608301526115586080830184611833565b73ffffffffffffffffffffffffffffffffffffffff841681528260208201526060604082015260006118f16060830184611833565b95945050505050565b60006020828403121561190c57600080fd5b815161191781611562565b9392505050565b60008251611930818460208701611807565b9190910192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff8516815283602082015260606040820152600061155860608301848661193a565b600073ffffffffffffffffffffffffffffffffffffffff808a1683528089166020840152808816604084015280871660608401525084608083015260c060a0830152611a0960c08301848661193a565b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201526000611a456040830184611833565b949350505050565b600060208284031215611a5f57600080fd5b8151801515811461191757600080fdfea164736f6c634300080d000a"

var abiL2StandardBridgeStr = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_l1Token","type":"address"},{"indexed":true,"internalType":"address","name":"_l2Token","type":"address"},{"indexed":true,"internalType":"address","name":"_from","type":"address"},{"indexed":false,"internalType":"address","name":"_to","type":"address"},{"indexed":false,"internalType":"uint256","name":"_amount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"_data","type":"bytes"}],"name":"DepositFailed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_l1Token","type":"address"},{"indexed":true,"internalType":"address","name":"_l2Token","type":"address"},{"indexed":true,"internalType":"address","name":"_from","type":"address"},{"indexed":false,"internalType":"address","name":"_to","type":"address"},{"indexed":false,"internalType":"uint256","name":"_amount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"_data","type":"bytes"}],"name":"DepositFinalized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_l1Token","type":"address"},{"indexed":true,"internalType":"address","name":"_l2Token","type":"address"},{"indexed":true,"internalType":"address","name":"_from","type":"address"},{"indexed":false,"internalType":"address","name":"_to","type":"address"},{"indexed":false,"internalType":"uint256","name":"_amount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"_data","type":"bytes"}],"name":"WithdrawalInitiated","type":"event"},{"inputs":[],"name":"crossLayerWitness","outputs":[{"internalType":"contract ICrossLayerWitness","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_l1Token","type":"address"},{"internalType":"address","name":"_l2Token","type":"address"},{"internalType":"address","name":"_from","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"finalizeERC20Deposit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_from","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"finalizeETHDeposit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_l2witness","type":"address"},{"internalType":"address","name":"_l1TokenBridge","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"l1TokenBridge","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_l2Token","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"withdrawETH","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"_to","type":"address"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"withdrawETHTo","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"_l2Token","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"withdrawTo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`
