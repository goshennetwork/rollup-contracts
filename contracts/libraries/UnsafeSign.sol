// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

library UnsafeSign {
    // order of secp256k1 curve
    uint256 internal constant ORDER = 115792089237316195423570985008687907852837564279074904382605163141518161494337;
    // half order of secp256k1 curve
    uint256 internal constant HALF_ORDER = ORDER >> 1;

    // primitive element's( known as 'g') x in secp256k1 curve point
    uint256 internal constant GX = 55066263022277343669578718895168534326250603453777594175500187360389116729240;

    // the address of G
    address internal constant GADDR = address(0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf);
    // the address of G*2
    address internal constant G2ADDR = address(0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF);

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
        uint256 order = ORDER; // cache here to reduce code size
        uint256 e = uint256(signedHash) % order;
        // make sure not overflow
        uint256 s = GX;
        unchecked {
            if (s + e < s) {
                // overflow use inverse num
                s = order - ((order - s) + (order - e));
            } else {
                //not overflow just calc
                s = (s + e) % order;
            }
        }
        uint64 v = 27;
        if (s > HALF_ORDER) {
            // s is the scale of the curve(just like private key), so when s beyond half order, just use inverse element
            v = 28;
            s = order - s;
        }
        // only happen when pv*r+e=order*integer
        require(s != 0, "zero s");
        // now add chainId
        v += 8 + chainId * 2;
        return (GX, s, v);
    }

    /// sign specific hash and chainId with privkey == 2, return r,s,v
    function Sign2(bytes32 signedHash, uint64 chainId)
        internal
        pure
        returns (
            uint256,
            uint256,
            uint64
        )
    {
        (uint256 r, uint256 s, uint64 v) = Sign(signedHash, chainId);
        uint256 order = ORDER; // cache here to reduce code size
        // make sure not overflow
        unchecked {
            if (s + GX < s) {
                // overflow use inverse num
                s = order - ((order - s) + (order - GX));
            } else {
                //not overflow just calc
                s = (s + GX) % order;
            }
        }
        require(s != 0, "zero s");

        return (r, s, v);
    }
}
