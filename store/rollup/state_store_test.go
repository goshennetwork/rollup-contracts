package rollup

import (
	"testing"

	"github.com/laizy/web3"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	store := NewStateMemStore()
	store.StoreInfo(&schema.StateChainInfo{TotalSize: 1})
	info := store.GetInfo()
	assert.Equal(t, uint64(1), info.TotalSize)

	store.StoreLastL1BlockHeight(1)
	lastHeight, err := store.GetLastL1BlockHeight()
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), lastHeight)

	store.putStateBatchInfo(genStatesBatch(0, [][32]byte{web3.Hash{1}, web3.Hash{2}}))
	state, err := store.GetState(0)
	assert.Nil(t, err)
	assert.Equal(t, web3.Hash{1}, state.BlockHash)
}

func genStatesBatch(startIndex uint64, blockHash [][32]byte) *binding.StateBatchAppendedEvent {
	return &binding.StateBatchAppendedEvent{
		StartIndex: startIndex,
		BlockHash:  blockHash,
	}
}
