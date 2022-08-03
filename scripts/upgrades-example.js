async function main() {
    // deploy
    const L1CrossLayerWitness = await ethers.getContractFactory('L1CrossLayerWitness');
    const l1CrossLayerWitness = await upgrades.deployProxy(L1CrossLayerWitness, ['0x0000000000000000000000000000000000000000']);
    console.log('deploy L1CrossLayerWitness at', l1CrossLayerWitness.address);

    // upgrades
    const l1CrossLayerWitnessV2 = await upgrades.upgradeProxy(l1CrossLayerWitness.address, L1CrossLayerWitness);
    console.log('upgraded, tx is', l1CrossLayerWitnessV2.deployTransaction.hash);
}

main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });