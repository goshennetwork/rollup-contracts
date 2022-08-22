pragma solidity ^0.8.0;

//import 'hardhat/console.sol';

contract PreCompile {
    address public signer;

    bytes32 public data;

    constructor(){
        setSigner(msg.sender);
    }

    function setSigner(address _signer) public {
        signer = _signer;
    }

    function writeData(bytes32 hash, bytes32 r, bytes32 s, uint8 v) public {
        checkSig(hash, r, s, v);
        data = hash;

        {
            uint32 rounds = 12;

            bytes32[2] memory h;
            h[0] = hex"48c9bdf26aa6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa5";
            h[1] = hex"d182e6ad7f5bbe511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b";

            bytes32[4] memory m;
            m[0] = hex"6162630000000000000000000000000000000000000000000000000000000000";
            m[1] = hex"000000000000cc00000000000000000000000000000000000000000000000000";
            m[2] = hex"0000000000000ddd000000000000000000000000000000000000000000000000";
            m[3] = hex"0000000000000000000000000000ff0000000000000000000000000000000000";

            bytes8[2] memory t;
            t[0] = hex"03000000";
            t[1] = hex"00023000";

            bool f = true;

            blake2F(rounds, h, m, t, f);
        }

        checkSha256(data);
        checkRipemd160(data);
        checkDataCopy(data);

        uint d = bigmodexp(uint(hash), uint(r), uint(s));
        uint256[2] memory addG = bn256Add(uint(hash), uint(r), uint(s), d);
        uint256[2] memory scalarMulG = bn256ScalarMul(uint(hash), uint(r), uint(s));
        bool pairing = bn256Pairing(uint(hash), uint(r), uint(s), d, addG[0], addG[1]);
    }


    function checkSig(bytes32 hash, bytes32 r, bytes32 s, uint8 v) public view {
        address signer1 = ecrecover(hash, v, r, s);
        require(signer1 == signer, "ecrecover failed");
        uint256[1] memory result;
        bool success;
        uint256[4] memory input = [uint(hash), uint(v), uint(r), uint(s)];
        assembly{
            success := staticcall(gas(), 0x01, input, 0x80, result, 0x20)
        }
        require(success, "call failed");
        address signer2 = address(uint160(result[0]));
        require(signer2 == signer, "ecrecover-precompile failed");
    }

    function checkSha256(bytes32 data) public view {
        bytes32 res1 = sha256(abi.encodePacked(data));
        uint256[1] memory input = [uint(data)];
        uint256[1] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x02, input, 0x20, result, 0x20)
        }
        require(success, "call failed");
        require(res1 == bytes32(result[0]), "sha256-precompile failed");
    }

    function checkRipemd160(bytes32 data) public view {
        bytes20 res1 = ripemd160(abi.encodePacked(data));
        uint256[1] memory input = [uint(data)];
        uint256[1] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x03, input, 0x20, result, 0x20)
        }
        require(success, "call failed");
        require(res1 == bytes20(uint160(result[0])), "ripemd160-precompile failed");
    }

    function checkDataCopy(bytes32 data) public view {
        uint256[1] memory input = [uint(data)];
        uint256[1] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x04, input, 0x20, result, 0x20)
        }
        require(success, "call failed");
        require(data == bytes32(result[0]), "ripemd160-precompile failed");
    }

    function bigmodexp(uint a, uint b, uint c) public view returns (uint) {
        uint256[3] memory input = [a, b, c];
        uint256[1] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x05, input, 0x60, result, 0x20)
        }
        require(success, "Bigmodexp failed");
        return result[0];
    }

    function bn256Add(uint a, uint b, uint c, uint d) public view returns (uint256[2] memory) {
        uint256[4] memory input = [a, b, c, d];
        uint256[2] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x06, input, 0x80, result, 0x40)
        }
        require(success, "bn256Add failed");
        return result;
    }

    function bn256ScalarMul(uint a, uint b, uint c) public view returns (uint256[2] memory) {
        uint256[3] memory input = [a, b, c];
        uint256[2] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x06, input, 0x60, result, 0x40)
        }
        require(success, "bn256ScalarMul failed");
        return result;
    }

    function bn256Pairing(uint a, uint b, uint c, uint d, uint e, uint f) public view returns (bool) {
        uint256[6] memory input = [a, b, c, d, e, f];
        uint256[1] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x06, input, 0x12, result, 0x20)
        }
        require(success, "bn256Pairing failed");
        return result[0] == 1;
    }

    function blake2F(uint32 rounds, bytes32[2] memory h, bytes32[4] memory m, bytes8[2] memory t, bool f) public view returns (bytes32[2] memory) {
        bytes memory input = abi.encodePacked(rounds, h[0], h[1], m[0], m[1], m[2], m[3], t[0], t[1], f);
        bytes32[2] memory result;
        bool success;
        assembly{
            success := staticcall(gas(), 0x09, add(input, 32), 0xd5, result, 0x40)
        }
        require(success, "blake2F failed");
        return result;
    }

}
