package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
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

	// create and register pulp codec
	pulp := hmTypes.GetPulpInstance()

	// create your application type
	var app = &HeimdallApp{
		cdc:           cdc,
		BaseApp:       bam.NewBaseApp(AppName, logger, db, hmTypes.RLPTxDecoder(pulp), baseAppOptions...),
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

// MakeCodec create codec
func MakeCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)

	// custom types
	checkpoint.RegisterWire(cdc)

	cdc.Seal()
	return cdc
}

// MakePulp creates pulp codec and registers custom types for decoder
func MakePulp() *hmTypes.Pulp {
	pulp := hmTypes.GetPulpInstance()

	// register custom type
	checkpoint.RegisterPulp(pulp)

	return pulp
}

// BeginBlocker executes before each block
func (app *HeimdallApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker executes on each end block
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
				app.masterKeeper.SetValidatorSetChangedFlag(ctx, false)
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
	var genesisState GenesisState
	err := json.Unmarshal(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	// set ACK count to 0
	// app.masterKeeper.InitACKCount(ctx)

	// initialize validator set
	newValidatorSet := tmTypes.ValidatorSet{}
	validatorUpdates := make([]abci.ValidatorUpdate, 1)

	for i, validator := range genesisState.Validators {
		// add val
		tmValidator := validator.ToTmValidator()
		if ok := newValidatorSet.Add(&tmValidator); !ok {
			panic(errors.New("Error while adding new validator"))
		} else {
			// convert to Validator Update
			updateVal := abci.ValidatorUpdate{
				Power:  tmValidator.VotingPower,
				PubKey: tmTypes.TM2PB.PubKey(tmValidator.PubKey),
			}
			validatorUpdates[i] = updateVal
		}
	}

	// Initial validator set log
	logger.Info("Initial validator set", "size", newValidatorSet.Size())

	// update validator set in store
	app.masterKeeper.UpdateValidatorSetInStore(ctx, newValidatorSet)

	// increment accumulator
	app.masterKeeper.IncreamentAccum(ctx, 1)

	// udpate validators
	return abci.ResponseInitChain{
		Validators: validatorUpdates,
	}
}

func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmTypes.GenesisValidator, err error) {
	//ctx := app.NewContext(true, abci.Header{})
	return appState, validators, err
}

// todo try to move this to helper , since it uses checkpoint it causes cycle import error RN
func GetExtraData(_checkpoint hmTypes.CheckpointBlockHeader, ctx sdk.Context) []byte {
	logger.Debug("Creating extra data", "startBlock", _checkpoint.StartBlock, "endBlock", _checkpoint.EndBlock, "roothash", _checkpoint.RootHash)

	// craft a message
	msg := checkpoint.NewMsgCheckpointBlock(_checkpoint.Proposer, _checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash)

	// decoding transaction
	tx := hmTypes.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		logger.Error("Error decoding transaction data", "error", err)
	}

	return txBytes
}

// PrepareAndSendCheckpoint prepares all the data required for sending checkpoint and sends tx to rootchain
func PrepareAndSendCheckpoint(ctx sdk.Context, keeper common.Keeper) {
	// fetch votes from block header
	var votes []tmTypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		logger.Error("Error while unmarshalling vote", "error", err)
	}

	// TODO sort sigs before sending
	// get sigs from votes
	sigs := helper.GetSigs(votes)

	// Getting latest checkpoint data from store using height as key and unmarshall
	_checkpoint, err := keeper.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to unmarshall checkpoint from buffer while preparing checkpoint tx", "error", err, "height", ctx.BlockHeight())
		return
	}

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
			logger.Info("Sending Valid Checkpoint...")
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
