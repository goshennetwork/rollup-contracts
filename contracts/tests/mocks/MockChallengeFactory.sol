pragma solidity ^0.8.0;

import "../../interfaces/IChallengeFactory.sol";
import "../../interfaces/IStakingManager.sol";
import "../../interfaces/IStateTransition.sol";
import "../../interfaces/IStateCommitChain.sol";

//
contract MockChallengeFactory is IChallengeFactory {
    IStakingManager iStakingManager;
    IStateTransition iExecutor;
    IStateCommitChain iScc;

    //
    function init(
        address _sm,
        address _iexector,
        address _iscc
    ) external {
        require(address(iStakingManager == 0), "already init");
        iStakingManager = IStakingManager(_sm);
        iExecutor = IStateTransition(_iexector);
        iScc = IStateCommitChain(_iscc);
    }

    //
    function stakingManager() external view override returns (IStakingManager) {
        return iStakingManager;
    }

    function executor() external view override returns (IStateTransition) {
        return iExecutor;
    }

    function scc() external view override returns (IStateCommitChain) {
        return iScc;
    }

    function dao() external view returns (address) {
        return address(0xdead6666);
    }

    function isChallengeContract(address _addr) external view override returns (bool) {
        return true;
    }
}
