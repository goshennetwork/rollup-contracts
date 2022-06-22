pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";

contract TestDAO is TestBase, DAO {
    address testAcc1 = address(0x1111);
    address testAcc2 = address(0x1112);
    address DAOOwner = address(0x999999);

    function setUp() public {
        vm.startPrank(DAOOwner);
        _initialize();

        vm.stopPrank();
    }

    function testSetSequencerWhitelist() public {
        vm.startPrank(DAOOwner);
        vm.expectEmit(true, true, true, true, address(dao));
        emit SequencerWhitelistUpdated(testAcc1, true);
        dao.setSequencerWhitelist(testAcc1, true);
        require(dao.sequencerWhitelist(testAcc1), "setSequencerWhitelist failed");
        vm.expectEmit(true, true, true, true, address(dao));
        emit SequencerWhitelistUpdated(testAcc2, true);
        dao.setSequencerWhitelist(testAcc2, true);
        require(dao.sequencerWhitelist(testAcc2), "setSequencerWhitelist failed");

        emit SequencerWhitelistUpdated(testAcc1, false);
        dao.setSequencerWhitelist(testAcc1, false);
        require(!dao.sequencerWhitelist(testAcc1), "setSequencerWhitelist failed");
        vm.expectEmit(true, true, true, true, address(dao));
        emit SequencerWhitelistUpdated(testAcc2, false);
        dao.setSequencerWhitelist(testAcc2, false);
        require(!dao.sequencerWhitelist(testAcc2), "setSequencerWhitelist failed");
    }

    function testSetSequencerWhitelistWithCallerIsNotOwner() public {
        vm.startPrank(testAcc1);
        vm.expectRevert("Ownable: caller is not the owner");
        dao.setSequencerWhitelist(testAcc1, true);
    }

    function testSetProposerWhitelist() public {
        vm.startPrank(DAOOwner);
        vm.expectEmit(true, true, true, true, address(dao));
        emit ProposerWhitelistUpdated(testAcc1, true);
        dao.setProposerWhitelist(testAcc1, true);
        require(dao.proposerWhitelist(testAcc1), "setProposerWhitelist failed");
        vm.expectEmit(true, true, true, true, address(dao));
        emit ProposerWhitelistUpdated(testAcc2, true);
        dao.setProposerWhitelist(testAcc2, true);
        require(dao.proposerWhitelist(testAcc2), "setProposerWhitelist failed");

        vm.expectEmit(true, true, true, true, address(dao));
        emit ProposerWhitelistUpdated(testAcc1, false);
        dao.setProposerWhitelist(testAcc1, false);
        require(!dao.proposerWhitelist(testAcc1), "setProposerWhitelist failed");
        vm.expectEmit(true, true, true, true, address(dao));
        emit ProposerWhitelistUpdated(testAcc2, false);
        dao.setProposerWhitelist(testAcc2, false);
        require(!dao.proposerWhitelist(testAcc2), "setProposerWhitelist failed");
    }

    function testSetProposerWhitelistWithCallerIsNotOwner() public {
        vm.startPrank(testAcc1);
        vm.expectRevert("Ownable: caller is not the owner");
        dao.setProposerWhitelist(testAcc1, true);
    }

    function testSetChallengerWhitelist() public {
        vm.startPrank(DAOOwner);
        vm.expectEmit(true, true, true, true, address(dao));
        emit ChallengerWhitelistUpdated(testAcc1, true);
        dao.setChallengerWhitelist(testAcc1, true);
        require(dao.challengerWhitelist(testAcc1), "setChallengerWhitelist failed");
        vm.expectEmit(true, true, true, true, address(dao));
        emit ChallengerWhitelistUpdated(testAcc2, true);
        dao.setChallengerWhitelist(testAcc2, true);
        require(dao.challengerWhitelist(testAcc2), "setChallengerWhitelist failed");

        vm.expectEmit(true, true, true, true, address(dao));
        emit ChallengerWhitelistUpdated(testAcc1, false);
        dao.setChallengerWhitelist(testAcc1, false);
        require(!dao.challengerWhitelist(testAcc1), "setChallengerWhitelist failed");
        vm.expectEmit(true, true, true, true, address(dao));
        emit ChallengerWhitelistUpdated(testAcc2, false);
        dao.setChallengerWhitelist(testAcc2, false);
        require(!dao.challengerWhitelist(testAcc2), "setChallengerWhitelist failed");
    }

    function testSetChallengerWhitelistWithCallerIsNotOwner() public {
        vm.startPrank(testAcc1);
        vm.expectRevert("Ownable: caller is not the owner");
        dao.setChallengerWhitelist(testAcc1, true);
    }

    function testTransferERC20() public {
        vm.startPrank(DAOOwner);
        feeToken.transfer(address(dao), 10 ether);
        uint256 daoBal = feeToken.balanceOf(address(dao));
        uint256 testAcc2Bal = feeToken.balanceOf(testAcc2);
        uint256 amount = 1 ether;
        dao.transferERC20(IERC20(feeToken), testAcc2, amount);
        uint256 daoAfterBal = feeToken.balanceOf(address(dao));
        uint256 testAcc2AfterBal = feeToken.balanceOf(testAcc2);
        require(daoBal - amount == daoAfterBal);
        require(testAcc2Bal + amount == testAcc2AfterBal);
    }

    function testTransferERC20WithCallerIsNotOwner() public {
        vm.startPrank(DAOOwner);
        feeToken.transfer(address(dao), 10 ether);
        vm.stopPrank();
        vm.startPrank(testAcc1);
        vm.expectRevert("Ownable: caller is not the owner");
        dao.transferERC20(IERC20(feeToken), testAcc2, 1 ether);
    }

    function testTransferERC20WithBalanceNotEnough() public {
        vm.startPrank(DAOOwner);
        feeToken.transfer(address(dao), 1 ether);
        vm.expectRevert("ERC20: transfer amount exceeds balance");
        dao.transferERC20(IERC20(feeToken), testAcc2, 10 ether);
    }
}
