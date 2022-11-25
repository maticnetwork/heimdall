package checkpoint

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	cmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ACKCountKey         = []byte{0x11} // key to store ACK count
	BufferCheckpointKey = []byte{0x12} // Key to store checkpoint in buffer
	CheckpointKey       = []byte{0x13} // prefix key for when storing checkpoint after ACK
	LastNoACKKey        = []byte{0x14} // key to store last no-ack

	//#########Milestone Keys#################

	MilestoneKey          = []byte{0x20} // Key to store milestone
	CountKey              = []byte{0x30} //Key to store the count
	MilestoneNoAckKey     = []byte{0x40} //Key to store the NoAckMilestone
	MilestoneLastNoAckKey = []byte{0x50} //Key to store the Latest NoAckMilestone
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

// AddCheckpoint adds checkpoint into final blocks
func (k *Keeper) AddCheckpoint(ctx sdk.Context, checkpointNumber uint64, checkpoint hmTypes.Checkpoint) error {
	key := GetCheckpointKey(checkpointNumber)
	if err := k.addCheckpoint(ctx, key, checkpoint); err != nil {
		return err
	}

	k.Logger(ctx).Info("Adding good checkpoint to state", "checkpoint", checkpoint, "checkpointNumber", checkpointNumber)

	return nil
}

// SetCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) SetCheckpointBuffer(ctx sdk.Context, checkpoint hmTypes.Checkpoint) error {
	err := k.addCheckpoint(ctx, BufferCheckpointKey, checkpoint)
	if err != nil {
		return err
	}

	return nil
}

// addCheckpoint adds checkpoint to store
func (k *Keeper) addCheckpoint(ctx sdk.Context, key []byte, checkpoint hmTypes.Checkpoint) error {
	store := ctx.KVStore(k.storeKey)

	// create Checkpoint block and marshall
	out, err := k.cdc.MarshalBinaryBare(checkpoint)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling checkpoint", "error", err)
		return err
	}

	// store in key provided
	store.Set(key, out)

	return nil
}

// GetCheckpointByNumber to get checkpoint by checkpoint number
func (k *Keeper) GetCheckpointByNumber(ctx sdk.Context, number uint64) (hmTypes.Checkpoint, error) {
	store := ctx.KVStore(k.storeKey)
	checkpointKey := GetCheckpointKey(number)

	var _checkpoint hmTypes.Checkpoint

	if store.Has(checkpointKey) {
		if err := k.cdc.UnmarshalBinaryBare(store.Get(checkpointKey), &_checkpoint); err != nil {
			return _checkpoint, err
		}

		return _checkpoint, nil
	}

	return _checkpoint, errors.New("Invalid checkpoint Index")
}

// GetCheckpointList returns all checkpoints with params like page and limit
func (k *Keeper) GetCheckpointList(ctx sdk.Context, page uint64, limit uint64) ([]hmTypes.Checkpoint, error) {
	store := ctx.KVStore(k.storeKey)

	// create headers
	var checkpoints []hmTypes.Checkpoint

	// have max limit
	if limit > 20 {
		limit = 20
	}

	// get paginated iterator
	iterator := hmTypes.KVStorePrefixIteratorPaginated(store, CheckpointKey, uint(page), uint(limit))

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		var checkpoint hmTypes.Checkpoint
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &checkpoint); err == nil {
			checkpoints = append(checkpoints, checkpoint)
		}
	}

	return checkpoints, nil
}

// GetLastCheckpoint gets last checkpoint, checkpoint number = TotalACKs
func (k *Keeper) GetLastCheckpoint(ctx sdk.Context) (hmTypes.Checkpoint, error) {
	store := ctx.KVStore(k.storeKey)
	acksCount := k.GetACKCount(ctx)

	lastCheckpointKey := acksCount

	// fetch checkpoint and unmarshall
	var _checkpoint hmTypes.Checkpoint

	// no checkpoint received
	// header key
	headerKey := GetCheckpointKey(lastCheckpointKey)
	if store.Has(headerKey) {
		err := k.cdc.UnmarshalBinaryBare(store.Get(headerKey), &_checkpoint)
		if err != nil {
			k.Logger(ctx).Error("Unable to fetch last checkpoint from store", "key", lastCheckpointKey, "acksCount", acksCount)
			return _checkpoint, err
		} else {
			return _checkpoint, nil
		}
	}

	return _checkpoint, cmn.ErrNoCheckpointFound(k.Codespace())
}

// GetCheckpointKey appends prefix to checkpointNumber
func GetCheckpointKey(checkpointNumber uint64) []byte {
	checkpointNumberBytes := []byte(strconv.FormatUint(checkpointNumber, 10))
	return append(CheckpointKey, checkpointNumberBytes...)
}

// HasStoreValue check if value exists in store or not
func (k *Keeper) HasStoreValue(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

// FlushCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(BufferCheckpointKey)
}

// GetCheckpointFromBuffer gets checkpoint in buffer
func (k *Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (*hmTypes.Checkpoint, error) {
	store := ctx.KVStore(k.storeKey)

	// checkpoint block header
	var checkpoint hmTypes.Checkpoint

	if store.Has(BufferCheckpointKey) {
		// Get checkpoint and unmarshall
		err := k.cdc.UnmarshalBinaryBare(store.Get(BufferCheckpointKey), &checkpoint)
		return &checkpoint, err
	}

	return nil, errors.New("No checkpoint found in buffer")
}

// SetLastNoAck set last no-ack object
func (k *Keeper) SetLastNoAck(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(timestamp, 10))
	// set no-ack
	store.Set(LastNoACKKey, value)
}

// GetLastNoAck returns last no ack
func (k *Keeper) GetLastNoAck(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	// check if ack count is there
	if store.Has(LastNoACKKey) {
		// get current ACK count
		result, err := strconv.ParseUint(string(store.Get(LastNoACKKey)), 10, 64)
		if err == nil {
			return result
		}
	}

	return 0
}

// GetCheckpoints get checkpoint all checkpoints
func (k *Keeper) GetCheckpoints(ctx sdk.Context) []hmTypes.Checkpoint {
	store := ctx.KVStore(k.storeKey)
	// get checkpoint header iterator
	iterator := sdk.KVStorePrefixIterator(store, CheckpointKey)
	defer iterator.Close()

	// create headers
	var headers []hmTypes.Checkpoint

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		var checkpoint hmTypes.Checkpoint
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &checkpoint); err == nil {
			headers = append(headers, checkpoint)
		}
	}

	return headers
}

//
// Ack count
//

// GetACKCount returns current ACK count
func (k Keeper) GetACKCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	// check if ack count is there
	if store.Has(ACKCountKey) {
		// get current ACK count
		ackCount, err := strconv.ParseUint(string(store.Get(ACKCountKey)), 10, 64)
		if err != nil {
			k.Logger(ctx).Error("Unable to convert key to int")
		} else {
			return ackCount
		}
	}

	return 0
}

// UpdateACKCountWithValue updates ACK with value
func (k Keeper) UpdateACKCountWithValue(ctx sdk.Context, value uint64) {
	store := ctx.KVStore(k.storeKey)

	// convert
	ackCount := []byte(strconv.FormatUint(value, 10))

	// update
	store.Set(ACKCountKey, ackCount)
}

// UpdateACKCount updates ACK count by 1
func (k Keeper) UpdateACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	// get current ACK Count
	ACKCount := k.GetACKCount(ctx)

	// increment by 1
	ACKs := []byte(strconv.FormatUint(ACKCount+1, 10))

	// update
	store.Set(ACKCountKey, ACKs)
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// ///////////Milestone Functions////////////////////
// AddMilestone adds milestone into final blocks
func (k *Keeper) AddMilestone(ctx sdk.Context, milestone hmTypes.Milestone) error {

	milestoneNumber := k.GetMilestoneCount(ctx) + 1 //GetCount gives the number of previous milestone

	key := GetMilestoneKey(milestoneNumber)
	if err := k.addMilestone(ctx, key, milestone); err != nil {
		return err
	}

	pruningNumber := milestoneNumber - 100

	k.PruneMilestone(ctx, pruningNumber) //Prune the old milestone to reduce the memory consumption
	k.SetMilestoneCount(ctx, milestoneNumber)
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
	Count := k.GetMilestoneCount(ctx)

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

// SetCount set the count number
func (k *Keeper) SetMilestoneCount(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(number, 10))
	// set no-ack
	store.Set(CountKey, value)
}

// GetCount returns count
func (k *Keeper) GetMilestoneCount(ctx sdk.Context) uint64 {
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
func (k *Keeper) SetNoAckMilestone(ctx sdk.Context, milestoneId string) {
	store := ctx.KVStore(k.storeKey)

	milestoneNoAckKey := GetMilestoneNoAckKey(milestoneId)
	value := []byte(milestoneId)

	// set no-ack-milestone
	store.Set(milestoneNoAckKey, value)
	store.Set(MilestoneLastNoAckKey, value)
}

// GetLastNoAckMilestone returns last no ack milestone
func (k *Keeper) GetLastNoAckMilestone(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	// check if ack count is there
	if store.Has(MilestoneLastNoAckKey) {
		// get current ACK count
		result := string(store.Get(MilestoneLastNoAckKey))
		return result
	}
	return ""
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
