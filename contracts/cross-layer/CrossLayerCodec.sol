// SPDX-License-Identifier: GPL v3
pragma solidity ^0.8.0;

import "../libraries/BytesSlice.sol";

library CrossLayerCodec {
    using BytesSlice for Slice;

    function crossLayerMessageHash(
        address _target,
        address _sender,
        uint64 _messageIndex,
        bytes memory _message
    ) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(_target, _sender, _messageIndex, _message));
    }

    function encodeL1ToL2CallData(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex,
        bytes32 _mmrRoot,
        uint64 _mmrSize
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "relayMessage(address,address,bytes,uint64,bytes32,uint64)",
                _target,
                _sender,
                _message,
                _messageIndex,
                _mmrRoot,
                _mmrSize
            );
    }

    function encodeL2ToL1CallData(
        address _target,
        address _sender,
        bytes memory _message,
        uint64 _messageIndex
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "relayMessage(address,address,bytes,uint64)",
                _target,
                _sender,
                _message,
                _messageIndex
            );
    }
}
