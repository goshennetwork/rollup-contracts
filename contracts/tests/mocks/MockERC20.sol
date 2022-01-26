pragma solidity ^0.8.0;
import "@openzeppelin/token/ERC20/ERC20.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("mock", "test") {
        _mint(msg.sender, 0x100_000_000_000_000_000);
    }
}
