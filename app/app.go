package app

import (
	"bytes"
	"encoding/json"
	"errors"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	AppName = "Heimdall"

	// internals
	maxGasPerBlock   int64 = 1000000  // 1 Million
	maxBytesPerBlock int64 = 22020096 // 21 MB
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
	caller       helper.ContractCaller
}

var logger = helper.Logger.With("module", "app")

// NewHeimdallApp creates heimdall app
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

	app.masterKeeper = common.NewKeeper(app.cdc, app.keyMaster, app.keyStaker, app.keyCheckpoint, common.DefaultCodespace)

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("we got error", "Error", err)
		cmn.Exit(err.Error())
	}
	app.caller = contractCallerObj

	// register message routes
	app.Router().AddRoute("checkpoint", checkpoint.NewHandler(app.masterKeeper, &app.caller))
	app.Router().AddRoute("staking", staking.NewHandler(app.masterKeeper, &app.caller))
	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.beginBlocker)
	app.SetEndBlocker(app.endBlocker)
	app.SetAnteHandler(auth.NewAnteHandler())

	// mount the multistore and load the latest state
	app.MountStores(app.keyMain, app.keyCheckpoint, app.keyStaker)
	err = app.LoadLatestVersion(app.keyMain)
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
	staking.RegisterWire(cdc)

	cdc.Seal()
	return cdc
}

// MakePulp creates pulp codec and registers custom types for decoder
func MakePulp() *hmTypes.Pulp {
	pulp := hmTypes.GetPulpInstance()

	// register custom type
	checkpoint.RegisterPulp(pulp)
	staking.RegisterPulp(pulp)

	return pulp
}

// BeginBlocker executes before each block
func (app *HeimdallApp) beginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker executes on each end block
func (app *HeimdallApp) endBlocker(ctx sdk.Context, x abci.RequestEndBlock) abci.ResponseEndBlock {
	var valUpdates []abci.ValidatorUpdate

	if ctx.BlockHeader().NumTxs > 0 {

		// --- Start update to new validators

		currentValidatorSet := app.masterKeeper.GetValidatorSet(ctx)
		currentValidatorSetCopy := currentValidatorSet.Copy()
		allValidators := app.masterKeeper.GetAllValidators(ctx)
		//validatorToSignerMap := app.masterKeeper.GetValidatorToSignerMap(ctx)
		ackCount := app.masterKeeper.GetACKCount(ctx)

		// apply updates
		helper.UpdateValidators(
			&currentValidatorSet, // pointer to current validator set -- UpdateValidators will modify it
			allValidators,        // All validators
			ackCount,             // ack count
		)

		// update validator set in store
		err := app.masterKeeper.UpdateValidatorSetInStore(ctx, currentValidatorSet)
		if err != nil {
			logger.Error("Unable to update validator set in state", "Error", err)
		} else {
			// remove all stale validators
			for _, validator := range currentValidatorSetCopy.Validators {
				// validator update
				val := abci.ValidatorUpdate{
					Power:  0,
					PubKey: validator.PubKey.ABCIPubKey(),
				}
				valUpdates = append(valUpdates, val)
			}

			// add new validators
			currentValidatorSet := app.masterKeeper.GetValidatorSet(ctx)
			for _, validator := range currentValidatorSet.Validators {
				// validator update
				val := abci.ValidatorUpdate{
					Power:  int64(validator.Power),
					PubKey: validator.PubKey.ABCIPubKey(),
				}
				valUpdates = append(valUpdates, val)
			}
		}

		// --- End update validators

		// check if ACK is present in cache
		if app.masterKeeper.GetCheckpointCache(ctx, common.CheckpointACKCacheKey) {
			logger.Info("Checkpoint ACK processed in block", "CheckpointACKProcessed", app.masterKeeper.GetCheckpointCache(ctx, common.CheckpointACKCacheKey))

			// --- Start update proposer

			// increment accum
			app.masterKeeper.IncreamentAccum(ctx, 1)

			// log new proposer
			vs := app.masterKeeper.GetValidatorSet(ctx)
			newProposer := vs.GetProposer()
			logger.Debug(
				"New proposer selected",
				"validator", newProposer.ID,
				"signer", newProposer.Signer.String(),
				"power", newProposer.Power,
			)

			// --- End update proposer

			// clear ACK cache
			app.masterKeeper.FlushACKCache(ctx)
		}

		// check if checkpoint is present in cache
		if app.masterKeeper.GetCheckpointCache(ctx, common.CheckpointCacheKey) {
			logger.Info("Checkpoint processed in block", "CheckpointProcessed", app.masterKeeper.GetCheckpointCache(ctx, common.CheckpointCacheKey))
			// Send Checkpoint to Rootchain
			PrepareAndSendCheckpoint(ctx, app.masterKeeper, app.caller)
			// clear Checkpoint cache
			app.masterKeeper.FlushCheckpointCache(ctx)
		}
	}

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		ValidatorUpdates: valUpdates,
	}
}

func (app *HeimdallApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	logger.Info("Loading validators from genesis and setting defaults")
	var genesisState GenesisState
	err := json.Unmarshal(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}
	var isGenesis bool
	if len(genesisState.CurrentValSet.Validators) == 0 {
		isGenesis = true
	} else {
		isGenesis = false
	}

	valSet, valUpdates := app.GetValidatorsFromGenesis(ctx, &genesisState, genesisState.AckCount)
	if len(valSet.Validators) == 0 {
		panic(errors.New("no valid validators found"))
	}

	var currentValSet hmTypes.ValidatorSet
	if isGenesis {
		currentValSet = valSet
	} else {
		currentValSet = genesisState.CurrentValSet
	}

	// TODO match valSet and genesisState.CurrentValSet for difference in accum
	// update validator set in store
	err = app.masterKeeper.UpdateValidatorSetInStore(ctx, currentValSet)
	if err != nil {
		logger.Error("Unable to marshall validator set while adding in store", "Error", err)
		panic(err)
	}

	// increment accumulator if starting from genesis
	if isGenesis {
		app.masterKeeper.IncreamentAccum(ctx, 1)
	}

	// Set initial ack count
	app.masterKeeper.UpdateACKCountWithValue(ctx, genesisState.AckCount)

	// Add checkpoint in buffer
	app.masterKeeper.SetCheckpointBuffer(ctx, genesisState.BufferedCheckpoint)

	// Set Caches
	app.SetCaches(ctx, &genesisState)

	// Set last no-ack
	app.masterKeeper.SetLastNoAck(ctx, genesisState.LastNoACK)

	// Add all headers
	app.InsertHeaders(ctx, &genesisState)

	logger.Info("adding new validators","updates",valUpdates[0].PubKey)
	// TODO make sure old validtors dont go in validator updates ie deactivated validators have to be removed
	// udpate validators
	return abci.ResponseInitChain{
		// validator updates
		Validators: valUpdates,

		// consensus params
		//ConsensusParams: &abci.ConsensusParams{
		//	Block: &abci.BlockParams{
		//		MaxBytes: maxBytesPerBlock,
		//		MaxGas:   maxGasPerBlock,
		//	},
		//	Evidence: &abci.EvidenceParams{},
		//},
	}
}

// returns validator genesis/existing from genesis state
// TODO add check from main chain if genesis information is right,else people can create invalid genesis and distribute
func (app *HeimdallApp) GetValidatorsFromGenesis(ctx sdk.Context, genesisState *GenesisState, ackCount uint64) (newValSet hmTypes.ValidatorSet, valUpdates []abci.ValidatorUpdate) {
	if len(genesisState.GenValidators) > 0 {
		logger.Debug("Loading genesis validators")
		for _, validator := range genesisState.GenValidators {
			hmValidator := validator.HeimdallValidator()
			logger.Debug("gen validator","gen",validator,"hmVal",hmValidator)
			if ok := hmValidator.ValidateBasic(); !ok {
				logger.Error("Invalid validator properties", "validator", hmValidator)
				return
			}
			if !hmValidator.IsCurrentValidator(genesisState.AckCount) {
				logger.Error("Genesis validators should be current validators", "FaultyValidator", hmValidator)
				return
			}
			if ok := newValSet.Add(&hmValidator); !ok {
				panic(errors.New("error while adding genesis validator"))
			} else {
				// Add individual validator to state
				app.masterKeeper.AddValidator(ctx, hmValidator)
				// convert to Validator Update
				updateVal := abci.ValidatorUpdate{
					Power:  int64(validator.Power),
					PubKey: validator.PubKey.ABCIPubKey(),
				}

				// Add validator to validator updated to be processed below
				valUpdates = append(valUpdates, updateVal)
			}
		}
		logger.Debug("Adding validators to state", "ValidatorSet", newValSet, "ValUpdates", valUpdates)
		return
	}

	// read validators
	logger.Debug("Loading validators from state-dump")
	for _, validator := range genesisState.Validators {
		if !validator.ValidateBasic() {
			logger.Error("Invalid validator properties", "validator", validator)
			return
		}
		if ok := newValSet.Add(&validator); !ok {
			panic(errors.New("Error while addings new validator"))
		} else {
			// Add individual validator to state
			app.masterKeeper.AddValidator(ctx, validator)

			// check if validator is current validator
			// add to val updates else skip
			if validator.IsCurrentValidator(ackCount) {
				// convert to Validator Update
				updateVal := abci.ValidatorUpdate{
					Power:  int64(validator.Power),
					PubKey: validator.PubKey.ABCIPubKey(),
				}

				// Add validator to validator updated to be processed below
				valUpdates = append(valUpdates, updateVal)
			}
		}
	}
	logger.Debug("Adding validators to state", "ValidatorSet", newValSet, "ValUpdates", valUpdates)
	return
}

// Set caches like checkpoint and checkpointACK cache
// Incase user needs to retry sending last checkpoint or sending ACK
func (app *HeimdallApp) SetCaches(ctx sdk.Context, genesisState *GenesisState) {
	if genesisState.CheckpointCache {
		logger.Debug("Found checkpoint cache", "CheckpointCache", genesisState.CheckpointCache)
		app.masterKeeper.SetCheckpointCache(ctx, common.DefaultValue)
		return
	}
	if genesisState.CheckpointACKCache {
		logger.Debug("Found checkpoint ACK cache", "CheckpointACKCache", genesisState.CheckpointACKCache)
		app.masterKeeper.SetCheckpointAckCache(ctx, common.DefaultValue)
		return
	}
}

// Insert headers into state
func (app *HeimdallApp) InsertHeaders(ctx sdk.Context, genesisState *GenesisState) {
	if len(genesisState.Headers) != 0 {
		logger.Debug("Trying to add successfull checkpoints", "NoOfHeaders", len(genesisState.Headers))
		if int(genesisState.AckCount) != len(genesisState.Headers) {
			logger.Error("Number of headers and ack count do not match", "HeaderCount", len(genesisState.Headers), "AckCount", genesisState.AckCount)
			panic(errors.New("Incorrect state in state-dump , Please Check "))
		}
		for i, header := range genesisState.Headers {
			checkpointHeaderIndex := helper.GetConfig().ChildBlockInterval * (uint64(i) + 1)
			app.masterKeeper.AddCheckpoint(ctx, checkpointHeaderIndex, header)
		}
	}
	return
}

// ExportAppStateAndValidators export app state and validators
func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmTypes.GenesisValidator, err error) {
	return appState, validators, err
}

// GetExtraData get extra data for checkpoint
func GetExtraData(_checkpoint hmTypes.CheckpointBlockHeader, ctx sdk.Context) []byte {
	logger.Debug("Creating extra data", "startBlock", _checkpoint.StartBlock, "endBlock", _checkpoint.EndBlock, "roothash", _checkpoint.RootHash, "timestamp", _checkpoint.TimeStamp)

	// craft a message
	txBytes, err := helper.CreateTxBytes(
		checkpoint.NewMsgCheckpointBlock(
			_checkpoint.Proposer,
			_checkpoint.StartBlock,
			_checkpoint.EndBlock,
			_checkpoint.RootHash,
			_checkpoint.TimeStamp,
		),
	)
	if err != nil {
		logger.Error("Error decoding transaction data", "error", err)
	}

	return txBytes[hmTypes.PulpHashLength:]
}

// PrepareAndSendCheckpoint prepares all the data required for sending checkpoint and sends tx to rootchain
func PrepareAndSendCheckpoint(ctx sdk.Context, keeper common.Keeper, caller helper.ContractCaller) {
	// fetch votes from block header
	var votes []tmTypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		logger.Error("Error while unmarshalling vote", "error", err)
	}

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
	lastblock, err := caller.CurrentChildBlock()
	if err != nil {
		logger.Error("Could not fetch last block from mainchain", "error", err)
		panic(err)
	}

	// get validator address
	validatorAddress := ethCommon.BytesToAddress(helper.GetPubKey().Address().Bytes())

	// check if we are proposer
	if bytes.Equal(keeper.GetCurrentProposer(ctx).Signer.Bytes(), validatorAddress.Bytes()) {
		logger.Info("We are proposer! Validating if checkpoint needs to be pushed", "commitedLastBlock", lastblock, "startBlock", _checkpoint.StartBlock)
		// check if we need to send checkpoint or not
		if ((lastblock + 1) == _checkpoint.StartBlock) || (lastblock == 0 && _checkpoint.StartBlock == 0) {
			logger.Info("Sending valid checkpoint", "startBlock", _checkpoint.StartBlock)
			caller.SendCheckpoint(helper.GetVoteBytes(votes, ctx), sigs, extraData)
		} else if lastblock > _checkpoint.StartBlock {
			logger.Info("Start block does not match, checkpoint already sent", "commitedLastBlock", lastblock, "startBlock", _checkpoint.StartBlock)
		} else if lastblock > _checkpoint.EndBlock {
			logger.Info("Checkpoint already sent", "commitedLastBlock", lastblock, "startBlock", _checkpoint.StartBlock)
		} else {
			logger.Info("No need to send checkpoint")
		}
	} else {
		logger.Info("We are not proposer", "proposer", keeper.GetCurrentProposer(ctx), "validator", validatorAddress.String())
	}
}
