// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../interfaces/IChallenge.sol";
import "../interfaces/IChallengeFactory.sol";
import "./DisputeTree.sol";
import "../libraries/UnsafeMath.sol";

import "@openzeppelin/contracts/interfaces/IERC20.sol";

contract Challenge is IChallenge {
    using DisputeTree for mapping(uint256 => DisputeTree.DisputeNode);

    mapping(uint128 => bytes32) public stepState; /// @dev step number => state
    IChallengeFactory public factory;
    uint256 public override minChallengerDeposit;

    //so the last step and 0 step's state is not in node's state root.
    mapping(uint256 => DisputeTree.DisputeNode) public override disputeTree;
    //record every challenger last select node key in disputeTree.
    mapping(address => uint256) public override lastSelectedNodeKey;
    SystemInfo public override systemInfo;
    ChallengeStage public override stage;
    ClaimStatus public override claimStatus;
    // who start challenge.
    address creator;
    //at which l1 block number, the game timeout.
    uint256 public override expireAfterBlock;
    //fixme: evaluate timeout more legitimate. The dispute solver can delay the challenge by provide step ((1<<256) -1),and choose deadline to repond, and responsible challenger respond in next block ,so the system judge will delay 256*(timeout+1)+timeout,if timeout is 100 this roughly 4.5 Days!
    uint256 proposerTimeLimit;
    //amount challenge get from dispute proposer.
    uint256 rewardAmount;

    /** challenge game have 3 stage now:
     * stage1: game started by challenger, proposer need to init game info.
     * stage2: proposer and challengers find out one step and challenger prove this one step is wrong.
     * stage3: challenge game over.Now challenger have to claim out the payback(proposer get reward immediately when game over, but challenger have to wait to claim).
     * note: in stage1&stage2, proposer can make challenge game "stuck"(not participate in time).
     */
    modifier stage1() {
        require(stage == ChallengeStage.Started, "only started stage");
        _;
    }

    modifier stage2() {
        require(stage == ChallengeStage.Running, "only running stage");
        _;
    }

    modifier stage3() {
        require(stage == ChallengeStage.Finished, "only finished stage");
        _;
    }

    modifier beforeBlockConfirmed() {
        require(!stateConfirmed(), "block confirmed");
        _;
    }

    modifier afterBlockConfirmed() {
        require(stateConfirmed(), "block not confirmed");
        _;
    }

    modifier onlyProposer() {
        require(msg.sender == systemInfo.stateInfo.proposer, "only proposer");
        _;
    }

    function stateConfirmed() public view override returns (bool) {
        return factory.rollupStateChain().isStateConfirmed(systemInfo.stateInfo);
    }

    //when create, creator should deposit at this contract.
    function create(
        bytes32 _systemStartState,
        address _creator,
        uint256 _proposerTimeLimit,
        Types.StateInfo memory _stateInfo,
        uint256 _minChallengerDeposit
    ) external override {
        require(stage == ChallengeStage.Uninitialized, "initialized");
        factory = IChallengeFactory(msg.sender);
        systemInfo.systemStartState = _systemStartState;
        stepState[0] = _systemStartState;
        creator = _creator;
        proposerTimeLimit = _proposerTimeLimit;
        expireAfterBlock = block.number + proposerTimeLimit;
        systemInfo.stateInfo = _stateInfo;
        minChallengerDeposit = _minChallengerDeposit;
        //started
        stage = ChallengeStage.Started;
        //emit by challengeFactory
        //emit ChallengeStarted(_blockN, _proposer, _systemStartState, _systemEndState, expireAfterBlock);
    }

    function initialize(
        uint64 endStep,
        bytes32 _systemEndState,
        bytes32[MidSteps] calldata _subStates
    ) external override stage1 onlyProposer {
        uint128 _endStep = endStep;
        //in start period.
        /// @notice if system's step is less than N_SECTION, the proposer will always lose in challenge, but that's ok because
        /// when it happens, the system must exist a huge bug
        require(block.number <= expireAfterBlock && _endStep > MidSteps + 1 && _systemEndState != 0, "wrong context");
        systemInfo.systemEndState = _systemEndState;
        systemInfo.endStep = _endStep;
        stepState[_endStep] = _systemEndState;
        factory.executor().verifyFinalState(_systemEndState, systemInfo.stateInfo.blockHash);

        uint256 rootKey = DisputeTree.encodeNodeKey(0, _endStep);
        disputeTree[rootKey] = DisputeTree.DisputeNode({
            parent: rootKey, // we take node.parent != 0 as initialized node. so set root's parent pointer to self
            challenger: creator,
            expireAfterBlock: block.number + proposerTimeLimit
        });
        /// @dev the last step is equal to end system state
        uint128 _stepLower = 0;
        for (uint128 i = 0; i < MidSteps; i += 1) {
            bytes32 _stateRoot = _subStates[i];
            require(_stateRoot != 0, "wrong state root");

            uint128 _stepUpper = DisputeTree.midStep(MidSteps, i, 0, endStep);
            stepState[_stepUpper] = _stateRoot;
            uint256 _tempNodeKey = DisputeTree.encodeNodeKey(_stepLower, _stepUpper);
            disputeTree[_tempNodeKey].parent = rootKey;
            _stepLower = _stepUpper;
        }
        disputeTree[DisputeTree.encodeNodeKey(_stepLower, _endStep)].parent = rootKey;

        lastSelectedNodeKey[creator] = rootKey;
        stage = ChallengeStage.Running;
        //notify challengers in this game.
        emit ChallengeInitialized(_endStep, _subStates);
    }

    function revealSubStates(uint256[] calldata _nodeKeys, bytes32[MidSteps][] calldata _stateRoots)
        external
        override
        beforeBlockConfirmed
        stage2
        onlyProposer
    {
        require(_nodeKeys.length == _stateRoots.length && _nodeKeys.length > 0, "illegal length");
        for (uint256 i = 0; i < _nodeKeys.length; i = UnsafeMath.unsafeIncrement(i)) {
            uint256 _nodeKey = _nodeKeys[i];
            uint256 expireBlock = disputeTree[_nodeKey].expireAfterBlock;
            require(
                disputeTree[_nodeKey].parent != 0 && expireBlock > block.number && expireBlock < type(uint128).max,
                "empty parent or expired or already revealed"
            );
            (uint128 _stepStart, uint128 _stepEnd) = DisputeTree.decodeNodeKey(_nodeKey);
            uint128 _stepLower = _stepStart;
            for (uint128 j = 0; j < MidSteps; j += 1) {
                /// @dev change tempNSection if remained step num less than N_SECTION
                uint128 _stepUpper = DisputeTree.midStep(MidSteps, j, _stepStart, _stepEnd);
                if (_stepUpper != _stepLower) {
                    uint256 _tempNodeKey = DisputeTree.encodeNodeKey(_stepLower, _stepUpper);
                    disputeTree[_tempNodeKey].parent = _nodeKey;

                    bytes32 _stateRoot = _stateRoots[i][j];
                    require(_stateRoot != 0, "wrong state root");
                    stepState[_stepUpper] = _stateRoot;
                    _stepLower = _stepUpper;
                }
            }
            disputeTree[DisputeTree.encodeNodeKey(_stepLower, _stepEnd)].parent = _nodeKey;
            // mark this node will not expire any more
            disputeTree[_nodeKey].expireAfterBlock += type(uint128).max;
        }

        emit MidStateRevealed(_nodeKeys, _stateRoots);
    }

    function proposerTimeout(uint256 _nodeKey) external override beforeBlockConfirmed {
        if (stage == ChallengeStage.Finished) {
            //challenge game finished, just return.
            return;
        } else if (stage == ChallengeStage.Started) {
            //proposer initialize timeout.
            require(block.number > expireAfterBlock, "initialize challenge info not timeout");
        } else if (stage == ChallengeStage.Running) {
            //proposer reveal timeout
            /// @notice make sure the node is exist and need to be revealed
            DisputeTree.DisputeNode memory _node = disputeTree[_nodeKey];
            require(_node.parent != 0 && _node.expireAfterBlock > 0, "no such node");
            require(block.number > _node.expireAfterBlock, "not timeout yet");
        }
        _challengeSuccess();
        emit ProposerTimeout(_nodeKey);
    }

    function selectDisputeBranch(uint256[] calldata _parentNodeKeys, uint128[] calldata _Nth)
        external
        override
        beforeBlockConfirmed
        stage2
    {
        require(_parentNodeKeys.length > 0 && _parentNodeKeys.length == _Nth.length, "inconsistent length");
        uint256[] memory _childKeys = new uint256[](_parentNodeKeys.length);
        uint256 _expireAfterBlock = block.number + proposerTimeLimit;
        for (uint256 i = 0; i < _parentNodeKeys.length; i = UnsafeMath.unsafeIncrement(i)) {
            uint256 _parentNodeKey = _parentNodeKeys[i];
            uint128 _n = _Nth[i];
            uint256 _childKey = disputeTree.addNewChild(
                MidSteps + 1,
                _n,
                _parentNodeKey,
                _expireAfterBlock,
                msg.sender
            );
            uint256 _lastSelect = lastSelectedNodeKey[msg.sender];
            if (_lastSelect == 0) {
                // first select, need deposit.
                IERC20 depositToken = factory.stakingManager().token();
                require(depositToken.transferFrom(msg.sender, address(this), minChallengerDeposit));
            } else {
                //can only select last's child node
                require(DisputeTree.isChildNode(_lastSelect, _childKey));
            }
            lastSelectedNodeKey[msg.sender] = _childKey;
            _childKeys[i] = _childKey;
        }
        emit DisputeBranchSelected(msg.sender, _childKeys, _expireAfterBlock);
    }

    function execOneStepTransition(uint256 _leafNodeKey) external beforeBlockConfirmed stage2 {
        (uint128 _stepLower, uint128 _stepUpper) = DisputeTree.decodeNodeKey(_leafNodeKey);
        require(disputeTree[_leafNodeKey].parent != 0 && _stepUpper == 1 + _stepLower, "not one step node");
        bytes32 _startState = stepState[_stepLower];
        bytes32 _endState = stepState[_stepUpper];
        require(_startState != 0 && _endState != 0, "not provided yet");

        bytes32 executedRoot = factory.executor().executeNextStep(_startState);
        require(executedRoot != _endState, "state transition is right");
        _challengeSuccess();
        emit OneStepTransition(_stepLower, _endState, executedRoot);
    }

    function claimProposerWin() external override afterBlockConfirmed stage2 {
        _proposerSuccess();
        IERC20 token = factory.stakingManager().token();
        uint256 _amount = token.balanceOf(address(this));
        address _proposer = systemInfo.stateInfo.proposer;
        require(token.transfer(_proposer, _amount), "transfer failed");
        emit ProposerWin(_proposer, _amount);
    }

    function claimChallengerWin(address _challenger, Types.StateInfo memory _stateInfo) external override stage3 {
        IERC20 token = factory.stakingManager().token();
        if (claimStatus == ClaimStatus.UnClaimed) {
            //if not claimed, then claim
            uint256 _before = token.balanceOf(address(this));
            factory.stakingManager().claim(systemInfo.stateInfo.proposer, _stateInfo);
            uint256 _now = token.balanceOf(address(this));
            rewardAmount = _now - _before;
            //claim over
            claimStatus = ClaimStatus.Over;
        }
        if (systemInfo.endStep == 0) {
            //not initial.
            uint256 _amount = token.balanceOf(address(this));
            //transfer to creator.
            require(token.transfer(creator, _amount), "transfer failed");
            return;
        }

        uint256 _rootKey = DisputeTree.encodeNodeKey(0, systemInfo.endStep);
        (uint256 _nodeKey, uint64 _depth, bool oneBranch) = disputeTree.getFirstLeafNode(MidSteps + 1, _rootKey);
        if (oneBranch) {
            _divideTheCake(_nodeKey, _depth, _challenger, token);
        } else {
            //more than one branch
            uint256 _amount = token.balanceOf(address(this));
            //transfer to DAO
            require(token.transfer(factory.dao(), _amount), "transfer failed");
        }
    }

    //divide the cake at specific branch provided lowest node address.
    function _divideTheCake(
        uint256 _lowestNodeKey,
        uint64 _depth,
        address _challenger,
        IERC20 token
    ) internal {
        require(lastSelectedNodeKey[_challenger] != 0, "you can't eat cake");
        require(rewardAmount > 0, "no cake");
        uint256 _canWithdraw = minChallengerDeposit;
        uint64 _amount = _depth;
        //pay back deposit
        // vi = (i+k) / [n*(n+1)/2 + nk] , k = 10, n = 50, v0 = 10/(25*51+ 500) = 1/355, vn/v0 = 6
        uint256 _scale;
        uint256 _k = 10;
        uint256 _pieces = (((1 + _amount) * _amount) / 2) + (_amount * _k);
        uint256 _correctNodeKey = _lowestNodeKey;
        while (_correctNodeKey != 0) {
            DisputeTree.DisputeNode storage node = disputeTree[_correctNodeKey];
            //first pay back,and record the amount of gainer.
            if (_challenger == node.challenger) {
                _scale += (_amount + _k);
            }
            _amount--;
            if (node.parent == _correctNodeKey) {
                //reach the root
                break;
            }
            _correctNodeKey = node.parent;
        }
        _canWithdraw += (_scale * rewardAmount) / _pieces;
        lastSelectedNodeKey[_challenger] = 0;
        require(token.transfer(_challenger, _canWithdraw), "transfer failed");
    }

    function canClaimTheCake(address _challenger) public view returns (bool) {
        return lastSelectedNodeKey[_challenger] != 0;
    }

    //finish game and rollback the dispute l2 block & slash the dispute proposer.
    function _challengeSuccess() internal {
        stage = ChallengeStage.Finished;
        factory.rollupStateChain().rollbackStateBefore(systemInfo.stateInfo);
        factory.stakingManager().slash(
            systemInfo.stateInfo.index,
            systemInfo.stateInfo.blockHash,
            systemInfo.stateInfo.proposer
        );
    }

    function _proposerSuccess() internal {
        stage = ChallengeStage.Finished;
        //only challenger need to claim
        claimStatus = ClaimStatus.Over;
    }
}
