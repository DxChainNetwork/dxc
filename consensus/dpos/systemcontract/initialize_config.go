package systemcontract

import (
	"math/big"
)

// using for Validators contract's initialize
var (
	InitRate    = uint8(70)
	InitDeposit = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1000))
)
