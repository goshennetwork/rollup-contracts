package rollup

import (
	"encoding/binary"
	"fmt"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type L1WitnessStore struct {
	store schema.KeyValueDB
	mmr   *MMR
}

func NewL1WitnessStore(db schema.KeyValueDB) *L1WitnessStore {
	return &L1WitnessStore{
		store: db,
		mmr:   NewL1MMR(db),
	}
}

func (self *L1WitnessStore) StoreSentMessage(msgs []*binding.MessageSentEvent) error {
	num := self.TotalMessage()
	tree := self.mmr.GetCompactMerkleTree()
	sink := codec.NewZeroCopySink(nil)
	for _, msg := range msgs {
		if msg.MessageIndex > num { //ignore duplicated msg
			return fmt.Errorf("mismatch message, want %d, but %d", num, msg.MessageIndex)
		}
		if msg.MessageIndex == num {
			tree.AppendHash(getMsgHash(sink, msg))
			sink.Reset()
			key := genL1SentMessageKey(msg.MessageIndex)
			self.store.Put(key, codec.SerializeToBytes(&schema.CrossLayerSentMessage{
				BlockNumber:  msg.Raw.BlockNumber,
				MessageIndex: msg.MessageIndex,
				Target:       msg.Target,
				Sender:       msg.Sender,
				MMRRoot:      msg.MmrRoot,
				Message:      msg.Message,
			}))
			num++
		}
	}
	self.mmr.StoreCompactMerkleTree(tree)
	self.StoreTotalMessage(num)
	return nil
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

func (self *L1WitnessStore) StoreTotalMessage(num uint64) {
	var v [8]byte
	binary.BigEndian.PutUint64(v[:], num)
	self.store.Put(schema.L1WitnessSentMessageNumPrefix, v[:])
}

func (self *L1WitnessStore) TotalMessage() uint64 {
	v, err := self.store.Get(schema.L1WitnessSentMessageNumPrefix)
	utils.Ensure(err)
	if len(v) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(v)
}

func genL1SentMessageKey(msgIndex uint64) []byte {
	key := make([]byte, 9)
	key[0] = schema.L1WitnessSentMessageKey
	binary.BigEndian.PutUint64(key[1:], msgIndex)
	return key
}
