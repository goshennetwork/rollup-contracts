require("@nomiclabs/hardhat-waffle");
require('@openzeppelin/hardhat-upgrades');
const env = require('./.env.json');


// You need to export an object to set up your config
// Go to https://hardhat.org/config/ to learn more

const PRIV_1 = env.PRIVATE_KEY_1;

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
    solidity: {
        version: '0.8.13',
        settings: {
            optimizer: {
                enabled: true,
                runs: 999999
            }
        }
    },
    metadata: {
        // do not include the metadata hash, since this is machine dependent
        // and we want all generated code to be deterministic
        // https://docs.soliditylang.org/en/v0.8.6/metadata.html
        bytecodeHash: 'none'
    },
    networks: {
        hardhat: {
            allowUnlimitedContractSize: false
        },
        testnet: {
            url: 'http://172.168.3.70:8501',
            accounts: [`0x${PRIV_1}`]
        },
        mumbai: {
            url: 'https://rpc.ankr.com/polygon_mumbai',
            accounts: { // this mnemonic is invalid checksum
                mnemonic: env.mnemonic
            },
        },
        kavatest: {
            url: 'https://evm.testnet.kava.io',
            accounts: { // this mnemonic is invalid checksum
                mnemonic: env.mnemonic
            },
        },
        l2dev: {
            url: 'http://192.168.6.237:23333',
            accounts: { // this mnemonic is invalid checksum
                mnemonic: env.mnemonic
            },
        }
    }
};
