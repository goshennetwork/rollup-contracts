// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/IRollupStateChain.sol";
import "../interfaces/IAddressResolver.sol";
import "../interfaces/IChainStorageContainer.sol";

contract RollupStateChainSpecialCase is Initializable {
    IAddressResolver addressResolver;
    //the window to fraud proof
    uint256 public fraudProofWindow;


    /**
     * when special case happend , dao will try to make sure system safe
     */
    event SpecialCaseRollbacked(uint64 indexed _stateIndex);

    function rollbackSpecialCase(uint64 size) public {
        require(msg.sender == address(addressResolver.dao()), "only dao allowed2");
        addressResolver.rollupStateChainContainer().resize(size);
        emit SpecialCaseRollbacked(size);
    }
}