// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.13;
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/ForgeVM.sol";
import "../interfaces/IL1StandardBridge.sol";
import "../bridge/L1StandardBridge.sol";
import "../bridge/L2StandardBridge.sol";
import "./TestBase.sol";
import "./MockContract.sol";
import "../cross-layer/CrossLayerContext.sol";
import "../state-machine/StateTransition.sol";
import "../challenge/Challenge.sol";
import "../challenge/ChallengeFactory.sol";


/*  test code(block number & rpc should update):

forge test --fork-url http://172.168.3.70:8501 
--fork-block-number 608381 --match-contract "testUpgradeForkL1"*/

contract testUpgradeForkL1 is TestBase {
    address owner;
    address proxyAdminAddrl1;
    address proxyAdminAddrl2;
    function setUp() public {
        owner = address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266);
        proxyAdminAddrl1 = address(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512);
    }

/* 1.test L1code */
    function testUpgradeL1StandardBridge() public {
        // get l1 bridge contract & new mockl1bridge contract
        address l1BridgeProxy = address(0xD8a5a9b31c3C0232E196d518E89Fd8bF83AcAd43);
        MockL1Bridge newl1StandardBridge = new MockL1Bridge();
        // upgrade l1bridge contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1BridgeProxy)), address(newl1StandardBridge)
            );
        vm.stopPrank();
        newl1StandardBridge = MockL1Bridge(l1BridgeProxy);
        // test call newL1Bridge
        require(newl1StandardBridge.return1() == 1, "upgrade fail 1");
    }

    function testUpgradeAddressManager() public {
        // get l1 AddressManager contract & new mockl1AddressManager contract
        address l1AddressManager = address(0x2bdCC0de6bE1f7D2ee689a0342D76F52E8EFABa3);
        MockAddressManager newAddressManager = new MockAddressManager();

        // upgrade l1AddressManager contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1AddressManager)), address(newAddressManager)
            );
        vm.stopPrank();

        newAddressManager = MockAddressManager(l1AddressManager);
        // test call newAddressManager
        require(newAddressManager.return6() == 6, "upgrade fail 6");
    }

    function testUpgradeInputChainStorage() public {
        // get l1 InputChainStorage contract & new mockInputChainStorage contract
        address l1InputChainStorage = address(0x04C89607413713Ec9775E14b954286519d836FEf);
        MockChainStorageContainer newl1InputChainStorage = new MockChainStorageContainer();

        // upgrade l1InputChainStorage contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1InputChainStorage)), address(newl1InputChainStorage)
            );
        vm.stopPrank();

        newl1InputChainStorage = MockChainStorageContainer(l1InputChainStorage);
        // test call newl1InputChainStorage
        require(newl1InputChainStorage.return7() == 7, "upgrade fail 7");
    }

    function testUpgradeStateChainStorage() public {
        // get l1 StateChainStorage contract & new mockStateChainStorage contract
        address l1StateChainStorage = address(0xdbC43Ba45381e02825b14322cDdd15eC4B3164E6);
        MockChainStorageContainer newl1StateChainStorage = new MockChainStorageContainer();

        // upgrade l1StateChainStorage contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1StateChainStorage)), address(newl1StateChainStorage)
            );
        vm.stopPrank();

        newl1StateChainStorage = MockChainStorageContainer(l1StateChainStorage);
        // test call newl1StateChainStorage
        require(newl1StateChainStorage.return7() == 7, "upgrade fail 7");
    }

    function testUpgradeRollupInputChain() public {
        // get l1 RollupInputChain contract & new mockRollupInputChain contract
        address l1RollupInputChain = address(0x1fA02b2d6A771842690194Cf62D91bdd92BfE28d);
        MockRollupInputChain newl1RollupInputChain = new MockRollupInputChain();

        // upgrade l1RollupInputChain contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1RollupInputChain)), address(newl1RollupInputChain)
            );
        vm.stopPrank();

        newl1RollupInputChain = MockRollupInputChain(l1RollupInputChain);
        // test call newl1RollupInputChain
        require(newl1RollupInputChain.return8() == 8, "upgrade fail 8");
    }

    function testUpgradeRollupStateChain() public {
        // get l1 RollupStateChain contract & new mockRollupStateChain contract
        address l1RollupStateChain = address(0x1fA02b2d6A771842690194Cf62D91bdd92BfE28d);
        MockRollupStateChain newl1RollupStateChain = new MockRollupStateChain();

        // upgrade l1RollupStateChain contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1RollupStateChain)), address(newl1RollupStateChain)
            );
        vm.stopPrank();

        newl1RollupStateChain = MockRollupStateChain(l1RollupStateChain);
        // test call newl1RollupStateChain
        require(newl1RollupStateChain.return9() == 9, "upgrade fail 9");
    }

    function testUpgradeL1CrossLayerWitness() public {
        // get l1 CrossLayerWitness contract & new mockL1CrossLayerWitness contract
        address l1CrossLayerWitness = address(0x1fA02b2d6A771842690194Cf62D91bdd92BfE28d);
        MockL1CrossLayerWitness newL1CrossLayerWitness = new MockL1CrossLayerWitness();

        // upgrade L1CrossLayerWitness contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1CrossLayerWitness)), address(newL1CrossLayerWitness)
            );
        vm.stopPrank();

        newL1CrossLayerWitness = MockL1CrossLayerWitness(l1CrossLayerWitness);
        // test call newL1CrossLayerWitness
        require(newL1CrossLayerWitness.return3() == 3, "upgrade fail 3");
    }

    function testUpgradeL1StakingManager() public {
        // get l1 StakingManager contract & new mockStakingManager contract
        address L1StakingManager = address(0x922D6956C99E12DFeB3224DEA977D0939758A1Fe);
        MockStakingManager newL1StakingManager = new MockStakingManager();

        // upgrade L1StakingManager contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(L1StakingManager)), address(newL1StakingManager)
            );
        vm.stopPrank();

        newL1StakingManager = MockStakingManager(L1StakingManager);
        // test call newL1StakingManager
        require(newL1StakingManager.return10() == 10, "upgrade fail 10");
    }


    function testUpgradeL1ChallengeFactory() public {
        // get l1 ChallengeFactory contract & new mockChallengeFactory contract
        address l1ChallengeFactory = address(0xB0D4afd8879eD9F52b28595d31B441D079B2Ca07);
        MockChallengeFactory newL1ChallengeFactory = new MockChallengeFactory();

        // upgrade l1ChallengeFactory contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1ChallengeFactory)), address(newL1ChallengeFactory)
            );
        vm.stopPrank();

        newL1ChallengeFactory = MockChallengeFactory(l1ChallengeFactory);
        // test call newL1ChallengeFactory
        require(newL1ChallengeFactory.return12() == 12, "upgrade fail 12");
    }

    function testUpgradeL1DAO() public {
        // get l1 DAO contract & new mockChallengeFactory contract
        address l1DAO = address(0xB0D4afd8879eD9F52b28595d31B441D079B2Ca07);
        MockDAO newL1DAO = new MockDAO();

        // upgrade l1DAO contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1DAO)), address(newL1DAO)
            );
        vm.stopPrank();

        newL1DAO = MockDAO(l1DAO);
        // test call newL1DAO
        require(newL1DAO.return5() == 5, "upgrade fail 5");
    }

    function testUpgradeL1StateTransition() public {
        // get l1 StateTransition contract & new mockStateTransition contract
        address l1StateTransition = address(0x04C89607413713Ec9775E14b954286519d836FEf);
        MockStateTransition newL1StateTransition = new MockStateTransition();

        // upgrade l1StateTransition contract
        vm.startPrank(owner);
        ProxyAdmin(proxyAdminAddrl1).upgrade(TransparentUpgradeableProxy(
            payable(l1StateTransition)), address(newL1StateTransition)
            );
        vm.stopPrank();

        newL1StateTransition = MockStateTransition(l1StateTransition);
        // test call newL1StateTransition
        require(newL1StateTransition.return11() == 11, "upgrade fail 5");
    }

    // test newChallenge error , later can test
    // function testUpgradeL1newChallenge() public {
    //     // get l1 ChallengeLogic contract & new mockChallenge contract
    //     address L1ChallengeLogic = address(0xFD471836031dc5108809D173A067e8486B9047A3);
    //     address ChallengeBeacon = address(0xcbEAF3BDe82155F56486Fb5a1072cb8baAf547cc);

    //     TestERC20 feeToken = TestERC20(0x7bc06c482DEAd17c0e297aFbC32f6e63d3846650);
    //     ChallengeFactory challengeFactory = ChallengeFactory(0xB0D4afd8879eD9F52b28595d31B441D079B2Ca07);
    //     DAO dao = DAO(0x162A433068F51e18b7d13932F27e66a3f99E6890);
    //     IChainStorageContainer stateChainStorage = IChainStorageContainer(0xdbC43Ba45381e02825b14322cDdd15eC4B3164E6);
    //     IChainStorageContainer inputChainStorage = IChainStorageContainer(0x04C89607413713Ec9775E14b954286519d836FEf);
    //     StakingManager L1StakingManager = StakingManager(0x922D6956C99E12DFeB3224DEA977D0939758A1Fe);
        
    //     bytes memory _data;
    //     address oldChallenge = address(new BeaconProxy(ChallengeBeacon, _data));
    //     MockProposer proposer = new MockProposer();

    //     MockChallenge newL1Challenge = new MockChallenge();
        
    //     vm.startPrank(owner);
    //     feeToken.approve(address(challengeFactory), 100 ether);

    //     dao.setChallengerWhitelist(owner, true);
    //     dao.setProposerWhitelist(address(proposer), true);
    //     feeToken.transfer(address(proposer), 100 ether);

    //     proposer.setStakingManager(L1StakingManager);
    //     proposer.approve(feeToken, address(L1StakingManager), 1 ether);
    //     proposer.deposit();

    //     Types.StateInfo memory challengeStateinfo = Types.StateInfo(bytes32("0x1"), 310, 1, address(proposer));
    //     Types.StateInfo memory parentStateinfo = Types.StateInfo(bytes32("0x1"), 309, 1, address(proposer));
        
    //     vm.stopPrank();

    //     // rollupStateChain call stateChainStorage
    //     vm.startPrank(0xc351628EB244ec633d5f21fBD6621e1a683B1181);
    //     stateChainStorage.append(Types.hash(parentStateinfo));
    //     stateChainStorage.append(Types.hash(challengeStateinfo));
    //     vm.stopPrank();

    //     // rollupInputChain call inputChainStorageAddr
    //     vm.startPrank(0x1fA02b2d6A771842690194Cf62D91bdd92BfE28d);
    //     inputChainStorage.append(Types.hash(parentStateinfo));
    //     inputChainStorage.append(Types.hash(challengeStateinfo));
    //     vm.stopPrank();

    //     vm.startPrank(owner);
    //     vm.warp(0);
    //     challengeFactory.newChallenge(challengeStateinfo, parentStateinfo);
    //     vm.stopPrank();
    // }
    
}