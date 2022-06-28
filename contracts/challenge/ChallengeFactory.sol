// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IChallengeFactory.sol";
import "../interfaces/IAddressResolver.sol";
import "./Challenge.sol";
import "@openzeppelin/contracts/proxy/beacon/BeaconProxy.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/IAddressManager.sol";

contract ChallengeFactory is IChallengeFactory, Initializable {
    using Types for Types.StateInfo;
    mapping(address => bool) contracts;
    mapping(bytes32 => address) challengedStates;
    IAddressResolver resolver;
    uint256 public blockLimitPerRound;
    address public override challengeBeacon;
    uint256 public challengerDeposit;

    function initialize(
        IAddressResolver _resolver,
        address _beacon,
        uint256 _blockLimitPerRound,
        uint256 _challengerDeposit
    ) public initializer {
        resolver = _resolver;
        challengeBeacon = _beacon;
        blockLimitPerRound = _blockLimitPerRound;
        challengerDeposit = _challengerDeposit;
    }

    function newChallange(
        //when create, creator should deposit at this contract.
        Types.StateInfo memory _challengedStateInfo,
        Types.StateInfo memory _parentStateInfo
    ) public {
        require(resolver.dao().challengerWhitelist(msg.sender), "only challenger");
        bytes32 _hash = _challengedStateInfo.hash();
        require(challengedStates[_hash] != address(0), "already challenged");
        require(resolver.rollupStateChain().verifyStateInfo(_challengedStateInfo), "wrong stateInfo");
        require(!resolver.rollupStateChain().isStateConfirmed(_challengedStateInfo), "state confirmed");
        require(resolver.rollupStateChain().verifyStateInfo(_parentStateInfo), "wrong stateInfo");
        require(_parentStateInfo.index + 1 == _challengedStateInfo.index, "wrong parent stateInfo");
        bytes32 _inputHash = resolver.rollupInputChain().getInputHash(_challengedStateInfo.index);
        bytes32 _systemStartState = resolver.stateTransition().generateStartState(
            _inputHash,
            _challengedStateInfo.index,
            _parentStateInfo.blockHash
        );
        bytes memory _data;
        address newChallenge = address(new BeaconProxy(challengeBeacon, _data));
        contracts[newChallenge] = true;
        challengedStates[_hash] = newChallenge;
        IChallenge(newChallenge).create(
            _systemStartState,
            msg.sender,
            blockLimitPerRound,
            _challengedStateInfo,
            challengerDeposit
        );
        emit ChallengeStarted(
            _challengedStateInfo.index,
            _challengedStateInfo.proposer,
            _systemStartState,
            block.number + blockLimitPerRound,
            newChallenge
        );
    }

    function getChallengedContract(bytes32 _stateInfoHash) public view returns (address) {
        address _c = challengedStates[_stateInfoHash];
        require(_c != address(0), "not challenged");
        return _c;
    }

    function stakingManager() public view returns (IStakingManager) {
        return resolver.stakingManager();
    }

    function executor() public view returns (IStateTransition) {
        return resolver.stateTransition();
    }

    function rollupStateChain() public view returns (IRollupStateChain) {
        return resolver.rollupStateChain();
    }

    function dao() public view returns (address) {
        return address(resolver.dao());
    }

    function isChallengeContract(address _addr) public view returns (bool) {
        return contracts[_addr];
    }
}
