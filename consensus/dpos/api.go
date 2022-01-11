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
	"github.com/DxChainNetwork/dxc/common/hexutil"

	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/consensus"
	"github.com/DxChainNetwork/dxc/consensus/dpos/systemcontract"
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
	InitBlock   *big.Int
	Guarantee   common.Address
	UpdateBlock *big.Int
	Status      uint8
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

// GetCurrentEpochValidators return current epoch validators
func (api *API) GetCurrentEpochValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	validators := systemcontract.NewValidators()
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return []common.Address{}, errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return []common.Address{}, err
	}
	curValidators, err := validators.GetCurrentEpochValidators(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []common.Address{}, err
	}
	return curValidators, nil
}

// GetEffictiveValidators return all effictive validators
func (api *API) GetEffictiveValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	validators := systemcontract.NewValidators()
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return []common.Address{}, errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return []common.Address{}, err
	}
	curValidators, err := validators.GetEffictiveValidators(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return []common.Address{}, err
	}
	return curValidators, nil
}

// GetMinDeposit return the minimum stake amount
func (api *API) GetMinDeposit(number *rpc.BlockNumber) (*big.Int, error) {
	base := systemcontract.NewBase()
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return big.NewInt(0), errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return big.NewInt(0), err
	}
	minDeposit, err := base.GetMinDeposit(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return big.NewInt(0), err
	}
	return minDeposit, nil

}

// GetAddressProposalSets return the address proposal id
func (api *API) GetAddressProposalSets(addr common.Address, page *big.Int, size *big.Int) ([]string, error) {
	proposals := systemcontract.NewProposals()
	var header *types.Header
	header = api.chain.CurrentHeader()
	if header == nil {
		return []string{}, errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return []string{}, err
	}
	proposalIds, err := proposals.AddressProposalSets(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, page, size)
	if err != nil {
		return []string{}, err
	}
	var newProposalIds []string
	for i := 0; i < len(proposalIds); i++ {
		id := hexutil.Encode(proposalIds[i][0:len(proposalIds[i])])
		newProposalIds = append(newProposalIds, id)
	}
	return newProposalIds, nil
}

// GetAllProposalSets return all proposals id
func (api *API) GetAllProposalSets(page *big.Int, size *big.Int) ([]string, error) {
	proposals := systemcontract.NewProposals()
	var header *types.Header
	header = api.chain.CurrentHeader()
	if header == nil {
		return []string{}, errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return []string{}, err
	}
	proposalIds, err := proposals.AllProposalSets(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, page, size)
	if err != nil {
		return []string{}, err
	}
	var newProposalIds []string
	for i := 0; i < len(proposalIds); i++ {
		id := hexutil.Encode(proposalIds[i][0:len(proposalIds[i])])
		newProposalIds = append(newProposalIds, id)
	}
	return newProposalIds, nil
}

// GetAllProposals return all proposals
func (api *API) GetAllProposals(page *big.Int, size *big.Int) ([]ProposalInfo, error) {
	proposals := systemcontract.NewProposals()
	var header *types.Header
	header = api.chain.CurrentHeader()
	if header == nil {
		return []ProposalInfo{}, errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return []ProposalInfo{}, err
	}
	proposalInfos, err := proposals.AllProposals(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, page, size)
	if err != nil {
		return []ProposalInfo{}, err
	}
	var newProposals []ProposalInfo
	for i := 0; i < len(proposalInfos); i++ {
		detail := ProposalInfo{
			Id:          hexutil.Encode(proposalInfos[i].Id[0:len(proposalInfos[i].Id)]),
			Proposer:    proposalInfos[i].Proposer,
			UpdateBlock: proposalInfos[i].UpdateBlock,
			PType:       proposalInfos[i].PType,
			Guarantee:   proposalInfos[i].Guarantee,
			Deposit:     proposalInfos[i].Deposit,
			Details:     proposalInfos[i].Details,
			InitBlock:   proposalInfos[i].InitBlock,
			Rate:        proposalInfos[i].Rate,
			Status:      proposalInfos[i].Status,
		}
		newProposals = append(newProposals, detail)
	}
	return newProposals, nil
}

// GetAddressProposals return the address proposals
func (api *API) GetAddressProposals(addr common.Address, page *big.Int, size *big.Int) ([]ProposalInfo, error) {
	proposals := systemcontract.NewProposals()
	var header *types.Header
	header = api.chain.CurrentHeader()
	if header == nil {
		return []ProposalInfo{}, errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return []ProposalInfo{}, err
	}
	proposalInfos, err := proposals.AddressProposals(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr, page, size)
	if err != nil {
		return []ProposalInfo{}, err
	}
	var newProposals []ProposalInfo
	for i := 0; i < len(proposalInfos); i++ {
		detail := ProposalInfo{
			Id:          hexutil.Encode(proposalInfos[i].Id[0:len(proposalInfos[i].Id)]),
			Proposer:    proposalInfos[i].Proposer,
			UpdateBlock: proposalInfos[i].UpdateBlock,
			PType:       proposalInfos[i].PType,
			Guarantee:   proposalInfos[i].Guarantee,
			Deposit:     proposalInfos[i].Deposit,
			Details:     proposalInfos[i].Details,
			InitBlock:   proposalInfos[i].InitBlock,
			Rate:        proposalInfos[i].Rate,
			Status:      proposalInfos[i].Status,
		}
		newProposals = append(newProposals, detail)
	}
	return newProposals, nil
}

// GetProposalCount return all proposal count
func (api *API) GetProposalCount() (*big.Int, error) {
	proposals := systemcontract.NewProposals()
	var header *types.Header
	header = api.chain.CurrentHeader()
	if header == nil {
		return big.NewInt(0), errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := proposals.ProposalCount(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
}

// GetAddressProposalCount return the address proposal count
func (api *API) GetAddressProposalCount(addr common.Address) (*big.Int, error) {
	proposals := systemcontract.NewProposals()
	var header *types.Header
	header = api.chain.CurrentHeader()
	if header == nil {
		return big.NewInt(0), errUnknownBlock
	}
	state, err := api.dpos.stateFn(header.Root)
	if err != nil {
		return big.NewInt(0), err
	}
	count, err := proposals.AddressProposalCount(state, header, newChainContext(api.chain, api.dpos), api.dpos.chainConfig, addr)
	if err != nil {
		return big.NewInt(0), err
	}
	return count, nil
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
