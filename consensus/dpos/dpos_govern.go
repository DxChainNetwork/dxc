package dpos

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DxChainNetwork/dxc/accounts"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/common/hexutil"
	"github.com/DxChainNetwork/dxc/consensus"
	"github.com/DxChainNetwork/dxc/consensus/dpos/systemcontract"
	"github.com/DxChainNetwork/dxc/consensus/dpos/vmcaller"
	"github.com/DxChainNetwork/dxc/core/state"
	"github.com/DxChainNetwork/dxc/core/types"
	"github.com/DxChainNetwork/dxc/core/vm"
	"github.com/DxChainNetwork/dxc/log"
	"github.com/DxChainNetwork/dxc/rlp"
	"math"
	"math/big"
)

// Proposal is the system governance proposal info.
type Proposal struct {
	Id     *big.Int
	Action *big.Int
	From   common.Address
	To     common.Address
	Value  *big.Int
	Data   []byte
}

func (d *Dpos) getPassedProposalCount(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) (uint32, error) {

	method := "getPassedProposalCount"
	data, err := d.abi[systemcontract.SysGovContractName].Pack(method)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalCount", "error", err)
		return 0, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &systemcontract.SysGovContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// use parent
	result, err := vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig)
	if err != nil {
		return 0, err
	}

	// unpack data
	ret, err := d.abi[systemcontract.SysGovContractName].Unpack(method, result)
	if err != nil {
		return 0, err
	}
	if len(ret) != 1 {
		return 0, errors.New("invalid output length")
	}
	count, ok := ret[0].(uint32)
	if !ok {
		return 0, errors.New("invalid count format")
	}

	return count, nil
}

func (d *Dpos) getPassedProposalByIndex(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, idx uint32) (*Proposal, error) {

	method := "getPassedProposalByIndex"
	data, err := d.abi[systemcontract.SysGovContractName].Pack(method, idx)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalByIndex", "error", err)
		return nil, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &systemcontract.SysGovContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// use parent
	result, err := vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig)
	if err != nil {
		return nil, err
	}

	// unpack data
	prop := &Proposal{}
	err = d.abi[systemcontract.SysGovContractName].UnpackIntoInterface(prop, method, result)
	if err != nil {
		return nil, err
	}

	return prop, nil
}

//finishProposalById
func (d *Dpos) finishProposalById(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, id *big.Int) error {
	method := "finishProposalById"
	data, err := d.abi[systemcontract.SysGovContractName].Pack(method, id)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalByIndex", "error", err)
		return err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &systemcontract.SysGovContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// execute message without a transaction
	state.Prepare(common.Hash{}, 0)
	_, err = vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dpos) executeProposal(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, prop *Proposal, totalTxIndex int) (*types.Transaction, *types.Receipt, error) {
	// Even if the miner is not `running`, it's still working,
	// the 'miner.worker' will try to FinalizeAndAssemble a block,
	// in this case, the signTxFn is not set. A `non-miner node` can't execute system governance proposal.
	if d.signTxFn == nil {
		return nil, nil, errors.New("signTxFn not set")
	}

	propRLP, err := rlp.EncodeToBytes(prop)
	if err != nil {
		return nil, nil, err
	}
	//make system governance transaction
	nonce := state.GetNonce(d.validator)
	amout := prop.Value
	if d.chainConfig.IsSophon(header.Number) {
		// fix bug
		amout = new(big.Int)
	}
	tx := types.NewTransaction(nonce, systemcontract.SysGovToAddr, amout, header.GasLimit, new(big.Int), propRLP)
	tx, err = d.signTxFn(accounts.Account{Address: d.validator}, tx, chain.Config().ChainID)
	if err != nil {
		return nil, nil, err
	}
	//add nonce for validator
	state.SetNonce(d.validator, nonce+1)
	receipt := d.executeProposalMsg(chain, header, state, prop, totalTxIndex, tx.Hash(), common.Hash{})

	return tx, receipt, nil
}

func (d *Dpos) replayProposal(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, prop *Proposal, totalTxIndex int, tx *types.Transaction) (*types.Receipt, error) {
	sender, err := types.Sender(d.signer, tx)
	if err != nil {
		return nil, err
	}
	if sender != header.Coinbase {
		return nil, errors.New("invalid sender for system governance transaction")
	}
	propRLP, err := rlp.EncodeToBytes(prop)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(propRLP, tx.Data()) {
		return nil, fmt.Errorf("data missmatch, proposalID: %s, rlp: %s, txHash:%s, txData:%s", prop.Id.String(), hexutil.Encode(propRLP), tx.Hash().String(), hexutil.Encode(tx.Data()))
	}
	//make system governance transaction
	nonce := state.GetNonce(sender)
	//add nonce for validator
	state.SetNonce(sender, nonce+1)
	receipt := d.executeProposalMsg(chain, header, state, prop, totalTxIndex, tx.Hash(), header.Hash())

	return receipt, nil
}

func (d *Dpos) executeProposalMsg(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, prop *Proposal, totalTxIndex int, txHash, bHash common.Hash) *types.Receipt {
	var receipt *types.Receipt
	action := prop.Action.Uint64()
	switch action {
	case 0:
		// evm action.
		receipt = d.executeEvmCallProposal(chain, header, state, prop, totalTxIndex, txHash, bHash)
	case 1:
		// delete code action
		ok := state.Erase(prop.To)
		receipt = types.NewReceipt([]byte{}, ok != true, header.GasUsed)
		log.Info("executeProposalMsg", "action", "erase", "id", prop.Id.String(), "to", prop.To, "txHash", txHash.String(), "success", ok)
	default:
		receipt = types.NewReceipt([]byte{}, true, header.GasUsed)
		log.Warn("executeProposalMsg failed, unsupported action", "action", action, "id", prop.Id.String(), "from", prop.From, "to", prop.To, "value", prop.Value.String(), "data", hexutil.Encode(prop.Data), "txHash", txHash.String())
	}

	receipt.TxHash = txHash
	receipt.BlockHash = bHash
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(state.TxIndex())

	return receipt
}

// the returned value should not nil.
func (d *Dpos) executeEvmCallProposal(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, prop *Proposal, totalTxIndex int, txHash, bHash common.Hash) *types.Receipt {
	// actually run the governance message
	msg := vmcaller.NewLegacyMessage(prop.From, &prop.To, 0, prop.Value, header.GasLimit, new(big.Int), prop.Data, false)
	state.Prepare(txHash, totalTxIndex)
	_, err := vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig)

	// governance message will not actually consumes gas
	receipt := types.NewReceipt([]byte{}, err != nil, header.GasUsed)
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = state.GetLogs(txHash, bHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	log.Info("executeProposalMsg", "action", "evmCall", "id", prop.Id.String(), "from", prop.From, "to", prop.To, "value", prop.Value.String(), "data", hexutil.Encode(prop.Data), "txHash", txHash.String(), "err", err)

	return receipt
}

// Methods for debug trace

// ApplySysTx applies a system-transaction using a given evm,
// the main purpose of this method is for tracing a system-transaction.
func (d *Dpos) ApplySysTx(evm *vm.EVM, state *state.StateDB, txIndex int, sender common.Address, tx *types.Transaction) (ret []byte, vmerr error, err error) {
	var prop = &Proposal{}
	if err = rlp.DecodeBytes(tx.Data(), prop); err != nil {
		return
	}
	evm.Context.ExtraValidator = nil
	nonce := evm.StateDB.GetNonce(sender)
	//add nonce for validator
	evm.StateDB.SetNonce(sender, nonce+1)

	action := prop.Action.Uint64()
	switch action {
	case 0:
		// evm action.
		// actually run the governance message
		msg := vmcaller.NewLegacyMessage(prop.From, &prop.To, 0, prop.Value, tx.Gas(), new(big.Int), prop.Data, false)
		state.Prepare(tx.Hash(), txIndex)
		evm.TxContext = vm.TxContext{
			Origin:   msg.From(),
			GasPrice: new(big.Int).Set(msg.GasPrice()),
		}
		ret, _, vmerr = evm.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())
		state.Finalise(true)
	case 1:
		// delete code action
		_ = state.Erase(prop.To)
	default:
		vmerr = errors.New("unsupported action")
	}
	return
}
