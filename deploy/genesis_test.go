package deploy

import (
	"fmt"
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/config"
)

func TestBuildL2GenesisData(t *testing.T) {
	conf := &config.L2GenesisConfig{
		FeeCollectorOwner: web3.Address{1, 2, 3},
		BridgeBalance:     10000000000,
		L2ContractAddressConfig: &config.L2ContractAddressConfig{
			ProxyAdmin:               web3.Address{1, 2, 3},
			L2FeeCollector:           web3.Address{4, 5, 6},
			L2CrossLayerWitness:      web3.Address{7, 8, 9},
			L2CrossLayerWitnessLogic: web3.Address{10, 11, 12},
			L2StandardBridge:         web3.Address{13, 14, 15},
			L2StandardBridgeLogic:    web3.Address{16, 17, 18},
		},
	}
	data := BuildL2GenesisData(conf, web3.Address{19})
	fmt.Println(utils.JsonString(data))
}
