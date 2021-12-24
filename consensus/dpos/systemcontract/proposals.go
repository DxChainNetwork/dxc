package systemcontract

import (
	"errors"
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

type Proposals struct {
	abi          abi.ABI
	contractAddr common.Address
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
	if proposalsIds, ok := ret[0].([][4]byte); !ok {
		return [][4]byte{}, errors.New("invalid AddressProposalSets result format")
	} else {
		log.Info("AddressProposalSets result", "address", addr.String(), "result", proposalsIds)
		return proposalsIds, nil
	}

}
