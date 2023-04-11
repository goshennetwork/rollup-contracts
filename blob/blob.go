package blob

import (
	"encoding/binary"
	"errors"

	"github.com/goshennetwork/rollup-contracts/blob/kzg"
	"github.com/goshennetwork/rollup-contracts/blob/params"
	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils/codec"
	"github.com/protolambda/go-kzg/bls"
	codec2 "github.com/protolambda/ztyp/codec"
)

const BytesPerBlob = 31 * params.FieldElementsPerBlob

func encodeLenAndAlign(data []byte) []byte {
	lenData := uint32(len(data))
	numBlobs := (lenData + 4 + BytesPerBlob - 1) / BytesPerBlob
	output := make([]byte, numBlobs*BytesPerBlob)
	binary.BigEndian.PutUint32(output, lenData)
	copy(output[4:], data)
	return output
}

func decodeLen(data []byte) ([]byte, error) {
	lenData := binary.BigEndian.Uint32(data)
	if lenData+4 > uint32(len(data)) {
		return nil, errors.New("wrong blob format: data len mismatch")
	}
	return data[4 : 4+lenData], nil
}

func Encode(data []byte) (ret []Blob) {
	reader := codec.NewZeroCopyReader(encodeLenAndAlign(data))
	for reader.Len() > 0 {
		blob := Blob{}
		for i := 0; i < params.FieldElementsPerBlob; i++ {
			val := reader.ReadBytes(31)
			copy(blob[i][:], val)
		}
		ret = append(ret, blob)
	}

	return ret
}

func Decode(blobs []Blob) ([]byte, error) {
	sink := codec.NewZeroCopySink(nil)
	for _, blob := range blobs {
		for _, v := range blob {
			if v[31] != 0 {
				return nil, errors.New("wrong blob format: elem too large")
			}
			sink.WriteBytes(v[:31])
		}
	}
	return decodeLen(sink.Bytes())
}

// Blob data
type Blob [params.FieldElementsPerBlob]BLSFieldElement

func (blob *Blob) Deserialize(dr *codec2.DecodingReader) error {
	if blob == nil {
		return errors.New("cannot decode ssz into nil Blob")
	}
	for i := uint64(0); i < params.FieldElementsPerBlob; i++ {
		// TODO: do we want to check if each field element is within range?
		if _, err := dr.Read(blob[i][:]); err != nil {
			return err
		}
	}
	return nil
}

func (blob *Blob) Serialize(w *codec2.EncodingWriter) error {
	for i := range blob {
		if err := w.Write(blob[i][:]); err != nil {
			return err
		}
	}
	return nil
}

type BLSFieldElement [32]byte

func (blob *Blob) ComputeCommitment() (commitment KZGCommitment, ok bool) {
	frs := make([]bls.Fr, len(blob))
	for i, elem := range blob {
		if !bls.FrFrom32(&frs[i], elem) {
			return KZGCommitment{}, false
		}
	}
	// data is presented in eval form
	commitmentG1 := kzg.BlobToKzg(frs)
	var out KZGCommitment
	copy(out[:], bls.ToCompressedG1(commitmentG1))
	return out, true
}

// Compressed BLS12-381 G1 element
type KZGCommitment [48]byte

func (kzg KZGCommitment) ComputeVersionedHash() web3.Hash {
	h := crypto.Keccak256Hash(kzg[:])
	h[0] = params.BlobCommitmentVersionKZG
	return web3.Hash(h)
}

func (self *BlobWithCommitment) Serialization(sink *codec.ZeroCopySink) {
	for _, v := range self.Blob {
		sink.WriteBytes(v[:])
	}
	sink.WriteBytes(self.Commitment[:])
}

func (self *BlobWithCommitment) DeSerialization(source *codec.ZeroCopySource) error {
	reader := source.Reader()
	for i := 0; i < params.FieldElementsPerBlob; i++ {
		self.Blob[i] = [32]byte(reader.ReadHash())
	}
	copy(self.Commitment[:], reader.ReadBytes(48))
	return reader.Error()
}

// BlobWithCommitment store every blob with commitment to reduce the cost of computing commitment
type BlobWithCommitment struct {
	Blob       Blob
	Commitment KZGCommitment
}
