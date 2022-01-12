package ethapi

import (
	"context"
	"errors"
	"github.com/DxChainNetwork/dxc/accounts"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/common/hexutil"
	"github.com/DxChainNetwork/dxc/consensus/dpos/systemcontract"
	"github.com/DxChainNetwork/dxc/log"
	"math/big"
)

var (
	errFormatFrom  = errors.New("format from address error")
	errFormatValue = errors.New("format value error")
	errFormatGas   = errors.New("format gas fee error")
	errFormatNonce = errors.New("format nonce error")
)

// PublicDposTxAPI exposes the dpos tx methods for the RPC interface
type PublicDposTxAPI struct {
	b         Backend
	nonceLock *AddrLocker
}

// NewPublicDposTxAPI construct a PublicDposTxAPI object
func NewPublicDposTxAPI(b Backend, nonceLock *AddrLocker) *PublicDposTxAPI {
	return &PublicDposTxAPI{b, nonceLock}
}

// prepareAccount check from address in args
func (pd *PublicDposTxAPI) prepareAccount(args *TransactionArgs, from *common.Address) error {
	if from != nil {
		account := accounts.Account{Address: *from}
		_, err := pd.b.AccountManager().Find(account)
		if err != nil {
			return err
		}
		args.From = from
		return nil
	}

	if wallets := pd.b.AccountManager().Wallets(); len(wallets) > 0 {
		if walletAccounts := wallets[0].Accounts(); len(walletAccounts) > 0 {
			*args.From = walletAccounts[0].Address
		}
	}
	return errors.New("need unlock or add wallet to send transaction")
}

func (pd *PublicDposTxAPI) prepareTxArgs(args *TransactionArgs, fields map[string]interface{}) error {
	fromaddr := common.Address{}
	if from, ok := fields["from"]; ok {
		if fromStr, ok := from.(string); ok {
			fromaddr = common.HexToAddress(fromStr)
		} else {
			log.Error("prepareTxArgs", "from", from)
			return errFormatFrom
		}
	}
	if err := pd.prepareAccount(args, &fromaddr); err != nil {
		return err
	}

	if value, ok := fields["value"]; ok {
		if valueStr, ok := value.(string); ok {
			value, ok := big.NewInt(0).SetString(valueStr, 10)
			if !ok {
				return errFormatValue
			}
			args.Value = (*hexutil.Big)(value)
		} else {
			log.Error("prepareTxArgs", "value", value)
			return errFormatValue
		}
	}

	if gas, ok := fields["gas"]; ok {
		if gasFloat64, ok := gas.(float64); ok {
			gasUint64 := uint64(gasFloat64)
			args.Gas = (*hexutil.Uint64)(&gasUint64)
		} else {
			log.Error("prepareTxArgs", "gas", gas)
			return errFormatGas
		}
	}

	if nonce, ok := fields["nonce"]; ok {
		if nonceFloat64, ok := nonce.(float64); ok {
			nonceUint64 := uint64(nonceFloat64)
			args.Nonce = (*hexutil.Uint64)(&nonceUint64)
		} else {
			log.Error("prepareTxArgs", "nonce", nonce)
			return errFormatNonce
		}
	}

	return nil
}

// sendDposTx sign and submit transaction
func (pd *PublicDposTxAPI) sendDposTx(ctx context.Context, args *TransactionArgs) (common.Hash, error) {
	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, pd.b); err != nil {
		return common.Hash{}, err
	}
	// Assemble the transaction and sign with the wallet
	tx := args.ToTransaction()

	account := accounts.Account{Address: *args.From}
	wallet, err := pd.b.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}

	signed, err := wallet.SignTx(account, tx, pd.b.ChainConfig().ChainID)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.From, "to", args.To, "value", args.Value.ToInt(), "err", err)
		return common.Hash{}, err
	}
	txHash, err := SubmitTransaction(ctx, pd.b, signed)
	if err != nil {
		return common.Hash{}, err
	}
	return txHash, nil
}

// InitProposal initProposal function of Proposal contract
func (pd *PublicDposTxAPI) InitProposal(fields map[string]interface{}) (common.Hash, error) {
	ctx := context.Background()
	args := new(TransactionArgs)
	args.To = &systemcontract.ProposalsContractAddr

	if err := pd.prepareTxArgs(args, fields); err != nil {
		return common.Hash{}, err
	}

	pd.nonceLock.LockAddr(*args.From)
	defer pd.nonceLock.UnlockAddr(*args.From)

	//price, err := pd.b.SuggestGasTipCap(ctx)
	//if err != nil {
	//	return common.Hash{}, err
	//}
	//args.GasPrice = (*hexutil.Big)(price)

	//args.Gas = (*hexutil.Uint64)(uint64(1000000))

	//nonce, err := pd.b.GetPoolNonce(ctx, *args.From)
	//if err != nil {
	//	return common.Hash{}, err
	//}
	//args.Nonce = (*hexutil.Uint64)(&nonce)

	method := "initProposal"
	abiMap := systemcontract.GetInteractiveABI()

	pTypeUint8, rateUint8, detailsStr := uint8(0), uint8(0), ""

	if pType, ok := fields["pType"]; ok {
		if pTypeFloat64, ok := pType.(float64); !ok {
			return common.Hash{}, errors.New("format pType error")
		} else {
			pTypeUint8 = uint8(pTypeFloat64)
		}
	} else {
		return common.Hash{}, errors.New("fields not contain pType")
	}
	if rate, ok := fields["rate"]; ok {
		if rateFloat64, ok := rate.(float64); !ok {
			return common.Hash{}, errors.New("format rate error")
		} else {
			rateUint8 = uint8(rateFloat64)
		}
	} else {
		return common.Hash{}, errors.New("fields not contain rate")
	}

	if datails, ok := fields["datails"]; ok {
		if detailsStr, ok = datails.(string); !ok {
			return common.Hash{}, errors.New("format details error")
		}
	}

	data, err := abiMap[systemcontract.ProposalsContractName].Pack(method, pTypeUint8, rateUint8, detailsStr)
	if err != nil {
		return common.Hash{}, err
	}
	args.Data = (*hexutil.Bytes)(&data)

	txHash, err := pd.sendDposTx(ctx, args)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}
