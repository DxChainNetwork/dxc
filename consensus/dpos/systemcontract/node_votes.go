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

type VotesRewardRedeemInfo struct {
	Validator           common.Address
	ValidatorName       string
	ValidatorRate       uint8
	ValidatorTotalVotes *big.Int
	Amount              *big.Int
	PendingReward       *big.Int
	PendingRedeem       *big.Int
	LockRedeemEpochs    []*big.Int
	LockRedeemVotes     []*big.Int
}

// NewNodeVotes return Proposals contract instance
func NewNodeVotes() *NodeVotes {
	return &NodeVotes{
		abi:          abiMap[NodeVotesContractName],
		contractAddr: NodeVotesContractAddr,
	}
}

// PendingVoteReward function pendingReward
func (n *NodeVotes) PendingVoteReward(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, val common.Address, voter common.Address) (*big.Int, error) {
	method := "pendingReward"
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

// PendingVoteRedeem function PendingRedeem
func (n *NodeVotes) PendingVoteRedeem(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, val common.Address, voter common.Address) (*big.Int, error) {
	method := "pendingRedeem"
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

// VotesRewardRedeemInfo function votesRewardRedeemInfo
func (n *NodeVotes) VotesRewardRedeemInfo(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, val common.Address, voter common.Address) (*VotesRewardRedeemInfo, error) {
	method := "votesRewardRedeemInfo"
	data, err := n.abi.Pack(method, val, voter)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return &VotesRewardRedeemInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return &VotesRewardRedeemInfo{}, err
	}
	rewardRedeemInfo := &VotesRewardRedeemInfo{}
	err = n.abi.UnpackIntoInterface(rewardRedeemInfo, method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return &VotesRewardRedeemInfo{}, err
	}
	return rewardRedeemInfo, nil
}

// VotesRewardRedeemInfoWithPage function votesRewardRedeemInfoWithPage
func (n *NodeVotes) VotesRewardRedeemInfoWithPage(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, voter common.Address, page *big.Int, size *big.Int) ([]VotesRewardRedeemInfo, error) {
	method := "votesRewardRedeemInfoWithPage"
	data, err := n.abi.Pack(method, voter, page, size)

	if err != nil {
		log.Error("NodeVotes Pack error", "method", method, "error", err)
		return []VotesRewardRedeemInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &n.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("NodeVotes contract execute error", "method", method, "error", err)
		return []VotesRewardRedeemInfo{}, err
	}
	var rewardRedeemInfos []VotesRewardRedeemInfo
	err = n.abi.UnpackIntoInterface(&rewardRedeemInfos, method, result)
	if err != nil {
		log.Error("NodeVotes contract Unpack error", "method", method, "error", err, "result", result)
		return []VotesRewardRedeemInfo{}, err
	}
	return rewardRedeemInfos, nil
}
