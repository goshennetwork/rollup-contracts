set -ex

files="
AddressManager.sol/AddressManager.json
ChainStorageContainer.sol/ChainStorageContainer.json
Challenge.sol/Challenge.json
ERC20.sol/ERC20.json
L1CrossLayerWitness.sol/L1CrossLayerWitness.json
L2CrossLayerWitness.sol/L2CrossLayerWitness.json
L1StandardBridge.sol/L1StandardBridge.json
L2StandardBridge.sol/L2StandardBridge.json
RollupInputChain.sol/RollupInputChain.json
RollupStateChain.sol/RollupStateChain.json
StakingManager.sol/StakingManager.json
TestERC20.sol/TestERC20json
ChallengeFactory.sol/ChallengeFactory.json
UpgradeableBeacon.sol/UpgradeableBeacon.json
ProxyAdmin.sol/ProxyAdmin.json
TransparentUpgradeableProxy.sol/TransparentUpgradeableProxy.json
Whitelist.sol/.json
L2FeeCollector.sol/L2FeeCollector.json
TestL2ERC20.sol/TestL2ERC20.json
L2StandardERC20.sol/L2StandardERC20.json
"
prefix="out/"
for file in $files ; do
  go-web3 --output binding/ --package binding --source $prefix$file
done