// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface IAddressManager {
    event AddressUpdated(string _name, address _old, address _new);

    /**
     * @dev Set new address related name
     * @param _name Contract name to related
     * @param _addr Contract address
     * @notice Revert when contract name already set its address
     */
    function newAddr(string memory _name, address _addr) external;

    /**
     * @dev update new address related name
     * @param _name Contract name to related
     * @param _addr Contract address
     * @notice Revert when contract name has not set its address before
     */
    function updateAddr(string memory _name, address _addr) external;
}
