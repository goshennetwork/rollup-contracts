package txmanager

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/event"
	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils"
)

const (
	ETHPriceBump    = 10
	GoshenPriceBump = 1
)

type TxWithContext struct {
	Tx         *web3.Transaction
	PriceLimit uint64
	ReCalcGas  bool
	Info       string
}

type Config struct {
	Period         time.Duration
	ChangeInterval time.Duration
	ListenInterval time.Duration
}

func DefaultCfg() Config {
	return Config{
		Period:         30 * time.Minute,
		ChangeInterval: 25 * time.Second,
		ListenInterval: 10 * time.Second,
	}
}

// TxManager concurrent use
type TxManager struct {
	cfg             Config
	Signer          *contract.Signer
	confirmHeight   uint64
	priceBump       *big.Int
	confirmHashFeed event.Feed
	errorTxFeed     event.Feed

	nonce *uint64
	sync.Mutex
}

type TxConfirmEvent struct {
	Tx   *web3.Transaction
	Info string
}

type TxErrorEvent struct {
	Tx  *web3.Transaction
	Err error
}

func NewTxManager(cfg Config, signer *contract.Signer, confirmHeight uint64, priceBump *big.Int) *TxManager {
	r := &TxManager{cfg, signer, confirmHeight, priceBump, event.Feed{}, event.Feed{}, nil, sync.Mutex{}}
	return r
}

func (t *TxManager) SubscribeConfirmEvent(ch chan<- *TxConfirmEvent) event.Subscription {
	return t.confirmHashFeed.Subscribe(ch)
}

func (t *TxManager) SubscribeErrorEvent(ch chan<- *TxErrorEvent) event.Subscription {
	return t.errorTxFeed.Subscribe(ch)
}

func (t *TxManager) resetNonce() {
	t.Lock()
	defer t.Unlock()
	t.nonce = nil
}

func (t *TxManager) nextNonce() (uint64, error) {
	t.Lock()
	defer t.Unlock()

	if t.nonce == nil {
		nonce, err := t.Signer.Eth().GetNonce(t.Signer.Address(), web3.Latest)
		if err != nil {
			log.Errorf("get nonce: %s", err)
		}
		t.nonce = &nonce
		return *t.nonce, nil
	}
	*t.nonce++
	return *t.nonce, nil
}

// ReConstruct only change price and gaslimit
func (t *TxManager) ReConstruct(txp *TxWithContext) error {
	tx := CopyTx(txp.Tx)

	price, err := t.Signer.Eth().GasPrice()
	if err != nil {
		return fmt.Errorf("get gasPrice: %w", err)
	}
	if price < tx.GasPrice { //new gas price is lower than old, just ignore
		return nil
	}

	factor := new(big.Int).Add(big.NewInt(100), t.priceBump)
	bumpPrice100 := new(big.Int).Add(factor.Mul(factor, new(big.Int).SetUint64(tx.GasPrice)), big.NewInt(99))
	bumpPrice := bumpPrice100.Div(bumpPrice100, big.NewInt(100))

	if bumpPrice.Uint64() >= price { // only use bump price if it is higher current price, otherwise it is useless
		price = bumpPrice.Uint64()
	}

	// if the price is larger than price limit, do not change price
	if txp.PriceLimit != 0 && price > txp.PriceLimit {
		return fmt.Errorf("over priceLimit: priceLimit: %d, got: %d", txp.PriceLimit, price)
	}

	if txp.ReCalcGas {
		gas, err := t.EstimateGas(tx)
		if err != nil { /// inner execute error or maybe just network error.todo: split 2 kind of error
			return err
		}
		tx.Gas = gas
	}
	tx.GasPrice = price
	txp.Tx = t.Signer.SignTx(tx)
	return nil
}

func (t *TxManager) EstimateGas(tx *web3.Transaction) (uint64, error) {
	return t.Signer.Client.Eth().EstimateGas(tx.ToCallMsg())
}

func CopyTx(tx *web3.Transaction) *web3.Transaction {
	return &web3.Transaction{
		From:     tx.From,
		To:       tx.To,
		Input:    tx.Input,
		GasPrice: tx.GasPrice,
		Gas:      tx.Gas,
		Value:    tx.Value,
		Nonce:    tx.Nonce,
	}
}

func (t *TxManager) waitAndChange(ctx context.Context, txp *TxWithContext) {
	ticker := time.NewTicker(t.cfg.ChangeInterval)
	defer ticker.Stop()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, t.cfg.Period)
	defer cancel()

	var (
		materErr error
	)
	defer func() {
		if materErr != nil {
			t.errorTxFeed.Send(&TxErrorEvent{txp.Tx, materErr})
		}
	}()

	confirmCh, txCh := make(chan *web3.Transaction, 1), make(chan *web3.Transaction, 1)
	listen := func(txCh <-chan *web3.Transaction) {
		ticker := time.NewTicker(t.cfg.ListenInterval)
		var txs [2]*web3.Transaction
		i := 0
		for {
			select {
			case tx := <-txCh: //only need to record 2 hash because of async
				txs[i%2] = tx
				i++
			case <-ctxWithTimeout.Done():
				// time out
				log.Errorf("tx: %s timeout", utils.JsonStr(txs))
				return
			case <-ticker.C:
				for _, tx := range txs {
					if tx == nil {
						continue
					}
					receipt, err := t.Signer.Eth().GetTransactionReceipt(tx.Hash())
					if err != nil {
						log.Errorf("get tx receipt: %s", err)
						continue
					}
					headBlock, err := t.Signer.Eth().BlockNumber()
					if err != nil {
						log.Errorf("get blockNumber: %s", err)
						continue
					}
					if receipt == nil {
						continue
					}
					if receipt.Status == 1 && receipt.BlockNumber+t.confirmHeight <= headBlock { // confirmed
						confirmCh <- tx
						return
					}
				}

			}
		}
	}

	go listen(txCh)
	txCh <- txp.Tx

	for {
		select {
		case <-ticker.C:
			// time to reconstruct
			if err := t.ReConstruct(txp); err != nil {
				log.Errorf("reconstruct tx: %s", err)
				continue
			}
			_, err := t.Signer.Eth().SendRawTransaction(txp.Tx.MarshalRLP())
			if err != nil {
				log.Errorf("sendRawTransaction: %s", err)
				continue
			}
			txCh <- txp.Tx
		case tx := <-confirmCh:
			//confirmed
			log.Infof("tx: %s confirmed", tx.Hash())
			// notify subscription
			t.confirmHashFeed.Send(&TxConfirmEvent{tx, txp.Info})
			return
		case <-ctx.Done():
			materErr = errors.New("timeout")
			log.Error(materErr.Error())
			return
		}

	}

}

func (t *TxManager) SyncSendTx(tx *web3.Transaction, reCalcGas bool, info string, priceLimits ...uint64) error {
	nextNonce, err := t.nextNonce()
	if err != nil {
		return err
	}
	tx.Nonce = nextNonce
	tx = t.Signer.SignTx(tx)

	priceLimit := uint64(0)
	if len(priceLimits) > 0 {
		priceLimit = priceLimits[0]
	}
	ctx, _ := context.WithTimeout(context.Background(), t.cfg.Period)

	_, err = t.Signer.Eth().SendRawTransaction(tx.MarshalRLP())
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "nonce too low"):
			//nonce error, so it will never be success, just reset nonce and return
			t.resetNonce()
		}
		return fmt.Errorf("sendRawTransaction: %s", err)
	}
	t.waitAndChange(ctx, &TxWithContext{tx, priceLimit, reCalcGas, info})
	return nil
}

func (t *TxManager) AsyncSendTx(tx *web3.Transaction, reCalcGas bool, info string, priceLimits ...uint64) error {
	nonce, err := t.nextNonce()
	if err != nil {
		return err
	}
	tx.Nonce = nonce
	tx = t.Signer.SignTx(tx)

	priceLimit := uint64(0)
	if len(priceLimits) > 0 {
		priceLimit = priceLimits[0]
	}
	ctx, _ := context.WithTimeout(context.Background(), t.cfg.Period)

	_, err = t.Signer.Eth().SendRawTransaction(tx.MarshalRLP())
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "nonce too low"):
			//nonce error, so it will never be success, just reset nonce and return
			t.resetNonce()
		}
		return fmt.Errorf("sendRawTransaction: %s", err)
	}
	go t.waitAndChange(ctx, &TxWithContext{tx, priceLimit, reCalcGas, info})
	return nil
}
