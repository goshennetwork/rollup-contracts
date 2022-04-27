package rollup

// format: queueNum(uint64) + queueStart(uint64) + batchNum(uint64) + batch0Time(uint64) +
// batchLeftTimeDiff([]uint32) + batchesData
type RollupChainInput struct {
	QueueNum          uint64
	QueueStart        uint64
	BatchNum          uint64
	Batch0Time        uint64
	BatchLeftTimeDiff []uint32
	BatchesData       []byte
}
