package schema

import (
	"testing"

	"github.com/goshennetwork/rollup-contracts/merkle"
	"github.com/laizy/web3"
	"github.com/stretchr/testify/assert"
)

func TestCodec(t *testing.T) {
	var l1 = &merkle.CompactMerkleTree{}
	b := SerializeCompactMerkleTree(l1)
	i, h, err := DeserializeCompactMerkleTree(b)
	assert.Equal(t, uint64(0), i)
	assert.Nil(t, err)
	assert.Equal(t, []web3.Hash{}, h)
}
