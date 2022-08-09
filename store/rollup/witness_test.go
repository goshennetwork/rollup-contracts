package rollup

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/merkle"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
	"github.com/stretchr/testify/assert"
)

func newL2WitnessStore(db schema.KeyValueDB) *L2WitnessStore {
	return NewL2WitnessStore(db)
}

func newL1WitnessStore(db schema.KeyValueDB) *L1WitnessStore {
	return NewL1WitnessStore(db)
}

func TestZeroCopyAndAbiEncodePacked(t *testing.T) {
	target := web3.HexToAddress("0xEC9C107cf2D52B4E771301c3d702196D2e163bDC")
	msgSender := web3.HexToAddress("0x9A2900E4b204E31dD58eCc8F276808169D8E4A1b")
	msgIndex := uint64(777777777)
	msg := []byte("asdfafdfasfasdfaddfadjfatjydfagjfgajkdakljfakdlgajkhgasjhgajg")
	sink := codec.NewZeroCopySink(nil)
	sink.WriteAddress(target)
	sink.WriteAddress(msgSender)
	sink.WriteUint64BE(msgIndex)
	sink.WriteBytes(msg)
	t.Logf("%x", sink.Bytes())
	t.Log(crypto.Keccak256Hash(sink.Bytes()).String())
}

func TestL1Witness(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	db := overlaydb.NewOverlayDB(storage.NewFakeDB())

	l1Witness := newL1WitnessStore(db)
	mmrStore := NewL1MMR(db)

	m := 5
	n := 10
	msgs := genRandomSentMessage(n)
	l1Witness.StoreSentMessage(msgs)
	proof, err := mmrStore.GetCompactMerkleTree().InclusionProof(uint64(m), uint64(n))
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range proof {
		t.Log(p.String())
	}
	msgHash := getMsgHash(codec.NewZeroCopySink(nil), msgs[m])
	verifier := merkle.NewMerkleVerifier()
	err = verifier.VerifyLeafHashInclusion(msgHash, uint64(m), proof, mmrStore.GetCompactMerkleTree().Root(),
		mmrStore.GetCompactMerkleTree().TreeSize())
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, false, l1Witness.store, msgs)

	relayedMsg := genRandomRelayedMessage(10)
	l1Witness.StoreRelayedMessage(relayedMsg)
	for _, msg := range relayedMsg {
		got, err := l1Witness.GetRelayedMessage(msg.MessageIndex)
		assert.Nil(t, err)
		want := &schema.MessageRelayedEvent{msg.MessageIndex, msg.MsgHash}
		assert.True(t, reflect.DeepEqual(got, want))
	}

	relayFailedMsg := genRandomRelayFailedMessage(10)
	l1Witness.StoreRelayFailedMessage(relayFailedMsg)
	for _, msg := range relayFailedMsg {
		got, err := l1Witness.GetRelayFailedMessage(msg.MessageIndex)
		assert.Nil(t, err)
		want := &schema.MessageRelayFailedEvent{msg.MessageIndex, msg.MsgHash, msg.MmrSize, msg.MmrRoot}
		assert.True(t, reflect.DeepEqual(got, want))
	}
}

func TestL2Witness(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	db := overlaydb.NewOverlayDB(storage.NewFakeDB())
	l2Witness := newL2WitnessStore(db)
	mmrStore := NewL2MMR(db)
	m := 5
	n := 10
	msgs := genRandomSentMessage(n)
	l2Witness.StoreSentMessage(msgs)

	proof, err := mmrStore.GetCompactMerkleTree().InclusionProof(uint64(m), uint64(n))
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range proof {
		t.Log(p.String())
	}
	msgHash := getMsgHash(codec.NewZeroCopySink(nil), msgs[m])
	t.Log(msgHash.String())
	verifier := merkle.NewMerkleVerifier()
	err = verifier.VerifyLeafHashInclusion(msgHash, uint64(m), proof, mmrStore.GetCompactMerkleTree().Root(),
		mmrStore.GetCompactMerkleTree().TreeSize())
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, true, l2Witness.store, msgs)
}

func genRandomSentMessage(length int) []*binding.MessageSentEvent {
	result := make([]*binding.MessageSentEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.MessageSentEvent{
			MessageIndex: rand.Uint64(),
			Raw: &web3.Log{
				BlockNumber: rand.Uint64(),
			},
		}
		_, _ = rand.Read(evt.Target[:])
		_, _ = rand.Read(evt.Sender[:])
		_, _ = rand.Read(evt.MmrRoot[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		message := make([]byte, 64)
		_, _ = rand.Read(message[:])
		evt.Message = message
		result = append(result, evt)
	}
	return result
}

func genRandomRelayedMessage(length int) []*binding.MessageRelayedEvent {
	result := make([]*binding.MessageRelayedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.MessageRelayedEvent{
			MessageIndex: rand.Uint64(),
			Raw: &web3.Log{
				BlockNumber: rand.Uint64(),
			},
		}
		_, _ = rand.Read(evt.MsgHash[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		result = append(result, evt)
	}
	return result
}

func genRandomRelayFailedMessage(length int) []*binding.MessageRelayFailedEvent {
	result := make([]*binding.MessageRelayFailedEvent, 0)
	for i := 0; i < length; i++ {
		evt := &binding.MessageRelayFailedEvent{
			MessageIndex: rand.Uint64(),
			Raw: &web3.Log{
				BlockNumber: rand.Uint64(),
			},
			MmrSize: rand.Uint64(),
		}
		_, _ = rand.Read(evt.MsgHash[:])
		_, _ = rand.Read(evt.Raw.TransactionHash[:])
		_, _ = rand.Read(evt.MmrRoot[:])
		result = append(result, evt)
	}
	return result
}

func assertEqual(t *testing.T, isL2 bool, store schema.KeyValueDB, msgs []*binding.MessageSentEvent) {
	for _, msg := range msgs {
		key := genL1SentMessageKey(msg.MessageIndex)
		if isL2 {
			key = genL2SentMessageKey(msg.MessageIndex)
		}
		data, err := store.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		source := codec.NewZeroCopySource(data)
		newMsg := &schema.CrossLayerSentMessage{}
		err = newMsg.Deserialization(source)
		if err != nil {
			t.Fatal(err)
		}
		if newMsg.Target != msg.Target || newMsg.Sender != msg.Sender || newMsg.MessageIndex != msg.MessageIndex ||
			newMsg.MMRRoot != newMsg.MMRRoot || !bytes.Equal(newMsg.Message, msg.Message) {
			t.Fatal("failed")
		}
	}
}
