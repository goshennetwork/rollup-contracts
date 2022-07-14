pragma solidity ^0.8.0;

import "./L2FeeCollector.sol";
import "../interfaces/ForgeVM.sol";
import "../test-helper/TestERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract TestL2FeeCollector {
    ForgeVM public constant vm = ForgeVM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    L2FeeCollector l2FeeCollector;
    TestERC20 testErc20;
    address sender = address(0x7878);
    address toAddr = address(0x8787);

    function setUp() public {
        vm.startPrank(sender);
        l2FeeCollector = new L2FeeCollector();
        testErc20 = new TestERC20("test token", "test");
        vm.stopPrank();
    }

    function testWithdrawEth() public {
        vm.deal(address(l2FeeCollector), 10 ether);
        vm.startPrank(sender);
        uint256 senderBal = sender.balance;
        uint256 l2FeeCollectorBal = address(l2FeeCollector).balance;
        uint256 amount = 1 ether;
        l2FeeCollector.withdrawEth(amount);
        uint256 senderAfterBal = sender.balance;
        uint256 l2FeeCollectorAfterBal = address(l2FeeCollector).balance;
        require(senderBal + amount == senderAfterBal);
        require(l2FeeCollectorBal - amount == l2FeeCollectorAfterBal);
    }

    function testWithdrawEthWithCallerNotOwner() public {
        vm.deal(address(l2FeeCollector), 10 ether);
        vm.startPrank(toAddr);
        uint256 amount = 1 ether;
        vm.expectRevert(bytes("Ownable: caller is not the owner"));
        l2FeeCollector.withdrawEth(amount);
    }

    function testWithdrawEthWithL2FeeETHNotEnough() public {
        vm.deal(address(l2FeeCollector), 1 ether);
        vm.startPrank(sender);
        uint256 amount = 2 ether;
        vm.expectRevert(bytes("Address: insufficient balance"));
        l2FeeCollector.withdrawEth(amount);
    }

    function testWithdrawEthTo() public {
        vm.deal(address(l2FeeCollector), 10 ether);
        vm.startPrank(sender);
        uint256 toAddrBal = toAddr.balance;
        uint256 l2FeeCollectorBal = address(l2FeeCollector).balance;
        uint256 amount = 1 ether;
        l2FeeCollector.withdrawEthTo(payable(toAddr), amount);
        uint256 toAddrAfterBal = toAddr.balance;
        uint256 l2FeeCollectorAfterBal = address(l2FeeCollector).balance;
        require(toAddrBal + amount == toAddrAfterBal);
        require(l2FeeCollectorBal - amount == l2FeeCollectorAfterBal);
    }

    function testWithdrawEthToWhitZeroToAddr() public {
        toAddr = address(0);
        vm.deal(address(l2FeeCollector), 10 ether);
        vm.startPrank(sender);
        uint256 l2FeeCollectorBal = address(l2FeeCollector).balance;
        uint256 amount = 1 ether;
        l2FeeCollector.withdrawEthTo(payable(toAddr), amount);
        uint256 l2FeeCollectorAfterBal = address(l2FeeCollector).balance;
        require(l2FeeCollectorBal - amount == l2FeeCollectorAfterBal);
    }

    function testWithdrawEthToWithCallerNotOwner() public {
        vm.deal(address(l2FeeCollector), 10 ether);
        vm.startPrank(toAddr);
        uint256 amount = 1 ether;
        vm.expectRevert(bytes("Ownable: caller is not the owner"));
        l2FeeCollector.withdrawEthTo(payable(toAddr), amount);
    }

    function testWithdrawERC20() public {
        vm.startPrank(sender);
        testErc20.transfer(address(l2FeeCollector), 10 ether);
        uint256 amount = 1 ether;
        uint256 senderBal = testErc20.balanceOf(sender);
        uint256 l2FeeCollectorBal = testErc20.balanceOf(address(l2FeeCollector));
        l2FeeCollector.withdrawERC20(IERC20(testErc20), amount);
        uint256 senderAfterBal = testErc20.balanceOf(sender);
        uint256 l2FeeCollectorAfterBal = testErc20.balanceOf(address(l2FeeCollector));
        require(senderBal + amount == senderAfterBal);
        require(l2FeeCollectorBal - amount == l2FeeCollectorAfterBal);
    }

    function testWithdrawERC20WithWrongTokenAddr() public {
        vm.startPrank(sender);
        testErc20.transfer(address(l2FeeCollector), 10 ether);
        uint256 amount = 1 ether;
        vm.expectRevert(bytes("Address: call to non-contract"));
        l2FeeCollector.withdrawERC20(IERC20(toAddr), amount);
    }

    function testWithdrawERC20WithL2FeeETHNotEnough() public {
        vm.startPrank(sender);
        uint256 amount = 2 ether;
        vm.expectRevert(bytes("ERC20: transfer amount exceeds balance"));
        l2FeeCollector.withdrawERC20(IERC20(testErc20), amount);
    }

    function testWithdrawERC20WithCallerNotOwner() public {
        vm.startPrank(sender);
        testErc20.transfer(address(l2FeeCollector), 10 ether);
        vm.stopPrank();
        vm.startPrank(toAddr);
        uint256 amount = 1 ether;
        vm.expectRevert(bytes("Ownable: caller is not the owner"));
        l2FeeCollector.withdrawERC20(IERC20(testErc20), amount);
    }

    function testWithdrawERC20To() public {
        vm.startPrank(sender);
        testErc20.transfer(address(l2FeeCollector), 10 ether);
        uint256 amount = 1 ether;
        uint256 toAddrBal = testErc20.balanceOf(toAddr);
        uint256 l2FeeCollectorBal = testErc20.balanceOf(address(l2FeeCollector));
        l2FeeCollector.withdrawERC20To(IERC20(testErc20), toAddr, amount);
        uint256 toAddrAfterBal = testErc20.balanceOf(toAddr);
        uint256 l2FeeCollectorAfterBal = testErc20.balanceOf(address(l2FeeCollector));
        require(toAddrBal + amount == toAddrAfterBal);
        require(l2FeeCollectorBal - amount == l2FeeCollectorAfterBal);
    }

    function testWithdrawERC20ToWithCallerNotOwner() public {
        vm.startPrank(sender);
        testErc20.transfer(address(l2FeeCollector), 10 ether);
        vm.stopPrank();
        vm.startPrank(toAddr);
        uint256 amount = 1 ether;
        vm.expectRevert(bytes("Ownable: caller is not the owner"));
        l2FeeCollector.withdrawERC20To(IERC20(testErc20), toAddr, amount);
    }

    function testWithdrawERC20ToWhitZeroToAddr() public {
        toAddr = address(0);
        vm.startPrank(sender);
        testErc20.transfer(address(l2FeeCollector), 10 ether);
        uint256 amount = 1 ether;
        vm.expectRevert(bytes("ERC20: transfer to the zero address"));
        l2FeeCollector.withdrawERC20To(IERC20(testErc20), toAddr, amount);
    }
}
