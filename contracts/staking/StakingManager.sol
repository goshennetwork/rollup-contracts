// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/IStakingManager.sol";
import "@openzeppelin/contracts/interfaces/IERC20.sol";
import "../libraries/Types.sol";
import "../libraries/console.sol";
import "../interfaces/IAddressResolver.sol";

contract StakingManager is IStakingManager, Initializable {
    IAddressResolver resolver;
    mapping(address => StakingInfo) public getStakingInfo;
    //price should never change, unless every stakingInfo record the relating info of price.
    uint256 public price;

    modifier onlyChallenge() {
        require(resolver.challengeFactory().isChallengeContract(msg.sender), "only challenge contract permitted");
        _;
    }

    function initialize(address _resolver, uint256 _price) public initializer {
        resolver = IAddressResolver(_resolver);
        price = _price;
    }

    /// The token address used for staking.
    function token() public view returns (IERC20) {
        return resolver.feeToken();
    }

    function deposit() external override {
        StakingInfo storage senderStaking = getStakingInfo[msg.sender];
        require(senderStaking.state == StakingState.UNSTAKED, "only unStaked user can deposit");
        require(resolver.feeToken().transferFrom(msg.sender, address(this), price), "transfer failed");
        senderStaking.state = StakingState.STAKING;
        emit Deposited(msg.sender, price);
    }

    function isStaking(address _who) external view override returns (bool) {
        return getStakingInfo[_who].state == StakingState.STAKING;
    }

    function startWithdrawal() external override {
        StakingInfo storage senderStake = getStakingInfo[msg.sender];
        require(senderStake.state == StakingState.STAKING, "not in staking");
        senderStake.state = StakingState.WITHDRAWING;
        senderStake.needConfirmedHeight = resolver.rollupStateChain().totalSubmittedState();
        emit WithdrawStarted(msg.sender, senderStake.needConfirmedHeight);
    }

    function finalizeWithdrawal(Types.StateInfo memory _stateInfo) external override {
        StakingInfo storage senderStake = getStakingInfo[msg.sender];
        require(senderStake.state == StakingState.WITHDRAWING, "not in withdrawing");
        _assertStateIsConfirmed(resolver.rollupStateChain(), senderStake.needConfirmedHeight, _stateInfo);
        senderStake.state = StakingState.UNSTAKED;
        resolver.feeToken().transfer(msg.sender, price);
        emit WithdrawFinalized(msg.sender, price);
    }

    function slash(
        uint64 _chainHeight,
        bytes32 _blockHash,
        address _proposer
    ) external override onlyChallenge {
        StakingInfo storage proposerStake = getStakingInfo[_proposer];
        //unstaked is not allowed
        require(proposerStake.state != StakingState.UNSTAKED, "unStaked unexpected");
        if (proposerStake.firstSlashTime == 0) {
            proposerStake.firstSlashTime = uint64(block.timestamp);
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

    function claim(address _proposer, Types.StateInfo memory _stateInfo) external override onlyChallenge {
        StakingInfo storage proposerStake = getStakingInfo[_proposer];
        require(proposerStake.state == StakingState.SLASHING, "not in slashing");
        IRollupStateChain _stateChain = resolver.rollupStateChain();
        require(_stateChain.verifyStateInfo(_stateInfo), "incorrect state info");
        _assertStateIsConfirmed(_stateChain, proposerStake.earliestChallengeHeight, _stateInfo);
        require(_stateInfo.blockHash != proposerStake.earliestChallengeBlockHash, "unused challenge");
        resolver.feeToken().transfer(msg.sender, price);
        proposerStake.state = StakingState.UNSTAKED;
        //// make info that will effect slash clean.
        proposerStake.earliestChallengeHeight = 0;
        proposerStake.firstSlashTime = 0;
        emit DepositClaimed(_proposer, msg.sender, price);
    }

    function claimToGovernance(address _proposer, Types.StateInfo memory _stateInfo) external override {
        StakingInfo storage proposerStake = getStakingInfo[_proposer];
        require(proposerStake.state == StakingState.SLASHING, "not in slashing");
        IRollupStateChain _stateChain = resolver.rollupStateChain();
        require(_stateChain.verifyStateInfo(_stateInfo), "incorrect state info");
        _assertStateIsConfirmed(_stateChain, proposerStake.earliestChallengeHeight, _stateInfo);
        require(_stateInfo.blockHash == proposerStake.earliestChallengeBlockHash, "useful challenge");
        address _dao = resolver.dao();
        resolver.feeToken().transfer(_dao, price);
        proposerStake.state = StakingState.UNSTAKED;
        //// make info that will effect slash clean.
        proposerStake.earliestChallengeHeight = 0;
        proposerStake.firstSlashTime = 0;
        emit DepositClaimed(_proposer, _dao, price);
    }

    function _assertStateIsConfirmed(
        IRollupStateChain _stateChain,
        uint256 _index,
        Types.StateInfo memory _stateInfo
    ) internal view {
        require(_stateChain.verifyStateInfo(_stateInfo), "incorrect state info");
        require(_stateChain.isStateConfirmed(_stateInfo), "provided state not confirmed");
        require(_stateInfo.index == _index, "should provide wanted state info");
    }
}
