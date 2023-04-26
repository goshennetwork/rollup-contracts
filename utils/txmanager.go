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

func (t *TxManager) WaitAndChangeTxn(txn *contract.Txn, flexNonce bool, reGas ...bool) (web3.Hash, error) {
	tx, err := txn.ToTransaction()
	if err != nil {
		return web3.Hash{}, err
	}

	return t.WaitAndChange(tx, flexNonce, reGas...)
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

func (t *TxManager) SendTx(tx *web3.Transaction) error {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	defer cancel()

	quitCh, confirmCh := make(chan struct{}, 1), make(chan struct{}, 1)
	listener := func(ctx context.Context, quit <-chan struct{}, confirm chan struct{}, hash web3.Hash) {
		ticker := time.NewTicker(20 * time.Second)
		for {
			select {
			case <-ctx.Done():
				// time out
				log.Errorf("tx: %s timeout", hash)
				return
			case <-quit:
				return
			case <-ticker.C:
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
					confirm <- struct{}{}
					return
				}

			}
		}
	}

	hash, err := t.Eth().SendRawTransaction(tx.MarshalRLP())
	if err != nil {
		return fmt.Errorf("sendRawTransaction: %w", err)
	}

	for range ticker.C {

	}

}

func (t *TxManager) WaitAndChange(tx *web3.Transaction, flexNonce bool, reGas ...bool) (web3.Hash, error) {
	/// make sure tx nonce is confirmed, otherwise may stuck for previous tx.
	tx.Nonce = t.GetNonce()
	//make sure price do not above priceLimit
	if t.priceLimit != nil && tx.GasPrice > t.priceLimit.Uint64() {
		tx.GasPrice = t.priceLimit.Uint64()
	}
	tx = t.SignTx(tx)

	if newTx, err := t.SendTx(tx, flexNonce, reGas...); err != nil { //at this time, nonce error could not indicate the tx confirmed, so just return error
		//tx err
		return tx.Hash(), fmt.Errorf("sendTx: %w", err)
	} else {
		tx = newTx
	}
	// store 2 hash, because async of txpool, maybe miner already chose the tx and is packing, but the txpool have not changed yet
	var sentTx [2]web3.Hash
	sentTx[0] = tx.Hash()

	ddl := time.NewTicker(3 * time.Minute)
	defer ddl.Stop()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	nonceErr := 0
	i := 0
	for {
		select {
		case <-ddl.C:
			newTx, err := t.ReConstruct(tx, flexNonce, reGas...)
			if err != nil {
				if !errors.Is(err, ERR_NONCE) {
					return tx.Hash(), fmt.Errorf("reConstruct tx: %w", err)
				}
				/// if nonce error more than 1, just return error
				if nonceErr >= 1 {
					return tx.Hash(), ERR_NONCE
				}
				nonceErr++
				/// maybe tx already confirmed, check it
				continue
			}

			//recontruct success, do not replace old tx until client accept it
			if newTx, err := t.SendTx(newTx, flexNonce, reGas...); err != nil {
				//tx err
				log.Errorf("sendRawTransaction: %w", err)
				continue
			} else {
				i++
				//now replace success
				tx = newTx
				sentTx[i%2] = tx.Hash()
			}
			log.Info("transaction replayed", "hash", tx.Hash())
		case <-ticker.C:
			for _, hash := range sentTx {
				if hash == (web3.Hash{}) {
					continue
				}
				r, err := t.Client.Eth().GetTransactionReceipt(hash)
				if err != nil {
					log.Error("getTransactionReceipt", "err", err)
					continue
				}

				if r != nil {
					nonceErr = 0 //find out receipt, clean nonce error
					//check confirmed Block
					blockNumber := t.BlockNumber()
					if blockNumber < r.BlockNumber { // wired, maye rollback or rpc client balancing
						continue
					}
					confirms := blockNumber - r.BlockNumber + 1
					if confirms < t.confirmHeight { // not confirmed yet
						log.Warn("tx confirming", "hash", hash, "confirm block number", confirms)
						continue
					}
					log.Infof("tx %s confirmed", hash)
					return hash, nil
				}
			}
		}
	}
}
