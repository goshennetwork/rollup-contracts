pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/interfaces/IERC20.sol";

contract L2FeeCollector is Ownable {

    receive() external payable {}

    function withdrawEth() public onlyOwner {
        address _owner = owner();
        _owner.transfer(this.value);
    }

    // in case we can receive other tokens as fee
    function withdrawERC20(IERC20 token) public onlyOwner {
        address _owner = owner();
        uint balance = token.balanceOf(address(this));
        token.transfer(_owner, balance);
    }
}
