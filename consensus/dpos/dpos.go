// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package dpos implements the proof-of-stake-authority consensus engine.
package dpos

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/DxChainNetwork/dxc/accounts"
	"github.com/DxChainNetwork/dxc/accounts/abi"
	"github.com/DxChainNetwork/dxc/common"
	"github.com/DxChainNetwork/dxc/consensus"
	"github.com/DxChainNetwork/dxc/consensus/dpos/systemcontract"
	"github.com/DxChainNetwork/dxc/consensus/dpos/vmcaller"
	"github.com/DxChainNetwork/dxc/consensus/misc"
	"github.com/DxChainNetwork/dxc/core/state"
	"github.com/DxChainNetwork/dxc/core/types"
	"github.com/DxChainNetwork/dxc/crypto"
	"github.com/DxChainNetwork/dxc/ethdb"
	"github.com/DxChainNetwork/dxc/log"
	"github.com/DxChainNetwork/dxc/metrics"
	"github.com/DxChainNetwork/dxc/params"
	"github.com/DxChainNetwork/dxc/rlp"
	"github.com/DxChainNetwork/dxc/rpc"
	"github.com/DxChainNetwork/dxc/trie"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory

	wiggleTime    = 500 * time.Millisecond // Random delay (per validator) to allow concurrent validators
	maxValidators = 99                     // Max validators allowed sealing.

	inmemoryBlacklist = 21 // Number of recent blacklist snapshots to keep in memory
)

type blacklistDirection uint

const (
	DirectionFrom blacklistDirection = iota
	DirectionTo
	DirectionBoth
)

// Dpos delegated proof-of-stake protocol constants.
var (
	// TODO: update epochLength
	epochLength = uint64(20) // Default number of blocks after which to checkpoint and reset the pending votes

	extraVanity = 32                     // Fixed number of extra-data prefix bytes reserved for validator vanity
	extraSeal   = crypto.SignatureLength // Fixed number of extra-data suffix bytes reserved for validator seal

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.

	diffInTurn = big.NewInt(2) // Block difficulty for in-turn signatures
	diffNoTurn = big.NewInt(1) // Block difficulty for out-of-turn signatures
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the validator vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	// errExtraValidators is returned if non-checkpoint block contain validator data in
	// their extra-data fields.
	errExtraValidators = errors.New("non-checkpoint block contains extra validator list")

	// errInvalidExtraValidators is returned if validator data in extra-data field is invalid.
	errInvalidExtraValidators = errors.New("invalid extra validators in extra data field")

	// errInvalidCheckpointValidators is returned if a checkpoint block contains an
	// invalid list of validators (i.e. non-divisible by 20 bytes).
	errInvalidCheckpointValidators = errors.New("invalid validator list on checkpoint block")

	// errMismatchingCheckpointValidators is returned if a checkpoint block contains a
	// list of validators different from the one the local node calculated.
	errMismatchingCheckpointValidators = errors.New("mismatching validator list on checkpoint block")

	// errInvalidMixDigest is returned if a block's mix digest is non-zero.
	errInvalidMixDigest = errors.New("non-zero mix digest")

	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")

	// errInvalidDifficulty is returned if the difficulty of a block neither 1 or 2.
	errInvalidDifficulty = errors.New("invalid difficulty")

	// errWrongDifficulty is returned if the difficulty of a block doesn't match the
	// turn of the validator.
	errWrongDifficulty = errors.New("wrong difficulty")

	// errInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")

	// ErrInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")

	// errUnauthorizedValidator is returned if a header is signed by a non-authorized entity.
	errUnauthorizedValidator = errors.New("unauthorized validator")

	// errRecentlySigned is returned if a header is signed by an authorized entity
	// that already signed a header recently, thus is temporarily not allowed to.
	errRecentlySigned = errors.New("recently signed")

	// errInvalidValidatorLen is returned if validators length is zero or bigger than maxValidators.
	errInvalidValidatorsLength = errors.New("invalid validators length")

	// errInvalidCoinbase is returned if the coinbase isn't the validator of the block.
	errInvalidCoinbase = errors.New("invalid coinbase")

	errInvalidSysGovCount = errors.New("invalid system governance tx count")
)

var (
	getblacklistTimer = metrics.NewRegisteredTimer("dpos/blacklist/get", nil)
	getRulesTimer     = metrics.NewRegisteredTimer("dpos/eventcheckrules/get", nil)
)

// StateFn gets state by the state root hash.
type StateFn func(hash common.Hash) (*state.StateDB, error)

// ValidatorFn hashes and signs the data to be signed by a backing account.
type ValidatorFn func(validator accounts.Account, mimeType string, message []byte) ([]byte, error)
type SignTxFn func(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	// If the signature's already cached, return that
	hash := header.Hash()
	if address, known := sigcache.Get(hash); known {
		return address.(common.Address), nil
	}
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var validator common.Address
	copy(validator[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, validator)
	return validator, nil
}

// Dpos is the proof-of-stake-authority consensus engine proposed to support the
// Ethereum testnet following the Ropsten attacks.
type Dpos struct {
	chainConfig *params.ChainConfig // ChainConfig to execute evm
	config      *params.DposConfig  // Consensus engine configuration parameters
	db          ethdb.Database      // Database to store and retrieve snapshot checkpoints

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	blacklists      *lru.Cache // blacklists caches recent blacklist to speed up transactions validation
	blLock          sync.Mutex // Make sure only get blacklist once for each block
	eventCheckRules *lru.Cache // eventCheckRules caches recent EventCheckRules to speed up log validation
	rulesLock       sync.Mutex // Make sure only get eventCheckRules once for each block

	proposals map[common.Address]bool // Current list of proposals we are pushing

	signer types.Signer // the signer instance to recover tx sender

	validator common.Address // Ethereum address of the signing key
	signFn    ValidatorFn    // Validator function to authorize hashes with
	signTxFn  SignTxFn
	lock      sync.RWMutex // Protects the validator fields

	stateFn StateFn // Function to get state by state root

	abi map[string]abi.ABI // Interactive with system contracts

	chain consensus.ChainHeaderReader // chain is only for reading parent headers when getting blacklist and rules

	// The fields below are for testing only
	fakeDiff bool // Skip difficulty verifications
}

// New creates a Dpos proof-of-stake-authority consensus engine with the initial
// validators set to the ones provided by the user.
func New(chainConfig *params.ChainConfig, db ethdb.Database) *Dpos {
	// Set any missing consensus parameters to their defaults
	conf := *chainConfig.Dpos
	log.Info("epoch info", "epoch", conf.Epoch)
	if conf.Epoch == 0 {
		conf.Epoch = epochLength
	}
	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)
	blacklists, _ := lru.New(inmemoryBlacklist)
	rules, _ := lru.New(inmemoryBlacklist)

	abi := systemcontract.GetInteractiveABI()

	return &Dpos{
		chainConfig:     chainConfig,
		config:          &conf,
		db:              db,
		recents:         recents,
		signatures:      signatures,
		blacklists:      blacklists,
		eventCheckRules: rules,
		proposals:       make(map[common.Address]bool),
		abi:             abi,
		signer:          types.LatestSignerForChainID(chainConfig.ChainID),
	}
}

func (d *Dpos) SetChain(chain consensus.ChainHeaderReader) {
	d.chain = chain
}

// SetStateFn sets the function to get state.
func (d *Dpos) SetStateFn(fn StateFn) {
	d.stateFn = fn
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (d *Dpos) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
	// return ecrecover(header, d.signatures)
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (d *Dpos) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return d.verifyHeader(chain, header, nil)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (d *Dpos) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := d.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (d *Dpos) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()

	// Don't waste time checking blocks from the future
	if header.Time > uint64(time.Now().Unix()) {
		return consensus.ErrFutureBlock
	}
	// Check that the extra-data contains the vanity, validators and signature.
	if len(header.Extra) < extraVanity {
		return errMissingVanity
	}
	if len(header.Extra) < extraVanity+extraSeal {
		return errMissingSignature
	}

	// check extra data
	isEpoch := number%d.config.Epoch == 0
	// Ensure that the extra-data contains a validator list on checkpoint, but none otherwise
	validatorsBytes := len(header.Extra) - extraVanity - extraSeal
	// TODO check
	//if !isEpoch && validatorsBytes != 0 {
	//	return errExtraValidators
	//}
	// Ensure that the validator bytes length is valid
	if isEpoch && validatorsBytes%common.AddressLength != 0 {
		return errExtraValidators
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (common.Hash{}) {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if header.UncleHash != uncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 && header.Difficulty == nil {
		return errInvalidDifficulty
	}
	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	}
	// If all checks passed, validate any special fields for hard forks
	if err := misc.VerifyForkHashes(chain.Config(), header, false); err != nil {
		return err
	}
	// All basic checks passed, verify cascading fields
	return d.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (d *Dpos) verifyCascadingFields(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}

	if parent.Time+d.config.Period > header.Time {
		return ErrInvalidTimestamp
	}

	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	if !chain.Config().IsLondon(header.Number) {
		// Verify BaseFee not present before EIP-1559 fork.
		if header.BaseFee != nil {
			return fmt.Errorf("invalid baseFee before fork: have %d, want <nil>", header.BaseFee)
		}
		if err := misc.VerifyGaslimit(parent.GasLimit, header.GasLimit); err != nil {
			return err
		}
	} else if err := misc.VerifyEip1559Header(chain.Config(), parent, header); err != nil {
		// Verify the header's EIP-1559 attributes.
		return err
	}

	// All basic checks passed, verify the seal and return
	return d.verifySeal(chain, header, parents)
}

// snapshot retrieves the authorization snapshot at a given point in time.
func (d *Dpos) snapshot(chain consensus.ChainHeaderReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)
	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := d.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(d.config, d.signatures, d.db, hash); err == nil {
				log.Trace("Loaded voting snapshot from disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}
		// If we're at the genesis, snapshot the initial state. Alternatively if we're
		// at a checkpoint block without a parent (light client CHT), or we have piled
		// up more headers than allowed to be reorged (chain reinit from a freezer),
		// consider the checkpoint trusted and snapshot it.
		if number == 0 || (number%d.config.Epoch == 0 && (len(headers) > params.FullImmutabilityThreshold || chain.GetHeaderByNumber(number-1) == nil)) {
			checkpoint := chain.GetHeaderByNumber(number)
			if checkpoint != nil {
				hash := checkpoint.Hash()

				validators := make([]common.Address, (len(checkpoint.Extra)-extraVanity-extraSeal)/common.AddressLength)
				for i := 0; i < len(validators); i++ {
					copy(validators[i][:], checkpoint.Extra[extraVanity+i*common.AddressLength:])
				}
				snap = newSnapshot(d.config, d.signatures, number, hash, validators)
				if err := snap.store(d.db); err != nil {
					return nil, err
				}
				log.Info("Stored checkpoint snapshot to disk", "number", number, "hash", hash)
				break
			}
		}
		// No snapshot for this header, gather the header and move backward
		var header *types.Header
		if len(parents) > 0 {
			// If we have explicit parents, pick from there (enforced)
			header = parents[len(parents)-1]
			if header.Hash() != hash || header.Number.Uint64() != number {
				return nil, consensus.ErrUnknownAncestor
			}
			parents = parents[:len(parents)-1]
		} else {
			// No explicit parents (or no more left), reach out to the database
			header = chain.GetHeader(hash, number)
			if header == nil {
				return nil, consensus.ErrUnknownAncestor
			}
		}
		headers = append(headers, header)
		number, hash = number-1, header.ParentHash
	}
	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}
	snap, err := snap.apply(headers, chain, parents)
	if err != nil {
		return nil, err
	}
	d.recents.Add(snap.Hash, snap)

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if err = snap.store(d.db); err != nil {
			return nil, err
		}
		log.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}
	return snap, err
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (d *Dpos) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (d *Dpos) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	return d.verifySeal(chain, header, nil)
}

// verifySeal checks whether the signature contained in the header satisfies the
// consensus protocol requirements. The method accepts an optional list of parent
// headers that aren't yet part of the local blockchain to generate the snapshots
// from.
func (d *Dpos) verifySeal(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := d.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// Resolve the authorization key and check against validators
	signer, err := ecrecover(header, d.signatures)
	if err != nil {
		return err
	}
	if signer != header.Coinbase {
		return errInvalidCoinbase
	}

	if _, ok := snap.Validators[signer]; !ok {
		return errUnauthorizedValidator
	}

	for seen, recent := range snap.Recents {
		if recent == signer {
			// Validator is among recents, only fail if the current block doesn't shift it out
			if limit := uint64(len(snap.Validators)/2 + 1); seen > number-limit {
				return errRecentlySigned
			}
		}
	}

	// Ensure that the difficulty corresponds to the turn-ness of the signer
	if !d.fakeDiff {
		inturn := snap.inturn(header.Number.Uint64(), signer)
		if inturn && header.Difficulty.Cmp(diffInTurn) != 0 {
			return errWrongDifficulty
		}
		if !inturn && header.Difficulty.Cmp(diffNoTurn) != 0 {
			return errWrongDifficulty
		}
	}

	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (d *Dpos) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// If the block isn't a checkpoint, cast a random vote (good enough for now)
	header.Coinbase = d.validator
	header.Nonce = types.BlockNonce{}

	number := header.Number.Uint64()
	snap, err := d.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	log.Info("[Prepare]", "snap.validators", snap.validators(), "snap.number", snap.Number, "current header", number, "header.coinbase", header.Coinbase.String())

	// Set the correct difficulty
	header.Difficulty = calcDifficulty(snap, d.validator)

	// Bail out if we're unauthorized to sign a block
	if _, authorized := snap.Validators[header.Coinbase]; !authorized {
		return errUnauthorizedValidator
	}

	// Ensure the extra data has all its components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]
	// Mix digest is reserved for now, set to empty
	header.MixDigest = common.Hash{}

	// Ensure the timestamp has the correct delay
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Time = parent.Time + d.config.Period
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given.
func (d *Dpos) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs *[]*types.Transaction, uncles []*types.Header, receipts *[]*types.Receipt, systemTxs []*types.Transaction) error {
	// Initialize all system contracts at block 1.
	if header.Number.Cmp(common.Big1) == 0 {
		if err := d.initializeSystemContracts(chain, header, state); err != nil {
			log.Error("Initialize system contracts failed", "err", err)
			return err
		}
	}

	if header.Difficulty.Cmp(diffInTurn) != 0 {
		kickout, err := d.tryPunishValidator(chain, header, state)
		if err != nil {
			panic(err)
		}
		if kickout {
			newValidators, err := d.getCurEpochValidators(chain, header, state)
			if err != nil {
				panic(err)
			}
			log.Info("kickout validator", "header", header.Number.Uint64(), "newValidators", newValidators)

			for _, validator := range newValidators {
				header.Extra = append(header.Extra, validator.Bytes()...)
			}
		}
	}

	// avoid nil pointer
	if txs == nil {
		s := make([]*types.Transaction, 0)
		txs = &s
	}
	if receipts == nil {
		rs := make([]*types.Receipt, 0)
		receipts = &rs
	}

	// do epoch thing at the end, because it will update active validators
	if header.Number.Uint64()%d.config.Epoch == 0 {
		log.Info("[FinalizeAndAssemble]: update epoch", "update", true)
		if err := d.doSomethingAtEpoch(chain, header, state); err != nil {
			panic(err)
		}

		newEpochValidators, err := d.getCurEpochValidators(chain, header, state)
		if err != nil {
			panic(err)
		}
		log.Info("update epoch", "header", header.Number.Uint64(), "newEpochValidators", newEpochValidators)
		header.Extra = header.Extra[:extraVanity]
		for _, validator := range newEpochValidators {
			header.Extra = append(header.Extra, validator.Bytes()...)
		}

		header.Extra = append(header.Extra, make([]byte, extraSeal)...)

		validatorsBytes := make([]byte, len(newEpochValidators)*common.AddressLength)
		for i, validator := range newEpochValidators {
			copy(validatorsBytes[i*common.AddressLength:], validator.Bytes())
		}

		extraSuffix := len(header.Extra) - extraSeal
		if !bytes.Equal(header.Extra[extraVanity:extraSuffix], validatorsBytes) {
			return errInvalidExtraValidators
		}
	} else {
		header.Extra = append(header.Extra, make([]byte, extraSeal)...)
	}

	// deposit block reward
	if err := d.trySendBlockReward(chain, header, state); err != nil {
		return err
	}

	// No block rewards in PoA, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)

	return nil
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (d *Dpos) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (b *types.Block, rs []*types.Receipt, err error) {
	defer func() {
		if err != nil {
			log.Warn("FinalizeAndAssemble failed", "err", err)
		}
	}()
	// Initialize all system contracts at block 1.
	if header.Number.Cmp(common.Big1) == 0 {
		if err := d.initializeSystemContracts(chain, header, state); err != nil {
			panic(err)
		}
	}
	// punish validator if necessary
	if header.Difficulty.Cmp(diffInTurn) != 0 {
		kickout, err := d.tryPunishValidator(chain, header, state)
		if err != nil {
			panic(err)
		}
		if kickout {
			newValidators, err := d.getCurEpochValidators(chain, header, state)
			if err != nil {
				panic(err)
			}
			log.Info("kickout validator", "header", header.Number.Uint64(), "newValidators", newValidators)

			for _, validator := range newValidators {
				header.Extra = append(header.Extra, validator.Bytes()...)
			}
		}
	}
	header.Extra = header.Extra[:extraVanity]
	// do epoch thing at the end, because it will update active validators
	if header.Number.Uint64()%d.config.Epoch == 0 {
		log.Info("[FinalizeAndAssemble]: update epoch", "update", true)
		if err := d.doSomethingAtEpoch(chain, header, state); err != nil {
			panic(err)
		}

		newEpochValidators, err := d.getCurEpochValidators(chain, header, state)
		if err != nil {
			panic(err)
		}
		log.Info("update epoch info", "header", header.Number.Uint64(), "newEpochValidators", newEpochValidators)
		for _, validator := range newEpochValidators {
			header.Extra = append(header.Extra, validator.Bytes()...)
		}
	}
	header.Extra = append(header.Extra, make([]byte, extraSeal)...)

	// deposit block reward
	if err := d.trySendBlockReward(chain, header, state); err != nil {
		panic(err)
	}

	// No block rewards in PoA, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, new(trie.Trie)), receipts, nil
}

func (d *Dpos) trySendBlockReward(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) error {
	if header.Coinbase == common.BigToAddress(big.NewInt(0)) {
		return nil
	}

	s := systemcontract.NewSystemRewards()
	// get Block Reward
	epochInfo, err := s.GetEpochInfo(state, header, newChainContext(chain, d), d.chainConfig, new(big.Int).Div(header.Number, new(big.Int).SetUint64(d.config.Epoch)))
	if err != nil {
		log.Error("GetEpochInfo error ", "error", err)
	}

	log.Info("distributeBlockReward", "BlockReward", epochInfo.BlockReward, "FeeRecoder", state.GetBalance(consensus.FeeRecoder))
	totalReward := new(big.Int).Add(epochInfo.BlockReward, state.GetBalance(consensus.FeeRecoder))
	rewardToFoundation := new(big.Int).Div(new(big.Int).Mul(totalReward, big.NewInt(5)), big.NewInt(100))
	rewardToMiner := new(big.Int).Sub(totalReward, rewardToFoundation)

	state.AddBalance(foundationAddress, rewardToFoundation)
	state.AddBalance(systemcontract.SystemRewardsContractAddr, rewardToMiner)
	log.Info("distributeBlockReward", "foundation", state.GetBalance(foundationAddress), "sysAddr", state.GetBalance(systemcontract.SystemRewardsContractAddr))

	// reset tx fee recoder balance
	state.SetBalance(consensus.FeeRecoder, common.Big0)

	method := "distributeBlockReward"
	data, err := d.abi[systemcontract.SystemRewardsContractName].Pack(method, rewardToMiner)
	if err != nil {
		log.Error("Can't pack data for distributeBlockReward", "err", err)
		return err
	}

	nonce := state.GetNonce(header.Coinbase)
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &systemcontract.SystemRewardsContractAddr, nonce, new(big.Int), math.MaxUint64, new(big.Int), data, true)
	if _, err := vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig); err != nil {
		return err
	}

	return nil
}

func (d *Dpos) tryPunishValidator(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) (bool, error) {
	number := header.Number.Uint64()
	snap, err := d.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return false, err
	}
	validators := snap.validators()
	outTurnValidator := validators[number%uint64(len(validators))]
	// check sigend recently or not
	signedRecently := false
	for _, recent := range snap.Recents {
		if recent == outTurnValidator {
			signedRecently = true
			break
		}
	}
	if !signedRecently {
		if kickout, err := d.punishValidator(outTurnValidator, chain, header, state); err != nil {
			return kickout, err
		}
	}

	return false, nil
}

func (d *Dpos) punishValidator(validator common.Address, chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) (bool, error) {

	method := "punish"
	data, err := d.abi[systemcontract.SystemRewardsContractName].Pack(method, validator)
	if err != nil {
		log.Error("punish failed", "error", err)
		return false, err
	}
	nonce := state.GetNonce(header.Coinbase)
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &systemcontract.SystemRewardsContractAddr, nonce, new(big.Int), math.MaxUint64, new(big.Int), data, true)
	// use parent
	result, err := vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig)
	if err != nil {
		return false, err
	}
	ret, err := d.abi[systemcontract.SystemRewardsContractName].Unpack(method, result)
	if err != nil {
		log.Error("punish unpack failed", "error", err)
		return false, err
	}
	kickout, ok := ret[0].(bool)
	if !ok {
		return false, errors.New("punish result format error")
	}

	log.Info("punish result", "validator", validator.String(), "kickout", kickout)
	return kickout, nil
}

func (d *Dpos) doSomethingAtEpoch(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) error {
	if header.Coinbase == common.BigToAddress(common.Big0) {
		return nil
	}

	data, err := d.abi[systemcontract.ValidatorsContractName].Pack("tryElect")
	if err != nil {
		log.Error("tryElect Pack error", "error", err)
		return err
	}

	nonce := state.GetNonce(header.Coinbase)
	msg := vmcaller.NewLegacyMessage(header.Coinbase, &systemcontract.ValidatorsContractAddr, nonce, new(big.Int), math.MaxInt64, new(big.Int), data, true)
	if _, err := vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig); err != nil {
		log.Error("tryElect execute error", "error", err)
		return err
	}

	return nil
}

// initializeSystemContracts initializes all genesis system contracts.
func (d *Dpos) initializeSystemContracts(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) error {
	if header.Coinbase == common.BigToAddress(big.NewInt(0)) {
		return nil
	}

	snap, err := d.snapshot(chain, 0, header.ParentHash, nil)
	if err != nil {
		return err
	}

	genesisValidators := snap.validators()
	if len(genesisValidators) == 0 || len(genesisValidators) > maxValidators {
		return errInvalidValidatorsLength
	}

	method := "initialize"
	contracts := []struct {
		addr    common.Address
		packFun func() ([]byte, error)
	}{
		{systemcontract.ProposalsContractAddr, func() ([]byte, error) {
			return d.abi[systemcontract.ProposalsContractName].Pack(method, systemcontract.ValidatorsContractAddr)
		}},
		{systemcontract.SystemRewardsContractAddr, func() ([]byte, error) {
			return d.abi[systemcontract.SystemRewardsContractName].Pack(method, systemcontract.ValidatorsContractAddr, systemcontract.NodeVotesContractAddr)
		}},
		{systemcontract.ValidatorsContractAddr, func() ([]byte, error) {
			return d.abi[systemcontract.ValidatorsContractName].Pack(method, systemcontract.ProposalsContractAddr, systemcontract.SystemRewardsContractAddr, systemcontract.NodeVotesContractAddr, systemcontract.InitValAddress, systemcontract.InitDeposit, systemcontract.InitRate)
		}},
		{systemcontract.NodeVotesContractAddr, func() ([]byte, error) {
			return d.abi[systemcontract.NodeVotesContractName].Pack(method, systemcontract.ValidatorsContractAddr, systemcontract.SystemRewardsContractAddr)
		}},
		{systemcontract.AddressListContractAddr, func() ([]byte, error) {
			return d.abi[systemcontract.AddressListContractName].Pack(method, systemcontract.DevAdmin)
		}},
		{systemcontract.AddressListContractAddr, func() ([]byte, error) {
			return d.abi[systemcontract.AddressListContractName].Pack("initializeV2")
		}},
	}

	for _, contract := range contracts {
		data, err := contract.packFun()
		if err != nil {
			return err
		}
		nonce := state.GetNonce(header.Coinbase)

		var msg types.Message

		if contract.addr == systemcontract.ValidatorsContractAddr {
			msg = vmcaller.NewLegacyMessage(header.Coinbase, &contract.addr, nonce, systemcontract.InitDeposit, math.MaxUint64, new(big.Int), data, true)
		} else {
			msg = vmcaller.NewLegacyMessage(header.Coinbase, &contract.addr, nonce, new(big.Int), math.MaxUint64, new(big.Int), data, true)
		}

		if _, err := vmcaller.ExecuteMsg(msg, state, header, newChainContext(chain, d), d.chainConfig); err != nil {
			log.Error("initializeSystemContracts execute error", "contract", contract.addr.String())
			return err
		}
	}

	return nil
}

// get current epoch validators after try elect
func (d *Dpos) getCurEpochValidators(chain consensus.ChainHeaderReader, header *types.Header, statedb *state.StateDB) ([]common.Address, error) {

	method := "getCurEpochValidators"
	data, err := d.abi[systemcontract.ValidatorsContractName].Pack(method)
	if err != nil {
		log.Error("Can't pack data for curEpochValidators", "error", err)
		return []common.Address{}, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &systemcontract.ValidatorsContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, newChainContext(chain, d), d.chainConfig)
	if err != nil {
		return []common.Address{}, err
	}

	// unpack data
	ret, err := d.abi[systemcontract.ValidatorsContractName].Unpack(method, result)
	if err != nil {
		return []common.Address{}, err
	}
	if len(ret) != 1 {
		return []common.Address{}, errors.New("invalid params length")
	}
	validators, ok := ret[0].([]common.Address)
	if !ok {
		return []common.Address{}, errors.New("invalid validators format")
	}
	return validators, err
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (d *Dpos) Authorize(validator common.Address, signFn ValidatorFn, signTxFn SignTxFn) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.validator = validator
	d.signFn = signFn
	d.signTxFn = signTxFn
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (d *Dpos) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	header := block.Header()

	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// For 0-period chains, refuse to seal empty blocks (no reward but would spin sealing)
	if d.config.Period == 0 && len(block.Transactions()) == 0 {
		log.Info("Sealing paused, waiting for transactions")
		return nil
	}
	// Don't hold the val fields for the entire sealing procedure
	d.lock.RLock()
	val, signFn := d.validator, d.signFn
	d.lock.RUnlock()

	// Bail out if we're unauthorized to sign a block
	snap, err := d.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}
	if _, authorized := snap.Validators[val]; !authorized {
		return errUnauthorizedValidator
	}
	// If we're amongst the recent validators, wait for the next block
	for seen, recent := range snap.Recents {
		if recent == val {
			// Validator is among recents, only wait if the current block doesn't shift it out
			if limit := uint64(len(snap.Validators)/2 + 1); number < limit || seen > number-limit {
				log.Info("Signed recently, must wait for others")
				return nil
			}
		}
	}

	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := time.Unix(int64(header.Time), 0).Sub(time.Now()) // nolint: gosimple
	if header.Difficulty.Cmp(diffNoTurn) == 0 {
		// It's not our turn explicitly to sign, delay it a bit
		wiggle := time.Duration(len(snap.Validators)/2+1) * wiggleTime
		delay += time.Duration(rand.Int63n(int64(wiggle)))

		log.Trace("Out-of-turn signing requested", "wiggle", common.PrettyDuration(wiggle))
	}
	// Sign all the things!
	sighash, err := signFn(accounts.Account{Address: val}, accounts.MimetypeDpos, DposRLP(header))
	if err != nil {
		return err
	}
	copy(header.Extra[len(header.Extra)-extraSeal:], sighash)
	// Wait until sealing is terminated or delay timeout.
	log.Trace("Waiting for slot to sign and propagate", "delay", common.PrettyDuration(delay))
	go func() {
		select {
		case <-stop:
			return
		case <-time.After(delay):
		}

		select {
		case results <- block.WithSeal(header):
		default:
			log.Warn("Sealing result is not read by miner", "sealhash", SealHash(header))
		}
	}()

	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have:
// * DIFF_NOTURN(2) if BLOCK_NUMBER % validator_COUNT != validator_INDEX
// * DIFF_INTURN(1) if BLOCK_NUMBER % validator_COUNT == validator_INDEX
func (d *Dpos) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	snap, err := d.snapshot(chain, parent.Number.Uint64(), parent.Hash(), nil)
	if err != nil {
		return nil
	}
	return calcDifficulty(snap, d.validator)
}

func calcDifficulty(snap *Snapshot, validator common.Address) *big.Int {
	if snap.inturn(snap.Number+1, validator) {
		return new(big.Int).Set(diffInTurn)
	}
	return new(big.Int).Set(diffNoTurn)
}

// SealHash returns the hash of a block prior to it being sealed.
func (d *Dpos) SealHash(header *types.Header) common.Hash {
	return SealHash(header)
}

// Close implements consensus.Engine. It's a noop for dpos as there are no background threads.
func (d *Dpos) Close() error {
	return nil
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the validator voting.
func (d *Dpos) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "dpos",
		Version:   "1.0",
		Service:   &API{chain: chain, dpos: d},
		Public:    false,
	}}
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header)
	hasher.Sum(hash[:0])
	return hash
}

// DposRLP returns the rlp bytes which needs to be signed for the proof-of-stake-authority
// sealing. The RLP to sign consists of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func DposRLP(header *types.Header) []byte {
	b := new(bytes.Buffer)
	encodeSigHeader(b, header)
	return b.Bytes()
}

func encodeSigHeader(w io.Writer, header *types.Header) {
	err := rlp.Encode(w, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-crypto.SignatureLength], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

func (d *Dpos) PreHandle(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) error {
	//if d.chainConfig.RedCoastBlock != nil && d.chainConfig.RedCoastBlock.Cmp(header.Number) == 0 {
	//	return systemcontract.ApplySystemContractUpgrade(systemcontract.SysContractV1, state, header, newChainContext(chain, d), d.chainConfig)
	//}
	//if d.chainConfig.SophonBlock != nil && d.chainConfig.SophonBlock.Cmp(header.Number) == 0 {
	//	return systemcontract.ApplySystemContractUpgrade(systemcontract.SysContractV2, state, header, newChainContext(chain, d), d.chainConfig)
	//}
	return nil
}

// IsSysTransaction checks whether a specific transaction is a system transaction.
func (d *Dpos) IsSysTransaction(sender common.Address, tx *types.Transaction, header *types.Header) (bool, error) {
	if tx.To() == nil {
		return false, nil
	}

	to := tx.To()
	if sender == header.Coinbase && *to == systemcontract.SysGovToAddr && tx.GasPrice().Sign() == 0 {
		return true, nil
	}
	// Make sure the miner can NOT call the system contract through a normal transaction.
	if sender == header.Coinbase && *to == systemcontract.SysGovContractAddr {
		return true, nil
	}
	return false, nil
}

// CanCreate determines where a given address can create a new contract.
//
// This will query the system Developers contract, by DIRECTLY to get the target slot value of the contract,
// it means that it's strongly relative to the layout of the Developers contract's state variables
func (d *Dpos) CanCreate(state consensus.StateReader, addr common.Address, height *big.Int) bool {
	if d.chainConfig.IsRedCoast(height) && d.config.EnableDevVerification {
		if isDeveloperVerificationEnabled(state) {
			slot := calcSlotOfDevMappingKey(addr)
			valueHash := state.GetState(systemcontract.AddressListContractAddr, slot)
			// none zero value means true
			return valueHash.Big().Sign() > 0
		}
	}
	return true
}

// ValidateTx do a consensus-related validation on the given transaction at the given header and state.
// the parentState must be the state of the header's parent block.
func (d *Dpos) ValidateTx(sender common.Address, tx *types.Transaction, header *types.Header, parentState *state.StateDB) error {
	// Must use the parent state for current validation,
	// so we must start the validation after redCoastBlock
	if d.chainConfig.RedCoastBlock != nil && d.chainConfig.RedCoastBlock.Cmp(header.Number) < 0 {
		m, err := d.getBlacklist(header, parentState)
		if err != nil {
			return err
		}
		if d, exist := m[sender]; exist && (d != DirectionTo) {
			log.Trace("Hit blacklist", "tx", tx.Hash().String(), "addr", sender.String(), "direction", d)
			return types.ErrAddressDenied
		}
		if to := tx.To(); to != nil {
			if d, exist := m[*to]; exist && (d != DirectionFrom) {
				log.Trace("Hit blacklist", "tx", tx.Hash().String(), "addr", to.String(), "direction", d)
				return types.ErrAddressDenied
			}
		}
	}
	return nil
}

func (d *Dpos) getBlacklist(header *types.Header, parentState *state.StateDB) (map[common.Address]blacklistDirection, error) {
	defer func(start time.Time) {
		getblacklistTimer.UpdateSince(start)
	}(time.Now())

	if v, ok := d.blacklists.Get(header.ParentHash); ok {
		return v.(map[common.Address]blacklistDirection), nil
	}

	d.blLock.Lock()
	defer d.blLock.Unlock()
	if v, ok := d.blacklists.Get(header.ParentHash); ok {
		return v.(map[common.Address]blacklistDirection), nil
	}

	// if the last updates is long ago, we don't need to get blacklist from the contract.
	if d.chainConfig.SophonBlock != nil && header.Number.Cmp(d.chainConfig.SophonBlock) > 0 {
		num := header.Number.Uint64()
		lastUpdated := lastBlacklistUpdatedNumber(parentState)
		if num >= 2 && num > lastUpdated+1 {
			parent := d.chain.GetHeader(header.ParentHash, num-1)
			if parent != nil {
				if v, ok := d.blacklists.Get(parent.ParentHash); ok {
					m := v.(map[common.Address]blacklistDirection)
					d.blacklists.Add(header.ParentHash, m)
					return m, nil
				}
			} else {
				log.Error("Unexpected error when getBlacklist, can not get parent from chain", "number", num, "blockHash", header.Hash(), "parentHash", header.ParentHash)
			}
		}
	}

	// can't get blacklist from cache, try to call the contract
	alABI := d.abi[systemcontract.AddressListContractName]
	get := func(method string) ([]common.Address, error) {
		ret, err := d.commonCallContract(header, parentState, alABI, systemcontract.AddressListContractAddr, method, 1)
		if err != nil {
			log.Error(fmt.Sprintf("%s failed", method), "err", err)
			return nil, err
		}

		blacks, ok := ret[0].([]common.Address)
		if !ok {
			return []common.Address{}, errors.New("invalid blacklist format")
		}
		return blacks, nil
	}
	froms, err := get("getBlacksFrom")
	if err != nil {
		return nil, err
	}
	tos, err := get("getBlacksTo")
	if err != nil {
		return nil, err
	}

	m := make(map[common.Address]blacklistDirection)
	for _, from := range froms {
		m[from] = DirectionFrom
	}
	for _, to := range tos {
		if _, exist := m[to]; exist {
			m[to] = DirectionBoth
		} else {
			m[to] = DirectionTo
		}
	}
	d.blacklists.Add(header.ParentHash, m)
	return m, nil
}

func (d *Dpos) CreateEvmExtraValidator(header *types.Header, parentState *state.StateDB) types.EvmExtraValidator {
	if d.chainConfig.SophonBlock != nil && d.chainConfig.SophonBlock.Cmp(header.Number) < 0 {
		blacks, err := d.getBlacklist(header, parentState)
		if err != nil {
			log.Error("getBlacklist failed", "err", err)
			return nil
		}
		rules, err := d.getEventCheckRules(header, parentState)
		if err != nil {
			log.Error("getEventCheckRules failed", "err", err)
			return nil
		}
		return &blacklistValidator{
			blacks: blacks,
			rules:  rules,
		}
	}
	return nil
}

func (d *Dpos) getEventCheckRules(header *types.Header, parentState *state.StateDB) (map[common.Hash]*EventCheckRule, error) {
	defer func(start time.Time) {
		getRulesTimer.UpdateSince(start)
	}(time.Now())

	if v, ok := d.eventCheckRules.Get(header.ParentHash); ok {
		return v.(map[common.Hash]*EventCheckRule), nil
	}

	d.rulesLock.Lock()
	defer d.rulesLock.Unlock()
	if v, ok := d.eventCheckRules.Get(header.ParentHash); ok {
		return v.(map[common.Hash]*EventCheckRule), nil
	}

	// if the last updates is long ago, we don't need to get blacklist from the contract.
	num := header.Number.Uint64()
	lastUpdated := lastRulesUpdatedNumber(parentState)
	if num >= 2 && num > lastUpdated+1 {
		parent := d.chain.GetHeader(header.ParentHash, num-1)
		if parent != nil {
			if v, ok := d.eventCheckRules.Get(parent.ParentHash); ok {
				m := v.(map[common.Hash]*EventCheckRule)
				d.eventCheckRules.Add(header.ParentHash, m)
				return m, nil
			}
		} else {
			log.Error("Unexpected error when getEventCheckRules, can not get parent from chain", "number", num, "blockHash", header.Hash(), "parentHash", header.ParentHash)
		}
	}

	// can't get blacklist from cache, try to call the contract
	alABI := d.abi[systemcontract.AddressListContractName]
	method := "getRuleByIndex"
	get := func(i uint32) (common.Hash, int, common.AddressCheckType, error) {
		ret, err := d.commonCallContract(header, parentState, alABI, systemcontract.AddressListContractAddr, method, 3, i)
		if err != nil {
			return common.Hash{}, 0, common.CheckNone, err
		}
		sig := ret[0].([32]byte)
		idx := ret[1].(*big.Int).Uint64()
		ct := ret[2].(uint8)

		return sig, int(idx), common.AddressCheckType(ct), nil
	}

	cnt, err := d.getEventCheckRulesLen(header, parentState)
	if err != nil {
		log.Error("getEventCheckRulesLen failed", "err", err)
		return nil, err
	}
	rules := make(map[common.Hash]*EventCheckRule)
	for i := 0; i < cnt; i++ {
		sig, idx, ct, err := get(uint32(i))
		if err != nil {
			log.Error("getRuleByIndex failed", "index", i, "number", num, "blockHash", header.Hash(), "err", err)
			return nil, err
		}
		rule, exist := rules[sig]
		if !exist {
			rule = &EventCheckRule{
				EventSig: sig,
				Checks:   make(map[int]common.AddressCheckType),
			}
			rules[sig] = rule
		}
		rule.Checks[idx] = ct
	}

	d.eventCheckRules.Add(header.ParentHash, rules)
	return rules, nil
}

func (d *Dpos) getEventCheckRulesLen(header *types.Header, parentState *state.StateDB) (int, error) {
	ret, err := d.commonCallContract(header, parentState, d.abi[systemcontract.AddressListContractName], systemcontract.AddressListContractAddr, "rulesLen", 1)
	if err != nil {
		return 0, err
	}
	ln, ok := ret[0].(uint32)
	if !ok {
		return 0, fmt.Errorf("unexpected output type, value: %v", ret[0])
	}
	return int(ln), nil
}

func (d *Dpos) commonCallContract(header *types.Header, statedb *state.StateDB, contractABI abi.ABI, addr common.Address, method string, expectResultLen int, args ...interface{}) ([]interface{}, error) {
	data, err := contractABI.Pack(method, args...)
	if err != nil {
		log.Error("Can't pack data ", "method", method, "err", err)
		return nil, err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &addr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// Note: It's safe to use minimalChainContext for executing AddressListContract
	result, err := vmcaller.ExecuteMsg(msg, statedb, header, newMinimalChainContext(d), d.chainConfig)
	if err != nil {
		return nil, err
	}

	// unpack data
	ret, err := contractABI.Unpack(method, result)
	if err != nil {
		return nil, err
	}
	if len(ret) != expectResultLen {
		return nil, errors.New("invalid result length")
	}
	return ret, nil
}

// Since the state variables are as follows:
//    bool public initialized;
//    bool public enabled;
//    address public admin;
//    address public pendingAdmin;
//    mapping(address => bool) private devs;
//
// according to [Layout of State Variables in Storage](https://docs.soliditylang.org/en/v0.8.4/internals/layout_in_storage.html),
// and after optimizer enabled, the `initialized`, `enabled` and `admin` will be packed, and stores at slot 0,
// `pendingAdmin` stores at slot 1, and the position for `devs` is 2.
func isDeveloperVerificationEnabled(state consensus.StateReader) bool {
	compactValue := state.GetState(systemcontract.AddressListContractAddr, common.Hash{})
	// Layout of slot 0:
	// [0   -    9][10-29][  30   ][    31     ]
	// [zero bytes][admin][enabled][initialized]
	enabledByte := compactValue.Bytes()[common.HashLength-2]
	return enabledByte == 0x01
}

func calcSlotOfDevMappingKey(addr common.Address) common.Hash {
	p := make([]byte, common.HashLength)
	binary.BigEndian.PutUint16(p[common.HashLength-2:], uint16(systemcontract.DevMappingPosition))
	return crypto.Keccak256Hash(addr.Hash().Bytes(), p)
}

func lastBlacklistUpdatedNumber(state consensus.StateReader) uint64 {
	value := state.GetState(systemcontract.AddressListContractAddr, systemcontract.BlackLastUpdatedNumberPosition)
	return value.Big().Uint64()
}

func lastRulesUpdatedNumber(state consensus.StateReader) uint64 {
	value := state.GetState(systemcontract.AddressListContractAddr, systemcontract.RulesLastUpdatedNumberPosition)
	return value.Big().Uint64()
}
