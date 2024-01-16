// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TestERC20 is ERC20 {
    uint8 private immutable decimal;

    constructor(string memory name_, string memory symbol_, uint8 decimals_) ERC20(name_, symbol_) {
        decimal = decimals_;
        _mint(msg.sender, 10000000 * (10 ** uint256(decimal)));
    }

    function decimals() public view override returns (uint8) {
        return decimal;
    }
}
