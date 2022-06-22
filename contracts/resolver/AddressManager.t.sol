pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";
import "../state-machine/StateTransition.sol";

contract TestAddressManager is AddressManager {
    ForgeVM public constant vm = ForgeVM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

    using MerkleMountainRange for CompactMerkleTree;
    CompactMerkleTree _trees;

    AddressManager addressManager;
    RollupStateChain rollupStateChain1;
    RollupInputChain rollupInputChain1;
    L1CrossLayerWitness l1CrossLayerWitness1;
    L2CrossLayerWitness l2CrossLayerWitness2;
    TestERC20 feeToken;
    StakingManager stakingManager1;
    ProxyAdmin proxyAdmin;
    uint256 constant fraudProofWindow = 3;
    address challengerFactory;
    DAO dao_;
    address stateStorage;
    address inputStorage;
    address owner_ = address(0x8888);

    //    event AddressSet(string _name, address _old, address _new);

    function _initialize() internal {
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
        l1CrossLayerWitness1 = L1CrossLayerWitness(address(proxy));

        // deploy L2CrossLayerWitness
        L2CrossLayerWitness l2CrossLayerWitnessLogic = new L2CrossLayerWitness();
        proxy = new TransparentUpgradeableProxy(
            address(l2CrossLayerWitnessLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(L2CrossLayerWitness.initialize.selector, address(addressManager))
        );
        l2CrossLayerWitness2 = L2CrossLayerWitness(address(proxy));

        feeToken = new TestERC20("test token", "test");

        RollupStateChain rollupStateChainLogic = new RollupStateChain();
        proxy = new TransparentUpgradeableProxy(
            address(rollupStateChainLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(RollupStateChain.initialize.selector, address(addressManager), fraudProofWindow)
        );
        rollupStateChain1 = RollupStateChain(address(proxy));

        // TODO: use normal challenge factory
        challengerFactory = address(new MockChallengeFactory());

        // deploy dao
        DAO daoLogic = new DAO();
        proxy = new TransparentUpgradeableProxy(
            address(daoLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(DAO.initialize.selector)
        );
        dao_ = DAO(address(proxy));

        // deploy staking manager
        StakingManager stakingManagerLogic = new StakingManager();
        proxy = new TransparentUpgradeableProxy(
            address(stakingManagerLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(
                StakingManager.initialize.selector,
                address(dao_),
                challengerFactory,
                address(rollupStateChain1),
                address(feeToken),
                1 ether
            )
        );
        stakingManager1 = StakingManager(address(proxy));

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
        rollupInputChain1 = RollupInputChain(address(proxy));

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
        stateStorage = address(address(proxy));

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
        inputStorage = address(address(proxy));
    }

    function setUp() public {
        vm.startPrank(owner_);
        _initialize();

        vm.expectEmit(true, true, true, true, address(addressManager));
        emit AddressSet(AddressName.ROLLUP_INPUT_CHAIN, address(0), address(rollupInputChain1));
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN, address(rollupInputChain1));
        vm.expectEmit(true, true, true, true, address(addressManager));

        emit AddressSet(AddressName.STAKING_MANAGER, address(0), address(stakingManager1));
        addressManager.setAddress(AddressName.STAKING_MANAGER, address(stakingManager1));
        vm.expectEmit(true, true, true, true, address(addressManager));

        emit AddressSet(AddressName.ROLLUP_STATE_CHAIN_CONTAINER, address(0), address(stateStorage));
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN_CONTAINER, stateStorage);
        vm.expectEmit(true, true, true, true, address(addressManager));

        emit AddressSet(AddressName.ROLLUP_INPUT_CHAIN_CONTAINER, address(0), address(inputStorage));
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN_CONTAINER, inputStorage);
        vm.expectEmit(true, true, true, true, address(addressManager));

        emit AddressSet(AddressName.ROLLUP_STATE_CHAIN, address(0), address(rollupStateChain1));
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN, address(rollupStateChain1));
        vm.expectEmit(true, true, true, true, address(addressManager));

        emit AddressSet(AddressName.L1_CROSS_LAYER_WITNESS, address(0), address(l1CrossLayerWitness1));
        addressManager.setAddress(AddressName.L1_CROSS_LAYER_WITNESS, address(l1CrossLayerWitness1));

        vm.expectEmit(true, true, true, true, address(addressManager));
        emit AddressSet(AddressName.L2_CROSS_LAYER_WITNESS, address(0), address(l2CrossLayerWitness2));
        addressManager.setAddress(AddressName.L2_CROSS_LAYER_WITNESS, address(l2CrossLayerWitness2));

        vm.expectEmit(true, true, true, true, address(addressManager));
        emit AddressSet(AddressName.DAO, address(0), address(dao_));
        addressManager.setAddress(AddressName.DAO, address(dao_));

        vm.expectEmit(true, true, true, true, address(addressManager));
        emit AddressSet(AddressName.CHALLENGE_FACTORY, address(0), address(challengerFactory));
        addressManager.setAddress(AddressName.CHALLENGE_FACTORY, challengerFactory);
        vm.stopPrank();
    }

    function testSerAddress() public {
        vm.startPrank(owner_);
        vm.expectEmit(true, true, true, true, address(addressManager));
        emit AddressSet("name", address(0), address(0x8));
        addressManager.setAddress("name", address(0x8));
    }

    function testSerAddressWithCallerIsNotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        addressManager.setAddress("name", address(0x8));
    }

    function testSerAddressWithZeroAddr() public {
        vm.startPrank(owner_);
        vm.expectRevert("empty addr");
        addressManager.setAddress("name", address(0));
    }

    function testSerAddressBatch() public {
        vm.startPrank(owner_);
        vm.expectEmit(true, true, true, true, address(addressManager));
        emit AddressSet("name", address(0), address(0x81));
        emit AddressSet("name1", address(0), address(0x82));
        emit AddressSet("name2", address(0), address(0x83));
        string[] memory _names = new string[](3);
        _names[0] = "name";
        _names[1] = "name1";
        _names[2] = "name2";
        address[] memory addrs = new address[](3);
        addrs[0] = address(0x81);
        addrs[1] = address(0x82);
        addrs[2] = address(0x83);
        addressManager.setAddressBatch(_names, addrs);
    }

    function testSerAddressBatchWithLengthFailed() public {
        vm.startPrank(owner_);
        string[] memory _names = new string[](2);
        _names[0] = "name";
        _names[1] = "name1";
        address[] memory addrs = new address[](3);
        addrs[0] = address(0x81);
        addrs[1] = address(0x82);
        addrs[2] = address(0x83);
        vm.expectRevert("length mismatch");
        addressManager.setAddressBatch(_names, addrs);
    }

    function testSerAddressBatchWithCallerIsNotOwner() public {
        string[] memory _names = new string[](3);
        _names[0] = "name";
        _names[1] = "name1";
        _names[2] = "name2";
        address[] memory addrs = new address[](3);
        addrs[0] = address(0x81);
        addrs[1] = address(0x82);
        addrs[2] = address(0x83);
        vm.expectRevert("Ownable: caller is not the owner");
        addressManager.setAddressBatch(_names, addrs);
    }

    function testGetAddr() public {
        address addr = addressManager.getAddr(AddressName.CHALLENGE_FACTORY);
        require(addr == address(challengerFactory));
        addr = addressManager.getAddr("mock");
        require(addr == address(0));
    }

    function testResolve() public {
        address addr = addressManager.resolve(AddressName.CHALLENGE_FACTORY);
        require(addr == address(challengerFactory));
        vm.expectRevert("no name saved");
        addressManager.resolve("mock");
    }

    function testDao() public {
        IDAO dao = addressManager.dao();
        require(dao == dao_);
    }

    function testRollupInputChain() public {
        IRollupInputChain rollupInputChain = addressManager.rollupInputChain();
        require(rollupInputChain == rollupInputChain1);
    }

    function testRollupInputChainContainer() public {
        IChainStorageContainer rollupInputChainContainer = addressManager.rollupInputChainContainer();
        require(address(rollupInputChainContainer) == address(inputStorage));
    }

    function testRollupStateChain() public {
        IRollupStateChain rollupStateChain = addressManager.rollupStateChain();
        require(address(rollupStateChain) == address(rollupStateChain1));
    }

    function testRollupStateChainContainer() public {
        IChainStorageContainer rollupStateChainContainer = addressManager.rollupStateChainContainer();
        require(address(rollupStateChainContainer) == address(stateStorage));
    }

    function testStakingManager() public {
        IStakingManager stakingManager = addressManager.stakingManager();
        require(stakingManager == stakingManager1);
    }

    function testChallengeFactory() public {
        IChallengeFactory challengeFactory = addressManager.challengeFactory();
        require(address(challengeFactory) == address(challengerFactory));
    }

    function testL1CrossLayerWitness() public {
        IL1CrossLayerWitness l1CrossLayerWitness = addressManager.l1CrossLayerWitness();
        require(l1CrossLayerWitness == l1CrossLayerWitness1);
    }

    function testL2CrossLayerWitness() public {
        IL2CrossLayerWitness l2CrossLayerWitness = addressManager.l2CrossLayerWitness();
        require(l2CrossLayerWitness == l2CrossLayerWitness2);
    }

    function testStateTransition() public {
        StateTransition _stateTransition = new StateTransition();
        vm.startPrank(owner_);
        addressManager.setAddress(AddressName.STATE_TRANSITION, address(_stateTransition));
        IStateTransition stateTransition = addressManager.stateTransition();
        require(stateTransition == _stateTransition);
    }
}
