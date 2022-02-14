pragma solidity >0.5.0 <0.8.0;

import "./MerkleTrie.sol";

interface MTrie {
    function update(
        bytes memory _key,
        bytes memory _value,
        bytes32 _root,
        bytes32 _expectRoot
    ) external;

    function get(bytes memory _key, bytes32 _root) external returns (bytes memory);

    function insertTrieNode(bytes calldata anything) external;
}

contract MerkleTemp is MTrie {
    function update(
        bytes memory _key,
        bytes memory _value,
        bytes32 _root,
        bytes32 _expectRoot
    ) external override {
        bytes32 _getRoot = Lib_MerkleTrie.update(_key, _value, _root);
        require(_getRoot == _expectRoot, "not equal");
    }

    function get(bytes memory _key, bytes32 _root) external override returns (bytes memory) {
        (bool exist, bytes memory _res) = Lib_MerkleTrie.get(_key, _root);
        require(exist, "not exist");
        return _res;
    }

    function insertTrieNode(bytes calldata anything) external override {
        Lib_MerkleTrie.GetTrie()[keccak256(anything)] = anything;
    }
}
