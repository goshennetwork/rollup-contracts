package utils

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
)

var ERR_NONCE = errors.New("nonce error")

type TxWithPriceLimit struct {
	Tx         *web3.Transaction
	PriceLimit uint64
}

const ETHPriceBump = 10
const GoshenPriceBump = 1

//TxManager concurrent use
type TxManager struct {
	*contract.Signer
	confirmHeight uint64
	priceBump     *big.Int

	nonce *uint64
}

func NewTxManager(signer *contract.Signer, confirmHeight uint64, priceBump *big.Int) *TxManager {
	r := &TxManager{signer, confirmHeight, priceBump, nil}
	return r
}

func (t *TxManager) nextNonce() (uint64, error) {
	if t.nonce == nil {
		nonce, err := t.Eth().GetNonce(t.Signer.Address(), web3.Latest)
		if err != nil {
			log.Errorf("get nonce: %w", err)
		}
		t.nonce = &nonce
		return *t.nonce, nil
	}
	*t.nonce++
	return *t.nonce, nil
}

// ReConstruct only change price and gaslimit
func (t *TxManager) ReConstruct(txp *TxWithPriceLimit) error {
	tx := CopyTx(txp.Tx)
	gas, err := t.EstimateGas(tx)
	if err != nil { /// inner execute error or maybe just network error.todo: split 2 kind of error
		return err
	}

	price, err := t.Eth().GasPrice()
	if err != nil {
		return fmt.Errorf("get gasPrice: %w", err)
	}
	if price < tx.GasPrice { //new gas price is lower than old, just ignore
		return nil
	}

	factor := new(big.Int).Add(big.NewInt(100), t.priceBump)
	bumpPrice100 := factor.Mul(factor, new(big.Int).SetUint64(tx.GasPrice))
	bumpPrice := bumpPrice100.Div(bumpPrice100, big.NewInt(100))

	if bumpPrice.Uint64() >= price { // only use bump price if it is higher current price, otherwise it is useless
		price = bumpPrice.Uint64()
	}

	// if the price is larger than price limit, do not change price
	if txp.PriceLimit != 0 && price > txp.PriceLimit {
		return fmt.Errorf("over priceLimit: priceLimit: %d, got: %d", txp.PriceLimit, price)
	}

	tx.Gas = gas
	tx.GasPrice = price
	txp.Tx = t.SignTx(tx)
	return nil
}

func (t *TxManager) EstimateGas(tx *web3.Transaction) (uint64, error) {
	return t.Client.Eth().EstimateGas(tx.ToCallMsg())
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

func (t *TxManager) sendTx(ctx context.Context, txp *TxWithPriceLimit) {
	timer := time.NewTimer(0)
	defer timer.Stop()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 6*time.Hour)
	defer cancel()

	confirmCh, hashCh := make(chan web3.Hash, 1), make(chan web3.Hash, 1)
	listen := func(hashCh <-chan web3.Hash) {
		ticker := time.NewTicker(20 * time.Second)
		var hashes [2]web3.Hash
		i := 0
		for {
			select {
			case hash := <-hashCh: //only need to record 2 hash because of async
				hashes[i%2] = hash
				i++
			case <-ctxWithTimeout.Done():
				// time out
				log.Errorf("tx: %s timeout", hashes)
				return
			case <-ticker.C:
				for _, hash := range hashes {
					if hash == (web3.Hash{}) {
						continue
					}
					receipt, err := t.Eth().GetTransactionReceipt(hash)
					if err != nil {
						log.Errorf("get tx receipt: %w", err)
						continue
					}
					headBlock, err := t.Eth().BlockNumber()
					if err != nil {
						log.Errorf("get blockNumber: %w", err)
						continue
					}
					if receipt.Status == 1 && receipt.BlockNumber+t.confirmHeight <= headBlock { // confirmed
						confirmCh <- hash
						return
					}
				}

			}
		}
	}

	go listen(hashCh)

	first := true

	for {
		select {
		case <-timer.C:
			timer.Reset(5 * time.Minute)
			if !first {
				// time to reconstruct
				if err := t.ReConstruct(txp); err != nil {
					log.Errorf("reconstruct tx: %w", err)
					continue
				}
			} else {
				first = false
			}
			hash, err := t.Eth().SendRawTransaction(txp.Tx.MarshalRLP())
			if err != nil {
				log.Errorf("sendRawTransaction: %w", err)
				continue
			}
			hashCh <- hash
		case hash := <-confirmCh:
			//confirmed
			log.Infof("tx: %s confirmed", hash)
			return
		case <-ctx.Done():
			log.Error("timeout")
			return
		}

	}

}

func (t *TxManager) AsyncSendTx(tx *web3.Transaction, priceLimits ...uint64) error {
	nextNonce, err := t.nextNonce()
	if err != nil {
		return err
	}
	t.Nonce = nextNonce
	tx = t.SignTx(tx)

	priceLimit := uint64(0)
	if len(priceLimits) > 0 {
		priceLimit = priceLimits[0]
	}
	ctx, _ := context.WithTimeout(context.Background(), 6*time.Hour)

	go t.sendTx(ctx, &TxWithPriceLimit{tx, priceLimit})
	return nil
}
