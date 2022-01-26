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
    IChallenge challenge;
    MockChallenger c1;
    MockChallenger c2;
    MockChallenger c3;
    MockChallenger c4;
    MockProposer p;
    bytes32 fake32 = bytes32(uint256(0xff));

    function setUp() public {
        MockChallengeFactory factory = new MockChallengeFactory();
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
        //start
        challenge = factory.newChallengeWithProposer(address(c1), address(p));
        c1.setChallenge(challenge);
        c2.setChallenge(challenge);
        c3.setChallenge(challenge);
        c4.setChallenge(challenge);
        p.setChallenge(challenge);
        erc20.transfer(address(c1), 0x500_000_000_000);
        erc20.transfer(address(c2), 0x500_000_000_000);
        erc20.transfer(address(c3), 0x500_000_000_000);
        erc20.transfer(address(p), 0x500_000_000_000);
    }

    function init0_5() internal {
        //0->5
        p.initialize(5, fake32);
    }

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

    function revealEmptyChild(
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
        _roots[0] = 0;
        p.revealMidStates(_nodeKeys, _roots);
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
     * 7:proposer failed when timeout.
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
    function testTimeoutGameOver() public {
        vm.roll(100);
        challenge.proposerTimeout(0);
    }

    //7
    function testFailTimeoutGameOver() public {
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
        uint256[] memory _nodeKeys = new uint256[](1);
        bool[] memory _isLefts = new bool[](1);
        _nodeKeys[0] = DisputeTree.encodeNodeKey(0, 5);
        c4.selectDisputeBranch(_nodeKeys, _isLefts);
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

    function testSelect1() public {
        init0_5();
        select(0, 5, true);
    }

    /**
     * reveal task in stage 2:
     * 1:
     */

    function testReveal1() public {
        init0_5();
        select(0, 5, true);
        revealChild(0, 5, true);
    }

    function testRevertRevealEmptyState() public {
        init0_5();
        vm.expectRevert("dup new state or zero new state");
        revealEmptyChild(0, 5, true);
    }

    function testRevertRevealDup() public {
        init0_5();
        revealChild(0, 5, true);
        vm.expectRevert("dup new state or zero new state");
        revealChild(0, 5, true);
    }

    function testRevealAhead() public {
        init0_5();
        revealChild(0, 5, true);
        revealChild(0, 2, true);
        select(0, 5, true);
        select(0, 2, true);
    }

    function testRevertSelectOneStep() public {
        init0_5();
        revealChild(0, 5, true);
        revealChild(0, 2, true);
        select(0, 5, true);
        select(0, 2, true);
        vm.expectRevert("one step have no child");
        select(0, 1, true);
    }
}
