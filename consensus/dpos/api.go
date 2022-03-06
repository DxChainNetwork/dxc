// Copyright 2017 The go-ethereum Authors
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

package dpos

import (
	"fmt"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/common/hexutil"
	"github.com/DxChainNetwork/dxc/consensus"
	"github.com/DxChainNetwork/dxc/consensus/dpos/systemcontract"
	"github.com/DxChainNetwork/dxc/core/state"
	"github.com/DxChainNetwork/dxc/core/types"
	"github.com/DxChainNetwork/dxc/rpc"
	"math/big"
)

// API is a user facing RPC API to allow controlling the validator and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain consensus.ChainHeaderReader
	dpos  *Dpos
}

type ProposalInfo struct {
	Id          string
	Proposer    common.Address
	PType       uint8
	Deposit     *big.Int
	Rate        uint8
	Details     string
	Name        string
	InitBlock   *big.Int
	Guarantee   common.Address
	UpdateBlock *big.Int
	Status      uint8
}

func (api *API) GetHeaderAndState(number *rpc.BlockNumber) (*types.Header, *state.StateDB, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, nil, errUnknownBlock
	}
	statedb, err := api.dpos.stateFn(header.Root)
	return header, statedb, err
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetValidators retrieves the list of authorized validators at the specified block.
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the validators from its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.validators(), nil
}

// GetValidatorsAtHash retrieves the list of authorized validators at the specified block.
func (api *API) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.validators(), nil
}

func (api *API) GetBaseInfos(number *rpc.BlockNumber) (map[string]interface{}, error) {
	base := systemcontract.NewBase()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return map[string]interface{}{}, err
	}
	infos, err := base.GetBaseInfos(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return infos, nil
}

// GetValidator return the validator of address
func (api *API) GetValidator(addr common.Address, number *rpc.BlockNumber) (*systemcontract.Validator, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return &systemcontract.Validator{}, err
	}
	val, err := validators.GetValidator(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return &systemcontract.Validator{}, err
	}
	return val, nil
}

// GetTotalDeposit return total deposit
func (api *API) GetTotalDeposit(number *rpc.BlockNumber) (*big.Int, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	deposit, err := validators.TotalDeposit(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return big.NewInt(0), err
	}
	return deposit, nil
}

// GetCurrentEpochValidators return current epoch validators
func (api *API) GetCurrentEpochValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []common.Address{}, err
	}
	curValidators, err := validators.GetCurrentEpochValidators(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []common.Address{}, err
	}
	return curValidators, nil
}

// GetEffictiveValidators return all effictive validators
func (api *API) GetEffictiveValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []common.Address{}, err
	}

	size := big.NewInt(50)
	var allValidators []common.Address

	count, err := validators.EffictiveValsLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []common.Address{}, err
	}

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		allValidators, err = validators.GetEffictiveValidatorsWithPage(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, page, size)
		if err != nil {
			return []common.Address{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			voters, err := validators.GetEffictiveValidatorsWithPage(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, big.NewInt(i), size)
			if err != nil {
				return []common.Address{}, err
			}
			allValidators = append(allValidators, voters...)
		}
	}

	return allValidators, nil
}

// GetInvalidValidators return all invalid validators
func (api *API) GetInvalidValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []common.Address{}, err
	}

	size := big.NewInt(50)
	var invalidValidators []common.Address

	count, err := validators.InvalidValsLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []common.Address{}, err
	}

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		invalidValidators, err = validators.GetInvalidValidatorsWithPage(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, page, size)
		if err != nil {
			return []common.Address{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			voters, err := validators.GetInvalidValidatorsWithPage(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, big.NewInt(i), size)
			if err != nil {
				return []common.Address{}, err
			}
			invalidValidators = append(invalidValidators, voters...)
		}
	}

	return invalidValidators, nil
}

// GetCancelQueueValidators return all canceling queue validators
func (api *API) GetCancelQueueValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []common.Address{}, err
	}
	cancelingValidators, err := validators.GetCancelQueueValidators(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []common.Address{}, err
	}
	return cancelingValidators, nil
}

// GetValidatorVoters return the address voter
func (api *API) GetValidatorVoters(addr common.Address, number *rpc.BlockNumber) ([]common.Address, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []common.Address{}, err
	}

	size := big.NewInt(50)
	var allVoters []common.Address

	count, err := validators.ValidatorVotersLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return []common.Address{}, err
	}

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		allVoters, err = validators.GetValidatorVoters(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, page, size)
		if err != nil {
			return []common.Address{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			voters, err := validators.GetValidatorVoters(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, big.NewInt(i), size)
			if err != nil {
				return []common.Address{}, err
			}
			allVoters = append(allVoters, voters...)
		}
	}

	return allVoters, nil
}

// EffictiveValsLength return effictive validators length
func (api *API) EffictiveValsLength(number *rpc.BlockNumber) (*big.Int, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := validators.EffictiveValsLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// InvalidValsLength return invalid validators length
func (api *API) InvalidValsLength(number *rpc.BlockNumber) (*big.Int, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := validators.InvalidValsLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// CancelQueueValidatorsLength return cancel queue validators length
func (api *API) CancelQueueValidatorsLength(number *rpc.BlockNumber) (*big.Int, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := validators.CancelQueueValidatorsLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// ValidatorVotersLength return the validator voters length
func (api *API) ValidatorVotersLength(addr common.Address, number *rpc.BlockNumber) (*big.Int, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := validators.ValidatorVotersLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// IsEffictiveValidator return the address is validator
func (api *API) IsEffictiveValidator(addr common.Address, number *rpc.BlockNumber) (bool, error) {
	validators := systemcontract.NewValidators()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return false, err
	}
	val, err := validators.IsEffictiveValidator(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return false, err
	}
	return val, nil
}

// Proposals

// GetAddressProposalSets return the address proposal id
func (api *API) GetAddressProposalSets(addr common.Address, number *rpc.BlockNumber) ([]string, error) {
	proposals := systemcontract.NewProposals()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []string{}, err
	}
	count, err := proposals.AddressProposalCount(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return []string{}, err
	}

	size := big.NewInt(50)
	var allSets [][4]byte

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		allSets, err = proposals.AddressProposalSets(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, page, size)
		if err != nil {
			return []string{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			proposalIds, err := proposals.AddressProposalSets(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, big.NewInt(i), size)
			if err != nil {
				return []string{}, err
			}
			allSets = append(allSets, proposalIds...)
		}
	}

	newProposalIds := []string{}
	for i := 0; i < len(allSets); i++ {
		id := hexutil.Encode(allSets[i][0:len(allSets[i])])
		newProposalIds = append(newProposalIds, id)
	}
	return newProposalIds, nil
}

// GetAllProposalSets return all proposals id
func (api *API) GetAllProposalSets(number *rpc.BlockNumber) ([]string, error) {
	proposals := systemcontract.NewProposals()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []string{}, err
	}

	count, err := proposals.ProposalCount(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []string{}, err
	}

	size := big.NewInt(50)
	var allSets [][4]byte

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		allSets, err = proposals.AllProposalSets(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, page, size)
		if err != nil {
			return []string{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			proposalIds, err := proposals.AllProposalSets(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, big.NewInt(i), size)
			if err != nil {
				return []string{}, err
			}
			allSets = append(allSets, proposalIds...)
		}
	}

	newProposalIds := []string{}
	for i := 0; i < len(allSets); i++ {
		id := hexutil.Encode(allSets[i][0:len(allSets[i])])
		newProposalIds = append(newProposalIds, id)
	}
	return newProposalIds, nil

}

// GetAllProposals return all proposals
func (api *API) GetAllProposals(number *rpc.BlockNumber) ([]ProposalInfo, error) {
	proposals := systemcontract.NewProposals()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []ProposalInfo{}, err
	}

	count, err := proposals.ProposalCount(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []ProposalInfo{}, err
	}

	size := big.NewInt(50)
	var allProposals []systemcontract.ProposalInfo

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		allProposals, err = proposals.AllProposals(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, page, size)
		if err != nil {
			return []ProposalInfo{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			proposals, err := proposals.AllProposals(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, big.NewInt(i), size)
			if err != nil {
				return []ProposalInfo{}, err
			}
			allProposals = append(allProposals, proposals...)
		}
	}

	newProposals := []ProposalInfo{}
	for i := 0; i < len(allProposals); i++ {
		detail := ProposalInfo{
			Id:          hexutil.Encode(allProposals[i].Id[0:len(allProposals[i].Id)]),
			Proposer:    allProposals[i].Proposer,
			UpdateBlock: allProposals[i].UpdateBlock,
			PType:       allProposals[i].PType,
			Guarantee:   allProposals[i].Guarantee,
			Deposit:     allProposals[i].Deposit,
			Details:     allProposals[i].Details,
			Name:        allProposals[i].Name,
			InitBlock:   allProposals[i].InitBlock,
			Rate:        allProposals[i].Rate,
			Status:      allProposals[i].Status,
		}
		newProposals = append(newProposals, detail)
	}
	return newProposals, nil

}

// GetProposal return the proposal of id
func (api *API) GetProposal(id string, number *rpc.BlockNumber) (*ProposalInfo, error) {
	proposals := systemcontract.NewProposals()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return &ProposalInfo{}, err
	}
	proposalInfo, err := proposals.GetProposal(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, id)
	if err != nil {
		return &ProposalInfo{}, err
	}
	detail := &ProposalInfo{
		Id:          hexutil.Encode(proposalInfo.Id[0:len(proposalInfo.Id)]),
		Proposer:    proposalInfo.Proposer,
		UpdateBlock: proposalInfo.UpdateBlock,
		PType:       proposalInfo.PType,
		Guarantee:   proposalInfo.Guarantee,
		Deposit:     proposalInfo.Deposit,
		Details:     proposalInfo.Details,
		Name:        proposalInfo.Name,
		InitBlock:   proposalInfo.InitBlock,
		Rate:        proposalInfo.Rate,
		Status:      proposalInfo.Status,
	}
	return detail, nil
}

// GetAddressProposals return the address proposals
func (api *API) GetAddressProposals(addr common.Address, number *rpc.BlockNumber) ([]ProposalInfo, error) {
	proposals := systemcontract.NewProposals()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []ProposalInfo{}, err
	}

	var allProposals []systemcontract.ProposalInfo
	size := big.NewInt(50)

	count, err := proposals.AddressProposalCount(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return []ProposalInfo{}, err
	}

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		allProposals, err = proposals.AddressProposals(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, page, size)
		if err != nil {
			return []ProposalInfo{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			proposalIds, err := proposals.AddressProposals(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, big.NewInt(i), size)
			if err != nil {
				return []ProposalInfo{}, err
			}
			allProposals = append(allProposals, proposalIds...)
		}
	}

	var newProposals []ProposalInfo
	for i := 0; i < len(allProposals); i++ {
		detail := ProposalInfo{
			Id:          hexutil.Encode(allProposals[i].Id[0:len(allProposals[i].Id)]),
			Proposer:    allProposals[i].Proposer,
			UpdateBlock: allProposals[i].UpdateBlock,
			PType:       allProposals[i].PType,
			Guarantee:   allProposals[i].Guarantee,
			Deposit:     allProposals[i].Deposit,
			Details:     allProposals[i].Details,
			Name:        allProposals[i].Name,
			InitBlock:   allProposals[i].InitBlock,
			Rate:        allProposals[i].Rate,
			Status:      allProposals[i].Status,
		}
		newProposals = append(newProposals, detail)
	}
	return newProposals, nil
}

// GetProposalCount return all proposal count
func (api *API) GetProposalCount(number *rpc.BlockNumber) (*big.Int, error) {
	proposals := systemcontract.NewProposals()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := proposals.ProposalCount(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// GetAddressProposalCount return the address proposal count
func (api *API) GetAddressProposalCount(addr common.Address, number *rpc.BlockNumber) (*big.Int, error) {
	proposals := systemcontract.NewProposals()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := proposals.AddressProposalCount(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// NodeVotes

// PendingVoteReward return the voter vote the validator rewards
func (api *API) PendingVoteReward(val common.Address, voter common.Address, number *rpc.BlockNumber) (*big.Int, error) {
	nodeVotes := systemcontract.NewNodeVotes()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	value, err := nodeVotes.PendingVoteReward(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, val, voter)
	if err != nil {
		return big.NewInt(0), err
	}
	return value, nil
}

// PendingVoteRedeem return the voter redeem validators voters
func (api *API) PendingVoteRedeem(val common.Address, voter common.Address, number *rpc.BlockNumber) (*big.Int, error) {
	nodeVotes := systemcontract.NewNodeVotes()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	value, err := nodeVotes.PendingVoteRedeem(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, val, voter)
	if err != nil {
		return big.NewInt(0), err
	}
	return value, nil
}

// VoteListLength return the voter vote list length
func (api *API) VoteListLength(addr common.Address, number *rpc.BlockNumber) (*big.Int, error) {
	nodeVotes := systemcontract.NewNodeVotes()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := nodeVotes.VoteListLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// VotesRewardRedeemInfo votesRewardRedeemInfo
func (api *API) VotesRewardRedeemInfo(val common.Address, voter common.Address, number *rpc.BlockNumber) (*systemcontract.VotesRewardRedeemInfo, error) {
	nodeVotes := systemcontract.NewNodeVotes()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return &systemcontract.VotesRewardRedeemInfo{}, err
	}

	info, err := nodeVotes.VotesRewardRedeemInfo(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, val, voter)
	if err != nil {
		return &systemcontract.VotesRewardRedeemInfo{}, err
	}
	return info, nil
}

// VotesRewardRedeemInfos nodevotes.VotesRewardRedeemInfos
func (api *API) VotesRewardRedeemInfos(voter common.Address, number *rpc.BlockNumber) ([]systemcontract.VotesRewardRedeemInfo, error) {
	nodeVotes := systemcontract.NewNodeVotes()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []systemcontract.VotesRewardRedeemInfo{}, err
	}

	size := big.NewInt(50)
	var allInfos []systemcontract.VotesRewardRedeemInfo

	count, err := nodeVotes.VoteListLength(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, voter)
	if err != nil {
		return []systemcontract.VotesRewardRedeemInfo{}, err
	}

	if count.Cmp(size) <= 0 {
		page := big.NewInt(1)
		allInfos, err = nodeVotes.VotesRewardRedeemInfoWithPage(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, voter, page, size)
		if err != nil {
			return []systemcontract.VotesRewardRedeemInfo{}, err
		}
	} else {
		var res big.Int
		div := res.Div(count, size)
		div = res.Add(div, big.NewInt(1))
		for i := int64(1); i <= div.Int64(); i++ {
			infos, err := nodeVotes.VotesRewardRedeemInfoWithPage(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, voter, big.NewInt(i), size)
			if err != nil {
				return []systemcontract.VotesRewardRedeemInfo{}, err
			}
			allInfos = append(allInfos, infos...)
		}
		return allInfos, nil
	}

	return allInfos, nil
}

// systemRewards

type SysRewardsInfo struct {
	Epochs            []*big.Int
	ValidatorRewards  []*big.Int
	DelegatorsRewards []*big.Int
	Rates             []uint
	PendingReward     *big.Int
	FrozenReward      *big.Int
	RewardPerVote     *big.Int
}

// EpochInfo return the epoch info
func (api *API) EpochInfo(epoch *big.Int, number *rpc.BlockNumber) (*systemcontract.EpochInfo, error) {
	systemRewards := systemcontract.NewSystemRewards()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return &systemcontract.EpochInfo{}, err
	}
	epochInfo, err := systemRewards.GetEpochInfo(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, epoch)
	if err != nil {
		return &systemcontract.EpochInfo{}, err
	}
	return epochInfo, nil
}

// KickoutInfo return kickout addresses in epoch
func (api *API) KickoutInfo(epoch *big.Int, number *rpc.BlockNumber) ([]common.Address, error) {
	systemRewards := systemcontract.NewSystemRewards()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return []common.Address{}, err
	}
	vals, err := systemRewards.KickoutInfo(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, epoch)
	if err != nil {
		return []common.Address{}, err
	}
	return vals, nil
}

// ValidatorRewardsInfo return the sys reward info
func (api *API) ValidatorRewardsInfo(addr common.Address, number *rpc.BlockNumber) (*SysRewardsInfo, error) {
	systemRewards := systemcontract.NewSystemRewards()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return &SysRewardsInfo{}, err
	}
	rewardInfo, err := systemRewards.ValidatorRewardsInfo(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return &SysRewardsInfo{}, err
	}
	newInfo := SysRewardsInfo{
		Epochs:            rewardInfo.Epochs,
		ValidatorRewards:  rewardInfo.ValidatorRewards,
		DelegatorsRewards: rewardInfo.DelegatorsRewards,
		PendingReward:     rewardInfo.PendingReward,
		FrozenReward:      rewardInfo.FrozenReward,
		RewardPerVote:     rewardInfo.RewardPerVote,
	}

	for i := 0; i < len(rewardInfo.Rates); i++ {
		newInfo.Rates = append(newInfo.Rates, uint(rewardInfo.Rates[i]))
	}

	return &newInfo, nil
}

// ValidatorRewardInfoByEpoch return the address and the epoch reward info
func (api *API) ValidatorRewardInfoByEpoch(addr common.Address, epoch *big.Int, number *rpc.BlockNumber) (*systemcontract.Reward, error) {
	systemRewards := systemcontract.NewSystemRewards()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return &systemcontract.Reward{}, err
	}
	rewards, err := systemRewards.GetValRewardInfoByEpoch(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, epoch)
	if err != nil {
		return &systemcontract.Reward{}, err
	}
	return rewards, nil
}

// PendingValidatorReward return the address reward
func (api *API) PendingValidatorReward(addr common.Address, number *rpc.BlockNumber) (map[string]*big.Int, error) {
	systemRewards := systemcontract.NewSystemRewards()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return map[string]*big.Int{}, err
	}
	avaliable, frozen, err := systemRewards.PendingValidatorReward(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return map[string]*big.Int{}, err
	}
	result := make(map[string]*big.Int)
	result["avaliable"] = avaliable
	result["frozen"] = frozen
	return result, nil
}

// PunishInfo punishInfo function of systemRewards contract
func (api *API) PunishInfo(addr common.Address, epoch *big.Int, number *rpc.BlockNumber) (*big.Int, error) {
	systemRewards := systemcontract.NewSystemRewards()
	header, statedb, err := api.GetHeaderAndState(number)
	if err != nil {
		return big.NewInt(0), err
	}
	punishCount, err := systemRewards.PunishInfo(statedb, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, epoch)
	if err != nil {
		return big.NewInt(0), err
	}

	return punishCount, nil
}

type status struct {
	InturnPercent float64                `json:"inturnPercent"`
	SigningStatus map[common.Address]int `json:"sealerActivity"`
	NumBlocks     uint64                 `json:"numBlocks"`
}

// Status returns the status of the last N blocks,
// - the number of active validators,
// - the number of validators,
// - the percentage of in-turn blocks
func (api *API) Status() (*status, error) {
	var (
		numBlocks = uint64(64)
		header    = api.chain.CurrentHeader()
		diff      = uint64(0)
		optimals  = 0
	)
	snap, err := api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	var (
		validators = snap.validators()
		end        = header.Number.Uint64()
		start      = end - numBlocks
	)
	if numBlocks > end {
		start = 1
		numBlocks = end - start
	}
	signStatus := make(map[common.Address]int)
	for _, s := range validators {
		signStatus[s] = 0
	}
	for n := start; n < end; n++ {
		h := api.chain.GetHeaderByNumber(n)
		if h == nil {
			return nil, fmt.Errorf("missing block %d", n)
		}
		if h.Difficulty.Cmp(diffInTurn) == 0 {
			optimals++
		}
		diff += h.Difficulty.Uint64()
		sealer, err := api.dpos.Author(h)
		if err != nil {
			return nil, err
		}
		signStatus[sealer]++
	}
	return &status{
		InturnPercent: float64(100*optimals) / float64(numBlocks),
		SigningStatus: signStatus,
		NumBlocks:     numBlocks,
	}, nil
}
