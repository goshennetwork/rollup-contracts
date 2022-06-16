// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../bridge/L1StandardBridge.sol";
import "../test-helper/TestBase.sol";
import "../test-helper/TestERC20.sol";

contract TestL1StandardBridge is TestBase, L1StandardBridge {
    L1StandardBridge l1StandardBridge;
    TestERC20 testErc20;
    address mockL2Token = address(0x666666);
    address l2MockBridgeAddr = address(0x1111);
    address sender = address(0x88888);
    address toAddr = address(0x99999);

    function setUp() public {
        _initialize();
        vm.startPrank(sender);
        L1StandardBridge l1StandardBridgeLogic = new L1StandardBridge();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(l1StandardBridgeLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(L1StandardBridge.initialize.selector, address(l1CrossLayerWitness), l2MockBridgeAddr)
        );
        l1StandardBridge = L1StandardBridge(payable(proxy));
        require(l1StandardBridge.l2TokenBridge() == l2MockBridgeAddr);
        testErc20 = new TestERC20("test token", "test");
        testErc20.approve(address(l1StandardBridge), 10 ether);
        vm.stopPrank();
    }

    function testDepositETH() public {
        vm.deal(sender, 10);
        vm.startPrank(sender, sender);
        uint256 senderBal = sender.balance;
        uint256 amount = 10;
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ETHDepositInitiated(sender, sender, amount, "0x01");
        l1StandardBridge.depositETH{ value: amount }("0x01");
        uint256 l1StandardBridgeBal = address(l1StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal - senderAfterBal == amount, "testDepositETH failed");
        require(l1StandardBridgeBal == amount, "testDepositETH failed");
    }

    function testDepositETHWithZeroValue() public {
        // test deposit amount == 0
        vm.startPrank(sender, sender);
        uint256 l1StandardBridgeBal = address(l1StandardBridge).balance;
        uint256 senderBal = sender.balance;
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ETHDepositInitiated(sender, sender, 0, "0x01");
        l1StandardBridge.depositETH("0x01");
        uint256 l1StandardBridgeAfterBal = address(l1StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal == senderAfterBal, "testDepositETH failed");
        require(l1StandardBridgeBal == l1StandardBridgeAfterBal, "testDepositETH failed");
    }

    function testFailDepositETH() public {
        // test amount > sender.balance
        vm.deal(sender, 10);
        vm.startPrank(sender, sender);
        l1StandardBridge.depositETH{ value: 20 }("0x01");
    }

    function testDepositETHTo() public {
        vm.deal(sender, 10);
        vm.startPrank(sender, sender);
        uint256 senderBal = sender.balance;
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ETHDepositInitiated(sender, toAddr, 10, "0x01");
        l1StandardBridge.depositETHTo{ value: 10 }(toAddr, "0x01");
        uint256 l1StandardBridgeBal = address(l1StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal - senderAfterBal == 10, "testDepositETH failed");
        require(l1StandardBridgeBal == 10, "testDepositETH failed");
    }

    function testDepositETHToWithZeroValue() public {
        // test toAddr = address(0)
        toAddr = address(0);
        vm.deal(sender, 10);
        vm.startPrank(sender, sender);
        uint256 senderBal = sender.balance;
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ETHDepositInitiated(sender, toAddr, 10, "0x01");
        l1StandardBridge.depositETHTo{ value: 10 }(toAddr, "0x01");
        uint256 l1StandardBridgeBal = address(l1StandardBridge).balance;
        uint256 senderAfterBal = sender.balance;
        require(senderBal - senderAfterBal == 10, "testDepositETH failed");
        require(l1StandardBridgeBal == 10, "testDepositETH failed");
    }

    function testDepositERC20() public {
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        uint256 l1StandardBridgeBal = testErc20.balanceOf(address(l1StandardBridge));
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ERC20DepositInitiated(address(testErc20), mockL2Token, sender, sender, 1 ether, "0x01");
        l1StandardBridge.depositERC20(address(testErc20), mockL2Token, 1 ether, "0x01");
        uint256 senderAfterBal = testErc20.balanceOf(sender);
        uint256 l1StandardBridgeAfterBal = testErc20.balanceOf(address(l1StandardBridge));
        require(senderBal - 1 ether == senderAfterBal, "DepositERC20 failed");
        require(l1StandardBridgeAfterBal - 1 ether == l1StandardBridgeBal, "DepositERC20 failed");
    }

    function testDepositERC20WithZeroValue() public {
        // test amount == 0
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        uint256 l1StandardBridgeBal = testErc20.balanceOf(address(l1StandardBridge));
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ERC20DepositInitiated(address(testErc20), mockL2Token, sender, sender, 0, "0x01");
        l1StandardBridge.depositERC20(address(testErc20), mockL2Token, 0, "0x01");
        uint256 senderAfterBal = testErc20.balanceOf(sender);
        uint256 l1StandardBridgeAfterBal = testErc20.balanceOf(address(l1StandardBridge));
        require(senderBal == senderAfterBal, "DepositERC20 failed");
        require(l1StandardBridgeAfterBal == l1StandardBridgeBal, "DepositERC20 failed");
    }

    function testFailDepositERC20() public {
        // test amount > sender.balance
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        l1StandardBridge.depositERC20(address(testErc20), mockL2Token, senderBal + 1, "0x01");
    }

    function testDepositERC20To() public {
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        uint256 l1StandardBridgeBal = testErc20.balanceOf(address(l1StandardBridge));
        // amount == 0
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ERC20DepositInitiated(address(testErc20), mockL2Token, sender, toAddr, 0, "0x01");
        l1StandardBridge.depositERC20To(address(testErc20), mockL2Token, toAddr, 0 ether, "0x01");
        // amount == 0.5 ether
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ERC20DepositInitiated(address(testErc20), mockL2Token, sender, toAddr, 0.5 ether, "0x01");
        l1StandardBridge.depositERC20To(address(testErc20), mockL2Token, toAddr, 0.5 ether, "0x01");
        uint256 senderAfterBal = testErc20.balanceOf(sender);
        uint256 l1StandardBridgeAfterBal = testErc20.balanceOf(address(l1StandardBridge));
        require(senderBal - 0.5 ether == senderAfterBal, "DepositERC20To failed1");
        require(l1StandardBridgeAfterBal - 0.5 ether == l1StandardBridgeBal, "DepositERC20To failed2");
        uint256 deposit = l1StandardBridge.deposits(address(testErc20), mockL2Token);
        require(deposit == 0.5 ether, "DepositERC20To failed");
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ERC20DepositInitiated(address(testErc20), mockL2Token, sender, toAddr, 1.5 ether, "0x01");
        l1StandardBridge.depositERC20To(address(testErc20), mockL2Token, toAddr, 1.5 ether, "0x01");
        uint256 senderAfterBal2 = testErc20.balanceOf(sender);
        uint256 l1StandardBridgeAfterBal2 = testErc20.balanceOf(address(l1StandardBridge));
        require(senderAfterBal - 1.5 ether == senderAfterBal2, "DepositERC20To failed3");
        require(l1StandardBridgeAfterBal2 - 1.5 ether == l1StandardBridgeAfterBal, "DepositERC20To failed4");
        deposit = l1StandardBridge.deposits(address(testErc20), mockL2Token);
        require(deposit == 2 ether, "DepositERC20To failed5");
    }

    function testFailDepositERC20To() public {
        // amount > testErc20.balanceOf(sender)
        vm.startPrank(sender, sender);
        uint256 senderBal = testErc20.balanceOf(sender);
        l1StandardBridge.depositERC20To(address(testErc20), mockL2Token, toAddr, senderBal + 1, "0x01");
    }

    function testFinalizeETHWithdrawal() public {
        vm.deal(address(l1StandardBridge), 10 ether);
        uint256 l1StandardBridgeBal = address(l1StandardBridge).balance;
        uint256 toAddrBal = toAddr.balance;
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeETHWithdrawal(address,address,uint256,bytes)",
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ETHWithdrawalFinalized(sender, toAddr, 1 ether, "0x01");
        callRelayMessage(1, address(l1StandardBridge), l2MockBridgeAddr, signatureWithData);
        uint256 l1StandardBridgeAfterBal = address(l1StandardBridge).balance;
        uint256 toAddrAfterBal = toAddr.balance;
        require(toAddrBal + 1 ether == toAddrAfterBal, "testFinalizeETHWithdrawal failed");
        require(l1StandardBridgeBal - 1 ether == l1StandardBridgeAfterBal, "testFinalizeETHWithdrawal failed");
    }

    function testFailFinalizeETHWithdrawalWithL1BridgeETHNotEnough() public {
        // l1StandardBridge ETH not enough   tx revert
        uint256 l1StandardBridgeBal = address(l1StandardBridge).balance;
        require(l1StandardBridgeBal < 1 ether, "l1StandardBridgeBal failed");
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeETHWithdrawal(address,address,uint256,bytes)",
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        callRelayMessage(1, address(l1StandardBridge), l2MockBridgeAddr, signatureWithData);
    }

    function testFinalizeERC20Withdrawal() public {
        vm.startPrank(sender);
        testErc20.transfer(address(l1StandardBridge), 10 ether);
        vm.stopPrank();
        vm.startPrank(sender, sender);
        l1StandardBridge.depositERC20To(address(testErc20), mockL2Token, toAddr, 1 ether, "0x01");
        uint256 deposit = l1StandardBridge.deposits(address(testErc20), mockL2Token);
        require(deposit >= 1 ether, "DepositERC20To failed");
        vm.stopPrank();
        uint256 l1StandardBridgeBal = testErc20.balanceOf(address(l1StandardBridge));
        uint256 toAddrBal = testErc20.balanceOf(toAddr);
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeERC20Withdrawal(address,address,address,address,uint256,bytes)",
            address(testErc20),
            mockL2Token,
            sender,
            toAddr,
            1 ether,
            "0x01"
        );
        vm.expectEmit(true, true, true, true, address(l1StandardBridge));
        emit ERC20WithdrawalFinalized(address(testErc20), mockL2Token, sender, toAddr, 1 ether, "0x01");
        callRelayMessage(1, address(l1StandardBridge), l2MockBridgeAddr, signatureWithData);
        uint256 l1StandardBridgeAfterBal = testErc20.balanceOf(address(l1StandardBridge));
        uint256 toAddrAfterBal = testErc20.balanceOf(toAddr);
        require(toAddrBal + 1 ether == toAddrAfterBal, "testFinalizeERC20Withdrawal failed");
        require(l1StandardBridgeBal - 1 ether == l1StandardBridgeAfterBal, "testFinalizeERC20Withdrawal failed");
    }

    function testFailFinalizeERC20WithdrawalWithL1BridgeERC20TokenNotEnough() public {
        // test l1StandardBridge testErc20 Token not enough
        vm.startPrank(sender, sender);
        l1StandardBridge.depositERC20To(address(testErc20), mockL2Token, toAddr, 10 ether, "0x01");
        vm.stopPrank();
        vm.startPrank(address(l1StandardBridge), address(l1StandardBridge));
        testErc20.transfer(address(sender), 9 ether);
        uint256 l1StandardBridgeBal = testErc20.balanceOf(address(l1StandardBridge));
        require(l1StandardBridgeBal <= 1 ether, "l1StandardBridgeBal > amount");
        uint256 deposit = l1StandardBridge.deposits(address(testErc20), mockL2Token);
        require(deposit >= 10 ether, "DepositERC20To failed");
        vm.stopPrank();
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeERC20Withdrawal(address,address,address,address,uint256,bytes)",
            address(testErc20),
            mockL2Token,
            sender,
            toAddr,
            1.5 ether,
            "0x01"
        );
        callRelayMessage(1, address(l1StandardBridge), l2MockBridgeAddr, signatureWithData);
    }

    function testFailFinalizeERC20WithdrawalWithDepositNotEnough() public {
        // test l1StandardBridge.deposits(address(testErc20), mockL2Token) < amount
        vm.startPrank(sender, sender);
        l1StandardBridge.depositERC20To(address(testErc20), mockL2Token, toAddr, 1 ether, "0x01");
        testErc20.transfer(address(l1StandardBridge), 9 ether);
        uint256 l1StandardBridgeBal = testErc20.balanceOf(address(l1StandardBridge));
        require(l1StandardBridgeBal >= 1.5 ether, "l1StandardBridgeBal < amount");
        uint256 deposit = l1StandardBridge.deposits(address(testErc20), mockL2Token);
        require(deposit >= 1 ether, "DepositERC20To failed");
        vm.stopPrank();
        bytes memory signatureWithData = abi.encodeWithSignature(
            "finalizeERC20Withdrawal(address,address,address,address,uint256,bytes)",
            address(testErc20),
            mockL2Token,
            sender,
            toAddr,
            1.5 ether,
            "0x01"
        );
        callRelayMessage(1, address(l1StandardBridge), l2MockBridgeAddr, signatureWithData);
    }
}
