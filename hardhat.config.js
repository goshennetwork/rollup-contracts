require("@nomiclabs/hardhat-waffle");
require('@openzeppelin/hardhat-upgrades');

// You need to export an object to set up your config
// Go to https://hardhat.org/config/ to learn more

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
    }
};
