// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IStakingManager.sol";
import "../interfaces/IExecutor.sol";
import "../interfaces/IStateCommitChain.sol";

interface IChallengeFactory {
    function stakingManager() public view returns (IStakingManager);

    function executor() public view returns (IExecutor);

    function scc() public view returns (IStateCommitChain);
}
