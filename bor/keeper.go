package bor

import (
	"errors"
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

var (
	LastSpanIDKey         = []byte{0x35} // Key to store last span start block
	SpanPrefixKey         = []byte{0x36} // prefix key to store span
	LastProcessedEthBlock = []byte{0x38} // key to store last processed eth block for seed
	SpanLastProducerKey   = []byte{0x39} // key to store last producer of the span
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
	contractCaller helper.ContractCaller
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
	caller helper.ContractCaller,
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

func GetSpanLastProducerKey(id uint64) []byte {
	idBytes := sdk.Uint64ToBigEndian(id)
	return append(SpanLastProducerKey, idBytes...)
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
		proposer     hmTypes.Validator
		err          error
	)

	if ctx.BlockHeight() < helper.GetNeedANameHeight() {
		newProducers, err = k.SelectNextProducers(ctx, seed, nil)
		if err != nil {
			return err
		}

		// increment last eth block
		k.IncrementLastEthBlock(ctx)
	} else {
		// fetch span(id - 2)
		var lastSpan *hmTypes.Span
		spanId := id - 2
		if id < 2 {
			spanId = id - 1
		}

		lastSpan, err = k.GetSpan(ctx, spanId)
		if err != nil {
			return err
		}

		vals := make([]hmTypes.Validator, 0, len(lastSpan.ValidatorSet.Validators))
		for _, val := range lastSpan.ValidatorSet.Validators {
			vals = append(vals, *val)
		}

		// select next producers
		newProducers, err = k.SelectNextProducers(ctx, seed, vals)
		if err != nil {
			return err
		}

		// get last producer of the last span
		lastProd, err := k.GetSpanLastProducer(ctx, id-1)
		if err != nil {
			return err
		}

		ctx.Logger().Info("!!!!Fetched last producer for span", "span", id-1, "producer id", lastProd.Signer, "ID", lastProd.ID)

		prods := make([]*hmTypes.Validator, 0, len(newProducers))
		for _, val := range newProducers {
			prods = append(prods, &val)
		}

		borParams := k.GetParams(ctx)
		valSet := hmTypes.NewValidatorSet(prods)
		valSet.IncrementProposerPriority(int(borParams.SprintDuration))
		proposer = *valSet.GetProposer()
		ctx.Logger().Info("!!!!Calculated last producer for span", "span", id, "producer id", proposer.Signer, "ID", proposer.ID)

		if lastProd.ID == proposer.ID {
			// if last producer is the same, then rotate proposer
			valSet.IncrementProposerPriority(1)
			prods = valSet.Validators

			for i, val := range prods {
				newProducers[i] = *val
			}

			proposer = *valSet.GetProposer()
			ctx.Logger().Info("!!!!Rotated last producer for span", "span", id, "producer id", proposer.Signer, "ID", proposer.ID)

		}

		if err = k.StoreSpanLastProducer(ctx, id, proposer); err != nil {
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

	return k.AddNewSpan(ctx, newSpan)
}

// SelectNextProducers selects producers for next span
func (k *Keeper) SelectNextProducers(ctx sdk.Context, seed common.Hash, prevVals []hmTypes.Validator) (vals []hmTypes.Validator, err error) {
	// spanEligibleVals are current validators who are not getting deactivated in between next span
	spanEligibleVals := k.sk.GetSpanEligibleValidators(ctx)
	producerCount := k.GetParams(ctx).ProducerCount

	// if producers to be selected is more than current validators no need to select/shuffle
	if len(spanEligibleVals) <= int(producerCount) {
		return spanEligibleVals, nil
	}

	if prevVals != nil {
		// rollback voting powers for the selection algorithm
		spanEligibleVals = k.rollbackVotingPowers(spanEligibleVals, prevVals)
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

func (k Keeper) GetNextSpanSeed(ctx sdk.Context, id uint64) (common.Hash, error) {
	var (
		blockHeader *ethTypes.Header
		lastSpan    *hmTypes.Span
		err         error
	)

	if ctx.BlockHeader().Height < helper.GetNeedANameHeight() {
		lastEthBlock := k.GetLastEthBlock(ctx)
		// increment last processed header block number
		newEthBlock := lastEthBlock.Add(lastEthBlock, big.NewInt(1))
		k.Logger(ctx).Debug("newEthBlock to generate seed", "newEthBlock", newEthBlock)

		// fetch block header from mainchain
		var err error
		blockHeader, err = k.contractCaller.GetMainChainBlock(newEthBlock)
		if err != nil {
			k.Logger(ctx).Error("Error fetching block header from mainchain while calculating next span seed", "error", err)
			return common.Hash{}, err
		}
	} else {
		spanId := id - 2
		if id < 2 {
			spanId = id - 1
		}
		lastSpan, err = k.GetSpan(ctx, spanId)
		if err != nil {
			k.Logger(ctx).Error("Error fetching span while calculating next span seed", "error", err)
			return common.Hash{}, err
		}

		borBlock := lastSpan.EndBlock
		if spanId == id-1 {
			borBlock = lastSpan.StartBlock
		}

		blockHeader, err = k.contractCaller.GetMaticChainBlock(big.NewInt(int64(borBlock)))
		if err != nil {
			k.Logger(ctx).Error("Error fetching block header from bor chain while calculating next span seed", "error", err, "block", borBlock)
			return common.Hash{}, err
		}
	}

	return blockHeader.Hash(), nil
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

// GetSpanLastProducer gets producer of the last sprint of the given span
func (k *Keeper) GetSpanLastProducer(ctx sdk.Context, id uint64) (hmTypes.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	lastProdKey := GetSpanLastProducerKey(id)

	// At the beginning there will be no lastProdKey present
	// so we need to initialize the store with the last producer of the last span
	if !store.Has(lastProdKey) {
		lastSpan, err := k.GetSpan(ctx, id)
		if err != nil {
			return hmTypes.Validator{}, err
		}

		prods := make([]*hmTypes.Validator, 0, len(lastSpan.SelectedProducers))
		for _, val := range lastSpan.SelectedProducers {
			prods = append(prods, &val)
		}

		borParams := k.GetParams(ctx)
		valSet := hmTypes.NewValidatorSet(prods)
		valSet.IncrementProposerPriority(int(borParams.SprintDuration))
		if err := k.StoreSpanLastProducer(ctx, id, *valSet.GetProposer()); err != nil {
			return hmTypes.Validator{}, err
		}

		ctx.Logger().Info("!!!!Initialized last producer for span", "span", id, "producer id", valSet.GetProposer().Signer, "ID", valSet.GetProposer().ID)
		return *valSet.GetProposer(), nil
	}

	var producer hmTypes.Validator
	if err := k.cdc.UnmarshalBinaryBare(store.Get(lastProdKey), &producer); err != nil {
		return hmTypes.Validator{}, err
	}

	return hmTypes.Validator{}, nil

}

// StoreSpanLastProducer stores producer of the last sprint of the given span
func (k *Keeper) StoreSpanLastProducer(ctx sdk.Context, id uint64, producer hmTypes.Validator) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(producer)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling validator", "error", err)
		return err
	}

	// store set span id
	store.Set(GetSpanLastProducerKey(id), out)
	return nil
}

// rollbackVotingPowers rolls back voting powers of validators from a previous snapshot of validators
func (k *Keeper) rollbackVotingPowers(valsFrom, valsTo []hmTypes.Validator) []hmTypes.Validator {
	idToVP := make(map[uint64]int64)
	for _, val := range valsFrom {
		idToVP[val.ID.Uint64()] = val.VotingPower
	}

	for _, valB := range valsTo {
		if vp, ok := idToVP[valB.ID.Uint64()]; ok {
			idToVP[valB.ID.Uint64()] = vp
		}
	}

	for i := range valsFrom {
		valsFrom[i].VotingPower = idToVP[valsFrom[i].ID.Uint64()]
	}

	return valsFrom
}
