// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface IAddressManager {
    event AddressSet(string _name, address _old, address _new);

    /**
     * @dev set new address related name
     * @param _name Contract name to related
     * @param _addr Contract address
     * @notice Revert when contract name has not set its address before
     */
    function setAddress(string memory _name, address _addr) external;
}
