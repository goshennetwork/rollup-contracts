pragma solidity ^0.8.0;

import "../../interfaces/IChallengeFactory.sol";
import "../../interfaces/IStakingManager.sol";
import "../../interfaces/IStateTransition.sol";
import "../../interfaces/IStateCommitChain.sol";
import "../../interfaces/IChallenge.sol";
import "../../challenge/Challenge.sol";

contract MockChallengeFactory is IChallengeFactory {
    IStakingManager public override stakingManager;
    IStateTransition public override executor;
    IStateCommitChain public override scc;

    function init(
        address _sm,
        address _executor,
        address _scc
    ) external {
        require(address(stakingManager) == address(0), "already init");
        stakingManager = IStakingManager(_sm);
        executor = IStateTransition(_executor);
        scc = IStateCommitChain(_scc);
    }

    function dao() external view override returns (address) {
        return address(0xdead6666);
    }

    function isChallengeContract(address _addr) external view override returns (bool) {
        return true;
    }

    function newChallengeWithProposer(address _creator, address _proposer) external returns (IChallenge) {
        bytes32 _ff = bytes32(uint256(0xff));
        IChallenge _c = new Challenge();
        stakingManager.token().transferFrom(_creator, address(_c), _c.minChallengerDeposit());
        _c.create(1, _proposer, _ff, _ff, _creator, 50);
        return _c;
    }
}
