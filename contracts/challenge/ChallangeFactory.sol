pragma solidity ^0.8.0;

import "../interfaces/IChallengeFactory.sol";
import "../interfaces/IAddressResolver.sol";
import "./Challenge.sol";
import "@openzeppelin/contracts/proxy/beacon/IBeacon.sol";
import "@openzeppelin/contracts/proxy/beacon/BeaconProxy.sol";

contract ChallengeFactory is IChallengeFactory, IBeacon {
    mapping(address => bool) contracts;
    mapping(uint64 => address) challengedStates;
    IAddressResolver resolver;
    IChallenge challenge;
    uint256 immutable proposerTimeLimit;
    //fixme: flows need more evaluation.
    uint256 public constant minChallengerDeposit = 0.1 ether;

    constructor(uint256 _proposerTimeLimit) {
        proposerTimeLimit = _proposerTimeLimit;
        challenge = new Challenge();
    }

    function newChallange(
        //when create, creator should deposit at this contract.
        Types.StateInfo memory _challengedStateInfo,
        Types.StateInfo memory _parentStateInfo
    ) public returns (bool) {
        require(challengedStates[_challengedStateInfo.index] != address(0), "already challenged");
        require(resolver.rollupStateChain().verifyStateInfo(_challengedStateInfo), "wrong stateInfo");
        require(resolver.rollupStateChain().isStateConfirmed(_challengedStateInfo), "state confirmed");
        require(resolver.rollupStateChain().verifyStateInfo(_parentStateInfo), "wrong stateInfo");
        require(resolver.rollupStateChain().isStateConfirmed(_parentStateInfo), "state confirmed");
        require(_parentStateInfo.index + 1 == _challengedStateInfo.index, "wrong parent stateInfo");
        bytes32 _inputHash = resolver.rollupInputChain().getInputHash(_challengedStateInfo.index);
        bytes32 _systemStartState = resolver.stateTransition().generateStartState(
            _inputHash,
            _challengedStateInfo.index,
            _parentStateInfo.blockHash
        );
        bytes memory _data;
        address newChallenge = address(new BeaconProxy(address(this), _data));
        contracts[newChallenge] = true;
        challengedStates[_challengedStateInfo.index] = newChallenge;
        IChallenge(newChallenge).create(
            _challengedStateInfo.index,
            _systemStartState,
            msg.sender,
            proposerTimeLimit,
            _challengedStateInfo,
            minChallengerDeposit
        );
        return true;
    }

    function getChallengedContract(uint64 _stateIndex) public view returns (address) {
        address _c = challengedStates[_stateIndex];
        require(_c != address(0), "not challenged");
        return _c;
    }

    function implementation() public view returns (address) {
        return address(challenge);
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
        return resolver.dao();
    }

    function isChallengeContract(address _addr) public view returns (bool) {
        return contracts[_addr];
    }
}
