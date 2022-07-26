// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.13;
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract MockL1Bridge is Initializable{
    function return1() public pure returns(uint){
        return 1;
    }
}

contract MockL2Bridge is Initializable{
    function return2() public pure returns(uint){
        return 2;
    }
}

contract MockL1CrossLayerWitness is Initializable{
    function return3() public pure returns(uint){
        return 3;
    }
}

contract MockL2CrossLayerWitness is Initializable{
    function return4() public pure returns(uint){
        return 4;
    }
}

contract MockDAO is Initializable{
    function return5() public pure returns(uint){
        return 5;
    }
}

contract MockAddressManager is Initializable{
    function return6() public pure returns(uint){
        return 6;
    }
}

contract MockChainStorageContainer is Initializable{
    function return7() public pure returns(uint){
        return 7;
    }
}

contract MockRollupInputChain is Initializable{
    function return8() public pure returns(uint){
        return 8;
    }
}

contract MockRollupStateChain is Initializable{
    function return9() public pure returns(uint){
        return 9;
    }
}

contract MockStakingManager is Initializable{
    function return10() public pure returns(uint){
        return 10;
    }
}

contract MockStateTransition is Initializable{
    function return11() public pure returns(uint){
        return 11;
    }
}