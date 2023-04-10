# proposer
`proposer` is the ethereum account that can push state to `RollupStateChain` by `appendStateBatch`
required:
- deposited first at `StakingManager`, make sure the malicious proposer will get punished
- is included at `DAO`'s proposerWhiteList, to guarantee the system safety at this period

the step to be a proposer:
- make sure the account have enough token at `FeeToken`, and deposit at `StakingManager`
- contact with the manager of `DAO` to make the previus account to be involved in proposer white list.