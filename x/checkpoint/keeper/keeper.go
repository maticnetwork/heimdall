package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	// "github.com/maticnetwork/heimdall/chainmanager"
	// "github.com/maticnetwork/heimdall/staking"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ACKCountKey         = []byte{0x11} // key to store ACK count
	BufferCheckpointKey = []byte{0x12} // Key to store checkpoint in buffer
	CheckpointKey       = []byte{0x13} // prefix key for when storing checkpoint after ACK
	LastNoACKKey        = []byte{0x14} // key to store last no-ack
)

// ModuleCommunicator manages different module interaction
type ModuleCommunicator interface {
	GetAllDividendAccounts(ctx sdk.Context) []hmTypes.DividendAccount
}

type (
	Keeper struct {
		cdc                codec.Marshaler
		storeKey           sdk.StoreKey
		memKey             sdk.StoreKey
		paramSubspace      paramtypes.Subspace
		moduleCommunicator ModuleCommunicator
		//TODO: add staking and chainmanager
		// sk staking.Keeper
		// ck chainmanager.Keeper
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey, memKey sdk.StoreKey,
	paramstore paramtypes.Subspace,
	// stakingKeeper staking.Keeper,
	// chainKeeper chainmanager.Keeper,
	moduleCommunicator ModuleCommunicator,
) *Keeper {
	// set KeyTable if it has not already been set
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		// sk:                 stakingKeeper,
		// ck:                 chainKeeper,
		moduleCommunicator: moduleCommunicator,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
