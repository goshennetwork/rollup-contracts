package blob

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

/*
element_0[0..0,uint32(rawLength),byte(version)]

element_1...element_4095: store bytes31(data),so make sure filedElement is less than module, because the kzg use little endien encode, so last byte set to 0.
*/

const BLOB_VERSION = 0

/// one field is reserved for head element
const DataElementNum = params.FieldElementsPerBlob - 1
const MaxDataByte = DataElementNum * 31 /// every data element store 31 byte, the first byte is always zero

func Encode(data []byte) (ret []*types.Blob, err error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	byteNum := 0
	for len(data) > 0 {
		switch len(data) > MaxDataByte {
		case true:
			byteNum = MaxDataByte
		case false: //not overhead, just store length
			byteNum = len(data)
		}

		/// first element is head element for storing global info
		/// write length to head
		blob := &types.Blob{}
		WriteHeadElement(blob, BLOB_VERSION, uint32(byteNum))
		for i := 0; i < (byteNum+30)/31; i += 1 {
			_ = WriteDataElement(blob, i+1, data[31*i:])
		}
		ret = append(ret, blob)
		data = data[byteNum:]
	}
	return ret, nil
}

func WriteHeadElement(blob *types.Blob, version byte, length uint32) {
	blob[0][0] = version
	binary.BigEndian.PutUint32(blob[0][1:5], length)
}

func ReadHeadElement(blob *types.Blob) (version byte, length uint32) {
	headElement := blob[0]
	return headElement[0], binary.BigEndian.Uint32(headElement[1:5])
}

//WriteDataElement write first 31 byte to the data element
func WriteDataElement(blob *types.Blob, index int, data []byte) error {
	if index < 0 || index > DataElementNum+1 {
		return fmt.Errorf("wrong data index: %d", index)
	}
	var b [32]byte
	copy(b[:31], data)
	blob[index] = b
	return nil
}

func ReadDataElement(blob *types.Blob, index int) (ret []byte, next bool, err error) {
	if index <= 0 || index > DataElementNum+1 {
		return ret, next, fmt.Errorf("wrong data index: %d", index)
	}
	v, l := ReadHeadElement(blob)
	if v != BLOB_VERSION {
		return ret, next, fmt.Errorf("wrong version")
	}
	if uint32(index) > (l+30)/31 {
		return ret, next, fmt.Errorf("no data element at provided index")
	}

	next = uint32(index+1) <= (l+30)/31

	source := blob[index]
	if source[31] > 0 { ///last byte should always be zero
		return ret, false, fmt.Errorf("data element first byte not zeroed")
	}
	maxNum := index * 31
	dest := make([]byte, 31)
	copy(dest[:], source[:])

	if l < uint32(maxNum) { // cut off the padding zero
		if (uint32(maxNum) - l) > 31 {
			//should never happen, already checked before
			panic(1)
		}
		dest = dest[:31-(uint32(maxNum)-l)]
	}

	return dest, next, nil
}

func Decode(blobs []*types.Blob) ([]byte, error) {
	var dataNum uint32
	// check version and calc the length
	for _, blob := range blobs {
		version, length := ReadHeadElement(blob)
		if version != BLOB_VERSION {
			return nil, fmt.Errorf("unsupported version")
		}
		dataNum += length
	}

	r := make([]byte, 0, dataNum)
	/// read data element
	for _, blob := range blobs {
		for j := range blob {
			if j == 0 { //ignore head element
				continue
			}
			data, next, err := ReadDataElement(blob, j)
			if err != nil {
				return nil, fmt.Errorf("read data element: %w", err)
			}
			r = append(r, data[:]...)
			if !next { //next element have no data, just break
				break
			}
		}
	}
	return r, nil
}
