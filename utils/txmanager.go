package utils

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/laizy/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
)

var ERR_NONCE = errors.New("nonce error")

//TxManager concurrent use
type TxManager struct {
	*contract.Signer
	confirmHeight uint64
	/// if set, only lower price will send tx
	priceLimit *big.Int
}

func NewTxManager(signer *contract.Signer, confirmHeight uint64, priceLimit ...*big.Int) *TxManager {
	r := &TxManager{signer, confirmHeight, nil}
	if len(priceLimit) > 0 {
		r.priceLimit = new(big.Int).Set(priceLimit[0])
	}
	return r
}

func (t *TxManager) BlockNumber() uint64 {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for range timer.C {
		timer.Reset(5 * time.Second)
		n, err := t.Client.Eth().BlockNumber()
		if err == nil {
			return n
		}
		log.Error("get blockNumber", "err", err)
		continue
	}
	panic(1)
}

func (t *TxManager) GetNonce() uint64 {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for range timer.C {
		timer.Reset(5 * time.Second)
		n, err := t.Client.Eth().GetNonce(t.Address(), web3.Latest)
		if err == nil {
			return n
		}
		log.Error("get nonce", "err", err)
		continue
	}
	panic(1)
}

func (t *TxManager) EstimateGas(tx *web3.Transaction) (uint64, error) {
	return t.Client.Eth().EstimateGas(tx.ToCallMsg())
}

func (t *TxManager) GetPrice() uint64 {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for range timer.C {
		timer.Reset(10 * time.Second)
		p, err := t.Client.Eth().GasPrice()
		if err == nil {
			if t.priceLimit != nil && p > t.priceLimit.Uint64() {
				//if price above wanted price just wait
				log.Infof("current price too high, price limit: %f gwei, current price: %f gwei", ToGwei(t.priceLimit), ToGwei(new(big.Int).SetUint64(p)))
				continue
			}
			return p
		}
		log.Error("get gasPrice", "err", err)
		continue
	}
	panic(1)
}

func (t *TxManager) WaitAndChangeTxn(txn *contract.Txn, flexNonce bool, reGas ...bool) (web3.Hash, error) {
	return t.WaitAndChange(txn.MustToTransaction(), flexNonce, reGas...)
}

// ReContruct will reContruct tx with new gasPrice, if flexNonce set to true, will update tx nonce, if want to regas the tx, it will update tx gaslimit
// this function will validate tx locally, if failed, just return error
func (t *TxManager) ReConstruct(tx *web3.Transaction, flexNonce bool, reGas ...bool) error {
	nonce := t.GetNonce()
	if !flexNonce && nonce != tx.Nonce { // nonce not equal, maybe tx is confirmed or user send another tx.
		return ERR_NONCE
	}
	//flex nonce just reset tx nonce to new nonce
	tx.Nonce = nonce
	tx.GasPrice = t.GetPrice()
	if len(reGas) > 0 && reGas[0] { //re calc the gas
		gas, err := t.EstimateGas(tx)
		if err != nil { /// inner execute error or maybe just network error.todo: split 2 kind of error
			return err
		}
		tx.Gas = gas
	} else { //just execute locally to check whether success
		//now try to execute locally, if tx is failed, just return error, indicating this tx will fail
		result, _ := t.ExecuteTxn(tx)
		if result.Err != nil {
			b, _ := tx.MarshalJSON()
			log.Errorf("tx execute failed, tx: %s", string(b))
			return fmt.Errorf("execution reverted: %s", result.RevertReason)
		}
	}
	return nil
}

///fixme: do not use tx.Hash() to get hash, it has cache.
func (t *TxManager) WaitAndChange(tx *web3.Transaction, flexNonce bool, reGas ...bool) (web3.Hash, error) {
	/// make sure tx nonce is confirmed, otherwise may stuck for previous tx.
	tx.Nonce = t.GetNonce()
	tx = t.SignTx(tx)
	var err error
	var hash web3.Hash
	if hash, err = t.Client.Eth().SendRawTransaction(tx.MarshalRLP()); err != nil {
		//tx err
		return hash, fmt.Errorf("sendRawTransaction: %w", err)
	}
	b, _ := tx.MarshalJSON()
	log.Info("send tx", "hash", hash, "tx", string(b))

	ddl := time.NewTicker(3 * time.Minute)
	defer ddl.Stop()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	nonceErr := 0
	for {
		select {
		case <-ddl.C:
			err = t.ReConstruct(tx, flexNonce, reGas...)
			if err != nil {
				if !errors.Is(err, ERR_NONCE) {
					return web3.Hash{}, fmt.Errorf("reConstruct tx: %w", err)
				}
				/// if nonce error more than 1, just return error
				if nonceErr >= 1 {
					return hash, ERR_NONCE
				}
				nonceErr++
				/// maybe tx already confirmed, check it
				continue
			}
			if newHash, err := t.Client.Eth().SendRawTransaction(t.SignTx(tx).MarshalRLP()); err != nil {
				//tx err
				log.Errorf("sendRawTransaction: %w", err)
				continue
			} else {
				//now replace success, just listen this hash
				hash = newHash
			}
			b, _ := tx.MarshalJSON()
			log.Info("replay transaction", "hash", hash, "tx", string(b))
		case <-ticker.C:
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
