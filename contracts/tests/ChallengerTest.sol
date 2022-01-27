pragma solidity ^0.8.0;

import "../interfaces/IChallenge.sol";
import "../challenge/Challenge.sol";
import "./mocks/MockChallengeFactory.sol";
import "./mocks/MockStateCommitChain.sol";
import "./mocks/MockStakingManager.sol";
import "./mocks/MockStateTransition.sol";
import "./mocks/MockChallenger.sol";
import "./mocks/MockProposer.sol";
import "../challenge/DisputeTree.sol";

import "./mocks/MockERC20.sol";

interface Vm {
    // Set block.height (newHeight)
    function roll(uint256) external;

    function expectRevert(bytes calldata c) external;
}

contract TestChallenger {
    Vm vm = Vm(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    MockChallengeFactory factory;
    IChallenge challenge;
    MockChallenger c1;
    MockChallenger c2;
    MockChallenger c3;
    MockChallenger c4;
    MockProposer p;
    bytes32 fake32 = bytes32(uint256(0xff));

    function setUp() public {
        factory = new MockChallengeFactory();
        IStateCommitChain scc = new MockStateCommitChain();
        ERC20 erc20 = new MockERC20();
        IStakingManager sm = new MockStakingManager(address(erc20));
        IStateTransition executor = new MockStateTransition();
        factory.init(address(sm), address(executor), address(scc));
        c1 = new MockChallenger();
        c2 = new MockChallenger();
        c3 = new MockChallenger();
        c4 = new MockChallenger();
        p = new MockProposer();
        erc20.transfer(address(c1), 1 ether);
        erc20.transfer(address(c2), 1 ether);
        erc20.transfer(address(c3), 1 ether);
        erc20.transfer(address(p), 1 ether);
        erc20.transfer(address(sm), 1 ether);
        //start
        c1.approve(erc20, address(factory), 0.01 ether);
        challenge = factory.newChallengeWithProposer(address(c1), address(p));
        c2.approve(erc20, address(challenge), 0.01 ether);
        c3.approve(erc20, address(challenge), 0.01 ether);
        c4.approve(erc20, address(challenge), 0.01 ether);
        c1.setChallenge(challenge);
        c2.setChallenge(challenge);
        c3.setChallenge(challenge);
        c4.setChallenge(challenge);
        p.setChallenge(challenge);
    }

    //stage 1 -> 2
    function init0_5() internal {
        //0->5
        p.initialize(5, fake32);
    }

    //stage 2
    function select(
        uint128 _s,
        uint128 _e,
        bool _isLeft
    ) internal {
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](1);
        _nodeKeys[0] = DisputeTree.encodeNodeKey(_s, _e);
        _isLefts[0] = _isLeft;
        c1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    //stage 2
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
        p.revealMidStates(_nodeKeys, _roots);
    }

    //stage 2->3
    function exec(uint128 _start) internal {
        require(_start < _start + 1);
        //over flow
        uint256 _nodeKey = DisputeTree.encodeNodeKey(_start, _start + 1);
        challenge.execOneStepTransition(_nodeKey);
    }

    /**
     * @dev stage1 situation:
     * 1:only proposer can init
     * 2:only init once
     * because init pass 2 params: `endStep` & `midStateRoot`,so
     * 3: `endStep` must more than 1
     * 4: `midStateRoot` can't be 0
     * 5: function have to prove the endState is right and output is right.Protected by executor.No testing here
     * 6:proposer must init in time
     * if proposer doesn't participate in time,game stuck and proposer can be claimed.
     * 7:if proposer time out when init, anyone can make challenge win.
     */

    //1:only proposer can init
    function testFailNotProposer() public {
        challenge.initialize(2, fake32);
    }

    //2:only init once
    function testFailInitDup() public {
        p.initialize(2, fake32);
        p.initialize(3, fake32);
    }

    //3: `endStep` must more than 1
    function testFailInit1step() public {
        //1 end step not allowed
        p.initialize(1, fake32);
    }

    //4: `midStateRoot` can't be 0
    function testRevertInit0state() public {
        //0 mid state not allowed
        vm.expectRevert("0 system state root is illegal");
        p.initialize(2, 0);
    }

    //6:proposer must init in time
    function testFailRevertInitTimeout() public {
        vm.roll(100);
        init0_5();
    }

    //7:proposer failed when timeout.
    function testGameOverInitErr() public {
        vm.roll(100);
        challenge.proposerTimeout(0);
    }

    //7
    function testFailGameOverInitOk() public {
        challenge.proposerTimeout(0);
    }

    //normal
    function testInit() public {
        init0_5();
    }

    /**
     * @dev stage2, stage 2 do 2 tasks:challenger select one branch;proposer reveal mid state.
     * 1:challenger can only chose exist parent node to derive new selected node.
     * 2:can't derive old exist node.
     * 3:challenger can only select one branch
     * 4:challenger have to deposit to challenge contract.
     * 5:challenger can't supply empty parentKey or inconsistent length between parentKeys and leftChild flags.
     * 6:one step can't drive child anymore
     */

    //1
    function testRevertInvalidParent() public {
        init0_5();
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](1);
        _nodeKeys[0] = DisputeTree.encodeNodeKey(0, 2);
        vm.expectRevert("parent not exist");
        c1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    //2
    function testFailDeriveOldChild() public {
        init0_5();
        select(0, 5, true);
        select(0, 5, true);
    }

    //3
    function testRevertSelect2Branch() public {
        init0_5();
        select(0, 5, true);
        vm.expectRevert("you can only select one branch");
        select(0, 5, false);
    }

    //4
    function testFailInsufficientWhenSelect() public {
        init0_5();
        c1 = c4;
        //change sender
        select(0, 5, true);
    }

    //5
    function testRevertEmptyParent() public {
        init0_5();
        uint256[] memory _nodeKeys = new uint256[](0);
        bool[] memory _isLefts = new bool[](0);
        vm.expectRevert("inconsistent length");
        c1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    //5
    function testRevertInconsistentParent() public {
        init0_5();
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](2);
        vm.expectRevert("inconsistent length");
        c1.selectDisputeBranch(_nodeKeys, _isLefts);
    }

    //6
    function testRevertSelectOneStep() public {
        init0_5();
        select(0, 5, true);
        revealChild(0, 5, true);
        select(0, 2, true);
        revealChild(0, 2, true);
        vm.expectRevert("one step have no child");
        select(0, 1, true);
    }

    function testSelect1() public {
        init0_5();
        select(0, 5, true);
    }

    /**
     * reveal task in stage 2:
     * 1: only proposer can reveal
     * 2: only reveal in time
     * 3: only reveal specific state once
     * 4: cant reveal state as 0
     * 5: can't provide empty nodeKeys num or inconsistent num with stateRoots num
     * 6：can reveal in advance
     * 7: if proposer time out when revealing(except one step), anyone can make challenger win in this game.
     */

    //1
    function testRevertInvalidRevealer() public {
        init0_5();
        p = MockProposer(address(challenge));
        //now send is 0x000000..00.
        vm.expectRevert("only proposer");
        revealChild(0, 5, true);
    }

    //2
    function testRevertRevealTimeout() public {
        init0_5();
        select(0, 5, true);
        vm.roll(100);
        vm.expectRevert("time out");
        revealChild(0, 5, true);
    }

    //3
    function testRevertRevealDup() public {
        init0_5();
        revealChild(0, 5, true);
        vm.expectRevert("dup new state or zero new state");
        revealChild(0, 5, true);
    }

    //4
    function testRevertRevealEmptyState() public {
        init0_5();
        vm.expectRevert("dup new state or zero new state");
        fake32 = 0;
        //0
        revealChild(0, 5, true);
    }

    //5
    function testRevertRevealEmptyNum() public {
        init0_5();
        uint256[] memory _nodeKeys = new uint256[](0);
        bytes32[] memory _roots = new bytes32[](0);
        vm.expectRevert("illegal length");
        p.revealMidStates(_nodeKeys, _roots);
    }

    //5
    function testRevertRevealInconsistentNum() public {
        init0_5();
        uint256[] memory _nodeKeys = new uint256[](1);
        bytes32[] memory _roots = new bytes32[](0);
        _nodeKeys[0] = 1;
        vm.expectRevert("illegal length");
        p.revealMidStates(_nodeKeys, _roots);
    }

    //6
    function testRevealAhead() public {
        init0_5();
        revealChild(0, 2, true);
        revealChild(0, 5, true);
    }

    //7
    function testGameOverRevealTimeout() public {
        init0_5();
        select(0, 5, true);
        vm.roll(100);
        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 2);
        challenge.proposerTimeout(_nodeKey);
    }

    //7
    function testRevertGameOverOneStep() public {
        init0_5();
        revealChild(0, 2, true);
        revealChild(0, 5, true);
        select(0, 5, true);
        select(0, 2, true);
        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 1);
        vm.roll(100);
        vm.expectRevert("one step don't need to prove");
        challenge.proposerTimeout(_nodeKey);
    }

    //7
    function testRevertGameOverButRevealOk() public {
        init0_5();
        select(0, 5, true);
        revealChild(0, 5, true);
        vm.roll(100);
        uint256 _nodeKey = DisputeTree.encodeNodeKey(0, 2);
        vm.expectRevert("mid state root is revealed");
        challenge.proposerTimeout(_nodeKey);
    }

    //normal reveal
    function testReveal1() public {
        init0_5();
        select(0, 5, true);
        revealChild(0, 5, true);
    }

    /**
     * @dev in stage,if challenger finally find out one step, they can exec the program to prove he is right and end the game.
     * 1: must provide exist one step node.
     * 2: the one step is actually right:start systemState is right, end systemState is wrong.
     */
    //1
    function testRevertExecEmptyNode() public {
        init0_5();
        vm.expectRevert("not one step node");
        exec(0);
    }

    //1
    function testRevertExecWrongNode() public {
        init0_5();
        select(0, 5, true);
        revealChild(0, 5, true);
        vm.expectRevert("not one step node");
        challenge.execOneStepTransition(DisputeTree.encodeNodeKey(0, 2));
    }

    //2
    function testExec() public {
        init0_5();
        select(0, 5, true);
        revealChild(0, 5, true);
        revealChild(0, 2, true);
        select(0, 2, true);
        exec(0);
    }

    /**
     * @dev in stage 3.only challenger need to claim out payback.and game over may at stage 1 or 2
     * 1：in stage 1,proposer not init, so there is only one challenger called `creator`.Dup claim have no profit.
     * 2: in stage 2,proposer reveal timeout, so there may be multi challengers,but only true challenger can eat cake(dup eat is allowed too).
     * 3: only one branch is allowed.otherwise, transfer reward to dao
     */

    //1
    function testClaimWinAtStage1() public {
        vm.roll(100);
        challenge.proposerTimeout(0);
        require(factory.stakingManager().token().balanceOf(address(c1)) < 1 ether, "not 0");
        challenge.claimChallengerWin(address(0));
        require(factory.stakingManager().token().balanceOf(address(c1)) > 1 ether, "not add");
    }

    function testClaimDupAtStage1() public {
        vm.roll(100);
        challenge.proposerTimeout(0);
        challenge.claimChallengerWin(address(0));
        uint256 _old = factory.stakingManager().token().balanceOf(address(c1));
        require(_old > 1 ether, "not add");
        challenge.claimChallengerWin(address(0));
        require(factory.stakingManager().token().balanceOf(address(c1)) == _old, "wired added");
    }

    //2
    function testClaimWinAtStage2() public {
        init0_5();
        MockChallenger _old = c1;
        c1 = c2;
        select(0, 5, true);
        revealChild(0, 5, true);
        c1 = c3;
        select(0, 2, true);
        exec(0);
        c1 = _old;
        IERC20 token = factory.stakingManager().token();
        require(token.balanceOf(address(c1)) < 1 ether, "c1 not 0");
        challenge.claimChallengerWin(address(c1));
        require(token.balanceOf(address(c1)) > 1 ether, "c1 not add");
        require(token.balanceOf(address(c2)) < 1 ether, "c2 not 0");
        challenge.claimChallengerWin(address(c2));
        require(token.balanceOf(address(c2)) > 1 ether, "c2 not add");
        require(token.balanceOf(address(c3)) < 1 ether, "c3 not 0");
        challenge.claimChallengerWin(address(c3));
        require(token.balanceOf(address(c3)) > 1 ether, "c3 not 0");
        //because there exists 3 nodes, so the cake pieces is 1+2+3=6,so reminding less than 6 wei.
        require(token.balanceOf(address(challenge)) < 6, "not consume out");
    }

    //2
    function testRevertWiredChallenger() public {
        init0_5();
        select(0, 5, true);
        vm.roll(100);
        challenge.proposerTimeout(DisputeTree.encodeNodeKey(0, 2));
        vm.expectRevert("you can't eat cake");
        challenge.claimChallengerWin(address(c2));
    }

    //2
    function testRevertDupEat() public {
        init0_5();
        select(0, 5, true);
        vm.roll(100);
        challenge.proposerTimeout(DisputeTree.encodeNodeKey(0, 2));
        challenge.claimChallengerWin(address(c1));
        vm.expectRevert("you can't eat cake");
        challenge.claimChallengerWin(address(c1));
    }

    //3
    function test2Branch() public {
        init0_5();
        MockChallenger old = c1;
        select(0, 5, true);
        c1 = c2;
        select(0, 5, false);
        vm.roll(100);
        c1 = old;
        IERC20 token = factory.stakingManager().token();
        challenge.proposerTimeout(DisputeTree.encodeNodeKey(0, 2));
        require(token.balanceOf(factory.dao()) == 0, "dao not 0");
        challenge.claimChallengerWin(address(c1));
        require(token.balanceOf(address(c1)) < 1 ether);
        require(token.balanceOf(factory.dao()) > 0, "dao not added");
        require(token.balanceOf(address(challenge)) == 0, "challenge contract not 0");
    }

    /**
     * @dev every challenge game is limit by a specific time point.Contract called confirmedBlock(l1).After confirmedBlock, proposer win.
     * there is 3 way game ends:proposer timeout,challenger exec success or proposer win.
     */
    function testRevertConfirmed1() public {
        vm.roll(10001);
        vm.expectRevert("block confirmed");
        challenge.proposerTimeout(0);
    }

    function testRevertConfirmed2() public {
        vm.roll(10001);
        vm.expectRevert("block confirmed");
        exec(0);
    }

    function testRevertBeforeConfirmed() public {
        vm.expectRevert("block not confirmed");
        challenge.claimProposerWin();
    }
}
