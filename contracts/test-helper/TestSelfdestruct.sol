pragma solidity ^0.8.0;

contract Test {

    uint public constant a = 10;

    address payable public owner;

    constructor(){
        owner = payable(0xEB285F24b2676d14Eb3d924E3Cc180115A5303C9);
    }

    function destruct() public {
        selfdestruct(owner);
    }

}
