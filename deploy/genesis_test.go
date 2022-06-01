package deploy

import (
	"fmt"
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
)

func TestBuildL2GenesisData(t *testing.T) {
	conf := &GenesisConfig{
		FeeCollectorOwner:   web3.Address{1, 2, 3},
		FeeCollector:        web3.Address{4, 5, 6},
		L2CrossLayerWitness: web3.Address{7, 8, 9},
		WitnessBalance:      web3.Ether(10000000000),
	}
	data := BuildL2GenesisData(conf)
	fmt.Println(utils.JsonString(data))
}
