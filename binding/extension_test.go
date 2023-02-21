package binding

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ontology-layer-2/rollup-contracts/blob"
	"github.com/stretchr/testify/assert"
)

func TestBrotliEncode(t *testing.T) {
	txdata := types.LegacyTx{}
	testCase := []*RollupInputBatches{
		{ //only 2 sub batch
			SubBatches: []*SubBatch{
				{
					Txs: []*types.Transaction{types.NewTx(&txdata)},
				},
				{
					Txs: []*types.Transaction{types.NewTx(&txdata)},
				},
			},
			Version: BrotliEncodeType,
		},
		{ //only queue
			QueueNum: 1,
			Version:  BrotliEncodeType,
		},
	}
	for i, tcase := range testCase {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			code := tcase.Encode()
			t.Log(code)
			d := new(RollupInputBatches)
			if err := d.Decode(code); err != nil {
				t.Fatal(err)
			}
			//no batch will not store version, just simple set to zerp
			if len(tcase.SubBatches) == 0 {
				tcase.Version = NormalEncodeType
			}
			wanted, _ := rlp.EncodeToBytes(tcase)
			got, _ := rlp.EncodeToBytes(d)
			if !reflect.DeepEqual(got, wanted) {
				t.Fatal("txs")
			}
		})

	}
}

func TestInputBatchCodec(t *testing.T) {
	txdata := types.LegacyTx{}
	testCase := []*RollupInputBatches{
		{ //only 2 sub batch
			SubBatches: []*SubBatch{
				{
					Txs: []*types.Transaction{types.NewTx(&txdata)},
				},
				{
					Txs: []*types.Transaction{types.NewTx(&txdata)},
				},
			},
		},
		{ //only queue
			QueueNum: 1,
		},
	}
	for i, tcase := range testCase {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			code := tcase.Encode()
			t.Log(code)
			d := new(RollupInputBatches)
			if err := d.Decode(code, nil); err != nil {
				t.Fatal(err)
			}
			wanted, _ := rlp.EncodeToBytes(tcase)
			got, _ := rlp.EncodeToBytes(d)
			if !reflect.DeepEqual(got, wanted) {
				t.Fatal("txs")
			}
		})

	}
}

func TestBlobCodec(t *testing.T) {
	txdata := types.LegacyTx{}
	testCase := []*RollupInputBatches{
		{ //only 2 sub batch
			SubBatches: []*SubBatch{
				{
					Txs: []*types.Transaction{types.NewTx(&txdata)},
				},
				{
					Txs: []*types.Transaction{types.NewTx(&txdata)},
				},
			},
			Version: BlobEnabledMask | BrotliEnabledMask,
		},
		{ //only queue
			QueueNum: 1,
			Version:  BlobEnabledMask,
		},
	}

	mockOracle := blob.NewMockOracle()
	for i, tcase := range testCase {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			code := tcase.Encode()
			d := new(RollupInputBatches)
			d.Version = tcase.Version
			blobs, err := tcase.Blobs()
			assert.NoError(t, err)
			for _, v := range blobs {
				///save commitment to oracle
				c, ok := v.ComputeCommitment()
				assert.True(t, ok)
				if err := mockOracle.VerifyAndRecordBlob(c.ComputeVersionedHash(), c, v); err != nil {
					t.Fatal(err)
				}
			}

			if err := d.Decode(code, mockOracle); err != nil {
				t.Fatal(err)
			}

			wanted, _ := rlp.EncodeToBytes(tcase)
			got, _ := rlp.EncodeToBytes(d)
			if !reflect.DeepEqual(got, wanted) {
				t.Fatal("txs")
			}
		})

	}
}
