// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "./UnsafeSign.sol";
import "./RLPWriter.sol";
import "./console.sol";

contract TestUnsafeSign {
    function testSign() public {
        for (uint256 i = 0; i < 100000; i++) {
            bytes32 signedHash = keccak256(abi.encode(i, "test"));
            (uint256 r, uint256 s, uint64 v) = UnsafeSign.Sign(signedHash, 1);
            uint64 _pureV = v - 2 * 1 - 8;
            require(_pureV <= 28, "invalid v");
            address sender = ecrecover(signedHash, uint8(_pureV), bytes32(r), bytes32(s));
            require(sender == UnsafeSign.SENDER, "wrong sender");
        }
    }

    function testTx() public {
        bytes memory data = new bytes(0);
        bytes[] memory _rlpList = getRlpList(0, 0, address(0), data);
        bytes32 _signTxHash = keccak256(RLPWriter.writeList(_rlpList));
        require(_signTxHash == keccak256(RLPWriter.writeList(_rlpList)));
        require(
            _signTxHash == bytes32(uint256(0x8c6115c6530a74eb5904bc51bcc0c8777c2e6144f20c04821ad703e301eef28c)),
            "wrong signed hash"
        );
        (uint256 r, uint256 s, uint64 v) = UnsafeSign.Sign(_signTxHash, 1337);
        //now change rsv value in tx to calc tx's hash
        _rlpList[6] = RLPWriter.writeUint(v);
        _rlpList[7] = RLPWriter.writeUint(r);
        _rlpList[8] = RLPWriter.writeUint(s);
        bytes32 _txHash = keccak256(RLPWriter.writeList(_rlpList));
        require(
            _txHash == bytes32(uint256(0xbabdfa0c36828141cf854c840a1b55103b949bb06407396a28bcdaabf5d6aad9)),
            "wrong tx hash"
        );
    }

    //encode tx params: sender, to, gasLimit, data, nonce, r,s,v and gasPrice(1 GWEI), value(0), chainId
    //sender used to recognize tx from L1CrossLayerWitness
    function getRlpList(
        uint64 _nonce,
        uint64 _gasLimit,
        address _target,
        bytes memory _data
    ) internal pure returns (bytes[] memory) {
        bytes[] memory list = new bytes[](9);
        list[0] = RLPWriter.writeUint(uint256(_nonce));
        list[1] = RLPWriter.writeUint(1_000_000_000);
        list[2] = RLPWriter.writeUint(uint256(_gasLimit));
        list[3] = RLPWriter.writeAddress(_target);
        list[4] = RLPWriter.writeUint(0);
        list[5] = RLPWriter.writeBytes(_data);
        list[6] = RLPWriter.writeUint(1337);
        list[7] = abi.encodePacked(bytes1(0x80));
        list[8] = abi.encodePacked(bytes1(0x80));
        return list;
    }
}
