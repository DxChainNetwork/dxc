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

type Validators struct {
	abi          abi.ABI
	contractAddr common.Address
}

func NewValidators() *Validators {
	return &Validators{
		abi:          abiMap[ValidatorsContractName],
		contractAddr: ValidatorsContractAddr,
	}
}

func (v *Validators) GetCurrentEpochValidators(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) ([]common.Address, error) {
	method := "getCurEpochValidators"
	data, err := v.abi.Pack(method)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return []common.Address{}, err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return []common.Address{}, err
	}

	vals, ok := ret[0].([]common.Address)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return []common.Address{}, err
	}

	return vals, nil
}

func (v *Validators) GetEffictiveValidators(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) ([]common.Address, error) {
	method := "getEffictiveValidators"
	data, err := v.abi.Pack(method)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return []common.Address{}, err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return []common.Address{}, err
	}
	vals, ok := ret[0].([]common.Address)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return []common.Address{}, err
	}

	return vals, nil
}
