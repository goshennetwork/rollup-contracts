package staking

import (
	"math/big"

	"github.com/laizy/web3/utils/common"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/utils/common/math"
	"github.com/ontology-layer-2/rollup-contracts/tests/contracts/build"
)

type StakingManager web3.Address

func (self StakingManager) Address() web3.Address {
	return web3.Address(self)
}

//function deposit() external;
func (self StakingManager) Deposit(sender evm.AccountRef, vm *evm.EVM) error {
	input := build.StakingManagerAbi().Methods["deposit"].ID()
	_, _, err := vm.Call(sender, self.Address(), input, math.MaxUint64, new(big.Int))
	return err
}

//function isStaking(address _who) external view returns (bool);
func (self StakingManager) IsStaking(sender evm.AccountRef, vm *evm.EVM, who web3.Address) (out bool, err error) {
	method := build.StakingManagerAbi().Methods["isStaking"]
	input := method.MustEncodeIDAndInput(who)
	ret, _, err := vm.Call(sender, self.Address(), input, math.MaxUint64, new(big.Int))
	if err != nil {
		return false, err
	}
	i, err := method.Outputs.Decode(ret)
	if err != nil {
		return false, err
	}
	err = contracts.Decode(i, &out)
	return
}

//function startWithdrawal() external;
func (self StakingManager) StartWithdrawal(sender evm.AccountRef, vm *evm.EVM) error {
	input := build.StakingManagerAbi().Methods["startWithdrawal"].ID()
	_, _, err := vm.Call(sender, self.Address(), input, math.MaxUint64, new(big.Int))
	return err
}

//function finalizeWithdrawal(Types.StateInfo memory _stateInfo) external;
func (self StakingManager) FinalizeWithdrawal(sender evm.AccountRef, vm *evm.EVM, stateInfo build.StateInfo) error {
	input := build.StakingManagerAbi().Methods["finalizeWithdrawal"].MustEncodeIDAndInput(stateInfo)
	_, _, err := vm.Call(sender, self.Address(), input, math.MaxUint64, new(big.Int))
	return err
}

//function slash(
//        uint64 _chainHeight,
//        bytes32 _stateRoot,
//        address _proposer
//    ) external;
func (self StakingManager) Slash(sender evm.AccountRef, vm *evm.EVM, chainHeight uint64, stateRoot common.Hash, proposer web3.Address) error {
	method := build.StakingManagerAbi().Methods["slash"]
	input := method.MustEncodeIDAndInput(chainHeight, stateRoot, proposer)
	_, _, err := vm.Call(sender, self.Address(), input, math.MaxUint64, new(big.Int))
	return err
}

//function claim(address _proposer, Types.StateInfo memory _stateInfo) external;
func (self StakingManager) Claim(sender evm.AccountRef, vm *evm.EVM, proposer web3.Address, stateInfo build.StateInfo) error {
	input := build.StakingManagerAbi().Methods["claim"].MustEncodeIDAndInput(proposer, stateInfo)
	_, _, err := vm.Call(sender, self.Address(), input, math.MaxUint64, new(big.Int))
	return err
}

// function claimToGovernance(address _proposer, Types.StateInfo memory _stateInfo) external;
func (self StakingManager) ClaimToGovernance(sender evm.AccountRef, vm *evm.EVM, proposer web3.Address, stateInfo build.StateInfo) error {
	input := build.StakingManagerAbi().Methods["claimToGovernance"].MustEncodeIDAndInput(proposer, stateInfo)
	_, _, err := vm.Call(sender, self.Address(), input, math.MaxUint64, new(big.Int))
	return err
}
