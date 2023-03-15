package blob

import (
	"fmt"
)

type MockOracle struct {
	Blobs map[[32]byte]*Blob
}

func NewMockOracle() *MockOracle {
	return &MockOracle{
		Blobs: make(map[[32]byte]*Blob),
	}
}

func (self *MockOracle) VerifyAndRecordBlob(version [32]byte, commitment [48]byte, blob *Blob) error {
	c, ok := blob.ComputeCommitment()
	if !ok {
		return fmt.Errorf("can't generate commitment")
	}
	if c != commitment {
		return fmt.Errorf("inconsistent commitment")
	}
	if c.ComputeVersionedHash() != version {
		return fmt.Errorf("inconsistent version")
	}
	var cb Blob
	cb = *blob
	self.Blobs[version] = &cb
	return nil
}

func (self *MockOracle) GetBlobsWithCommitmentVersions(versions ...[32]byte) ([]*Blob, error) {
	ret := make([]*Blob, len(versions))
	for i, v := range versions {
		b := self.Blobs[v]
		if b == nil {
			return nil, fmt.Errorf("no blob, version hash: %x", v)
		}
		ret[i] = b
	}

	return ret, nil
}
