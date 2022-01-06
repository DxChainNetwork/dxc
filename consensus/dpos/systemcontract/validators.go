package systemcontract

import (
	"github.com/DxChainNetwork/dxc/accounts/abi"
	"github.com/DxChainNetwork/dxc/common"
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
