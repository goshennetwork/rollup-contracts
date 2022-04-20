// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IRollupStateChain {
    /**
     * @dev Check the provided stateInfo whether confirmed, but not guarantee the correctness of stateInfo
     * @param _stateInfo State info to check.
     * @return _confirmed Whether or not the given state info is confirmed
     */
    function isStateConfirmed(Types.StateInfo memory _stateInfo) external view returns (bool _confirmed);

    /**
     * @dev Verify provided info, it checks info's index and hash
     * @param _stateInfo State info in state chain
     * @return Return true if state info is indeed in state chain
     */
    function verifyStateInfo(Types.StateInfo memory _stateInfo) external view returns (bool);

    ///emit when appendStates, anyone can check the block hash and open a challenge.
    event StateBatchAppended(
        uint64 indexed _startIndex,
        address indexed _proposer,
        uint64 _timestamp,
        bytes32[] _blockHash
    );

    /**
     * @dev Appends a list of block hash to the state chain. Only staking sender permitted
     * @param _blockHashes A list of state (we now store block hash).
     * @param _startIndex First block hash's index
     * @notice Revert if:
     * - _totalStates not equal to state chain size
     * - sender isn't staking
     * - pending states will beyond input chain size
     */
    function appendStateBatch(bytes32[] memory _blockHashes, uint64 _startIndex) external;

    event StateRollbacked(uint64 indexed _stateIndex, bytes32 indexed _blockHash);

    /**
     * @dev Cut state chain at specific state.Only called by Challenge contract
     * @param _stateInfo State info to cut the state chain
     * @notice Revert if:
     * - caller isn't Challenge contract
     * - invalid stateInfo
     * - stateInfo already confirmed;
     */
    function rollbackStateBefore(Types.StateInfo memory _stateInfo) external;

    ///get total number of submitted state
    function totalSubmittedState() external view returns (uint64);
}
