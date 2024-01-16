// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../token/L2StandardERC20.sol";

contract L2TestERC20 is L2StandardERC20 {
    uint8 private immutable decimal;

    constructor(address _l2Bridge, address _l1Token, string memory _name, string memory _symbol, uint8 decimals_)
        L2StandardERC20(_l2Bridge, _l1Token, _name, _symbol)
    {
        decimal = decimals_;
    }

    function decimals() public view override returns (uint8) {
        return decimal;
    }
}
