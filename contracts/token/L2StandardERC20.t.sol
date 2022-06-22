pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";
import "./L2StandardERC20.sol";

contract TestL2StandardERC20 is TestBase {
    L2StandardERC20 l2StandardERC20;
    address sender = address(0x898989);
    address mockL1Token = address(0x898988);
    address mockL2Bridge = address(0x898987);

    event Transfer(address indexed from, address indexed to, uint256 value);

    function setUp() public {
        vm.startPrank(sender);
        l2StandardERC20 = new L2StandardERC20(mockL2Bridge, mockL1Token, "test token", "tt");
        vm.stopPrank();
    }

    function testSupportsInterface() public {
        bytes4 interfaceId = bytes4(keccak256("supportsInterface(bytes4)"));
        bool res = l2StandardERC20.supportsInterface(interfaceId);
        require(res);
        interfaceId = bytes4(keccak256("test()"));
        res = l2StandardERC20.supportsInterface(interfaceId);
        require(!res);
    }

    function testMint() public {
        vm.startPrank(mockL2Bridge);
        address toAddr = address(0x2828);
        vm.expectEmit(true, true, true, true, address(l2StandardERC20));
        emit Transfer(address(0), toAddr, 999999999999 ether);
        l2StandardERC20.mint(toAddr, 999999999999 ether);
        uint256 toAddrBal = l2StandardERC20.balanceOf(toAddr);
        require(toAddrBal == 999999999999 ether);
        vm.stopPrank();
        vm.startPrank(sender);
        vm.expectRevert("Only L2 Bridge allowed");
        l2StandardERC20.mint(toAddr, 1 ether);
    }

    function testBurn() public {
        vm.startPrank(mockL2Bridge);
        address toAddr = address(0x2828);
        l2StandardERC20.mint(toAddr, 9999 ether);
        l2StandardERC20.burn(toAddr, 8888 ether);
        uint256 toAddrBal = l2StandardERC20.balanceOf(toAddr);
        require(toAddrBal == 1111 ether);
        vm.expectRevert("ERC20: burn amount exceeds balance");
        l2StandardERC20.burn(toAddr, 2222 ether);
        vm.stopPrank();
        vm.startPrank(sender);
        vm.expectRevert("Only L2 Bridge allowed");
        l2StandardERC20.burn(toAddr, 1 ether);
    }
}
