package blob

import (
	"fmt"

	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
)

type LocalOracle struct {
	Diskdb PersistStore
}

func (self *LocalOracle) GetBlobWithCommitment(versionHash [32]byte) (Blob, KZGCommitment, error) {
	v, err := self.Diskdb.Get(versionHash[:])
	if err != nil {
		return Blob{}, KZGCommitment{}, err
	}

	ret := BlobWithCommitment{}
	utils.Ensure(ret.DeSerialization(codec.NewZeroCopySource(v)))
	return ret.Blob, ret.Commitment, nil

}

func (self *LocalOracle) StoreBlobWithCommitment(versionHash [32]byte, commitment KZGCommitment, blob Blob) {
	utils.Ensure(self.Diskdb.Put(versionHash[:], codec.SerializeToBytes(&BlobWithCommitment{blob, commitment})))
}

func (self *LocalOracle) GetBlobsWithCommitmentVersions(versionHashes ...[32]byte) ([]Blob, []KZGCommitment, error) {
	retBlob := make([]Blob, len(versionHashes))
	retCommitment := make([]KZGCommitment, len(versionHashes))
	for i, v := range versionHashes {
		blob, commitment, err := self.GetBlobWithCommitment(v)
		if err != nil {
			return nil, nil, err
		}
		retBlob[i] = blob
		retCommitment[i] = commitment
	}
	return retBlob, retCommitment, nil
}

type LocalCachedOracle struct {
	diskdb PersistStore
	remote BlobOracle
}

func NewLocalCachedOracle(diskdb PersistStore, remote BlobOracle) *LocalCachedOracle {
	return &LocalCachedOracle{diskdb: diskdb, remote: remote}
}

func (self *LocalCachedOracle) VerifyAndRecordBlob(version [32]byte, commitment [48]byte, blob *Blob) error {
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
	(&LocalOracle{self.diskdb}).StoreBlobWithCommitment(version, commitment, *blob)
	return nil
}

func (self *LocalCachedOracle) GetBlobsWithCommitmentVersions(versionHashes ...[32]byte) ([]Blob, []KZGCommitment, error) {
	retBlob := make([]Blob, len(versionHashes))
	retCommitment := make([]KZGCommitment, len(versionHashes))
	/// first try to load from disk
	for i, versionHash := range versionHashes {
		if blob, commitment, err := (&LocalOracle{self.diskdb}).GetBlobWithCommitment(versionHash); err == nil {
			retBlob[i] = blob
			retCommitment[i] = commitment
			continue
		}
		//not find in local diskdb, try to get from remote
		blobs, commitments, err := self.remote.GetBlobsWithCommitmentVersions(versionHash)
		if err != nil {
			return nil, nil, fmt.Errorf("get version from remote: %w", err)
		}
		if len(blobs) != 1 || len(commitments) != 1 {
			return nil, nil, fmt.Errorf("blob and commiement should all be 1")
		}
		//now try to store in local
		if err := self.VerifyAndRecordBlob(versionHash, commitments[0], &blobs[0]); err != nil {
			//get fake blob with commitment
			return nil, nil, fmt.Errorf("verify failed: %w", err)
		}
		retBlob[i] = blobs[0]
		retCommitment[i] = commitments[0]
	}
	return retBlob, retCommitment, nil
}
