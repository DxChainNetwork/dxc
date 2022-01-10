package systemcontract

import (
	"errors"
	"github.com/DxChainNetwork/dxc/accounts/abi"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/common/hexutil"
	"github.com/DxChainNetwork/dxc/consensus/dpos/vmcaller"
	"github.com/DxChainNetwork/dxc/core"
	"github.com/DxChainNetwork/dxc/core/state"
	"github.com/DxChainNetwork/dxc/core/types"
	"github.com/DxChainNetwork/dxc/log"
	"github.com/DxChainNetwork/dxc/params"
	"math"
	"math/big"
)

type Proposals struct {
	abi          abi.ABI
	contractAddr common.Address
}

type ProposalInfo struct {
	Id          [4]byte
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

// NewProposals return Proposals contract instance
func NewProposals() *Proposals {
	return &Proposals{
		abi:          abiMap[ProposalsContractName],
		contractAddr: ProposalsContractAddr,
	}
}

// InitProposal function initProposal
func (p *Proposals) InitProposal(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, value *big.Int, pType uint8, rate uint8, details string) error {
	method := "initProposal"
	if data, err := p.abi.Pack(method, pType, rate, details); err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return err
	} else {
		nonce := statedb.GetNonce(header.Coinbase)
		log.Info("InitProposal getBalance", "addr", header.Coinbase, "balance", statedb.GetBalance(header.Coinbase))
		msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, nonce, value, math.MaxUint64, new(big.Int), data, true)
		if _, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config); err != nil {
			log.Error("InitProposal execute", "error", err)
			return err
		}
	}

	return nil
}

func (p *Proposals) AddressProposalSets(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page int64, size int64) ([]string, error) {
	method := "addressProposalSets"
	data, err := p.abi.Pack(method, addr, big.NewInt(page), big.NewInt(size))

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return []string{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AddressProposalSets result", "error", err)
		return []string{}, err
	}

	ret, err := p.abi.Unpack(method, result)

	if err != nil {
		return []string{}, err
	}
	if proposalsIds, ok := ret[0].([][4]byte); !ok {
		return []string{}, errors.New("invalid AddressProposalSets result format")
	} else {
		log.Info("AddressProposalSets result", "address", addr.String(), "result", proposalsIds)
		var newProposalsIds []string
		for i := 0; i < len(proposalsIds); i++ {
			id := hexutil.Encode(proposalsIds[i][0:len(proposalsIds[i])])
			newProposalsIds = append(newProposalsIds, id)
		}
		return newProposalsIds, nil
	}

}

func (p *Proposals) AllProposalSets(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, page int64, size int64) ([]string, error) {
	method := "allProposalSets"
	data, err := p.abi.Pack(method, big.NewInt(page), big.NewInt(size))

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return []string{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AllProposalSets result", "error", err)
		return []string{}, err
	}

	ret, err := p.abi.Unpack(method, result)
	log.Info("allProposals result", "ret", ret)
	if err != nil {
		return []string{}, err
	}
	if proposalsIds, ok := ret[0].([][4]byte); !ok {
		return []string{}, errors.New("invalid AddressProposalSets result format")
	} else {
		//lenth := len(proposalsIds)
		var newProposalsIds []string
		for i := 0; i < len(proposalsIds); i++ {
			id := hexutil.Encode(proposalsIds[i][0:len(proposalsIds[i])])
			newProposalsIds = append(newProposalsIds, id)
		}
		log.Info("AllProposalSets result", "result", newProposalsIds)
		return newProposalsIds, nil
	}

}

func (p *Proposals) AddressProposals(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page int64, size int64) (*ProposalInfo, error) {
	method := "addressProposals"
	data, err := p.abi.Pack(method, addr, big.NewInt(page), big.NewInt(size))

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return &ProposalInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AddressProposals result", "error", err)
		return &ProposalInfo{}, err
	}
	proposalInfo := &ProposalInfo{}
	err = p.abi.UnpackIntoInterface(proposalInfo, method, result)
	if err != nil {
		return &ProposalInfo{}, err
	}
	return proposalInfo, nil
}

func (p *Proposals) AllProposals(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, page int64, size int64) ([]ProposalInfo, error) {
	method := "allProposals"
	data, err := p.abi.Pack(method, big.NewInt(page), big.NewInt(size))

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return []ProposalInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("allProposals result", "error", err)
		return []ProposalInfo{}, err
	}
	ret, err := p.abi.Unpack(method, result)
	if err != nil {
		log.Error("allProposals result", "error", err)
		return []ProposalInfo{}, err
	}
	log.Info("allProposals result", "ret", ret[0])
	if proposalsInfos, ok := ret[0].([]ProposalInfo); !ok {
		return []ProposalInfo{}, errors.New("invalid AllProposals result format")
	} else {
		//lenth := len(proposalsIds)
		log.Info("allProposals result", "result", proposalsInfos)
		return proposalsInfos, nil
	}

}
