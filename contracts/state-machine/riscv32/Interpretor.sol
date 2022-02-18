// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "./Instruction.sol";
import "../MachineState.sol";
import "./Register.sol";
import "../MemoryLayout.sol";

contract Interpretor {
    MachineState public mstate;

    constructor(address state) {
        mstate = MachineState(state);
    }

    //WARNNING: this is only for testing RV32I system.
    function start(uint32 _entrypoint) public {
        mstate.writeRegister(0, Register.REG_PC, _entrypoint);
        bytes32 root = 0;
        for (bool halted = false; !halted; ) {
            (root, halted) = step(root);
        }
    }

    function step(bytes32 root) public returns (bytes32, bool) {
        uint32 currPC = mstate.readRegister(root, Register.REG_PC);
        uint32 inst = mstate.readMemory(root, currPC);
        uint8 op = Instruction.opcode(inst);
        uint32 nextPC = currPC + 4;
        if (op == Instruction.OP_R_TYPE) {
            (, uint8 rd, uint8 fn3, uint32 vrs1, uint32 vrs2, uint8 fn7) = Instruction.decodeRType(inst);
            uint256 fn = (uint256(fn3) << 8) + uint256(fn7);
            vrs1 = mstate.readRegister(root, vrs1); // reuse register id as value to avoid stack too deep error
            vrs2 = mstate.readRegister(root, vrs2);
            if (fn == (0 << 8) + 0) {
                unchecked {
                    vrs1 += vrs2;
                }
            } else if (fn == (0 << 8) + 32) {
                unchecked {
                    vrs1 -= vrs2;
                }
            } else if (fn == 1 << 8) {
                vrs1 = vrs1 << (vrs2 & 0x1f); // sll shift left logical
            } else if (fn == 2 << 8) {
                vrs1 = int32(vrs1) < int32(vrs2) ? 1 : 0; // slt set less than
            } else if (fn == 3 << 8) {
                vrs1 = (vrs1 < vrs2) ? 1 : 0; // sltu set less than unsigned
            } else if (fn == 4 << 8) {
                vrs1 = vrs1 ^ vrs2; // xor
            } else if (fn == 5 << 8) {
                vrs1 = vrs1 >> (vrs2 & 0x1f); // srl: shift right logical
            } else if (fn == (5 << 8) + 32) {
                vrs1 = uint32(int32(vrs1) >> (vrs2 & 0x1f)); // sra: shift arithmetic
            } else if (fn == 6 << 8) {
                vrs1 = vrs1 | vrs2;
            } else if (fn == 7 << 8) {
                vrs1 = vrs1 & vrs2;
            } else {
                nextPC = MemoryLayout.HaltMagic;
            }
            root = mstate.writeRegister(root, rd, vrs1);
        } else if (op == Instruction.OP_I_FENCE_TYPE) {
            // fence: nop
        } else if (op == Instruction.OP_I_CSR_TYPE) {
            (, uint8 rd, uint8 fn3, uint8 rs1, uint32 csr) = Instruction.decodeIType(inst);
            if (fn3 == 0) {
                // environment call/break
                if (csr == 0) {
                    // ecall
                    // WARNNING: TESTING
                    uint32 _a0 = mstate.readRegister(root, Register.REG_A0);
                    if (_a0 != 1) {
                        revert("failed");
                    }
                    nextPC = MemoryLayout.HaltMagic;
                } else if (csr == 1) {
                    // ebreak: nop
                } else {
                    nextPC = MemoryLayout.HaltMagic; // invalid
                }
            } else if (fn3 == 1) {
                //csrrw control status register read & write
                nextPC = MemoryLayout.HaltMagic;
            } else if (fn3 == 2) {
                //csrrs control status register read & set bit
                nextPC = MemoryLayout.HaltMagic;
            } else if (fn3 == 3) {
                //csrrc control status register read & clear bit
                nextPC = MemoryLayout.HaltMagic;
            } else if (fn3 == 5) {
                //the flow 3 instruction rs1 means zimm
                //csrrwi control status register read & write immediate
                nextPC = MemoryLayout.HaltMagic;
            } else if (fn3 == 6) {
                //csrrsi control status register read & set bit immediate
                nextPC = MemoryLayout.HaltMagic;
            } else if (fn3 == 7) {
                //csrrci control sttus register read & clear bit immediate
                nextPC = MemoryLayout.HaltMagic;
            } else {
                nextPC = MemoryLayout.HaltMagic; // invalid for RV32I
            }
        } else if (op == Instruction.OP_I_JALR_TYPE) {
            // JALR rd rs1 imm : rd = pc + 4, pc = rs1 + sext(imm), pc[0] = 0
            (, uint8 rd, uint8 fn3, uint8 rs1, uint32 imm) = Instruction.decodeIType(inst);
            uint32 vrs1 = mstate.readRegister(root, rs1);
            if (fn3 == 0) {
                root = mstate.writeRegister(root, rd, nextPC);
                unchecked {
                    vrs1 += imm;
                }
                nextPC = vrs1 & (~uint32(1)); // reset lowest bit
            } else {
                nextPC = MemoryLayout.HaltMagic;
            }
        } else if (op == Instruction.OP_I_LOAD_TYPE) {
            (, uint8 rd, uint8 fn3, uint8 rs1, uint32 imm) = Instruction.decodeIType(inst);
            uint32 vrs1 = mstate.readRegister(root, rs1);
            unchecked {
                vrs1 += imm;
            }
            if (fn3 == 0) {
                vrs1 = uint32(int32(int8(uint8(mstate.readMemoryByte(root, vrs1))))); // load byte and sign extension
            } else if (fn3 == 1) {
                bytes2 half = mstate.readMemoryBytes2(root, vrs1);
                vrs1 = uint32(int32(int16(BytesEndian.bytes2ToUint16(half)))); // load halfword and sign extension
            } else if (fn3 == 2) {
                vrs1 = mstate.readMemory(root, vrs1); // load word
            } else if (fn3 == 4) {
                vrs1 = uint32(uint8(mstate.readMemoryByte(root, vrs1))); // load byte unsigned
            } else if (fn3 == 5) {
                bytes2 half = mstate.readMemoryBytes2(root, vrs1);
                vrs1 = uint32(BytesEndian.bytes2ToUint16(half)); // load halfword unsigned
            } else {
                nextPC = MemoryLayout.HaltMagic; // invalid for rv32I
            }
            root = mstate.writeRegister(root, rd, vrs1);
        } else if (op == Instruction.OP_I_CALC_TYPE) {
            (, uint8 rd, uint8 fn3, uint8 rs1, uint32 imm) = Instruction.decodeIType(inst);
            uint32 vrs1 = mstate.readRegister(root, rs1);
            if (fn3 == 0) {
                unchecked {
                    vrs1 += imm; // addi
                }
            } else if (fn3 == 2) {
                vrs1 = (int32(vrs1) < int32(imm)) ? 1 : 0; // slti
            } else if (fn3 == 3) {
                vrs1 = (vrs1 < imm) ? 1 : 0; // sltiu
            } else if (fn3 == 4) {
                vrs1 ^= imm; // xori
            } else if (fn3 == 6) {
                vrs1 |= imm; // ori
            } else if (fn3 == 7) {
                vrs1 &= imm; // andi
            } else if (fn3 == 1) {
                // slli shift left logical immediate
                if (imm >> 5 != 0) {
                    nextPC = MemoryLayout.HaltMagic;
                }
                unchecked {
                    vrs1 <<= imm & 0x1f;
                }
            } else if (fn3 == 5) {
                //srli/srai
                uint32 shift = imm & 0x1f;
                uint32 imm7 = imm >> 5;
                if (imm7 == 0) {
                    // srli shift right logical immediate
                    vrs1 >>= shift;
                } else if (imm7 == 32) {
                    // srai shift right arithmetic immediate
                    vrs1 = uint32(int32(vrs1) >> shift);
                } else {
                    nextPC = MemoryLayout.HaltMagic; // invalid for RV32I
                }
            } else {
                nextPC = MemoryLayout.HaltMagic; // invalid for RV32I
            }
            root = mstate.writeRegister(root, rd, vrs1);
        } else if (op == Instruction.OP_S_TYPE) {
            (, uint8 fn3, uint8 rs1, uint8 rs2, uint32 imm) = Instruction.decodeSType(inst);
            uint32 vrs1 = mstate.readRegister(root, rs1);
            bytes4 vrs2 = mstate.readRegisterBytes4(root, rs2);
            unchecked {
                vrs1 += imm;
            }
            if (fn3 == 0) {
                root = mstate.writeMemoryByte(root, vrs1, bytes1(vrs2));
            } else if (fn3 == 1) {
                root = mstate.writeMemoryBytes2(root, vrs1, bytes2(vrs2));
            } else if (fn3 == 2) {
                root = mstate.writeMemoryBytes4(root, vrs1, vrs2);
            } else {
                nextPC = MemoryLayout.HaltMagic; // invalid for RV32I
            }
        } else if (op == Instruction.OP_B_TYPE) {
            (, uint8 _fn3, uint8 _rs1, uint8 _rs2, uint32 _offset) = Instruction.decodeBType(inst);
            uint32 vrs1 = mstate.readRegister(root, _rs1);
            uint32 vrs2 = mstate.readRegister(root, _rs2);
            uint32 vpc = mstate.readRegister(root, Register.REG_PC);
            if (_fn3 == 2 || _fn3 == 3) {
                nextPC = MemoryLayout.HaltMagic;
            } else if (
                (_fn3 == 0 && vrs1 == vrs2) ||
                (_fn3 == 1 && vrs1 != vrs2) ||
                (_fn3 == 4 && int32(vrs1) < int32(vrs2)) ||
                (_fn3 == 5 && int32(vrs1) >= int32(vrs2)) ||
                (_fn3 == 6 && vrs1 < vrs2) ||
                (_fn3 == 7 && vrs1 >= vrs2)
            ) {
                unchecked {
                    nextPC = vpc + _offset;
                }
            }
        } else if (op == Instruction.OP_U_LUI_TYPE) {
            (, uint8 rd, uint32 imm) = Instruction.decodeUType(inst);
            root = mstate.writeRegister(root, rd, imm);
            // LUI rd imm
        } else if (op == Instruction.OP_U_AUIPC_TYPE) {
            (, uint8 rd, uint32 imm) = Instruction.decodeUType(inst);
            root = mstate.writeRegister(root, rd, currPC + imm); // auipc rd imm
        } else if (op == Instruction.OP_J_JAL_TYPE) {
            // JAL rd imm : rd = pc + 4, pc = pc + imm
            (, uint8 rd, uint32 imm) = Instruction.decodeJType(inst);
            root = mstate.writeRegister(root, rd, nextPC);
            nextPC = currPC + imm;
        } else {
            nextPC = MemoryLayout.HaltMagic; // invalid opcode for RV32I
        }

        root = mstate.writeRegister(root, Register.REG_PC, nextPC);
        //instruction step success, increment step num
        uint32 _num = mstate.readRegister(root, Register.REG_Counter);
        root = mstate.writeRegister(root, Register.REG_Counter, _num + 1);
        return (root, nextPC == MemoryLayout.HaltMagic);
    }
}
