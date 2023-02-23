package blob

import "github.com/ethereum/go-ethereum/core/types"

type BlobOracle interface {
	GetBlobsWithCommitmentVersions(versions ...[32]byte) ([]*types.Blob, error)
}
