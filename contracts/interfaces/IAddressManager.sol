// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

interface IAddressManager {
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

    /**
     * @dev Get contract address by contract name,If contract name not saved yet, it will return empty address
     * @param _name Contract name to related
     * @return Contract address related by name, if name not exist, it will return empty address
     */
    function getAddr(string memory _name) external view returns (address);
}
