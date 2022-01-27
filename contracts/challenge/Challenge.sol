pragma solidity ^0.8.0;

import "../interfaces/IChallenge.sol";
import "../interfaces/IChallengeFactory.sol";
import "./DisputeTree.sol";
import { IERC20 } from "@openzeppelin/interfaces/IERC20.sol";

contract Challenge is IChallenge {
    using DisputeTree for mapping(uint256 => DisputeTree.DisputeNode);

    IChallengeFactory factory;
    //fixme: flows need more evaluation.
    uint256 public constant override minChallengerDeposit = 0.01 ether;

    //so the last step and 0 step's state is not in node's state root.
    mapping(uint256 => DisputeTree.DisputeNode) public disputeTree;
    //record every challenger last select node key in disputeTree.
    mapping(address => uint256) lastSelectedNodeKey;
    SystemInfo systemInfo;
    State state;
    ClaimStatus claimStatus;
    // who start challenge.
    address creator;
    //at which l1 block number, the game timeout.
    uint256 expireAfterBlock;
    //fixme: evaluate timeout more legitimate. The dispute solver can delay the challenge by provide step ((1<<256) -1),and choose deadline to repond, and responsible challenger respond in next block ,so the system judge will delay 256*(timeout+1)+timeout,if timeout is 100 this roughly 4.5 Days!
    uint256 proposerTimeLimit;
    //record whether the challenged block can process.
    uint256 confirmedBlock;
    //amount challenge get from dispute proposer.
    uint256 rewardAmount;

    /** challenge game have 3 stage now:
     * stage1: game started by challenger, proposer need to init game info.
     * stage2: now proposer and multi challengers have to work together to find out one step and challenger prove this one step is wrong.
     * stage3: challenge game over.Now challenger have to claim out the payback(proposer get reward immediately when game over, but challenger have to wait to claim).
     * note: in stage1&stage2, proposer can make challenge game "stuck"(not participate in time).
     */
    modifier stage1() {
        require(state == State.Started, "only started stage");
        _;
    }

    modifier stage2() {
        require(state == State.Running, "only running stage");
        _;
    }

    modifier stage3() {
        require(state == State.Finished, "only finished stage");
        _;
    }

    modifier beforeBlockConfirmed() {
        require(block.number <= confirmedBlock, "block confirmed");
        _;
    }

    modifier afterBlockConfirmed() {
        require(block.number > confirmedBlock, "block not confirmed");
        _;
    }

    //when create, creator should deposit at this contract.
    function create(
        uint256 _blockN,
        address _proposer,
        bytes32 _systemStartState,
        bytes32 _systemEndState,
        bytes32 _outputRoot,
        address _creator,
        uint256 _proposerTimeLimit
    ) external override {
        factory = IChallengeFactory(msg.sender);
        systemInfo.blockNumber = _blockN;
        systemInfo.proposer = _proposer;
        systemInfo.systemStartState = _systemStartState;
        systemInfo.systemEndState = _systemEndState;
        systemInfo.outputRoot = _outputRoot;
        creator = _creator;
        expireAfterBlock = block.number + proposerTimeLimit;
        proposerTimeLimit = _proposerTimeLimit;
        (, , , , confirmedBlock) = factory.scc().getBlockInfo(_blockN);
        //started
        state = State.Started;
        //emit by challengeFactory
        //emit ChallengeStarted(_blockN, _proposer, _systemStartState, _systemEndState, expireAfterBlock);
    }

    function initialize(uint128 _endStep, bytes32 _midSystemState) external override stage1 {
        //in start period.
        require(
            block.number <= expireAfterBlock && msg.sender == systemInfo.proposer && _endStep > 1, //larger than 1
            "wrong context"
        );
        factory.executor().verifyFinalState(systemInfo.systemEndState, systemInfo.outputRoot);
        require(_midSystemState != 0, "0 system state root is illegal");

        systemInfo.endStep = _endStep;
        uint256 rootKey = DisputeTree.encodeNodeKey(0, _endStep);
        disputeTree[rootKey] = DisputeTree.DisputeNode({
            parent: rootKey, // we take node.parent != 0 as initialized node. so set root's parent pointer to self
            challenger: creator,
            expireAfterBlock: block.number + proposerTimeLimit,
            midStateRoot: _midSystemState
        });

        lastSelectedNodeKey[creator] = rootKey;
        state = State.Running;
        //notify other's participate in this game.
        emit ChallengeInitialized(_endStep, _midSystemState);
    }

    function revealMidStates(uint256[] calldata _nodeKeys, bytes32[] calldata _stateRoots)
        external
        override
        beforeBlockConfirmed
        stage2
    {
        require(msg.sender == systemInfo.proposer, "only proposer");
        require(_nodeKeys.length == _stateRoots.length && _nodeKeys.length > 0, "illegal length");

        for (uint256 i = 0; i < _stateRoots.length; i++) {
            bytes32 stateRoot = _stateRoots[i];
            DisputeTree.DisputeNode storage node = disputeTree[_nodeKeys[i]];
            // 不需要提前检查节点是否存在，允许proposer提前披露状态。减少交互次数
            if (node.parent != 0) {
                //not exist
                require(block.number <= node.expireAfterBlock, "time out");
            }
            require(node.midStateRoot == 0 && stateRoot != 0, "wrong state root");
            node.midStateRoot = stateRoot;
        }
        emit MidStateRevealed(_nodeKeys, _stateRoots);
    }

    function proposerTimeout(uint256 _nodeKey) external override beforeBlockConfirmed {
        if (state == State.Finished) {
            //challenge game finished, just return.
            return;
        } else if (state == State.Started) {
            //proposer initialize timeout.
            require(block.number > expireAfterBlock, "initialize challenge info not timeout");
        } else if (state == State.Running) {
            //proposer reveal timeout
            DisputeTree.DisputeNode storage nextNode = disputeTree[_nodeKey];
            (uint128 stepLower, uint128 stepUpper) = DisputeTree.decodeNodeKey(_nodeKey);
            require(nextNode.parent > 0 && stepUpper > stepLower + 1, "one step don't need to prove");
            require(block.number > nextNode.expireAfterBlock, "report mid state not timeout");
            require(nextNode.midStateRoot == 0, "mid state root is revealed");
        }
        _challengeSuccess();
        emit ProposerTimeout(_nodeKey);
    }

    function selectDisputeBranch(uint256[] calldata _parentNodeKeys, bool[] calldata _isLefts)
        external
        override
        beforeBlockConfirmed
        stage2
    {
        require(_parentNodeKeys.length > 0 && _parentNodeKeys.length == _isLefts.length, "inconsistent length");
        uint256[] memory _childKeys = new uint256[](_parentNodeKeys.length);
        uint256 _expireAfterBlock = block.number + proposerTimeLimit;
        for (uint256 i = 0; i < _parentNodeKeys.length; i++) {
            uint256 _parentNodeKey = _parentNodeKeys[i];
            bool _isLeft = _isLefts[i];
            uint256 _childKey = disputeTree.addNewChild(_parentNodeKey, _isLeft, _expireAfterBlock, msg.sender);
            uint256 _lastSelect = lastSelectedNodeKey[msg.sender];
            if (_lastSelect == 0) {
                //first select, need deposit.
                factory.stakingManager().token().transferFrom(msg.sender, address(this), minChallengerDeposit);
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
        bytes32 _startState = _stepLower == 0
            ? systemInfo.systemStartState
            : disputeTree[DisputeTree.searchNodeWithMidStep(0, systemInfo.endStep, _stepLower)].midStateRoot;
        bytes32 _endState = _stepUpper == systemInfo.endStep
            ? systemInfo.systemEndState
            : disputeTree[DisputeTree.searchNodeWithMidStep(0, systemInfo.endStep, _stepUpper)].midStateRoot;
        bytes32 executedRoot = factory.executor().executeNextStep(_startState);
        require(executedRoot != _endState, "state transition is right");
        _challengeSuccess();
        emit OneStepTransition(_stepLower, _endState, executedRoot);
    }

    function claimProposerWin() external override afterBlockConfirmed stage2 {
        _proposerSuccess();
        IERC20 token = factory.stakingManager().token();
        uint256 _amount = token.balanceOf(address(this));
        //todo: maybe burn some amount
        token.transfer(systemInfo.proposer, _amount);
        emit ProposerWin(systemInfo.proposer, _amount);
    }

    //if unclaimed, claim and
    function claimChallengerWin(address _challenger) external override stage3 {
        if (claimStatus == ClaimStatus.UnClaimed) {
            //if not claimed, then claim
            uint256 _before = factory.stakingManager().token().balanceOf(address(this));
            factory.stakingManager().claim(systemInfo.proposer);
            uint256 _now = factory.stakingManager().token().balanceOf(address(this));
            rewardAmount = _now - _before;
            //claim over
            claimStatus = ClaimStatus.Over;
        }
        if (systemInfo.endStep == 0) {
            //not initial.
            IERC20 token = factory.stakingManager().token();
            uint256 _amount = token.balanceOf(address(this));
            //transfer to creator.
            token.transfer(creator, _amount);
            return;
        }

        (uint256 _nodeKey, bool oneBranch) = disputeTree.getFirstLeafNode(
            DisputeTree.encodeNodeKey(0, systemInfo.endStep)
        );
        if (oneBranch) {
            _divideTheCake(_nodeKey, _challenger);
        } else {
            //more than one branch
            IERC20 token = factory.stakingManager().token();
            uint256 _amount = token.balanceOf(address(this));
            //transfer to DAO
            token.transfer(factory.dao(), _amount);
        }
    }

    //divide the cake at specific branch provided lowest node address.
    function _divideTheCake(uint256 _lowestNodeKey, address _challenger) internal {
        require(lastSelectedNodeKey[_challenger] != 0, "you can't eat cake");
        IERC20 token = factory.stakingManager().token();
        require(rewardAmount > 0, "no cake");
        uint256 _canWithdraw;
        uint256 _correctNodeAddr = _lowestNodeKey;
        uint256 _amount = 0;
        uint256 _rootKey = DisputeTree.encodeNodeKey(0, systemInfo.endStep);
        bool haveDeposited;
        while (_correctNodeAddr != 0) {
            DisputeTree.DisputeNode storage node = disputeTree[_correctNodeAddr];
            //pay back challenger's deposit
            address _gainer = node.challenger;
            if (_gainer == _challenger) {
                //only pay back once,because challenger can select different nodes in one branch.
                haveDeposited = true;
            }
            _amount++;
            if (_correctNodeAddr == node.parent) {
                //reach the root;
                break;
            }
            _correctNodeAddr = node.parent;
        }
        if (haveDeposited) {
            //pay back
            _canWithdraw += minChallengerDeposit;
        }
        if (_amount == 1) {
            //only root node.pay all reward to it.
            if (_challenger == disputeTree[_rootKey].challenger) {
                _canWithdraw += rewardAmount;
            }
        } else {
            //Now just divide remaining to pieces.Assume there are 5 gainer.so divide to 5+4+3+2+1=15.so first gainer get 5/15,second gainer get 4/15,
            //next gainer get 3/15,next gainer get 2/15,last gainer get 1/15.they eat all cake!but assume there is 256 gainer.the last gainer gain 1/32896,maybe not meet the gas cost he consumes.
            //todo: so maybe we should pay back the gas cost to gainer, and then divide the cake.
            uint256 _pieces = ((1 + _amount) * _amount) / 2;
            _correctNodeAddr = _lowestNodeKey;
            while (_correctNodeAddr != 0) {
                DisputeTree.DisputeNode storage node = disputeTree[_correctNodeAddr];
                //first pay back,and record the amount of gainer.
                if (_challenger == node.challenger) {
                    _canWithdraw += (_amount * rewardAmount) / _pieces;
                }
                _amount--;
                if (node.parent == _correctNodeAddr) {
                    //reach the root
                    break;
                }
                _correctNodeAddr = node.parent;
            }
        }
        token.transfer(_challenger, _canWithdraw);
        lastSelectedNodeKey[_challenger] = 0;
    }

    //finish game and rollback the dispute l2 block & slash the dispute proposer.
    function _challengeSuccess() internal {
        _setWinner(true);
        factory.scc().rollbackBlockBefore(systemInfo.blockNumber);
        factory.stakingManager().slash(systemInfo.blockNumber, systemInfo.systemEndState, systemInfo.proposer);
    }

    function _proposerSuccess() internal {
        _setWinner(false);
    }

    function _setWinner(bool isChallengerWin) internal {
        state = State.Finished;
        if (!isChallengerWin) {
            //only challenger need to claim
            claimStatus = ClaimStatus.Over;
        }
    }
}
