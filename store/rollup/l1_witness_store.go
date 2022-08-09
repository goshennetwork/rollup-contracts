package rollup

import (
	"encoding/binary"

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

func (self *L1WitnessStore) StoreSentMessage(msgs []*binding.MessageSentEvent) {
	tree := self.mmr.GetCompactMerkleTree()
	sink := codec.NewZeroCopySink(nil)
	for _, msg := range msgs {
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
	}
	self.mmr.StoreCompactMerkleTree(tree)
}

func (self *L1WitnessStore) StoreRelayedMessage(msgs []*binding.MessageRelayedEvent) {
	for _, msg := range msgs {
		e := &schema.MessageRelayedEvent{MessageIndex: msg.MessageIndex, MsgHash: msg.MsgHash}
		self.store.Put(genL1RelayedMessageKey(msg.MessageIndex), codec.SerializeToBytes(e))
	}
}

func (self *L1WitnessStore) GetRelayedMessage(msgIndex uint64) (*schema.MessageRelayedEvent, error) {
	v, err := self.store.Get(genL1RelayedMessageKey(msgIndex))
	utils.Ensure(err)
	if len(v) == 0 {
		return nil, schema.ErrNotFound
	}
	e := &schema.MessageRelayedEvent{}
	if err := e.Deserialization(codec.NewZeroCopySource(v)); err != nil {
		return nil, err
	}
	return e, nil
}

func (self *L1WitnessStore) StoreRelayFailedMessage(msgs []*binding.MessageRelayFailedEvent) {
	for _, msg := range msgs {
		e := &schema.MessageRelayFailedEvent{msg.MessageIndex, msg.MsgHash, msg.MmrSize, msg.MmrRoot}
		self.store.Put(genL1RelayFailedMessageKey(msg.MessageIndex), codec.SerializeToBytes(e))
	}
}

func (self *L1WitnessStore) GetRelayFailedMessage(msgIndex uint64) (*schema.MessageRelayFailedEvent, error) {
	v, err := self.store.Get(genL1RelayFailedMessageKey(msgIndex))
	utils.Ensure(err)
	if len(v) == 0 {
		return nil, schema.ErrNotFound
	}
	e := &schema.MessageRelayFailedEvent{}
	if err := e.Deserialization(codec.NewZeroCopySource(v)); err != nil {
		return nil, err
	}
	return e, nil
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

func genL1SentMessageKey(msgIndex uint64) []byte {
	key := make([]byte, 9)
	key[0] = schema.L1WitnessSentMessageKey
	binary.BigEndian.PutUint64(key[1:], msgIndex)
	return key
}

func genL1RelayedMessageKey(msgIndex uint64) []byte {
	key := make([]byte, 9)
	key[0] = schema.L1WitnessRelayedMessageKey
	binary.BigEndian.PutUint64(key[1:], msgIndex)
	return key
}

func genL1RelayFailedMessageKey(msgIndex uint64) []byte {
	key := make([]byte, 9)
	key[0] = schema.L1WitnessRelayFailedMessageKey
	binary.BigEndian.PutUint64(key[1:], msgIndex)
	return key
}
