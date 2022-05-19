// SPDX-License-Identifier: MIT
// Copyright 2020-2021 Optimism
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/introspection/IERC165.sol";

interface IL2StandardERC20 is IERC20, IERC165 {
    function l1Token() external returns (address);

    function mint(address _to, uint256 _amount) external;

    function burn(address _from, uint256 _amount) external;
}
