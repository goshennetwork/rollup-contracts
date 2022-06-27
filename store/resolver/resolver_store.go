package resolver

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type AddressManager struct {
	store schema.KeyValueDB
}

func NewStore(db schema.KeyValueDB) *AddressManager {
	return &AddressManager{
		store: db,
	}
}

func (self *AddressManager) StoreLastL1BlockHeight(lastEndHeight uint64) {
	self.store.Put(schema.AddressManagerLastL1BlockHeightKey, codec.NewZeroCopySink(nil).WriteUint64(lastEndHeight).Bytes())
}

func (self *AddressManager) GetLastL1BlockHeight() (uint64, error) {
	v, err := self.store.Get(schema.AddressManagerLastL1BlockHeightKey)
	if err != nil {
		return 0, err
	}
	if len(v) == 0 {
		return 0, nil
	}
	return codec.NewZeroCopySource(v).ReadUint64()
}

func (self *AddressManager) SetAddress(name string, addr web3.Address) {
	self.store.Put(genAddrKey(name), addr.Bytes())
}

func (self *AddressManager) GetAddress(name string) (addr web3.Address, err error) {
	var data []byte
	data, err = self.store.Get(genAddrKey(name))
	if err != nil {
		return
	}
	if len(data) == 0 {
		err = schema.ErrNotFound
		return
	}
	source := codec.NewZeroCopySource(data)
	return source.ReadAddress()
}

func genAddrKey(name string) []byte {
	var k [33]byte
	k[0] = schema.AddressNamePrefix
	copy(k[1:], crypto.Keccak256([]byte(name)))
	return k[:]
}
