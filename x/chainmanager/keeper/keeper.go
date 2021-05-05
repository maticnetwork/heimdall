package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/maticnetwork/heimdall/helper"

	"github.com/maticnetwork/heimdall/x/chainmanager/types"
)

type (
	Keeper struct {
		cdc            codec.Marshaler
		storeKey       sdk.StoreKey
		paramSubspace  paramtypes.Subspace
		contractCaller helper.ContractCaller
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, paramSubspace paramtypes.Subspace, caller helper.ContractCaller) Keeper {
	if !paramSubspace.HasKeyTable() {
		paramSubspace = paramSubspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramSubspace:  paramSubspace,
		contractCaller: caller,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetParams sets the chainmanager module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	k.paramSubspace.SetParamSet(ctx, params)
}

// GetParams gets the chainmanager module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

//
// proposer
//

// GetBlockProposer returns block proposer
func (k Keeper) GetBlockProposer(ctx sdk.Context) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.ProposerKey()) {
		return sdk.AccAddress{}, false
	}

	bz := store.Get(types.ProposerKey())
	return sdk.AccAddress(bz), true
}

// SetBlockProposer sets block proposer
func (k Keeper) SetBlockProposer(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposerKey(), addr.Bytes())
}

// RemoveBlockProposer removes block proposer from store
func (k Keeper) RemoveBlockProposer(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ProposerKey())
}
