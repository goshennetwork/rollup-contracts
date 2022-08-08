// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.13;
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/ForgeVM.sol";
import "../interfaces/IL1StandardBridge.sol";
import "../bridge/L1StandardBridge.sol";
import "../bridge/L2StandardBridge.sol";
import "./TestBase.sol";
import "./MockContract.sol";
import "../cross-layer/CrossLayerContext.sol";
import "../state-machine/StateTransition.sol";
import "../challenge/Challenge.sol";
import "../challenge/ChallengeFactory.sol";

contract testUpgradeForkL2 is TestBase {
    address owner;
    address proxyAdminAddrl1;
    address proxyAdminAddrl2;
    function setUp() public {
        owner = address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266);
        proxyAdminAddrl1 = address(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512);
        proxyAdminAddrl2 = address(0xa0Ee7A142d267C1f36714E4a8F75612F20a79720);
    }
//test code (block number & rpc should update):

// forge test --fork-url http://172.168.3.73:8545 --fork-block-number 1393 
// -m "testUpgradeForkL2" -v 

/* 2.test L2code */
    function testUpgradeForkL2CrossLayerWitness() public {
        // get l2 CrossLayerWitness contract & new mockL2CrossLayerWitness contract
        address l2CrossLayerWitness = 0x2210000000000000000000000000000000000221 ;
        MockL2CrossLayerWitness newl2CrossLayerWitness = new MockL2CrossLayerWitness();

        // upgrade l2CrossLayerWitness contract
        vm.startPrank(proxyAdminAddrl2);
        TransparentUpgradeableProxy(payable(l2CrossLayerWitness)).upgradeTo(address(newl2CrossLayerWitness));
        vm.stopPrank();

        newl2CrossLayerWitness = MockL2CrossLayerWitness(l2CrossLayerWitness);
        // test call newl2CrossLayerWitness
        require(newl2CrossLayerWitness.return4() == 4, "upgrade fail 4");
    }

    function testUpgradeForkL2StandardBridge() public {
        // get l2 StandardBridge contract & new mockL2StandardBridge contract
        address L2StandardBridge = 0x2210000000000000000000000000000000000221 ;
        MockL2Bridge newL2StandardBridge = new MockL2Bridge();

        // upgrade L2StandardBridge contract
        vm.startPrank(proxyAdminAddrl2);
        TransparentUpgradeableProxy(payable(L2StandardBridge)).upgradeTo(address(newL2StandardBridge));
        vm.stopPrank();

        newL2StandardBridge = MockL2Bridge(L2StandardBridge);
        // test call newL2StandardBridge
        require(newL2StandardBridge.return2() == 2, "upgrade fail 2");
    }

    
}