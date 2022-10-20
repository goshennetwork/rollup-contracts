package sync_service

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/ontology-layer-2/rollup-contracts/store"
	"github.com/ontology-layer-2/rollup-contracts/store/leveldbstore"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type Client interface {
	GetBlockByNumber(i uint64) (web3.Hash, uint64)
	BlockNumber() uint64
}

func (c *MockClient) GetBlockByNumber(i uint64) (web3.Hash, uint64) {
	return c.hashes[i], 1
}

func (c *MockClient) BlockNumber() uint64 {
	seed := rand.New(rand.NewSource(time.Now().Unix() + int64(c.count)))
	c.count += 1
	if c.count > 2 && seed.Intn(2) == 1 { //50% to rollback
		var h web3.Hash
		seed.Read(h[:])
		c.hashes[c.count-1] = h
	}
	return c.count
}

type MockClient struct {
	hashes [10001]web3.Hash
	count  uint64
}

func NewRandomMockClient() *MockClient {
	seed := rand.New(rand.NewSource(time.Now().Unix()))
	c := &MockClient{}
	for i, h := range c.hashes {
		seed.Read(h[:])
		c.hashes[i] = h
	}
	return c
}

func NewMemStorage() *store.Storage {
	return store.NewStorage(leveldbstore.NewMemLevelDBStore())
}

type rollbackService struct {
	l1client Client
	db       *store.Storage
}

func (self *rollbackService) syncL1Contracts(startHeight, endHeight uint64) error {
	fmt.Println("start: ", startHeight, "end: ", endHeight)
	overlay := self.db.Writer()
	//only write ez key
	overlay.L2Client().StoreTotalCheckedBatchNum(endHeight)
	overlay.L2Client().StoreCheckedBlockNum(endHeight, endHeight+1)
	//
	bhash, btime := self.l1client.GetBlockByNumber(endHeight)
	//now check point
	highestCheckpointInfo := overlay.GetHighestL1CheckPointInfo()
	dirtyk, dirtyv := overlay.Dirty()
	if highestCheckpointInfo == nil { //first, just record as highest check point info
		overlay.SetHighestL1CheckPointInfo(&schema.L1CheckPointInfo{startHeight, endHeight, dirtyk, dirtyv})
	} else {
		pendingCheckpoint := overlay.GetPendingL1CheckPointInfo()
		if pendingCheckpoint == nil {
			//open a new pend
			pendingCheckpoint = &schema.L1CheckPointInfo{startHeight, endHeight + 1, nil, nil}
		} else {
			//check consistence of pending key
			if pendingCheckpoint.EndPoint != startHeight { //wired should never happen
				//not consistence
				panic(1)
			}
		}
		pendingCheckpoint.DirtyKey = append(pendingCheckpoint.DirtyKey, dirtyk...)
		pendingCheckpoint.DirtyValue = append(pendingCheckpoint.DirtyValue, dirtyv...)
		pendingCheckpoint.EndPoint = endHeight + 1
		if pendingCheckpoint.OldEnough() { //reached, just make pending to highest
			overlay.SetHighestL1CheckPointInfo(pendingCheckpoint)
			//remove pending info, next pending start is end +1
			overlay.SetPendingL1CheckPointInfo(&schema.L1CheckPointInfo{endHeight + 1, endHeight + 1, nil, nil})
		} else { //not reach height, just add to pending
			overlay.SetPendingL1CheckPointInfo(pendingCheckpoint)
		}
	}
	overlay.SetLastSyncedL1Timestamp(btime)
	overlay.SetLastSyncedL1Height(endHeight)
	overlay.SetLastSyncedL1Hash(bhash)
	overlay.Commit()
	return nil
}

func (self rollbackService) run() {
	lastHeight := self.db.GetLastSyncedL1Height()
	isSetup := lastHeight == 0
	round := 0
	startHeight := lastHeight + 1
	for {
		if startHeight >= 10000 {
			return
		}

		if self.db.L2Client().GetTotalCheckedBatchNum() != startHeight-1 {
			fmt.Println("self.db.L2Client().GetTotalCheckedBatchNum(): ", self.db.L2Client().GetTotalCheckedBatchNum())
			fmt.Println("startHeight-1: ", startHeight-1)
			panic(2)
		}

		if self.db.L2Client().GetTotalCheckedBlockNum(startHeight-1) != startHeight {
			panic(2)
		}

		l1Height := self.l1client.BlockNumber()
		fmt.Println("l1Height: ", l1Height, "startHeight", startHeight)
		endHeight, err := CalcEndBlock(startHeight, l1Height)
		if err != nil {
			continue
		}
		//be sure setup first 2 round will not roll back.
		if isSetup && round < 2 { //ez first 2 block
			round++
			endHeight = startHeight
		}
		//now check whether reorg first
		lastHash := self.db.GetLastSyncedL1Hash()
		bhash, _ := self.l1client.GetBlockByNumber(startHeight - 1)
		if lastHash != (web3.Hash{}) && bhash != lastHash { //reorg happen, just rollback to former 32 block simply
			writer := self.db.Writer()
			startHeight = RollBack(writer)
			fmt.Println("roll back")
			lastEnd := startHeight - 1
			bhash, btime := self.l1client.GetBlockByNumber(lastEnd)
			writer.SetLastSyncedL1Height(lastEnd)
			writer.SetLastSyncedL1Timestamp(btime)
			writer.SetLastSyncedL1Hash(bhash)
			writer.SetL1DbVersion(writer.GetL1DbVersion() + 1)
			writer.Commit()
			continue
		}

		err = self.syncL1Contracts(startHeight, endHeight)
		if err != nil {
			time.Sleep(15 * time.Second)
			continue
		}
		startHeight = endHeight + 1
	}
}

func TestRollback(t *testing.T) {
	s := &rollbackService{NewRandomMockClient(), NewMemStorage()}
	s.run()
}
