package systemcontract

import (
	"github.com/DxChainNetwork/dxc/accounts/abi"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/params"
	"math/big"
	"strings"
)

// ValidatorsABI contains all methods to interactive with validator contracts.
const ValidatorsABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"},{"indexed":false,"internalType":"uint256","name":"_deposit","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_rate","type":"uint256"}],"name":"LogAddValidator","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"}],"name":"LogKickoutValidator","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"}],"name":"LogRedeemValidator","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"_validatorCount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_deposit","type":"uint256"}],"name":"LogTryElect","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"},{"indexed":false,"internalType":"uint256","name":"_lockEnd","type":"uint256"}],"name":"LogUnstakeValidator","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"},{"indexed":false,"internalType":"uint256","name":"_deposit","type":"uint256"}],"name":"LogUpdateValidatorDeposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"},{"indexed":false,"internalType":"uint8","name":"_rate","type":"uint8"}],"name":"LogUpdateValidatorRate","type":"event"},{"inputs":[],"name":"BLACK_HOLE_ADDRESS","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"CancelQueueValidatorsLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"EPOCH_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_PROPOSAL_DETAIL_LENGTH","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_PUNISH_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_RATE","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_RATE_OF_CHANGE","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATORS_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_L4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MEDIUM_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_DEPOSIT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_RATE","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"PROPOSAL_DURATION","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"RATE_SET_LOCK_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV1_TO_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV2_TO_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV3_TO_LV4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV4_TO_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_OVER_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_UNDER_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"VALIDATOR_UNSTAKE_LOCK_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"VOTE_CANCEL_BLOCK","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"},{"internalType":"uint256","name":"_deposit","type":"uint256"},{"internalType":"uint8","name":"_rate","type":"uint8"}],"name":"addValidatorFromProposal","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"_voter","type":"address"},{"internalType":"address","name":"_val","type":"address"},{"internalType":"uint256","name":"_votes","type":"uint256"}],"name":"cancelVoteValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"curEpochValidators","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"curEpochValidatorsIdMap","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"currentEpoch","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"effictiveValsLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCancelQueueValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurEpochValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getEffictiveValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getInvalidValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"},{"internalType":"uint256","name":"page","type":"uint256"},{"internalType":"uint256","name":"size","type":"uint256"}],"name":"getVoters","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_proposal","type":"address"},{"internalType":"address","name":"_sysReward","type":"address"},{"internalType":"address","name":"_nodeVote","type":"address"},{"internalType":"address","name":"_initVal","type":"address"},{"internalType":"uint256","name":"_initDeposit","type":"uint256"},{"internalType":"uint8","name":"_initRate","type":"uint8"}],"name":"initialize","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"invalidValsLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"isEffictiveValidator","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"kickoutValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"nodeVote","outputs":[{"internalType":"contract INodeVote","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"proposals","outputs":[{"internalType":"contract IProposals","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"redeem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"sysRewards","outputs":[{"internalType":"contract ISystemRewards","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalDeposit","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"tryElect","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unstake","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_deposit","type":"uint256"}],"name":"updateValidatorDeposit","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"uint8","name":"_rate","type":"uint8"}],"name":"updateValidatorRate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"validatorToVotersLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"validators","outputs":[{"internalType":"enum Validators.ValidatorStatus","name":"status","type":"uint8"},{"internalType":"uint256","name":"deposit","type":"uint256"},{"internalType":"uint8","name":"rate","type":"uint8"},{"internalType":"uint256","name":"totalVotes","type":"uint256"},{"internalType":"uint256","name":"unstakeLockingEndBlock","type":"uint256"},{"internalType":"uint256","name":"rateSettLockingEndBlock","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_voter","type":"address"},{"internalType":"address","name":"_val","type":"address"},{"internalType":"uint256","name":"_votes","type":"uint256"}],"name":"voteValidator","outputs":[],"stateMutability":"payable","type":"function"}]`

const ProposalsABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"id","type":"bytes32"},{"indexed":true,"internalType":"address","name":"proposer","type":"address"},{"indexed":false,"internalType":"uint256","name":"block","type":"uint256"}],"name":"LogCancelProposal","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"id","type":"bytes32"},{"indexed":true,"internalType":"address","name":"guarantee","type":"address"},{"indexed":false,"internalType":"uint256","name":"block","type":"uint256"}],"name":"LogGuarantee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"id","type":"bytes32"},{"indexed":true,"internalType":"address","name":"proposer","type":"address"},{"indexed":false,"internalType":"uint256","name":"block","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"deposit","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"rate","type":"uint256"}],"name":"LogInitProposal","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"id","type":"bytes32"},{"indexed":true,"internalType":"address","name":"proposer","type":"address"},{"indexed":false,"internalType":"uint256","name":"block","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"deposit","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"rate","type":"uint256"}],"name":"LogUpdateProposal","type":"event"},{"inputs":[],"name":"BLACK_HOLE_ADDRESS","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"EPOCH_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_PROPOSAL_DETAIL_LENGTH","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_PUNISH_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_RATE","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_RATE_OF_CHANGE","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATORS_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_L4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MEDIUM_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_DEPOSIT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_RATE","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"PROPOSAL_DURATION","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"RATE_SET_LOCK_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV1_TO_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV2_TO_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV3_TO_LV4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV4_TO_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_OVER_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_UNDER_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"VALIDATOR_UNSTAKE_LOCK_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"VOTE_CANCEL_BLOCK","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"val","type":"address"}],"name":"addressProposalCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"val","type":"address"},{"internalType":"uint256","name":"page","type":"uint256"},{"internalType":"uint256","name":"size","type":"uint256"}],"name":"addressProposalSets","outputs":[{"internalType":"bytes4[]","name":"","type":"bytes4[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"val","type":"address"},{"internalType":"uint256","name":"page","type":"uint256"},{"internalType":"uint256","name":"size","type":"uint256"}],"name":"addressProposals","outputs":[{"components":[{"internalType":"bytes4","name":"id","type":"bytes4"},{"internalType":"address","name":"proposer","type":"address"},{"internalType":"enum Proposals.ProposalType","name":"pType","type":"uint8"},{"internalType":"uint256","name":"deposit","type":"uint256"},{"internalType":"uint8","name":"rate","type":"uint8"},{"internalType":"string","name":"details","type":"string"},{"internalType":"uint256","name":"initBlock","type":"uint256"},{"internalType":"address","name":"guarantee","type":"address"},{"internalType":"uint256","name":"updateBlock","type":"uint256"},{"internalType":"enum Proposals.ProposalStatus","name":"status","type":"uint8"}],"internalType":"struct Proposals.ProposalInfo[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"page","type":"uint256"},{"internalType":"uint256","name":"size","type":"uint256"}],"name":"allProposalSets","outputs":[{"internalType":"bytes4[]","name":"","type":"bytes4[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"page","type":"uint256"},{"internalType":"uint256","name":"size","type":"uint256"}],"name":"allProposals","outputs":[{"components":[{"internalType":"bytes4","name":"id","type":"bytes4"},{"internalType":"address","name":"proposer","type":"address"},{"internalType":"enum Proposals.ProposalType","name":"pType","type":"uint8"},{"internalType":"uint256","name":"deposit","type":"uint256"},{"internalType":"uint8","name":"rate","type":"uint8"},{"internalType":"string","name":"details","type":"string"},{"internalType":"uint256","name":"initBlock","type":"uint256"},{"internalType":"address","name":"guarantee","type":"address"},{"internalType":"uint256","name":"updateBlock","type":"uint256"},{"internalType":"enum Proposals.ProposalStatus","name":"status","type":"uint8"}],"internalType":"struct Proposals.ProposalInfo[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes4","name":"id","type":"bytes4"}],"name":"cancelProposal","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"currentEpoch","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes4","name":"id","type":"bytes4"}],"name":"guarantee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"enum Proposals.ProposalType","name":"pType","type":"uint8"},{"internalType":"uint8","name":"rate","type":"uint8"},{"internalType":"string","name":"details","type":"string"}],"name":"initProposal","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"_validator","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"proposalCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes4","name":"","type":"bytes4"}],"name":"proposalInfos","outputs":[{"internalType":"bytes4","name":"id","type":"bytes4"},{"internalType":"address","name":"proposer","type":"address"},{"internalType":"enum Proposals.ProposalType","name":"pType","type":"uint8"},{"internalType":"uint256","name":"deposit","type":"uint256"},{"internalType":"uint8","name":"rate","type":"uint8"},{"internalType":"string","name":"details","type":"string"},{"internalType":"uint256","name":"initBlock","type":"uint256"},{"internalType":"address","name":"guarantee","type":"address"},{"internalType":"uint256","name":"updateBlock","type":"uint256"},{"internalType":"enum Proposals.ProposalStatus","name":"status","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"proposals","outputs":[{"internalType":"bytes4","name":"","type":"bytes4"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes4","name":"id","type":"bytes4"},{"internalType":"uint8","name":"rate","type":"uint8"},{"internalType":"uint256","name":"deposit","type":"uint256"},{"internalType":"string","name":"details","type":"string"}],"name":"updateProposal","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"validators","outputs":[{"internalType":"contract IValidators","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

const NodeVotesABI = `[{"anonymous": false,"inputs": [{"indexed": true,"internalType": "address","name": "voter","type": "address"},{"indexed": true,"internalType": "address","name": "validator","type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "votes",
          "type": "uint256"
        }
      ],
      "name": "LogCancelVote",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "voter",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "validator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "reward",
          "type": "uint256"
        }
      ],
      "name": "LogEarn",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "voter",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "votes",
          "type": "uint256"
        }
      ],
      "name": "LogRedeem",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "voter",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "validator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "votes",
          "type": "uint256"
        }
      ],
      "name": "LogVote",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "BLACK_HOLE_ADDRESS",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "EPOCH_BLOCKS",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_LEVEL_VALIDATOR_COUNT",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_PROPOSAL_DETAIL_LENGTH",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_PUNISH_COUNT",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_RATE",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_RATE_OF_CHANGE",
      "outputs": [
        {
          "internalType": "uint16",
          "name": "",
          "type": "uint16"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_VALIDATORS_COUNT",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_VALIDATOR_COUNT_L4",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_VALIDATOR_COUNT_LV1",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_VALIDATOR_COUNT_LV2",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MAX_VALIDATOR_COUNT_LV3",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MEDIUM_LEVEL_VALIDATOR_COUNT",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MIN_DEPOSIT",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MIN_LEVEL_VALIDATOR_COUNT",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MIN_RATE",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "PROPOSAL_DURATION",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "RATE_SET_LOCK_BLOCKS",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "REWARD_DEPOSIT_FROM_LV1_TO_LV2",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "REWARD_DEPOSIT_FROM_LV2_TO_LV3",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "REWARD_DEPOSIT_FROM_LV3_TO_LV4",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "REWARD_DEPOSIT_FROM_LV4_TO_LV5",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "REWARD_DEPOSIT_OVER_LV5",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "REWARD_DEPOSIT_UNDER_LV1",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "TOTAL_DEPOSIT_LV1",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "TOTAL_DEPOSIT_LV2",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "TOTAL_DEPOSIT_LV3",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "TOTAL_DEPOSIT_LV4",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "TOTAL_DEPOSIT_LV5",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "VALIDATOR_UNSTAKE_LOCK_BLOCKS",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "VOTE_CANCEL_BLOCK",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        }
      ],
      "name": "cancelVote",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "page",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "size",
          "type": "uint256"
        }
      ],
      "name": "cancelVoteValidatorList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        }
      ],
      "name": "cancelVoteValidatorListLength",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "currentEpoch",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        }
      ],
      "name": "earn",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_validator",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_sysReward",
          "type": "address"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_voter",
          "type": "address"
        }
      ],
      "name": "pendingRedeem",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_voter",
          "type": "address"
        }
      ],
      "name": "pendingVoteReward",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address[]",
          "name": "vals",
          "type": "address[]"
        }
      ],
      "name": "redeem",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "page",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "size",
          "type": "uint256"
        }
      ],
      "name": "redeemInfo",
      "outputs": [
        {
          "components": [
            {
              "internalType": "address",
              "name": "validator",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "votes",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "redeemBlock",
              "type": "uint256"
            }
          ],
          "internalType": "struct NodeVotes.RedeemVoterInfo[]",
          "name": "",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "sysRewards",
      "outputs": [
        {
          "internalType": "contract ISystemRewards",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "validators",
      "outputs": [
        {
          "internalType": "contract IValidators",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        }
      ],
      "name": "vote",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "page",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "size",
          "type": "uint256"
        }
      ],
      "name": "voteList",
      "outputs": [
        {
          "components": [
            {
              "internalType": "address",
              "name": "validator",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "votes",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "rewardDebt",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "updateRewardEpoch",
              "type": "uint256"
            }
          ],
          "internalType": "struct NodeVotes.VoteInfo[]",
          "name": "",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_val",
          "type": "address"
        }
      ],
      "name": "voteListLength",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "stateMutability": "payable",
      "type": "receive"
    }
  ]`

const SystemRewardsABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"},{"indexed":false,"internalType":"uint256","name":"_valReward","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_delegatorReward","type":"uint256"}],"name":"LogDistributeBlockReward","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"},{"indexed":false,"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"LogEarnValidatorReward","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_val","type":"address"},{"indexed":false,"internalType":"bool","name":"_kickout","type":"bool"},{"indexed":false,"internalType":"uint256","name":"_reward","type":"uint256"}],"name":"LogPunish","type":"event"},{"inputs":[],"name":"BLACK_HOLE_ADDRESS","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"EPOCH_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_PROPOSAL_DETAIL_LENGTH","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_PUNISH_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_RATE","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_RATE_OF_CHANGE","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATORS_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_L4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_VALIDATOR_COUNT_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MEDIUM_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_DEPOSIT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_LEVEL_VALIDATOR_COUNT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MIN_RATE","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"PROPOSAL_DURATION","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"RATE_SET_LOCK_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV1_TO_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV2_TO_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV3_TO_LV4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_FROM_LV4_TO_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_OVER_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REWARD_DEPOSIT_UNDER_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV2","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV3","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV4","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOTAL_DEPOSIT_LV5","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"VALIDATOR_UNSTAKE_LOCK_BLOCKS","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"VOTE_CANCEL_BLOCK","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"currentEpoch","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_reward","type":"uint256"}],"name":"distributeBlockReward","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"earnValReward","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"earnValRewardFromValidatorC","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"epochs","outputs":[{"internalType":"uint256","name":"blockReward","type":"uint256"},{"internalType":"uint256","name":"tvl","type":"uint256"},{"internalType":"uint256","name":"validatorCount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"getValRewardEpochs","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"},{"internalType":"uint256","name":"_epoch","type":"uint256"}],"name":"getValRewardInfoByEpoch","outputs":[{"components":[{"internalType":"uint256","name":"pengingValidatorReward","type":"uint256"},{"internalType":"uint256","name":"pengingDelegatorsReward","type":"uint256"},{"internalType":"uint256","name":"totalVotes","type":"uint256"},{"internalType":"uint256","name":"cancelVotes","type":"uint256"}],"internalType":"struct SystemRewards.Reward","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_validator","type":"address"},{"internalType":"address","name":"_node","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"nodeVoteC","outputs":[{"internalType":"contract INodeVote","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"pendingValReward","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"pendingVoterReward","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"punish","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"punishInfo","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"sysRewards","outputs":[{"internalType":"uint256","name":"nextValRewardEpochIndex","type":"uint256"},{"internalType":"uint256","name":"nextDelegatorsRewardEpochIndex","type":"uint256"},{"internalType":"uint256","name":"accRewardPerVote","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"updateCancelVote","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_tvl","type":"uint256"},{"internalType":"uint256","name":"_valCount","type":"uint256"}],"name":"updateEpochWhileElect","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"},{"internalType":"uint256","name":"_deposit","type":"uint256"}],"name":"updateValDeposit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_val","type":"address"}],"name":"updateVoterReward","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"validatorC","outputs":[{"internalType":"contract IValidators","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

const SysGovABI = `[{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"finishProposalById","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint32","name":"index","type":"uint32"}],"name":"getPassedProposalByIndex","outputs":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"uint256","name":"action","type":"uint256"},{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getPassedProposalCount","outputs":[{"internalType":"uint32","name":"","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_admin","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

const AddressListABI = `[{"inputs":[],"name":"blackLastUpdatedNumber","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"devVerifyEnabled","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBlacksFrom","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBlacksTo","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint32","name":"i","type":"uint32"}],"name":"getRuleByIndex","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"},{"internalType":"uint128","name":"","type":"uint128"},{"internalType":"enum AddressList.CheckType","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"initializeV2","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_admin","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"isDeveloper","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rulesLastUpdatedNumber","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"rulesLen","outputs":[{"internalType":"uint32","name":"","type":"uint32"}],"stateMutability":"view","type":"function"}]`

// DevMappingPosition is the position of the state variable `devs`.
// Since the state variables are as follows:
//    bool public initialized;
//    bool public devVerifyEnabled;
//    address public admin;
//    address public pendingAdmin;
//
//    mapping(address => bool) private devs;
//
//    //NOTE: make sure this list is not too large!
//    address[] blacksFrom;
//    address[] blacksTo;
//    mapping(address => uint256) blacksFromMap;      // address => index+1
//    mapping(address => uint256) blacksToMap;        // address => index+1
//
//    uint256 public blackLastUpdatedNumber; // last block number when the black list is updated
//    uint256 public rulesLastUpdatedNumber;  // last block number when the rules are updated
//    // event check rules
//    EventCheckRule[] rules;
//    mapping(bytes32 => mapping(uint128 => uint256)) rulesMap;   // eventSig => checkIdx => indexInArray+1
//
// according to [Layout of State Variables in Storage](https://docs.soliditylang.org/en/v0.8.4/internals/layout_in_storage.html),
// and after optimizer enabled, the `initialized`, `enabled` and `admin` will be packed, and stores at slot 0,
// `pendingAdmin` stores at slot 1, so the position for `devs` is 2.
const DevMappingPosition = 2

var (
	BlackLastUpdatedNumberPosition = common.BytesToHash([]byte{0x07})
	RulesLastUpdatedNumberPosition = common.BytesToHash([]byte{0x08})
)

var (
	ValidatorsContractName    = "Validators"
	ProposalsContractName     = "Proposals"
	NodeVotesContractName     = "NodeVotes"
	SystemRewardsContractName = "SystemRewards"
	AddressListContractName   = "address_list"

	SysGovContractName = "governance"

	ValidatorsContractAddr    = common.HexToAddress("0x0000000000000000000000000000000000fff001")
	ProposalsContractAddr     = common.HexToAddress("0x0000000000000000000000000000000000fff002")
	NodeVotesContractAddr     = common.HexToAddress("0x0000000000000000000000000000000000fff003")
	SystemRewardsContractAddr = common.HexToAddress("0x0000000000000000000000000000000000fff004")
	AddressListContractAddr   = common.HexToAddress("0x0000000000000000000000000000000000fff005")
	SysGovContractAddr        = common.HexToAddress("0x0000000000000000000000000000000000fff006")

	// SysGovToAddr is the To address for the system governance transaction, NOT contract address
	SysGovToAddr = common.HexToAddress("0x000000000000000000000000000000000000ffff")

	abiMap map[string]abi.ABI
)

func init() {
	abiMap = make(map[string]abi.ABI, 0)
	tmpABI, _ := abi.JSON(strings.NewReader(ValidatorsABI))
	abiMap[ValidatorsContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(ProposalsABI))
	abiMap[ProposalsContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(NodeVotesABI))
	abiMap[NodeVotesContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(SystemRewardsABI))
	abiMap[SystemRewardsContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(AddressListABI))
	abiMap[AddressListContractName] = tmpABI

	tmpABI, _ = abi.JSON(strings.NewReader(SysGovABI))
	abiMap[SysGovContractName] = tmpABI

}

func GetInteractiveABI() map[string]abi.ABI {
	return abiMap
}

func GetValidatorAddr(blockNum *big.Int, config *params.ChainConfig) *common.Address {
	return &ValidatorsContractAddr
}
