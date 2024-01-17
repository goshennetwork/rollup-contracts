// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../libraries/Types.sol";
import "../interfaces/IRollupStateChain.sol";
import "../interfaces/IRollupInputChain.sol";
import "../interfaces/IStakingManager.sol";
import "../interfaces/IChallengeFactory.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IChainStorageContainer.sol";

contract RollupStateChain is IRollupStateChain, Initializable {
    using Types for Types.StateInfo;

    IAddressResolver public resolver;
    //the window to fraud proof
    uint256 public fraudProofWindow;

    function initialize(address _addressResolver, uint256 _fraudProofWindow) public initializer {
        resolver = IAddressResolver(_addressResolver);
        fraudProofWindow = _fraudProofWindow;
    }

    function isStateConfirmed(Types.StateInfo memory _stateInfo) public view returns (bool _confirmed) {
        return (_stateInfo.timestamp + fraudProofWindow) <= block.timestamp;
    }

    function verifyStateInfo(Types.StateInfo memory _stateInfo) public view returns (bool) {
        IChainStorageContainer _chain = resolver.rollupStateChainContainer();
        return _stateInfo.index < _chain.chainSize() && _chain.get(_stateInfo.index) == _stateInfo.hash();
    }

    function appendStateBatch(bytes32[] memory _blockHashes, uint64 _startAt) public {
        require(resolver.whitelist().canPropose(msg.sender), "only proposer");
        IChainStorageContainer _chain = resolver.rollupStateChainContainer();
        // in case of duplicated
        require(_startAt == _chain.chainSize(), "start pos mismatch");

        // Proposers must in staking
        require(resolver.stakingManager().isStaking(msg.sender), "unstaked");
        require(_blockHashes.length > 0, "no block hashes");

        require(
            _chain.chainSize() + _blockHashes.length <= resolver.rollupInputChain().chainHeight(),
            "exceed input chain height"
        );
        uint64 _now = uint64(block.timestamp);
        Types.StateInfo memory _stateInfo;

        uint64 _pendingIndex = _startAt;
        _stateInfo.timestamp = _now;
        _stateInfo.proposer = msg.sender;
        for (uint256 i = 0; i < _blockHashes.length; i++) {
            _stateInfo.blockHash = _blockHashes[i];
            _stateInfo.index = _pendingIndex;
            _chain.append(_stateInfo.hash());
            _pendingIndex++;
        }
        emit StateBatchAppended(msg.sender, _startAt, _now, _blockHashes);
    }

    //must check not confirmed yet
    function rollbackStateBefore(Types.StateInfo memory _stateInfo) public {
        require(
            resolver.challengeFactory().isChallengeContract(msg.sender), "only permitted by challenge contract"
        );
        require(verifyStateInfo(_stateInfo), "invalid state info");
        require(!isStateConfirmed(_stateInfo), "state confirmed");
        resolver.rollupStateChainContainer().resize(_stateInfo.index);
        emit StateRollbacked(_stateInfo.index, _stateInfo.blockHash);
    }

    function totalSubmittedState() external view returns (uint64) {
        return resolver.rollupStateChainContainer().chainSize();
    }
}
