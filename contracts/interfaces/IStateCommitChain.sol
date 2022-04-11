// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IStateCommitChain {
    /**
     * @dev Check the provided stateInfo whether inside fraud proof window, but not  guarantee the correctness of sateInfo
     * @param _stateInfo State info to check.
     * @return _inside Whether or not the given state info is inside the fraud proof window.
     */
    function insideFraudProofWindow(Types.StateInfo memory _stateInfo) external view returns (bool _inside);

    /**
     * @dev Verify provided info, it checkes info's index and hash
     * @param _stateInfo State info in state chain
     * @return Return true if state info is indeed in state chain
     */
    function verifyStateInfo(Types.StateInfo memory _stateInfo) external view returns (bool);

    ///emit when appendStates, anyone can check the block hash and open a challenge.
    event Appended(
        uint64 indexed _startIndex,
        bytes32[] indexed _blockHash,
        address indexed _proposer,
        uint64 _timestamp
    );

    /**
     * @dev Appends a list of block hash to the state chain.
     * @param _blockHashes A list of state (we now store block hash).
     * @param _totalStates Total states stored in state chain
     */
    function appendStates(bytes32[] memory _blockHashes, uint64 _totalStates) external;

    event Deleted(uint64 indexed _stateIndex, bytes32 indexed _blockHash);

    /**
     * @dev Cut state chain at specific state.=
     * @param _stateInfo State info to cut the state chain
     */
    function deleteState(Types.StateInfo memory _stateInfo) external;

    ///get state chain height
    function chainHeight() external view returns (uint64);
}
