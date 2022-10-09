async function depositEthToL2() {
    let value = ethers.utils.parseEther('1');
    let L1StandardBridge = await ethers.getContractFactory('L1StandardBridge')
    const cfg = require('../rollup-config.json');
    let l1StandardBridge = L1StandardBridge.attach(cfg.L1Addresses.L1StandardBridge);
    let tx = await l1StandardBridge.depositETH(Buffer.from(''), {
        value: value
    });
    console.log(tx.hash);
}

depositEthToL2().then();