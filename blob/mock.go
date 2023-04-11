package blob

import (
	"fmt"
	"sync"
)

type MockOracle struct {
	blobWithCommitment map[[32]byte]*BlobWithCommitment
	rwlock             sync.RWMutex
}

func NewMockOracle() *MockOracle {
	return &MockOracle{
		blobWithCommitment: make(map[[32]byte]*BlobWithCommitment),
	}
}

func (self *MockOracle) write(k [32]byte, v *BlobWithCommitment) {
	self.rwlock.Lock()
	defer self.rwlock.Unlock()
	self.blobWithCommitment[k] = v
}

func (self *MockOracle) read(k [32]byte) *BlobWithCommitment {
	self.rwlock.RLock()
	defer self.rwlock.RUnlock()
	return self.blobWithCommitment[k]
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
	cb := BlobWithCommitment{
		Blob:       *blob,
		Commitment: commitment,
	}

	self.write(version, &cb)
	return nil
}

func (self *MockOracle) GetBlobsWithCommitmentVersions(versions ...[32]byte) ([]Blob, []KZGCommitment, error) {
	retBlob := make([]Blob, len(versions))
	retCommitment := make([]KZGCommitment, len(versions))
	for i, v := range versions {
		blobWithCommitment := self.read(v)
		if blobWithCommitment == nil {
			return nil, nil, fmt.Errorf("no blob, version hash: %x", v)
		}
		retBlob[i] = blobWithCommitment.Blob
		retCommitment[i] = blobWithCommitment.Commitment

	}

	return retBlob, retCommitment, nil
}
