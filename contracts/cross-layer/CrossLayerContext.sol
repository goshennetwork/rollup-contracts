// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../interfaces/ICrossLayerWitness.sol";

contract CrossLayerContext is Initializable {
    ICrossLayerWitness public crossLayerWitness;

    function initialize(address _witness) public initializer {
        crossLayerWitness = ICrossLayerWitness(_witness);
    }

    modifier ensureCrossLayerSender(address _sourceLayerSender) {
        require(msg.sender == address(crossLayerWitness), "no permission");

        address crossLayersender = crossLayerWitness.crossLayerSender();
        require(crossLayersender == _sourceLayerSender, "wrong cross layer sender");

        _;
    }

    function sendCrossLayerMessage(address _target, bytes memory _message) internal {
        crossLayerWitness.sendMessage(_target, _message);
    }
}
