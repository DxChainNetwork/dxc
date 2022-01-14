// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// package web3ext contains geth specific web3.js extensions.
package web3ext

var Modules = map[string]string{
	"admin":    AdminJs,
	"clique":   CliqueJs,
	"dpos":     DposJs,
	"ethash":   EthashJs,
	"debug":    DebugJs,
	"eth":      EthJs,
	"miner":    MinerJs,
	"net":      NetJs,
	"personal": PersonalJs,
	"rpc":      RpcJs,
	"txpool":   TxpoolJs,
	"les":      LESJs,
	"vflux":    VfluxJs,
}

const CliqueJs = `
web3._extend({
	property: 'clique',
	methods: [
		new web3._extend.Method({
			name: 'getSnapshot',
			call: 'clique_getSnapshot',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getSnapshotAtHash',
			call: 'clique_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getSigners',
			call: 'clique_getSigners',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getSignersAtHash',
			call: 'clique_getSignersAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'propose',
			call: 'clique_propose',
			params: 2
		}),
		new web3._extend.Method({
			name: 'discard',
			call: 'clique_discard',
			params: 1
		}),
		new web3._extend.Method({
			name: 'status',
			call: 'clique_status',
			params: 0
		}),
		new web3._extend.Method({
			name: 'getSigner',
			call: 'clique_getSigner',
			params: 1,
			inputFormatter: [null]
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'proposals',
			getter: 'clique_proposals'
		}),
	]
});
`

const DposJs = `
web3._extend({
	property: 'dpos',
	methods: [
		new web3._extend.Method({
			name: 'getSnapshot',
			call: 'dpos_getSnapshot',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getSnapshotAtHash',
			call: 'dpos_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getValidators',
			call: 'dpos_getValidators',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getValidatorsAtHash',
			call: 'dpos_getValidatorsAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'minDeposit',
			call: 'dpos_getMinDeposit',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'validator',
			call: 'dpos_getValidator',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'totalDeposit',
			call: 'dpos_getTotalDeposit',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'currentEpochValidators',
			call: 'dpos_getCurrentEpochValidators',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'effictiveValidators',
			call: 'dpos_getEffictiveValidators',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'cancelQueueValidators',
			call: 'dpos_getCancelQueueValidators',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'invalidValidators',
			call: 'dpos_getInvalidValidators',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'getVoters',
			call: 'dpos_getVoters',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 4
		}),
		new web3._extend.Method({
			name: 'effictiveValsLength',
			call: 'dpos_getEffictiveValsLength',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'invalidValsLength',
			call: 'dpos_getInvalidValsLength',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'cancelQueueValsLength',
			call: 'dpos_getCancelQueueValidatorsLength',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'validatorVotersLength',
			call: 'dpos_getValidatorToVotersLength',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'isEffictiveValidator',
			call: 'dpos_getIsEffictiveValidator',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'allProposalIds',
			call: 'dpos_getAllProposalSets',
			inputFormatter: [null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 3
		}),
		new web3._extend.Method({
			name: 'addressProposalIds',
			call: 'dpos_getAddressProposalSets',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 4
		}),
		new web3._extend.Method({
			name: 'allProposals',
			call: 'dpos_getAllProposals',
			inputFormatter: [null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 3
		}),
		new web3._extend.Method({
			name: 'proposal',
			call: 'dpos_getProposal',
			inputFormatter: [null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'addressProposals',
			call: 'dpos_getAddressProposals',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 4
		}),
		new web3._extend.Method({
			name: 'allProposalsCount',
			call: 'dpos_getProposalCount',
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			params: 1
		}),
		new web3._extend.Method({
			name: 'addressProposalsCount',
			call: 'dpos_getAddressProposalCount',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'pendingVoterReward',
			call: 'dpos_getPendingVoteReward',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 3
		}),
		new web3._extend.Method({
			name: 'pendingVoterRedeem',
			call: 'dpos_getPendingRedeem',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'voteListLength',
			call: 'dpos_getVoteListLength',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'cancelVoteValidatorListLength',
			call: 'dpos_getCancelVoteValidatorListLength',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'cancelVoteValidatorList',
			call: 'dpos_getCancelVoteValidatorList',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 4
		}),
		new web3._extend.Method({
			name: 'voteList',
			call: 'dpos_getVoteList',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 4
		}),
		new web3._extend.Method({
			name: 'voterRedeemInfo',
			call: 'dpos_getRedeemInfo',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,null,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 4
		}),
		new web3._extend.Method({
			name: 'epochInfo',
			call: 'dpos_getEpochInfo',
			inputFormatter: [null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'getValRewardIndex',
			call: 'dpos_getSysRewards',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'getValRewardEpochs',
			call: 'dpos_getValRewardEpochs',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'getValRewardInfoByEpoch',
			call: 'dpos_getValRewardInfoByEpoch',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,null,web3._extend.formatters.inputBlockNumberFormatter],
			params: 3
		}),
		new web3._extend.Method({
			name: 'pendingValReward',
			call: 'dpos_getPendingValReward',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'pendingValidatorToVoterReward',
			call: 'dpos_getPendingVoterReward',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,web3._extend.formatters.inputBlockNumberFormatter],
			params: 2
		}),
		new web3._extend.Method({
			name: 'initProposal',
			call: 'dpos_initProposal',
			inputFormatter: [null,null,null,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 4
		}),
		new web3._extend.Method({
			name: 'updateProposal',
			call: 'dpos_updateProposal',
			inputFormatter: [null,null,null,null,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 5
		}),
		new web3._extend.Method({
			name: 'cancelProposal',
			call: 'dpos_cancelProposal',
			inputFormatter: [null,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
		new web3._extend.Method({
			name: 'guarantee',
			call: 'dpos_guarantee',
			inputFormatter: [null,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
		new web3._extend.Method({
			name: 'updateValidatorDeposit',
			call: 'dpos_updateValidatorDeposit',
			inputFormatter: [null,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
		new web3._extend.Method({
			name: 'updateValidatorRate',
			call: 'dpos_updateValidatorRate',
			inputFormatter: [null,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
		new web3._extend.Method({
			name: 'validatorUnstake',
			call: 'dpos_validatorUnstake',
			inputFormatter: [function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 1
		}),
		new web3._extend.Method({
			name: 'validatorRedeem',
			call: 'dpos_validatorRedeem',
			inputFormatter: [function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 1
		}),
		new web3._extend.Method({
			name: 'earnValReward',
			call: 'dpos_earnValReward',
			inputFormatter: [function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 1
		}),
		new web3._extend.Method({
			name: 'earnVoterReward',
			call: 'dpos_earnVoterReward',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
		new web3._extend.Method({
			name: 'vote',
			call: 'dpos_vote',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
		new web3._extend.Method({
			name: 'cancelVote',
			call: 'dpos_cancelVote',
			inputFormatter: [web3._extend.formatters.inputAddressFormatter,function(options) {
				options = options == undefined? {} : options
				return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
		new web3._extend.Method({
			name: 'voterRedeem',
			call: 'dpos_voterRedeem',
			inputFormatter: [function(vals) {
				var formatVals = [];
				for (var i = 0; i < vals.length; i++) { formatVals.push(web3._extend.formatters.inputAddressFormatter(vals[i])); }
				return formatVals;
			}, function(options) {
					options = options == undefined? {} : options
					return web3._extend.formatters.inputCallFormatter(options)
			}],
			params: 2
		}),
	],
});
`

const EthashJs = `
web3._extend({
	property: 'ethash',
	methods: [
		new web3._extend.Method({
			name: 'getWork',
			call: 'ethash_getWork',
			params: 0
		}),
		new web3._extend.Method({
			name: 'getHashrate',
			call: 'ethash_getHashrate',
			params: 0
		}),
		new web3._extend.Method({
			name: 'submitWork',
			call: 'ethash_submitWork',
			params: 3,
		}),
		new web3._extend.Method({
			name: 'submitHashrate',
			call: 'ethash_submitHashrate',
			params: 2,
		}),
	]
});
`

const AdminJs = `
web3._extend({
	property: 'admin',
	methods: [
		new web3._extend.Method({
			name: 'addPeer',
			call: 'admin_addPeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'removePeer',
			call: 'admin_removePeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'addTrustedPeer',
			call: 'admin_addTrustedPeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'removeTrustedPeer',
			call: 'admin_removeTrustedPeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'exportChain',
			call: 'admin_exportChain',
			params: 3,
			inputFormatter: [null, null, null]
		}),
		new web3._extend.Method({
			name: 'importChain',
			call: 'admin_importChain',
			params: 1
		}),
		new web3._extend.Method({
			name: 'sleepBlocks',
			call: 'admin_sleepBlocks',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startHTTP',
			call: 'admin_startHTTP',
			params: 5,
			inputFormatter: [null, null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'stopHTTP',
			call: 'admin_stopHTTP'
		}),
		// This method is deprecated.
		new web3._extend.Method({
			name: 'startRPC',
			call: 'admin_startRPC',
			params: 5,
			inputFormatter: [null, null, null, null, null]
		}),
		// This method is deprecated.
		new web3._extend.Method({
			name: 'stopRPC',
			call: 'admin_stopRPC'
		}),
		new web3._extend.Method({
			name: 'startWS',
			call: 'admin_startWS',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'stopWS',
			call: 'admin_stopWS'
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'nodeInfo',
			getter: 'admin_nodeInfo'
		}),
		new web3._extend.Property({
			name: 'peers',
			getter: 'admin_peers'
		}),
		new web3._extend.Property({
			name: 'datadir',
			getter: 'admin_datadir'
		}),
	]
});
`

const DebugJs = `
web3._extend({
	property: 'debug',
	methods: [
		new web3._extend.Method({
			name: 'accountRange',
			call: 'debug_accountRange',
			params: 6,
			inputFormatter: [web3._extend.formatters.inputDefaultBlockNumberFormatter, null, null, null, null, null],
		}),
		new web3._extend.Method({
			name: 'printBlock',
			call: 'debug_printBlock',
			params: 1,
			outputFormatter: console.log
		}),
		new web3._extend.Method({
			name: 'getBlockRlp',
			call: 'debug_getBlockRlp',
			params: 1
		}),
		new web3._extend.Method({
			name: 'testSignCliqueBlock',
			call: 'debug_testSignCliqueBlock',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null],
		}),
		new web3._extend.Method({
			name: 'setHead',
			call: 'debug_setHead',
			params: 1
		}),
		new web3._extend.Method({
			name: 'seedHash',
			call: 'debug_seedHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'dumpBlock',
			call: 'debug_dumpBlock',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'chaindbProperty',
			call: 'debug_chaindbProperty',
			params: 1,
			outputFormatter: console.log
		}),
		new web3._extend.Method({
			name: 'chaindbCompact',
			call: 'debug_chaindbCompact',
		}),
		new web3._extend.Method({
			name: 'verbosity',
			call: 'debug_verbosity',
			params: 1
		}),
		new web3._extend.Method({
			name: 'vmodule',
			call: 'debug_vmodule',
			params: 1
		}),
		new web3._extend.Method({
			name: 'backtraceAt',
			call: 'debug_backtraceAt',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'stacks',
			call: 'debug_stacks',
			params: 0,
			outputFormatter: console.log
		}),
		new web3._extend.Method({
			name: 'freeOSMemory',
			call: 'debug_freeOSMemory',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'setGCPercent',
			call: 'debug_setGCPercent',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'memStats',
			call: 'debug_memStats',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'gcStats',
			call: 'debug_gcStats',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'cpuProfile',
			call: 'debug_cpuProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startCPUProfile',
			call: 'debug_startCPUProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'stopCPUProfile',
			call: 'debug_stopCPUProfile',
			params: 0
		}),
		new web3._extend.Method({
			name: 'goTrace',
			call: 'debug_goTrace',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startGoTrace',
			call: 'debug_startGoTrace',
			params: 1
		}),
		new web3._extend.Method({
			name: 'stopGoTrace',
			call: 'debug_stopGoTrace',
			params: 0
		}),
		new web3._extend.Method({
			name: 'blockProfile',
			call: 'debug_blockProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setBlockProfileRate',
			call: 'debug_setBlockProfileRate',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeBlockProfile',
			call: 'debug_writeBlockProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'mutexProfile',
			call: 'debug_mutexProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setMutexProfileFraction',
			call: 'debug_setMutexProfileFraction',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeMutexProfile',
			call: 'debug_writeMutexProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeMemProfile',
			call: 'debug_writeMemProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'traceBlock',
			call: 'debug_traceBlock',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockFromFile',
			call: 'debug_traceBlockFromFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBadBlock',
			call: 'debug_traceBadBlock',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'standardTraceBadBlockToFile',
			call: 'debug_standardTraceBadBlockToFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'standardTraceBlockToFile',
			call: 'debug_standardTraceBlockToFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockByNumber',
			call: 'debug_traceBlockByNumber',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockByHash',
			call: 'debug_traceBlockByHash',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceTransaction',
			call: 'debug_traceTransaction',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceCall',
			call: 'debug_traceCall',
			params: 3,
			inputFormatter: [null, null, null]
		}),
		new web3._extend.Method({
			name: 'preimage',
			call: 'debug_preimage',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getBadBlocks',
			call: 'debug_getBadBlocks',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'storageRangeAt',
			call: 'debug_storageRangeAt',
			params: 5,
		}),
		new web3._extend.Method({
			name: 'getModifiedAccountsByNumber',
			call: 'debug_getModifiedAccountsByNumber',
			params: 2,
			inputFormatter: [null, null],
		}),
		new web3._extend.Method({
			name: 'getModifiedAccountsByHash',
			call: 'debug_getModifiedAccountsByHash',
			params: 2,
			inputFormatter:[null, null],
		}),
		new web3._extend.Method({
			name: 'freezeClient',
			call: 'debug_freezeClient',
			params: 1,
		}),
	],
	properties: []
});
`

const EthJs = `
web3._extend({
	property: 'eth',
	methods: [
		new web3._extend.Method({
			name: 'chainId',
			call: 'eth_chainId',
			params: 0
		}),
		new web3._extend.Method({
			name: 'sign',
			call: 'eth_sign',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Method({
			name: 'resend',
			call: 'eth_resend',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, web3._extend.utils.fromDecimal, web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Method({
			name: 'signTransaction',
			call: 'eth_signTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'estimateGas',
			call: 'eth_estimateGas',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputCallFormatter, web3._extend.formatters.inputBlockNumberFormatter],
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Method({
			name: 'submitTransaction',
			call: 'eth_submitTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'fillTransaction',
			call: 'eth_fillTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'getHeaderByNumber',
			call: 'eth_getHeaderByNumber',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getHeaderByHash',
			call: 'eth_getHeaderByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getBlockByNumber',
			call: 'eth_getBlockByNumber',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, function (val) { return !!val; }]
		}),
		new web3._extend.Method({
			name: 'getBlockByHash',
			call: 'eth_getBlockByHash',
			params: 2,
			inputFormatter: [null, function (val) { return !!val; }]
		}),
		new web3._extend.Method({
			name: 'getRawTransaction',
			call: 'eth_getRawTransactionByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getRawTransactionFromBlock',
			call: function(args) {
				return (web3._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'eth_getRawTransactionByBlockHashAndIndex' : 'eth_getRawTransactionByBlockNumberAndIndex';
			},
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, web3._extend.utils.toHex]
		}),
		new web3._extend.Method({
			name: 'getProof',
			call: 'eth_getProof',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null, web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'createAccessList',
			call: 'eth_createAccessList',
			params: 2,
			inputFormatter: [null, web3._extend.formatters.inputBlockNumberFormatter],
		}),
		new web3._extend.Method({
			name: 'feeHistory',
			call: 'eth_feeHistory',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputBlockNumberFormatter, null]
		}),
		new web3._extend.Method({
			name: 'getSysTransactionsByBlockNumber',
			call: 'eth_getSysTransactionsByBlockNumber',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getSysTransactionsByBlockHash',
			call: 'eth_getSysTransactionsByBlockHash',
			params: 1
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'pendingTransactions',
			getter: 'eth_pendingTransactions',
			outputFormatter: function(txs) {
				var formatted = [];
				for (var i = 0; i < txs.length; i++) {
					formatted.push(web3._extend.formatters.outputTransactionFormatter(txs[i]));
					formatted[i].blockHash = null;
				}
				return formatted;
			}
		}),
		new web3._extend.Property({
			name: 'maxPriorityFeePerGas',
			getter: 'eth_maxPriorityFeePerGas',
			outputFormatter: web3._extend.utils.toBigNumber
		}),
		new web3._extend.Property({
			name: 'gasPricePrediction',
			getter: 'eth_gasPricePrediction'
		}),
	]
});
`

const MinerJs = `
web3._extend({
	property: 'miner',
	methods: [
		new web3._extend.Method({
			name: 'start',
			call: 'miner_start',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'stop',
			call: 'miner_stop'
		}),
		new web3._extend.Method({
			name: 'setEtherbase',
			call: 'miner_setEtherbase',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter]
		}),
		new web3._extend.Method({
			name: 'setExtra',
			call: 'miner_setExtra',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setGasPrice',
			call: 'miner_setGasPrice',
			params: 1,
			inputFormatter: [web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Method({
			name: 'setGasLimit',
			call: 'miner_setGasLimit',
			params: 1,
			inputFormatter: [web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Method({
			name: 'setRecommitInterval',
			call: 'miner_setRecommitInterval',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'getHashrate',
			call: 'miner_getHashrate'
		}),
	],
	properties: []
});
`

const NetJs = `
web3._extend({
	property: 'net',
	methods: [],
	properties: [
		new web3._extend.Property({
			name: 'version',
			getter: 'net_version'
		}),
	]
});
`

const PersonalJs = `
web3._extend({
	property: 'personal',
	methods: [
		new web3._extend.Method({
			name: 'importRawKey',
			call: 'personal_importRawKey',
			params: 2
		}),
		new web3._extend.Method({
			name: 'sign',
			call: 'personal_sign',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Method({
			name: 'ecRecover',
			call: 'personal_ecRecover',
			params: 2
		}),
		new web3._extend.Method({
			name: 'openWallet',
			call: 'personal_openWallet',
			params: 2
		}),
		new web3._extend.Method({
			name: 'deriveAccount',
			call: 'personal_deriveAccount',
			params: 3
		}),
		new web3._extend.Method({
			name: 'signTransaction',
			call: 'personal_signTransaction',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, null]
		}),
		new web3._extend.Method({
			name: 'unpair',
			call: 'personal_unpair',
			params: 2
		}),
		new web3._extend.Method({
			name: 'initializeWallet',
			call: 'personal_initializeWallet',
			params: 1
		})
	],
	properties: [
		new web3._extend.Property({
			name: 'listWallets',
			getter: 'personal_listWallets'
		}),
	]
})
`

const RpcJs = `
web3._extend({
	property: 'rpc',
	methods: [],
	properties: [
		new web3._extend.Property({
			name: 'modules',
			getter: 'rpc_modules'
		}),
	]
});
`

const TxpoolJs = `
web3._extend({
	property: 'txpool',
	methods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'content',
			getter: 'txpool_content'
		}),
		new web3._extend.Property({
			name: 'inspect',
			getter: 'txpool_inspect'
		}),
		new web3._extend.Property({
			name: 'status',
			getter: 'txpool_status',
			outputFormatter: function(status) {
				status.pending = web3._extend.utils.toDecimal(status.pending);
				status.queued = web3._extend.utils.toDecimal(status.queued);
				return status;
			}
		}),
		new web3._extend.Method({
			name: 'contentFrom',
			call: 'txpool_contentFrom',
			params: 1,
		}),
		new web3._extend.Property({
			name: 'jamIndex',
			getter: 'txpool_jamIndex'
		}),
	]
});
`

const LESJs = `
web3._extend({
	property: 'les',
	methods:
	[
		new web3._extend.Method({
			name: 'getCheckpoint',
			call: 'les_getCheckpoint',
			params: 1
		}),
		new web3._extend.Method({
			name: 'clientInfo',
			call: 'les_clientInfo',
			params: 1
		}),
		new web3._extend.Method({
			name: 'priorityClientInfo',
			call: 'les_priorityClientInfo',
			params: 3
		}),
		new web3._extend.Method({
			name: 'setClientParams',
			call: 'les_setClientParams',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setDefaultParams',
			call: 'les_setDefaultParams',
			params: 1
		}),
		new web3._extend.Method({
			name: 'addBalance',
			call: 'les_addBalance',
			params: 2
		}),
	],
	properties:
	[
		new web3._extend.Property({
			name: 'latestCheckpoint',
			getter: 'les_latestCheckpoint'
		}),
		new web3._extend.Property({
			name: 'checkpointContractAddress',
			getter: 'les_getCheckpointContractAddress'
		}),
		new web3._extend.Property({
			name: 'serverInfo',
			getter: 'les_serverInfo'
		}),
	]
});
`

const VfluxJs = `
web3._extend({
	property: 'vflux',
	methods:
	[
		new web3._extend.Method({
			name: 'distribution',
			call: 'vflux_distribution',
			params: 2
		}),
		new web3._extend.Method({
			name: 'timeout',
			call: 'vflux_timeout',
			params: 2
		}),
		new web3._extend.Method({
			name: 'value',
			call: 'vflux_value',
			params: 2
		}),
	],
	properties:
	[
		new web3._extend.Property({
			name: 'requestStats',
			getter: 'vflux_requestStats'
		}),
	]
});
`
