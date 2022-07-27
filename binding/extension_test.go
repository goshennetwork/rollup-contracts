package binding

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestInputBatchCodec(t *testing.T) {
	txData := &types.LegacyTx{}

	testCase := []*RollupInputBatches{
		{ //only 2 sub batch
			SubBatches: []*SubBatch{
				{
					Txs: []*types.Transaction{types.NewTx(txData)},
				},
				{
					Txs: []*types.Transaction{types.NewTx(txData)},
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
			d := new(RollupInputBatches)
			if err := d.Decode(code); err != nil {
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
