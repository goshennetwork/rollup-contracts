// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

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
    AddressManager addressManager;
    RollupStateChain rollupStateChain;
    RollupInputChain rollupInputChain;
    L1CrossLayerWitness l1CrossLayerWitness;
    L2CrossLayerWitness l2CrossLayerWitness;
    TestERC20 feeToken;
    StakingManager stakingManager;
    uint256 constant fraudProofWindow = 3;
    address challengerFactory;
    DAO dao;

    function initialize() internal {
        addressManager = new AddressManager();
        addressManager.initialize();
        l1CrossLayerWitness = new L1CrossLayerWitness();
        l1CrossLayerWitness.initialize(address(addressManager));
        l2CrossLayerWitness = new L2CrossLayerWitness();
        feeToken = new TestERC20("test token", "test");
        rollupStateChain = new RollupStateChain();
        rollupStateChain.initialize(address(addressManager), fraudProofWindow);
        challengerFactory = address(new MockChallengeFactory());
        stakingManager = new StakingManager();
        dao = new DAO();
        dao.initialize();
        stakingManager.initialize(
            address(dao),
            challengerFactory,
            address(rollupStateChain),
            address(feeToken),
            1 ether
        );
        rollupInputChain = new RollupInputChain();
        rollupInputChain.initialize(address(addressManager), 15000000, 3000000);
        ChainStorageContainer stateStorageContainer = new ChainStorageContainer();
        stateStorageContainer.initialize(AddressName.ROLLUP_STATE_CHAIN, address(addressManager));
        address stateStorage = address(stateStorageContainer);
        ChainStorageContainer inputStorageContainer = new ChainStorageContainer();
        inputStorageContainer.initialize(AddressName.ROLLUP_INPUT_CHAIN, address(addressManager));
        address inputStorage = address(inputStorageContainer);
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
}

contract MockChallengeFactory {
    function isChallengeContract(address _addr) external view returns (bool) {
        return _addr == address(this);
    }
}
