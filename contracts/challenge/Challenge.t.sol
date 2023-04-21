// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./Challenge.sol";
import "../interfaces/ForgeVM.sol";
import "../libraries/console.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MockStateTransition {
    ///@dev validate final state, revert if final state is not halt or output inconsistent root
    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) public {}
}

contract MockStakingManager {
    ERC20 erc20;

    constructor() {
        erc20 = new ERC20("", "");
    }

    function token() public view returns (IERC20) {
        console.log(address(this));
        return IERC20(address(erc20));
    }
}

contract MockStateChain {
    function isStateConfirmed(Types.StateInfo memory _stateInfo) external view returns (bool _confirmed) {
        return false;
    }
}

contract MockChallengeFactory {
    MockStateTransition st;
    MockStakingManager sm;
    MockStateChain sc;

    constructor() {
        st = new MockStateTransition();
        sm = new MockStakingManager();
        sc = new MockStateChain();
    }

    function executor() public view returns (IStateTransition) {
        return IStateTransition(address(st));
    }

    function stakingManager() public view returns (IStakingManager) {
        return IStakingManager(address(sm));
    }

    function rollupStateChain() public view returns (IRollupStateChain) {
        return IRollupStateChain(address(sc));
    }
}

contract TestChallenge {
    bytes32 fakeBytes32 = bytes32(uint256(0xdead));
    ForgeVM public constant vm = ForgeVM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    MockChallengeFactory factory;
    Challenge challenge;
    uint256 stepNum = 1 << 32;
    uint128 N_SECTION = 7;

    address creator = address(0xd000);
    address proposer = address(0xd001);

    function setUp() public {
        factory = new MockChallengeFactory();
        vm.startPrank(address(factory));
        challenge = new Challenge();
        Types.StateInfo memory _info;
        _info.proposer = proposer;
        challenge.create(fakeBytes32, creator, 0xffff, _info, 0);
        vm.stopPrank();

        /// @dev initialize
        vm.startPrank(proposer);
        challenge.initialize(uint64(stepNum), fakeBytes32, genStates(N_SECTION - 1));
        vm.stopPrank();
    }

    function testRevealCost() public {
        vm.startPrank(creator);
        challenge.selectDisputeBranch(
            gen256List(DisputeTree.encodeNodeKey(0, uint128(stepNum))),
            gen128List(uint128(0))
        );
        vm.stopPrank();

        vm.startPrank(proposer);
        uint256 key;
        (, key) = DisputeTree.nSection(N_SECTION, 0, 0, uint128(stepNum));
        uint256[] memory _ks = gen256List(key);
        bytes32[] memory _ss = genBytes32List(fakeBytes32, N_SECTION - 1);
        uint256 _before = gasleft();
        challenge.revealSubStates(_ks, _ss);
        console.log(_before - gasleft());
        vm.stopPrank();
    }

    function gen128List(uint128 _i) internal returns (uint128[] memory) {
        uint128[] memory _info = new uint128[](1);
        _info[0] = _i;
        return _info;
    }

    function gen256List(uint256 _i) internal returns (uint256[] memory) {
        uint256[] memory _info = new uint256[](1);
        _info[0] = _i;
        return _info;
    }

    function genBytes32List(bytes32 _i) internal returns (bytes32[] memory) {
        bytes32[] memory _info = new bytes32[](1);
        _info[0] = _i;
        return _info;
    }

    function genBytes32List(bytes32 _i, uint256 num) internal returns (bytes32[] memory) {
        bytes32[] memory _info = new bytes32[](num);
        for (uint256 i = 0; i < num; i++) {
            _info[i] = _i;
        }
        return _info;
    }

    function genStates(uint256 num) internal view returns (bytes32[] memory) {
        bytes32[] memory _s = new bytes32[](num);
        for (uint256 i = 0; i < num; i++) {
            _s[i] = fakeBytes32;
        }
        return _s;
    }
}