// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

library UnsafeSign {
    //order of secp256k1 curve
    uint256 internal constant ORDER = 115792089237316195423570985008687907852837564279074904382605163141518161494337;
    //half order of secp256k1 curve
    uint256 internal constant HALF_ORDER =
        115792089237316195423570985008687907852837564279074904382605163141518161494337 >> 1;

    //primitive element's( known as 'g') x in secp256k1 curve point
    uint256 internal constant Primitive_ELEMENT_X =
        55066263022277343669578718895168534326250603453777594175500187360389116729240;

    //the scale of r_point in secp256k1 curve (r_point=g_point^k, so when k=1, r_point=g_point,r_x(known as 'r')=g_x)
    uint256 internal constant PRIVE_K = 1;

    //inverse in multy ring, this is used for verify signature, g^(inv*k)=g,so inv * k =1 +N*integer, when k =1, inv simply calc to 1
    uint256 internal constant MOD_INVERSE = 1;

    address internal constant SENDER = address(0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf);
    uint256 internal constant PRIVATE_KEY = 1;
    //just r_x
    uint256 internal constant MUL_R_WITH_PRIVATE =
        55066263022277343669578718895168534326250603453777594175500187360389116729240;

    ///@dev sign specific hash and chainId, return r,s,v
    function Sign(bytes32 signedHash, uint64 chainId)
        internal
        pure
        returns (
            uint256,
            uint256,
            uint64
        )
    {
        uint256 e = uint256(signedHash) % ORDER;
        //make sure not overflow
        uint256 s = MUL_R_WITH_PRIVATE;
        unchecked {
            if (s + e < s) {
                //overflow use inverse num
                s = ORDER - ((ORDER - s) + (ORDER - e));
            } else {
                //not overflow just calc
                s = s + e;
            }
        }
        //mod inverse is 1, just ignore, now s=k_inverse*(dr+e)modN
        s = s % ORDER;
        bool vice = false;
        if (s > HALF_ORDER) {
            //s is the scale of the curv(just like private key), so when s beyond half order, just use inverse element
            vice = true;
            s = ORDER - s;
        }
        //only happen when pv*r+e=order*integer
        //todo: maybe use more than just one k to avoid special signed hahs attack?
        require(s != 0, "zero s");
        uint64 v = 27;
        if (vice) {
            v = 28;
        }
        //now add chainId
        v += 8 + chainId * 2;
        return (Primitive_ELEMENT_X, s, v);
    }
}
