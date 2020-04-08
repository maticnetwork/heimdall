package slashing

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"

	// "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	_ module.AppModule            = AppModule{}
	_ module.AppModuleBasic       = AppModuleBasic{}
	_ hmTypes.HeimdallModuleBasic = AppModule{}
	// _ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the slashing module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the slashing module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the slashing module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	// RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the slashing
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the slashing module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// VerifyGenesis performs verification on auth module state.
func (AppModuleBasic) VerifyGenesis(bz map[string]json.RawMessage) error {
	// var chainManagertData chainmanagerTypes.GenesisState
	// errcm := chainmanagerTypes.ModuleCdc.UnmarshalJSON(bz[chainmanagerTypes.ModuleName], &chainManagertData)
	// if errcm != nil {
	// 	return errcm
	// }

	// var data types.GenesisState
	// err := types.ModuleCdc.UnmarshalJSON(bz[types.ModuleName], &data)

	// if err != nil {
	// 	return err
	// }
	// return verifyGenesis(data, chainManagertData)
	// TODO - slashing changes
	return nil
}

// RegisterRESTRoutes registers the REST routes for the slashing module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	// rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the slashing module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	// return cli.GetTxCmd(cdc)
	// TODO - slashing
	return nil
}

// GetQueryCmd returns no root query command for the slashing module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// return cli.GetQueryCmd(StoreKey, cdc)
	// TODO - slashing
	return nil
}

//____________________________________________________________________________

// AppModule implements an application module for the slashing module.
type AppModule struct {
	AppModuleBasic

	// keeper        Keeper
	// accountKeeper types.AccountKeeper
	// bankKeeper    types.BankKeeper
	// stakingKeeper stakingkeeper.Keeper
}

// // NewAppModule creates a new AppModule object
// func NewAppModule(keeper Keeper, ak types.AccountKeeper, bk types.BankKeeper, sk stakingkeeper.Keeper) AppModule {
// 	return AppModule{
// 		AppModuleBasic: AppModuleBasic{},
// 		keeper:         keeper,
// 		accountKeeper:  ak,
// 		bankKeeper:     bk,
// 		stakingKeeper:  sk,
// 	}
// }

// TODO - remove this
// NewAppModule creates a new AppModule object
func NewAppModule() AppModule {
	return AppModule{}
}

// Name returns the slashing module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the slashing module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the slashing module.
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the slashing module.
func (am AppModule) NewHandler() sdk.Handler {
	// return NewHandler(am.keeper)
	// TODO - slashing handler
	return nil
}

// QuerierRoute returns the slashing module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the slashing module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	// return NewQuerier(am.keeper)
	// TODO - slashing
	return nil
}

// InitGenesis performs genesis initialization for the slashing module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	// var genesisState types.GenesisState
	// ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	// InitGenesis(ctx, am.keeper, am.stakingKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the auth
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	// gs := ExportGenesis(ctx, am.keeper)
	// return types.ModuleCdc.MustMarshalJSON(gs)
	return nil
}

// BeginBlock returns the begin blocker for the slashing module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	// BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the slashing module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
