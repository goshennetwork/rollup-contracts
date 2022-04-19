// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IAddressResolver.sol";
import "../interfaces/IAddressManager.sol";
import "./AddressName.sol";

contract AddressResolver is IAddressResolver {
    IAddressManager addressManager;

    constructor(address _addressManager) {
        addressManager = IAddressManager(_addressManager);
    }

    function resolve(string memory _name) public view returns (address) {
        address _addr = addressManager.getAddr(_name);
        require(_addr != address(0), "no name saved");
        return _addr;
    }

    function ctc() public view returns (IRollupInputChain) {
        return IRollupInputChain(resolve(AddressName.CTC));
    }

    function ctcContainer() public view returns (IChainStorageContainer) {
        return IChainStorageContainer(resolve(AddressName.CTC_CONTAINER));
    }

    function scc() public view returns (IRollupStateChain) {
        return IRollupStateChain(resolve(AddressName.SCC));
    }

    function sccContainer() public view returns (IChainStorageContainer) {
        return IChainStorageContainer(resolve(AddressName.SCC_CONTAINER));
    }

    function stakingManager() public view returns (IStakingManager) {
        return IStakingManager(resolve(AddressName.STAKING_MANAGER));
    }

    function challengeFactory() public view returns (IChallengeFactory) {
        return IChallengeFactory(resolve(AddressName.CHALLENGE_FACTORY));
    }

    function l1CrossLayerMessageWitness() public view returns (address) {
        return resolve(AddressName.L1_CROSS_LAYER_MESSAGE_WITNESS);
    }
}
