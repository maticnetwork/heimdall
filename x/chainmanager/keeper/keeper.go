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
