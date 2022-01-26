pragma solidity ^0.8.0;

import "../../interfaces/IStateTransition.sol";

contract MockStateTransition is IStateTransition {
    function executeNextStep(bytes32 stateHash) external override returns (bytes32 nextStateHash) {
        return bytes32(uint256(0x1234432112344321));
    }

    function generateStartState(
        uint256 blockNumber,
        bytes32 parentHash,
        bytes32 txhash,
        bytes32 coinbase,
        uint256 gasLimit,
        uint256 timestemp
    ) external override returns (bytes32) {
        return bytes32(uint256(0x0001000200030004));
    }

    function verifyFinalState(bytes32 finalState, bytes32 outputRoot) external override {}
}
