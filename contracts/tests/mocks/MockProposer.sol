pragma solidity ^0.8.0;

import "../../interfaces/IChallenge.sol";
import "../../interfaces/IStakingManager.sol";
import "@openzeppelin/interfaces/IERC20.sol";

contract MockProposer {
    IChallenge challenge;
    IStakingManager sm;

    function setChallenge(IChallenge _c) external {
        challenge = _c;
    }

    function initialize(uint128 _endStep, bytes32 _midSystemState) external {
        challenge.initialize(_endStep, bytes32(uint256(0x6666)),_midSystemState);
    }

    function revealMidStates(uint256[] calldata _nodeKeys, bytes32[] calldata _stateRoots) external {
        challenge.revealMidStates(_nodeKeys, _stateRoots);
    }

    function approve(
        IERC20 _token,
        address _spender,
        uint256 _amount
    ) external {
        _token.approve(_spender, _amount);
    }

    ///staking manager action of proposer
    function setStakingManager(IStakingManager _s)external{
        sm=_s;
    }

    function deposit()public{
        sm.deposit();
    }

    function startWithdrawal()public {
        sm.startWithdrawal();
    }

    function finalizeWithdrawal()public{
        sm.finalizeWithdrawal();
    }
}
