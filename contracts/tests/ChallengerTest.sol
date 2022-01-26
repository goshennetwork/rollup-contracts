pragma solidity ^0.8.0;

import "../interfaces/IChallenge.sol";
import "../challenge/Challenge.sol";
import "./mocks/MockChallengeFactory.sol";
import "./mocks/MockStateCommitChain.sol";
import "./mocks/MockStakingManager.sol";

//import "@openzeppelin/token/ERC20/ERC20.sol";

contract TestChallenger {
    Challenge challenge1;

    function setUp() public {
        //        IChallengeFactory factory = new MockChallengeFactory();
        //        IStateCommitChain scc = new MockStateCommitChain();
        //        ERC20 erc20 = new ERC20();
        //        IStakingManager sm = new MockStakingManager(address(erc20));
        challenge1 = new Challenge();
    }
}
