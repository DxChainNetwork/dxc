package systemcontract

import (
	"github.com/DxChainNetwork/dxc/accounts/abi"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/consensus/dpos/vmcaller"
	"github.com/DxChainNetwork/dxc/core"
	"github.com/DxChainNetwork/dxc/core/state"
	"github.com/DxChainNetwork/dxc/core/types"
	"github.com/DxChainNetwork/dxc/log"
	"github.com/DxChainNetwork/dxc/params"
	"math"
	"math/big"
)

type NodeVotes struct {
	abi          abi.ABI
	contractAddr common.Address
}

type VoteInfo struct {
	Validator         common.Address
	Votes             *big.Int
	RewardDebt        *big.Int
	UpdateRewardEpoch *big.Int
}

type RedeemVoterInfo struct {
	Validator   common.Address
	Votes       *big.Int
	RedeemBlock *big.Int
}

// NewNodeVotes return Proposals contract instance
func NewNodeVotes() *NodeVotes {
	return &NodeVotes{
		abi:          abiMap[NodeVotesContractName],
		contractAddr: NodeVotesContractAddr,
	}
}

// PendingVoteReward function PendingVoteReward
func (n *NodeVotes) PendingVoteReward(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, val common.Address, voter common.Address) (*big.Int, error) {
	method := "pendingVoteReward"
	data, err := n.abi.Pack(method, val, voter)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := n.abi.Unpack(method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	value, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("NodeVotes contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return value, nil
}

// PendingRedeem function PendingRedeem
func (n *NodeVotes) PendingRedeem(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, voter common.Address) (*big.Int, error) {
	method := "pendingRedeem"
	data, err := n.abi.Pack(method, voter)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := n.abi.Unpack(method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	value, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("NodeVotes contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return value, nil
}

// VoteListLength function VoteListLength
func (n *NodeVotes) VoteListLength(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, voter common.Address) (*big.Int, error) {
	method := "voteListLength"
	data, err := n.abi.Pack(method, voter)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := n.abi.Unpack(method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	count, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("NodeVotes contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return count, nil
}

// CancelVoteValidatorListLength function CancelVoteValidatorListLength
func (n *NodeVotes) CancelVoteValidatorListLength(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, voter common.Address) (*big.Int, error) {
	method := "cancelVoteValidatorListLength"
	data, err := n.abi.Pack(method, voter)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := n.abi.Unpack(method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	count, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("NodeVotes contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return count, nil
}

// CancelVoteValidatorList function CancelVoteValidatorList
func (n *NodeVotes) CancelVoteValidatorList(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page *big.Int, size *big.Int) ([]common.Address, error) {
	method := "cancelVoteValidatorList"
	data, err := n.abi.Pack(method, addr, page, size)
	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return []common.Address{}, err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return []common.Address{}, err
	}

	ret, err := n.abi.Unpack(method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return []common.Address{}, err
	}
	voters, ok := ret[0].([]common.Address)
	if !ok {
		log.Error("NodeVotes contract format result error", "method", method, "error", err)
		return []common.Address{}, err
	}

	return voters, nil
}

// VoteList function VoteList
func (n *NodeVotes) VoteList(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page *big.Int, size *big.Int) ([]VoteInfo, error) {
	method := "voteList"
	data, err := n.abi.Pack(method, addr, page, size)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return []VoteInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return []VoteInfo{}, err
	}
	var voteInfo []VoteInfo
	err = n.abi.UnpackIntoInterface(&voteInfo, method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return []VoteInfo{}, err
	}
	return voteInfo, nil
}

// RedeemInfo function RedeemInfo
func (n *NodeVotes) RedeemInfo(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page *big.Int, size *big.Int) ([]RedeemVoterInfo, error) {
	method := "redeemInfo"
	data, err := n.abi.Pack(method, addr, page, size)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return []RedeemVoterInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return []RedeemVoterInfo{}, err
	}
	var redeemVoterInfo []RedeemVoterInfo
	err = n.abi.UnpackIntoInterface(&redeemVoterInfo, method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return []RedeemVoterInfo{}, err
	}
	return redeemVoterInfo, nil
}
