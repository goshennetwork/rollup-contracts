// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./IChainStorageContainer.sol";
import "./IStakingManager.sol";
import "./IChallengeFactory.sol";
import "./ICanonicalTransactionChain.sol";

///@dev resolver only read address
interface IAddressResolver {
    ///Get address related with name
    ///@notice Revert if wanted contract have no address recorded
    function resolve(string memory _name) external view returns (address);

    ///Get CanonicalTransactionChain contract
    function ctc() external view returns (ICanonicalTransactionChain);

    ///Get ChainStorageContainer of CanonicalTransactionChain contract
    function ctcContainer() external view returns (IChainStorageContainer);

    ///Get StateCommitChain contract
    function scc() external view returns (IStateCommitChain);

    ///Get ChainStorageContainer of StateCommitChain contract
    function sccContainer() external view returns (IChainStorageContainer);

    ///Get StakingManager contract
    function stakingManager() external view returns (IStakingManager);

    ///Get ChallengeFactory contract
    function challengeFactory() external view returns (IChallengeFactory);

    ///get L1CrossDomain contract address
    function l1CrossDomainAddr() external view returns (address);
}
