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

type Base struct {
	abi          abi.ABI
	contractAddr common.Address
}

// NewBase return Base contract instance
func NewBase() *Base {
	return &Base{
		abi:          abiMap[ValidatorsContractName],
		contractAddr: ValidatorsContractAddr,
	}
}

// GetMinDeposit `MIN_DEPOSIT`
func (b *Base) GetMinDeposit(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (*big.Int, error) {
	method := "MIN_DEPOSIT"
	data, err := b.abi.Pack(method)

	if err != nil {
		log.Error("can't pack Base contract method", "method", method)
		return new(big.Int), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &b.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("GetMinDeposit result", "error", err)
		return new(big.Int), err
	}

	//unpack data
	ret, err := b.abi.Unpack(method, result)
	if err != nil {
		log.Error("GetMinDeposit Unpack", "error", err, "result", result)
		return new(big.Int), err
	}
	minDeposit, ok := ret[0].(*big.Int)
	if !ok {
		return new(big.Int), errors.New("invalid minDeposit format")
	}

	log.Info("get Base contract result", "method", method, "result", minDeposit.String())

	return minDeposit, err
}

func (b *Base) GetMinRate(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (uint8, error) {
	method := "MIN_RATE"
	data, err := b.abi.Pack(method)

	if err != nil {
		log.Error("can't pack Base contract method", "method", method)
		return uint8(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &b.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("GetMinRate result", "error", err)
		return uint8(0), err
	}

	//unpack data
	ret, err := b.abi.Unpack(method, result)
	if err != nil {
		log.Error("GetMinRate Unpack", "error", err, "result", result)
		return uint8(0), err
	}
	minRate, ok := ret[0].(uint8)
	if !ok {
		return uint8(0), errors.New("invalid minRate format")
	}

	log.Info("get Base contract result", "method", method, "result", minRate)

	return minRate, err
}
