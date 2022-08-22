pragma solidity ^0.8.0;

import "../interfaces/ForgeVM.sol";
import "./PreCompile.sol";

contract TestPreCompile {
    ForgeVM public constant vm = ForgeVM(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

    PreCompile public precompile;

    uint256 private privKey = 0xdf57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e;

    function setUp() public {
        vm.startPrank(address(0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199));
        precompile = new PreCompile();
    }

    function testCheckSig() public {
        bytes32 hash = keccak256(abi.encodePacked("test precomile ecrecover"));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privKey, hash);
        precompile.checkSig(hash, r, s, v);
    }

    function testCheckSha256() public {
        bytes32 hash = keccak256(abi.encodePacked("test precomile sha256"));
        precompile.checkSha256(hash);
    }

    function testCheckRipemd160() public {
        bytes32 hash = keccak256(abi.encodePacked("test precomile ripemd160"));
        precompile.checkRipemd160(hash);
    }

    function testCheckDataCopy() public {
        bytes32 hash = keccak256(abi.encodePacked("test precomile datacopy"));
        precompile.checkDataCopy(hash);
    }

    function testBlakeF2() public {
        uint32 rounds = 12;

        bytes32[2] memory h;
        h[0] = hex"48c9bdf267e6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa5";
        h[1] = hex"d182e6ad7f520e511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b";

        bytes32[4] memory m;
        m[0] = hex"6162630000000000000000000000000000000000000000000000000000000000";
        m[1] = hex"0000000000000000000000000000000000000000000000000000000000000000";
        m[2] = hex"0000000000000000000000000000000000000000000000000000000000000000";
        m[3] = hex"0000000000000000000000000000000000000000000000000000000000000000";

        bytes8[2] memory t;
        t[0] = hex"03000000";
        t[1] = hex"00000000";

        bool f = true;

        bytes32[2] memory result = precompile.blake2F(rounds, h, m, t, f);
        require(result[0] == bytes32(0xba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d1));
        require(result[1] == bytes32(0x7d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923));
    }
}
