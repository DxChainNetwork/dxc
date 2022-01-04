package types

import "github.com/DxChainNetwork/dxc/common"

// EvmExtraValidator contains some extra validations to a transaction,
// and the validator is used inside the evm.
type EvmExtraValidator interface {
	// IsAddressDenied returns whether an address is denied.
	IsAddressDenied(address common.Address, cType common.AddressCheckType) bool
	// IsLogDenied returns whether a log (contract event) is denied.
	IsLogDenied(log *Log) bool
}
