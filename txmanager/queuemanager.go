package txmanager

import (
	"sync"

	"github.com/laizy/log"
	"github.com/laizy/web3"
)

type TxStatus byte

const (
	NotFound TxStatus = iota
	Pending
	Finished

	Err
)

type QueueManager[T comparable] struct {
	TxManager *TxManager
	m         map[T]TxStatus
	getKey    func(tx *web3.Transaction) T

	quit chan struct{}
	sync.RWMutex
}

func NewQueueManager[T comparable](txManager *TxManager, getKey func(tx *web3.Transaction) T) *QueueManager[T] {
	return &QueueManager[T]{TxManager: txManager, m: make(map[T]TxStatus), getKey: getKey, quit: make(chan struct{}, 1)}
}

func (self *QueueManager[T]) Start() {
	go self.run()
}

func (self *QueueManager[T]) run() {
	confirmCh, errCh := make(chan *TxConfirmEvent, 1), make(chan *TxErrorEvent, 1)
	sub := self.TxManager.SubscribeConfirmEvent(confirmCh)
	defer sub.Unsubscribe()

	sub2 := self.TxManager.SubscribeErrorEvent(errCh)
	defer sub2.Unsubscribe()
	for {
		select {
		case <-self.quit:
			return
		case e := <-confirmCh:
			if len(e.Info) > 0 {
				log.Info(e.Info)
			}
			self.Confirm(e.Tx)
		case e := <-errCh:
			if e.Err != nil {
				log.Warn(e.Err.Error())
			}
			self.Error(e.Tx)
		}
	}
}

func (self *QueueManager[T]) Close() {
	close(self.quit)
}

func (self *QueueManager[T]) Send(tx *web3.Transaction, reCalcGas bool, info string, priceLimits ...uint64) error {
	self.Lock()
	defer self.Unlock()
	switch self.m[self.getKey(tx)] {
	case Pending:
		log.Info("task is pending in queue")
		return nil
	case Finished:
		log.Info("task is finished")
		return nil
	case Err:
		log.Warn("task failed, retry...")
	}

	self.m[self.getKey(tx)] = Pending
	return self.TxManager.AsyncSendTx(tx, reCalcGas, info, priceLimits...)
}

func (self *QueueManager[T]) Confirm(tx *web3.Transaction) {
	self.Lock()
	defer self.Unlock()
	self.m[self.getKey(tx)] = Finished
}

func (self *QueueManager[T]) Error(tx *web3.Transaction) {
	self.Lock()
	defer self.Unlock()
	self.m[self.getKey(tx)] = Err
}
