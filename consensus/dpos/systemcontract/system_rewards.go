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
	log.Info("GetEpochInfo", "epoch", epoch)
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

	log.Info("SystemRewards contract execute result", "method", method, "blockReward", epochInfo.BlockReward)

	return epochInfo, nil
}
