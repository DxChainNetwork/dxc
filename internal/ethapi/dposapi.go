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
func (pd *PublicDposTxAPI) prepareAccount(args *TransactionArgs) error {
	if args.From != nil && *args.From != common.BigToAddress(common.Big0) {
		account := accounts.Account{Address: *args.From}
		_, err := pd.b.AccountManager().Find(account)
		if err != nil {
			return err
		}
		return nil
	}

	if wallets := pd.b.AccountManager().Wallets(); len(wallets) > 0 {
		if walletAccounts := wallets[0].Accounts(); len(walletAccounts) > 0 {
			args.From = &walletAccounts[0].Address
			return nil
		}
	}
	return errors.New("need unlock or add wallet to send transaction")
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
func (pd *PublicDposTxAPI) InitProposal(pType uint8, rate uint8, details string, args *TransactionArgs) (common.Hash, error) {
	ctx := context.Background()
	args.To = &systemcontract.ProposalsContractAddr

	if err := pd.prepareAccount(args); err != nil {
		return common.Hash{}, err
	}

	pd.nonceLock.LockAddr(*args.From)
	defer pd.nonceLock.UnlockAddr(*args.From)

	log.Info("initProposal", "ptype", pType, "rate", rate, "details", details)

	method := "initProposal"
	abiMap := systemcontract.GetInteractiveABI()

	data, err := abiMap[systemcontract.ProposalsContractName].Pack(method, pType, rate, details)
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

// UpdateProposal updateProposal function of Proposal contract
func (pd *PublicDposTxAPI) UpdateProposal(id string, rate uint8, deposit *big.Int, details string, args *TransactionArgs) (common.Hash, error) {
	ctx := context.Background()
	args.To = &systemcontract.ProposalsContractAddr

	if err := pd.prepareAccount(args); err != nil {
		return common.Hash{}, err
	}

	pd.nonceLock.LockAddr(*args.From)
	defer pd.nonceLock.UnlockAddr(*args.From)

	log.Info("updateProposal", "id", id, "rate", rate, "deposit", deposit, "details", details)

	method := "updateProposal"
	abiMap := systemcontract.GetInteractiveABI()

	idBytes, err := hexutil.Decode(id)
	if err != nil {
		return common.Hash{}, err
	}
	var idByte4 [4]byte
	copy(idByte4[:], idBytes[:4])

	data, err := abiMap[systemcontract.ProposalsContractName].Pack(method, idByte4, rate, deposit, details)
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

// CancelProposal cancelProposal function of Proposal contract
func (pd *PublicDposTxAPI) CancelProposal(id string, args *TransactionArgs) (common.Hash, error) {
	ctx := context.Background()
	args.To = &systemcontract.ProposalsContractAddr

	if err := pd.prepareAccount(args); err != nil {
		return common.Hash{}, err
	}

	pd.nonceLock.LockAddr(*args.From)
	defer pd.nonceLock.UnlockAddr(*args.From)

	log.Info("cancelProposal", "id", id)

	method := "cancelProposal"
	abiMap := systemcontract.GetInteractiveABI()

	idBytes, err := hexutil.Decode(id)
	if err != nil {
		return common.Hash{}, err
	}
	var idByte4 [4]byte
	copy(idByte4[:], idBytes[:4])

	data, err := abiMap[systemcontract.ProposalsContractName].Pack(method, idByte4)
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

// Guarantee guarantee function of Proposal contract
func (pd *PublicDposTxAPI) Guarantee(id string, args *TransactionArgs) (common.Hash, error) {
	ctx := context.Background()
	args.To = &systemcontract.ProposalsContractAddr

	if err := pd.prepareAccount(args); err != nil {
		return common.Hash{}, err
	}

	pd.nonceLock.LockAddr(*args.From)
	defer pd.nonceLock.UnlockAddr(*args.From)

	log.Info("Guarantee Proposal", "id", id)

	method := "guarantee"
	abiMap := systemcontract.GetInteractiveABI()

	idBytes, err := hexutil.Decode(id)
	if err != nil {
		return common.Hash{}, err
	}
	var idByte4 [4]byte
	copy(idByte4[:], idBytes[:4])

	data, err := abiMap[systemcontract.ProposalsContractName].Pack(method, idByte4)
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
