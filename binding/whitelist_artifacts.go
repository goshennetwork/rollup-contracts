package binding

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abiWhitelist *abi.ABI

// WhitelistAbi returns the abi of the Whitelist contract
func WhitelistAbi() *abi.ABI {
	return abiWhitelist
}

var binWhitelist []byte

// WhitelistBin returns the bin of the Whitelist contract
func WhitelistBin() []byte {
	return binWhitelist
}

var binRuntimeWhitelist []byte

// WhitelistBinRuntime returns the runtime bin of the Whitelist contract
func WhitelistBinRuntime() []byte {
	return binRuntimeWhitelist
}

func init() {
	var err error
	abiWhitelist, err = abi.NewABI(abiWhitelistStr)
	if err != nil {
		panic(fmt.Errorf("cannot parse Whitelist abi: %v", err))
	}
	if len(binWhitelistStr) != 0 {
		binWhitelist, err = hex.DecodeString(binWhitelistStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse Whitelist bin: %v", err))
		}
	}
	if len(binRuntimeWhitelistStr) != 0 {
		binRuntimeWhitelist, err = hex.DecodeString(binRuntimeWhitelistStr[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse Whitelist bin runtime: %v", err))
		}
	}
}

var binWhitelistStr = "0x608060405234801561001057600080fd5b506108d6806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806372a1d4641161005b57806372a1d464146100f157806392b5d19014610114578063c4d66de814610127578063e9ed9b641461013a57600080fd5b80633537129f1461008257806342b4632e146100b957806355dba28a146100dc575b600080fd5b6100a561009036600461084a565b60016020526000908152604090205460ff1681565b604051901515815260200160405180910390f35b6100a56100c736600461084a565b60026020526000908152604090205460ff1681565b6100ef6100ea36600461086e565b61014d565b005b6100a56100ff36600461084a565b60036020526000908152604090205460ff1681565b6100ef61012236600461086e565b610306565b6100ef61013536600461084a565b6104b2565b6100ef61014836600461086e565b610679565b60048054604080517f4162169f000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff90921692634162169f9282820192602092908290030181865afa1580156101ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101de91906108ac565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610277576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6f6e6c792064616f20616c6c6f7765640000000000000000000000000000000060448201526064015b60405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660008181526001602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f523de897ac29dfd3280af48d1f26b7222ffeb76627593665a907495621512e5e91015b60405180910390a15050565b60048054604080517f4162169f000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff90921692634162169f9282820192602092908290030181865afa158015610373573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061039791906108ac565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461042b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6f6e6c792064616f20616c6c6f77656400000000000000000000000000000000604482015260640161026e565b73ffffffffffffffffffffffffffffffffffffffff821660008181526003602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f32bae78da1582e04b3d20a2d58c706339a9fe9531d524129ab14e0979dc1ca9c91016102fa565b600054610100900460ff16158080156104d25750600054600160ff909116105b806104ec5750303b1580156104ec575060005460ff166001145b610578576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161026e565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156105d657600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8416179055801561067557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498906020016102fa565b5050565b60048054604080517f4162169f000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff90921692634162169f9282820192602092908290030181865afa1580156106e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061070a91906108ac565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461079e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6f6e6c792064616f20616c6c6f77656400000000000000000000000000000000604482015260640161026e565b73ffffffffffffffffffffffffffffffffffffffff821660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f5df38d395edc15b669d646569bd015513395070b5b4deb8a16300abb060d1b5a91016102fa565b73ffffffffffffffffffffffffffffffffffffffff8116811461084757600080fd5b50565b60006020828403121561085c57600080fd5b813561086781610825565b9392505050565b6000806040838503121561088157600080fd5b823561088c81610825565b9150602083013580151581146108a157600080fd5b809150509250929050565b6000602082840312156108be57600080fd5b81516108678161082556fea164736f6c634300080d000a"

var binRuntimeWhitelistStr = "0x608060405234801561001057600080fd5b506004361061007d5760003560e01c806372a1d4641161005b57806372a1d464146100f157806392b5d19014610114578063c4d66de814610127578063e9ed9b641461013a57600080fd5b80633537129f1461008257806342b4632e146100b957806355dba28a146100dc575b600080fd5b6100a561009036600461084a565b60016020526000908152604090205460ff1681565b604051901515815260200160405180910390f35b6100a56100c736600461084a565b60026020526000908152604090205460ff1681565b6100ef6100ea36600461086e565b61014d565b005b6100a56100ff36600461084a565b60036020526000908152604090205460ff1681565b6100ef61012236600461086e565b610306565b6100ef61013536600461084a565b6104b2565b6100ef61014836600461086e565b610679565b60048054604080517f4162169f000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff90921692634162169f9282820192602092908290030181865afa1580156101ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101de91906108ac565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610277576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6f6e6c792064616f20616c6c6f7765640000000000000000000000000000000060448201526064015b60405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660008181526001602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f523de897ac29dfd3280af48d1f26b7222ffeb76627593665a907495621512e5e91015b60405180910390a15050565b60048054604080517f4162169f000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff90921692634162169f9282820192602092908290030181865afa158015610373573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061039791906108ac565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461042b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6f6e6c792064616f20616c6c6f77656400000000000000000000000000000000604482015260640161026e565b73ffffffffffffffffffffffffffffffffffffffff821660008181526003602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f32bae78da1582e04b3d20a2d58c706339a9fe9531d524129ab14e0979dc1ca9c91016102fa565b600054610100900460ff16158080156104d25750600054600160ff909116105b806104ec5750303b1580156104ec575060005460ff166001145b610578576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161026e565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156105d657600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8416179055801561067557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498906020016102fa565b5050565b60048054604080517f4162169f000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff90921692634162169f9282820192602092908290030181865afa1580156106e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061070a91906108ac565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461079e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6f6e6c792064616f20616c6c6f77656400000000000000000000000000000000604482015260640161026e565b73ffffffffffffffffffffffffffffffffffffffff821660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f5df38d395edc15b669d646569bd015513395070b5b4deb8a16300abb060d1b5a91016102fa565b73ffffffffffffffffffffffffffffffffffffffff8116811461084757600080fd5b50565b60006020828403121561085c57600080fd5b813561086781610825565b9392505050565b6000806040838503121561088157600080fd5b823561088c81610825565b9150602083013580151581146108a157600080fd5b809150509250929050565b6000602082840312156108be57600080fd5b81516108678161082556fea164736f6c634300080d000a"

var abiWhitelistStr = `[{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"canChallenge","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"canPropose","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"canSequence","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract IAddressResolver","name":"_resolver","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"challenger","type":"address"},{"internalType":"bool","name":"enabled","type":"bool"}],"name":"setChallenger","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"proposer","type":"address"},{"internalType":"bool","name":"enabled","type":"bool"}],"name":"setProposer","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sequencer","type":"address"},{"internalType":"bool","name":"enabled","type":"bool"}],"name":"setSequencer","outputs":[],"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"challenger","type":"address"},{"indexed":false,"internalType":"bool","name":"enabled","type":"bool"}],"name":"ChallengerUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"proposer","type":"address"},{"indexed":false,"internalType":"bool","name":"enabled","type":"bool"}],"name":"ProposerUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"submitter","type":"address"},{"indexed":false,"internalType":"bool","name":"enabled","type":"bool"}],"name":"SequencerUpdated","type":"event"}]`
