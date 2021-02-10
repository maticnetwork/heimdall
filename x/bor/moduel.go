package bor

import (
	"encoding/json"
	"fmt"

	"github.com/maticnetwork/heimdall/helper"

	hmTypes "github.com/maticnetwork/heimdall/types"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/maticnetwork/heimdall/x/bor/client/cli"
	"github.com/maticnetwork/heimdall/x/bor/keeper"
	"github.com/maticnetwork/heimdall/x/bor/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

// AppModuleBasic defines the basic application module used by the auth module.

// AppModuleBasic implements the AppModuleBasic interface for the capability module.
type AppModuleBasic struct {
	cdc codec.Marshaler
}

func NewAppModuleBasic(cdc codec.Marshaler) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	types.RegisterCodec(amino)
}

func (a AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (a AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
	panic("implement me")
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, serveMux *runtime.ServeMux) {
	panic("implement me")
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements an application module for the gov module.

type AppModule struct {
	AppModuleBasic

	keeper         keeper.Keeper
	contractCaller helper.IContractCaller
}

func NewAppModule(cdc codec.Marshaler, keeper keeper.Keeper, contractCaller helper.IContractCaller) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		contractCaller: contractCaller,
	}
}
func (a AppModule) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(message, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.ValidateGenesis(genState)
}

func (a AppModule) InitGenesis(context sdk.Context, cdc codec.JSONMarshaler, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)
	InitGenesis(context, a.keeper, genState)
	return []abci.ValidatorUpdate{}
}

func (a AppModule) ExportGenesis(context sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	genState := ExportGenesis(context, a.keeper)
	return cdc.MustMarshalJSON(genState)
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
}

func (a AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(a.keeper))
}

func (a AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (a AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(a.keeper, amino)
}

func (a AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(a.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(a.keeper))
}

func (a AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {
}

func (a AppModule) EndBlock(context sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (a AppModule) NewSideTxHandler() hmTypes.SideTxHandler {
	return NewSideTxHandler(a.keeper, a.contractCaller)
}

// NewPostTxHandler side tx handler
func (a AppModule) NewPostTxHandler() hmTypes.PostTxHandler {
	return NewPostTxHandler(a.keeper)
}
