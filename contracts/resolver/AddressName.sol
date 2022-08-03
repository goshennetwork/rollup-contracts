// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

library AddressName {
    ///DAO
    string constant DAO = "DAO";
    bytes32 constant DAO_HASH = keccak256("DAO");
    ///RollupInputChain
    string constant ROLLUP_INPUT_CHAIN = "RollupInputChain";
    bytes32 constant ROLLUP_INPUT_CHAIN_HASH = keccak256("RollupInputChain");
    ///ChainStorageContainer of RollupInputChain
    string constant ROLLUP_INPUT_CHAIN_CONTAINER = "RollupInputChainContainer";
    bytes32 constant ROLLUP_INPUT_CHAIN_CONTAINER_HASH = keccak256("RollupInputChainContainer");
    ///RollupStateChain
    string constant ROLLUP_STATE_CHAIN = "RollupStateChain";
    bytes32 constant ROLLUP_STATE_CHAIN_HASH = keccak256("RollupStateChain");
    ///ChainStorageContainer of RollupStateChain
    string constant ROLLUP_STATE_CHAIN_CONTAINER = "RollupStateChainContainer";
    bytes32 constant ROLLUP_STATE_CHAIN_CONTAINER_HASH = keccak256("RollupStateChainContainer");
    ///StakingManager
    string constant STAKING_MANAGER = "StakingManager";
    bytes32 constant STAKING_MANAGER_HASH = keccak256("StakingManager");
    ///ChallengeFactory
    string constant CHALLENGE_FACTORY = "ChallengeFactory";
    bytes32 constant CHALLENGE_FACTORY_HASH = keccak256("ChallengeFactory");
    ///L1CrossLayerWitness
    string constant L1_CROSS_LAYER_WITNESS = "L1CrossLayerWitness";
    bytes32 constant L1_CROSS_LAYER_WITNESS_HASH = keccak256("L1CrossLayerWitness");
    ///L2CrossLayerWitness
    string constant L2_CROSS_LAYER_WITNESS = "L2CrossLayerWitness";
    bytes32 constant L2_CROSS_LAYER_WITNESS_HASH = keccak256("L2CrossLayerWitness");
    ///StateTransition
    string constant STATE_TRANSITION = "StateTransition";
    bytes32 constant STATE_TRANSITION_HASH = keccak256("StateTransition");
    ///L1StandardBridge
    string constant L1_STANDARD_BRIDGE = "L1StandardBridge";
    bytes32 constant L1_STANDARD_BRIDGE_HASH = keccak256("L1StandardBridge");
    ///ChallengeBeacon
    string constant CHALLENGE_BEACON = "ChallengeBeacon";
    bytes32 constant CHALLENGE_BEACON_HASH = keccak256("ChallengeBeacon");
    ///FeeToken
    string constant FEE_TOKEN = "FeeToken";
    bytes32 constant FEE_TOKEN_HASH = keccak256("FeeToken");
    ///MachineState
    string constant MACHINE_STATE = "MachineState";
    bytes32 constant MACHINE_STATE_HASH = keccak256("MachineState");
    ///Whitelist
    string constant WHITELIST = "Whitelist";
    bytes32 constant WHITELIST_HASH = keccak256("Whitelist");
}
