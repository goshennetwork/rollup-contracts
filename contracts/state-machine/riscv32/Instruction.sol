// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../../libraries/BytesEndian.sol";

library Instruction {
    //opcode info
    //R type
    uint32 internal constant OP_R_TYPE = 51;
    //I type
    uint32 internal constant OP_I_CSR_TYPE = 115;
    uint32 internal constant OP_I_FENCE_TYPE = 15;
    uint32 internal constant OP_I_JALR_TYPE = 103;
    uint32 internal constant OP_I_LOAD_TYPE = 3;
    uint32 internal constant OP_I_CALC_TYPE = 19;
    //S type
    uint32 internal constant OP_S_TYPE = 35;
    //B type
    uint32 internal constant OP_B_TYPE = 99;
    //U type
    uint32 internal constant OP_U_LUI_TYPE = 55; //load upper immediate
    uint32 internal constant OP_U_AUIPC_TYPE = 23; //add upper immediate PC
    //J type
    uint32 internal constant OP_J_JAL_TYPE = 111; //jump and link

    function opcode(uint32 inst) internal pure returns (uint8) {
        return uint8(inst) & 0x7f;
    }

    /**
     * R-type: [7]funct7 + [5]rs2 + [5]rs1 + [3]func3 + [5]rd + [7]opcode
     */
    function decodeRType(uint32 inst)
        internal
        pure
        returns (
            uint8 op,
            uint8 rd,
            uint8 f3,
            uint8 rs1,
            uint8 rs2,
            uint8 f7
        )
    {
        op = uint8(inst) & 0x7f;
        rd = uint8(inst >> 7) & 0x1f;
        f3 = uint8(inst >> (7 + 5)) & 0x07;
        rs1 = uint8(inst >> (7 + 5 + 3)) & 0x1f;
        rs2 = uint8(inst >> (7 + 5 + 3 + 5)) & 0x1f;
        f7 = uint8(inst >> (7 + 5 + 3 + 5 + 5));
        return (op, rd, f3, rs1, rs2, f7);
    }

    /**
     * I-type: [12]immediate[11:0] + [5]rs1 + [3]funct3 + [5]rd + [7]opcode
     */
    function decodeIType(uint32 inst)
        internal
        pure
        returns (
            uint8 op,
            uint8 rd,
            uint8 f3,
            uint8 rs1,
            uint32 imm
        )
    {
        op = uint8(inst) & 0x7f;
        rd = uint8(inst >> 7) & 0x1f;
        f3 = uint8(inst >> (7 + 5)) & 0x07;
        rs1 = uint8(inst >> (7 + 5 + 3)) & 0x1f;
        imm = inst >> (7 + 5 + 3 + 5);
        imm = uint32(int32(imm << 20) >> 20); // sign extension;
        return (op, rd, f3, rs1, imm);
    }

    /**
     * S-type: [7]imm[11:5] + [5]rs2 + [5]rs1 + [3]funct3 + [5]imm[0:4] + [7]opcode
     */
    function decodeSType(uint32 inst)
        internal
        pure
        returns (
            uint8 op,
            uint8 f3,
            uint8 rs1,
            uint8 rs2,
            uint32 imm
        )
    {
        uint8 imm1;
        uint8 imm2;
        (op, imm1, f3, rs1, rs2, imm2) = decodeRType(inst);
        imm = uint32(imm1) + (uint32(imm2) << 5);
        imm = uint32(int32(imm << 20) >> 20);
    }

    /**
     * B-type: imm[12]imm[10:5] + [5]rs2 + [5]rs1 + [3]funct3 + [5]imm[4:1]imm[11] + [7]opcode
     */
    function decodeBType(uint32 inst)
        internal
        pure
        returns (
            uint8 op,
            uint8 f3,
            uint8 rs1,
            uint8 rs2,
            uint32 imm
        )
    {
        uint8 imm1;
        uint8 imm2;
        (op, imm1, f3, rs1, rs2, imm2) = decodeRType(inst);
        uint32 bit11 = (imm1 & 1);
        uint32 imm04 = imm1 - bit11;
        uint32 bit12 = (imm2 >> 6) & 1;
        imm = (bit12 << 12) + (bit11 << 11) + ((imm2 & 0x3f) << 5) + imm04;
        imm = uint32(int32(imm << 20) >> 20);
    }

    /**
     * U-type: [20]imm[31:12] + [5]rd + [7]opcpde
     */
    function decodeUType(uint32 inst)
        internal
        pure
        returns (
            uint8 op,
            uint8 rd,
            uint32 imm
        )
    {
        op = uint8(inst) & 0x7f;
        rd = uint8(inst >> 7) & 0x1f;
        imm = (inst >> 12) << 12;
    }

    /**
     * J-type: [1]imm[20] + [10]imm[10:1] + [1]imm[11] + [8]imm[19:12] + [5]rd + [7]opcode
     */
    function decodeJType(uint32 inst)
        internal
        pure
        returns (
            uint8 op,
            uint8 rd,
            uint32 imm
        )
    {
        op = uint8(inst) & 0x7f;
        rd = uint8(inst >> 7) & 0x1f;
        uint32 imm12_19 = (inst >> 12) & 0xff;
        uint32 bit11 = (inst >> 20) & 1;
        uint32 imm1_10 = (inst >> 21) & 0x03ff;
        uint32 bit20 = inst >> 31;
        imm = (bit20 << 20) + (imm12_19 << 12) + (bit11 << 11) + (imm1_10 << 1);
        imm = uint32(int32(imm << 20) >> 20);
    }
}
