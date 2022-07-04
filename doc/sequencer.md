# Sequenceer

一个可以对RollupInputChain合约进行appendBatch操作的账户即为`sequencer`,其需要两个权限:
- 需要在`DAO`的白名单内，以便在现阶段保证足够的安全性
- 需要在`StakingManager`中进行抵押，以便对欺诈的sequencer进行经济惩罚

设置权限的步骤:
1: 账户A需要保证在`FeeToken`ERC20合约中有足够的可用于后续操作的token（数额可查询`StakingManager`的`price()`方法）。
2: 账户A需要调用`StakingManager`的`deposit()`进行质押。
3: 账户A需要联系控制`DAO`合约的管理员调用其`setSequencerWhitelist`方法将账户A加入白名单。