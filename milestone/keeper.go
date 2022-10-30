package milestone

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/chainmanager"
	cmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/milestone/types"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	DefaultValue          = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag
	MilestoneKey          = []byte{0x20} // Key to store milestone
	CountKey              = []byte{0x30} //Key to store the count
	MilestoneNoAckKey     = []byte{0x21} //Key to store the NoAckMilestone
	MilestoneLastNoAckKey = []byte{0x22} //Key to store the Latest NoAckMilestone
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

// AddMilestone adds milestone into final blocks
func (k *Keeper) AddMilestone(ctx sdk.Context, milestone hmTypes.Milestone) error {

	milestoneNumber := k.GetCount(ctx) + 1 //GetCount gives the number of previous milestone

	key := GetMilestoneKey(milestoneNumber)
	if err := k.addMilestone(ctx, key, milestone); err != nil {
		return err
	}

	pruningNumber := milestoneNumber - 100

	k.PruneMilestone(ctx, pruningNumber) //Prune the old milestone to reduce the memory consumption
	k.SetCount(ctx, milestoneNumber)
	k.Logger(ctx).Info("Adding good milestone to state", "milestone", milestone, "milestoneNumber", milestoneNumber)

	return nil
}

// addMilestone adds milestone to store
func (k *Keeper) addMilestone(ctx sdk.Context, key []byte, milestone hmTypes.Milestone) error {
	store := ctx.KVStore(k.storeKey)

	// create Checkpoint block and marshall
	out, err := k.cdc.MarshalBinaryBare(milestone)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling milestone", "error", err)
		return err
	}

	// store in key provided
	store.Set(key, out)

	return nil
}

// GetMilestoneKey appends prefix to milestoneNumber
func GetMilestoneKey(milestoneNumber uint64) []byte {
	milestoneNumberBytes := []byte(strconv.FormatUint(milestoneNumber, 10))
	return append(MilestoneKey, milestoneNumberBytes...)
}

// GetMilestoneByNumber to get milestone by milestone number
func (k *Keeper) GetMilestoneByNumber(ctx sdk.Context, number uint64) (*hmTypes.Milestone, error) {
	store := ctx.KVStore(k.storeKey)
	milestoneKey := GetMilestoneKey(number)

	var _milestone hmTypes.Milestone

	if store.Has(milestoneKey) {
		if err := k.cdc.UnmarshalBinaryBare(store.Get(milestoneKey), &_milestone); err != nil {
			return nil, err
		}

		return &_milestone, nil
	}

	return nil, errors.New("Invalid milestone Index")
}

// GetLastMilestone gets last milestone, milestone number = GetCount()
func (k *Keeper) GetLastMilestone(ctx sdk.Context) (*hmTypes.Milestone, error) {
	store := ctx.KVStore(k.storeKey)
	Count := k.GetCount(ctx)

	lastMilestoneKey := GetMilestoneKey(Count)

	// fetch milestone and unmarshall
	var _milestone hmTypes.Milestone

	if store.Has(lastMilestoneKey) {
		err := k.cdc.UnmarshalBinaryBare(store.Get(lastMilestoneKey), &_milestone)
		if err != nil {
			k.Logger(ctx).Error("Unable to fetch last milestone from store", "number", Count)
			return nil, err
		} else {
			return &_milestone, nil
		}
	}

	return nil, cmn.ErrNoMilestoneFound(k.Codespace())
}

// HasStoreValue check if value exists in store or not
func (k *Keeper) HasStoreValue(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
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

// SetCount set the count number
func (k *Keeper) SetCount(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(number, 10))
	// set no-ack
	store.Set(CountKey, value)
}

// GetCount returns count
func (k *Keeper) GetCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	// check if count is there
	if store.Has(CountKey) {
		// get current count
		result, err := strconv.ParseUint(string(store.Get(CountKey)), 10, 64)
		if err == nil {
			return result
		}
	}

	return uint64(0)
}

// FlushCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) PruneMilestone(ctx sdk.Context, number uint64) {

	store := ctx.KVStore(k.storeKey)
	if number <= 0 {
		return
	}

	milestoneKey := GetMilestoneKey(number)

	if !store.Has(milestoneKey) {
		return
	}

	store.Delete(milestoneKey)
}

// SetLastNoAck set last no-ack object
func (k *Keeper) SetNoAckMilestone(ctx sdk.Context, milestone hmTypes.Milestone) error {
	store := ctx.KVStore(k.storeKey)

	//milestoneNoAckKey := GetMilestoneNoAckKey(milestoneId)
	//value := []byte(milestoneId)

	// create Checkpoint block and marshall
	out, err := k.cdc.MarshalBinaryBare(milestone)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling milestone", "error", err)
		return err
	}

	// store in key provided
	store.Set(MilestoneLastNoAckKey, out)
	return nil
	// set no-ack-milestone
	//store.Set(milestoneNoAckKey, value)
	//store.Set(MilestoneLastNoAckKey, value)
}

// GetLastNoAckMilestone returns last no ack milestone
func (k *Keeper) GetLastNoAckMilestone(ctx sdk.Context) (*hmTypes.Milestone, error) {
	store := ctx.KVStore(k.storeKey)
	// check if ack count is there
	// if store.Has(MilestoneLastNoAckKey) {
	// 	// get current ACK count
	// 	result := string(store.Get(MilestoneLastNoAckKey))
	// 	return result
	// }

	var _milestone hmTypes.Milestone

	if store.Has(MilestoneLastNoAckKey) {
		err := k.cdc.UnmarshalBinaryBare(store.Get(MilestoneLastNoAckKey), &_milestone)
		if err != nil {
			k.Logger(ctx).Error("Unable to fetch last no ack milestone from store")
			return nil, err
		} else {
			return &_milestone, nil
		}
	}

	return nil, cmn.ErrNoMilestoneFound(k.Codespace())
}

// GetLastNoAckMilestone returns last no ack milestone
func (k *Keeper) GetNoAckMilestone(ctx sdk.Context, milestoneId string) bool {
	store := ctx.KVStore(k.storeKey)
	// check if No Ack Milestone is there
	if store.Has(GetMilestoneNoAckKey(milestoneId)) {
		return true
	}
	return false
}

// GetMilestoneKey appends prefix to milestoneNumber
func GetMilestoneNoAckKey(milestoneId string) []byte {
	milestoneNoAckBytes := []byte(milestoneId)
	return append(MilestoneNoAckKey, milestoneNoAckBytes...)
}
