// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "./RollupStateChain.sol";
import "./RollupInputChain.sol";
import "./ChainStorageContainer.sol";
import "../interfaces/IWhitelist.sol";
import "../libraries/console.sol";
import "./RollupStateChainSpecialCase.sol";

contract RollupStateChainSpecialUpgrade is Ownable {
    /**
     * when rollupstatechain upload error stateBatch ,
     * dao upgrade rollupstatechain to protect layer2 chain.
     */

    // must dao call this function
    function SpecialCaseUpgrade(
        address _dao,
        address _addressManager,
        address _proxyAdmin,
        uint64 n,
        address proposerError
    ) public onlyOwner {
        // get L1  contract address
        address dao = _dao;
        AddressManager addressManager = AddressManager(_addressManager);
        address rollupStateChain = address(addressManager.rollupStateChain());
        IWhitelist whitelist = addressManager.whitelist();
        ProxyAdmin proxyAdmin = ProxyAdmin(_proxyAdmin);
        address oldRollupStateChain = proxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(rollupStateChain))
        );

        // upgrade rollupStateChain --> rollupStateChainSpecialCase
        RollupStateChainSpecialCase newRollupStateChain = new RollupStateChainSpecialCase();
        proxyAdmin.upgrade(TransparentUpgradeableProxy(payable(rollupStateChain)), address(newRollupStateChain));
        newRollupStateChain = RollupStateChainSpecialCase(address(rollupStateChain));

        // dao ---> call newRollupStateChain ---> stateContainer ---> resize(n)
        newRollupStateChain.rollbackSpecialCase(n);
        // change errorProposer to blacklist
        whitelist.setProposer(proposerError, false);

        // change newRollupStateChain ---> upgrade to oldRollupStateChain
        proxyAdmin.upgrade(TransparentUpgradeableProxy(payable(rollupStateChain)), address(oldRollupStateChain));
        //set dao as addressManager owner & transfer ProxyAdminownership to dao
        //& transfer addressManagerOwnership to dao
        addressManager.setAddress(AddressName.DAO, dao);
        proxyAdmin.transferOwnership(dao);
        addressManager.transferOwnership(dao);
        require(proxyAdmin.owner() == dao, "proxyAdminaddr owner error");
        require(addressManager.dao() == dao, "dao address error");
        require(addressManager.owner() == dao, "addressManager owner error");
    }

    // if Upgrade function error , transfer ownership to dao
    function transferProxyAdminOwnership(address _dao, address _proxyAdmin) public onlyOwner {
        address dao = _dao;
        ProxyAdmin proxyAdmin = ProxyAdmin(_proxyAdmin);
        proxyAdmin.transferOwnership(dao);
    }

    function transferAddressManagerOwnership(address _dao, address _addressManager) public onlyOwner {
        address dao = _dao;
        AddressManager addressManager = AddressManager(_addressManager);
        addressManager.transferOwnership(dao);
    }

    function setAddress(address _dao, address _addressManager) public onlyOwner {
        address dao = _dao;
        AddressManager addressManager = AddressManager(_addressManager);
        addressManager.setAddress(AddressName.DAO, dao);
    }
}
