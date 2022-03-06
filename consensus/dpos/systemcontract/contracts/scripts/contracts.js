const Web3 = require('web3');
const { rpcs, genesisConfig, accounts, accountKeys } = require("../config.js")

// web3 rpc config; 
// different from hardhat console configuration
const web3 = new Web3(rpcs.node2); // node2

// system contract abi
const ValidatorsJson = require("../artifacts/contracts/Validators.sol/Validators.json")
const ProposalsJson = require("../artifacts/contracts/Proposals.sol/Proposals.json")
const NodeVotesJson = require("../artifacts/contracts/NodeVotes.sol/NodeVotes.json")
const SystemRewardsJson = require("../artifacts/contracts/SystemRewards.sol/SystemRewards.json")

// system contract instance
const validators = new web3.eth.Contract(ValidatorsJson.abi, genesisConfig.ValidatorsAddr)
const proposals = new web3.eth.Contract(ProposalsJson.abi, genesisConfig.ProposalsAddr)
const nodeVotes = new web3.eth.Contract(NodeVotesJson.abi, genesisConfig.NodeVotesAddr)
const systemRewards = new web3.eth.Contract(SystemRewardsJson.abi, genesisConfig.SystemRewardsAddr)

accountKeys.forEach(element => { web3.eth.accounts.wallet.add(element) });

module.exports = {
    web3,
    validators,
    proposals,
    nodeVotes,
    systemRewards,
    accountKeys,
    accounts,
}