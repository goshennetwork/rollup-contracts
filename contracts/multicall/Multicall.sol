// from https://github.com/Uniswap/v3-periphery/blob/79c708f357df69f7b3a494467e0f501810a11146/contracts/base/Multicall.sol
// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;
import "../interfaces/IMulticall.sol";

abstract contract Multicall is IMulticall {
    function multicall(bytes[] calldata data) public payable returns (bytes[] memory results) {
        results = new bytes[](data.length);
        for (uint256 i = 0; i < data.length; i++) {
            (bool success, bytes memory result) = address(this).delegatecall(data[i]);

            if (!success) {
                if (result.length < 68) revert();
                assembly {
                    result := add(result, 0x04)
                }
                revert(abi.decode(result, (string)));
            }
            results[i] = result;
        }
    }
}
