pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract L2FeeCollector is Ownable {
    receive() external payable {}

    function withdrawEth(uint256 _amount) public onlyOwner {
        Address.sendValue(payable(owner()), _amount);
    }

    function withdrawEthTo(address payable _to, uint256 _amount) public onlyOwner {
        Address.sendValue(_to, _amount);
    }

    // in case we can receive other tokens as fee
    function withdrawERC20(IERC20 token, uint256 _amount) public onlyOwner {
        SafeERC20.safeTransfer(token, owner(), _amount);
    }

    // in case we can receive other tokens as fee
    function withdrawERC20To(IERC20 token, address _to, uint256 _amount) public onlyOwner {
        SafeERC20.safeTransfer(token, _to, _amount);
    }
}
