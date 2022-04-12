// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./IChainStorageContainer.sol";
import "./IStakingManager.sol";
import "./IChallengeFactory.sol";
import "./ICanonicalTransactionChain.sol";

interface IAddressResolver {
    function resolve(string memory _name) external view returns (address);

    function ctc() external view returns (ICanonicalTransactionChain);

    function ctcContainer() external view returns (IChainStorageContainer);

    function scc() external view returns (IStateCommitChain);

    function sccContainer() external view returns (IChainStorageContainer);

    function stakingManager() external view returns (IStakingManager);

    function challengeFactory() external view returns (IChallengeFactory);

    function crossDomainAddr() external view returns (address);
}
