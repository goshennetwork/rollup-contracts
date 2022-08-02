// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../interfaces/IAddressResolver.sol";

contract Whitelist is Initializable {
    mapping(address => bool) public override canSequence;
    mapping(address => bool) public override canPropose;
    mapping(address => bool) public override canChallenge;

    IAddressResolver addressResolver;

    event SequencerUpdated(address submitter, bool enabled);
    event ProposerUpdated(address proposer, bool enabled);
    event ChallengerUpdated(address challenger, bool enabled);

    modifier onlyDAO() {
        require(msg.sender == address(addressResolver.dao()), "only dao allowed");
        _;
    }

    function initialize(IAddressResolver _resolver) public initializer {
        addressResolver = _resolver;
    }

    function setSequencer(address sequencer, bool enabled) public onlyDAO {
        canSequence[sequencer] = enabled;
        emit SequencerUpdated(sequencer, enabled);
    }

    function setProposer(address proposer, bool enabled) public onlyDAO {
        canPropose[proposer] = enabled;
        emit ProposerUpdated(proposer, enabled);
    }

    function setChallenger(address challenger, bool enabled) public onlyDAO {
        canChallenge[challenger] = enabled;
        emit ChallengerUpdated(challenger, enabled);
    }
}
