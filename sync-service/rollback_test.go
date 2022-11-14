package sync_service

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store"
	"github.com/ontology-layer-2/rollup-contracts/store/leveldbstore"
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
	if c.count > 2 && seed.Intn(10) == 1 { //10% to rollback
		var h web3.Hash
		seed.Read(h[:])
		c.hashes[c.count-1] = h
	}
	return c.count
}

type MockClient struct {
	hashes [50001]web3.Hash
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
	hash, t := self.l1client.GetBlockByNumber(endHeight) // get block first
	overlay1, overlay2, overlay3 := self.db.Writer(), self.db.Writer(), self.db.Writer()
	overlay1.L2Client().StoreTotalCheckedBatchNum(endHeight)
	overlay1.L2Client().StoreCheckedBlockNum(endHeight, endHeight+1)
	queues := make([]*binding.TransactionEnqueuedEvent, endHeight-startHeight+1)
	msgs := make([]*binding.MessageSentEvent, endHeight-startHeight+1)
	for i, _ := range queues {
		index := startHeight + uint64(i) - 1
		queues[i] = &binding.TransactionEnqueuedEvent{QueueIndex: index}
		msgs[i] = &binding.MessageSentEvent{MessageIndex: index, Raw: &web3.Log{}}
	}
	utils.Ensure(overlay2.InputChain().StoreEnqueuedTransaction(queues...))
	utils.Ensure(overlay3.L1CrossLayerWitness().StoreSentMessage(msgs))
	overlay1.StoreHighestL1CheckPointInfo1(startHeight)
	overlay2.StoreHighestL1CheckPointInfo2(startHeight)
	overlay3.StoreHighestL1CheckPointInfo3(startHeight)
	overlay3.Commit()
	overlay2.Commit()
	overlay1.SetLastSyncedL1Timestamp(t)
	overlay1.SetLastSyncedL1Height(endHeight)
	overlay1.SetLastSyncedL1Hash(hash)
	overlay1.Commit()
	return nil
}
func (self rollbackService) run() {
	lastHeight := self.db.GetLastSyncedL1Height()
	isSetup := lastHeight == 0
	round := 0
	startHeight := lastHeight + 1
	for {
		if startHeight >= 50000 {
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
			startHeight = RollBack(writer, 1)
			RollBack(writer, 2)
			RollBack(writer, 3)
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
