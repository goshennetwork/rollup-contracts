// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./IChainStorageContainer.sol";
import "./IStakingManager.sol";
import "./IChallengeFactory.sol";
import "./IRollupInputChain.sol";
import "./IL1CrossLayerMessageWitness.sol";
import "./IL2CrossLayerMessageWitness.sol";

///@dev resolver only read address
interface IAddressResolver {
    ///Get address related with name
    ///@notice Revert if wanted contract have no address recorded
    function resolve(string memory _name) external view returns (address);

    function dao() external view returns (address);

    ///Get RollupInputChain contract
    function rollupInputChain() external view returns (IRollupInputChain);

    ///Get ChainStorageContainer of RollupInputChain contract
    function rollupInputChainContainer() external view returns (IChainStorageContainer);

    ///Get RollupStateChain contract
    function rollupStateChain() external view returns (IRollupStateChain);

    ///Get ChainStorageContainer of RollupStateChain contract
    function rollupStateChainContainer() external view returns (IChainStorageContainer);

    ///Get StakingManager contract
    function stakingManager() external view returns (IStakingManager);

    ///Get ChallengeFactory contract
    function challengeFactory() external view returns (IChallengeFactory);

    ///get L1CrossLayerMessageWitness contract address
    function l1CrossLayerMessageWitness() external view returns (IL1CrossLayerMessageWitness);

    ///get L2CrossLayerMessageWitness contract address
    function l2CrossLayerMessageWitness() external view returns (IL2CrossLayerMessageWitness);
}
