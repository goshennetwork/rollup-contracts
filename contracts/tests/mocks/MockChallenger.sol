pragma solidity ^0.8.0;

import "../../interfaces/IChallenge.sol";
import "@openzeppelin/interfaces/IERC20.sol";
import "../../interfaces/IStakingManager.sol";


contract MockChallenger {
    IChallenge challenge;
    IStakingManager sm;

    function setChallenge(IChallenge _c) external {
        challenge = _c;
    }

    function selectDisputeBranch(uint256[] calldata _parentNodeKey, bool[] calldata _isLeft) external {
        challenge.selectDisputeBranch(_parentNodeKey, _isLeft);
    }

    function claimChallengerWin(address _challenger) external {
        challenge.claimChallengerWin(_challenger);
    }

    function approve(
        IERC20 _token,
        address _spender,
        uint256 _amount
    ) external {
        _token.approve(_spender, _amount);
    }

}
