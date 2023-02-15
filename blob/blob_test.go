package blob

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg"
	"github.com/protolambda/go-kzg/bls"
	"github.com/stretchr/testify/assert"
)

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
			b, err := Encode(testCase)
			assert.NoError(t, err, "encode")
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
			b, err := Encode([]byte(testCase))
			assert.NoError(t, err, "encode")
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
				assert.Equal(t, ch, h)
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
			b, err := Encode([]byte(testCase))
			assert.NoError(t, err, "encode")
			for i := range b {
				blob := b[i]
				commitment, ok := blob.ComputeCommitment()
				assert.True(t, ok, "compute commitment")

				/// verify commitments
				polynomial := make([]bls.Fr, len(blob))
				for i, elem := range blob {
					assert.True(t, bls.FrFrom32(&polynomial[i], elem))
				}

				x := uint64(2)
				var xFr bls.Fr
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

// Helper: Compute proof for polynomial
func ComputeProof(poly []bls.Fr, x uint64, crsG1 []bls.G1Point) *bls.G1Point {
	// divisor = [-x, 1]
	divisor := [2]bls.Fr{}
	var tmp bls.Fr
	bls.AsFr(&tmp, x)
	bls.SubModFr(&divisor[0], &bls.ZERO, &tmp)
	bls.CopyFr(&divisor[1], &bls.ONE)
	//for i := 0; i < 2; i++ {
	//	fmt.Printf("div poly %d: %s\n", i, FrStr(&divisor[i]))
	//}
	// quot = poly / divisor
	quotientPolynomial := polyLongDiv(poly, divisor[:])
	//for i := 0; i < len(quotientPolynomial); i++ {
	//	fmt.Printf("quot poly %d: %s\n", i, FrStr(&quotientPolynomial[i]))
	//}

	// evaluate quotient poly at shared secret, in G1
	return bls.LinCombG1(crsG1[:len(quotientPolynomial)], quotientPolynomial)
}

// Helper: Long polynomial division for two polynomials in coefficient form
func polyLongDiv(dividend []bls.Fr, divisor []bls.Fr) []bls.Fr {
	a := make([]bls.Fr, len(dividend))
	for i := 0; i < len(a); i++ {
		bls.CopyFr(&a[i], &dividend[i])
	}
	aPos := len(a) - 1
	bPos := len(divisor) - 1
	diff := aPos - bPos
	out := make([]bls.Fr, diff+1)
	for diff >= 0 {
		quot := &out[diff]
		polyFactorDiv(quot, &a[aPos], &divisor[bPos])
		var tmp, tmp2 bls.Fr
		for i := bPos; i >= 0; i-- {
			// In steps: a[diff + i] -= b[i] * quot
			// tmp =  b[i] * quot
			bls.MulModFr(&tmp, quot, &divisor[i])
			// tmp2 = a[diff + i] - tmp
			bls.SubModFr(&tmp2, &a[diff+i], &tmp)
			// a[diff + i] = tmp2
			bls.CopyFr(&a[diff+i], &tmp2)
		}
		aPos -= 1
		diff -= 1
	}
	return out
}

// Helper: invert the divisor, then multiply
func polyFactorDiv(dst *bls.Fr, a *bls.Fr, b *bls.Fr) {
	// TODO: use divmod instead.
	var tmp bls.Fr
	bls.InvModFr(&tmp, b)
	bls.MulModFr(dst, &tmp, a)
}
