package utils

import (
	"io/ioutil"
	"math/big"
	"os"

	"github.com/laizy/web3/utils"
)

func AtomicWriteFile(filePath string, data string) {
	if FileExisted(filePath) {
		filename := filePath + "~"
		err := ioutil.WriteFile(filename, []byte(data), 0644)
		utils.Ensure(err)
		err = os.Rename(filename, filePath)
		utils.Ensure(err)
	} else {
		err := ioutil.WriteFile(filePath, []byte(data), 0644)
		utils.Ensure(err)
	}
}

// FileExisted checks whether filename exists in filesystem
func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Amount(decimal int, f float64) *big.Int {
	f = f * 1e6
	fe6 := new(big.Int).SetUint64(uint64(f))
	fe6.Mul(fe6, Power(decimal))
	out := fe6.Div(fe6, big.NewInt(1e6))
	return out
}

func Power(decimal int) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
}

func FromRawAmount(decimal int, a *big.Int) float64 {
	//copy a
	aa := new(big.Int).Set(a)
	if aa.Sign() == -1 {
		aa.Mul(aa, big.NewInt(-1))
	}
	//a=a*1e6
	aa.Mul(aa, big.NewInt(1e6))
	//a=a/10^decimal
	aa.Div(aa, Power(decimal))
	fe6 := aa.Uint64()
	//a=a/1e6
	return float64(fe6) / 1e6
}

func ETH(f float64) *big.Int {
	return Amount(18, f)
}

func ToETH(a *big.Int) float64 {
	return FromRawAmount(18, a)
}

func Gwei(f float64) *big.Int {
	return Amount(9, f)
}

func ToGwei(a *big.Int) float64 {
	return FromRawAmount(9, a)
}
