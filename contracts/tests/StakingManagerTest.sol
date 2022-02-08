pragma solidity ^0.8.0;

import "../staking/StakingManager.sol";
import "./mocks/MockChallengeFactory.sol";
import "./mocks/MockStateCommitChain.sol";
import "./mocks/MockProposer.sol";
import "./mocks/MockERC20.sol";
import "./mocks/MockStateTransition.sol";



interface Vm {
    // Set block.height (newHeight)
    function roll(uint256) external;

    function expectRevert(bytes calldata c) external;
}

contract StakingManagerTest {
    Vm vm = Vm(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);
    address fakeDao = address(0x0123);
    StakingManager sm;

    MockProposer p1;
    MockProposer p2;

    function setUp() public {
        MockChallengeFactory factory = new MockChallengeFactory();
        IStateCommitChain scc = new MockStateCommitChain();
        ERC20 erc20 = new MockERC20();
        sm = new StakingManager(fakeDao, address(factory), address(scc), address(erc20), 2 ether);

        IStateTransition executor = new MockStateTransition();
        factory.init(address(sm), address(executor), address(scc));
        p1 = new MockProposer();
        p2 = new MockProposer();
        p1.setStakingManager(sm);
        p2.setStakingManager(sm);

        erc20.transfer(address(p1),2 ether);
        p1.approve(erc20,address(sm),4 ether);
    }

    //normal deposit
    function testDeposit()public{
        p1.deposit();
    }
    //invalid: deposit without balance
    function testFailDepositWithoutBalance()public{
        p2.deposit();
    }
    //deposit should change the staking state
    function testIsStaking()public{
        require(!sm.isStaking(address(p1)),"staking");
        p1.deposit();
        require(sm.isStaking(address(p1)),"not staking");
    }
    //normal withdraw
    function testWithdraw()public{
        p1.deposit();
        p1.startWithdrawal();
        require(!sm.isStaking(address(p1)),"staking");
    }
    //invalid: withdraw without deposit in advance
    function testFailWithdrawWithoutStaking()public{
        p1.startWithdrawal();
    }
    //invalid: withdraw duplicated
    function testFailWithdrawDup()public{
        p1.deposit();
        p1.startWithdrawal();
        p1.startWithdrawal();
    }
    //normal finalize after startWithdraw and the block have confirmed.
    function testFinalize()public{
        p1.deposit();
        p1.startWithdrawal();
        //make block confirmed.
        vm.roll(1);
        p1.finalizeWithdrawal();
        //new deposit is permitted.
        p1.deposit();
    }
    //invalid: finalize without deposit.
    function testFailFinalizeWithoutDeposit()public{
        p1.finalizeWithdrawal();
    }
    //invalid: finalize without start withdraw first
    function testFailFinalizeWithoutWithdraw()public{
        p1.deposit();
        p1.finalizeWithdrawal();
    }
    //invalid: finalize before block confirmed
    function testFailFinalizeWithoutConfirmed()public{
        p1.deposit();
        p1.startWithdrawal();
        p1.deposit();
    }
    ///what's more, proposer can be challenged to be fraudulent when in staking state or in withdrawing state.

    //normal slash by challenge contract.(mock challenge factory always think every address is challenge contract,so directly invoke staking manager)
    function testSlash()public{
        p1.deposit();
        //now there is only once.
        sm.slash(1,bytes32(uint256(0x12345678)),address(p1));
        require(!sm.isStaking(address(p1)),"staking");
    }

    //invalid: slash proposer that unStaked, which means the slashed proposer have not deposited yet.
    function testRevertSlashEmpty()public{
        vm.expectRevert("unStaked unexpected");
        sm.slash(1,0,address(p1));
    }
    //invalid: can't withdraw anymore after being slashed
    function testRevertFinalizeAfterSlashed()public{
        p1.deposit();
        p1.startWithdrawal();
        sm.slash(1,0,address(p1));
        //block confirmed
        vm.roll(1);
        vm.expectRevert("not in withdrawing");
        p1.finalizeWithdrawal();
    }

    //normal claim when slashed block is confirmed
    function testClaim()public{
        p1.deposit();
        //now there is only once.
        sm.slash(1,bytes32(uint256(0x12345678)),address(p1));
        //now block is confirmed
        vm.roll(1);
        sm.claim(address(p1));
    }

    //invalid: should claim when slashed block is confirmed
    function testFailClaimWithoutConfirmed()public{
        p1.deposit();
        sm.slash(1,bytes32(uint256(0x12345678)),address(p1));
        sm.claim(address(p1));
    }
    //invalid: only claim when slashing
    function testFailClaimWithoutSlash()public{
        sm.claim(address(p1));
    }
    //invalid: slash same end systemRoot means the slashed block is useless.So both proposer and challenger is fraudulent
    function testRevertClaimUnusedSlash()public{
        p1.deposit();
        sm.slash(1,MockStateCommitChain(address(sm.scc())).root(),address(p1));
        vm.roll(1);
        vm.expectRevert("unused challenge");
        sm.claim(address(p1));
    }














}
