package systemcontract

import (
	"github.com/DxChainNetwork/dxc/common"
	"math/big"
)

// using for Validators contract's initialize
var (
	InitValAddress = common.HexToAddress("0x1c0e8eaf42ec8d4010e960313248d2af95be7d34")
	InitRate       = uint8(70)
	InitDeposit    = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))
)
