set -ex

files="
AddressManager.sol/*
ChainStorageContainer.sol/*
Challenge.sol/*
ERC20.sol/*
L1CrossLayerWitness.sol/*
L2CrossLayerWitness.sol/*
RollupInputChain.sol/*
RollupStateChain.sol/*
StakingManager.sol/*
TestERC20.sol/*
ChallengeFactory.sol/*
UpgradeableBeacon.sol/*
ProxyAdmin.sol/*
TransparentUpgradeableProxy.sol/*
DAO.sol/*
L2FeeCollector.sol/*
"
prefix="out/"
for file in $files ; do
  go-web3 --output binding/ --package binding --source $prefix$file
done