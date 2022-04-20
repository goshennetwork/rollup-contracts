// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IStakingManager.sol";
import "@openzeppelin/contracts/interfaces/IERC20.sol";
import "../interfaces/IRollupStateChain.sol";
import "../interfaces/IChallengeFactory.sol";
import "../libraries/Types.sol";

contract StakingManager is IStakingManager {
    address private DAOAddress;
    IChallengeFactory challengeFactory;
    IRollupStateChain public rollupStateChain;
    IERC20 public override token;
    mapping(address => StakingInfo) stakingInfos;
    //price should never change, unless every stakingInfo record the relating info of price.
    uint256 public price;

    constructor(
        address _DAOAddress,
        address _challengeFactory,
        address _rollupStateChain,
        address _erc20,
        uint256 _price
    ) {
        DAOAddress = _DAOAddress;
        challengeFactory = IChallengeFactory(_challengeFactory);
        rollupStateChain = IRollupStateChain(_rollupStateChain);
        token = IERC20(_erc20);
        price = _price;
    }

    function deposit() external override {
        StakingInfo storage senderStaking = stakingInfos[msg.sender];
        require(senderStaking.state == StakingState.UNSTAKED, "only unstacked user can deposit");
        token.transferFrom(msg.sender, address(this), price);
        senderStaking.state = StakingState.STAKING;
        emit Deposited(msg.sender, price);
    }

    function isStaking(address _who) external view override returns (bool) {
        return stakingInfos[_who].state == StakingState.STAKING;
    }

    function startWithdrawal() external override {
        StakingInfo storage senderStake = stakingInfos[msg.sender];
        require(senderStake.state == StakingState.STAKING, "not in staking");
        senderStake.state = StakingState.WITHDRAWING;
        senderStake.needConfirmedHeight = rollupStateChain.totalSubmittedState();
        emit WithdrawStarted(msg.sender, senderStake.needConfirmedHeight);
    }

    function finalizeWithdrawal(Types.StateInfo memory _stateInfo) external override {
        StakingInfo storage senderStake = stakingInfos[msg.sender];
        require(senderStake.state == StakingState.WITHDRAWING, "not in withdrawing");
        _assertStateIsConfirmed(senderStake.needConfirmedHeight, _stateInfo);
        senderStake.state = StakingState.UNSTAKED;
        token.transfer(msg.sender, price);
        emit WithdrawFinalized(msg.sender, price);
    }

    function slash(
        uint64 _chainHeight,
        bytes32 _blockHash,
        address _proposer
    ) external override {
        StakingInfo storage proposerStake = stakingInfos[_proposer];
        //only challenge.
        require(challengeFactory.isChallengeContract(msg.sender), "only challenge contract permitted");
        //unstaked is not allowed
        require(proposerStake.state != StakingState.UNSTAKED, "unStaked unexpected");
        if (proposerStake.firstSlashTime == 0) {
            proposerStake.firstSlashTime = block.timestamp;
        }
        require(
            proposerStake.earliestChallengeHeight == 0 || _chainHeight < proposerStake.earliestChallengeHeight,
            "should be smaller than last lash"
        );
        proposerStake.earliestChallengeHeight = _chainHeight;
        proposerStake.earliestChallengeBlockHash = _blockHash;
        //set state to slashing
        proposerStake.state = StakingState.SLASHING;
        emit DepositSlashed(_proposer, msg.sender, _chainHeight, _blockHash);
    }

    function claim(address _proposer, Types.StateInfo memory _stateInfo) external override {
        StakingInfo storage proposerStake = stakingInfos[_proposer];
        //only challenge.
        require(challengeFactory.isChallengeContract(msg.sender), "only challenge contract permitted");
        require(proposerStake.state == StakingState.SLASHING, "not in slashing");
        require(rollupStateChain.verifyStateInfo(_stateInfo), "incorrect state info");
        _assertStateIsConfirmed(proposerStake.earliestChallengeHeight, _stateInfo);
        require(_stateInfo.blockHash != proposerStake.earliestChallengeBlockHash, "unused challenge");
        token.transfer(msg.sender, price);
        proposerStake.state = StakingState.UNSTAKED;
        emit DepositClaimed(_proposer, msg.sender, price);
    }

    function claimToGovernance(address _proposer, Types.StateInfo memory _stateInfo) external override {
        StakingInfo storage proposerStake = stakingInfos[_proposer];
        require(proposerStake.state == StakingState.SLASHING, "not in slashing");
        require(rollupStateChain.verifyStateInfo(_stateInfo), "incorrect state info");
        _assertStateIsConfirmed(proposerStake.earliestChallengeHeight, _stateInfo);
        require(_stateInfo.blockHash == proposerStake.earliestChallengeBlockHash, "useful challenge");
        token.transfer(DAOAddress, price);
        proposerStake.state = StakingState.UNSTAKED;
        emit DepositClaimed(_proposer, DAOAddress, price);
    }

    function _assertStateIsConfirmed(uint256 _index, Types.StateInfo memory _stateInfo) internal view {
        require(rollupStateChain.verifyStateInfo(_stateInfo), "incorrect state info");
        require(rollupStateChain.isStateConfirmed(_stateInfo), "provided state not confirmed");
        require(_stateInfo.index == _index, "should provide wanted state info");
    }
}
