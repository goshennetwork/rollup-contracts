package rollup

import (
	"encoding/binary"

	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/merkle"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type L2WitnessStore struct {
	store             schema.KeyValueDB
	compactMerkleTree *merkle.CompactMerkleTree
}

func NewL2WitnessStore(db schema.KeyValueDB, tree *merkle.CompactMerkleTree) *L2WitnessStore {
	return &L2WitnessStore{
		store:             db,
		compactMerkleTree: tree,
	}
}

func (self *L2WitnessStore) StoreSentMessage(msgs []*binding.MessageSentEvent) {
	sink := codec.NewZeroCopySink(nil)
	for _, msg := range msgs {
		hash := getMsgHash(sink, msg)
		self.compactMerkleTree.AppendHash(hash)
		//root := self.compactMerkleTree.Root()
		//fmt.Printf("store %s, root %s\n", hash.String(), root.String())
		sink.Reset()

		key := genL2SentMessageKey(msg.MessageIndex)
		self.store.Put(key, codec.SerializeToBytes(&schema.CrossLayerSentMessage{
			BlockNumber:  msg.Raw.BlockNumber,
			MessageIndex: msg.MessageIndex,
			Target:       msg.Target,
			Sender:       msg.Sender,
			MMRRoot:      msg.MmrRoot,
			Message:      msg.Message,
		}))
	}
}

func (self *L2WitnessStore) GetSentMessage(msgIndex uint64) (*schema.CrossLayerSentMessage, error) {
	key := genL2SentMessageKey(msgIndex)
	v, err := self.store.Get(key)
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, schema.ErrNotFound
	}
	source := codec.NewZeroCopySource(v)
	msg := &schema.CrossLayerSentMessage{}
	err = msg.Deserialization(source)
	return msg, err
}

func (self *L2WitnessStore) GetL2MMRProof(msgIndex uint64, size uint64) ([]web3.Hash, error) {
	if size == 0 {
		size = self.compactMerkleTree.TreeSize()
	}
	return self.compactMerkleTree.InclusionProof(msgIndex, size)
}

// crypto.Keccak256Hash(
//			msg.Target.Bytes(),
//			msg.Sender.Bytes(),
//			common.LeftPadBytes(new(big.Int).SetUint64(msg.MessageIndex).Bytes(), 64),
//			msg.Message,
//		)
func getMsgHash(sink *codec.ZeroCopySink, msg *binding.MessageSentEvent) web3.Hash {
	sink.WriteAddress(msg.Target)
	sink.WriteAddress(msg.Sender)
	var padding [24]byte
	sink.WriteBytes(padding[:])
	sink.WriteUint64BE(msg.MessageIndex)
	sink.WriteBytes(msg.Message)
	return crypto.Keccak256Hash(sink.Bytes())
}

func genL2SentMessageKey(msgIndex uint64) []byte {
	key := make([]byte, 9)
	key[0] = schema.L2WitnessSentMessageKey
	binary.BigEndian.PutUint64(key[1:], msgIndex)
	return key
}
