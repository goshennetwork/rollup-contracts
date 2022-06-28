package rollup

import (
	"encoding/binary"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/merkle"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type L1WitnessStore struct {
	store             schema.KeyValueDB
	compactMerkleTree *merkle.CompactMerkleTree
}

func NewL1WitnessStore(db schema.KeyValueDB, tree *merkle.CompactMerkleTree) *L1WitnessStore {
	return &L1WitnessStore{
		store:             db,
		compactMerkleTree: tree,
	}
}

func (self *L1WitnessStore) StoreSentMessage(msgs []*binding.MessageSentEvent) {
	for _, msg := range msgs {
		self.compactMerkleTree.AppendHash(msg.MsgHash())
		key := genL1SentMessageKey(msg.MessageIndex)
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

func (self *L1WitnessStore) GetL1CompactMerkleTree() (uint64, []web3.Hash, error) {
	v, err := self.store.Get(schema.L1CompactMerkleTreeKey)
	if err != nil {
		return 0, []web3.Hash{}, err
	}
	if len(v) == 0 {
		return 0, []web3.Hash{}, nil
	}
	return schema.DeserializeCompactMerkleTree(v)
}

func (self *L1WitnessStore) GetL1MMRProof(msgIndex uint64, size uint64) ([]web3.Hash, error) {
	if size == 0 {
		size = self.compactMerkleTree.TreeSize()
	}
	return self.compactMerkleTree.InclusionProof(msgIndex, size)
}

func (self *L1WitnessStore) GetSentMessage(msgIndex uint64) (*schema.CrossLayerSentMessage, error) {
	key := genL1SentMessageKey(msgIndex)
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

func genL1SentMessageKey(msgIndex uint64) []byte {
	key := make([]byte, 9)
	key[0] = schema.L1WitnessSentMessageKey
	binary.BigEndian.PutUint64(key[1:], msgIndex)
	return key
}
