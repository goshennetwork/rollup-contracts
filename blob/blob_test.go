package blob

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ontology-layer-2/rollup-contracts/blob/kzg"
	"github.com/ontology-layer-2/rollup-contracts/blob/params"
	"github.com/protolambda/go-kzg/bls"
	"github.com/stretchr/testify/assert"
)

/// one field is reserved for head element
const DataElementNum = params.FieldElementsPerBlob - 1
const MaxDataByte = DataElementNum * 31 /// every data element store 31 byte, the last byte is always zero

func genRandomData(length int) []byte {
	s := rand.NewSource(time.Now().Unix())
	r := make([]byte, length)
	rand.New(s).Read(r)
	return r
}

func TestEncode(t *testing.T) {
	var testCases = [][]byte{
		{1},
		{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		genRandomData(MaxDataByte),
		genRandomData(MaxDataByte + 1),
		genRandomData(rand.Intn(100 * MaxDataByte)),
		genRandomData(1 << 24), //16MB
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			b := Encode(testCase)
			decoded, err := Decode(b)
			assert.NoError(t, err, "decode")
			assert.Equal(t, testCase, decoded)
		})
	}
}

func TestCommit(t *testing.T) {
	var testCases = []string{
		"hello, world",
		"hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world",
		string(genRandomData(MaxDataByte)),
		string(genRandomData(MaxDataByte + 1)),
		string(genRandomData(1 << 24)),
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			b := Encode([]byte(testCase))
			for i := range b {
				blob := b[i]
				commitment, ok := blob.ComputeCommitment()
				assert.True(t, ok, "compute commitment")

				/// verify commitments
				frs := make([]bls.Fr, len(blob))
				for i, elem := range blob {
					assert.True(t, bls.FrFrom32(&frs[i], elem))
				}
				g1, err := bls.FromCompressedG1(commitment[:])
				assert.NoError(t, err, "de compress g1")
				err = kzg.VerifyBlobsLegacy([]*bls.G1Point{g1}, [][]bls.Fr{frs})
				if err != nil {
					t.Fatalf("bad verifyBlobs: %v", err)
				}
				ch := commitment.ComputeVersionedHash()

				h := crypto.Keccak256Hash(commitment[:])
				h[0] = 0x01
				assert.Equal(t, ch[:], h[:])
				//t.Log(ch)
			}
		})
	}
}

func TestProof(t *testing.T) {
	var testCases = []string{
		"hello, world",
		"hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world,hello world",
		string(genRandomData(MaxDataByte)),
		string(genRandomData(MaxDataByte + 1)),
		string(genRandomData(1 << 24)),
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			b := Encode([]byte(testCase))
			for i := range b {
				blob := b[i]
				commitment, ok := blob.ComputeCommitment()
				assert.True(t, ok, "compute commitment")

				/// verify commitments
				polynomial := make([]bls.Fr, len(blob))
				for i, elem := range blob {
					assert.True(t, bls.FrFrom32(&polynomial[i], elem))
				}

				var xFr bls.Fr
				x := uint64(2)
				bls.AsFr(&xFr, x)
				var value bls.Fr
				kzg.EvaluatePolyInEvaluationForm(&value, polynomial[:], &xFr)
				proof, err := kzg.ComputeProof(polynomial, &xFr)
				assert.NoError(t, err)

				// Verify kzg proof
				g1, err := bls.FromCompressedG1(commitment[:])
				assert.NoError(t, err, "de compress g1")
				if kzg.VerifyKzgProof(g1, &xFr, &value, proof) != true {
					t.Fatal("failed proof verification")
				}
			}
		})
	}
}
