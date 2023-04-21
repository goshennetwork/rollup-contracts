// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

uint128 constant MidSteps = 6;

interface IChallenge {
    //the info of rv32 system info.
    struct SystemInfo {
        Types.StateInfo stateInfo;
        //systemEndState index,must > 1.
        uint128 endStep;
        //systemStartState calculated by Executor.
        bytes32 systemStartState;
        bytes32 systemEndState;
    }

    enum ChallengeStage {
        Uninitialized,
        //challenge game started.
        Started,
        // challenge game initialized,and now need to find out one step.
        Running,
        // one step find out and challenger win, or proposer win, game over.
        Finished
    }
    //challenger whether need to claim dispute proposer's staking.
    enum ClaimStatus {
        //when game finished and challenger win, they have to claim the payback.
        UnClaimed,
        //claim over.
        Over
    }

    function minChallengerDeposit() external view returns (uint256);

    function stateConfirmed() external view returns (bool);

    /**
     * @dev Create challenge by challengeFactory.guarantee the info provided true.
     * @param _systemStartState System initial state of program, calculated by executor.
     * @param _creator Challenger who start challenge.
     * @param _proposerTimeLimit After how much l1 block, the proposer expired.
     * @param _stateInfo StateInfo contains the challenged block info, already confirmed by challengeFactory
     * @param _minChallengerDeposit floor price for challenge to engage challenge game
     * @notice revert when transfer failed
     */
    function create(
        bytes32 _systemStartState,
        address _creator,
        uint256 _proposerTimeLimit,
        Types.StateInfo memory _stateInfo,
        uint256 _minChallengerDeposit
    ) external;

    event ChallengeInitialized(uint128 _systemEndStep, bytes32[MidSteps] _subStates);

    /**
     * @dev Initialize challenge info, provide endStep and sub states of root node
     * @param endStep End step index of system state of program,must larger than 1.
     * @param _systemEndState End system state, 0 is illegal, and end state must be "correct"(the program is halt, and the output is consistent with outputRoot).
     * @param _subStates sub states of root node, 0 is illegal.
     * @notice required:
     * - 1.Only stage Started
     * - 2.Only proposer
     * - 3.Proposer must initialize in time; the step nums should larger than N Section of challenge game; end system state can not be zero; sub states num is N_Section-1
     * - 4.The final state is proven properly(proposer should provide needed preimage to StateMachine)
     * - 5.The sub branch node if root node is not exist(should never happen), and the state is not provided yet(should never happen), the sub state should not be zero
     */
    function initialize(
        uint64 endStep,
        bytes32 _systemEndState,
        bytes32[MidSteps] calldata _subStates
    ) external;

    event MidStateRevealed(uint256[] nodeKeys, bytes32[MidSteps][] stateRoots);

    /**
     * @dev Proposer reveal the node's branch step state.
     * @param _parentNodeKeys The parent node which need to reveal its branch step state
     * @param _stateRoots The branch step state of  parent node,0 state root is illegal.Duplicated step state is ignored, so when parent
     * is 0-5, and branch num is 3, only need to provide state of step 1 and step 2. step 5 is surely revealed in advance
     * @notice required:
     * - 1.The challenged state is in fraud proof window
     * - 2.There is at least one parent node key, and the parent node num should less than state num
     * - 3.Parent must exist and reveal time not beyond the expire time
     * - 4.The parent node must not node be one step node
     * - 5.Can not reveal a parent node's branch state twice
     * - 6.State revealed can't be zero
     */
    function revealSubStates(uint256[] calldata _parentNodeKeys, bytes32[MidSteps][] calldata _stateRoots) external;

    event ProposerTimeout(uint256 nodeKey);

    /**
     * @dev proposer only need to do 2 kinds of action:
     * 1: initialize challenge info after challenge started.
     * 2: provide step state when challenge running.
     * if one of these actions timeout, anyone can call this function to end challenge game
     * @param _nodeKey If proposer timeout doing 1st situation, node key can easily set to 0
     * Otherwise, in 2nd situation, _nodeKey is the disputeNode key which proposer didn't reveal its branch in expireTime
     */
    function proposerTimeout(uint256 _nodeKey) external;

    event DisputeBranchSelected(address indexed challenger, uint256[] nodeKey, uint256 expireAfterBlock);

    /**
     * @dev Anyone has deposited in this challengeGame can select one branch in dispute tree.which means selected dispute
     * nodes' start system state is right,and the nodes' end system state is wrong.
     * @param _parentNodeKeys The parent node key in disputeTree.When a dispute branch is selected, it must derived from exist larger disputeNode, we call it parent node
     * i.e. A node present 0->4 stateTransition, B node present 0->2 stateTransition,0->2 is driven by 0->4,so A is parent node.
     * @param _Nth the chosen node is index number of parent branch, if parent node is divided into 5 branch, the last branch is index 4.
     * @notice  required:
     * - 1.The challenged state is in fraud windows.
     * - 2.Parent node should not be one step node, one step can't be divided
     * - 3.A challenger can only select one branch, and challenger need to deposit at challenge contract first
     * - 4.Parent node should exist and has been revealed by provider,(if node is revealed by provider, the sub branch node will be created by provider).
     */
    function selectDisputeBranch(uint256[] calldata _parentNodeKeys, uint128[] calldata _Nth) external;

    event OneStepTransition(uint256 startStep, bytes32 revealedRoot, bytes32 executedRoot);

    /**
     * @dev Anyone can verify oneStepRuns in challenge period.if the result system state of one step is not same as proposer claimed. challenger win roll back the block and game over.It ensures:
     * 1: in interactive period.
     * 2: proposer is wrong.
     * 3: roll back the wrong block
     * @param _leafNodeKey is one step node key in dispute tree.
     * @notice Revert if node is not oneStepNode.
     */
    function execOneStepTransition(uint256 _leafNodeKey) external;

    event ProposerWin(address _winner, uint256 _amount);

    /**
     * @dev The only way can proposer win.This allows unfinished challenge game over after passing the challenge time.
     */
    function claimProposerWin() external;

    /**
     * @dev Challenger can get reward by this way if there exist only one branch.Otherwise, now just transfer to DAO.
     * @param _challenger Which challenger claims the payback
     * @param _stateInfo StateInfo to provide the New stateInfo
     */
    function claimChallengerWin(address _challenger, Types.StateInfo memory _stateInfo) external;

    /**
     * @dev get last selected node key of challenger
     * @return the node key of specific challenger
     * @notice the creator of challenger also recorded in here.
     */
    function lastSelectedNodeKey(address _challenger) external view returns (uint256);

    /**
     * @dev return challenge system info
     * @return stateInfo state info of challenged state
     * @return endStep system end step index
     * @return systemStartState system start state
     * @return systemEndState system end state
     */
    function systemInfo()
        external
        view
        returns (
            Types.StateInfo memory stateInfo,
            uint128 endStep,
            bytes32 systemStartState,
            bytes32 systemEndState
        );

    /**
     * @dev get dispute tree of specific node key
     * @param _nodeKey the node key of dispute node in tree
     * @return parent the node key of parent node,if this node is not exist, parent is zero
     * @return challenger the challenger who open this node
     * @return expireAfterBlock if l1 block number larger than this, then this node is expired
     */
    function disputeTree(uint256 _nodeKey)
        external
        view
        returns (
            uint256 parent,
            address challenger,
            uint256 expireAfterBlock
        );

    /// @return stage of this challenge
    function stage() external view returns (ChallengeStage);

    /// @return claim status of this challenge
    function claimStatus() external view returns (ClaimStatus);

    /// @return in stage 1, the challenge expired after this block.
    function expireAfterBlock() external view returns (uint256);

    /// @return ask provided challenger can claim the cake or not
    function canClaimTheCake(address _challenger) external view returns (bool);
}
