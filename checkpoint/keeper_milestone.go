package checkpoint

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	cmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	MilestoneKey          = []byte{0x20} // Key to store milestone
	CountKey              = []byte{0x30} //Key to store the count
	MilestoneNoAckKey     = []byte{0x40} //Key to store the NoAckMilestone
	MilestoneLastNoAckKey = []byte{0x50} //Key to store the Latest NoAckMilestone
	LastMilestoneTimeout  = []byte{0x60} //Key to store the Last Milestone Timeout
	BlockNumberKey        = []byte{0x70} //Key to store the count
)

// Logger returns a module-specific logger
func (k Keeper) MilestoneLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "Milestone")
}

// AddMilestone adds milestone in the store
func (k *Keeper) AddMilestone(ctx sdk.Context, milestone hmTypes.Milestone) error {
	milestoneNumber := k.GetMilestoneCount(ctx) + 1 //GetCount gives the number of previous milestone

	key := GetMilestoneKey(milestoneNumber)
	if err := k.addMilestone(ctx, key, milestone); err != nil {
		return err
	}

	pruningNumber := milestoneNumber - helper.MilestonePruneNumber

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

	var milestone hmTypes.Milestone

	if store.Has(milestoneKey) {
		if err := k.cdc.UnmarshalBinaryBare(store.Get(milestoneKey), &milestone); err != nil {
			return nil, err
		}

		return &milestone, nil
	}

	return nil, cmn.ErrInvalidMilestoneIndex(k.Codespace())
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

// GetMilestoneCount returns milestone count
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

// SetMilestoneBlockNumber set the block number when the latest milestone enter the handler
func (k *Keeper) SetMilestoneBlockNumber(ctx sdk.Context, number int64) {
	store := ctx.KVStore(k.storeKey)
	// convert block number to bytes
	value := []byte(strconv.FormatInt(number, 10))
	// set
	store.Set(BlockNumberKey, value)
}

// GetMilestoneBlockNumber returns the block number when the latest milestone enter the handler
func (k *Keeper) GetMilestoneBlockNumber(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	// check if block number is there
	if store.Has(BlockNumberKey) {
		// get the block number
		result, err := strconv.ParseInt(string(store.Get(BlockNumberKey)), 10, 64)
		if err == nil {
			return result
		}
	}

	return int64(0)
}

// PruneMilestone remove the milestone from the db based on number
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
	// check if MilestoneLastNoAckKey key exists
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
	return store.Has(GetMilestoneNoAckKey(milestoneId))
}

// SetLastMilestoneTimeout set lastMilestone timeout time
func (k *Keeper) SetLastMilestoneTimeout(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(timestamp, 10))
	// set no-ack
	store.Set(LastMilestoneTimeout, value)
}

// GetLastMilestoneTimeout returns lastMilestone timeout time
func (k *Keeper) GetLastMilestoneTimeout(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	//check if lastMilestoneTimeout key exists
	if store.Has(LastMilestoneTimeout) {
		// get last milestone timeout
		result, err := strconv.ParseUint(string(store.Get(LastMilestoneTimeout)), 10, 64)
		if err == nil {
			return result
		}
	}

	return 0
}

// GetMilestoneKey appends prefix to milestoneNumber
func GetMilestoneNoAckKey(milestoneId string) []byte {
	milestoneNoAckBytes := []byte(milestoneId)
	return append(MilestoneNoAckKey, milestoneNoAckBytes...)
}
