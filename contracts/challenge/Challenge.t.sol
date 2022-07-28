// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./Challenge.sol";
import "./ChallengeFactory.sol";
import "../interfaces/IChallenge.sol";
import "../interfaces/IChallengeFactory.sol";
import "./DisputeTree.sol";
import "../interfaces/ForgeVM.sol";
import "../libraries/Types.sol";
import "../test-helper/TestBase.sol";
import "../test-helper/TestERC20.sol";
import "./MockProposer.sol";
import "./MockChallenger.sol";

contract MockStateTransition {
    function generateStartState(
        bytes32 rollupInputHash,
        uint64 blockNumber,
        bytes32 parentBlockHash
    ) external pure returns (bytes32) {
        return keccak256(abi.encodePacked(rollupInputHash, blockNumber, parentBlockHash));
    }

    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external view {}

    function executeNextStep(bytes32 stateHash) external pure returns (bytes32) {
        stateHash = "0x0"; //solve warning

        // here is test "state transition is right"
        // when return 0xff , it will revert ; it's difficult to test
        // return bytes32(uint256(0xff));

        // normal case
        return bytes32(uint256(0x1234432112344321));
    }
}

contract TestChallenge is Challenge {
    address testAddress = address(0x7777);
    address testAddress2 = address(0x8888);
    ForgeVM public constant vm = ForgeVM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    bytes32 fake32 = bytes32(uint256(0xff));

    using Types for Types.StateInfo;
    using DisputeTree for mapping(uint256 => DisputeTree.DisputeNode);
    IChallenge challenge;
    ChallengeFactory challengeFactory;
    UpgradeableBeacon challengebeacon;
    AddressManager addressManager;
    DAO dao;
    RollupStateChain rollupstatechain;
    ChainStorageContainer stateStorageContainer;
    RollupInputChain rollupinputchain;
    ChainStorageContainer inputStorageContainer;
    StakingManager stakingManager;
    // TestERchallenger20 feeToken;

    TestERC20 feeToken;
    IChallenge CreateChallenge;
    address Newchallenge;
    //Mock Contract
    MockStateTransition stateTransition;
    MockChallenger challenger1;
    MockChallenger challenger2;
    MockChallenger challenger3;
    MockChallenger challenger4;
    MockProposer proposer;

    function setUp() public {
        vm.startPrank(testAddress);
        // deploy related contract
        addressManager = new AddressManager();
        addressManager.initialize();
        dao = new DAO();
        dao.initialize();
        rollupstatechain = new RollupStateChain();
        rollupstatechain.initialize(address(addressManager), 10);
        stateStorageContainer = new ChainStorageContainer();
        stateStorageContainer.initialize("RollupStateChain", address(addressManager));
        rollupinputchain = new RollupInputChain();
        rollupinputchain.initialize(address(addressManager), 15000000, 3000000, 1234);
        inputStorageContainer = new ChainStorageContainer();
        inputStorageContainer.initialize("RollupStateChain", address(addressManager));
        stateTransition = new MockStateTransition();

        // change addressManager.Address
        addressManager.setAddress("testAddress", testAddress);
        addressManager.setAddress(AddressName.DAO, address(dao));
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN, address(rollupstatechain));
        addressManager.setAddress(AddressName.ROLLUP_STATE_CHAIN_CONTAINER, address(stateStorageContainer));
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN_CONTAINER, address(inputStorageContainer));
        addressManager.setAddress(AddressName.ROLLUP_INPUT_CHAIN, address(rollupinputchain));
        addressManager.setAddress(AddressName.STATE_TRANSITION, address(stateTransition));
        // deploy challengeFactory
        challengeFactory = new ChallengeFactory();
        challenge = new Challenge();
        challengebeacon = new UpgradeableBeacon(address(challenge));
        challengeFactory.initialize(addressManager, address(challengebeacon), 10, 1);
        addressManager.setAddress(AddressName.CHALLENGE_FACTORY, address(challengeFactory));

        // deploy feeToken contract & (mint feeToken & approve feeToken) to testAddress
        feeToken = new TestERC20("test feeToken", "test",18);
        feeToken.transfer(address(testAddress2), 1000 ether);

        //create 3*challenger
        challenger1 = new MockChallenger();
        challenger2 = new MockChallenger();
        challenger3 = new MockChallenger();
        challenger4 = new MockChallenger();
        proposer = new MockProposer();
        dao.setProposerWhitelist(address(proposer), true);

        //transfer feeToken to challengers & proposer
        feeToken.transfer(address(challenger1), 10 ether);
        feeToken.transfer(address(challenger2), 10 ether);
        feeToken.transfer(address(challenger3), 10 ether);
        feeToken.transfer(address(challenger4), 10 ether);
        feeToken.transfer(address(proposer), 100 ether);

        // new stakingManager + setAddress
        stakingManager = new StakingManager();
        stakingManager.initialize(
            address(dao),
            address(challengeFactory),
            address(rollupstatechain),
            address(feeToken),
            1 ether
        );
        feeToken.transfer(address(stakingManager), 1 ether);
        addressManager.setAddress(AddressName.STAKING_MANAGER, address(stakingManager));
        //dao set challenger whitelist
        dao.setChallengerWhitelist(address(challenger1), true);
        dao.setChallengerWhitelist(address(challenger2), true);
        dao.setChallengerWhitelist(address(challenger3), true);
        dao.setChallengerWhitelist(address(challenger4), true);

        //proposer deposite
        proposer.setStakingManager(stakingManager);
        proposer.approve(feeToken, address(stakingManager), 1 ether);
        proposer.deposit();
        vm.stopPrank();

        //chainContainer append stateinfo
        vm.startPrank(address(rollupstatechain));
        Types.StateInfo memory challengeStateinfo = Types.StateInfo(bytes32("0x1"), 1, 1, address(proposer));
        Types.StateInfo memory parentStateinfo = Types.StateInfo(bytes32("0x1"), 0, 1, address(proposer));
        stateStorageContainer.append(Types.hash(parentStateinfo));
        stateStorageContainer.append(Types.hash(challengeStateinfo));
        inputStorageContainer.append(Types.hash(parentStateinfo));
        inputStorageContainer.append(Types.hash(challengeStateinfo));
        vm.stopPrank();
        // challengeFactory --create--> challenge
        vm.startPrank(address(challenger1));
        challenger1.approve(feeToken, address(challengeFactory), 10 ether);
        challenger1.approve(feeToken, address(stakingManager), 1 ether);
        challengeFactory.newChallenge(challengeStateinfo, parentStateinfo);
        bytes32 _hash = challengeStateinfo.hash();
        Newchallenge = challengeFactory.getChallengedContract(_hash);
        stakingManager.deposit();
        vm.stopPrank();

        // challenger approve feeToken to challenge
        challenger2.approve(feeToken, Newchallenge, 1 ether);
        challenger3.approve(feeToken, Newchallenge, 1 ether);

        challenger1.setChallenge(IChallenge(Newchallenge));
        challenger2.setChallenge(IChallenge(Newchallenge));
        challenger3.setChallenge(IChallenge(Newchallenge));
        challenger4.setChallenge(IChallenge(Newchallenge));
        proposer.setChallenge(IChallenge(Newchallenge));
    }

    /* test challenge create() 
   function create() has been tested in ChallengeFactory
   test Challenge --create-- Challenge
*/

    // function create() have been  tested in ChallengeFactory
    // when EOA call challenge.create() ;no revert
    function testChallengeCreateChallenge() public {
        vm.startPrank(testAddress);
        Types.StateInfo memory challengeStateinfo = Types.StateInfo(bytes32("0x1"), 1, 1, address(1));
        IChallenge(Newchallenge).create(bytes32(0), testAddress, 10, challengeStateinfo, 1);
        vm.stopPrank();
    }

    /* test challenge initialize() 
   1.test Fail
*/

    // when stage not stage1 ; revert ("only started stage")
    function testInitializeStageWrong() public {
        vm.startPrank(address(proposer));
        Challenge challenge2 = new Challenge();
        vm.expectRevert("only started stage");
        IChallenge(challenge2).initialize(2, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
    }

    // test Fail
    // when Initialize() with msg.sender  !=  Proposer ; revert("only proposer")
    function testInitializeNotProposer() public {
        vm.startPrank(testAddress2);
        vm.expectRevert("only proposer");
        IChallenge(Newchallenge).initialize(2, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
    }

    // test Fail Initialize()
    // when block.number > block.number + proposerTimeLimit ; revert("wrong context")
    function testInitializeWrongBlockNumer() public {
        vm.startPrank(address(proposer));
        vm.roll(block.number + 100);
        vm.expectRevert("wrong context");
        IChallenge(Newchallenge).initialize(2, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
    }

    // test Fail Initialize()
    // when endStep <= 1; revert("wrong context")
    function testInitializeWrongEndStep() public {
        vm.startPrank(address(proposer));
        vm.expectRevert("wrong context");
        IChallenge(Newchallenge).initialize(1, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
    }

    // test Fail Initialize()
    // when _midSystemState == 0; revert ("illegal state root")
    function testInitializeWrongMidsystemState() public {
        vm.startPrank(address(proposer));
        vm.expectRevert("illegal state root");
        IChallenge(Newchallenge).initialize(2, bytes32(uint256(0x6666)), 0);
        vm.stopPrank();
    }

    /* test challenge initialize() 
   2.test pass
*/

    // test pass Initialize()
    // test event
    function testInitializepass() public {
        vm.startPrank(address(proposer));
        vm.expectEmit(false, false, false, true);
        emit ChallengeInitialized(2, fake32);
        IChallenge(Newchallenge).initialize(2, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
    }

    /*test revealMidStates()
1.test fail */

    // test wrong length
    // when _nodeKeys.length == 0 , revert ("illegal length")
    function testRevealMidStatesLengthEqual0() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        uint256[] memory _nodeKeys = new uint256[](0);
        bytes32[] memory _roots = new bytes32[](0);
        vm.expectRevert("illegal length");
        IChallenge(Newchallenge).revealMidStates(_nodeKeys, _roots);
        vm.stopPrank();
    }

    // test wrong length
    // when _nodeKeys.length != _stateRoots.length , revert ("illegal length")
    function testRevealMidStatesWrongLength() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        uint256[] memory _nodeKeys = new uint256[](5);
        bytes32[] memory _roots = new bytes32[](4);
        vm.expectRevert("illegal length");
        IChallenge(Newchallenge).revealMidStates(_nodeKeys, _roots);
        vm.stopPrank();
    }

    // test time out
    // when block.number > node.expireAfterBlock , revert ("time out")
    function testRevealMidStatesTimeOut() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        uint256[] memory _nodeKeys = new uint256[](4);
        bytes32[] memory _roots = new bytes32[](4);
        vm.roll(100);
        uint256 rootKey = DisputeTree.encodeNodeKey(0, 5);
        _nodeKeys[0] = rootKey;
        vm.expectRevert("time out");
        IChallenge(Newchallenge).revealMidStates(_nodeKeys, _roots);
        vm.stopPrank();
    }

    // test wrong state root
    // when block.number > node.expireAfterBlock , revert ("time out")
    function testRevealMidStatesWrongStateRoot() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        uint256[] memory _nodeKeys = new uint256[](4);
        bytes32[] memory _roots = new bytes32[](4);
        vm.expectRevert("wrong state root");
        IChallenge(Newchallenge).revealMidStates(_nodeKeys, _roots);
        vm.stopPrank();
    }

    /*test revealMidStates()
2.test pass & test event */

    //helper-function: init -- stage 1 -> 2
    function init0_5() internal {
        //0->5
        proposer.initialize(5, fake32);
    }

    //helper-function: stage 2 selectDisputeBranch
    function select(
        uint128 _s,
        uint128 _e,
        bool _isLeft
    ) internal {
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](1);
        _nodeKeys[0] = DisputeTree.encodeNodeKey(_s, _e);
        _isLefts[0] = _isLeft;
        challenger1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    //helper-function: stage 2 revealMidStates
    function revealChild(
        uint128 _s,
        uint128 _e,
        bool _isLeft
    ) internal {
        uint256[] memory _nodeKeys = new uint256[](1);
        bytes32[] memory _roots = new bytes32[](1);
        uint128 _mid = DisputeTree.middle(_s, _e);
        if (_isLeft) {
            _e = _mid;
        } else {
            _s = _mid;
        }
        _nodeKeys[0] = DisputeTree.encodeNodeKey(_s, _e);
        _roots[0] = fake32;
        proposer.revealMidStates(_nodeKeys, _roots);
    }

    // test reveal pass & event
    function testRevealMidStatesPass() public {
        init0_5();
        select(0, 5, true);
        //reveal
        (uint128 _s, uint128 _e, bool _isLeft) = (0, 5, true);
        uint256[] memory _nodeKeys = new uint256[](1);
        bytes32[] memory _roots = new bytes32[](1);
        uint128 _mid = DisputeTree.middle(_s, _e);
        if (_isLeft) {
            _e = _mid;
        } else {
            _s = _mid;
        }
        _nodeKeys[0] = DisputeTree.encodeNodeKey(_s, _e);
        _roots[0] = fake32;
        vm.expectEmit(false, false, false, true);
        emit MidStateRevealed(_nodeKeys, _roots);
        proposer.revealMidStates(_nodeKeys, _roots);
    }

    /*test proposerTimeout
1.test fail */

    //when stage == ChallengeStage.Started && not TimeOut, revert
    function testProposerTimeoutStartedNotTimeOut() public {
        vm.expectRevert("initialize challenge info not timeout");
        IChallenge(Newchallenge).proposerTimeout(0);
    }

    //when stage == running && nextNode.parent == 0 , revert ("one step don't need to prove")
    function testProposerTimeoutRunningParentEqual0() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        vm.expectRevert("one step don't need to prove");
        IChallenge(Newchallenge).proposerTimeout(0);
    }

    //when stage == running &&  stepUpper == stepLower + 1 , revert ("one step don't need to prove")
    function testProposerTimeoutRunningOneStep() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        revealChild(0, 2, true);
        revealChild(0, 5, true);
        vm.stopPrank();
        vm.startPrank(address(challenger1));
        select(0, 5, true);
        select(0, 2, true);
        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 1);
        vm.roll(100);
        vm.stopPrank();
        vm.expectRevert("one step don't need to prove");
        IChallenge(Newchallenge).proposerTimeout(_nodeKey);
    }

    //when block.number <= nextNode.expireAfterBlock , revert ("report mid state not timeout")
    function testProposerTimeoutRunningMidStateNotTimeOut() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();

        vm.startPrank(address(proposer));
        revealChild(0, 5, true);
        vm.stopPrank();

        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 2);

        vm.expectRevert("report mid state not timeout");
        IChallenge(Newchallenge).proposerTimeout(_nodeKey);
    }

    //when nextNode.midStateRoot != 0 , revert ("mid state root is revealed")
    function testProposerTimeoutRunningMidStateRevealed() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();

        vm.startPrank(address(proposer));
        revealChild(0, 5, true);
        vm.stopPrank();

        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 2);
        vm.roll(100);
        vm.expectRevert("mid state root is revealed");
        IChallenge(Newchallenge).proposerTimeout(_nodeKey);
    }

    /*test proposerTimeout
2.test pass */

    //1.stage Started Pass
    function testProposerTimeoutStartedPass() public {
        vm.startPrank(address(Newchallenge));
        vm.roll(100);
        IChallenge(Newchallenge).proposerTimeout(0);
        vm.stopPrank();
    }

    //1.stage Running Pass
    function testProposerTimeoutRunningPass() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();
        vm.roll(100);
        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 2);
        IChallenge(Newchallenge).proposerTimeout(_nodeKey);
    }

    /*test selectDisputeBranch
1.test Fail */

    // 1:challenger can only chose exist parent node to derive new selected node.
    function testSelectDisputeBranchWithParentNotExist() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](1);
        _nodeKeys[0] = DisputeTree.encodeNodeKey(0, 2);
        vm.expectRevert("parent not exist");
        challenger1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    // 2:can't derive old exist node.
    function testSelectDisputeBranchWithDeriveOldChild() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.expectRevert("already init");
        select(0, 5, true);
        vm.stopPrank();
    }

    // 3:challenger can only select one branch
    function testFailSelectDisputeBranchWithSelect2Branch() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        select(0, 5, false);
        vm.stopPrank();
    }

    // 4:challenger have to deposit to challenge contract.
    function testFailSelectDisputeBranchWithInsufficientAllowance() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        //challenger4 unstake
        vm.startPrank(address(challenger4));
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](1);
        _nodeKeys[0] = DisputeTree.encodeNodeKey(0, 5);
        _isLefts[0] = true;
        challenger4.selectDisputeBranch(_nodeKeys, _isLefts);
        vm.stopPrank();
    }

    // 5:challenger can't supply empty parentKey or inconsistent length between parentKeys and leftChild flags.
    // when  _parentNodeKeys.length == 0 , revert
    function testSelectDisputeBranchWithLengthEqual0() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        uint256[] memory _nodeKeys = new uint256[](0);
        bool[] memory _isLefts = new bool[](0);
        vm.expectRevert("inconsistent length");
        challenger1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    // when  _nodeKeys.length != _stateRoots.length , revert
    function testSelectDisputeBranchWithLengthNotEquivalent() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](2);
        vm.expectRevert("inconsistent length");
        challenger1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    // 6:one step can't drive child anymore
    function testSelectDisputeBranchWithOneStep() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();

        vm.startPrank(address(proposer));
        revealChild(0, 5, true);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 2, true);
        vm.stopPrank();

        vm.startPrank(address(proposer));
        revealChild(0, 2, true);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        vm.expectRevert("one step have no child");
        select(0, 1, true);
        vm.stopPrank();
    }

    /*test selectDisputeBranch
2.test Pass */
    function testSelectDisputeBranchPass() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](1);
        _nodeKeys[0] = DisputeTree.encodeNodeKey(0, 5);
        _isLefts[0] = true;

        // test event
        uint256[] memory childkeys = new uint256[](1);
        childkeys[0] = 680564733841876926926749214863536422912;
        vm.expectEmit(true, false, false, true);
        emit DisputeBranchSelected(address(challenger1), childkeys, 11);

        //select
        challenger1.selectDisputeBranch(_nodeKeys, _isLefts);
        vm.stopPrank();
    }

    /*test execOneStepTransition
1.test Fail */

    //stage 2->3
    function exec(uint128 _start) internal {
        require(_start < _start + 1);
        //over flow
        uint256 _nodeKey = DisputeTree.encodeNodeKey(_start, _start + 1);
        IChallenge(Newchallenge).execOneStepTransition(_nodeKey);
    }

    // when disputeTree[_leafNodeKey].parent == 0 , revert("not one step node")
    function testExecOneStepTransitionWithNotOneStep() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.expectRevert("not one step node");
        exec(0);
    }

    // when _stepUpper > 1 + _stepLower , revert("not one step node")
    // (_stepUpper = 2 ; _stepLower = 0)
    function testExecOneStepTransitionWithWrongNode() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();

        vm.startPrank(address(proposer));
        revealChild(0, 5, true);
        vm.stopPrank();

        vm.expectRevert("not one step node");
        IChallenge(Newchallenge).execOneStepTransition(DisputeTree.encodeNodeKey(0, 2));
    }

    /*test execOneStepTransition
2.test Pass */
    // the one step is actually right:start systemState is right, end systemState is wrong.
    function testExecOneStepTransitionPass() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();

        vm.startPrank(address(proposer));
        revealChild(0, 5, true);
        vm.stopPrank();

        vm.startPrank(address(proposer));
        revealChild(0, 2, true);
        vm.stopPrank();

        vm.startPrank(address(challenger1));
        select(0, 2, true);
        vm.stopPrank();

        //test event
        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 1);
        vm.expectEmit(false, false, false, true);
        emit OneStepTransition(0, fake32, bytes32(uint256(0x1234432112344321)));
        IChallenge(Newchallenge).execOneStepTransition(_nodeKey);
    }

    /*test claimProposerWin
1.test Fail */

    // test modifier "afterBlockConfirmed"
    function testClaimProposerWinBeforeConfirmed() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        vm.expectRevert("block not confirmed");
        IChallenge(Newchallenge).claimProposerWin();
    }

    // test modifier "stage2"
    function testClaimProposerWinNotStage2() public {
        vm.warp(100);
        vm.expectRevert("only running stage");
        IChallenge(Newchallenge).claimProposerWin();
    }

    /*test claimProposerWin
2.test Pass */

    function testClaimProposerWinPass() public {
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        vm.warp(100);
        uint256 amount = feeToken.balanceOf(Newchallenge);
        uint256 balance = feeToken.balanceOf(address(proposer));
        vm.expectEmit(false, false, false, true);
        emit ProposerWin(address(proposer), amount);
        IChallenge(Newchallenge).claimProposerWin();
        require(feeToken.balanceOf(address(proposer)) == balance + amount, "proposer win & transfer error");
    }

    /*test claimChallengerWin*/

    // AtStage1: Challenger Win & proposer not init
    function testClaimChallengerWinAtStage1() public {
        //1.  Newchallenge.proposerTimeout(0)
        vm.roll(50);
        IChallenge(Newchallenge).proposerTimeout(0);

        //2. rollback & append StateInfo
        Types.StateInfo memory testStateinfo2 = Types.StateInfo(bytes32("0x3"), 1, 1, address(proposer));
        vm.startPrank(address(rollupstatechain));
        stateStorageContainer.append(Types.hash(testStateinfo2));
        vm.stopPrank();

        uint256 balance = feeToken.balanceOf(address(challenger1));
        uint256 amount = feeToken.balanceOf(address(Newchallenge));
        vm.warp(100);
        //3. Newchallenge.claimChallengerWin
        IChallenge(Newchallenge).claimChallengerWin(address(0), testStateinfo2);
        require(feeToken.balanceOf(address(challenger1)) == balance + amount + stakingManager.price(), "claim error");
    }

    //  Challenger Win & At Stage1 & Double Claim
    function testClaimChallengerWinWithDoubleClaimAtStage1() public {
        //1.  Newchallenge.proposerTimeout(0)
        vm.roll(50);
        IChallenge(Newchallenge).proposerTimeout(0);

        //2. rollback & append StateInfo
        Types.StateInfo memory testStateinfo2 = Types.StateInfo(bytes32("0x3"), 1, 1, address(proposer));
        vm.startPrank(address(rollupstatechain));
        stateStorageContainer.append(Types.hash(testStateinfo2));
        vm.stopPrank();

        uint256 balance = feeToken.balanceOf(address(challenger1));
        uint256 amount = feeToken.balanceOf(address(Newchallenge));
        vm.warp(100);
        //3. Newchallenge.claimChallengerWin
        IChallenge(Newchallenge).claimChallengerWin(address(0), testStateinfo2);

        uint256 _old = feeToken.balanceOf(address(challenger1));
        require(_old == balance + amount + stakingManager.price(), "claim error");
        //4. Newchallenge.claimChallengerWin *2
        IChallenge(Newchallenge).claimChallengerWin(address(0), testStateinfo2);
        require(feeToken.balanceOf(address(challenger1)) == _old, "wired added");
    }

    // AtStage2:proposer reveal timeout, so there may be multi challengers,
    //          but only true challenger can eat cake(dup eat is allowed too).
    // test (3 * Challenger)  divide  Cake
    function testClaimChallengerWinAtStage2With3ChallengerEatCake() public {
        // init
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        // challenger2 select
        MockChallenger _old = challenger1;
        challenger1 = challenger2;
        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();
        // revealChild
        vm.startPrank(address(proposer));
        revealChild(0, 5, true);
        vm.stopPrank();
        // challenger3 select
        challenger1 = challenger3;
        vm.startPrank(address(challenger1));
        select(0, 2, true);
        vm.stopPrank();
        // execOneStepTransition
        exec(0);
        challenger1 = _old;
        // roll back & append stateInfo
        Types.StateInfo memory testStateinfo2 = Types.StateInfo(bytes32("0x3"), 1, 1, address(proposer));
        vm.startPrank(address(rollupstatechain));
        stateStorageContainer.append(Types.hash(testStateinfo2));
        vm.stopPrank();
        vm.warp(100);

        require(feeToken.balanceOf(address(challenger1)) < 9 ether, "challenger1 not 0");
        IChallenge(Newchallenge).claimChallengerWin(address(challenger1), testStateinfo2);
        require(feeToken.balanceOf(address(challenger1)) > 9.3 ether, "challenger1 not add");

        require(feeToken.balanceOf(address(challenger2)) < 10 ether, "challenger2 not 0");
        IChallenge(Newchallenge).claimChallengerWin(address(challenger2), testStateinfo2);
        require(feeToken.balanceOf(address(challenger2)) > 10.3 ether, "challenger2 not add");
        require(feeToken.balanceOf(address(challenger3)) < 10 ether, "challenger3 not 0");
        IChallenge(Newchallenge).claimChallengerWin(address(challenger3), testStateinfo2);
        require(feeToken.balanceOf(address(challenger3)) > 10.3 ether, "challenger3 not 0");
        require(feeToken.balanceOf(Newchallenge) == 1, "not consume out");
    }

    // test challenger1 init & challenger2 claim
    function testClaimChallengerWinWithWrongChallenger() public {
        // init
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        // select
        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();
        // proposerTimeout
        vm.roll(100);
        IChallenge(Newchallenge).proposerTimeout(DisputeTree.encodeNodeKey(0, 2));

        // append stateInfo
        Types.StateInfo memory testStateinfo2 = Types.StateInfo(bytes32("0x3"), 1, 1, address(proposer));
        vm.startPrank(address(rollupstatechain));
        stateStorageContainer.append(Types.hash(testStateinfo2));
        vm.stopPrank();
        // WiredChallenger claim
        vm.warp(100);
        vm.expectRevert("you can't eat cake");
        IChallenge(Newchallenge).claimChallengerWin(address(challenger2), testStateinfo2);
    }

    // test challenger1 (eat cake *2)
    function testClaimChallengerWinWithDupEat() public {
        // init
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();
        // select
        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();
        // proposerTimeout
        vm.roll(100);
        IChallenge(Newchallenge).proposerTimeout(DisputeTree.encodeNodeKey(0, 2));
        // append stateInfo
        Types.StateInfo memory testStateinfo2 = Types.StateInfo(bytes32("0x3"), 1, 1, address(proposer));
        vm.startPrank(address(rollupstatechain));
        stateStorageContainer.append(Types.hash(testStateinfo2));
        vm.stopPrank();
        // WiredChallenger claim
        vm.warp(100);
        IChallenge(Newchallenge).claimChallengerWin(address(challenger1), testStateinfo2);

        vm.expectRevert("you can't eat cake");
        IChallenge(Newchallenge).claimChallengerWin(address(challenger1), testStateinfo2);
    }

    // test 2 Branch & transfer to dao
    function testClaimChallengerWinWith2Branch() public {
        // init
        vm.startPrank(address(proposer));
        IChallenge(Newchallenge).initialize(5, bytes32(uint256(0x6666)), fake32);
        vm.stopPrank();

        // challenger1 select
        MockChallenger old = challenger1;
        vm.startPrank(address(challenger1));
        select(0, 5, true);
        vm.stopPrank();

        // challenger2 select
        challenger1 = challenger2;
        vm.startPrank(address(challenger1));
        select(0, 5, false);
        vm.stopPrank();
        vm.roll(100);

        // Newchallenge - proposerTimeout
        challenger1 = old;
        IChallenge(Newchallenge).proposerTimeout(DisputeTree.encodeNodeKey(0, 2));

        // roll back , then  append stateInfo
        Types.StateInfo memory testStateinfo2 = Types.StateInfo(bytes32("0x3"), 1, 1, address(proposer));
        vm.startPrank(address(rollupstatechain));
        stateStorageContainer.append(Types.hash(testStateinfo2));
        vm.stopPrank();
        vm.warp(100);

        require(feeToken.balanceOf(address(dao)) == 0, "dao not 0");
        IChallenge(Newchallenge).claimChallengerWin(address(challenger1), testStateinfo2);

        require(feeToken.balanceOf(address(challenger1)) < 9 ether);
        require(feeToken.balanceOf(address(dao)) > 0, "dao not added");
        require(feeToken.balanceOf(Newchallenge) == 0, "challenge contract not 0");
    }
}
