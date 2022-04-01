// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IStakingManager.sol";
import "@openzeppelin/contracts/interfaces/IERC20.sol";
import "../interfaces/IStateCommitChain.sol";
import "../interfaces/IChallengeFactory.sol";
import "../libraries/OVMCodec.sol";

contract StakingManager is IStakingManager {
    address private DAOAddress;
    IChallengeFactory challengeFactory;
    IStateCommitChain public scc;
    IERC20 public override token;
    mapping(address => StakingInfo) stakingInfos;
    //price should never change, unless every stakingInfo record the relating info of price.
    uint256 public price;

    constructor(
        address _DAOAddress,
        address _challengeFactory,
        address _stateCommitChain,
        address _erc20,
        uint256 _price
    ) {
        DAOAddress = _DAOAddress;
        challengeFactory = IChallengeFactory(_challengeFactory);
        scc = IStateCommitChain(_stateCommitChain);
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
        senderStake.needConfirmedBlock = scc.getCurrentBlockHeight();
        emit WithdrawStarted(msg.sender, senderStake.needConfirmedBlock);
    }

    function finalizeWithdrawal(OVMCodec.ChainBatchHeader memory _batchHeader) external override {
        StakingInfo storage senderStake = stakingInfos[msg.sender];
        require(senderStake.state == StakingState.WITHDRAWING, "not in withdrawing");
        _assertBlockIsConfirmed(senderStake.needConfirmedBlock, _batchHeader);
        senderStake.state = StakingState.UNSTAKED;
        token.transfer(msg.sender, price);
        emit WithdrawFinalized(msg.sender, price);
    }

    function slash(
        uint256 _blockHeight,
        bytes32 _stateRoot,
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
            proposerStake.earliestChallengeBlock == 0 || _blockHeight < proposerStake.earliestChallengeBlock,
            "should be smaller than last lash"
        );
        proposerStake.earliestChallengeBlock = _blockHeight;
        proposerStake.earliestChallengeState = _stateRoot;
        //set state to slashing
        proposerStake.state = StakingState.SLASHING;
        emit DepositSlashed(_proposer, msg.sender, _blockHeight, _stateRoot);
    }

    function claim(
        address _proposer,
        OVMCodec.BlockInfo memory _blockInfo,
        OVMCodec.ChainBatchHeader memory _batchHeader,
        OVMCodec.ChainInclusionProof memory _proof
    ) external override {
        StakingInfo storage proposerStake = stakingInfos[_proposer];
        //only challenge.
        require(challengeFactory.isChallengeContract(msg.sender), "only challenge contract permitted");
        require(proposerStake.state == StakingState.SLASHING, "not in slashing");
        uint256 _earliestChallengeBlock = proposerStake.earliestChallengeBlock;
        require(scc.verifyStateCommitment(_blockInfo, _batchHeader, _proof), "incorrect block info");
        _assertBlockIsConfirmed(proposerStake.earliestChallengeBlock, _batchHeader);
        require(_blockInfo.blockHash != proposerStake.earliestChallengeState, "unused challenge");
        token.transfer(msg.sender, price);
        proposerStake.state = StakingState.UNSTAKED;
        emit DepositClaimed(_proposer, msg.sender, price);
    }

    function claimToGovernance(
        address _proposer,
        OVMCodec.BlockInfo memory _blockInfo,
        OVMCodec.ChainBatchHeader memory _batchHeader,
        OVMCodec.ChainInclusionProof memory _proof
    ) external override {
        StakingInfo storage proposerStake = stakingInfos[_proposer];
        require(proposerStake.state == StakingState.SLASHING, "not in slashing");
        uint256 _earliestChallengeBlock = proposerStake.earliestChallengeBlock;
        require(scc.verifyStateCommitment(_blockInfo, _batchHeader, _proof), "incorrect block info");
        _assertBlockIsConfirmed(proposerStake.earliestChallengeBlock, _batchHeader);
        require(_blockInfo.blockHash == proposerStake.earliestChallengeState, "useful challenge");
        token.transfer(DAOAddress, price);
        proposerStake.state = StakingState.UNSTAKED;
        emit DepositClaimed(_proposer, DAOAddress, price);
    }

    function _assertBlockIsConfirmed(uint256 _blockNumber, OVMCodec.ChainBatchHeader memory _batchHeader)
        internal
        view
    {
        require(scc.verifyBatchHeader(_batchHeader), "incorrect batch header");
        require(!scc.insideFraudProofWindow(_batchHeader), "provide batch not confirmed");
        require(
            _batchHeader.prevTotalElements + _batchHeader.batchSize > _blockNumber,
            "can proof block not confirmed yet"
        );
    }
}
