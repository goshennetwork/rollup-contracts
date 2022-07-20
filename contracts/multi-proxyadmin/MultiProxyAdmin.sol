// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "../multicall/Multicall.sol";

contract MultiProxyAdmin is ProxyAdmin, Multicall {}
