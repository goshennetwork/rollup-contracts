// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IMemory.sol";

library MemoryLayout {
    uint32 public constant InputHashPos = 0x30000000;
    uint32 public constant HasOuputMagicPos = 0x30000800;
    uint32 public constant OutputHashPos = 0x30000804;
    uint32 public constant PreimageHash = 0x30001000;
    uint32 public constant ImageSize = 0x31000000;
    uint32 public constant Image = 0x31000004;
    uint32 public constant RegStartPos = 0xc0000000;
    uint32 public constant RegZeroPos = RegStartPos;
    uint32 public constant RegLRPos = RegStartPos + 0x1f * 4;
    uint32 public constant RegPcPos = RegStartPos + 0x20 * 4;

    uint32 public constant HasOuputMagic = 0x1337f00d;
    uint32 public constant HaltMagic = 0x5EAD0000;

    function readRegPC(IMemory mem, bytes32 stateHash) internal returns (uint32) {
        return mem.read(stateHash, RegPcPos);
    }

    function writeInputHash(
        IMemory mem,
        bytes32 stateHash,
        bytes32 inputHash
    ) internal returns (bytes32) {
        return mem.writeBytes32(stateHash, InputHashPos, inputHash);
    }

    function hasOutput(IMemory mem, bytes32 stateHash) internal returns (bool) {
        return mem.read(stateHash, HasOuputMagicPos) == HasOuputMagic;
    }

    function isHalt(IMemory mem, bytes32 stateHash) internal returns (bool) {
        return readRegPC(mem, stateHash) == HaltMagic;
    }

    function readOutputRoot(IMemory mem, bytes32 stateHash) internal returns (bytes32) {
        require(hasOutput(mem, stateHash));

        return mem.readBytes32(stateHash, OutputHashPos);
    }
}
