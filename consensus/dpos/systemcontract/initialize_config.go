package systemcontract

import (
	"github.com/DxChainNetwork/dxc/common"
	"math/big"
)

// using for Validators contract's initialize
var (
	InitValAddress = common.HexToAddress("0xA88548E97AF8809aFaC9dC7a930650c117951059")
	InitRate       = uint8(70)
	InitDeposit    = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1000))
)
