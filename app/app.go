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
	maxGasPerBlock   sdk.Gas = 1000000  // 1 Million
	maxBytesPerBlock sdk.Gas = 22020096 // 21 MB
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

	app.masterKeeper = common.NewKeeper(app.cdc, app.keyMaster, app.keyStaker, app.keyCheckpoint, app.RegisterCodespace(common.DefaultCodespace))
	contractCallerObj, err := helper.NewContractCallerObj()
	if err != nil {
		cmn.Exit(err.Error())
	}
	// register message routes
	app.Router().AddRoute("checkpoint", checkpoint.NewHandler(app.masterKeeper, contractCallerObj))
	app.Router().AddRoute("staking", staking.NewHandler(app.masterKeeper, contractCallerObj))
	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.beginBlocker)
	app.SetEndBlocker(app.endBlocker)
	app.SetAnteHandler(auth.NewAnteHandler())

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
		validatorToSignerMap := app.masterKeeper.GetValidatorToSignerMap(ctx)
		ackCount := app.masterKeeper.GetACKCount(ctx)

		// apply updates
		helper.UpdateValidators(
			&currentValidatorSet, // pointer to current validator set -- UpdateValidators will modify it
			allValidators,        // All validators
			validatorToSignerMap, // validator to signer map
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
				"validator", newProposer.Address.String(),
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
			PrepareAndSendCheckpoint(ctx, app.masterKeeper)
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

	// initialize validator set
	newValidatorSet := hmTypes.ValidatorSet{}
	validatorUpdates := make([]abci.ValidatorUpdate, len(genesisState.Validators))

	for i, validator := range genesisState.Validators {
		hmValidator := validator.ToHeimdallValidator()

		if ok := newValidatorSet.Add(&hmValidator); !ok {
			panic(errors.New("Error while adding new validator"))
		} else {
			// Add individual validator to state
			app.masterKeeper.AddValidator(ctx, hmValidator)

			// convert to Validator Update
			updateVal := abci.ValidatorUpdate{
				Power:  int64(validator.Power),
				PubKey: validator.PubKey.ABCIPubKey(),
			}

			// Add validator to validator updated to be processed below
			validatorUpdates[i] = updateVal
		}
	}

	// Initial validator set log
	logger.Info("Initial validator set", "size", newValidatorSet.Size())

	// update validator set in store
	err = app.masterKeeper.UpdateValidatorSetInStore(ctx, newValidatorSet)
	if err != nil {
		logger.Error("Unable to marshall validator set while adding in store", "Error", err)
		panic(err)
	}

	// increment accumulator
	app.masterKeeper.IncreamentAccum(ctx, 1)

	//
	// Set initial ack count
	//
	app.masterKeeper.UpdateACKCountWithValue(ctx, genesisState.InitialAckCount)

	// udpate validators
	return abci.ResponseInitChain{
		// validator updates
		Validators: validatorUpdates,

		// consensus params
		ConsensusParams: &abci.ConsensusParams{
			BlockSize: &abci.BlockSize{
				MaxBytes: maxBytesPerBlock,
				MaxGas:   maxGasPerBlock,
			},
			EvidenceParams: &abci.EvidenceParams{},
		},
	}
}

// ExportAppStateAndValidators export app state and validators
func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmTypes.GenesisValidator, err error) {
	//ctx := app.NewContext(true, abci.Header{})
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
func PrepareAndSendCheckpoint(ctx sdk.Context, keeper common.Keeper) {
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
	lastblock, err := helper.CurrentChildBlock()
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
			helper.SendCheckpoint(helper.GetVoteBytes(votes, ctx), sigs, extraData)
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
