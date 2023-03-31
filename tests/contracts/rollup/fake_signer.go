package rollup

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/goshennetwork/rollup-contracts/tests/contracts"
	"github.com/laizy/web3"
)

const (
	GasPrice            = 1_000_000_000
	L1CrossLayerFakeKey = "0x01"
)

var L1CrossLayerFakeSender = web3.HexToAddress("0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf")

func CompleteTxData(target web3.Address, gasPrice, gasLimit uint64, data []byte, nonce uint64) *types.LegacyTx {
	to := common.Address(target)
	return &types.LegacyTx{To: &to, Gas: gasLimit, Data: data, Nonce: nonce, GasPrice: big.NewInt(0).SetUint64(gasPrice)}
}

func Sign(target web3.Address, gasPrice, gasLimit uint64, data []byte, nonce uint64, privateKey string) (*big.Int, *big.Int, *big.Int) {
	signer := types.NewEIP155Signer(new(big.Int).SetUint64(contracts.LocalL1ChainEnv.ChainConfig.L2ChainId))
	txData := CompleteTxData(target, gasPrice, gasLimit, data, nonce)
	signedHash := signer.Hash(types.NewTx(txData))
	//private can't same as 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80, which use by L1CrossLayerWitness
	r, s, v := GetRSV(signedHash, contracts.LocalL1ChainEnv.ChainConfig.L2ChainId, privateKey)
	return new(big.Int).SetBytes(r), new(big.Int).SetBytes(s), new(big.Int).SetBytes(v)
}

func GetRSV(hash [32]byte, chainId uint64, privKey string) ([]byte, []byte, []byte) {
	//k now set to 1
	k := big.NewInt(1)
	prv := new(big.Int).SetBytes(hexutil.MustDecode(privKey))
	//order of curve256
	N, _ := new(big.Int).SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)
	halfOrder := new(big.Int).Rsh(N, 1)
	//find inverse in multy ring, this is used for verify signature, g*(inv*k)=g,so inv * k =1 +N*integer, when k =1, inv simply calc to 1
	inv := new(big.Int).ModInverse(k, N)
	//this point is calc from k=1: privkey.Curve.ScalarBaseMult(new(big.Int).SetInt64(1).Bytes())
	r, _ := new(big.Int).SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240", 10)
	r.Mod(r, N)                         //round : R=r/N
	e := new(big.Int).SetBytes(hash[:]) //singed hash
	s := new(big.Int).Mul(prv, r)
	s.Add(s, e)
	s.Mul(s, inv)
	s.Mod(s, N) //s=( (d*r+e) * inv )%N
	postive := true
	//s is the scale of the curv(just like private key), so when s beyond half order, just use inverse element(just change y to positive)
	if s.Cmp(halfOrder) == 1 {
		postive = false
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		panic("calculated S is zero")
	}

	result := make([]byte, 1, 2*32+1)
	//positive flag
	result[0] = 27 + byte(0)
	if !postive {
		result[0] += 1
	}

	// Not sure this needs rounding but safer to do so.
	curvelen := (256 + 7) / 8

	// Pad R and S to curvelen if needed.
	bytelen := (r.BitLen() + 7) / 8
	if bytelen < curvelen {
		result = append(result,
			make([]byte, curvelen-bytelen)...)
	}
	result = append(result, r.Bytes()...)

	bytelen = (s.BitLen() + 7) / 8
	if bytelen < curvelen {
		result = append(result,
			make([]byte, curvelen-bytelen)...)
	}
	result = append(result, s.Bytes()...)

	term := byte(0)
	if result[0] == 28 {
		term = 1
	}
	result = append(result, term)[1:]
	vv := uint64(result[64]) + 35 + chainId*2
	_r := new(big.Int).SetBytes(result[:32]).Bytes() // used to clean leading zeros
	_s := new(big.Int).SetBytes(result[32:64]).Bytes()
	_v := new(big.Int).SetUint64(vv).Bytes()
	return _r, _s, _v
}
