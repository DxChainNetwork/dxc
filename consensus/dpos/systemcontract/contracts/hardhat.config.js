require("@nomiclabs/hardhat-waffle");
require("@nomiclabs/hardhat-truffle5");
require("@nomiclabs/hardhat-web3");

const Web3 = require("web3");
const fs = require("fs");
const { genesisConfig, rpcs, accounts, accountKeys } = require("./config.js")

// `hh accounts` 
task("accounts", "Prints accounts", async (_, { web3 }) => { console.log(await web3.eth.getAccounts()); });

const compiled = fs.existsSync("./artifacts")

// extend hre environment
extendEnvironment((hre) => {
  hre.Web3 = Web3
  hre.eth = hre.web3.eth;
  hre.genesisConfig = genesisConfig;

  if (compiled) {
    const ValidatorsJson = require("./artifacts/contracts/Validators.sol/Validators.json")
    const ProposalsJson = require("./artifacts/contracts/Proposals.sol/Proposals.json")
    const NodeVotesJson = require("./artifacts/contracts/NodeVotes.sol/NodeVotes.json")
    const SystemRewardsJson = require("./artifacts/contracts/SystemRewards.sol/SystemRewards.json")
    hre.validators = new hre.web3.eth.Contract(ValidatorsJson.abi, genesisConfig.ValidatorsAddr)
    hre.proposals = new hre.web3.eth.Contract(ProposalsJson.abi, genesisConfig.ProposalsAddr)
    hre.nodeVotes = new hre.web3.eth.Contract(NodeVotesJson.abi, genesisConfig.NodeVotesAddr)
    hre.systemRewards = new hre.web3.eth.Contract(SystemRewardsJson.abi, genesisConfig.SystemRewardsAddr)
  }

  // address alias
  hre.accounts = accounts
});

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  defaultNetwork: "localhost",
  networks: {
    localhost: {
      url: rpcs.localhost,
      accounts: accountKeys,
      // accounts: { mnemonic: mnemonic }
    },
  },

  solidity: {
    compilers: [
      {
        version: "0.8.11",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200
          }
        }
      }
    ],
  },

  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  },

  etherscan: {
    apiKey: process.env.API_KEY
  },
};
