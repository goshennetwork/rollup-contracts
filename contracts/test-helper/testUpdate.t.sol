// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.13;
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/ForgeVM.sol";
import "../bridge/L1StandardBridge.sol";
import "../bridge/L2StandardBridge.sol";
import "./TestBase.sol";
import "./MockContract.sol";
import "../cross-layer/CrossLayerContext.sol";
import "../state-machine/StateTransition.sol";

contract testUpdate is TestBase{

    function setUp() public {
        _initialize();
    }

    // 1.test Bridge
    function testUpgradeL1Bridge() public {
        // deploy L1Bridge & L2Bridge
        L1StandardBridge l1StandardBridge = new L1StandardBridge();
        L2StandardBridge l2StandardBridge = new L2StandardBridge();
        //deploy proxy & init (l1witness + l2TokenBridge)
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(l1StandardBridge),
            address(proxyAdmin),
            abi.encodeWithSelector(L1StandardBridge.initialize.selector, address(l1CrossLayerWitness), l2StandardBridge)
        );
        l1StandardBridge = L1StandardBridge(payable(proxy));

        // upgrade L1Bridge
        MockL1Bridge newL1Bridge = new MockL1Bridge();
        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newL1Bridge));
        vm.stopPrank();
        newL1Bridge = MockL1Bridge(address(proxy));
        // test call newL1Bridge
        require(newL1Bridge.return1() == 1 , "upgrade fail 1");
    }


    function testUpgradeL2Bridge() public {
        // deploy L1Bridge & L2Bridge
        L1StandardBridge l1StandardBridge = new L1StandardBridge();
        L2StandardBridge l2StandardBridge = new L2StandardBridge();
        //deploy proxy & init (l2witness + l1TokenBridge)
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(l2StandardBridge),
            address(proxyAdmin),
            abi.encodeWithSelector(L1StandardBridge.initialize.selector, address(l2CrossLayerWitness), l1StandardBridge)
        );
        l2StandardBridge = L2StandardBridge(payable(proxy));

        // upgrade L2Bridge
        MockL2Bridge newL2Bridge = new MockL2Bridge();
        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newL2Bridge));
        vm.stopPrank();
        newL2Bridge = MockL2Bridge(address(proxy));
        // test call newL2Bridge
        require(newL2Bridge.return2() == 2 , "upgrade fail 2");
    }

    //L2FeeCollector  can't  upgrade


    // 2.test cross-layer
    function testUpgradeL1CrossLayerWitness() public {
        // deploy L1CrossLayerWitness
        L1CrossLayerWitness l1CrossLayerWitness = new L1CrossLayerWitness();
        // deploy Transparent proxy & init ()
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(l1CrossLayerWitness),
            address(proxyAdmin),
            abi.encodeWithSelector(L1CrossLayerWitness.initialize.selector, address(addressManager))
        );
        l1CrossLayerWitness = L1CrossLayerWitness(payable(proxy));

        // upgrade L1CrossLayerWitness
        MockL1CrossLayerWitness newl1CrossLayerWitness = new MockL1CrossLayerWitness();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newl1CrossLayerWitness));
        vm.stopPrank();

        newl1CrossLayerWitness = MockL1CrossLayerWitness(address(proxy));
        // test call newl1CrossLayerWitness
        require(newl1CrossLayerWitness.return3() == 3 , "upgrade fail 3");
    }

    function testUpgradeL2CrossLayerWitness() public {
        // deploy L2CrossLayerWitness
        L2CrossLayerWitness l2CrossLayerWitness = new L2CrossLayerWitness();
        // deploy Transparent proxy & init ()
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(l2CrossLayerWitness),
            address(proxyAdmin),
            abi.encodeWithSelector(L2CrossLayerWitness.initialize.selector, address(addressManager))
        );
        l2CrossLayerWitness = L2CrossLayerWitness(payable(proxy));

        // upgrade L2CrossLayerWitness
        MockL2CrossLayerWitness newl2CrossLayerWitness = new MockL2CrossLayerWitness();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newl2CrossLayerWitness));
        vm.stopPrank();
        
        newl2CrossLayerWitness = MockL2CrossLayerWitness(address(proxy));
        // test call newl1CrossLayerWitness
        require(newl2CrossLayerWitness.return4() == 4 , "upgrade fail 4");
    }

    // 3.test DAO
    function testUpgradeDAO() public {
        // deploy DAO
        DAO dao = new DAO();
        // deploy Transparent proxy & init ()
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(dao),
            address(proxyAdmin),
            abi.encodeWithSelector(DAO.initialize.selector)
        );
        dao = DAO(address(proxy));

        // upgrade Dao
        MockDAO newDao = new MockDAO();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newDao));
        vm.stopPrank();
        
        newDao = MockDAO(address(proxy));
        // test call newDao
        require(newDao.return5() == 5 , "upgrade fail 5");
    }

    // 4.test AddressManager
    function testUpgradeAddressManager() public {
        // deploy AddressManager & proxy & init
        AddressManager addressManager = new AddressManager();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(addressManager),
            address(proxyAdmin),
            abi.encodeWithSelector(AddressManager.initialize.selector)
        );
        addressManager = AddressManager(address(proxy));

        // upgrade addressManager
        MockAddressManager newAddressManager = new MockAddressManager();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newAddressManager));
        vm.stopPrank();
        
        newAddressManager = MockAddressManager(address(proxy));
        // test call newDao
        require(newAddressManager.return6() == 6 , "upgrade fail 6");
    }

    // 5.test rollup
    function testUpgradeChainStorageContainer() public {
        // deploy ChainStorageContainer & proxy & init
        ChainStorageContainer chainStorageContainer = new ChainStorageContainer();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(chainStorageContainer),
            address(proxyAdmin),
            abi.encodeWithSelector(
                ChainStorageContainer.initialize.selector,
                AddressName.ROLLUP_STATE_CHAIN,
                address(addressManager)
            )
        );
        chainStorageContainer = ChainStorageContainer(address(proxy));

        // upgrade stateStorageContainer
        MockChainStorageContainer newChainStorageContainer = new MockChainStorageContainer();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newChainStorageContainer));
        vm.stopPrank();
        
        newChainStorageContainer = MockChainStorageContainer(address(proxy));
        // test call newDao
        require(newChainStorageContainer.return7() == 7 , "upgrade fail 7");
    }

    function testUpgradeRollupInputChain() public {
        // deploy RollupInputChain
        RollupInputChain rollupInputChain = new RollupInputChain();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(rollupInputChain),
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

        // upgrade rollupInputChain
        MockRollupInputChain newRollupInputChain = new MockRollupInputChain();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newRollupInputChain));
        vm.stopPrank();
        
        newRollupInputChain = MockRollupInputChain(address(proxy));
        // test call newRollupInputChain
        require(newRollupInputChain.return8() == 8 , "upgrade fail 8");
    }

    function testUpgradeRollupStateChain() public {
        // deploy RollupStateChain
        RollupStateChain rollupStateChain = new RollupStateChain();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(rollupStateChain),
            address(proxyAdmin),
            abi.encodeWithSelector(RollupStateChain.initialize.selector, address(addressManager), 3)
        );
        rollupStateChain = RollupStateChain(address(proxy));

        // upgrade RollupStateChain
        MockRollupStateChain newRollupStateChain = new MockRollupStateChain();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newRollupStateChain));
        vm.stopPrank();
        
        newRollupStateChain = MockRollupStateChain(address(proxy));
        // test call newRollupStateChain
        require(newRollupStateChain.return9() == 9 , "upgrade fail 9");
    }


    // 6.test staking
    function testUpgradeStakingManager() public {
        // deploy staking manager
        StakingManager stakingManager = new StakingManager();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(stakingManager),
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

        // upgrade StakingManager
        MockStakingManager newStakingManager = new MockStakingManager();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newStakingManager));
        vm.stopPrank();
        
        newStakingManager = MockStakingManager(address(proxy));
        // test call newStakingManager
        require(newStakingManager.return10() == 10 , "upgrade fail 10");
    }


    // 7.test state-machine
    function testUpgradeStateTransition() public {
        // deploy StateTransition & related contract
        StateTransition stateTransition = new StateTransition();
        MachineState machineState = new MachineState();

        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(stateTransition),
            address(proxyAdmin),
            abi.encodeWithSelector(
                stateTransition.initialize.selector,
                bytes32(0),
                IAddressResolver(addressManager),
                IMachineState(machineState)
            )
        );
        stateTransition = StateTransition(address(proxy));

        // upgrade StateTransition
        MockStateTransition newStateTransition = new MockStateTransition();

        vm.startPrank(ownerAddress);
        proxyAdmin.upgrade(proxy,address(newStateTransition));
        vm.stopPrank();
        
        newStateTransition = MockStateTransition(address(proxy));
        // test call newStateTransition
        require(newStateTransition.return11() == 11 , "upgrade fail 11");
    }

}