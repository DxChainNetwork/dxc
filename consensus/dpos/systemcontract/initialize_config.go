package systemcontract

import (
	"github.com/DxChainNetwork/dxc/common"
	"math/big"
)

// using for Validators contract's initialize
var (
	InitValAddress = common.HexToAddress("0x513A32afC0De7eE89CE8Dc8320B1F4535b539387")
	InitRate       = uint8(70)
	InitDeposit    = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))
)
