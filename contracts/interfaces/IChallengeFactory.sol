// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IStakingManager.sol";
import "../interfaces/IStateTransition.sol";
import "../interfaces/IRollupStateChain.sol";

interface IChallengeFactory {
    function stakingManager() external view returns (IStakingManager);

    function executor() external view returns (IStateTransition);

    function rollupStateChain() external view returns (IRollupStateChain);

    function dao() external view returns (address);

    function isChallengeContract(address _addr) external view returns (bool);
}
