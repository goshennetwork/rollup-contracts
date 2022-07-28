const config = require("./config/config.json");

async function main() {
    const decimals = 18;
    const AddressManager = await ethers.getContractFactory("AddressManager");
    const addressManager = await upgrades.deployProxy(AddressManager, []);
    console.log("sent AddressManager deploy tx, %s", addressManager.deployTransaction.hash);

    const L1CrossLayerWitness = await ethers.getContractFactory("L1CrossLayerWitness");
    const l1CrossLayerWitness = await upgrades.deployProxy(L1CrossLayerWitness, [addressManager.address]);
    console.log("sent L1CrossLayerWitness deploy tx, %s", l1CrossLayerWitness.deployTransaction.hash);

    const TestERC20 = await ethers.getContractFactory("TestERC20");
    let feeToken;
    if (config.feeToken) {
        feeToken = await TestERC20.attach(config.feeToken);
    } else {
        feeToken = await TestERC20.deploy("Test Fee Token", 'TFT', decimals);
        console.log("sent FeeToken deploy tx, %s", feeToken.deployTransaction.hash);
    }

    const RollupStateChain = await ethers.getContractFactory("RollupStateChain");
    const rollupStateChain = await upgrades.deployProxy(RollupStateChain, [addressManager.address, config.fraudProofWindow]);
    console.log("sent RollupStateChain deploy tx, %s", rollupStateChain.deployTransaction.hash);

    /* deploy challenge contracts */
    const Challenge = await ethers.getContractFactory("Challenge");
    const challenge = await Challenge.deploy();
    console.log("sent Challenge deploy tx, %s", challenge.deployTransaction.hash);
    // implementation must be deployed
    await challenge.deployed();
    console.log("Challenge deployed: %s", challenge.address);
    const UpgradeableBeacon = await ethers.getContractFactory("UpgradeableBeacon");
    const challengeBeacon = await UpgradeableBeacon.deploy(challenge.address);
    console.log("sent UpgradeableBeacon deploy tx, %s", challengeBeacon.deployTransaction.hash);
    const ChallengeFactory = await ethers.getContractFactory("ChallengeFactory");
    const challengeFactory = await upgrades.deployProxy(ChallengeFactory, [addressManager.address,
        challengeBeacon.address, config.blockLimitPerRound, ethers.utils.parseEther(config.challengerDeposit)
    ]);
    console.log("sent ChallengeFactory deploy tx, %s", challengeFactory.deployTransaction.hash);

    const DAO = await ethers.getContractFactory("DAO");
    const dao = await upgrades.deployProxy(DAO, []);
    console.log("sent DAO deploy tx, %s", dao.deployTransaction.hash);

    const StakingManager = await ethers.getContractFactory("StakingManager");
    const stakingManager = await upgrades.deployProxy(StakingManager, [dao.address, challengeFactory.address,
        rollupStateChain.address, feeToken.address, ethers.utils.parseEther(config.stakingPrice)
    ]);
    console.log("sent StakingManager deploy tx, %s", stakingManager.deployTransaction.hash);

    const RollupInputChain = await ethers.getContractFactory("RollupInputChain");
    const rollupInputChain = await upgrades.deployProxy(RollupInputChain, [addressManager.address, config.maxTxGasLimit,
        config.maxCrossLayerTxGasLimit, config.l2ChainId
    ]);
    console.log("sent RollupInputChain deploy tx, %s", rollupInputChain.deployTransaction.hash);

    const ChainStorageContainer = await ethers.getContractFactory("ChainStorageContainer");
    const stateStorageContainer = await upgrades.deployProxy(ChainStorageContainer,
        [config.addressName.ROLLUP_STATE_CHAIN, addressManager.address]);
    console.log("sent stateStorageContainer deploy tx, %s", stateStorageContainer.deployTransaction.hash);
    const inputStorageContainer = await upgrades.deployProxy(ChainStorageContainer,
        [config.addressName.ROLLUP_INPUT_CHAIN, addressManager.address]);
    console.log("sent inputStorageContainer deploy tx, %s", inputStorageContainer.deployTransaction.hash);

    /* deploy state transition */
    const MachineState = await ethers.getContractFactory("MachineState");
    const machineState = await MachineState.deploy();
    const StateTransition = await ethers.getContractFactory("StateTransition");
    const stateTransition = await upgrades.deployProxy(StateTransition, [config.imageStateRoot, addressManager.address,
        machineState.address
    ]);
    console.log("sent StateTransition deploy tx, %s", stateTransition.deployTransaction.hash);

    /* deploy bridge */
    const L1StandardBridge = await ethers.getContractFactory("L1StandardBridge");
    const l1StandardBridge = await upgrades.deployProxy(L1StandardBridge, [l1CrossLayerWitness.address, config.l2TokenBridge]);
    console.log("sent L1StandardBridge deploy tx, %s", l1StandardBridge.deployTransaction.hash);

    await addressManager.deployed();
    console.log("AddressManager deployed: %s", addressManager.address);
    /* config address manager */
    await addressManager.setAddress(config.addressName.ROLLUP_INPUT_CHAIN, rollupInputChain.address);
    await addressManager.setAddress(config.addressName.STAKING_MANAGER, stakingManager.address);
    await addressManager.setAddress(config.addressName.ROLLUP_STATE_CHAIN_CONTAINER, stateStorageContainer.address);
    await addressManager.setAddress(config.addressName.ROLLUP_INPUT_CHAIN_CONTAINER, inputStorageContainer.address);
    await addressManager.setAddress(config.addressName.ROLLUP_STATE_CHAIN, rollupStateChain.address);
    await addressManager.setAddress(config.addressName.L1_CROSS_LAYER_WITNESS, l1CrossLayerWitness.address);
    await addressManager.setAddress(config.addressName.L2_CROSS_LAYER_WITNESS, config.l2CrossLayerWitness);
    await addressManager.setAddress(config.addressName.DAO, dao.address);
    await addressManager.setAddress(config.addressName.CHALLENGE_FACTORY, challengeFactory.address);
    await addressManager.setAddress(config.addressName.STATE_TRANSITION, stateTransition.address);
    await addressManager.setAddress(config.addressName.L1_STANDARD_BRIDGE, l1StandardBridge.address);
    await addressManager.setAddress(config.addressName.CHALLENGE_BEACON, challengeBeacon.address);
    await addressManager.setAddress(config.addressName.FEE_TOKEN, feeToken.address);
    await addressManager.setAddress(config.addressName.MACHINE_STATE, machineState.address);

    /* wait contracts deployed */
    await dao.deployed();
    console.log("dao deployed: %s", challenge.address);
    await l1CrossLayerWitness.deployed();
    console.log("L1CrossLayerWitness deployed: %s", challenge.address);
    await feeToken.deployed();
    console.log("FeeToken deployed: %s", challenge.address);
    await rollupStateChain.deployed();
    console.log("RollupStateChain deployed: %s", challenge.address);
    await challengeBeacon.deployed();
    console.log("ChallengeBeacon deployed: %s", challenge.address);
    await challengeFactory.deployed();
    console.log("ChallengeFactory deployed: %s", challenge.address);
    await stakingManager.deployed();
    console.log("StakingManager deployed: %s", challenge.address);
    await rollupInputChain.deployed();
    console.log("RollupInputChain deployed: %s", challenge.address);
    await stateStorageContainer.deployed();
    console.log("stateStorageContainer deployed: %s", challenge.address);
    await inputStorageContainer.deployed();
    console.log("inputStorageContainer deployed: %s", challenge.address);
    await machineState.deployed();
    console.log("machineState deployed: %s", machineState.address);
    await stateTransition.deployed();
    console.log("stateTransition deployed: %s", stateTransition.address);
    await l1StandardBridge.deployed();
    console.log("l1StandardBridge deployed: %s", l1StandardBridge.address);

    const addresses = {
        DAO: dao.address,
        AddressManager: addressManager.address,
        L1CrossLayerWitness: l1CrossLayerWitness.address,
        FeeToken: feeToken.address,
        RollupStateChain: rollupStateChain.address,
        ChallengeLogic: challenge.address,
        ChallengeBeacon: challengeBeacon.address,
        ChallengeFactory: challengeFactory.address,
        StakingManager: stakingManager.address,
        RollupInputChain: rollupInputChain.address,
        StateChainStorage: stateStorageContainer.address,
        InputChainStorage: inputStorageContainer.address,
        MachineState: machineState.address,
        StateTransition: inputStorageContainer.address,
        L1StandardBridge: l1StandardBridge.address,
    }
    const fs = require('fs/promises');
    const filedata = JSON.stringify(addresses, "", " ");
    await fs.writeFile('./l1-contracts.json', filedata, err => {
        console.error(err);
        process.exit(1)
    });
    console.log('contracts deployed', JSON.stringify(addresses));
}

main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });