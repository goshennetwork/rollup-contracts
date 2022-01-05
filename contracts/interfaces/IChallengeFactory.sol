// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IStakingManager.sol";
import "../interfaces/IExecutor.sol";
import "../interfaces/IStateCommitChain.sol";

interface IChallengeFactory {
    function stakingManager() external view returns (IStakingManager);

    function executor() external view returns (IExecutor);

    function scc() external view returns (IStateCommitChain);

    function isChallengeContract(address _addr) external view returns (bool);
}
