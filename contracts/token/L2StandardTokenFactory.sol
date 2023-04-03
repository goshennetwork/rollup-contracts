// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./L2StandardERC20.sol";

contract L2StandardTokenFactory {
    address immutable l2StandardBridge;
    event StandardL2TokenCreated(address indexed _l1Token, address indexed _l2Token);

    constructor(address _l2StandardBridge) {
        l2StandardBridge = _l2StandardBridge;
    }

    /**
     * @dev Creates an instance of the standard ERC20 token on L2.
     * @param _l1Token Address of the corresponding L1 token.
     * @param _name ERC20 name.
     * @param _symbol ERC20 symbol.
     */
    function createStandardL2Token(
        address _l1Token,
        string memory _name,
        string memory _symbol
    ) external {
        require(_l1Token != address(0), "Must provide L1 token address");

        L2StandardERC20 l2Token = new L2StandardERC20(l2StandardBridge, _l1Token, _name, _symbol);

        emit StandardL2TokenCreated(_l1Token, address(l2Token));
    }
}
