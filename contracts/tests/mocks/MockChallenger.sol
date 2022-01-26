pragma solidity ^0.8.0;

import "../../interfaces/IChallenge.sol";

contract MockChallenger {
    IChallenge challenge;

    function setChallenge(IChallenge _c) external {
        challenge = _c;
    }

    function selectDisputeBranch(uint256[] calldata _parentNodeKey, bool[] calldata _isLeft) external {
        challenge.selectDisputeBranch(_parentNodeKey, _isLeft);
    }

    function claimChallengerWin() external {
        challenge.claimChallengerWin();
    }
}
