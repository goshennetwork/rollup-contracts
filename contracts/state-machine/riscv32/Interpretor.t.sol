// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

//import "../../libraries/console.sol";
import "../MachineState.sol";
import "./Interpretor.sol";
import "./Register.sol";
import "../../libraries/MerkleTrie.sol";

contract InterpretorTest {
    Interpretor exec;
    MachineState public mstate;
    bytes32 root;
    int32 constant minInt = int32(uint32(1 << 31)); //min int

    function setUp() public {
        mstate = new MachineState();
        exec = new Interpretor(address(mstate));
        root = MerkleTrie.KECCAK256_RLP_NULL_BYTES;
    }

    struct ExpectReg {
        uint32 register;
        int32 value;
    }
    struct ExpectMem {
        uint32 ptr;
        int32 value;
    }

    function resetPC() internal {
        root = mstate.writeRegister(root, Register.REG_PC, 0);
    }

    function checkRegState(ExpectReg memory e) internal returns (bool) {
        uint32 value = mstate.readRegister(root, e.register);
        return value == uint32(e.value);
    }

    function checkMemState(ExpectMem memory e) internal returns (bool) {
        uint32 value = mstate.readMemory(root, e.ptr);
        return value == uint32(e.value);
    }

    function checkInstruction(
        string memory raw,
        uint32 inst,
        ExpectReg memory e1,
        ExpectReg memory e2
    ) internal {
        execInstruction(raw, inst);
        require(checkRegState(e1), raw);
        require(checkRegState(e2), raw);
    }

    function execInstruction(string memory raw, uint32 inst) internal {
        //console.log(string(bytes.concat("exec: ", bytes(raw))));
        uint32 pc = mstate.readRegister(root, Register.REG_PC);
        root = mstate.writeMemory(root, pc, inst);
        (root, ) = exec.step(root);
    }

    function checkInstruction(
        string memory raw,
        uint32 inst,
        ExpectReg memory e
    ) internal {
        execInstruction(raw, inst);
        require(checkRegState(e), raw);
    }

    function checkInstruction(
        string memory raw,
        uint32 inst,
        ExpectMem memory e
    ) internal {
        execInstruction(raw, inst);
        require(checkMemState(e), raw);
    }

    function initRegister() public {
        for (uint32 i = 0; (i < 33); i++) {
            mstate.writeRegister(root, Register.REG_PC, i);
        }
    }

    function testExecInst() public {
        checkInstruction("li      a0,1", 0x00100513, ExpectReg(Register.REG_A0, 0 + 1)); //a0=0+1
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 0 + 1)); //a1=0+1
        checkInstruction("add     a0,a0,a1", 0x00b50533, ExpectReg(Register.REG_A0, 1 + 1)); //a0=1+1
        checkInstruction("addi    t0,t0,-97", 0xf9f28293, ExpectReg(Register.REG_T0, 0 - 97)); //t0=0 + (-97)
        checkInstruction("addi  t0, t0, 98", 0x06228293, ExpectReg(Register.REG_T0, -97 + 98)); //t0=(-97) + 98
        checkInstruction("sub     a0,a0,a1", 0x40b50533, ExpectReg(Register.REG_A0, 2 - 1)); //a0=2-1
        checkInstruction("sll     a0,a0,a1", 0x00b51533, ExpectReg(Register.REG_A0, 1 << 1)); //a0=1<<1
        checkInstruction("slt     a0,a0,a1", 0x00b52533, ExpectReg(Register.REG_A0, 0)); //postite 2>1
        checkInstruction("li      a2,-1", 0xfff00613, ExpectReg(Register.REG_A2, -1)); //a2=-1
        checkInstruction("li      a3,-2", 0xffe00693, ExpectReg(Register.REG_A3, -2)); //a3=-2
        checkInstruction("slt     a0,a2,a3", 0x00d62533, ExpectReg(Register.REG_A0, 0)); //-1>-2
        checkInstruction("sltu    a0,a2,a3", 0x00d63533, ExpectReg(Register.REG_A0, 0)); // ((1<<32)-1) > ((1<<32)-2)
        checkInstruction("li      a2,1", 0x00100613, ExpectReg(Register.REG_A2, 1)); //a2=1
        checkInstruction("sltu    a0,a2,a3", 0x00d63533, ExpectReg(Register.REG_A0, 1)); //1<((1<<32)-2)

        //xor
        checkInstruction("li      a0,1", 0x00100513, ExpectReg(Register.REG_A0, 1)); //a0=1
        checkInstruction("li      a1,2", 0x00200593, ExpectReg(Register.REG_A1, 2)); //a1=2
        checkInstruction("xor     a0,a0,a1", 0x00b54533, ExpectReg(Register.REG_A0, 3)); //a0= 0b_01 ^ 0b_10 = 0b_11
        checkInstruction("li      a0,-1", 0xfff00513, ExpectReg(Register.REG_A0, -1)); //a0= ((1<<32)-1)=0xff_ff_ff_ff
        checkInstruction("xor     a0,a0,a1", 0x00b54533, ExpectReg(Register.REG_A0, -1 - 2)); //a0=((1<<32)-1-2)=-3

        //srl
        checkInstruction("li      a0,1024", 0x40000513, ExpectReg(Register.REG_A0, 1024)); //a0=1024
        checkInstruction("li      a1,3", 0x00300593, ExpectReg(Register.REG_A1, 3)); //a1=3
        checkInstruction("srl     a0,a0,a1", 0x00b55533, ExpectReg(Register.REG_A0, 1 << 7));
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("srl     a0,a0,a1", 0x00b55533, ExpectReg(Register.REG_A0, 0)); //(1<<7)>>0xff
        checkInstruction("li      a0,-32", 0xfe000513, ExpectReg(Register.REG_A0, -32));
        checkInstruction("li      a1,3", 0x00300593, ExpectReg(Register.REG_A1, 3));
        checkInstruction("srl     a0,a0,a1", 0x00b55533, ExpectReg(Register.REG_A0, ((1 << 32) - 32) >> 3)); //((1<<32)-32)>>3

        //sra
        checkInstruction("li      a0,1024", 0x40000513, ExpectReg(Register.REG_A0, 1024)); //a0=1024
        checkInstruction("li      a1,3", 0x00300593, ExpectReg(Register.REG_A1, 3)); //a1=3
        checkInstruction("sra     a0,a0,a1", 0x40b55533, ExpectReg(Register.REG_A0, 1 << 7));
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("sra     a0,a0,a1", 0x40b55533, ExpectReg(Register.REG_A0, 0)); //(1<<7)>>0xff
        checkInstruction("li      a0,-32", 0xfe000513, ExpectReg(Register.REG_A0, -32));
        checkInstruction("li      a1,3", 0x00300593, ExpectReg(Register.REG_A1, 3));
        checkInstruction("sra     a0,a0,a1", 0x40b55533, ExpectReg(Register.REG_A0, -4)); //正常除法 -32/8
        checkInstruction("sra     a0,a0,a1", 0x40b55533, ExpectReg(Register.REG_A0, -1)); //-1

        //or
        checkInstruction("li      a0,1", 0x00100513, ExpectReg(Register.REG_A0, 1)); //a0=1
        checkInstruction("li      a1,2", 0x00200593, ExpectReg(Register.REG_A1, 2)); //a1=2
        checkInstruction("or      a0,a0,a1", 0x00b56533, ExpectReg(Register.REG_A0, 3)); //a0=0b_01 | 0b_10 = 0b_11
        checkInstruction("li      a0,-1", 0xfff00513, ExpectReg(Register.REG_A0, -1));
        checkInstruction("or      a0,a0,a1", 0x00b56533, ExpectReg(Register.REG_A0, -1));

        //and
        checkInstruction("li      a0,1", 0x00100513, ExpectReg(Register.REG_A0, 1)); //a0=1
        checkInstruction("li      a1,2", 0x00200593, ExpectReg(Register.REG_A1, 2)); //a1=2
        checkInstruction("and     a0,a0,a1", 0x00b57533, ExpectReg(Register.REG_A0, 0)); //a0=0b_01 & 0b_10=0b_00
        checkInstruction("li      a0,-1", 0xfff00513, ExpectReg(Register.REG_A0, -1));
        checkInstruction("and     a0,a0,a1", 0x00b57533, ExpectReg(Register.REG_A0, 2)); //a0=0xff_ff_ff_ff & 0b_10=0b_10

        //jalr
        checkInstruction("addi    a0,x0,5", 0x00500513, ExpectReg(Register.REG_A0, 5));
        resetPC();
        checkInstruction(
            "jalr    a1,8(a0);",
            0x008505e7,
            ExpectReg(Register.REG_A1, 4),
            ExpectReg(Register.REG_PC, 12)
        ); //omit 0b1

        //sb, sh, sw, lb, lh, lw,lbu,lhu
        checkInstruction("li      a0,-2", 0xffe00513, ExpectReg(Register.REG_A0, -2)); //0xff_ff_ff_fe
        checkInstruction("li      a1,1000", 0x3e800593, ExpectReg(Register.REG_A1, 1000));
        checkInstruction("sb      a0,0(a1)", 0x00a58023, ExpectMem(1000, -2 & ((1 << 8) - 1))); //0xfe
        checkInstruction("lb      a3,0(a1)", 0x00058683, ExpectReg(Register.REG_A3, -2)); //a3=-2
        checkInstruction("sh      a0,4(a1)", 0x00a59223, ExpectMem(1004, -2 & ((1 << 16) - 1))); //0xff_fe
        checkInstruction("lh      a3,4(a1)", 0x00459683, ExpectReg(Register.REG_A3, -2));
        checkInstruction("sw      a0,8(a1)", 0x00a5a423, ExpectMem(1008, -2)); //0xff_ff_ff_fe
        checkInstruction("lw      a3,8(a1)", 0x0085a683, ExpectReg(Register.REG_A3, -2));
        checkInstruction("lbu     a3,0(a1)", 0x0005c683, ExpectReg(Register.REG_A3, 0xfe));
        checkInstruction("lhu     a3,4(a1)", 0x0045d683, ExpectReg(Register.REG_A3, 0xfffe));

        //slti
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("slti    a0,a1,-2", 0xffe5a513, ExpectReg(Register.REG_A0, 0)); //1>-2
        checkInstruction("slti    a0,a1,2", 0x0025a513, ExpectReg(Register.REG_A0, 1)); //1<2
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("slti    a0,a1,-2", 0xffe5a513, ExpectReg(Register.REG_A0, 0)); //-1 > -2
        checkInstruction("slti    a0,a1,2", 0x0025a513, ExpectReg(Register.REG_A0, 1)); // -1 < 2

        //sltiu
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("sltiu   a0,a1,-2", 0xffe5b513, ExpectReg(Register.REG_A0, 1)); // 0b1<0xff_ff_ff_fe
        checkInstruction("sltiu   a0,a1,2", 0x0025b513, ExpectReg(Register.REG_A0, 1)); //1 < 2
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("sltiu   a0,a1,-2", 0xffe5b513, ExpectReg(Register.REG_A0, 0)); // 0xff_ff_ff_ff > 0xff_ff_ff_fe
        checkInstruction("sltiu   a0,a1,2", 0x0025b513, ExpectReg(Register.REG_A0, 0)); // 0xff_ff_ff_ff > 2

        //xori
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("xori    a0,a1,-1", 0xfff5c513, ExpectReg(Register.REG_A0, -2)); //0b1 ^ 0xff_ff_ff_ff=0x_ff_ff_ff_fe
        checkInstruction("xori    a0,a1,1", 0x0015c513, ExpectReg(Register.REG_A0, 0)); //0b1 ^ 0b1 =0
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("xori    a0,a1,-1", 0xfff5c513, ExpectReg(Register.REG_A0, 0)); //0xff_ff_ff_ff ^ 0xff_ff_ff_ff =0
        checkInstruction("xori    a0,a1,1", 0x0015c513, ExpectReg(Register.REG_A0, -2)); //0xff_ff_ff_ff ^ 0b1 =0xff_ff_ff_fe

        //ori
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("ori     a0,a1,-2", 0xffe5e513, ExpectReg(Register.REG_A0, -1)); //0b1 | 0xff_ff_ff_fe = 0xff_ff_ff_ff
        checkInstruction("ori     a0,a1,2", 0x0025e513, ExpectReg(Register.REG_A0, 3)); //0b_01 | 0b_10= 0b_11=3
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("ori     a0,a1,-2", 0xffe5e513, ExpectReg(Register.REG_A0, -1)); //0xff_ff_ff_ff | 0xff_ff_ff_fe=0xff_ff_ff_ff
        checkInstruction("ori     a0,a1,2", 0x0025e513, ExpectReg(Register.REG_A0, -1)); //0xff_ff_ff_ff | 0b_10 = 0xff_ff_ff_ff

        //andi
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("andi    a0,a1,2", 0x0025f513, ExpectReg(Register.REG_A0, 0)); //0b_01 & 0b_10=0
        checkInstruction("andi    a0,a1,-2", 0xffe5f513, ExpectReg(Register.REG_A0, 0)); //0b_01 & 0x_ff_ff_fe = 0
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("andi    a0,a1,2", 0x0025f513, ExpectReg(Register.REG_A0, 2)); //0xff_ff_ff_ff & 0b_10= 0b_10
        checkInstruction("andi    a0,a1,-2", 0xffe5f513, ExpectReg(Register.REG_A0, -2)); //0xff_ff_ff_ff & 0xff_ff_ff_fe = 0xff_ff_ff_fe

        // slli
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("slli    a0,a1,0x2", 0x00259513, ExpectReg(Register.REG_A0, 4)); //1 <<2
        int32 v = int32(uint32(1 << 31));
        checkInstruction("slli    a0,a1,0x1f", 0x01f59513, ExpectReg(Register.REG_A0, v)); //0x10_00_00_00
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("slli    a0,a1,0x2", 0x00259513, ExpectReg(Register.REG_A0, -4)); //0xff_ff_ff_ff <<2 = 0xff_ff_ff_fc
        checkInstruction("slli    a0,a1,0x1f", 0x01f59513, ExpectReg(Register.REG_A0, v)); //0xff_ff_ff_ff <<31=0x10_00_00_00

        //srli 位移超过32位非法，在编译时检查.
        checkInstruction("li      a1,8", 0x00800593, ExpectReg(Register.REG_A1, 8));
        checkInstruction("srli    a0,a1,0x1", 0x0015d513, ExpectReg(Register.REG_A0, 4)); //8>>1 =4
        checkInstruction("srli    a0,a1,0x4", 0x0045d513, ExpectReg(Register.REG_A0, 0)); //8>>4=0
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1));
        checkInstruction("srli    a0,a1,0x1f", 0x01f5d513, ExpectReg(Register.REG_A0, 1)); //0xff_ff_ff_ff>>31=0x01

        //srai 位移超过32位非法，在编译时检查.
        checkInstruction("li      a1,8", 0x00800593, ExpectReg(Register.REG_A1, 8));
        checkInstruction("srai    a0,a1,0x1", 0x4015d513, ExpectReg(Register.REG_A0, 4)); //8>>1=4
        checkInstruction("srai    a0,a1,0x4", 0x4045d513, ExpectReg(Register.REG_A0, 0)); //8>>4=0
        checkInstruction("li      a1,-8", 0xff800593, ExpectReg(Register.REG_A1, -8));
        checkInstruction("srai    a0,a1,0x1", 0x4015d513, ExpectReg(Register.REG_A0, -4)); //正常除法0xff_ff_ff_f8 SE>>1=0xff_ff_ff_fc
        checkInstruction("srai    a0,a1,0x1f", 0x41f5d513, ExpectReg(Register.REG_A0, -1)); //0xff_ff_ff_f8 SE>>31 = 0x_ff_ff_ff_ff

        //beq
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("li      a2,1", 0x00100613, ExpectReg(Register.REG_A2, 1));
        resetPC();
        checkInstruction("beq     a1,a2,12", 0x00c58663, ExpectReg(Register.REG_PC, 12)); //a1 == a2 ? pc+12:pc +4

        //bne
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1));
        checkInstruction("li      a2,2", 0x00200613, ExpectReg(Register.REG_A2, 2));
        resetPC();
        checkInstruction("bne     a1,a2,8", 0x00c59463, ExpectReg(Register.REG_PC, 8)); // a1 !=a2 ? pc+8:pc +4

        //blt
        checkInstruction("li      a1,1", 0x00100593, ExpectReg(Register.REG_A1, 1)); //a1=1
        checkInstruction("li      a2,-2", 0xffe00613, ExpectReg(Register.REG_A2, -2)); //a2=((1<<32)-2)
        resetPC();
        checkInstruction("blt     a1,a2,8", 0x00c5c463, ExpectReg(Register.REG_PC, 4)); // 1 > -2
        //bltu
        resetPC();
        checkInstruction("bltu    a1,a2,8", 0x00c5e463, ExpectReg(Register.REG_PC, 8)); //1 < ((1<<32)-2)

        //bge
        checkInstruction("li      a1,-1", 0xfff00593, ExpectReg(Register.REG_A1, -1)); //a1=0xffff_ffff
        checkInstruction("li      a2,1", 0x00100613, ExpectReg(Register.REG_A2, 1));
        resetPC();
        checkInstruction("bge     a1,a2,8", 0x00c5d863, ExpectReg(Register.REG_PC, 4)); //-1 < 1
        //bgeu
        resetPC();
        checkInstruction("bgeu    a1,a2,8", 0x00c5f463, ExpectReg(Register.REG_PC, 8)); //0xffff_ffff > 1

        //lui
        checkInstruction("lui     a0,0x1", 0x00001537, ExpectReg(Register.REG_A0, 0x1000));
        checkInstruction("lui     a1,0x0", 0x000005b7, ExpectReg(Register.REG_A1, 0));
        checkInstruction("lui     a3,0x80000", 0x800006b7, ExpectReg(Register.REG_A3, int32(uint32(0x80000000)))); //0x8000_0000

        //auipc
        resetPC();
        checkInstruction("auipc   a0,0x1", 0x00001517, ExpectReg(Register.REG_A0, 0x1000)); //0 + 1<<12

        //jal
        resetPC();
        checkInstruction("jal     a1,12", 0x00c005ef, ExpectReg(Register.REG_A1, 4), ExpectReg(Register.REG_PC, 12)); //a1=0 + 4, pc=0 +12
        checkInstruction("jal     a1,12", 0x00c005ef, ExpectReg(Register.REG_A1, 16), ExpectReg(Register.REG_PC, 24)); //a1=12+4, pc=12+12

        //mul
        initRegister();
        checkInstruction("lui     ra,0x8", 0x000080b7, ExpectReg(Register.REG_RA, 0x8_000));
        checkInstruction("addi    ra,ra,-512", 0xe0008093, ExpectReg(Register.REG_RA, 0x7_e00)); //ra=0x7e00
        checkInstruction("lui     sp,0xb6db7", 0xb6db7137, ExpectReg(Register.REG_SP, int32(uint32(0xb6_db7_000))));
        checkInstruction("addi    sp,sp,-585", 0xdb710113, ExpectReg(Register.REG_SP, int32(uint32(0xb6_db6_db7)))); //sp=0xb6_db6_db7
        checkInstruction("mul     a4,ra,sp", 0x02208733, ExpectReg(Register.REG_A4, 0x1200)); //a4=0x1_200 overflow

        //mulh
        initRegister();
        checkInstruction("li      ra,1", 0x00100093, ExpectReg(Register.REG_RA, 1)); //ra=1
        checkInstruction("li      sp,1", 0x00100113, ExpectReg(Register.REG_SP, 1)); //sp=1
        checkInstruction("mulh    a4,ra,sp", 0x02209733, ExpectReg(Register.REG_A4, 0)); //a4=0 high 32 bit is zero
        checkInstruction("li      sp,-1", 0xfff00113, ExpectReg(Register.REG_SP, -1)); //sp=-1
        checkInstruction("mulh    a4,ra,sp", 0x02209733, ExpectReg(Register.REG_A4, -1)); //a4=-1 high 32 bit is all 0xff_ff_ff_ff
        checkInstruction("li      ra,-1", 0xfff00093, ExpectReg(Register.REG_RA, -1)); //sp=-1
        checkInstruction("li      sp,-1", 0xfff00113, ExpectReg(Register.REG_SP, -1)); //sp=-1
        checkInstruction("mulh    a4,ra,sp", 0x02209733, ExpectReg(Register.REG_A4, 0)); //a4=0
        checkInstruction("lui     ra,0xff000", 0xff0000b7, ExpectReg(Register.REG_RA, int32(uint32(0xff_000_000)))); //ra=0xff_000_000
        checkInstruction("lui     sp,0xff000", 0xff000137, ExpectReg(Register.REG_SP, int32(uint32(0xff_000_000)))); //sp=0xff_000_000
        checkInstruction("mulh    a4,ra,sp", 0x02209733, ExpectReg(Register.REG_A4, 0x10_000)); //a4=0x10_000 high 32 bit

        //mlhsu
        initRegister();
        checkInstruction("li      ra,1", 0x00100093, ExpectReg(Register.REG_RA, 1)); //ra=1
        checkInstruction("li      sp,1", 0x00100113, ExpectReg(Register.REG_SP, 1)); //sp=1
        checkInstruction("mulhsu  a4,ra,sp", 0x0220a733, ExpectReg(Register.REG_A4, 0)); //a4=0 high 32 bit is zero
        checkInstruction("li      ra,-1", 0xfff00093, ExpectReg(Register.REG_RA, -1)); //sp=-1
        checkInstruction("li      sp,-1", 0xfff00113, ExpectReg(Register.REG_SP, -1)); //sp=-1
        checkInstruction("mulhsu  a4,ra,sp", 0x0220a733, ExpectReg(Register.REG_A4, -1)); //a4=0xff_ff_ff_ff
        checkInstruction("lui     ra,0xff000", 0xff0000b7, ExpectReg(Register.REG_RA, int32(uint32(0xff_000_000)))); //ra=0xff_000_000
        checkInstruction("lui     sp,0xff000", 0xff000137, ExpectReg(Register.REG_SP, int32(uint32(0xff_000_000)))); //sp=0xff_000_000
        checkInstruction("mulhsu  a4,ra,sp", 0x0220a733, ExpectReg(Register.REG_A4, int32(uint32(0xff_010_000)))); //a4=0xff_010_000

        //mulhu
        initRegister();
        checkInstruction("li      ra,1", 0x00100093, ExpectReg(Register.REG_RA, 1)); //ra=1
        checkInstruction("li      sp,1", 0x00100113, ExpectReg(Register.REG_SP, 1)); //sp=1
        checkInstruction("mulhu   a4,ra,sp", 0x0220b733, ExpectReg(Register.REG_A4, 0)); //a4=0 high 32 bit is zero
        checkInstruction("lui     ra,0xff000", 0xff0000b7, ExpectReg(Register.REG_RA, int32(uint32(0xff_000_000)))); //ra=0xff_000_000
        checkInstruction("lui     sp,0xff000", 0xff000137, ExpectReg(Register.REG_SP, int32(uint32(0xff_000_000)))); //sp=0xff_000_000
        checkInstruction("mulhu   a4,ra,sp", 0x0220b733, ExpectReg(Register.REG_A4, int32(uint32(0xfe_010_000)))); //a4=0xfe_010_000
        checkInstruction("li      ra,-1", 0xfff00093, ExpectReg(Register.REG_RA, -1)); //sp=-1
        checkInstruction("li      sp,-1", 0xfff00113, ExpectReg(Register.REG_SP, -1)); //sp=-1
        checkInstruction("mulhu   a4,ra,sp", 0x0220b733, ExpectReg(Register.REG_A4, -2)); //a4=-2, (2^32-1)*(2^32-1)=(2^32)*(2^32-2)+1

        //div
        initRegister();
        checkInstruction("li      ra,-20", 0xfec00093, ExpectReg(Register.REG_RA, -20)); //ra=-20
        checkInstruction("li      sp,6", 0x00600113, ExpectReg(Register.REG_SP, 6)); //sp=6
        checkInstruction("div     a4,ra,sp", 0x0220c733, ExpectReg(Register.REG_A4, -3)); //a4=-3
        checkInstruction("li      ra,0", 0x00000093, ExpectReg(Register.REG_RA, 0)); //ra=0
        checkInstruction("li      sp,0", 0x00000113, ExpectReg(Register.REG_SP, 0)); //sp=0
        checkInstruction("div     a4,ra,sp", 0x0220c733, ExpectReg(Register.REG_A4, -1)); //a4=-1 all num/0=-1
        checkInstruction("lui     ra,0x80000", 0x800000b7, ExpectReg(Register.REG_RA, minInt)); //ra=0x80_000_000 2^31 最小32位负数，只有符号位为1
        checkInstruction("li      sp,-1", 0xfff00113, ExpectReg(Register.REG_SP, -1)); //sp=-1
        checkInstruction("div     a4,ra,sp", 0x0220c733, ExpectReg(Register.REG_A4, int32(uint32(0x80_000_000)))); //a4=0x80_000_000 最小负数/-1 等于最小负数

        //divu
        initRegister();
        checkInstruction("li      ra,-20", 0xfec00093, ExpectReg(Register.REG_RA, -20)); //ra=(2^32)-20
        checkInstruction("li      sp,6", 0x00600113, ExpectReg(Register.REG_SP, 6)); //sp=6
        checkInstruction("divu    a4,ra,sp", 0x0220d733, ExpectReg(Register.REG_A4, int32(uint32(0x2a_aaa_aa7)))); //a4=0x2a_aaa_aa7
        checkInstruction("li      ra,0", 0x00000093, ExpectReg(Register.REG_RA, 0)); //ra=0
        checkInstruction("li      sp,0", 0x00000113, ExpectReg(Register.REG_SP, 0)); //sp=0
        checkInstruction("divu    a4,ra,sp", 0x0220d733, ExpectReg(Register.REG_A4, -1)); //a4=-1 all num/0=-1

        //rem
        initRegister();
        checkInstruction("li      ra,-20", 0xfec00093, ExpectReg(Register.REG_RA, -20)); //ra=-20
        checkInstruction("li      sp,6", 0x00600113, ExpectReg(Register.REG_SP, 6)); //sp=6
        checkInstruction("rem     a4,ra,sp", 0x0220e733, ExpectReg(Register.REG_A4, -2)); //a4=-2
        checkInstruction("li      ra,20", 0x01400093, ExpectReg(Register.REG_RA, 20)); //ra=20
        checkInstruction("li      sp,-6", 0xffa00113, ExpectReg(Register.REG_SP, -6)); //sp=-6
        checkInstruction("rem     a4,ra,sp", 0x0220e733, ExpectReg(Register.REG_A4, 2)); //a4=2
        checkInstruction("lui     ra,0x80000", 0x800000b7, ExpectReg(Register.REG_RA, minInt)); //ra=0x80_000_000 2^31 最小32位负数，只有符号位为1
        checkInstruction("li      sp,-1", 0xfff00113, ExpectReg(Register.REG_SP, -1)); //sp=-1
        checkInstruction("rem     a4,ra,sp", 0x0220e733, ExpectReg(Register.REG_A4, 0)); //a4=0
        checkInstruction("li      ra,0", 0x00000093, ExpectReg(Register.REG_RA, 0)); //ra=0
        checkInstruction("li      sp,0", 0x00000113, ExpectReg(Register.REG_SP, 0)); //sp=0
        checkInstruction("rem     a4,ra,sp", 0x0220e733, ExpectReg(Register.REG_A4, 0)); //a4=0 all num%0=0

        //remu
        initRegister();
        checkInstruction("li      ra,-20", 0xfec00093, ExpectReg(Register.REG_RA, -20)); //ra=(2^32)-20
        checkInstruction("li      sp,6", 0x00600113, ExpectReg(Register.REG_SP, 6)); //sp=6
        checkInstruction("remu    a4,ra,sp", 0x0220f733, ExpectReg(Register.REG_A4, 2)); //a4=2
        checkInstruction("lui     ra,0x80000", 0x800000b7, ExpectReg(Register.REG_RA, minInt)); //ra=0x80_000_000 2^31 最小32位负数，只有符号位为1
        checkInstruction("li      sp,-1", 0xfff00113, ExpectReg(Register.REG_SP, -1)); //sp=(2^32)-1
        checkInstruction("remu    a4,ra,sp", 0x0220f733, ExpectReg(Register.REG_A4, minInt)); //
        checkInstruction("li      ra,0", 0x00000093, ExpectReg(Register.REG_RA, 0)); //ra=0
        checkInstruction("li      sp,0", 0x00000113, ExpectReg(Register.REG_SP, 0)); //sp=0
        checkInstruction("remu    a4,ra,sp", 0x0220f733, ExpectReg(Register.REG_A4, 0)); // all num%0=0
    }
}
