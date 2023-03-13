async function main() {
    // deploy
    const L2FeeCollector = await ethers.getContractFactory('L2FeeCollector');
    const l2FeeCollector = await L2FeeCollector.attach("0xfee0000000000000000000000000000000000fee");

    const balance = await ethers.provider.getBalance(l2FeeCollector.address);
    console.log("fee collector balance: ", balance);
    const tx = await l2FeeCollector.withdrawEth(balance);
    console.log("withdrawEth, txHash: ", tx);
}
main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });