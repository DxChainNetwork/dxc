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

func (b *Base) GetBaseInfos(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (map[string]interface{}, error) {
	baseInfos := map[string]interface{}{}

	uint256Methods := []string{"BLOCK_SECONDS", "EPOCH_BLOCKS", "MIN_DEPOSIT",
		"MAX_VALIDATORS_COUNT", "MAX_VALIDATOR_DETAIL_LENGTH", "MAX_PUNISH_COUNT",
		"RATE_SET_LOCK_EPOCHS", "VALIDATOR_UNSTAKE_LOCK_EPOCHS", "PROPOSAL_DURATION_EPOCHS",
		"VALIDATOR_REWARD_LOCK_EPOCHS", "VOTE_CANCEL_EPOCHS", "TOTAL_DEPOSIT_LV1",
		"TOTAL_DEPOSIT_LV2", "TOTAL_DEPOSIT_LV3", "TOTAL_DEPOSIT_LV4", "TOTAL_DEPOSIT_LV5",
		"REWARD_DEPOSIT_UNDER_LV1", "REWARD_DEPOSIT_FROM_LV1_TO_LV2", "REWARD_DEPOSIT_FROM_LV2_TO_LV3",
		"REWARD_DEPOSIT_FROM_LV3_TO_LV4", "REWARD_DEPOSIT_FROM_LV4_TO_LV5", "REWARD_DEPOSIT_OVER_LV5",
		"MAX_VALIDATOR_COUNT_LV1", "MAX_VALIDATOR_COUNT_LV2", "MAX_VALIDATOR_COUNT_LV3", "MAX_VALIDATOR_COUNT_LV4",
		"MIN_LEVEL_VALIDATOR_COUNT", "MEDIUM_LEVEL_VALIDATOR_COUNT", "MAX_LEVEL_VALIDATOR_COUNT", "SAFE_MULTIPLIER"}

	uin8Methods := []string{"MIN_RATE", "MAX_RATE"}

	addressMethods := []string{"BLACK_HOLE_ADDRESS"}

	for i := 0; i < len(uint256Methods); i++ {
		method := uint256Methods[i]
		data, err := b.abi.Pack(method)

		if err != nil {
			log.Error("can't pack Base contract method", "method", method)
			return map[string]interface{}{}, err
		}

		msg := vmcaller.NewLegacyMessage(header.Coinbase, &b.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
		result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
		if err != nil {
			log.Error("GetBaseInfos result", "error", err)
			return map[string]interface{}{}, err
		}

		//unpack data
		ret, err := b.abi.Unpack(method, result)
		if err != nil {
			log.Error("GetBaseInfos Unpack", "error", err, "result", result)
			return map[string]interface{}{}, err
		}
		info, ok := ret[0].(*big.Int)
		if !ok {
			return map[string]interface{}{}, errors.New("invalid minRate format")
		}
		baseInfos[method] = info
	}

	for i := 0; i < len(uin8Methods); i++ {
		method := uin8Methods[i]
		data, err := b.abi.Pack(method)

		if err != nil {
			log.Error("can't pack Base contract method", "method", method)
			return map[string]interface{}{}, err
		}

		msg := vmcaller.NewLegacyMessage(header.Coinbase, &b.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
		result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
		if err != nil {
			log.Error("GetBaseInfos result", "error", err)
			return map[string]interface{}{}, err
		}

		//unpack data
		ret, err := b.abi.Unpack(method, result)
		if err != nil {
			log.Error("GetBaseInfos Unpack", "error", err, "result", result)
			return map[string]interface{}{}, err
		}
		info, ok := ret[0].(uint8)
		if !ok {
			return map[string]interface{}{}, errors.New("invalid minRate format")
		}
		baseInfos[method] = info
	}

	for i := 0; i < len(addressMethods); i++ {
		method := addressMethods[i]
		data, err := b.abi.Pack(method)

		if err != nil {
			log.Error("can't pack Base contract method", "method", method)
			return map[string]interface{}{}, err
		}

		msg := vmcaller.NewLegacyMessage(header.Coinbase, &b.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
		result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
		if err != nil {
			log.Error("GetBaseInfos result", "error", err)
			return map[string]interface{}{}, err
		}

		//unpack data
		ret, err := b.abi.Unpack(method, result)
		if err != nil {
			log.Error("GetBaseInfos Unpack", "error", err, "result", result)
			return map[string]interface{}{}, err
		}
		info, ok := ret[0].(common.Address)
		if !ok {
			return map[string]interface{}{}, errors.New("invalid minRate format")
		}
		baseInfos[method] = info
	}

	return baseInfos, nil
}
