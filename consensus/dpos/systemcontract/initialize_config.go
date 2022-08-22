package systemcontract

import (
	"encoding/json"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/common/math"
	"math/big"
	"strings"
)

// using for Validators contract's initialize
var (
	InitRate    = uint8(100)
	InitDeposit = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(40000000))
	InitName    = "dxc-validator-1"      // max bytes length: 100
	InitDetails = "initialize validator" // max bytes length: 10000
)

var MigrateOwner = common.HexToAddress("0xaC887b0C4277cDceb741A93c6e8516AcBA880018")

type MigrateAddrBalance struct {
	Address common.Address        `json:"address"`
	Balance *math.HexOrDecimal256 `json:"balance"`
}

var migrateAddrBalanceStr = `[

]`

func InitMigrateAddrBalance() ([]common.Address, []*big.Int) {
	var migrates []MigrateAddrBalance
	err := json.NewDecoder(strings.NewReader(migrateAddrBalanceStr)).Decode(&migrates)
	if err != nil {
		panic(err)
	}
	var addrs []common.Address
	var bals []*big.Int
	for i := 0; i < len(migrates); i++ {
		addrs = append(addrs, migrates[i].Address)
		bals = append(bals, (*big.Int)(migrates[i].Balance))
	}
	return addrs, bals
}
