// load env: export NODE_ENV=dev
require('dotenv-flow').config({ path: 'envs/', default_node_env: 'dev', silent: true });

// system contract address
const genesisConfig = {
    ValidatorsAddr: "0x0000000000000000000000000000000000fff001",
    ProposalsAddr: "0x0000000000000000000000000000000000fff002",
    NodeVotesAddr: "0x0000000000000000000000000000000000fff003",
    SystemRewardsAddr: "0x0000000000000000000000000000000000fff004",
    AddressListAddr: "0x0000000000000000000000000000000000fff005",
}
// RPC URLS
const rpcs = {
    localhost: "http://localhost:8545"
}

const accounts = process.env.ACCOUNTS ? process.env.ACCOUNTS.split(",") : []
// add accounts to wallet
// default: []
const accountKeys = process.env.PRIVATE_KEYS ? process.env.PRIVATE_KEYS.split(",") : []
// accountKeys.push("your personal privateKey")

module.exports = { genesisConfig, accountKeys, accounts, rpcs }