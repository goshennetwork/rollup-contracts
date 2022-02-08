pragma solidity ^0.8.0;

import "../../interfaces/IStateCommitChain.sol";

contract MockStateCommitChain is IStateCommitChain {
    bytes32 constant public blockHash=bytes32(uint256(0x3112ddaad));
    bytes32 constant public root=bytes32(uint256(0xd3d33d3d));
    address constant public proposer=address(0);
    uint256 constant public timestamp=0;
    uint256 constant public confirmedAfterBlock=10000;

    function getCurrentBlockHeight() external view override returns (uint256) {
        return 1;
    }

    function isBlockConfirmed(uint256 blockHeight) external view override returns (bool) {
        return block.number > 0;
    }

    function getBlockInfo(uint256 blockHeight)
    external
    view
    override
    returns (
        bytes32 _blockHash,
        bytes32 _root,
        address _proposer,
        uint256 _timestamp,
        uint256 _confirmedAfterBlock
    )
    {
        return (blockHash, root, proposer, timestamp, confirmedAfterBlock);
    }

    function rollbackBlockBefore(uint256 fraultBlock) external override {}
}
