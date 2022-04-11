// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import { Types } from "../libraries/Types.sol";
import { IStateCommitChain } from "../interfaces/IStateCommitChain.sol";
import { ICanonicalTransactionChain } from "../interfaces/ICanonicalTransactionChain.sol";
import { IStakingManager } from "../interfaces/IStakingManager.sol";
import { IChainStorageContainer } from "../interfaces/IChainStorageContainer.sol";
import "../interfaces/IChallengeFactory.sol";

contract StateCommitChain is IStateCommitChain {
    using Types for Types.StateInfo;
    IStakingManager stakingManager;
    IChallengeFactory challengeFactory;
    ICanonicalTransactionChain ctc;
    //the window to fraud proof
    uint256 public FRAUD_PROOF_WINDOW;
    //store state info as a chain
    bytes32[] stateChain;
    //the total num of states in state chain, it we cut the state chain simply change this num
    uint64 public override chainHeight;

    constructor(
        address _stakingManager,
        address _challengeFactory,
        address _ctc,
        uint256 _fraudProofWindow
    ) {
        stakingManager = IStakingManager(_stakingManager);
        challengeFactory = IChallengeFactory(_challengeFactory);
        ctc = ICanonicalTransactionChain(_ctc);
        FRAUD_PROOF_WINDOW = _fraudProofWindow;
    }

    function insideFraudProofWindow(Types.StateInfo memory _stateInfo) public view returns (bool _inside) {
        return (_stateInfo.timestamp + FRAUD_PROOF_WINDOW) > block.timestamp;
    }

    function verifyStateInfo(Types.StateInfo memory _stateInfo) public view returns (bool) {
        return _stateInfo.index < chainHeight && stateChain[_stateInfo.index] == _stateInfo.hash();
    }

    function appendStates(bytes32[] memory _blockHashes, uint64 _totalStates) public {
        //in case of duplicated
        require(_totalStates == chainHeight, "current length not equal, maybe others already appended");

        // Proposers must have previously staked at the BondManager
        require(stakingManager.isStaking(msg.sender), "Proposer should be staking");

        require(_blockHashes.length > 0, "no block hashes");

        require(
            chainHeight + _blockHashes.length <= ctc.chainHeight(),
            "Number of state info cannot exceed the tx chain height."
        );
        uint64 _now = uint64(block.timestamp);
        Types.StateInfo memory _stateInfo;
        for (uint256 i = 0; i < _blockHashes.length; i++) {
            _stateInfo.blockHash = _blockHashes[i];
            _stateInfo.timestamp = _now;
            _stateInfo.proposer = msg.sender;
            if (i + chainHeight < stateChain.length) {
                //has some unused storage, reuse it
                stateChain[i + chainHeight] = _stateInfo.hash();
            } else {
                //new storage
                stateChain.push(_stateInfo.hash());
            }
        }
        chainHeight += uint64(_blockHashes.length);
        emit Appended(_totalStates, _blockHashes, msg.sender, _now);
    }

    function deleteState(Types.StateInfo memory _stateInfo) public {
        require(challengeFactory.isChallengeContract(msg.sender), "only permitted by challenge contract");
        require(verifyStateInfo(_stateInfo), "invalid state info");
        require(insideFraudProofWindow(_stateInfo), "State info can only be deleted within the fraud proof window.");
        chainHeight = _stateInfo.index;
        emit Deleted(_stateInfo.index, _stateInfo.blockHash);
    }
}
