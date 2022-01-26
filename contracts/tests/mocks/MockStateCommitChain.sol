pragma solidity ^0.8.0;
import "../../interfaces/IStateCommitChain.sol";

contract MockStateCommitChain is IStateCommitChain {
    function getCurrentBlockHeight() external view override returns (uint256) {
        return 1;
    }

    function isBlockConfirmed(uint256 blockHeight) external view override returns (bool) {
        return false;
    }

    function getBlockInfo(uint256 blockHeight)
        external
        view
        override
        returns (
            bytes32 blockHash,
            bytes32 root,
            address proposer,
            uint256 timestamp,
            uint256 confirmedAfterBlock
        )
    {
        return (bytes32(uint256(0x3112ddaad)), bytes32(uint256(0xd3d33d3d)), address(0), 0, 10000);
    }

    function rollbackBlockBefore(uint256 fraultBlock) external override {}
}
