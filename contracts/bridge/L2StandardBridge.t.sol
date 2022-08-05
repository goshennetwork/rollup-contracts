// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../bridge/L2StandardBridge.sol";
import "../test-helper/TestBase.sol";
import "../token/L2StandardERC20.sol";

contract TestL2StandardBridge is TestBase, L2StandardBridge {
    L2StandardBridge l2StandardBridge;
    L2StandardERC20 testErc20;
    address mockL1Token = address(0x666666);
    address l1MockBridgeAddr = address(0x1111);
    address sender = address(0x88888);
    address toAddr = address(0x99999);

    function setUp() public {
        _initialize(sender);
        vm.startPrank(sender);
        L2StandardBridge l2StandardBridgeLogic = new L2StandardBridge();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(l2StandardBridgeLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(L2StandardBridge.initialize.selector, address(l2CrossLayerWitness), l1MockBridgeAddr)
        );
        l2StandardBridge = L2StandardBridge(payable(proxy));
        require(l2StandardBridge.l1TokenBridge() == l1MockBridgeAddr);
        testErc20 = new L2StandardERC20(address(l2StandardBridge), mockL1Token, "test token", "test");
        testErc20.approve(address(l2StandardBridge), 10 ether);
        vm.stopPrank();
    }

    function testWithdrawETH() public {
        vm.deal(sender, 10 ether);
        vm.startPrank(sender, sender);
        uint256 l2StandardBridgeBal = address(l2StandardBridge).balance;
        uint256 senderBal = sender.balance;
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit WithdrawalInitiated(address(0), address(0), sender, sender, 1 ether, "0x01");
        l2StandardBridge.withdrawETH{ value: 1 ether }("0x01");
        uint256 l2StandardBridgeAfterBal = address(l2StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal - senderAfterBal == 1 ether, "testWithdrawETH failed");
        require(l2StandardBridgeBal + 1 ether == l2StandardBridgeAfterBal, "testWithdrawETH failed");
    }

    function testWithdrawETHWithZeroValue() public {
        // test withdraw amount == 0
        vm.startPrank(sender, sender);
        uint256 l2StandardBridgeBal = address(l2StandardBridge).balance;
        uint256 senderBal = sender.balance;
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit WithdrawalInitiated(address(0), address(0), sender, sender, 0, "0x01");
        l2StandardBridge.withdrawETH("0x01");
        uint256 l2StandardBridgeAfterBal = address(l2StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal == senderAfterBal, "testWithdrawETH failed");
        require(l2StandardBridgeBal == l2StandardBridgeAfterBal, "testWithdrawETH failed");
    }

    function testFailWithdrawETH() public {
        // test amount > sender.balance
        vm.deal(sender, 10);
        vm.startPrank(sender, sender);
        l2StandardBridge.withdrawETH{ value: 20 }("0x01");
    }

    function testWithdrawETHTo() public {
        vm.deal(sender, 10);
        vm.startPrank(sender, sender);
        uint256 l2StandardBridgeBal = address(l2StandardBridge).balance;
        uint256 senderBal = sender.balance;
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit WithdrawalInitiated(address(0), address(0), sender, toAddr, 10, "0x01");
        l2StandardBridge.withdrawETHTo{ value: 10 }(toAddr, "0x01");
        uint256 l2StandardBridgeAfterBal = address(l2StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal - senderAfterBal == 10, "testWithdrawETHTo failed");
        require(l2StandardBridgeBal + 10 == l2StandardBridgeAfterBal, "testWithdrawETHTo failed");
    }

    function testWithdrawETHToWithZeroToAddr() public {
        // test toAddr = address(0)
        toAddr = address(0);
        vm.deal(sender, 10);
        vm.startPrank(sender, sender);
        uint256 senderBal = sender.balance;
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit WithdrawalInitiated(address(0), address(0), sender, toAddr, 10, "0x01");
        l2StandardBridge.withdrawETHTo{ value: 10 }(toAddr, "0x01");
        uint256 l2StandardBridgeBal = address(l2StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal - senderAfterBal == 10, "testWithdrawETHTo failed");
        require(l2StandardBridgeBal == 10, "testWithdrawETHTo failed");
    }

    function testWithdraw() public {
        vm.startPrank(address(l2StandardBridge));
        testErc20.mint(sender, 100 ether);
        vm.stopPrank();
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        uint256 l2StandardBridgeBal = testErc20.balanceOf(address(l2StandardBridge));
        require(l2StandardBridgeBal == 0, "testWithdraw failed");
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit WithdrawalInitiated(mockL1Token, address(testErc20), sender, sender, 1 ether, "0x01");
        l2StandardBridge.withdraw(address(testErc20), 1 ether, "0x01");
        uint256 senderAfterBal = testErc20.balanceOf(sender);
        uint256 l2StandardBridgeAfterBal = testErc20.balanceOf(address(l2StandardBridge));
        require(senderBal - 1 ether == senderAfterBal, "testWithdraw failed");
        require(l2StandardBridgeAfterBal == 0, "testWithdraw failed");
    }

    function testWithdrawWithZeroValue() public {
        // test amount == 0
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        uint256 l2StandardBridgeBal = testErc20.balanceOf(address(l2StandardBridge));
        require(l2StandardBridgeBal == 0, "testWithdraw failed");
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit WithdrawalInitiated(mockL1Token, address(testErc20), sender, sender, 0, "0x01");
        l2StandardBridge.withdraw(address(testErc20), 0, "0x01");
        uint256 senderAfterBal = testErc20.balanceOf(sender);
        uint256 l2StandardBridgeAfterBal = testErc20.balanceOf(address(l2StandardBridge));
        require(senderBal == senderAfterBal, "testWithdraw failed");
        require(l2StandardBridgeAfterBal == 0, "testWithdraw failed");
    }

    function testFailWithdraw() public {
        // test amount > sender.balance
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        l2StandardBridge.withdraw(address(testErc20), senderBal + 1, "0x01");
    }

    function testWithdrawTo() public {
        vm.startPrank(address(l2StandardBridge));
        testErc20.mint(sender, 100 ether);
        vm.stopPrank();
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        require(senderBal >= 10 ether);
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit WithdrawalInitiated(mockL1Token, address(testErc20), sender, toAddr, 10 ether, "0x01");
        l2StandardBridge.withdrawTo(address(testErc20), toAddr, 10 ether, "0x01");
        uint256 senderAfterBal = testErc20.balanceOf(sender);
        require(senderBal - 10 ether == senderAfterBal, "testWithdrawTo failed");
    }

    function testFailWithdrawTo() public {
        // amount > testErc20.balanceOf(sender)
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        require(senderBal == 0);
        l2StandardBridge.withdrawTo(address(testErc20), toAddr, 10 ether, "0x01");
    }

    function testFinalizeETHDeposit() public {
        vm.deal(address(l2StandardBridge), 10 ether);
        uint256 l2StandardBridgeBal = address(l2StandardBridge).balance;
        uint256 toAddrBal = toAddr.balance;
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeETHDeposit(address,address,uint256,bytes)",
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit DepositFinalized(address(0), address(0), sender, toAddr, 1 ether, "0x01");
        callRelayMessage(2, address(l2StandardBridge), l1MockBridgeAddr, signatureWithData);
        uint256 l2StandardBridgeAfterBal = address(l2StandardBridge).balance;
        uint256 toAddrAfterBal = toAddr.balance;
        require(toAddrBal + 1 ether == toAddrAfterBal, "testFinalizeETHDeposit failed");
        require(l2StandardBridgeBal - 1 ether == l2StandardBridgeAfterBal, "testFinalizeETHDeposit failed");
    }

    function testFailFinalizeETHDeposit() public {
        // l1StandardBridge ETH not enough   tx revert
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeETHDeposit(address,address,uint256,bytes)",
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        callRelayMessage(2, address(l2StandardBridge), l1MockBridgeAddr, signatureWithData);
    }

    function testFailFinalizeETHDepositWithL2BridgeETHNotEnough() public {
        vm.deal(address(l2StandardBridge), 1 ether);
        uint256 l2StandardBridgeBal = address(l2StandardBridge).balance;
        require(l2StandardBridgeBal == 1 ether, "l2StandardBridgeBal != 1 ether");
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeETHDeposit(address,address,uint256,bytes)",
            sender,
            toAddr,
            2 ether,
            "0x01"
        );
        callRelayMessage(2, address(l2StandardBridge), l1MockBridgeAddr, signatureWithData);
    }

    function testFinalizeERC20Deposit() public {
        vm.startPrank(address(l2StandardBridge));
        testErc20.mint(sender, 100 ether);
        vm.stopPrank();
        uint256 toAddrBal = testErc20.balanceOf(toAddr);
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeERC20Deposit(address,address,address,address,uint256,bytes)",
            mockL1Token,
            address(testErc20),
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit DepositFinalized(mockL1Token, address(testErc20), sender, toAddr, 1 ether, "0x01");
        callRelayMessage(2, address(l2StandardBridge), l1MockBridgeAddr, signatureWithData);
        uint256 toAddrAfterBal = testErc20.balanceOf(toAddr);
        require(toAddrBal + 1 ether == toAddrAfterBal, "testFinalizeERC20Deposit failed");
    }

    function testFailFinalizeERC20DepositWithZeroToAddr() public {
        toAddr = address(0);
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeERC20Deposit(address,address,address,address,uint256,bytes)",
            mockL1Token,
            address(testErc20),
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        callRelayMessage(2, address(l2StandardBridge), l1MockBridgeAddr, signatureWithData);
    }

    function testFinalizeERC20DepositWithWrongL1TokenAddr() public {
        mockL1Token = address(0x08989);
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeERC20Deposit(address,address,address,address,uint256,bytes)",
            mockL1Token,
            address(testErc20),
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        // check topic1, 2, 3, data && originating contract
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit DepositFailed(mockL1Token, address(testErc20), sender, toAddr, 1 ether, "0x01");
        callRelayMessage(2, address(l2StandardBridge), l1MockBridgeAddr, signatureWithData);
    }

    function testFinalizeERC20DepositWithWrongL2TokenAddr() public {
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeERC20Deposit(address,address,address,address,uint256,bytes)",
            mockL1Token,
            mockL1Token,
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        // check topic1, 2, 3, data && originating contract
        vm.expectEmit(true, true, true, true, address(l2StandardBridge));
        emit DepositFailed(mockL1Token, mockL1Token, sender, toAddr, 1 ether, "0x01");
        callRelayMessage(2, address(l2StandardBridge), l1MockBridgeAddr, signatureWithData);
    }
}
