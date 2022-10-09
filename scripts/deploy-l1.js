const config = require("./config/config.json");

async function main() {
    const decimals = 18;
    const AddressManager = await ethers.getContractFactory("AddressManager");
    const addressManager = await upgrades.deployProxy(AddressManager, []);
    console.log("sent AddressManager deploy tx, %s", addressManager.deployTransaction.hash);

    const L1CrossLayerWitness = await ethers.getContractFactory("L1CrossLayerWitness");
    const l1CrossLayerWitness = await upgrades.deployProxy(L1CrossLayerWitness, [], {
        initializer: false
    });
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
    const rollupStateChain = await upgrades.deployProxy(RollupStateChain, [], {
        initializer: false
    });
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
    const challengeFactory = await upgrades.deployProxy(ChallengeFactory, [], {
        initializer: false
    });
    console.log("sent ChallengeFactory deploy tx, %s", challengeFactory.deployTransaction.hash);
    const signers = await ethers.getSigners();
    const dao = await signers[0].getAddress();
    console.log("set dao address, %s", dao);
    const Whitelist = await ethers.getContractFactory("Whitelist");
    const whitelist = await upgrades.deployProxy(Whitelist, [], {
        initializer: false
    });
    console.log("sent Whitelist deploy tx, %s", whitelist.deployTransaction.hash);

    const StakingManager = await ethers.getContractFactory("StakingManager");
    const stakingManager = await upgrades.deployProxy(StakingManager, [], {
        initializer: false
    });
    console.log("sent StakingManager deploy tx, %s", stakingManager.deployTransaction.hash);

    const RollupInputChain = await ethers.getContractFactory("RollupInputChain");
    const rollupInputChain = await upgrades.deployProxy(RollupInputChain, [], {
        initializer: false
    });
    console.log("sent RollupInputChain deploy tx, %s", rollupInputChain.deployTransaction.hash);

    const ChainStorageContainer = await ethers.getContractFactory("ChainStorageContainer");
    const stateStorageContainer = await upgrades.deployProxy(ChainStorageContainer, [], {
        initializer: false
    });
    console.log("sent stateStorageContainer deploy tx, %s", stateStorageContainer.deployTransaction.hash);
    const inputStorageContainer = await upgrades.deployProxy(ChainStorageContainer, [], {
        initializer: false
    });
    console.log("sent inputStorageContainer deploy tx, %s", inputStorageContainer.deployTransaction.hash);

    /* deploy state transition */
    const MachineState = await ethers.getContractFactory("MachineState");
    const machineState = await MachineState.deploy();
    const StateTransition = await ethers.getContractFactory("StateTransition");
    const stateTransition = await upgrades.deployProxy(StateTransition, [], {
        initializer: false
    });
    console.log("sent StateTransition deploy tx, %s", stateTransition.deployTransaction.hash);

    /* deploy bridge */
    const L1StandardBridge = await ethers.getContractFactory("L1StandardBridge");
    const l1StandardBridge = await upgrades.deployProxy(L1StandardBridge, [], {
        initializer: false
    });
    console.log("sent L1StandardBridge deploy tx, %s", l1StandardBridge.deployTransaction.hash);

    await addressManager.deployed();
    console.log("AddressManager deployed: %s", addressManager.address);
    /* config address manager */
    const names = [
        config.addressName.ROLLUP_INPUT_CHAIN,
        config.addressName.STAKING_MANAGER,
        config.addressName.ROLLUP_STATE_CHAIN_CONTAINER,
        config.addressName.ROLLUP_INPUT_CHAIN_CONTAINER,
        config.addressName.ROLLUP_STATE_CHAIN,
        config.addressName.L1_CROSS_LAYER_WITNESS,
        config.addressName.L2_CROSS_LAYER_WITNESS,
        config.addressName.DAO,
        config.addressName.CHALLENGE_FACTORY,
        config.addressName.STATE_TRANSITION,
        config.addressName.L1_STANDARD_BRIDGE,
        config.addressName.CHALLENGE_BEACON,
        config.addressName.FEE_TOKEN,
        config.addressName.MACHINE_STATE,
        config.addressName.WHITELIST
    ];
    const addrs = [
        rollupInputChain.address,
        stakingManager.address,
        stateStorageContainer.address,
        inputStorageContainer.address,
        rollupStateChain.address,
        l1CrossLayerWitness.address,
        config.l2CrossLayerWitness,
        dao,
        challengeFactory.address,
        stateTransition.address,
        l1StandardBridge.address,
        challengeBeacon.address,
        feeToken.address,
        machineState.address,
        whitelist.address
    ];
    await addressManager.setAddressBatch(names, addrs);

    /* wait contracts deployed */
    await whitelist.deployed();
    console.log("whitelist deployed: %s", whitelist.address);
    await l1CrossLayerWitness.deployed();
    console.log("L1CrossLayerWitness deployed: %s", l1CrossLayerWitness.address);
    await feeToken.deployed();
    console.log("FeeToken deployed: %s", feeToken.address);
    await rollupStateChain.deployed();
    console.log("RollupStateChain deployed: %s", rollupStateChain.address);
    await challengeBeacon.deployed();
    console.log("ChallengeBeacon deployed: %s", challengeBeacon.address);
    await challengeFactory.deployed();
    console.log("ChallengeFactory deployed: %s", challengeFactory.address);
    await stakingManager.deployed();
    console.log("StakingManager deployed: %s", stakingManager.address);
    await rollupInputChain.deployed();
    console.log("RollupInputChain deployed: %s", rollupInputChain.address);
    await stateStorageContainer.deployed();
    console.log("stateStorageContainer deployed: %s", stateStorageContainer.address);
    await inputStorageContainer.deployed();
    console.log("inputStorageContainer deployed: %s", inputStorageContainer.address);
    await machineState.deployed();
    console.log("machineState deployed: %s", machineState.address);
    await stateTransition.deployed();
    console.log("stateTransition deployed: %s", stateTransition.address);
    await l1StandardBridge.deployed();
    console.log("l1StandardBridge deployed: %s", l1StandardBridge.address);

    /* initialize contract */
    let tx = await l1CrossLayerWitness.initialize(addressManager.address);
    console.log("l1CrossLayerWitness initialized, tx: %s", tx.hash);
    tx = await rollupStateChain.initialize(addressManager.address, config.fraudProofWindow);
    console.log("rollupStateChain initialized, tx: %s", tx.hash);
    tx = await challengeFactory.initialize(addressManager.address, challengeBeacon.address, config.blockLimitPerRound,
        ethers.utils.parseEther(config.challengerDeposit));
    console.log("challengeFactory initialized, tx: %s", tx.hash);
    tx = await whitelist.initialize(addressManager.address);
    console.log("whitelist initialized, tx: %s", tx.hash);
    tx = await stakingManager.initialize(addressManager.address, ethers.utils.parseEther(config.stakingPrice));
    console.log("stakingManager initialized, tx: %s", tx.hash);
    tx = await rollupInputChain.initialize(addressManager.address, config.maxTxGasLimit, config.maxCrossLayerTxGasLimit,
        config.l2ChainId);
    console.log("rollupInputChain initialized, tx: %s", tx.hash);
    tx = await stateStorageContainer.initialize(config.addressName.ROLLUP_STATE_CHAIN, addressManager.address);
    console.log("stateStorageContainer initialized, tx: %s", tx.hash);
    tx = await inputStorageContainer.initialize(config.addressName.ROLLUP_INPUT_CHAIN, addressManager.address);
    console.log("inputStorageContainer initialized, tx: %s", tx.hash);
    tx = await stateTransition.initialize(config.imageStateRoot, addressManager.address, machineState.address);
    console.log("stateTransition initialized, tx: %s", tx.hash);
    tx = await l1StandardBridge.initialize(l1CrossLayerWitness.address, config.l2TokenBridge);
    console.log("l1StandardBridge initialized, tx: %s", tx.hash);

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
        StateTransition: stateTransition.address,
        L1StandardBridge: l1StandardBridge.address,
        WhiteList: whitelist.address,
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