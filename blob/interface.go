package blob

type BlobOracle interface {
	GetBlobsWithCommitmentVersions(versions ...[32]byte) ([]Blob, []KZGCommitment, error)
}

type PersistStore interface {
	Put(key []byte, value []byte) error //Put the key-value pair to store
	Get(key []byte) ([]byte, error)     //Get the value if key in store
}
