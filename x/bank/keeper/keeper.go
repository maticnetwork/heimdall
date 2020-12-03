package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/bank/types"
	auth "github.com/maticnetwork/heimdall/x/auth/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keys for supply store
// Items are stored with the following key: values
//
// - 0x00: Supply
var (
	SupplyKey = []byte{0x00}
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		paramSubspace      paramtypes.Subspace
		ak auth.Keeper
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey, memKey sdk.StoreKey, ak auth.Keeper) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		ak: ak,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Context) (supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(SupplyKey)
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

// SetSupply sets the Supply to store
func (k Keeper) SetSupply(ctx sdk.Context, supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&supply)
	store.Set(SupplyKey, b)
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper Keeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSubspace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper Keeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSubspace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}