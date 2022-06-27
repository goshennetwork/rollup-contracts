package binding

import (
	"io"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils/codec"
)

type CrossLayerInfo struct {
	L1Token web3.Address
	L2Token web3.Address
	From    web3.Address
	To      web3.Address
	Amount  *big.Int
	Data    []byte
}

func (s *CrossLayerInfo) Serialization(sink *codec.ZeroCopySink) {
	sink.WriteAddress(s.L1Token)
	sink.WriteAddress(s.L2Token)
	sink.WriteAddress(s.From)
	sink.WriteAddress(s.To)
	sink.WriteVarBytes(s.Amount.Bytes())
	sink.WriteVarBytes(s.Data)
}

func (s *CrossLayerInfo) Deserialization(source *codec.ZeroCopySource) (err error) {
	s.L1Token, err = source.ReadAddress()
	if err != nil {
		return err
	}
	s.L2Token, err = source.ReadAddress()
	if err != nil {
		return err
	}
	s.From, err = source.ReadAddress()
	if err != nil {
		return err
	}
	s.To, err = source.ReadAddress()
	if err != nil {
		return err
	}
	var amountData []byte
	amountData, err = source.ReadVarBytes()
	if err != nil {
		return err
	}
	s.Amount = new(big.Int).SetBytes(amountData)
	s.Data, err = source.ReadVarBytes()
	return err
}

type CrossLayerInfos []*CrossLayerInfo

func (s CrossLayerInfos) Serialization(sink *codec.ZeroCopySink) {
	sink.WriteInt64(int64(len(s)))
	for _, evt := range s {
		sink.WriteVarBytes(codec.SerializeToBytes(evt))
	}
}

func DeserializationCrossLayerInfos(source *codec.ZeroCopySource) (s CrossLayerInfos, err error) {
	length, eof := source.NextInt64()
	if eof {
		return nil, io.ErrUnexpectedEOF
	}
	for i := int64(0); i < length; i++ {
		data, err := source.ReadVarBytes()
		if err != nil {
			return nil, err
		}
		evt := &CrossLayerInfo{}
		err = evt.Deserialization(codec.NewZeroCopySource(data))
		if err != nil {
			return nil, err
		}
		s = append(s, evt)
	}
	return s, nil
}

type TokenCrossLayerInfo interface {
	GetTokenCrossInfo() *CrossLayerInfo
}

func (evt *DepositFailedEvent) GetTokenCrossInfo() *CrossLayerInfo {
	return &CrossLayerInfo{
		L1Token: evt.L1Token,
		L2Token: evt.L2Token,
		From:    evt.From,
		To:      evt.To,
		Amount:  evt.Amount,
		Data:    evt.Data,
	}
}

func (evt *DepositFinalizedEvent) GetTokenCrossInfo() *CrossLayerInfo {
	return &CrossLayerInfo{
		L1Token: evt.L1Token,
		L2Token: evt.L2Token,
		From:    evt.From,
		To:      evt.To,
		Amount:  evt.Amount,
		Data:    evt.Data,
	}
}

func (evt *ERC20DepositInitiatedEvent) GetTokenCrossInfo() *CrossLayerInfo {
	return &CrossLayerInfo{
		L1Token: evt.L1Token,
		L2Token: evt.L2Token,
		From:    evt.From,
		To:      evt.To,
		Amount:  evt.Amount,
		Data:    evt.Data,
	}
}

func (evt *ERC20WithdrawalFinalizedEvent) GetTokenCrossInfo() *CrossLayerInfo {
	return &CrossLayerInfo{
		L1Token: evt.L1Token,
		L2Token: evt.L2Token,
		From:    evt.From,
		To:      evt.To,
		Amount:  evt.Amount,
		Data:    evt.Data,
	}
}

func (evt *WithdrawalInitiatedEvent) GetTokenCrossInfo() *CrossLayerInfo {
	return &CrossLayerInfo{
		L1Token: evt.L1Token,
		L2Token: evt.L2Token,
		From:    evt.From,
		To:      evt.To,
		Amount:  evt.Amount,
		Data:    evt.Data,
	}
}
