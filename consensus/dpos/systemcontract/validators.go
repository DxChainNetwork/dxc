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

type Validator struct {
	Status                  uint8
	Deposit                 *big.Int
	Rate                    uint8
	Name                    string
	Details                 string
	Votes                   *big.Int
	UnstakeLockingEndBlock  *big.Int
	RateSettLockingEndBlock *big.Int
}

func (v *Validators) GetCurrentEpochValidators(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) ([]common.Address, error) {
	method := "getCurEpochValidators"
	data, err := v.abi.Pack(method)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return []common.Address{}, err
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

func (v *Validators) GetEffictiveValidatorsWithPage(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, page *big.Int, size *big.Int) ([]common.Address, error) {
	method := "getEffictiveValidatorsWithPage"
	data, err := v.abi.Pack(method, page, size)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return []common.Address{}, err
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

// GetValidatorVoters get the address voters
func (v *Validators) GetValidatorVoters(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, page *big.Int, size *big.Int) ([]common.Address, error) {
	method := "getValidatorVoters"
	data, err := v.abi.Pack(method, addr, page, size)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return []common.Address{}, err
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
	voters, ok := ret[0].([]common.Address)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return []common.Address{}, err
	}

	return voters, nil
}

// GetInvalidValidatorsWithPage get invalid validators
func (v *Validators) GetInvalidValidatorsWithPage(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, page *big.Int, size *big.Int) ([]common.Address, error) {
	method := "getInvalidValidatorsWithPage"
	data, err := v.abi.Pack(method, page, size)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return []common.Address{}, err
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

// GetCancelQueueValidators get canceling queue validators
func (v *Validators) GetCancelQueueValidators(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) ([]common.Address, error) {
	method := "getCancelQueueValidators"
	data, err := v.abi.Pack(method)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return []common.Address{}, err
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

// EffictiveValsLength function EffictiveValsLength
func (v *Validators) EffictiveValsLength(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (*big.Int, error) {
	method := "effictiveValsLength"
	data, err := v.abi.Pack(method)

	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	count, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return count, nil

}

// InvalidValsLength function InvalidValsLength
func (v *Validators) InvalidValsLength(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (*big.Int, error) {
	method := "invalidValsLength"
	data, err := v.abi.Pack(method)

	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	count, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return count, nil

}

// CancelQueueValidatorsLength function CancelQueueValidatorsLength
func (v *Validators) CancelQueueValidatorsLength(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (*big.Int, error) {
	method := "CancelQueueValidatorsLength"
	data, err := v.abi.Pack(method)

	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	count, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return count, nil

}

// ValidatorVotersLength function validatorVotersLength
func (v *Validators) ValidatorVotersLength(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address) (*big.Int, error) {
	method := "validatorVotersLength"
	data, err := v.abi.Pack(method, addr)

	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	count, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return count, nil

}

// IsEffictiveValidator function IsEffictiveValidator
func (v *Validators) IsEffictiveValidator(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address) (bool, error) {
	method := "isEffictiveValidator"
	data, err := v.abi.Pack(method, addr)

	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return false, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return false, err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return false, err
	}
	val, ok := ret[0].(bool)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return false, err
	}
	return val, nil

}

// GetValidator function GetValidator
func (v *Validators) GetValidator(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address) (*Validator, error) {
	method := "validators"
	data, err := v.abi.Pack(method, addr)
	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return &Validator{}, err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return &Validator{}, err
	}
	val := &Validator{}
	err = v.abi.UnpackIntoInterface(val, method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return &Validator{}, err
	}
	return val, nil
}

// TotalDeposit function TotalDeposit
func (v *Validators) TotalDeposit(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (*big.Int, error) {
	method := "totalDeposit"
	data, err := v.abi.Pack(method)

	if err != nil {
		log.Error("Validators Pack error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &v.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("Validators contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	ret, err := v.abi.Unpack(method, result)
	if err != nil {
		log.Error("Validators contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	deposit, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("Validators contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	return deposit, nil

}
