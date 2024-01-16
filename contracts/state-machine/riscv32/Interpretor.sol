// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "./Instruction.sol";
import "../MachineState.sol";
import "./Register.sol";
import "../MemoryLayout.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "./Syscall.sol";
import "../../libraries/console.sol";
import "../../interfaces/IInterpretor.sol";

contract Interpretor is IInterpretor, Initializable {
    MachineState public mstate;

    function initialize(address state) public initializer {
        mstate = MachineState(state);
    }

    //WARNNING: this is only for testing RV32I system.
    function start(bytes32 _root, uint32 _entrypoint) public returns (bytes32, uint32, uint32) {
        _root = mstate.writeRegister(_root, Register.REG_PC, _entrypoint);
        uint32 _i;
        uint32 inst;
        for (bool halted = false; !halted;) {
            _i++;
            uint32 _pc = mstate.readRegister(_root, Register.REG_PC);
            if (_pc & 3 != 0) {
                bytes memory _b = abi.encodePacked("invalid pc, last inst is: ", Strings.toHexString(inst));
                revert(string(_b));
            }
            inst = mstate.readMemory(_root, _pc);
            (_root, halted) = step(_root);
        }
        return (_root, _i, inst);
    }

    function step(bytes32 root) public returns (bytes32, bool) {
        uint32 currPC = mstate.readRegister(root, Register.REG_PC);
        if (currPC == MemoryLayout.HaltMagic) {
            //already halt
            return (root, true);
        }
        uint32 nextPC = currPC + 4;
        uint32 inst = mstate.readMemory(root, currPC);
        uint8 op = Instruction.opcode(inst);
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
            } else if (fn == (0 << 8) + 1) {
                //mul 把寄存器x[rs2]乘到寄存器x[rs1]上，乘积写入 x[rd]。忽略算术溢出
                unchecked {
                    vrs1 = vrs1 * vrs2;
                }
            } else if (fn == (1 << 8) + 1) {
                //mulh 把寄存器 x[rs2]乘到寄存器x[rs1]上，都视为2的补码，将乘积的高位写入x[rd]
                unchecked {
                    vrs1 = uint32(uint64((int64(int32(vrs1)) * int64(int32(vrs2)))) >> 32);
                }
            } else if (fn == (2 << 8) + 1) {
                //mulhsu 把寄存器 x[rs2]乘到寄存器 x[rs1]上，x[rs1]为2的补码，x[rs2]为无符号数，将乘积的高位写入x[rd]。
                unchecked {
                    vrs1 = uint32(uint64((int64(int32(vrs1)) * int64(uint64(vrs2)))) >> 32);
                }
            } else if (fn == (3 << 8) + 1) {
                //mulhu 把寄存器x[rs2]乘到寄存器x[rs1]上，x[rs1]、x[rs2]均为无符号数，将乘积的高位写入x[rd]
                unchecked {
                    vrs1 = uint32((uint64(vrs1) * uint64(vrs2)) >> 32);
                }
            } else if (fn == (4 << 8) + 1) {
                //div 用寄存器x[rs1]的值除以寄存器x[rs2]的值，向零舍入，将这些数视为二进制补码，把商写入x[rd],软件层面检查除数为0的情况
                unchecked {
                    //ignore overflow
                    vrs1 = vrs2 == 0 ? uint32((1 << 32) - 1) : uint32(int32(vrs1) / int32(vrs2));
                }
            } else if (fn == (5 << 8) + 1) {
                //divu 用寄存器x[rs1]的值除以寄存器x[rs2]的值，向零舍入，将这些数视为无符号数，把商写入x[rd]
                vrs1 = vrs2 == 0 ? uint32((1 << 32) - 1) : vrs1 / vrs2;
            } else if (fn == (6 << 8) + 1) {
                //rem x[rs1]除以 x[rs2]，向0舍入，都视为2的补码，余数写入x[rd]
                vrs1 = vrs2 == 0 ? vrs1 : uint32(int32(vrs1) % int32(vrs2));
            } else if (fn == (7 << 8) + 1) {
                //remu x[rs1]除以x[rs2]，向0舍入，都视为无符号数，余数写入x[rd]
                vrs1 = vrs2 == 0 ? vrs1 : vrs1 % vrs2;
            } else {
                nextPC = MemoryLayout.HaltMagic;
            }
            root = mstate.writeRegister(root, rd, vrs1);
        } else if (op == Instruction.OP_I_FENCE_TYPE) {
            // fence: nop
        } else if (op == Instruction.OP_I_CSR_TYPE) {
            (,, uint8 fn3,, uint32 csr) = Instruction.decodeIType(inst);
            if (fn3 == 0) {
                // environment call/break
                if (csr == 0) {
                    //call
                    return handleSyscall(root, nextPC);
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
                (_fn3 == 0 && vrs1 == vrs2) || (_fn3 == 1 && vrs1 != vrs2) || (_fn3 == 4 && int32(vrs1) < int32(vrs2))
                    || (_fn3 == 5 && int32(vrs1) >= int32(vrs2)) || (_fn3 == 6 && vrs1 < vrs2)
                    || (_fn3 == 7 && vrs1 >= vrs2)
            ) {
                unchecked {
                    nextPC = vpc + _offset;
                }
            }
        } else if (op == Instruction.OP_U_LUI_TYPE) {
            (, uint8 rd, uint32 imm) = Instruction.decodeUType(inst);
            root = mstate.writeRegister(root, rd, imm); // LUI rd imm
        } else if (op == Instruction.OP_U_AUIPC_TYPE) {
            (, uint8 rd, uint32 imm) = Instruction.decodeUType(inst);
            unchecked {
                root = mstate.writeRegister(root, rd, currPC + imm); // auipc rd imm
            }
        } else if (op == Instruction.OP_J_JAL_TYPE) {
            // JAL rd imm : rd = pc + 4, pc = pc + imm
            (, uint8 rd, uint32 imm) = Instruction.decodeJType(inst);
            root = mstate.writeRegister(root, rd, nextPC);
            unchecked {
                nextPC = currPC + imm;
            }
        } else if (op == Instruction.OP_M_TYPE) {
            (root, nextPC) = handleAmo(root, nextPC, inst);
        } else {
            nextPC = MemoryLayout.HaltMagic; // invalid opcode for RV32I
        }

        root = mstate.writeRegister(root, Register.REG_PC, nextPC);
        return (root, nextPC == MemoryLayout.HaltMagic);
    }

    function handleSyscall(bytes32 _root, uint32 _nextPC) internal returns (bytes32, bool) {
        uint32 _systemNumer = mstate.readRegister(_root, Register.REG_A7);
        uint32 va0 = mstate.readRegister(_root, Register.REG_A0);
        if (_systemNumer == 0) {
            //pub fn input(hash: *mut u8);
            //get input hash, a0 put returned addr pos;write output in addr.
            _root = mstate.writeMemoryBytes32(_root, va0, mstate.readInput(_root));
        } else if (_systemNumer == 1) {
            //pub fn ret(hash: *const u8) -> !;
            //return, the program is over, a0 put state addr in memory.
            _root = mstate.writeOutput(_root, mstate.readMemoryBytes32(_root, va0));
            _nextPC = MemoryLayout.HaltMagic;
        } else if (_systemNumer == 2) {
            //pub fn preimage_len(hash: *const u8) -> usize
            //get preimage len, a0 put hash addr in memory;write out length in a0.
            bytes32 _hash = mstate.readMemoryBytes32(_root, va0);
            bytes memory _data = mstate.preimage(_hash);
            _root = mstate.writeRegister(_root, Register.REG_A0, uint32(_data.length));
        } else if (_systemNumer == 3) {
            //pub fn preimage_at(hash: *const u8, offset: usize) -> u32;
            //get preimage's 4 bytes at specific offset, a0 put hash addr, a1 put length of preimage;write out preimage in a0.
            bytes32 _hash = mstate.readMemoryBytes32(_root, va0);
            uint32 va1 = mstate.readRegister(_root, Register.REG_A1);
            uint32 data = mstate.preimageAt(_hash, va1);
            _root = mstate.writeRegister(_root, Register.REG_A0, data);
        } else if (_systemNumer == 4) {
            //pub fn panic(msg: *const u8, len: usize) -> !;
            //panic,a0 put the panic info start addr, a1 put length.program halt
            uint32 va1 = mstate.readRegister(_root, Register.REG_A1);
            revert(mstate.readMemoryString(_root, va0, va1));
        } else if (_systemNumer == 5) {
            //pub fn debug (msg: *const u8, len: usize);
            //debug,a0 put the debug info, a1 put the length.
            uint32 va1 = mstate.readRegister(_root, Register.REG_A1);
            console.logString(mstate.readMemoryString(_root, va0, va1));
        } else if (_systemNumer == 6) {
            /// hash, r, s: [u8;32], v: 0 or 1
            /// result: [u8;20]
            //pub fn ecrecover(result: *mut u8, hash: *const u8, r: *const u8, s: *const u8, v: u32)
            bytes32 hash = mstate.readMemoryBytes32(_root, mstate.readRegister(_root, Register.REG_A1));
            bytes32 r = mstate.readMemoryBytes32(_root, mstate.readRegister(_root, Register.REG_A2));
            bytes32 s = mstate.readMemoryBytes32(_root, mstate.readRegister(_root, Register.REG_A3));
            uint32 v = mstate.readRegister(_root, Register.REG_A4);
            address signer = ecrecover(hash, uint8(v + 27), r, s);
            _root = mstate.writeMemoryAddr(_root, va0, signer);
        } else {
            //invalid sys num
            _nextPC = MemoryLayout.HaltMagic;
        }
        _root = mstate.writeRegister(_root, Register.REG_PC, _nextPC);
        return (_root, _nextPC == MemoryLayout.HaltMagic);
    }

    function handleAmo(bytes32 root, uint32 nextPC, uint32 inst) internal returns (bytes32, uint32) {
        (, uint8 rd, uint8 fn3, uint32 rs1, uint32 rs2, uint8 fn7) = Instruction.decodeRType(inst);
        if (fn3 != 2) {
            nextPC = MemoryLayout.HaltMagic;
        }
        uint32 vrs1 = mstate.readRegister(root, rs1);
        uint32 t = mstate.readMemory(root, vrs1);
        uint32 vrs2 = mstate.readRegister(root, rs2);
        fn7 = fn7 >> 2;
        uint32 result = t;
        if (fn7 == 2 && rs2 == 0) {
            //lr.w 从内存中地址为 x[rs1]中加载四个字节，符号位扩展后写入 x[rd]，并对这个内存字注册保留。
            root = mstate.reserve(root, vrs1);
        } else if (fn7 == 3) {
            //sc.w 内存地址 x[rs1]上存在加载保留，将 x[rs2]寄存器中的 4 字节数存入该地址。如果存入成功， 向寄存器 x[rd]中存入 0，否则存入一个非 0 的错误码
            result = mstate.isReserved(root, vrs1) ? 0 : 1;
            if (result == 0) {
                t = vrs2;
                root = mstate.unReserve(root);
            }
        } else if (fn7 == 1) {
            //amoswap.w : rd = M[rs1]; swap(rd, rs2); M[rs1] = rd
            t = vrs2;
        } else if (fn7 == 0) {
            //amoadd.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t+x[rs2]，把 x[rd] 设为符号位扩展的 t
            unchecked {
                t = t + vrs2;
            }
        } else if (fn7 == 4) {
            //amoxor.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t 和 x[rs2]按位异 或的结果，把 x[rd]设为符号位扩展的 t。
            t = t ^ vrs2;
        } else if (fn7 == 12) {
            //amoand.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t 和 x[rs2]位与的 结果，把 x[rd]设为符号位扩展的 t
            t = t & vrs2;
        } else if (fn7 == 8) {
            //amoor.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t 和 x[rs2]位或的 结果，把 x[rd]设为符号位扩展的 t
            t = t | vrs2;
        } else if (fn7 == 16) {
            //amomin.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t 和 x[rs2]中较小 的一个（用二进制补码比较），把 x[rd]设为符号位扩展的 t
            t = int32(t) <= int32(vrs2) ? t : vrs2;
        } else if (fn7 == 20) {
            //amomax.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t 和 x[rs2]中较大的一个（用二进制补码比较），把 x[rd]设为符号位扩展的 t
            t = int32(t) >= int32(vrs2) ? t : vrs2;
        } else if (fn7 == 24) {
            //amominu.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t 和 x[rs2]中较小 的一个（用无符号比较），把 x[rd]设为符号位扩展的 t
            t = t <= vrs2 ? t : vrs2;
        } else if (fn7 == 28) {
            //amomaxu.w 将内存中地址为 x[rs1]中的字记为 t，把这个字变为 t 和 x[rs2]中 较大的一个（用无符号比较），把 x[rd]设为 t
            t = t >= vrs2 ? t : vrs2;
        } else {
            nextPC = MemoryLayout.HaltMagic;
        }
        root = mstate.writeMemory(root, vrs1, t);
        root = mstate.writeRegister(root, rd, result);
        return (root, nextPC);
    }
}
