# L2 Init State

## Local EVM

在LocalEVM执行下列操作

- 用go实现openzepplin-upgrades插件部署L2CrossWitness、L2StandardBridge的过程；
- 初始化L2StandardBridge
- 部署FeeCollector

## Genesis

- 配置L2CrossWitness、L2StandardBridge和FeeCollector的地址；
- 将三本合约的code、storage从Local EVM里面读取出来，在创世块中保存到对应的地址下；
- 设置L2StandardBridge的余额为10亿 ether
