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

type SystemRewards struct {
	abi          abi.ABI
	contractAddr common.Address
}

// EpochInfo struct `epoch` in SystemRewards contract
type EpochInfo struct {
	BlockReward    *big.Int
	Tvl            *big.Int
	ValidatorCount *big.Int
}

type Reward struct {
	ValidatorReward  *big.Int
	DelegatorsReward *big.Int
	Rate             uint8
}

type SysRewardsInfo struct {
	Epochs            []*big.Int
	ValidatorRewards  []*big.Int
	DelegatorsRewards []*big.Int
	Rates             []uint8
	PendingReward     *big.Int
	FrozenReward      *big.Int
	RewardPerVote     *big.Int
}

// NewSystemRewards return SystemRewards contract instance
func NewSystemRewards() *SystemRewards {
	return &SystemRewards{
		abi:          abiMap[SystemRewardsContractName],
		contractAddr: SystemRewardsContractAddr,
	}
}

// GetEpochInfo return epoch info
func (s *SystemRewards) GetEpochInfo(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, epoch *big.Int) (*EpochInfo, error) {
	method := "epochs"

	data, err := s.abi.Pack(method, epoch)
	if err != nil {
		log.Error("can't pack SystemRewards contract method", "method", method)
		return &EpochInfo{}, err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &s.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("SystemRewards contract execute error", "method", method, "error", err)
		return &EpochInfo{}, err
	}
	epochInfo := &EpochInfo{}
	err = s.abi.UnpackIntoInterface(epochInfo, method, result)
	if err != nil {
		log.Error("SystemRewards contract Unpack error", "method", method, "error", err, "result", result)
		return &EpochInfo{}, err
	}

	return epochInfo, nil
}

// KickoutInfo return kickout addresses in epoch
func (s *SystemRewards) KickoutInfo(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, epoch *big.Int) ([]common.Address, error) {
	method := "kickoutInfo"

	data, err := s.abi.Pack(method, epoch)
	if err != nil {
		log.Error("can't pack SystemRewards contract method", "method", method)
		return []common.Address{}, err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &s.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("SystemRewards contract execute error", "method", method, "error", err)
		return []common.Address{}, err
	}

	ret, err := s.abi.Unpack(method, result)
	if err != nil {
		log.Error("SystemRewards contract Unpack error", "method", method, "error", err, "result", result)
		return []common.Address{}, err
	}

	vals, ok := ret[0].([]common.Address)
	if !ok {
		log.Error("SystemRewards contract format result error", "method", method, "error", err)
		return []common.Address{}, err
	}

	return vals, nil
}

// GetValRewardInfoByEpoch return the address and the epoch reward info
func (s *SystemRewards) GetValRewardInfoByEpoch(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, epoch *big.Int) (*Reward, error) {
	method := "validatorEpochRewardInfo"

	data, err := s.abi.Pack(method, addr, epoch)
	if err != nil {
		log.Error("can't pack SystemRewards contract method", "method", method)
		return &Reward{}, err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &s.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("SystemRewards contract execute error", "method", method, "error", err)
		return &Reward{}, err
	}
	rewards := &Reward{}
	err = s.abi.UnpackIntoInterface(rewards, method, result)
	if err != nil {
		log.Error("SystemRewards contract Unpack error", "method", method, "error", err, "result", result)
		return &Reward{}, err
	}

	return rewards, nil
}

// PendingValidatorReward return the address reward
func (s *SystemRewards) PendingValidatorReward(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address) (*big.Int, *big.Int, error) {
	method := "pendingValidatorReward"
	data, err := s.abi.Pack(method, addr)
	if err != nil {
		log.Error("can't pack SystemRewards contract method", "method", method)
		return big.NewInt(0), big.NewInt(0), err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &s.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("SystemRewards contract execute error", "method", method, "error", err)
		return big.NewInt(0), big.NewInt(0), err
	}
	ret, err := s.abi.Unpack(method, result)
	if err != nil {
		log.Error("SystemRewards contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), big.NewInt(0), err
	}
	avaliable, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("SystemRewards contract format result error", "method", method, "error", err)
		return big.NewInt(0), big.NewInt(0), err
	}
	frozen, ok := ret[1].(*big.Int)
	if !ok {
		log.Error("SystemRewards contract format result error", "method", method, "error", err)
		return big.NewInt(0), big.NewInt(0), err
	}
	return avaliable, frozen, nil
}

// ValidatorRewardsInfo return the address sys rewards
func (s *SystemRewards) ValidatorRewardsInfo(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address) (*SysRewardsInfo, error) {
	method := "validatorRewardsInfo"

	data, err := s.abi.Pack(method, addr)
	if err != nil {
		log.Error("can't pack SystemRewards contract method", "method", method)
		return &SysRewardsInfo{}, err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &s.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("SystemRewards contract execute error", "method", method, "error", err)
		return &SysRewardsInfo{}, err
	}
	rewardInfo := &SysRewardsInfo{}
	err = s.abi.UnpackIntoInterface(rewardInfo, method, result)
	if err != nil {
		log.Error("SystemRewards contract Unpack error", "method", method, "error", err, "result", result)
		return &SysRewardsInfo{}, err
	}

	return rewardInfo, nil
}

// PunishInfo punishInfo function of systemRewards contract
func (s *SystemRewards) PunishInfo(statedb *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig, addr common.Address, epoch *big.Int) (*big.Int, error) {
	method := "punishInfo"

	data, err := s.abi.Pack(method, addr, epoch)
	if err != nil {
		log.Error("can't pack SystemRewards contract method", "method", method)
		return big.NewInt(0), err
	}
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &s.contractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, chainContext, config)
	if err != nil {
		log.Error("SystemRewards contract execute error", "method", method, "error", err)
		return big.NewInt(0), err
	}
	ret, err := s.abi.Unpack(method, result)
	if err != nil {
		log.Error("SystemRewards contract Unpack error", "method", method, "error", err, "result", result)
		return big.NewInt(0), err
	}
	punishCount, ok := ret[0].(*big.Int)
	if !ok {
		log.Error("SystemRewards contract format result error", "method", method, "error", err)
		return big.NewInt(0), err
	}

	return punishCount, nil
}
