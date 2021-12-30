package systemcontract

import (
	"github.com/DxChainNetwork/dxc/common"
	"math/big"
)

// using for Validators contract's initialize
var (
	InitValAddress = common.HexToAddress("0x6b246a131e2c59108c841e9bc725acad3d3ee3f0")
	InitRate       = uint8(70)
	InitDeposit    = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))
)
