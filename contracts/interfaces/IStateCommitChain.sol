// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IStateCommitChain {
    /**
     * @dev Check the provided stateInfo whether confirmed, but not  guarantee the correctness of sateInfo
     * @param _stateInfo State info to check.
     * @return _confirmed Whether or not the given state info is confirmed
     */
    function isStateConfirmed(Types.StateInfo memory _stateInfo) external view returns (bool _confirmed);

    /**
     * @dev Verify provided info, it checkes info's index and hash
     * @param _stateInfo State info in state chain
     * @return Return true if state info is indeed in state chain
     */
    function verifyStateInfo(Types.StateInfo memory _stateInfo) external view returns (bool);

    ///emit when appendStates, anyone can check the block hash and open a challenge.
    event StateAppended(uint64 indexed _startIndex, bytes32[] _blockHash, address indexed _proposer, uint64 _timestamp);

    /**
     * @dev Appends a list of block hash to the state chain.Only staking sender permitted
     * @param _blockHashes A list of state (we now store block hash).
     * @param _totalStates Total states stored in state chain
     * @notice Revert if:
     * - _totalStates if equal to chain size
     * - sender isn't staking
     * - pending states will beyond transaction chain size(every "block" in tx chain will finally drive a "block" in state chain)
     */
    function appendStates(bytes32[] memory _blockHashes, uint64 _totalStates) external;

    event StateRolledBackBefore(uint64 indexed _stateIndex, bytes32 indexed _blockHash);

    /**
     * @dev Cut state chain at specific state.Only called by Challenge contract
     * @param _stateInfo State info to cut the state chain
     * @notice Revert if:
     * - caller isn't Challenge contract
     * - invalid stateInfo
     * - stateInfo already confirmed;
     */
    function rollbackStateBefore(Types.StateInfo memory _stateInfo) external;

    ///get state chain height
    function chainHeight() external view returns (uint64);
}
