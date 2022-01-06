pragma solidity ^0.8.0;

interface IAddressResolver {
    //get stateCommitChain contract address.
    function scc() external view returns (address);

    function stakingManager() external view returns (address);

    function executor() external view returns (address);
}
