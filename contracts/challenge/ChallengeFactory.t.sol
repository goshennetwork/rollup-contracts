// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IAddressResolver.sol";
import "../resolver/AddressManager.sol";
import "./Challenge.sol";
import "./ChallengeFactory.sol";
import "@openzeppelin/contracts/proxy/beacon/BeaconProxy.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/IAddressManager.sol";
import "../interfaces/ForgeVM.sol";
import "../rollup/RollupStateChain.sol";
import "../rollup/RollupInputChain.sol";
import "../rollup/ChainStorageContainer.sol";
import "../test-helper/TestERC20.sol";
import "../staking/StakingManager.sol";
import "../dao/Whitelist.sol";

contract MockStateTransition {
    function generateStartState(
        bytes32 rollupInputHash,
        uint64 blockNumber,
        bytes32 parentBlockHash
    ) external pure returns (bytes32) {
        return keccak256(abi.encodePacked(rollupInputHash, blockNumber, parentBlockHash));
    }
}

contract TestChallengeFactory is ChallengeFactory {
    ForgeVM public constant vm = ForgeVM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    address testAddress = address(0x8888);
    address testAddress2 = address(0x9999);

    using Types for Types.StateInfo;

    ChallengeFactory challengeFactory;
    UpgradeableBeacon challengebeacon;
    AddressManager addressManager;
    Whitelist whitelist;
    RollupStateChain rollupstatechain;
    ChainStorageContainer stateStorageContainer;
    RollupInputChain rollupinputchain;
    ChainStorageContainer inputStorageContainer;
    MockStateTransition stateTransition;

    function setUp() public {
        vm.startPrank(testAddress);
        // deploy related contract
        addressManager = new AddressManager();
        addressManager.initialize();
        whitelist = new Whitelist();
        whitelist.initialize(IAddressResolver(address(addressManager)));
        rollupstatechain = new RollupStateChain();
        rollupstatechain.initialize(address(addressManager), 10);
        stateStorageContainer = new ChainStorageContainer();
        stateStorageContainer.initialize("testAddress", address(addressManager));
        rollupinputchain = new RollupInputChain();
        rollupinputchain.initialize(address(addressManager), 15000000, 3000000, 1234);
        inputStorageContainer = new ChainStorageContainer();
        inputStorageContainer.initialize("testAddress", address(addressManager));
        stateTransition = new MockStateTransition();

        // change addressManager.Address
        addressManager.setAddress("testAddress", testAddress);
        addressManager.setAddress(AddressName.DAO, testAddress);
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN, address(rollupstatechain));
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN_CONTAINER, address(stateStorageContainer));
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN_CONTAINER, address(inputStorageContainer));
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN, address(rollupinputchain));
        addressManager.setAddress(AddressName.STATE_TRANSITION, address(stateTransition));
        addressManager.setAddress(AddressName.WHITELIST, address(whitelist));
        // deploy challengeFactory
        challengeFactory = new ChallengeFactory();
        Challenge challenge = new Challenge();
        challengebeacon = new UpgradeableBeacon(address(challenge));
        challengeFactory.initialize(addressManager, address(challengebeacon), 10, 1);
        vm.stopPrank();

        // deploy token contract & (mint token & approve token) to testAddress2
        vm.startPrank(testAddress2);
        TestERC20 feeToken = new TestERC20("test token", "test", 18);
        feeToken.approve(address(challengeFactory), 100 ether);
        vm.stopPrank();

        vm.startPrank(testAddress);
        StakingManager stakingManager = new StakingManager();
        stakingManager.initialize(address(addressManager), 1 ether);
        addressManager.setAddress(AddressName.STAKING_MANAGER, address(stakingManager));
        addressManager.setAddress(AddressName.FEE_TOKEN, address(feeToken));

        vm.stopPrank();
    }

    /* Test newChallenge
   1.Test Fail */

    // test caller not challenger
    function testNewChallengeNotChallenger() public {
        Types.StateInfo memory stateinfo1 = Types.StateInfo(bytes32("0x0"), 1, 1, address(1));
        vm.expectRevert("only challenger");
        challengeFactory.newChallenge(stateinfo1, stateinfo1);
    }

    // test already challenged
    // create new challenge * 2
    function testNewChallengeAlreadyExist() public {
        vm.startPrank(testAddress);
        whitelist.setChallenger(testAddress2, true);
        Types.StateInfo memory challengeStateinfo = Types.StateInfo(bytes32("0x1"), 1, 1, address(1));
        Types.StateInfo memory parentStateinfo = Types.StateInfo(bytes32("0x1"), 0, 1, address(1));
        stateStorageContainer.append(Types.hash(parentStateinfo));
        stateStorageContainer.append(Types.hash(challengeStateinfo));
        inputStorageContainer.append(Types.hash(parentStateinfo));
        inputStorageContainer.append(Types.hash(challengeStateinfo));
        vm.stopPrank();

        vm.startPrank(testAddress2);
        challengeFactory.newChallenge(challengeStateinfo, parentStateinfo);
        vm.expectRevert("already challenged");
        challengeFactory.newChallenge(challengeStateinfo, parentStateinfo);
        vm.stopPrank();
    }

    // test wrong challenged stateInfo
    function testNewChallengeWrongStateInfo() public {
        vm.startPrank(testAddress);
        whitelist.setChallenger(testAddress2, true);
        vm.stopPrank();

        vm.startPrank(testAddress2);
        Types.StateInfo memory stateinfo1 = Types.StateInfo(bytes32("0x0"), 1, 1, address(1));
        vm.expectRevert("wrong stateInfo");
        challengeFactory.newChallenge(stateinfo1, stateinfo1);
        vm.stopPrank();
    }

    // test state confirmed
    function testNewChallengeStateConfirmed() public {
        vm.startPrank(testAddress);
        whitelist.setChallenger(testAddress2, true);
        Types.StateInfo memory stateinfo1 = Types.StateInfo(bytes32("0x0"), 0, 1, address(1));
        stateStorageContainer.append(Types.hash(stateinfo1));
        vm.stopPrank();

        vm.startPrank(testAddress2);
        vm.warp(11);
        vm.expectRevert("state confirmed");
        challengeFactory.newChallenge(stateinfo1, stateinfo1);
        vm.stopPrank();
    }

    // test wrong parent stateInfo
    function testNewChallengeWrongParentStateInfo() public {
        vm.startPrank(testAddress);
        whitelist.setChallenger(testAddress2, true);
        Types.StateInfo memory stateinfo1 = Types.StateInfo(bytes32("0x0"), 0, 1, address(1));
        Types.StateInfo memory stateinfo2 = Types.StateInfo(bytes32("0x0"), 10, 1, address(1));
        stateStorageContainer.append(Types.hash(stateinfo1));
        vm.stopPrank();

        vm.startPrank(testAddress2);
        vm.expectRevert("wrong stateInfo");
        challengeFactory.newChallenge(stateinfo1, stateinfo2);
        vm.stopPrank();
    }

    // test wrong parent stateInfo (invalid index)
    function testNewChallengeWrongParentStateInfoIndex() public {
        //rollupstateChain append (StateInfo*3)
        vm.startPrank(testAddress);
        whitelist.setChallenger(testAddress2, true);
        Types.StateInfo memory challengeStateinfo = Types.StateInfo(bytes32("0x0"), 2, 1, address(1));
        Types.StateInfo memory parentStateinfo = Types.StateInfo(bytes32("0x0"), 0, 1, address(1));
        stateStorageContainer.append(Types.hash(parentStateinfo));
        stateStorageContainer.append(bytes32("0x0"));
        stateStorageContainer.append(Types.hash(challengeStateinfo));
        vm.stopPrank();

        //parentStateInfo.index + 1 != challengedStateInfo.index
        vm.startPrank(testAddress2);
        vm.expectRevert("wrong parent stateInfo");
        challengeFactory.newChallenge(challengeStateinfo, parentStateinfo);
        vm.stopPrank();
    }

    /* test pass NewChallenge */
    function testNewChallengePass() public {
        vm.startPrank(testAddress);
        whitelist.setChallenger(testAddress2, true);
        Types.StateInfo memory challengeStateinfo = Types.StateInfo(bytes32("0x0"), 1, 1, address(1));
        Types.StateInfo memory parentStateinfo = Types.StateInfo(bytes32("0x0"), 0, 1, address(1));
        stateStorageContainer.append(Types.hash(parentStateinfo));
        stateStorageContainer.append(Types.hash(challengeStateinfo));
        inputStorageContainer.append(Types.hash(parentStateinfo));
        inputStorageContainer.append(Types.hash(challengeStateinfo));
        vm.stopPrank();

        vm.startPrank(testAddress2);
        challengeFactory.newChallenge(challengeStateinfo, parentStateinfo);
        vm.stopPrank();
    }
}
