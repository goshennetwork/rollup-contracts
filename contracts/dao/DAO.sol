// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../interfaces/IDAO.sol";

contract DAO is IDAO, OwnableUpgradeable {
    mapping(address => bool) public override sequencerWhitelist;
    mapping(address => bool) public override proposerWhitelist;
    mapping(address => bool) public override challengerWhitelist;

    event SequencerWhitelistUpdated(address submitter, bool enabled);
    event ProposerWhitelistUpdated(address proposer, bool enabled);
    event ChallengerWhitelistUpdated(address challenger, bool enabled);

    function initialize() public initializer {
        __Ownable_init();
    }

    function setSequencerWhitelist(address sequencer, bool enabled) public onlyOwner {
        sequencerWhitelist[sequencer] = enabled;
        emit SequencerWhitelistUpdated(sequencer, enabled);
    }

    function setProposerWhitelist(address proposer, bool enabled) public onlyOwner {
        proposerWhitelist[proposer] = enabled;
        emit ProposerWhitelistUpdated(proposer, enabled);
    }

    function setChallengerWhitelist(address challenger, bool enabled) public onlyOwner {
        challengerWhitelist[challenger] = enabled;
        emit ChallengerWhitelistUpdated(challenger, enabled);
    }

    // used to transfer fee token
    function transferERC20(IERC20 token, address to, uint256 amount) public onlyOwner {
        token.transfer(to, amount);
    }
}
