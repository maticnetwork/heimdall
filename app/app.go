package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	hmtypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
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
	keyMaster     *sdk.KVStoreKey

	keyStaker *sdk.KVStoreKey
	// manage getting and setting accounts
	//checkpointKeeper checkpoint.Keeper
	//stakerKeeper     staking.Keeper
	masterKeeper common.Keeper
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
		keyMaster:     sdk.NewKVStoreKey("master"),
	}

	// todo give every keeper its own codespace
	app.masterKeeper = common.NewKeeper(app.cdc, app.keyMaster, app.keyStaker, app.keyCheckpoint, app.RegisterCodespace(common.DefaultCodespace))
	// register message routes
	app.Router().AddRoute("checkpoint", checkpoint.NewHandler(app.masterKeeper))
	app.Router().AddRoute("staking", staking.NewHandler(app.masterKeeper))
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
	var valUpdates []abci.ValidatorUpdate

	if ctx.BlockHeader().NumTxs > 0 {
		// check if ACK is present in cache 
		if app.masterKeeper.GetCheckpointCache(ctx, common.CheckpointACKCacheKey) {
			// remove matured Validators
			app.masterKeeper.RemoveDeactivatedValidators(ctx)

			// check if validator set has changed 
			if app.masterKeeper.ValidatorSetChanged(ctx) {
				// GetAllValidators from store (includes previous validator set + updates)
				valUpdates = app.masterKeeper.GetAllValidators(ctx)
				
				// mark validator set changes have been sent to TM 
				app.masterKeeper.SetValidatorSetChangedFlag(ctx,false)
			}
			
			// clear ACK cache
			app.masterKeeper.SetCheckpointAckCache(ctx, common.EmptyBufferValue)
		}
		// check if checkpoint is present in cache 
		if app.masterKeeper.GetCheckpointCache(ctx, common.CheckpointCacheKey) {
			// Send Checkpoint to Rootchain
			PrepareAndSendCheckpoint(ctx, app.masterKeeper)

			// clear Checkpoint cache
			app.masterKeeper.SetCheckpointCache(ctx, common.EmptyBufferValue)
		}
	}

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		ValidatorUpdates: valUpdates,
	}
}

func (app *HeimdallApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	// set ACK count to 0
	app.masterKeeper.InitACKCount(ctx)

	// todo init validator set from store
	//app.stakerKeeper.UpdateValidatorSetInStore()

	return abci.ResponseInitChain{}
}

func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	//ctx := app.NewContext(true, abci.Header{})

	return appState, validators, err
}

// todo try to move this to helper , since it uses checkpoint it causes cycle import error RN
func GetExtraData(_checkpoint hmtypes.CheckpointBlockHeader, ctx sdk.Context) []byte {
	logger.Debug("Creating extra data", "startBlock", _checkpoint.StartBlock, "endBlock", _checkpoint.EndBlock, "roothash", _checkpoint.RootHash)

	// craft a message
	msg := checkpoint.NewMsgCheckpointBlock(_checkpoint.Proposer, _checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash)

	// decoding transaction
	tx := hmtypes.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		logger.Error("Error decoding transaction data", "error", err)
	}

	return txBytes
}

// todo try to move this to helper
// prepares all the data required for sending checkpoint and sends tx to rootchain
func PrepareAndSendCheckpoint(ctx sdk.Context, keeper common.Keeper) {
	// fetch votes from block header
	var votes []tmtypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		logger.Error("Error while unmarshalling vote", "error", err)
	}

	// todo sort sigs before sending
	// get sigs from votes
	sigs := helper.GetSigs(votes)

	// Getting latest checkpoint data from store using height as key and unmarshall
	_checkpoint, err := keeper.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to unmarshall checkpoint while fetching from buffer while preparing checkpoint tx for rootchain", "error", err, "height", ctx.BlockHeight())
		panic(err)
	} else {
		// Get extra data
		extraData := GetExtraData(_checkpoint, ctx)

		//fetch current child block from rootchain contract
		lastblock, err := helper.CurrentChildBlock()
		if err != nil {
			logger.Error("Could not fetch last block from mainchain", "Error", err)
			panic(err)
		}

		// get validator address
		validatorAddress := helper.GetPubKey().Address()

		// check if we are proposer
		if bytes.Equal(keeper.GetCurrentProposerAddress(ctx), validatorAddress.Bytes()) {
			logger.Info("You are proposer ! Validating if checkpoint needs to be pushed")

			// check if we need to send checkpoint or not
			if lastblock == _checkpoint.StartBlock {
				logger.Info("Sending Valid Checkpoint ...")
				helper.SendCheckpoint(helper.GetVoteBytes(votes, ctx), sigs, extraData)
			} else if lastblock > _checkpoint.StartBlock {
				logger.Info("Start block does not match,checkpoint already sent", "lastBlock", lastblock, "startBlock", _checkpoint.StartBlock)
			} else {
				logger.Error("Start Block Ahead of Rootchain header, chains out of sync , time to panic", "lastBlock", lastblock, "startBlock", _checkpoint.StartBlock)
				panic(fmt.Errorf("Ethereum Chain and Heimdall out of sync :("))
			}
		} else {
			logger.Info("You are not proposer", "Proposer", keeper.GetValidatorSet(ctx).Proposer.Address.String(), "You", validatorAddress.String())
		}

	}
}
