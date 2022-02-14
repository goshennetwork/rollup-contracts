pragma solidity >0.5.0 <0.8.0;
import "./MerkleTrie.sol";

contract MerkleCodecTest {
    /*because this merkle is used to state machine, so it's only need 2 method: get specific data of path
    and update the new trie root when leaf value is changed
    */
    bytes32 private constant emptyRoot = bytes32(hex"56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421");

    function testEmptyTrie() public {
        require(Lib_MerkleTrie.KECCAK256_RLP_NULL_BYTES == emptyRoot);
    }

    function checkBytes(bytes memory src, bytes memory dest) private {
        require(keccak256(src) == keccak256(dest), "bytes not equal");
    }

    function testNull() public {
        bytes memory _key = new bytes(32);
        bytes memory _value = "test";
        bytes32 _updatedRoot = Lib_MerkleTrie.update(_key, _value, emptyRoot);
        (bool _exist, bytes memory _result) = Lib_MerkleTrie.get(_key, _updatedRoot);
        require(_exist, "not exist");
        checkBytes(_value, _result);
    }
}
