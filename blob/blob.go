package blob

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

/*
element_0[byte(version),uint32(rawLength),0..0]

element_1...element_4095: store bytes31(data),so make sure filedElement is less than module, because the kzg use little endien encode, so last byte set to 0.
*/

const BLOB_VERSION = 0

/// one field is reserved for head element
const DataElementNum = params.FieldElementsPerBlob - 1
const MaxDataByte = DataElementNum * 31 /// every data element store 31 byte, the last byte is always zero

func Encode(data []byte) (ret []*types.Blob, err error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	head := true
	for len(data) > 0 {
		offset := 0
		blob := &types.Blob{}
		if head {
			/// first element is head element for storing global info
			/// write length to head
			WriteHeadElement(blob, BLOB_VERSION, uint32(len(data)))
			offset += 1
			head = false
		}
		byteNum := (params.FieldElementsPerBlob - offset) * 31
		if len(data) < byteNum {
			byteNum = len(data)
		}

		for i := 0; i < (byteNum+30)/31; i += 1 {
			_ = WriteDataElement(blob, i+offset, data[31*i:])
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
	var b [32]byte
	copy(b[:31], data)
	blob[index] = b
	return nil
}

func ReadData(blob *types.Blob, isHead bool) (data []byte) {
	l := params.FieldElementsPerBlob
	index := 0
	if isHead {
		l -= 1
		index += 1
	}
	data = make([]byte, l*31)
	offset := 0
	for ; index < len(blob); index++ {
		copy(data[offset:], blob[index][:31])
		offset += 31
	}
	return data
}
func ReadAll(blobs []*types.Blob) (data []byte, err error) {
	if len(blobs) == 0 {
		return nil, errors.New("no blobs")
	}
	version, l := ReadHeadElement(blobs[0])
	if version != BLOB_VERSION {
		return nil, errors.New("wrong version")
	}
	data = make([]byte, l)
	offset := 0
	for i, v := range blobs {
		d := ReadData(v, i == 0)
		copy(data[offset:], d)
		offset += len(d)
		if uint32(offset) >= l {
			break
		}
	}

	return data, nil
}
