# L2 Init State

## Local EVM

the operation at LocalEVM

- deploy the  L2CrossWitness, L2StandardBridge with openzepplin-upgrades but do not use proxyAdmin contract
- initialize L2StandardBridge
- deploy FeeCollector
- 5 contracts got
    - L2CrossWitness
    - L2StandardBridge
    - L2CrossWitnessLogic
    - L2StandardBridgeLogic
    - FeeCollector

## Genesis

- config the address of L2CrossWitness, L2StandardBridge, L2CrossWitnessLogic, L2StandardBridgeLogic,FeeCollector
- read the code and storage from evm and set to genesis block
- set the balance of L2StandardBridge to 1B Ether
