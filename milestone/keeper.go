package milestone

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/milestone/types"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag
	MilestoneKey = []byte{0x20} // Key to store milestone

)

// ModuleCommunicator manages different module interaction
type ModuleCommunicator interface {
	GetAllDividendAccounts(ctx sdk.Context) []hmTypes.DividendAccount
}

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	// staking keeper
	sk staking.Keeper
	ck chainmanager.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespace
	codespace sdk.CodespaceType
	// param space
	paramSpace subspace.Subspace

	// module communicator
	moduleCommunicator ModuleCommunicator
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	stakingKeeper staking.Keeper,
	chainKeeper chainmanager.Keeper,
	moduleCommunicator ModuleCommunicator,
) Keeper {
	keeper := Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		paramSpace:         paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:          codespace,
		sk:                 stakingKeeper,
		ck:                 chainKeeper,
		moduleCommunicator: moduleCommunicator,
	}

	return keeper
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// SetMilestone store the milestone in the store
func (k *Keeper) SetMilestone(ctx sdk.Context, milestone hmTypes.Milestone) error {
	err := k.AddMilestone(ctx, milestone)
	if err != nil {
		return err
	}

	return nil
}

// addMilestone adds milestone to store
func (k *Keeper) AddMilestone(ctx sdk.Context, milestone hmTypes.Milestone) error {
	store := ctx.KVStore(k.storeKey)

	// create milestone block and marshall
	out, err := k.cdc.MarshalBinaryBare(milestone)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling milestone", "error", err)
		return err
	}

	// store in key provided
	store.Set(MilestoneKey, out)

	return nil
}

// HasStoreValue check if value exists in store or not
func (k *Keeper) HasStoreValue(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

// GetMilestone gives the milestone
func (k *Keeper) GetMilestone(ctx sdk.Context) (*hmTypes.Milestone, error) {
	store := ctx.KVStore(k.storeKey)

	// milestone block header
	var milestone hmTypes.Milestone

	if store.Has(MilestoneKey) {
		// Get milestone and unmarshall
		err := k.cdc.UnmarshalBinaryBare(store.Get(MilestoneKey), &milestone)
		return &milestone, err
	}

	return nil, errors.New("No milestone found")
}

// Params

// SetParams sets the milestone module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the milestone module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetCount(ctx sdk.Context) (count uint64) {
	params := types.Params{}
	k.paramSpace.GetParamSet(ctx, &params)

	milestone, err := k.GetMilestone(ctx)
	if err != nil || milestone == nil {
		return 0
	}

	count = (milestone.EndBlock + 1) / params.SprintLength
	return
}
