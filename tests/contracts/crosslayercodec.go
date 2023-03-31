package contracts

import (
	"github.com/goshennetwork/rollup-contracts/binding"
	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils/codec"
)

func CrossLayerMessageHash(target, sender web3.Address, msgIndex uint64, msg []byte) web3.Hash {
	sink := codec.NewZeroCopySink(nil)
	var padding [24]byte
	sink.WriteAddress(target).WriteAddress(sender).WriteBytes(padding[:]).WriteUint64BE(msgIndex).WriteBytes(msg)
	return crypto.Keccak256Hash(sink.Bytes())
}

func EncodeL1ToL2CallData(target, sender web3.Address, msg []byte, msgIndex uint64, mmrRoot web3.Hash, mmrSize uint64) []byte {
	method := binding.L2CrossLayerWitnessAbi().Methods["relayMessage"]
	calldata := method.MustEncodeIDAndInput(target, sender, msg, msgIndex, mmrRoot, mmrSize)
	return calldata
}

func EncodeL2ToL1CallData(target, sender web3.Address, msg []byte, msgIndex uint64, mmrRoot web3.Hash, mmrSize uint64) []byte {
	method := binding.L1CrossLayerWitnessAbi().Methods["relayMessage"]
	calldata := method.MustEncodeIDAndInput(target, sender, msg, msgIndex, mmrRoot, mmrSize)
	return calldata
}
