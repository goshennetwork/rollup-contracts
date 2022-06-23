// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./console.sol";
import {CompactMerkleTree, MerkleMountainRange} from "./MerkleMountainRange.sol";

contract MMRTest {
    using MerkleMountainRange for CompactMerkleTree;
    CompactMerkleTree _trees;

    function getTreeSize() public view returns (uint64) {
        return _trees.treeSize;
    }

    function getRootHash() public view returns (bytes32) {
        return _trees.rootHash;
    }

    function append(bytes32 _leafHash) public {
        _trees.appendLeafHash(_leafHash);
    }

    function verifyProof(
        bytes32 _leafHash,
        uint64 _leafIndex,
        bytes32[] memory _proof,
        bytes32 _rootHash,
        uint64 _treeSize
    ) public pure {
        MerkleMountainRange.verifyLeafHashInclusion(_leafHash, _leafIndex, _proof, _rootHash, _treeSize);
    }

    function testAppend() public {
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 1, "0");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 1, "1");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 2, "2");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 1, "3");

        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 2, "4");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 2, "5");
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0));
        require(_trees.hashes.length == 3, "6");
    }

    function appendLeaf(uint64 _size) internal {
        for (uint64 i = 0; i < _size; i++) {
            MerkleMountainRange.appendLeafHash(_trees, bytes32(uint256(i)));
        }
    }

    function testZeroCopyAndAbiEncodePacked() public {
        address target = 0xEC9C107cf2D52B4E771301c3d702196D2e163bDC;
        address msgSender = 0x9A2900E4b204E31dD58eCc8F276808169D8E4A1b;
        uint64 msgIndex = 777777777;
        bytes memory msg = 'asdfafdfasfasdfaddfadjfatjydfagjfgajkdakljfakdlgajkhgasjhgajg';
        bytes memory data = abi.encodePacked(target, msgSender, msgIndex, msg);
        console.logBytes(data);
        console.logBytes32(keccak256(data));
    }

    function testVerify() public {
        // m=5, n=10
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0x656c98d56eadba8c4938fd4153bb51fd2c32f068c78594342e39fd8c1b632332));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0xc5078ae0bc75a0052209ebf1e0638ff2b824e3892e12f7d2863e7c62a3fe502e));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0xc2d9dcb829a4a878e5a18c6f3a4f25926dd1f1e51c3ed08d4b15e6474f179955));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0x99e57d9f68afe3e6fabf0f2b37b930a33d8631c23e653f03b62cd4745194eed4));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0xbfd88be2f23b6aa4d412e75ff774853b90ad8b4267ca99d2714dde4a706ecefa));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0xa91fdaa6209a0ab99d30f19f1327c55e12a5ac41f559fe9a6220c7abc00584a2));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0x9de0720cb4d747cad3702f50a6cdb35cf2f2738ab0843eacd6b4d158c0390bef));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0x867d11d93c3e54a3af819243a8813354286aeeb155835d7dda1754c95334a244));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0x4942f139a43e6502fbe3d6c72b1cd07c1c4daba4e0a77cd6cdc88ec1777045af));
        console.logBytes32(_trees.rootHash);
        MerkleMountainRange.appendLeafHash(_trees, bytes32(0x33004dc58f858443ceacfab70224ac91f2aebca48ab56ff258aee157fc825806));
        console.logBytes32(_trees.rootHash);
        bytes32[] memory _proof = new bytes32[](4);
        _proof[0] = bytes32(0xbfd88be2f23b6aa4d412e75ff774853b90ad8b4267ca99d2714dde4a706ecefa);
        _proof[1] = bytes32(0xecd38a5aa1d25ac31d019fce384d9502ac6abb9b04834998041fc094bd017acb);
        _proof[2] = bytes32(0xc71018f24e83677976bdf7e40941d5d54e713091339cb3f856390756b5af17b6);
        _proof[3] = bytes32(0xc88f10180627ad4cb58aace0f77c8d31e9c4741d830bd138c407d8768eddf04a);
        MerkleMountainRange.verifyLeafHashInclusion(
            bytes32(0xa91fdaa6209a0ab99d30f19f1327c55e12a5ac41f559fe9a6220c7abc00584a2),
            uint64(5), _proof, _trees.rootHash, _trees.treeSize);
    }
}
