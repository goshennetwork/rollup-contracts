# Sequenceer

`sequencer` is an ethereum account that have the right to push input of l2 system at `RollupInputChain` contract by `appendBatch`
- is included at `DAO`'s proposerWhiteList, to guarantee the system safety at this period
- deposited first at `StakingManager`, make sure the malicious sequencer will get punished

the step to be a sequencer:
1: the account should have the enough balance at `FeeToken`(the recommended amount can be queryed by `price()` method at `StakingManager` contract)
2: the account should invoke the `deposit()` method at `StakingManager` contract
3: the account should contract with the manager of `DAO` to make it involved in sequencer white list