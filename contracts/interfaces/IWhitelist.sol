// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../interfaces/IAddressResolver.sol";

interface IWhitelist {
    /// EVENT
    event SequencerUpdated(address submitter, bool enabled);
    event ProposerUpdated(address proposer, bool enabled);
    event ChallengerUpdated(address challenger, bool enabled);

    /// FUNCTION
    function canSequence(address addr) external view returns (bool);

    function canPropose(address addr) external view returns (bool);

    function canChallenge(address addr) external view returns (bool);

    function setSequencer(address sequencer, bool enabled) external;

    function setProposer(address proposer, bool enabled) external;

    function setChallenger(address challenger, bool enabled) external;
}
