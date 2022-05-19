// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IMachineState.sol";
import "./riscv32/Register.sol";

library MemoryLayout {
    uint32 public constant HaltMagic = 0x5EAD0000;

    function readRegPC(IMachineState mem, bytes32 stateHash) internal view returns (uint32) {
        return mem.readRegister(stateHash, Register.REG_PC);
    }

    function isHalt(IMachineState mem, bytes32 stateHash) internal view returns (bool) {
        return mem.readRegister(stateHash, Register.REG_PC) == HaltMagic;
    }

    function mustReadOutput(IMachineState mem, bytes32 stateHash) internal view returns (bytes32) {
        bytes32 _out = mem.readOutput(stateHash);
        require(_out != bytes32(0), "have no output root");
        return _out;
    }
}
