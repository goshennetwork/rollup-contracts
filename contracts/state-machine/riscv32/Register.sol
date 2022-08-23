// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../../libraries/BytesEndian.sol";
import "../../libraries/MerkleTrie.sol";
import "../../libraries/BytesSlice.sol";

library Register {
    uint32 internal constant REGISTER_NUM = 33;

    //reg offset info
    uint32 internal constant REG_X0 = 0; //x0 hardwired 0
    uint32 internal constant REG_RA = 1; //x1 return address
    uint32 internal constant REG_SP = 2; //x2 stack pointer
    uint32 internal constant REG_GP = 3; //x3 global pointer
    uint32 internal constant REG_TP = 4; //x4 thread pointer
    //temporary
    uint32 internal constant REG_T0 = 5; //x5
    uint32 internal constant REG_T1 = 6; //x6
    uint32 internal constant REG_T2 = 7; //x7
    //saved register,frame pointer
    uint32 internal constant REG_FP = 8; //x8 or s0
    //saved register
    uint32 internal constant REG_S1 = 9; //x9
    //function arguement,return value
    uint32 internal constant REG_A0 = 10; //x10
    uint32 internal constant REG_A1 = 11; //x11
    uint32 internal constant REG_A2 = 12; //x12
    uint32 internal constant REG_A3 = 13; //x13
    uint32 internal constant REG_A4 = 14; //x14
    uint32 internal constant REG_A5 = 15; //x15
    uint32 internal constant REG_A6 = 16; //x16
    uint32 internal constant REG_A7 = 17; //x17
    //saved register
    uint32 internal constant REG_S2 = 18; //x18
    uint32 internal constant REG_S3 = 19; //x19
    uint32 internal constant REG_S4 = 20; //x20
    uint32 internal constant REG_S5 = 21; //x21
    uint32 internal constant REG_S6 = 22; //x22
    uint32 internal constant REG_S7 = 23; //x23
    uint32 internal constant REG_S8 = 24; //x24
    uint32 internal constant REG_S9 = 25; //x25
    uint32 internal constant REG_S10 = 26; //x26
    uint32 internal constant REG_S11 = 27; //x27
    //temporary
    uint32 internal constant REG_T3 = 28; //x28
    uint32 internal constant REG_T4 = 29; //x29
    uint32 internal constant REG_T5 = 30; //x30
    uint32 internal constant REG_T6 = 31; //x31
    //pc
    uint32 internal constant REG_PC = 32; //x32

    function readRegisterBytes4(
        mapping(bytes32 => HashDB.Preimage) storage hashdb,
        bytes32 root,
        uint32 regid
    ) internal view returns (bytes4) {
        if (regid == REG_X0) {
            return bytes4(0);
        }
        (bool exists, bytes memory value) = MerkleTrie.get(hashdb, genRegisterKey(regid), root);
        return exists ? BytesSlice.bytesToBytes4(value) : bytes4(0);
    }

    function readRegister(
        mapping(bytes32 => HashDB.Preimage) storage hashdb,
        bytes32 root,
        uint32 regid
    ) internal view returns (uint32) {
        bytes4 result = readRegisterBytes4(hashdb, root, regid);
        return BytesEndian.bytes4ToUint32(result);
    }

    function genRegisterKey(uint32 regid) private pure returns (bytes memory) {
        return bytes.concat(bytes1(uint8(regid)));
    }

    function writeRegister(
        mapping(bytes32 => HashDB.Preimage) storage hashdb,
        bytes32 root,
        uint32 regid,
        uint32 value
    ) internal returns (bytes32) {
        return writeRegisterBytes4(hashdb, root, regid, BytesEndian.uint32ToLEBytes(value));
    }

    function writeRegisterBytes4(
        mapping(bytes32 => HashDB.Preimage) storage hashdb,
        bytes32 root,
        uint32 regid,
        bytes4 value
    ) internal returns (bytes32) {
        if (regid == REG_X0) {
            return root;
        }
        return MerkleTrie.update(hashdb, genRegisterKey(regid), BytesSlice.bytes4ToBytes(value), root);
    }
}
