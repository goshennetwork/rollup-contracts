set -ex

files="
AddressManager.sol/*
ChainStorageContainer.sol/*
Challenge.sol/*
ERC20.sol/*
L1CrossLayerWitness.sol/*
L2CrossLayerWitness.sol/*
L1StandardBridge.sol/*
L2StandardBridge.sol/*
RollupInputChain.sol/*
RollupStateChain.sol/*
StakingManager.sol/*
TestERC20.sol/*
ChallengeFactory.sol/*
UpgradeableBeacon.sol/*
ProxyAdmin.sol/*
TransparentUpgradeableProxy.sol/*
Whitelist.sol/*
L2FeeCollector.sol/*
TestL2ERC20.sol/*
L2StandardERC20.sol/*
"
prefix="out/"
for file in $files ; do
  go-web3 --output binding/ --package binding --source $prefix$file
done