package bor

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const maxSpanListLimit = 150 // a span is ~6 KB => we can fit 150 spans in 1 MB response
const blockAuthorsCollisionCheck = 10
const blockProducerMaxSpanLookback = 50

var (
	LastSpanIDKey         = []byte{0x35} // Key to store last span start block
	SpanPrefixKey         = []byte{0x36} // prefix key to store span
	LastProcessedEthBlock = []byte{0x38} // key to store last processed eth block for seed
	SeedLastProducerKey   = []byte{0x39} // key to store last producer of the span
)

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	sk  staking.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespace
	codespace sdk.CodespaceType
	// param space
	paramSpace subspace.Subspace
	// contract caller
	contractCaller helper.IContractCaller
	// chain manager keeper
	chainKeeper chainmanager.Keeper
}

// NewKeeper is the constructor of Keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	chainKeeper chainmanager.Keeper,
	stakingKeeper staking.Keeper,
	caller *helper.ContractCaller,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramSpace:     paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:      codespace,
		chainKeeper:    chainKeeper,
		sk:             stakingKeeper,
		contractCaller: caller,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetSpanKey appends prefix to start block
func GetSpanKey(id uint64) []byte {
	return append(SpanPrefixKey, []byte(strconv.FormatUint(id, 10))...)
}

func GetLastSeedProducer(id uint64) []byte {
	idBytes := sdk.Uint64ToBigEndian(id)
	return append(SeedLastProducerKey, idBytes...)
}

// SetContractCaller sets contract caller
func (k *Keeper) SetContractCaller(contractCaller helper.IContractCaller) {
	k.contractCaller = contractCaller
}

// AddNewSpan adds new span for bor to store
func (k *Keeper) AddNewSpan(ctx sdk.Context, span hmTypes.Span) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(span)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling span", "error", err)
		return err
	}

	// store set span id
	store.Set(GetSpanKey(span.ID), out)

	// update last span
	k.UpdateLastSpan(ctx, span.ID)

	return nil
}

// AddNewRawSpan adds new span for bor to store
func (k *Keeper) AddNewRawSpan(ctx sdk.Context, span hmTypes.Span) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(span)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling span", "error", err)
		return err
	}

	store.Set(GetSpanKey(span.ID), out)

	return nil
}

// GetSpan fetches span indexed by id from store
func (k *Keeper) GetSpan(ctx sdk.Context, id uint64) (*hmTypes.Span, error) {
	store := ctx.KVStore(k.storeKey)
	spanKey := GetSpanKey(id)

	// If we are starting from 0 there will be no spanKey present
	if !store.Has(spanKey) {
		return nil, errors.New("span not found for id")
	}

	var span hmTypes.Span
	if err := k.cdc.UnmarshalBinaryBare(store.Get(spanKey), &span); err != nil {
		return nil, err
	}

	return &span, nil
}

func (k *Keeper) HasSpan(ctx sdk.Context, id uint64) bool {
	store := ctx.KVStore(k.storeKey)
	spanKey := GetSpanKey(id)

	return store.Has(spanKey)
}

// GetAllSpans fetches all indexed by id from store
func (k *Keeper) GetAllSpans(ctx sdk.Context) (spans []*hmTypes.Span) {
	// iterate through spans and create span update array
	k.IterateSpansAndApplyFn(ctx, func(span hmTypes.Span) error {
		// append to list of validatorUpdates
		spans = append(spans, &span)
		return nil
	})

	return
}

// GetSpanList returns all spans with params like page and limit
func (k *Keeper) GetSpanList(ctx sdk.Context, page uint64, limit uint64) ([]hmTypes.Span, error) {
	store := ctx.KVStore(k.storeKey)

	// have max limit
	if limit > maxSpanListLimit {
		limit = maxSpanListLimit
	}

	// get paginated iterator
	iterator := hmTypes.KVStorePrefixIteratorPaginated(store, SpanPrefixKey, uint(page), uint(limit))

	// loop through validators to get valid validators
	var spans []hmTypes.Span

	for ; iterator.Valid(); iterator.Next() {
		var span hmTypes.Span
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &span); err == nil {
			spans = append(spans, span)
		}
	}

	return spans, nil
}

// GetLastSpan fetches last span using lastStartBlock
func (k *Keeper) GetLastSpan(ctx sdk.Context) (*hmTypes.Span, error) {
	store := ctx.KVStore(k.storeKey)

	var lastSpanID uint64

	if store.Has(LastSpanIDKey) {
		// get last span id
		var err error
		if lastSpanID, err = strconv.ParseUint(string(store.Get(LastSpanIDKey)), 10, 64); err != nil {
			return nil, err
		}
	}

	return k.GetSpan(ctx, lastSpanID)
}

// FreezeSet freezes validator set for next span
func (k *Keeper) FreezeSet(ctx sdk.Context, id uint64, startBlock uint64, endBlock uint64, borChainID string, seed common.Hash) error {
	var (
		newProducers []hmTypes.Validator
		err          error
	)

	if ctx.BlockHeight() < helper.GetJorvikHeight() {
		newProducers, err = k.SelectNextProducers(ctx, seed, nil)
		if err != nil {
			return err
		}

		// increment last eth block
		k.IncrementLastEthBlock(ctx)
	} else {
		var lastSpan *hmTypes.Span
		var lastSpanId uint64
		if id < 2 {
			lastSpanId = id - 1
		} else {
			lastSpanId = id - 2
		}

		lastSpan, err = k.GetSpan(ctx, lastSpanId)
		if err != nil {
			return err
		}

		prevVals := make([]hmTypes.Validator, 0, len(lastSpan.ValidatorSet.Validators))
		for _, val := range lastSpan.ValidatorSet.Validators {
			prevVals = append(prevVals, *val)
		}

		// select next producers
		newProducers, err = k.SelectNextProducers(ctx, seed, prevVals)
		if err != nil {
			return err
		}
	}

	// generate new span
	newSpan := hmTypes.NewSpan(
		id,
		startBlock,
		endBlock,
		k.sk.GetValidatorSet(ctx),
		newProducers,
		borChainID,
	)

	k.Logger(ctx).Info("Freezing new span", "id", id, "span", newSpan)

	return k.AddNewSpan(ctx, newSpan)
}

// SelectNextProducers selects producers for next span
func (k *Keeper) SelectNextProducers(ctx sdk.Context, seed common.Hash, prevVals []hmTypes.Validator) (vals []hmTypes.Validator, err error) {
	if ctx.BlockHeight() < helper.GetJorvikHeight() {
		prevVals = nil
	}

	// spanEligibleVals are current validators who are not getting deactivated in between next span
	spanEligibleVals := k.sk.GetSpanEligibleValidators(ctx)
	producerCount := k.GetParams(ctx).ProducerCount

	if producerCount > math.MaxInt64 {
		return nil, fmt.Errorf("producer count value out of range for int: %d", producerCount)
	}

	// if producers to be selected is more than current validators no need to select/shuffle
	if len(spanEligibleVals) <= int(producerCount) {
		return spanEligibleVals, nil
	}

	if len(prevVals) > 0 {
		// rollback voting powers for the selection algorithm
		spanEligibleVals = rollbackVotingPowers(ctx, spanEligibleVals, prevVals)
	}

	// TODO remove old selection algorithm
	// select next producers using seed as block header hash
	fn := SelectNextProducers
	if ctx.BlockHeight() < helper.GetNewSelectionAlgoHeight() {
		fn = XXXSelectNextProducers
	}

	newProducersIds, err := fn(seed, spanEligibleVals, producerCount)
	if err != nil {
		return vals, err
	}

	IDToPower := make(map[uint64]uint64)
	for _, ID := range newProducersIds {
		IDToPower[ID] = IDToPower[ID] + 1
	}

	for key, value := range IDToPower {
		if val, ok := k.sk.GetValidatorFromValID(ctx, hmTypes.NewValidatorID(key)); ok {
			if value > math.MaxInt64 {
				return nil, fmt.Errorf("value out of range for int64: %d", value)
			}
			val.VotingPower = int64(value)
			vals = append(vals, val)
		}
	}

	// sort by address
	vals = hmTypes.SortValidatorByAddress(vals)

	return vals, nil
}

// UpdateLastSpan updates the last span start block
func (k *Keeper) UpdateLastSpan(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(LastSpanIDKey, []byte(strconv.FormatUint(id, 10)))
}

// IncrementLastEthBlock increment last eth block
func (k *Keeper) IncrementLastEthBlock(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	lastEthBlock := big.NewInt(0)
	if store.Has(LastProcessedEthBlock) {
		lastEthBlock = lastEthBlock.SetBytes(store.Get(LastProcessedEthBlock))
	}

	store.Set(LastProcessedEthBlock, lastEthBlock.Add(lastEthBlock, big.NewInt(1)).Bytes())
}

// SetLastEthBlock sets last eth block number
func (k *Keeper) SetLastEthBlock(ctx sdk.Context, blockNumber *big.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(LastProcessedEthBlock, blockNumber.Bytes())
}

// GetLastEthBlock get last processed Eth block for seed
func (k *Keeper) GetLastEthBlock(ctx sdk.Context) *big.Int {
	store := ctx.KVStore(k.storeKey)

	lastEthBlock := big.NewInt(0)
	if store.Has(LastProcessedEthBlock) {
		lastEthBlock = lastEthBlock.SetBytes(store.Get(LastProcessedEthBlock))
	}

	return lastEthBlock
}

func (k *Keeper) GetNextSpanSeed(ctx sdk.Context, id uint64) (common.Hash, common.Address, error) {
	var (
		blockHeader *ethTypes.Header
		borBlock    uint64
		seedSpan    *hmTypes.Span
		err         error
		author      *common.Address
	)

	if ctx.BlockHeight() < helper.GetJorvikHeight() {
		lastEthBlock := k.GetLastEthBlock(ctx)
		// increment last processed header block number
		newEthBlock := lastEthBlock.Add(lastEthBlock, big.NewInt(1))
		k.Logger(ctx).Debug("newEthBlock to generate seed", "newEthBlock", newEthBlock)

		// fetch block header from mainchain
		var e error
		blockHeader, e = k.contractCaller.GetMainChainBlock(newEthBlock)
		if e != nil {
			k.Logger(ctx).Error("Error fetching block header from mainchain while calculating next span seed", "error", e)
			return common.Hash{}, common.Address{}, e
		}
		author = &common.Address{}
	} else {
		var seedSpanID uint64
		if id < 2 {
			seedSpanID = id - 1
		} else {
			seedSpanID = id - 2
		}
		seedSpan, err = k.GetSpan(ctx, seedSpanID)
		if err != nil {
			k.Logger(ctx).Error("Error fetching span while calculating next span seed", "error", err)
			return common.Hash{}, common.Address{}, err
		}

		borBlock, author, err = k.getBorBlockForSpanSeed(ctx, seedSpan, id)
		if err != nil {
			return common.Hash{}, common.Address{}, err
		}

		if borBlock > math.MaxInt64 {
			return common.Hash{}, common.Address{}, fmt.Errorf("bor block value out of range for int64: %d", borBlock)
		}

		blockHeader, err = k.contractCaller.GetMaticChainBlock(big.NewInt(int64(borBlock)))
		if err != nil {
			k.Logger(ctx).Error("Error fetching block header from bor chain while calculating next span seed", "error", err, "block", borBlock)
			return common.Hash{}, common.Address{}, err
		}

		if author == nil {
			k.Logger(ctx).Error("seed author is nil")
			return blockHeader.Hash(), common.Address{}, fmt.Errorf("seed author is nil")
		}

		k.Logger(ctx).Debug("fetched block for seed", "block", borBlock, "author", author, "span id", id)
	}

	return blockHeader.Hash(), *author, nil
}

// StoreSeedProducer stores producer of the block used for seed for the given span id
func (k *Keeper) StoreSeedProducer(ctx sdk.Context, id uint64, producer *common.Address) error {
	store := ctx.KVStore(k.storeKey)
	lastSeedKey := GetLastSeedProducer(id)

	if store.Has(lastSeedKey) {
		return errors.New("seed producer already stored")
	}

	store.Set(lastSeedKey, producer.Bytes())
	return nil
}

// GetSeedProducer gets producer of the block used for seed for the given span id
func (k *Keeper) GetSeedProducer(ctx sdk.Context, id uint64) (*common.Address, error) {
	store := ctx.KVStore(k.storeKey)
	lastSeedKey := GetLastSeedProducer(id)

	authorBytes := store.Get(lastSeedKey)
	if authorBytes == nil {
		return nil, nil //nolint: nilnil
	}

	author := common.BytesToAddress(authorBytes)

	return &author, nil
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the bor module's parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the bor module's parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

//
// Utils
//

// IterateSpansAndApplyFn iterates spans and apply the given function.
func (k *Keeper) IterateSpansAndApplyFn(ctx sdk.Context, f func(span hmTypes.Span) error) {
	store := ctx.KVStore(k.storeKey)

	// get span iterator
	iterator := sdk.KVStorePrefixIterator(store, SpanPrefixKey)
	defer iterator.Close()

	// loop through spans to get valid spans
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall span
		var result hmTypes.Span
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &result); err != nil {
			k.Logger(ctx).Error("Error UnmarshalBinaryBare", "error", err)
		}
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}

// getBorBlockForSpanSeed returns the bor block number and its producer whose hash is used as seed for the next span
func (k *Keeper) getBorBlockForSpanSeed(ctx sdk.Context, seedSpan *hmTypes.Span, proposedSpanID uint64) (uint64, *common.Address, error) {
	var (
		borBlock uint64
		author   *common.Address
		err      error
	)

	logger := k.Logger(ctx)

	logger.Debug("getting bor block for span seed", "span id", seedSpan.ID, "proposed span id", proposedSpanID)

	if proposedSpanID == 1 {
		borBlock = 1
		author, err = k.contractCaller.GetBorChainBlockAuthor(big.NewInt(int64(borBlock)))
		if err != nil {
			logger.Error("Error fetching first block for span seed", "error", err, "block", borBlock)
			return 0, nil, err
		}

		logger.Debug("returning first block author", "author", author, "block", borBlock)

		return borBlock, author, nil
	}

	uniqueAuthors := make(map[string]struct{})
	spanID := proposedSpanID - 1

	lastAuthor, err := k.GetSeedProducer(ctx, spanID)
	if err != nil {
		logger.Error("Error fetching last seed producer", "error", err, "span id", spanID)
		return 0, nil, err
	}

	// get seed block authors from last "blockProducerMaxSpanLookback" spans
	for i := 0; len(uniqueAuthors) < blockAuthorsCollisionCheck && i < blockProducerMaxSpanLookback; i++ {
		if spanID == 0 {
			break
		}

		author, err = k.GetSeedProducer(ctx, spanID)
		if err != nil {
			logger.Error("Error fetching span seed producer", "error", err, "span id", spanID)
			return 0, nil, err
		}

		if author == nil {
			logger.Info("GetSeedProducer returned empty value", "span id", spanID)
			break
		}

		uniqueAuthors[author.Hex()] = struct{}{}
		spanID--
	}

	logger.Debug("last authors", "authors", fmt.Sprintf("%+v", uniqueAuthors), "span id", seedSpan.ID)

	firstDiffFromLast := uint64(0)

	// try to find a seed block with an author not in the last "blockAuthorsCollisionCheck" spans
	borParams := k.GetParams(ctx)
	for borBlock = seedSpan.EndBlock; borBlock >= seedSpan.StartBlock; borBlock -= borParams.SprintDuration {
		if borBlock > math.MaxInt64 {
			return 0, nil, fmt.Errorf("bor block value out of range for int64: %d", borBlock)
		}
		author, err = k.contractCaller.GetBorChainBlockAuthor(big.NewInt(int64(borBlock)))
		if err != nil {
			logger.Error("Error fetching block author from bor chain while calculating next span seed", "error", err, "block", borBlock)
			return 0, nil, err
		}

		if _, exists := uniqueAuthors[author.Hex()]; !exists || len(seedSpan.ValidatorSet.Validators) == 1 {
			logger.Debug("got author", "author", author, "block", borBlock)
			return borBlock, author, nil
		}

		if firstDiffFromLast == 0 && lastAuthor != nil && author.Hex() != lastAuthor.Hex() {
			firstDiffFromLast = borBlock
		}
	}

	// if no unique author found, return the first different block author
	borBlock = firstDiffFromLast
	if borBlock == 0 {
		borBlock = seedSpan.EndBlock
	}

	if borBlock > math.MaxInt64 {
		return 0, nil, fmt.Errorf("bor block value out of range for int64: %d", borBlock)
	}

	author, err = k.contractCaller.GetBorChainBlockAuthor(big.NewInt(int64(borBlock)))
	if err != nil {
		logger.Error("Error fetching end block author from bor chain while calculating next span seed", "error", err, "block", borBlock)
		return 0, nil, err
	}

	logger.Debug("returning first different block author", "author", author, "block", borBlock)

	return borBlock, author, nil
}

// rollbackVotingPowers rolls back voting powers of validators from a previous snapshot of validators
func rollbackVotingPowers(ctx sdk.Context, valsNew, valsOld []hmTypes.Validator) []hmTypes.Validator {
	idToVP := make(map[uint64]int64)
	for _, val := range valsOld {
		idToVP[val.ID.Uint64()] = val.VotingPower
	}

	for i := range valsNew {
		if _, ok := idToVP[valsNew[i].ID.Uint64()]; ok {
			valsNew[i].VotingPower = idToVP[valsNew[i].ID.Uint64()]
		} else {
			valsNew[i].VotingPower = 0
		}
	}

	return valsNew
}
