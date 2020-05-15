package slashing

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/helper"
	slashingCli "github.com/maticnetwork/heimdall/slashing/client/cli"
	"github.com/maticnetwork/heimdall/slashing/client/rest"
	"github.com/maticnetwork/heimdall/slashing/simulation"
	"github.com/maticnetwork/heimdall/slashing/types"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmModule "github.com/maticnetwork/heimdall/types/module"
	simTypes "github.com/maticnetwork/heimdall/types/simulation"

	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
)

var (
	_ module.AppModule             = AppModule{}
	_ module.AppModuleBasic        = AppModuleBasic{}
	_ hmModule.HeimdallModuleBasic = AppModule{}
	_ hmModule.AppModuleSimulation = AppModule{}
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
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the slashing
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the slashing module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(data)
}

// VerifyGenesis performs verification on slashing module state.
func (AppModuleBasic) VerifyGenesis(bz map[string]json.RawMessage) error {
	var chainManagertData chainmanagerTypes.GenesisState
	errcm := chainmanagerTypes.ModuleCdc.UnmarshalJSON(bz[chainmanagerTypes.ModuleName], &chainManagertData)
	if errcm != nil {
		return errcm
	}

	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz[types.ModuleName], &data)

	if err != nil {
		return err
	}
	return nil
}

// RegisterRESTRoutes registers the REST routes for the slashing module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the slashing module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return slashingCli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the slashing module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return slashingCli.GetQueryCmd(cdc)

}

//____________________________________________________________________________

// AppModule implements an application module for the slashing module.
type AppModule struct {
	AppModuleBasic

	keeper         Keeper
	stakingKeeper  staking.Keeper
	contractCaller helper.IContractCaller
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, sk staking.Keeper, contractCaller helper.IContractCaller) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		stakingKeeper:  sk,
		contractCaller: contractCaller,
	}
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
	return NewHandler(am.keeper, am.contractCaller)
}

// QuerierRoute returns the slashing module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the slashing module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the slashing module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}

}

// ExportGenesis returns the exported genesis state as raw bytes for the auth
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the slashing module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the slashing module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the auth module
func (AppModule) GenerateGenesisState(simState *hmModule.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(simState hmModule.SimulationState) []simTypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized auth param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simTypes.ParamChange {
	return simulation.ParamChanges(r)
}

// RegisterStoreDecoder registers a decoder for auth module's types
func (AppModule) RegisterStoreDecoder(sdr hmModule.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.DecodeStore
}

// WeightedOperations doesn't return any auth module operation.
func (AppModule) WeightedOperations(_ hmModule.SimulationState) []simTypes.WeightedOperation {
	return nil
}

//
// Side module
//

// NewSideTxHandler side tx handler
func (am AppModule) NewSideTxHandler() hmTypes.SideTxHandler {
	return NewSideTxHandler(am.keeper, am.contractCaller)
}

// NewPostTxHandler side tx handler
func (am AppModule) NewPostTxHandler() hmTypes.PostTxHandler {
	return NewPostTxHandler(am.keeper, am.contractCaller)
}
