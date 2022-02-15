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
		msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, nonce, value, math.MaxUint64, new(big.Int), data, true)
		if _, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config); err != nil {
			log.Error("InitProposal execute", "error", err)
			return err
		}
	}

	return nil
}

// AddressProposalSets function AddressProposalSets
func (p *Proposals) AddressProposalSets(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page *big.Int, size *big.Int) ([][4]byte, error) {
	method := "addressProposalSets"
	data, err := p.abi.Pack(method, addr, page, size)

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return [][4]byte{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AddressProposalSets result", "error", err)
		return [][4]byte{}, err
	}

	ret, err := p.abi.Unpack(method, result)
	if err != nil {
		return [][4]byte{}, err
	}
	if proposalIds, ok := ret[0].([][4]byte); !ok {
		return [][4]byte{}, errors.New("invalid AddressProposalSets result format")
	} else {
		return proposalIds, nil
	}

}

// AllProposalSets function AllProposalSets
func (p *Proposals) AllProposalSets(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, page *big.Int, size *big.Int) ([][4]byte, error) {
	method := "allProposalSets"
	data, err := p.abi.Pack(method, page, size)

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return [][4]byte{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AllProposalSets result", "error", err)
		return [][4]byte{}, err
	}

	ret, err := p.abi.Unpack(method, result)
	if err != nil {
		return [][4]byte{}, err
	}
	if proposalIds, ok := ret[0].([][4]byte); !ok {
		return [][4]byte{}, errors.New("invalid AddressProposalSets result format")
	} else {
		return proposalIds, nil
	}
}

// AddressProposals function AddressProposals
func (p *Proposals) AddressProposals(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page *big.Int, size *big.Int) ([]ProposalInfo, error) {
	method := "addressProposals"
	data, err := p.abi.Pack(method, addr, page, size)

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return []ProposalInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AddressProposals result", "error", err)
		return []ProposalInfo{}, err
	}
	var proposalInfo []ProposalInfo
	err = p.abi.UnpackIntoInterface(&proposalInfo, method, result)
	if err != nil {
		log.Error("AddressProposals Unpack", "error", err)
		return []ProposalInfo{}, err
	}
	return proposalInfo, nil
}

// AllProposals function AllProposals
func (p *Proposals) AllProposals(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, page *big.Int, size *big.Int) ([]ProposalInfo, error) {
	method := "allProposals"
	data, err := p.abi.Pack(method, page, size)

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return []ProposalInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AllProposals result", "error", err)
		return []ProposalInfo{}, err
	}
	var proposalInfo []ProposalInfo
	err = p.abi.UnpackIntoInterface(&proposalInfo, method, result)
	if err != nil {
		log.Error("AllProposals Unpack", "error", err)
		return []ProposalInfo{}, err
	}
	return proposalInfo, nil
}

// AddressProposalCount function AddressProposalCount
func (p *Proposals) AddressProposalCount(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address) (*big.Int, error) {
	method := "addressProposalCount"
	data, err := p.abi.Pack(method, addr)

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AddressProposalCount result", "error", err)
		return big.NewInt(0), err
	}

	ret, err := p.abi.Unpack(method, result)
	if err != nil {
		return big.NewInt(0), err
	}
	if count, ok := ret[0].(*big.Int); !ok {
		return big.NewInt(0), errors.New("invalid AddressProposalCount result format")
	} else {
		return count, nil
	}
}

// ProposalCount function ProposalCount
func (p *Proposals) ProposalCount(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (*big.Int, error) {
	method := "proposalCount"
	data, err := p.abi.Pack(method)

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("AllProposalSets result", "error", err)
		return big.NewInt(0), err
	}

	ret, err := p.abi.Unpack(method, result)
	if err != nil {
		return big.NewInt(0), err
	}
	if count, ok := ret[0].(*big.Int); !ok {
		return big.NewInt(0), errors.New("invalid AddressProposalSets result format")
	} else {
		return count, nil
	}
}

// GetProposal function GetProposal
func (p *Proposals) GetProposal(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, id string) (*ProposalInfo, error) {
	method := "proposalInfos"
	idBytes, err := hexutil.Decode(id)
	if err != nil {
		return &ProposalInfo{}, err
	}
	var idByte4 [4]byte
	copy(idByte4[:], idBytes[:4])
	data, err := p.abi.Pack(method, idByte4)

	if err != nil {
		log.Error("can't pack Proposals contract method", "method", method)
		return &ProposalInfo{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &p.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("GetProposal result", "error", err)
		return &ProposalInfo{}, err
	}
	proposalInfo := &ProposalInfo{}
	err = p.abi.UnpackIntoInterface(proposalInfo, method, result)
	if err != nil {
		log.Error("GetProposal Unpack", "error", err)
		return &ProposalInfo{}, err
	}
	return proposalInfo, nil
}
