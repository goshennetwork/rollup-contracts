// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./RollupStateChain.sol";
import "./RollupInputChain.sol";
import "./ChainStorageContainer.sol";
import "../test-helper/TestBase.sol";
import "./RollupStateChainSpecialCase.sol";
import "./RollupStateChainSpecialUpgrade.sol";

contract TestRollupStateChainSpecialUpgrade is TestBase {

    function setUp() public {}

    /**
     * test rollupStateChainSpecialCase upgrade
     * must provide (proxyAdminaddr & addressManager & dao)
     */
    function testSpecialCaseUpgradePass() public {
        // get L1 contract
        address proxyAdminaddr = 0x5A0be95863ad38fC9De4cE121DA009321D9fBCF2 ;
        AddressManager addressManager = AddressManager(0x59bbFDD6f1DAAE70CEdEf6FBBa4f623353ed7f93);
        address dao = 0x5e3f6E5E8f2F8cB02f087aA573FadB09867fB09E;

        console.log("1.--------------");
        console.log("proxyAdminaddr owner:");
        console.log(ProxyAdmin(proxyAdminaddr).owner());
        console.log("dao addr:");
        console.log(addressManager.dao());
        console.log("addressManager owner:");
        console.log(addressManager.owner());

        vm.startPrank(dao);
        // deploy new rollupStateChain
        RollupStateChainSpecialUpgrade rollupStateChainSpecialUpgrade = new RollupStateChainSpecialUpgrade();
        // set new rollupStateChainSpecialUpgrade as dao 
        // & transfer ownership of (proxyAdminaddr & addressManager) to new contract
        addressManager.setAddress(AddressName.DAO, address(rollupStateChainSpecialUpgrade));
        ProxyAdmin(proxyAdminaddr).transferOwnership(address(rollupStateChainSpecialUpgrade));
        addressManager.transferOwnership(address(rollupStateChainSpecialUpgrade));

        console.log("2.--------------");
        console.log("proxyAdminaddr owner:");
        console.log(ProxyAdmin(proxyAdminaddr).owner());
        console.log("dao addr:");
        console.log(addressManager.dao());
        console.log("addressManager owner:");
        console.log(addressManager.owner());
        //upgrade rollupStateChain & set proposer blacklist & rollback 
        //& upgarde oldstatechain &transfer ownership & change dao address
        rollupStateChainSpecialUpgrade.SpecialCaseUpgrade(dao ,address(addressManager), proxyAdminaddr,2, address(1));
        console.log("upgrade over");
        require(ProxyAdmin(proxyAdminaddr).owner() == dao, "proxyAdminaddr owner error");
        require(addressManager.dao() == dao, "dao address error");
        require(addressManager.owner() == dao, "addressManager owner error");

        console.log("3.--------------");
        console.log("proxyAdminaddr owner:");
        console.log(ProxyAdmin(proxyAdminaddr).owner());
        console.log("dao addr:");
        console.log(addressManager.dao());
        console.log("addressManager owner:");
        console.log(addressManager.owner());

        vm.stopPrank();
    }

    // if special case , must ensure that ownership can be transferred out ;
    function testTransferOwnership() public {
        address proxyAdminaddr = 0x5A0be95863ad38fC9De4cE121DA009321D9fBCF2 ;
        AddressManager addressManager = AddressManager(0x59bbFDD6f1DAAE70CEdEf6FBBa4f623353ed7f93);
        address dao = 0x5e3f6E5E8f2F8cB02f087aA573FadB09867fB09E;

        vm.startPrank(dao);
        // deploy new rollupStateChain
        RollupStateChainSpecialUpgrade rollupStateChainSpecialUpgrade = new RollupStateChainSpecialUpgrade();

        //1.transfer ownership to contract
        addressManager.setAddress(AddressName.DAO, address(rollupStateChainSpecialUpgrade));
        ProxyAdmin(proxyAdminaddr).transferOwnership(address(rollupStateChainSpecialUpgrade));
        addressManager.transferOwnership(address(rollupStateChainSpecialUpgrade));

        console.log("1.--------------");
        console.log("proxyAdminaddr owner:");
        console.log(ProxyAdmin(proxyAdminaddr).owner());
        console.log("dao addr:");
        console.log(addressManager.dao());
        console.log("addressManager owner:");
        console.log(addressManager.owner());

        //2.call newContract to transfer ownership back
        rollupStateChainSpecialUpgrade.setAddress(dao, address(addressManager));
        rollupStateChainSpecialUpgrade.transferProxyAdminOwnership(dao, proxyAdminaddr);
        rollupStateChainSpecialUpgrade.transferAddressManagerOwnership(dao, address(addressManager));

        console.log("2.--------------");
        console.log("proxyAdminaddr owner:");
        console.log(ProxyAdmin(proxyAdminaddr).owner());
        console.log("dao addr:");
        console.log(addressManager.dao());
        console.log("addressManager owner:");
        console.log(addressManager.owner());
        
    }

}
