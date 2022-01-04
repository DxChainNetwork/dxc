package systemcontract

import (
	"bytes"
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
	"sort"
)

type ValidatorV0 struct {
	abi          abi.ABI
	contractAddr common.Address
}

func NewValidatorV0() *ValidatorV0 {
	return &ValidatorV0{
		abi:          abiMap[ValidatorsContractName],
		contractAddr: ValidatorsContractAddr,
	}
}

func (v *ValidatorV0) GetTopValidators(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) ([]common.Address, error) {
	method := "getTopValidators"
	data, err := v.abi.Pack(method)
	if err != nil {
		log.Error("Can't pack data for getTopValidators", "error", err)
		return []common.Address{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		return []common.Address{}, err
	}

	// unpack data
	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		return []common.Address{}, err
	}
	if len(ret) != 1 {
		return []common.Address{}, errors.New("invalid params length")
	}
	validators, ok := ret[0].([]common.Address)
	if !ok {
		return []common.Address{}, errors.New("invalid validators format")
	}
	sort.Slice(validators, func(i, j int) bool {
		return bytes.Compare(validators[i][:], validators[j][:]) < 0
	})
	return validators, err
}

func (v *ValidatorV0) GetValidatorFeeAddr(val common.Address, statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (common.Address, error) {
	method := "getValidatorInfo"
	data, err := v.abi.Pack(method, val)
	if err != nil {
		log.Error("Can't pack data for GetValidatorInfo", "error", err)
		return common.Address{}, err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// use parent
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		return common.Address{}, err
	}
	// unpack data
	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		return common.Address{}, err
	}
	feeAddr, ok := ret[0].(common.Address)
	if !ok {
		return common.Address{}, errors.New("invalid output")
	}
	return feeAddr, nil
}
