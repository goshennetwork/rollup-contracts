// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import { Types } from "../libraries/Types.sol";
import { IStateCommitChain } from "../interfaces/IStateCommitChain.sol";
import { ICanonicalTransactionChain } from "../interfaces/ICanonicalTransactionChain.sol";
import { IStakingManager } from "../interfaces/IStakingManager.sol";
import "../interfaces/IChallengeFactory.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IChainStorageContainer.sol";

contract StateCommitChain is IStateCommitChain {
    using Types for Types.StateInfo;
    IAddressResolver addressResolver;
    //the window to fraud proof
    uint256 public FRAUD_PROOF_WINDOW;

    constructor(address _addressResolver, uint256 _fraudProofWindow) {
        addressResolver = IAddressResolver(_addressResolver);
        FRAUD_PROOF_WINDOW = _fraudProofWindow;
    }

    function isStateConfirmed(Types.StateInfo memory _stateInfo) public view returns (bool _confirmed) {
        return (_stateInfo.timestamp + FRAUD_PROOF_WINDOW) <= block.timestamp;
    }

    function verifyStateInfo(Types.StateInfo memory _stateInfo) public view returns (bool) {
        IChainStorageContainer _chain = addressResolver.sccContainer();
        return _stateInfo.index < _chain.chainSize() && _chain.get(_stateInfo.index) == _stateInfo.hash();
    }

    function appendStates(bytes32[] memory _blockHashes, uint64 _totalStates) public {
        IChainStorageContainer _chain = addressResolver.sccContainer();
        //in case of duplicated
        require(_totalStates == _chain.chainSize(), "current length not equal, maybe others already appended");

        // Proposers must in staking
        require(addressResolver.stakingManager().isStaking(msg.sender), "Proposer should be staking");

        require(_blockHashes.length > 0, "no block hashes");

        require(
            _chain.chainSize() + _blockHashes.length <= addressResolver.ctc().chainHeight(),
            "Number of state info cannot exceed the tx chain height."
        );
        uint64 _now = uint64(block.timestamp);
        Types.StateInfo memory _stateInfo;
        for (uint256 i = 0; i < _blockHashes.length; i++) {
            _stateInfo.blockHash = _blockHashes[i];
            _stateInfo.timestamp = _now;
            _stateInfo.proposer = msg.sender;
            _stateInfo.index = _totalStates;
            _chain.append(_stateInfo.hash());
            _totalStates++;
        }
        emit StateAppended(_totalStates, _blockHashes, msg.sender, _now);
    }

    function rollbackStateBefore(Types.StateInfo memory _stateInfo) public {
        require(
            addressResolver.challengeFactory().isChallengeContract(msg.sender),
            "only permitted by challenge contract"
        );
        require(verifyStateInfo(_stateInfo), "invalid state info");
        require(!isStateConfirmed(_stateInfo), "State info can only be deleted without confirmed");
        addressResolver.sccContainer().resize(_stateInfo.index);
        emit StateRolledBackBefore(_stateInfo.index, _stateInfo.blockHash);
    }

    function chainHeight() public view returns (uint64) {
        return addressResolver.sccContainer().chainSize();
    }
}
