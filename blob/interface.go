package blob

type BlobOracle interface {
	GetBlobsWithCommitmentVersions(versions ...[32]byte) ([]Blob, []KZGCommitment, error)
}
