package rollup

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/stretchr/testify/assert"
)

func TestCtcStore_Enqueue(t *testing.T) {
	ctc := NewInputMemStore()
	rand.Seed(time.Now().Unix())
	data := uint64(rand.Int63n(100000))
	ctc.StoreEnqueuedTransaction(genQueueElement(0, data, data))
	ele, err := ctc.GetEnqueuedTransaction(0)
	assert.Nil(t, err)
	assert.Equal(t, ele.Timestamp, data)
	ctc.StoreEnqueuedTransaction(genQueueElement(1, data, data))
}

func TestCtcStore_AppendSequencerBatch(t *testing.T) {
	ctc := NewInputMemStore()
	ctc.StoreEnqueuedTransaction(genQueueElement(0, 1, 1))
	ctc.StoreEnqueuedTransaction(genQueueElement(1, 1, 1))
	ctc.StoreSequencerBatches(0, genTransactionBatchInfo(0, 2, 0))
}

func genTransactionBatchInfo(batchIndex, batchSize, prevTotalElements uint64) *binding.InputBatchAppendedEvent {
	return &binding.InputBatchAppendedEvent{
		Index:           batchIndex,
		StartQueueIndex: prevTotalElements,
		QueueNum:        0,
	}
}

func genQueueElement(index, timestamp, blockNumber uint64) *binding.TransactionEnqueuedEvent {
	return &binding.TransactionEnqueuedEvent{
		QueueIndex: index,
		Timestamp:  timestamp,
	}
}
