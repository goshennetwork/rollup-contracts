package binding

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"

	"github.com/andybalholm/brotli"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/blob"
)

const (
	NormalEncodeType = iota
	BrotliEncodeType
)

const BrotliEnabledMask = 1
const BlobEnabledMask = 1 << 7

func BrotliEnabled(version uint8) bool {
	return version&BrotliEnabledMask > 0
}

func BlobEnabled(version uint8) bool {
	return version&BlobEnabledMask > 0
}

func (self *RollupInputBatches) BlobEnabled() bool {
	return BlobEnabled(self.Version)
}

func (self *RollupInputBatches) BrotliEnabled() bool {
	return BrotliEnabled(self.Version)
}

type CounterRead struct {
	r       io.Reader
	counter uint64
	rawNum  uint64
}

func NewCounterReader(r io.Reader, limit uint64) io.Reader {
	return &CounterRead{r, 0, limit}
}

func (self *CounterRead) Read(b []byte) (n int, e error) {
	n, err := self.r.Read(b)
	self.counter += uint64(n)
	if self.counter > self.rawNum {
		return n, fmt.Errorf("out of rawNum")
	}
	return n, err
}

// format: batchIndex(uint64)+ queueNum(uint64) + queueStartIndex(uint64) + subBatchNum(uint64) + subBatch0Time(uint64) +
// subBatchLeftTimeDiff([]uint32) + batchesData
// batchesData: version(0) + rlp([][]transaction)
type RollupInputBatches struct {
	//BatchIndex ignored when calc hash, because its useless in l2 system
	BatchIndex uint64
	QueueNum   uint64
	QueueStart uint64
	SubBatches []*SubBatch
	Version    byte
}

type SubBatch struct {
	Timestamp uint64
	Txs       []*types.Transaction
}

func (self *RollupInputBatches) Calldata() []byte {
	//function appendInputBatch() public
	funcSelecter := RollupInputChainAbi().Methods["appendInputBatch"].ID()
	return append(funcSelecter, self.Encode()...)
}

// AppendBatch sends a appendBatch transaction in the solidity contract
func (_a *RollupInputChain) AppendInputBatches(batches *RollupInputBatches) *contract.Txn {
	txn := _a.c.Txn("appendInputBatch")
	txn.Data = batches.Calldata()

	return txn
}

func (self *RollupInputBatches) TxsInfo() (startTimestamp uint64, diffs []uint32, txs [][]*types.Transaction) {
	startTimestamp = self.SubBatches[0].Timestamp
	txs = append(txs, self.SubBatches[0].Txs)
	for i := 1; i < len(self.SubBatches); i++ {
		b := self.SubBatches[i]
		prev := self.SubBatches[i-1].Timestamp
		//equal happens when l1 block timestamp not refresh yet
		utils.EnsureTrue(b.Timestamp >= prev && prev+math.MaxUint32 > b.Timestamp)
		diff := b.Timestamp - prev
		diffs = append(diffs, uint32(diff))
		txs = append(txs, b.Txs)
	}
	return
}

func (self *RollupInputBatches) EncodeWithoutIndex(version byte) []byte {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteUint64BE(self.QueueNum).WriteUint64BE(self.QueueStart)
	batchNum := uint64(len(self.SubBatches))
	if batchNum < 1 {
		return sink.WriteUint64BE(0).Bytes()
	}
	startTime, diffs, txs := self.TxsInfo()
	sink.WriteUint64BE(batchNum).WriteUint64BE(startTime)
	for _, dfff := range diffs {
		sink.WriteUint32BE(dfff)
	}
	sink.WriteByte(version) //write version
	rlpTx, err := rlp.EncodeToBytes(txs)
	if err != nil {
		panic(err)
	}
	code := rlpTx
	if BrotliEnabled(version) { // if brotli enabled, compress
		//now compress rlp code
		buffer := &bytes.Buffer{}
		writer := brotli.NewWriter(buffer)
		_, err := writer.Write(rlpTx)
		utils.Ensure(err)
		utils.Ensure(writer.Flush())
		utils.Ensure(writer.Close())
		code = buffer.Bytes() //now write encoded
	}

	if self.BlobEnabled() { //just need to append blob num and version hash
		blobs := blob.Encode(code)
		//write blob num
		///check blob num, make sure byte will not change num value
		if len(blobs) != int(byte(len(blobs))) {
			panic(1)
		}
		sink.WriteByte(byte(len(blobs)))
		for _, v := range blobs {
			commitment, ok := v.ComputeCommitment()
			utils.EnsureTrue(ok)
			sink.WriteBytes(commitment.ComputeVersionedHash().Bytes())
		}
		return sink.Bytes()
	}

	sink.WriteBytes(code)
	return sink.Bytes()
}

func (self *RollupInputBatches) Blobs() ([]*blob.Blob, error) {
	if !self.BlobEnabled() {
		return nil, errors.New("do not support blob")
	}
	if len(self.SubBatches) == 0 { //no blob
		return nil, nil
	}
	_, _, txs := self.TxsInfo()
	rlpTx, err := rlp.EncodeToBytes(txs)
	if err != nil {
		panic(err)
	}
	code := rlpTx
	if self.BrotliEnabled() { // if brotli enabled, compress
		//now compress rlp code
		buffer := &bytes.Buffer{}
		writer := brotli.NewWriter(buffer)
		_, err := writer.Write(rlpTx)
		utils.Ensure(err)
		utils.Ensure(writer.Flush())
		utils.Ensure(writer.Close())
		code = buffer.Bytes() //now write encoded
	}

	//just need to append blob num and version hash
	blobs := blob.Encode(code)
	//write blob num
	///check blob num, make sure byte will not change num value
	if len(blobs) != int(byte(len(blobs))) {
		panic(1)
	}
	return blobs, nil
}

func (self *RollupInputBatches) Encode() []byte {
	sink := codec.NewZeroCopySink(nil)
	sink.WriteUint64BE(self.BatchIndex)
	dataWithoutIndex := self.EncodeWithoutIndex(self.Version)
	return append(sink.Bytes(), dataWithoutIndex...)
}

// InputBatchHash get input hash, ignore first 8 byte
func (self *RollupInputBatches) InputBatchHash() web3.Hash {
	return crypto.Keccak256Hash(self.EncodeWithoutIndex(self.Version))
}

func (self *RollupInputBatches) InputHash(queueHash web3.Hash) web3.Hash {
	return crypto.Keccak256Hash(self.InputBatchHash().Bytes(), queueHash.Bytes())
}

func safeAdd(x, y uint64) uint64 {
	utils.EnsureTrue(y < math.MaxUint64-x)
	return x + y
}

func (self *RollupInputBatches) DecodeWithoutIndex(b []byte, oracle ...blob.BlobOracle) error {
	reader := codec.NewZeroCopyReader(b)
	self.QueueNum = reader.ReadUint64BE()
	self.QueueStart = reader.ReadUint64BE()
	batchNum := reader.ReadUint64BE()
	if batchNum == 0 {
		//check length
		if reader.Len() != 0 {
			return fmt.Errorf("wrong b length")
		}
		return reader.Error()
	}
	batchTime := reader.ReadUint64BE()
	batchesTime := []uint64{batchTime}
	for i := uint64(0); i < batchNum-1; i++ {
		batchTime = safeAdd(batchTime, uint64(reader.ReadUint32BE()))
		if reader.Error() != nil {
			return reader.Error()
		}
		batchesTime = append(batchesTime, batchTime)
	}

	version := reader.ReadUint8()
	// filter out err first, cause follows may change reader
	if reader.Error() != nil {
		return reader.Error()
	}

	self.Version = version
	if self.BlobEnabled() { //blob append blob num and versionHash
		if len(oracle) == 0 {
			return errors.New("no blob oracle")
		}
		blobNum := reader.ReadUint8()
		versionHashes := make([][32]byte, blobNum)
		for i := uint8(0); i < blobNum; i++ {
			versionHashes[i] = reader.ReadHash()
			if reader.Error() != nil {
				return reader.Error()
			}
		}
		blobs, err := oracle[0].GetBlobsWithCommitmentVersions(versionHashes...)
		if err != nil {
			return fmt.Errorf("get blobs with commitment version: %w", err)
		}
		data, err := blob.Decode(blobs)
		if err != nil {
			return fmt.Errorf("decode blobs: %w", err)
		}
		reader = codec.NewZeroCopyReader(data)
	}

	if self.BrotliEnabled() {
		if reader.Len() == 0 {
			return fmt.Errorf("no brotli code")
		}
		brotliCode := reader.ReadBytes(reader.Len())
		//now limit to 4 mb
		rlpcode, err := ioutil.ReadAll(NewCounterReader(brotli.NewReader(bytes.NewReader(brotliCode)), 4*1024*1024))
		if err != nil {
			return err
		}
		//now transfer rlp code to reader
		reader = codec.NewZeroCopyReader(rlpcode)
	}

	rawBatchesData := reader.ReadBytes(reader.Len())
	if reader.Error() != nil {
		return reader.Error()
	}

	txs := make([][]*types.Transaction, 0)
	err := rlp.DecodeBytes(rawBatchesData, &txs)
	if err != nil {
		return err
	}

	if uint64(len(txs)) != batchNum {
		return fmt.Errorf("inconsistent batch num with tx")
	}
	for i, b := range txs {
		self.SubBatches = append(self.SubBatches, &SubBatch{
			Timestamp: batchesTime[i],
			Txs:       b,
		})
	}
	return nil
}

// decode batch info and check in info correctness
func (self *RollupInputBatches) Decode(b []byte, oracle ...blob.BlobOracle) error {
	reader := codec.NewZeroCopyReader(b[:8])
	self.BatchIndex = reader.ReadUint64BE()
	return self.DecodeWithoutIndex(b[8:], oracle...)
}
