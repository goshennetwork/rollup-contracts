package blob

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/laizy/web3/utils/codec"
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

func Encode(data []byte) (ret []*types.Blob) {
	reader := codec.NewZeroCopyReader(encodeLenAndAlign(data))
	for reader.Len() > 0 {
		blob := &types.Blob{}
		for i := 0; i < params.FieldElementsPerBlob; i++ {
			val := reader.ReadBytes(31)
			copy(blob[i][:], val)
		}
		ret = append(ret, blob)
	}

	return ret
}

func Decode(blobs []*types.Blob) ([]byte, error) {
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
