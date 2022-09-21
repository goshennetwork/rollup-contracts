async function depositEthToL2() {
    let value = ethers.utils.parseEther('1');
    let L1StandardBridge = await ethers.getContractFactory('L1StandardBridge')
    const l1Contracts = require('../l1-contracts.json');
    let l1StandardBridge = L1StandardBridge.attach(l1Contracts.L1StandardBridge);
    let tx = await l1StandardBridge.depositETH(Buffer.from(''), {value: value});
    console.log(tx.hash);
}

depositEthToL2().then();