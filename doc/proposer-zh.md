# proposer
`proposer`为可以向`RollupStateChain`进行`appendStateBatch`操作的以太坊账户
其需要两个前提:
- 需要在`StakingManager`中进行质押，保证对作恶的proposer进行经济惩罚
- 需要在`DAO`中的proposerWhiteList内，保证现阶段的安全性

进行操作:
- 确保账户在`FeeToken`的erc20代币足够并对`StakingManager`合约进行`deposite()`质押。
- 联系`DAO`管理员将对应账户加入白名单。