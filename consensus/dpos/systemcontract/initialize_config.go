package systemcontract

import (
	"math/big"
)

// using for Validators contract's initialize
var (
	InitRate    = uint8(100)
	InitDeposit = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1000))
	InitName    = "Validator1"           // max bytes length: 100
	InitDetails = "initialize validator" // max bytes length: 10000
)
