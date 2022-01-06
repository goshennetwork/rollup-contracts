pragma solidity ^0.8.0;
pragma abicoder v2;

import { IAddressResolver } from "../interfaces/IAddressResolver.sol";
import { IExecutor } from "../interfaces/IExecutor.sol";
import "../interfaces/IChallenge.sol";
import "../interfaces/IChallengeFactory.sol";
import "../interfaces/IStakingManager.sol";
import { DisputeTree } from "./DisputeTree.sol";
import { IERC20 } from "@openzeppelin/interfaces/IERC20.sol";

contract Challenge is IChallenge {
    using DisputeTree for mapping(uint256 => DisputeTree.DisputeNode);

    //so the last step and 0 step's state is not in node's state root.
    mapping(uint256 => DisputeTree.DisputeNode) public disputeTree;
    //record every challenger last select node key in disputeTree.
    mapping(address => uint256) lastSelectedNodeKey;
    SystemInfo systemInfo;
    // who start challenge.
    address creator;
    //at which l1 block number, the game timeout.
    uint256 expireAfterBlock;
    //after game finished, recognize which roles win the game.
    bool isChallengerWin;
    State state;
    WithdrawStatus withdrawStatus;

    IChallengeFactory factory;
    IStakingManager public stakingManger;
    IStateCommitChain public scc;
    //run the one step
    IExecutor public executor;
    //fixme: flows need more evaluation.
    uint256 constant MinChallengerDeposit = 0.01 ether;
    //fixme: evaluate timeout more legitimate. The dispute solver can delay the challenge by provide step ((1<<256) -1),and choose deadline to repond, and responsible challenger respond in next block ,so the system judge will delay 256*(timeout+1)+timeout,if timeout is 100 this roughly 4.5 Days!
    uint256 proposerTimeLimit;

    function create(
        IAddressResolver _addressResolver,
        uint256 _blockN,
        address _proposer,
        bytes32 _systemStartState,
        bytes32 _systemEndState,
        address _creator,
        uint256 _proposerTimeLimit
    ) external override {
        factory = IChallengeFactory(msg.sender);
        stakingManger = IStakingManager(_addressResolver.stakingManager());
        executor = IExecutor(_addressResolver.executor());
        scc = IStateCommitChain(_addressResolver.scc());

        systemInfo.blockNumber = _blockN;
        systemInfo.proposer = _proposer;
        systemInfo.systemStartState = _systemStartState;
        systemInfo.systemEndState = _systemEndState;
        creator = _creator;
        proposerTimeLimit = _proposerTimeLimit;
        expireAfterBlock = block.number + proposerTimeLimit;
        state = State.Started;
        emit ChallengeStarted(_blockN, _proposer, _systemStartState, _systemEndState, expireAfterBlock);
    }

    function getHeadLeafNode() public view returns (uint256) {
        return disputeTree.getFirstLeafNode();
    }

    function initialize(
        address _sender,
        uint128 _endStep,
        bytes32 _midSystemState
    ) external override {
        //only challengeFactory.
        require(msg.sender == address(factory), "only challenge factory can initialize");
        //in start period.
        require(
            block.number <= expireAfterBlock &&
                state == State.Started &&
                _sender == systemInfo.proposer &&
                _endStep > 1, //larger than 1
            "wrong context"
        );
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

    function revealMidStates(uint256[] calldata _nodeKeys, bytes32[] calldata _stateRoots) external override {
        require(
            state == State.Running && msg.sender == proposer && block.number < interactiveDeadline(),
            "wrong context"
        );
        require(_nodeKeys.length == _stateRoots.length && _nodeKeys.length > 0, "illegal length");

        for (uint256 i = 0; i < _stateRoots.length; i++) {
            bytes32 stateRoot = _stateRoots[i];
            DisputeTree.DisputeNode storage node = disputeTree[_nodeKeys[i]];
            // 不需要提前检查节点是否存在，允许proposer提前披露状态
            require(block.number <= node.expireAfterBlock, "time out");
            require(node.midStateRoot == 0 && stateRoot != 0);
            node.midStateRoot = stateRoot;
            node.expireAfterBlock = uint256(-1);
        }
        emit MidStateRevealed(_nodeKeys, _stateRoots);
    }

    function proposerTimeout(uint256 _nodeKey) external override {
        //should inside fraud proof window.
        _requireInsideProofWindow(systemInfo.blockNumber);
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
        }
        _challengeSuccess();
        emit ProposerTimeout(_nodeKey);
    }

    function selectDisputeBranch(uint256 _parentNodeKey, bool _isLeft) external override {
        //should inside fraud proof window.
        _requireInsideProofWindow(systemInfo.blockNumber);

        require(state == State.Running, "not running");
        uint256 _expireAfterBlock = block.number + proposerTimeLimit;
        uint256 _childKey = disputeTree.addNewChild(_parentNodeKey, _isLeft, _expireAfterBlock, msg.sender);
        uint256 _lastSelect = lastSelectedNodeKey[msg.sender];
        if (_lastSelect == 0) {
            //first select, need deposit.
            stakingManger.token().transferFrom(msg.sender, address(this), MinChallengerDeposit);
        } else {
            //can only select last's child node
            require(DisputeTree.isChildNode(_lastSelect, _childKey));
        }
        challengerLastAgree[msg.sender] = _childKey;

        emit DisputeBranchSelected(_childKey, _expireAfterBlock);
    }

    function execOneStepTransition(uint256 _leafNodeKey) external {
        //inside window
        _requireInsideProofWindow(systemInfo.blockNumber);
        require(state == State.Running, "wrong context");
        (uint128 _stepLower, uint128 _stepUpper) = DisputeTree.decodeNodeKey(_leafNodeKey);
        require(disputeTree[_leafNodeKey].parent != 0 && _stepUpper == 1 + _stepLower, "not one step node");
        bytes32 _startState = _stepLower == 0
            ? systemInfo.systemStartState
            : disputeTree[DisputeTree.searchNodeWithMidStep(0, systemInfo.endStep, _stepLower)].midStateRoot;
        bytes32 _endState = _stepUpper == endStep
            ? systemEndState
            : disputeTree[DisputeTree.searchNodeWithMidStep(0, systemInfo.endStep, _stepUpper)].midStateRoot;
        bytes32 executedRoot = factory.executor().executeNextStep(_startState);
        require(executedRoot != _endState, "state transition is right");
        _challengeSuccess();
        emit OneStepTransition(_stepLower, _endState, executedRoot);
    }

    function claimProposerWin() external override {
        //challenged block confirmed
        _requirePassProofWindow();
        require(state != State.Finished, "wrong context");
        _setWinner(false);
        withdrawStatus = WithdrawStatus.Over;
        IERC20 token = factory.stakingManger().token();
        uint256 _amount = token.balanceOf(address(this));
        token.transfer(systemInfo.proposer, _amount);
        emit ProposerWin(systemInfo.proposer, _amount);
    }

    function claimChallengerWin() external override {
        require(state == State.Finished && isChallengerWin, "wrong context");
        if (withdrawStatus == WithdrawStatus.UnClaimed) {
            //if not claimed claim
            factory.stakingManager().claim(systemInfo.proposer);
            withdrawStatus = WithdrawStatus.Over;
        }
        if (systemInfo.endStep == 0) {
            //not initial.
            IERC20 token = factory.stakingManger().token();
            uint256 _amount = token.balanceOf(address(this));
            //transfer to creator.
            token().transfer(creator, _amount);
            return;
        }

        (uint256 _nodeKey, bool oneBranch) = disputeTree.getFirstLeafNode(
            DisputeTree.encodeNodeKey(0, systemInfo.endStep)
        );
        if (oneBranch) {
            _divideTheCake(_nodeKey);
        } else {
            //more than one branch
            IERC20 token = factory.stakingManger().token();
            uint256 _amount = token.balanceOf(address(this));
            //transfer to DAO
            token().transfer(factory.dao(), _amount);
        }
    }

    //divide the cake at specific branch provided lowest node address.
    function _divideTheCake(uint256 _lowestNodeKey) internal {
        require(lastSelectedNodeKey[msg.sender] != 0, "you can't eat cake");
        IERC20 token = factory.stakingManger().token();
        uint256 _cake = token.balanceOf(address(this));
        uint256 _canWithdraw;
        require(_cake > 0, "no cake");
        uint256 _correctNodeAddr = _lowestNodeKey;
        uint256 _amount = 0;
        uint256 _rootKey = DisputeTree.encodeNodeKey(0, systemInfo.endStep);
        while (_correctNodeAddr != 0) {
            //first pay back,and record the amount of gainer.
            address _gainer = disputeTree[_correctNodeAddr].challenger;
            uint256 _eat = MinChallengerDeposit;
            _cake -= _eat;
            if (_gainer == msg.sender) {
                _canWithdraw += _eat;
            }
            _amount++;
            if (_correctNodeAddr == disputeTree[_correctNodeAddr].parent) {
                //reach the root;
                break;
            }
            _correctNodeAddr = disputeTree[_correctNodeAddr].parent;
        }

        if (_amount == 1) {
            //only root node.
            if (msg.sender == disputeTree[_rootKey].challenger) {
                _canWithdraw += _cake;
                _cake = 0;
            }
        } else {
            //Now just divide remaining to pieces.Assume there are 5 gainer.so divide to 5+4+3+2+1=15.so first gainer get 5/15,second gainer get 4/15,
            //next gainer get 3/15,next gainer get 2/15,last gainer get 1/15.they eat all cake!but assume there is 256 gainer.the last gainer gain 1/32896,maybe not meet the gas cost he consumes.
            //todo: so maybe we should pay back the gas cost to gainer, and then divide the cake.
            uint256 _pieces = ((1 + _amount) * _amount) / 2;
            _correctNodeAddr = _lowestNodeKey;
            for (_correctNodeAddr != 0; ; ) {
                //first pay back,and record the amount of gainer.
                address _gainer = disputeTree[_correctNodeAddr].challenger;
                if (msg.sender == _gainer) {
                    _canWithdraw += (_amount * _cake) / _pieces;
                }
                _amount--;
                _correctNodeAddr = disputeTree[_correctNodeAddr].parent;
            }
        }
        token.transfer(msg.sender, _canWithdraw);
        lastSelectedNodeKey[msg.sender] = 0;
    }

    function _challengeSuccess() internal {
        _setWinner(true);
        factory.scc().rollbackBlockBefore(systemInfo.blockNumber);
        factory.stakingManager().slash(systemInfo.blockNumber, systemInfo.systemEndState, systemInfo.proposer);
    }

    function _setWinner(bool isChallengerWin) internal {
        state = State.Finished;
        withdrawStatus = WithdrawStatus.UnClaimed;
        isChallengerWin = isChallengerWin;
    }

    function _requireInsideProofWindow(uint256 _blockN) internal view {
        require(!scc.isBlockConfirmed(_blockN), "block confirmed");
    }

    function _requirePassProofWindow(uint256 _blockN) internal view {
        require(scc.isBlockConfirmed(_blockN), "block not confirmed");
    }
}
