package app

import (
	"fmt"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	jsoniter "github.com/json-iterator/go"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/chainmanager"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/clerk"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	gov "github.com/maticnetwork/heimdall/gov"
	govTypes "github.com/maticnetwork/heimdall/gov/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/params"
	paramsClient "github.com/maticnetwork/heimdall/params/client"
	"github.com/maticnetwork/heimdall/params/subspace"
	paramsTypes "github.com/maticnetwork/heimdall/params/types"
	"github.com/maticnetwork/heimdall/sidechannel"
	sidechannelTypes "github.com/maticnetwork/heimdall/sidechannel/types"
	"github.com/maticnetwork/heimdall/slashing"
	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	"github.com/maticnetwork/heimdall/staking"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	"github.com/maticnetwork/heimdall/supply"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
	"github.com/maticnetwork/heimdall/topup"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	"github.com/maticnetwork/heimdall/types"
	hmModule "github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/version"
)

const (
	// AppName denotes app name
	AppName = "Heimdall"
)

// Assertion for Heimdall app
var _ App = &HeimdallApp{}

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		params.AppModuleBasic{},
		sidechannel.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		supply.AppModuleBasic{},
		chainmanager.AppModuleBasic{},
		staking.AppModuleBasic{},
		checkpoint.AppModuleBasic{},
		bor.AppModuleBasic{},
		clerk.AppModuleBasic{},
		topup.AppModuleBasic{},
		slashing.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsClient.ProposalHandler),
	)

	// module account permissions
	maccPerms = map[string][]string{
		authTypes.FeeCollectorName: nil,
		govTypes.ModuleName:        {},
	}
)

// HeimdallApp main heimdall app
type HeimdallApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]subspace.Subspace

	// side router
	sideRouter types.SideRouter

	// keepers
	SidechannelKeeper sidechannel.Keeper
	AccountKeeper     auth.AccountKeeper
	BankKeeper        bank.Keeper
	SupplyKeeper      supply.Keeper
	GovKeeper         gov.Keeper
	ChainKeeper       chainmanager.Keeper
	CheckpointKeeper  checkpoint.Keeper
	StakingKeeper     staking.Keeper
	BorKeeper         bor.Keeper
	ClerkKeeper       clerk.Keeper
	TopupKeeper       topup.Keeper
	SlashingKeeper    slashing.Keeper

	// param keeper
	ParamsKeeper params.Keeper

	// contract keeper
	caller helper.ContractCaller

	//  total coins supply
	TotalCoinsSupply sdk.Coins

	// the module manager
	mm *module.Manager

	// simulation module manager
	sm *hmModule.SimulationManager
}

var logger = helper.Logger.With("module", "app")

//
// Module communicator
//

// ModuleCommunicator retriever
type ModuleCommunicator struct {
	App *HeimdallApp
}

// GetACKCount returns ack count
func (d ModuleCommunicator) GetACKCount(ctx sdk.Context) uint64 {
	return d.App.CheckpointKeeper.GetACKCount(ctx)
}

// IsCurrentValidatorByAddress check if validator is current validator
func (d ModuleCommunicator) IsCurrentValidatorByAddress(ctx sdk.Context, address []byte) bool {
	return d.App.StakingKeeper.IsCurrentValidatorByAddress(ctx, address)
}

// GetAllDividendAccounts fetches all dividend accounts from topup module
func (d ModuleCommunicator) GetAllDividendAccounts(ctx sdk.Context) []types.DividendAccount {
	return d.App.TopupKeeper.GetAllDividendAccounts(ctx)
}

// GetValidatorFromValID get validator from validator id
func (d ModuleCommunicator) GetValidatorFromValID(ctx sdk.Context, valID types.ValidatorID) (validator types.Validator, ok bool) {
	return d.App.StakingKeeper.GetValidatorFromValID(ctx, valID)
}

// SetCoins sets coins
func (d ModuleCommunicator) SetCoins(ctx sdk.Context, addr types.HeimdallAddress, amt sdk.Coins) sdk.Error {
	return d.App.BankKeeper.SetCoins(ctx, addr, amt)
}

// GetCoins gets coins
func (d ModuleCommunicator) GetCoins(ctx sdk.Context, addr types.HeimdallAddress) sdk.Coins {
	return d.App.BankKeeper.GetCoins(ctx, addr)
}

// SendCoins transfers coins
func (d ModuleCommunicator) SendCoins(ctx sdk.Context, fromAddr types.HeimdallAddress, toAddr types.HeimdallAddress, amt sdk.Coins) sdk.Error {
	return d.App.BankKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// CreateValidatorSigningInfo used by slashing module
func (d ModuleCommunicator) CreateValidatorSigningInfo(ctx sdk.Context, valID types.ValidatorID, valSigningInfo types.ValidatorSigningInfo) {
	d.App.SlashingKeeper.SetValidatorSigningInfo(ctx, valID, valSigningInfo)
}

//
// Heimdall app
//

// NewHeimdallApp creates heimdall app
func NewHeimdallApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *HeimdallApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// set prefix
	config := sdk.GetConfig()
	config.Seal()

	// base app
	bApp := bam.NewBaseApp(AppName, logger, db, authTypes.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(nil)
	bApp.SetAppVersion(version.Version)

	// keys
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		sidechannelTypes.StoreKey,
		authTypes.StoreKey,
		bankTypes.StoreKey,
		supplyTypes.StoreKey,
		govTypes.StoreKey,
		chainmanagerTypes.StoreKey,
		stakingTypes.StoreKey,
		slashingTypes.StoreKey,
		checkpointTypes.StoreKey,
		borTypes.StoreKey,
		clerkTypes.StoreKey,
		topupTypes.StoreKey,
		paramsTypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)

	// create heimdall app
	var app = &HeimdallApp{
		cdc:       cdc,
		BaseApp:   bApp,
		keys:      keys,
		tkeys:     tkeys,
		subspaces: make(map[string]subspace.Subspace),
	}

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(app.cdc, keys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey], paramsTypes.DefaultCodespace)
	app.subspaces[sidechannelTypes.ModuleName] = app.ParamsKeeper.Subspace(sidechannelTypes.DefaultParamspace)
	app.subspaces[authTypes.ModuleName] = app.ParamsKeeper.Subspace(authTypes.DefaultParamspace)
	app.subspaces[bankTypes.ModuleName] = app.ParamsKeeper.Subspace(bankTypes.DefaultParamspace)
	app.subspaces[supplyTypes.ModuleName] = app.ParamsKeeper.Subspace(supplyTypes.DefaultParamspace)
	app.subspaces[govTypes.ModuleName] = app.ParamsKeeper.Subspace(govTypes.DefaultParamspace).WithKeyTable(govTypes.ParamKeyTable())
	app.subspaces[chainmanagerTypes.ModuleName] = app.ParamsKeeper.Subspace(chainmanagerTypes.DefaultParamspace)
	app.subspaces[stakingTypes.ModuleName] = app.ParamsKeeper.Subspace(stakingTypes.DefaultParamspace)
	app.subspaces[slashingTypes.ModuleName] = app.ParamsKeeper.Subspace(slashingTypes.DefaultParamspace)
	app.subspaces[checkpointTypes.ModuleName] = app.ParamsKeeper.Subspace(checkpointTypes.DefaultParamspace)
	app.subspaces[borTypes.ModuleName] = app.ParamsKeeper.Subspace(borTypes.DefaultParamspace)
	app.subspaces[clerkTypes.ModuleName] = app.ParamsKeeper.Subspace(clerkTypes.DefaultParamspace)
	app.subspaces[topupTypes.ModuleName] = app.ParamsKeeper.Subspace(topupTypes.DefaultParamspace)

	//
	// Contract caller
	//

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.caller = contractCallerObj

	//
	// module communicator
	//

	moduleCommunicator := ModuleCommunicator{App: app}

	//
	// keepers
	//

	// create side channel keeper
	app.SidechannelKeeper = sidechannel.NewKeeper(
		app.cdc,
		keys[sidechannelTypes.StoreKey], // target store
		app.subspaces[sidechannelTypes.ModuleName],
		common.DefaultCodespace,
	)

	// create chain keeper
	app.ChainKeeper = chainmanager.NewKeeper(
		app.cdc,
		keys[chainmanagerTypes.StoreKey], // target store
		app.subspaces[chainmanagerTypes.ModuleName],
		common.DefaultCodespace,
		app.caller,
	)

	// account keeper
	app.AccountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[authTypes.StoreKey], // target store
		app.subspaces[authTypes.ModuleName],
		authTypes.ProtoBaseAccount, // prototype
	)

	app.StakingKeeper = staking.NewKeeper(
		app.cdc,
		keys[stakingTypes.StoreKey], // target store
		app.subspaces[stakingTypes.ModuleName],
		common.DefaultCodespace,
		app.ChainKeeper,
		moduleCommunicator,
	)

	app.SlashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashingTypes.StoreKey], // target store
		app.StakingKeeper,
		app.subspaces[slashingTypes.ModuleName],
		common.DefaultCodespace,
		app.ChainKeeper,
	)

	// bank keeper
	app.BankKeeper = bank.NewKeeper(
		app.cdc,
		keys[bankTypes.StoreKey], // target store
		app.subspaces[bankTypes.ModuleName],
		bankTypes.DefaultCodespace,
		app.AccountKeeper,
		moduleCommunicator,
	)

	// bank keeper
	app.SupplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supplyTypes.StoreKey], // target store
		app.subspaces[supplyTypes.ModuleName],
		maccPerms,
		app.AccountKeeper,
		app.BankKeeper,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.
		AddRoute(govTypes.RouterKey, govTypes.ProposalHandler).
		AddRoute(paramsTypes.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))

	app.GovKeeper = gov.NewKeeper(
		app.cdc,
		keys[govTypes.StoreKey],
		app.subspaces[govTypes.ModuleName],
		app.SupplyKeeper,
		app.StakingKeeper,
		govTypes.DefaultCodespace,
		govRouter,
	)

	app.CheckpointKeeper = checkpoint.NewKeeper(
		app.cdc,
		keys[checkpointTypes.StoreKey], // target store
		app.subspaces[checkpointTypes.ModuleName],
		common.DefaultCodespace,
		app.StakingKeeper,
		app.ChainKeeper,
		moduleCommunicator,
	)

	app.BorKeeper = bor.NewKeeper(
		app.cdc,
		keys[borTypes.StoreKey], // target store
		app.subspaces[borTypes.ModuleName],
		common.DefaultCodespace,
		app.ChainKeeper,
		app.StakingKeeper,
		app.caller,
	)

	app.ClerkKeeper = clerk.NewKeeper(
		app.cdc,
		keys[clerkTypes.StoreKey], // target store
		app.subspaces[clerkTypes.ModuleName],
		common.DefaultCodespace,
		app.ChainKeeper,
	)

	// may be need signer
	app.TopupKeeper = topup.NewKeeper(
		app.cdc,
		keys[topupTypes.StoreKey],
		app.subspaces[topupTypes.ModuleName],
		topupTypes.DefaultCodespace,
		app.ChainKeeper,
		app.BankKeeper,
		app.StakingKeeper,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		sidechannel.NewAppModule(app.SidechannelKeeper),
		auth.NewAppModule(app.AccountKeeper, &app.caller, []authTypes.AccountProcessor{
			supplyTypes.AccountProcessor,
		}),
		bank.NewAppModule(app.BankKeeper, &app.caller),
		supply.NewAppModule(app.SupplyKeeper, &app.caller),
		gov.NewAppModule(app.GovKeeper, app.SupplyKeeper),
		chainmanager.NewAppModule(app.ChainKeeper, &app.caller),
		staking.NewAppModule(app.StakingKeeper, &app.caller),
		slashing.NewAppModule(app.SlashingKeeper, app.StakingKeeper, &app.caller),
		checkpoint.NewAppModule(app.CheckpointKeeper, app.StakingKeeper, app.TopupKeeper, &app.caller),
		bor.NewAppModule(app.BorKeeper, &app.caller),
		clerk.NewAppModule(app.ClerkKeeper, &app.caller),
		topup.NewAppModule(app.TopupKeeper, &app.caller),
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		sidechannelTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		govTypes.ModuleName,
		chainmanagerTypes.ModuleName,
		supplyTypes.ModuleName,
		stakingTypes.ModuleName,
		slashingTypes.ModuleName,
		checkpointTypes.ModuleName,
		borTypes.ModuleName,
		clerkTypes.ModuleName,
		topupTypes.ModuleName,
	)

	// register message routes and query routes
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// side router
	app.sideRouter = types.NewSideRouter()
	for _, m := range app.mm.Modules {
		if m.Route() != "" {
			if sm, ok := m.(hmModule.SideModule); ok {
				app.sideRouter.AddRoute(m.Route(), &types.SideHandlers{
					SideTxHandler: sm.NewSideTxHandler(),
					PostTxHandler: sm.NewPostTxHandler(),
				})
			}
		}
	}

	app.sideRouter.Seal()

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = hmModule.NewSimulationManager(
		auth.NewAppModule(app.AccountKeeper, &app.caller, []authTypes.AccountProcessor{
			supplyTypes.AccountProcessor,
		}),

		slashing.NewAppModule(app.SlashingKeeper, app.StakingKeeper, &app.caller),
		chainmanager.NewAppModule(app.ChainKeeper, &app.caller),
		topup.NewAppModule(app.TopupKeeper, &app.caller),
		staking.NewAppModule(app.StakingKeeper, &app.caller),
		checkpoint.NewAppModule(app.CheckpointKeeper, app.StakingKeeper, app.TopupKeeper, &app.caller),
		bank.NewAppModule(app.BankKeeper, &app.caller),
	)
	app.sm.RegisterStoreDecoders()

	// mount the multistore and load the latest state
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// perform initialization logic
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.AccountKeeper,
			app.ChainKeeper,
			app.SupplyKeeper,
			&app.caller,
			auth.DefaultSigVerificationGasConsumer,
		),
	)
	// side-tx processor
	app.SetPostDeliverTxHandler(app.PostDeliverTxHandler)
	app.SetBeginSideBlocker(app.BeginSideBlocker)
	app.SetDeliverSideTxHandler(app.DeliverSideTxHandler)

	// load latest version
	err = app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()

	return app
}

// MakeCodec create codec
func MakeCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	ModuleBasics.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// Name returns the name of the App
func (app *HeimdallApp) Name() string { return app.BaseApp.Name() }

// InitChainer initializes chain
func (app *HeimdallApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	// get validator updates
	if err := ModuleBasics.ValidateGenesis(genesisState); err != nil {
		panic(err)
	}

	// check fee collector module account
	if moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, authTypes.FeeCollectorName); moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", authTypes.FeeCollectorName))
	}

	// init genesis
	app.mm.InitGenesis(ctx, genesisState)

	stakingState := stakingTypes.GetGenesisStateFromAppState(genesisState)
	checkpointState := checkpointTypes.GetGenesisStateFromAppState(genesisState)

	// check if validator is current validator
	// add to val updates else skip
	var valUpdates []abci.ValidatorUpdate

	for _, validator := range stakingState.Validators {
		if validator.IsCurrentValidator(checkpointState.AckCount) {
			// convert to Validator Update
			updateVal := abci.ValidatorUpdate{
				Power:  validator.VotingPower,
				PubKey: validator.PubKey.ABCIPubKey(),
			}
			// Add validator to validator updated to be processed below
			valUpdates = append(valUpdates, updateVal)
		}
	}

	// TODO make sure old validators dont go in validator updates ie deactivated validators have to be removed
	// update validators
	return abci.ResponseInitChain{
		// validator updates
		Validators: valUpdates,
	}
}

// BeginBlocker application updates every begin block
func (app *HeimdallApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.AccountKeeper.SetBlockProposer(
		ctx,
		types.BytesToHeimdallAddress(req.Header.GetProposerAddress()),
	)

	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker executes on each end block
func (app *HeimdallApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	// transfer fees to current proposer
	if proposer, ok := app.AccountKeeper.GetBlockProposer(ctx); ok {
		moduleAccount := app.SupplyKeeper.GetModuleAccount(ctx, authTypes.FeeCollectorName)

		amount := moduleAccount.GetCoins().AmountOf(authTypes.FeeToken)
		if !amount.IsZero() {
			coins := sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: amount}}
			if err := app.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, authTypes.FeeCollectorName, proposer, coins); err != nil {
				logger.Error("EndBlocker | SendCoinsFromModuleToAccount", "Error", err)
			}
		}

		// remove block proposer
		app.AccountKeeper.RemoveBlockProposer(ctx)
	}

	var tmValUpdates []abci.ValidatorUpdate

	// --- Start update to new validators
	currentValidatorSet := app.StakingKeeper.GetValidatorSet(ctx)
	allValidators := app.StakingKeeper.GetAllValidators(ctx)
	ackCount := app.CheckpointKeeper.GetACKCount(ctx)

	// get validator updates
	setUpdates := helper.GetUpdatedValidators(
		&currentValidatorSet, // pointer to current validator set -- UpdateValidators will modify it
		allValidators,        // All validators
		ackCount,             // ack count
	)

	if len(setUpdates) > 0 {
		// create new validator set
		if err := currentValidatorSet.UpdateWithChangeSet(setUpdates); err != nil {
			// return with nothing
			logger.Error("Unable to update current validator set", "Error", err)
			return abci.ResponseEndBlock{}
		}

		//Hardfork to remove the rotation of validator list on stake update
		if ctx.BlockHeight() < helper.GetAalborgHardForkHeight() {
			// increment proposer priority
			currentValidatorSet.IncrementProposerPriority(1)
		}

		// validator set change
		logger.Debug("[ENDBLOCK] Updated current validator set", "proposer", currentValidatorSet.GetProposer())

		// save set in store
		if err := app.StakingKeeper.UpdateValidatorSetInStore(ctx, currentValidatorSet); err != nil {
			// return with nothing
			logger.Error("Unable to update current validator set in state", "Error", err)
			return abci.ResponseEndBlock{}
		}

		// convert updates from map to array
		for _, v := range setUpdates {
			tmValUpdates = append(tmValUpdates, abci.ValidatorUpdate{
				Power:  v.VotingPower,
				PubKey: v.PubKey.ABCIPubKey(),
			})
		}
	}

	// Change root chain contract addresses if required
	if chainManagerAddressMigration, found := helper.GetChainManagerAddressMigration(ctx.BlockHeight()); found {
		params := app.ChainKeeper.GetParams(ctx)

		params.ChainParams.MaticTokenAddress = chainManagerAddressMigration.MaticTokenAddress
		params.ChainParams.StakingManagerAddress = chainManagerAddressMigration.StakingManagerAddress
		params.ChainParams.RootChainAddress = chainManagerAddressMigration.RootChainAddress
		params.ChainParams.SlashManagerAddress = chainManagerAddressMigration.SlashManagerAddress
		params.ChainParams.StakingInfoAddress = chainManagerAddressMigration.StakingInfoAddress
		params.ChainParams.StateSenderAddress = chainManagerAddressMigration.StateSenderAddress

		// update chain manager state
		app.ChainKeeper.SetParams(ctx, params)
		logger.Info("Updated chain manager state", "params", params)
	}

	// end block
	app.mm.EndBlock(ctx, req)

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		ValidatorUpdates: tmValUpdates,
	}
}

// LoadHeight loads a particular height
func (app *HeimdallApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *HeimdallApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supplyTypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns HeimdallApp's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *HeimdallApp) Codec() *codec.Codec {
	return app.cdc
}

// SetCodec set codec to app
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *HeimdallApp) SetCodec(cdc *codec.Codec) {
	app.cdc = cdc
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *HeimdallApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *HeimdallApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *HeimdallApp) GetSubspace(moduleName string) subspace.Subspace {
	return app.subspaces[moduleName]
}

// GetSideRouter returns side-tx router
func (app *HeimdallApp) GetSideRouter() types.SideRouter {
	return app.sideRouter
}

// SetSideRouter sets side-tx router
// Testing ONLY
func (app *HeimdallApp) SetSideRouter(r types.SideRouter) {
	app.sideRouter = r
}

// GetModuleManager returns module manager
//
// NOTE: This is solely to be used for testing purposes.
func (app *HeimdallApp) GetModuleManager() *module.Manager {
	return app.mm
}

// SimulationManager implements the SimulationApp interface
func (app *HeimdallApp) SimulationManager() *hmModule.SimulationManager {
	return app.sm
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}
