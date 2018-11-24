package app

import (
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	hmtypes "github.com/maticnetwork/heimdall/types"
)

const (
	AppName = "Heimdall"
)

type HeimdallApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the multistore
	keyMain       *sdk.KVStoreKey
	keyCheckpoint *sdk.KVStoreKey
	keyStake      *sdk.KVStoreKey

	keyStaker *sdk.KVStoreKey
	// manage getting and setting accounts
	checkpointKeeper checkpoint.Keeper
	stakerKeeper     staking.Keeper
}

var logger = helper.Logger.With("module", "app")

func NewHeimdallApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *HeimdallApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create your application type
	var app = &HeimdallApp{
		cdc:           cdc,
		BaseApp:       bam.NewBaseApp(AppName, logger, db, hmtypes.RLPTxDecoder(), baseAppOptions...),
		keyMain:       sdk.NewKVStoreKey("main"),
		keyCheckpoint: sdk.NewKVStoreKey("checkpoint"),
		keyStaker:     sdk.NewKVStoreKey("staker"),
	}

	app.checkpointKeeper = checkpoint.NewKeeper(app.cdc, app.keyCheckpoint, app.RegisterCodespace(checkpoint.DefaultCodespace))
	app.stakerKeeper = staking.NewKeeper(app.cdc, app.keyStaker, app.RegisterCodespace(checkpoint.DefaultCodespace))
	// register message routes
	app.Router().AddRoute("checkpoint", checkpoint.NewHandler(app.checkpointKeeper))
	app.Router().AddRoute("staking", staking.NewHandler(app.stakerKeeper))
	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyCheckpoint, app.keyStaker)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()
	return app
}

func MakeCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	checkpoint.RegisterWire(cdc)
	// register custom type

	cdc.Seal()
	return cdc
}

func (app *HeimdallApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

func (app *HeimdallApp) EndBlocker(ctx sdk.Context, x abci.RequestEndBlock) abci.ResponseEndBlock {
	var validators []abci.ValidatorUpdate

	if ctx.BlockHeader().NumTxs > 0 {
		if app.checkpointKeeper.GetCheckpointCache(ctx, checkpoint.CheckpointACKCacheKey) {
			// remove matured Validators
			app.stakerKeeper.RemoveDeactivatedValidators(ctx)

			// fetch validators from store
			//validators:=app.stakerKeeper.GetAllValidators(ctx)
			// todo populate validators and send to TM

			// clear ACK cache
			app.checkpointKeeper.SetCheckpointAckCache(ctx, checkpoint.EmptyBufferValue)
		}
		if app.checkpointKeeper.GetCheckpointCache(ctx, checkpoint.CheckpointCacheKey) {
			// Send Checkpoint to Rootchain
			helper.PrepareAndSendCheckpoint(ctx, app.checkpointKeeper)

			// clear Checkpoint cache
			app.checkpointKeeper.SetCheckpointCache(ctx, checkpoint.EmptyBufferValue)
		}
	}

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		ValidatorUpdates: validators,
	}
}

func (app *HeimdallApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	// set ACK count to 0
	app.checkpointKeeper.InitACKCount(ctx)

	// init validator set from store

	return abci.ResponseInitChain{}
}

func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	//ctx := app.NewContext(true, abci.Header{})

	return appState, validators, err
}
