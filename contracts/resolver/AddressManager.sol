// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IAddressManager.sol";
import "./AddressName.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IL1StandardBridge.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";
import "../test-helper/TestERC20.sol";
import "../interfaces/IMachineState.sol";

contract AddressManager is IAddressManager, IAddressResolver, OwnableUpgradeable {
    mapping(bytes32 => address) public getAddrByHash;

    function initialize() public initializer {
        __Ownable_init();
    }

    function setAddress(string memory _name, address _addr) public onlyOwner {
        bytes32 _hash = hash(_name);
        address _old = _setAddress(_hash, _addr);
        emit AddressSet(_name, _old, _addr);
    }

    function _setAddress(bytes32 _hash, address _addr) internal returns (address) {
        require(_addr != address(0), "empty addr");
        address _old = getAddrByHash[_hash];
        getAddrByHash[_hash] = _addr;
        return _old;
    }

    function setAddressBatch(string[] calldata _names, address[] calldata _addrs) public onlyOwner {
        uint256 _len = _names.length;
        require(_len == _addrs.length, "length mismatch");
        for (uint256 i = 0; i < _len; i++) {
            string calldata _name = _names[i];
            address _addr = _addrs[i];
            bytes32 _hash = hash(_name);
            address _old = _setAddress(_hash, _addr);
            emit AddressSet(_name, _old, _addr);
        }
    }

    function getAddr(string memory _name) public view returns (address) {
        return getAddrByHash[hash(_name)];
    }

    function resolve(string memory _name) public view returns (address) {
        address _addr = getAddr(_name);
        require(_addr != address(0), "no name saved");
        return _addr;
    }

    function dao() public view returns (IDAO) {
        return IDAO(getAddrByHash[AddressName.DAO_HASH]);
    }

    function rollupInputChain() public view returns (IRollupInputChain) {
        return IRollupInputChain(getAddrByHash[AddressName.ROLLUP_INPUT_CHAIN_HASH]);
    }

    function rollupInputChainContainer() public view returns (IChainStorageContainer) {
        return IChainStorageContainer(getAddrByHash[AddressName.ROLLUP_INPUT_CHAIN_CONTAINER_HASH]);
    }

    function rollupStateChain() public view returns (IRollupStateChain) {
        return IRollupStateChain(getAddrByHash[AddressName.ROLLUP_STATE_CHAIN_HASH]);
    }

    function rollupStateChainContainer() public view returns (IChainStorageContainer) {
        return IChainStorageContainer(getAddrByHash[AddressName.ROLLUP_STATE_CHAIN_CONTAINER_HASH]);
    }

    function stakingManager() public view returns (IStakingManager) {
        return IStakingManager(getAddrByHash[AddressName.STAKING_MANAGER_HASH]);
    }

    function challengeFactory() public view returns (IChallengeFactory) {
        return IChallengeFactory(getAddrByHash[AddressName.CHALLENGE_FACTORY_HASH]);
    }

    function l1CrossLayerWitness() public view returns (IL1CrossLayerWitness) {
        return IL1CrossLayerWitness(getAddrByHash[AddressName.L1_CROSS_LAYER_WITNESS_HASH]);
    }

    function l2CrossLayerWitness() public view returns (IL2CrossLayerWitness) {
        return IL2CrossLayerWitness(getAddrByHash[AddressName.L2_CROSS_LAYER_WITNESS_HASH]);
    }

    function stateTransition() public view returns (IStateTransition) {
        return IStateTransition(getAddrByHash[AddressName.STATE_TRANSITION_HASH]);
    }

    function l1StandardBridge() public view returns (IL1StandardBridge) {
        return IL1StandardBridge(getAddrByHash[AddressName.L1_STANDARD_BRIDGE_HASH]);
    }

    function challengeBeacon() public view returns (UpgradeableBeacon) {
        return UpgradeableBeacon(getAddrByHash[AddressName.CHALLENGE_BEACON_HASH]);
    }

    function feeToken() public view returns (IERC20) {
        return TestERC20(getAddrByHash[AddressName.FEE_TOKEN_HASH]);
    }

    function machineState() public view returns (IMachineState) {
        return IMachineState(getAddrByHash[AddressName.MACHINE_STATE_HASH]);
    }

    function hash(string memory _name) internal pure returns (bytes32) {
        return keccak256(bytes(_name));
    }
}
