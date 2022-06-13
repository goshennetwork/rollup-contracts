// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import "../resolver/AddressManager.sol";
import "../resolver/AddressName.sol";
import "../staking/StakingManager.sol";
import "./TestERC20.sol";
import "../rollup/RollupStateChain.sol";
import "../rollup/RollupInputChain.sol";
import "../rollup/ChainStorageContainer.sol";
import "../cross-layer/L1CrossLayerWitness.sol";
import "../cross-layer/L2CrossLayerWitness.sol";
import "../interfaces/ForgeVM.sol";
import "../libraries/Types.sol";
import "../dao/DAO.sol";

contract TestBase {
    ForgeVM public constant vm = ForgeVM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

    using MerkleMountainRange for CompactMerkleTree;
    CompactMerkleTree _trees;

    AddressManager addressManager;
    RollupStateChain rollupStateChain;
    RollupInputChain rollupInputChain;
    L1CrossLayerWitness l1CrossLayerWitness;
    L2CrossLayerWitness l2CrossLayerWitness;
    TestERC20 feeToken;
    StakingManager stakingManager;
    ProxyAdmin proxyAdmin;
    uint256 constant fraudProofWindow = 3;
    address challengerFactory;
    DAO dao;

    function initialize() internal {
        // deploy proxy admin

        proxyAdmin = new ProxyAdmin();
        // deploy AddressManager
        AddressManager addressManagerLogic = new AddressManager();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(addressManagerLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(AddressManager.initialize.selector)
        );
        addressManager = AddressManager(address(proxy));

        // deploy L1CrossLayerWitness
        L1CrossLayerWitness l1CrossLayerWitnessLogic = new L1CrossLayerWitness();
        proxy = new TransparentUpgradeableProxy(
            address(l1CrossLayerWitnessLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(L1CrossLayerWitness.initialize.selector, address(addressManager))
        );
        l1CrossLayerWitness = L1CrossLayerWitness(address(proxy));

        // deploy L2CrossLayerWitness
        L2CrossLayerWitness l2CrossLayerWitnessLogic = new L2CrossLayerWitness();
        proxy = new TransparentUpgradeableProxy(
            address(l2CrossLayerWitnessLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(L2CrossLayerWitness.initialize.selector, address(addressManager))
        );
        l2CrossLayerWitness = L2CrossLayerWitness(address(proxy));

        feeToken = new TestERC20("test token", "test");

        RollupStateChain rollupStateChainLogic = new RollupStateChain();
        proxy = new TransparentUpgradeableProxy(
            address(rollupStateChainLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(RollupStateChain.initialize.selector, address(addressManager), fraudProofWindow)
        );
        rollupStateChain = RollupStateChain(address(proxy));

        // TODO: use normal challenge factory
        challengerFactory = address(new MockChallengeFactory());

        // deploy dao
        DAO daoLogic = new DAO();
        proxy = new TransparentUpgradeableProxy(
            address(daoLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(DAO.initialize.selector)
        );
        dao = DAO(address(proxy));

        // deploy staking manager
        StakingManager stakingManagerLogic = new StakingManager();
        proxy = new TransparentUpgradeableProxy(
            address(stakingManagerLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(
                StakingManager.initialize.selector,
                address(dao),
                challengerFactory,
                address(rollupStateChain),
                address(feeToken),
                1 ether
            )
        );
        stakingManager = StakingManager(address(proxy));

        // deploy RollupInputChain
        RollupInputChain rollupInputChainLogic = new RollupInputChain();
        proxy = new TransparentUpgradeableProxy(
            address(rollupInputChainLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(
                RollupInputChain.initialize.selector,
                address(addressManager),
                15000000,
                3000000,
                1234
            )
        );
        rollupInputChain = RollupInputChain(address(proxy));

        // deploy ChainStorageContainer
        ChainStorageContainer stateStorageContainer = new ChainStorageContainer();
        proxy = new TransparentUpgradeableProxy(
            address(stateStorageContainer),
            address(proxyAdmin),
            abi.encodeWithSelector(
                ChainStorageContainer.initialize.selector,
                AddressName.ROLLUP_STATE_CHAIN,
                address(addressManager)
            )
        );
        address stateStorage = address(address(proxy));

        // deploy ChainStorageContainer
        ChainStorageContainer inputStorageContainer = new ChainStorageContainer();
        proxy = new TransparentUpgradeableProxy(
            address(inputStorageContainer),
            address(proxyAdmin),
            abi.encodeWithSelector(
                ChainStorageContainer.initialize.selector,
                AddressName.ROLLUP_INPUT_CHAIN,
                address(addressManager)
            )
        );
        address inputStorage = address(address(proxy));

        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN, address(rollupInputChain));
        addressManager.setAddress(AddressName.STAKING_MANAGER, address(stakingManager));
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN_CONTAINER, stateStorage);
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN_CONTAINER, inputStorage);
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN, address(rollupStateChain));
        addressManager.setAddress(AddressName.L1_CROSS_LAYER_WITNESS, address(l1CrossLayerWitness));
        addressManager.setAddress(AddressName.L2_CROSS_LAYER_WITNESS, address(l2CrossLayerWitness));
        addressManager.setAddress(AddressName.DAO, address(dao));
        addressManager.setAddress(AddressName.CHALLENGE_FACTORY, challengerFactory);
    }


    function init(address target, address sender, bytes memory signatureWithData) internal {
        bytes32 _hash = CrossLayerCodec.crossLayerMessageHash(
            target,
            sender,
            0,
            signatureWithData
        );
        MerkleMountainRange.appendLeafHash(_trees, _hash);
        bytes32[] memory _proof;
        bytes[] memory list = new bytes[](15);
        list[13] = abi.encodePacked(_trees.rootHash);
        list[14] = abi.encodePacked(_trees.treeSize);
        bytes[] memory encodedList = new bytes[](15);
        for (uint256 i = 0; i < list.length; i++) {
            encodedList[i] = RLPWriter.writeBytes(list[i]);
        }
        bytes memory rlpData = RLPWriter.writeList(encodedList);
        Types.StateInfo memory stateInfo;
        stateInfo.blockHash = keccak256(rlpData);
        vm.startPrank(address(rollupStateChain));
        addressManager.rollupStateChainContainer().append(Types.hash(stateInfo));
        vm.warp(3);
        vm.stopPrank();
        vm.startPrank(address(addressManager));
        bool success = l1CrossLayerWitness.relayMessage(
            target,
            sender,
            signatureWithData,
            0,
            rlpData,
            stateInfo,
            _proof
        );
        require(success, "call relayMessage failed");
    }

}

contract MockChallengeFactory {
    function isChallengeContract(address _addr) external view returns (bool) {
        return _addr == address(this);
    }
}
