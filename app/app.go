package app

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	ethCommon "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"

	hmAuth "github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	// AppName denotes app name
	AppName = "Heimdall"
	// ABCIPubKeyTypeSecp256k1 denotes pub key type
	ABCIPubKeyTypeSecp256k1 = "secp256k1"
	// internals
	maxGasPerBlock   int64 = 10000000 // 10 Million
	maxBytesPerBlock int64 = 22020096 // 21 MB
)

// HeimdallApp main heimdall app
type HeimdallApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the multistore
	keyAccount       *sdk.KVStoreKey
	keyGov           *sdk.KVStoreKey
	keyCheckpoint    *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	keyBor           *sdk.KVStoreKey
	keyMain          *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tKeyParams       *sdk.TransientStoreKey

	accountKeeper       auth.AccountKeeper
	bankKeeper          bank.Keeper
	govKeeper           gov.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	paramsKeeper        params.Keeper

	checkpointKeeper checkpoint.Keeper
	stakingKeeper    staking.Keeper
	borKeeper        bor.Keeper

	// masterKeeper common.Keeper
	caller helper.ContractCaller
}

var logger = helper.Logger.With("module", "app")

// NewHeimdallApp creates heimdall app
func NewHeimdallApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *HeimdallApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create and register pulp codec
	pulp := authTypes.GetPulpInstance()

	// set prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(hmTypes.PrefixAccAddr, hmTypes.PrefixAccPub)
	config.SetBech32PrefixForValidator(hmTypes.PrefixValAddr, hmTypes.PrefixValPub)
	config.SetBech32PrefixForConsensusNode(hmTypes.PrefixConsAddr, hmTypes.PrefixConsPub)
	config.Seal()

	// create your application type
	var app = &HeimdallApp{
		cdc:        cdc,
		BaseApp:    bam.NewBaseApp(AppName, logger, db, authTypes.RLPTxDecoder(pulp), baseAppOptions...),
		keyMain:    sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount: sdk.NewKVStoreKey(authTypes.StoreKey),
		// keyGov:        sdk.NewKVStoreKey(gov.StoreKey),
		keyCheckpoint:    sdk.NewKVStoreKey(checkpointTypes.StoreKey),
		keyStaking:       sdk.NewKVStoreKey(stakingTypes.StoreKey),
		keyBor:           sdk.NewKVStoreKey(borTypes.StoreKey),
		keyFeeCollection: sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyParams:        sdk.NewKVStoreKey("params"),
		tKeyParams:       sdk.NewTransientStoreKey("transient_params"),
	}

	// define param keeper
	app.paramsKeeper = params.NewKeeper(cdc, app.keyParams, app.tKeyParams)

	// account keeper
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount, // target store
		app.paramsKeeper.Subspace(authTypes.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	// bank keeper
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	// fee collector
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		app.cdc,
		app.keyFeeCollection,
	)

	// app.govKeeper = gov.NewKeeper(
	// 	app.cdc,
	// 	app.keyGov,
	// 	app.paramsKeeper, app.paramsKeeper.Subspace(gov.DefaultParamspace), app.bankKeeper, &stakingKeeper,
	// 	gov.DefaultCodespace,
	// )

	app.stakingKeeper = staking.NewKeeper(
		app.cdc,
		app.keyStaking,
		app.paramsKeeper.Subspace(stakingTypes.DefaultParamspace),
		common.DefaultCodespace,
	)
	app.checkpointKeeper = checkpoint.NewKeeper(
		app.cdc,
		app.stakingKeeper,
		app.keyCheckpoint,
		app.paramsKeeper.Subspace(checkpointTypes.DefaultParamspace),
		common.DefaultCodespace,
	)
	app.borKeeper = bor.NewKeeper(
		app.cdc,
		app.stakingKeeper,
		app.keyBor,
		app.paramsKeeper.Subspace(borTypes.DefaultParamspace),
		common.DefaultCodespace,
	)

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		cmn.Exit(err.Error())
	}
	app.caller = contractCallerObj

	// register message routes
	app.Router().AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper))
	app.Router().AddRoute(checkpointTypes.RouterKey, checkpoint.NewHandler(app.checkpointKeeper, &app.caller))
	app.Router().AddRoute(stakingTypes.RouterKey, staking.NewHandler(app.stakingKeeper, &app.caller))
	app.Router().AddRoute(borTypes.RouterKey, bor.NewHandler(app.borKeeper))

	// query routes
	app.QueryRouter().AddRoute(authTypes.QuerierRoute, auth.NewQuerier(app.accountKeeper))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.beginBlocker)
	app.SetEndBlocker(app.endBlocker)
	app.SetAnteHandler(hmAuth.NewAnteHandler())

	// mount the multistore and load the latest state
	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyCheckpoint,
		app.keyStaking,
		app.keyBor,
		app.keyFeeCollection,
		app.keyParams,
		app.tKeyParams,
	)
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

	auth.RegisterCodec(cdc)
	// bank.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)

	// custom types
	checkpoint.RegisterWire(cdc)
	staking.RegisterWire(cdc)
	bor.RegisterWire(cdc)

	cdc.Seal()
	return cdc
}

// MakePulp creates pulp codec and registers custom types for decoder
func MakePulp() *authTypes.Pulp {
	pulp := authTypes.GetPulpInstance()

	// register custom type
	checkpoint.RegisterPulp(pulp)
	staking.RegisterPulp(pulp)
	bor.RegisterPulp(pulp)
	return pulp
}

// BeginBlocker executes before each block
func (app *HeimdallApp) beginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker executes on each end block
func (app *HeimdallApp) endBlocker(ctx sdk.Context, x abci.RequestEndBlock) abci.ResponseEndBlock {
	var valUpdates = make(map[hmTypes.ValidatorID]abci.ValidatorUpdate)
	if ctx.BlockHeader().NumTxs > 0 {
		// --- Start update to new validators
		currentValidatorSet := app.stakingKeeper.GetValidatorSet(ctx)
		currentValidatorSetCopy := currentValidatorSet.Copy()
		allValidators := app.stakingKeeper.GetAllValidators(ctx)
		ackCount := app.stakingKeeper.GetACKCount(ctx)

		// apply updates
		helper.UpdateValidators(
			&currentValidatorSet, // pointer to current validator set -- UpdateValidators will modify it
			allValidators,        // All validators
			ackCount,             // ack count
		)

		// update validator set in store
		err := app.stakingKeeper.UpdateValidatorSetInStore(ctx, currentValidatorSet)
		if err != nil {
			logger.Error("Unable to update validator set in state", "Error", err)
		} else {
			// remove all stale validators
			for _, validator := range currentValidatorSetCopy.Validators {
				val := abci.ValidatorUpdate{
					Power:  0,
					PubKey: validator.PubKey.ABCIPubKey(),
				}
				// validator update
				valUpdates[validator.ID] = val
			}
			// add new validators
			currentValidatorSet := app.stakingKeeper.GetValidatorSet(ctx)
			for _, validator := range currentValidatorSet.Validators {
				val := abci.ValidatorUpdate{
					Power:  int64(validator.Power),
					PubKey: validator.PubKey.ABCIPubKey(),
				}
				// validator update
				valUpdates[validator.ID] = val
			}
			// --- End update validators

			// check if ACK is present in cache
			if app.checkpointKeeper.GetCheckpointCache(ctx, checkpoint.CheckpointACKCacheKey) {
				logger.Info("Checkpoint ACK processed in block", "CheckpointACKProcessed", app.checkpointKeeper.GetCheckpointCache(ctx, checkpoint.CheckpointACKCacheKey))

				// --- Start update proposer

				// increment accum
				app.stakingKeeper.IncreamentAccum(ctx, 1)

				// log new proposer
				vs := app.stakingKeeper.GetValidatorSet(ctx)
				newProposer := vs.GetProposer()
				logger.Debug(
					"New proposer selected",
					"validator", newProposer.ID,
					"signer", newProposer.Signer.String(),
					"power", newProposer.Power,
				)
				// --- End update proposer

				// clear ACK cache
				app.checkpointKeeper.FlushACKCache(ctx)
			}

			// check if checkpoint is present in cache
			if app.checkpointKeeper.GetCheckpointCache(ctx, checkpoint.CheckpointCacheKey) {
				logger.Info("Checkpoint processed in block", "CheckpointProcessed", true)
				// collect and update sigs in span
				// Send Checkpoint to Rootchain
				PrepareAndSendCheckpoint(ctx, app.checkpointKeeper, app.stakingKeeper, app.caller)
				// clear Checkpoint cache
				app.checkpointKeeper.FlushCheckpointCache(ctx)
			}

			if app.borKeeper.GetSpanCache(ctx) {
				logger.Info("Propose Span processed in block", "ProposeSpanProcesses", true)
				// TODO Send proof to bor chain
				// app.borKeeper.AddSigs(ctx, ctx.BlockHeader().Votes)
				// get sigs from votes
				var votes []tmTypes.Vote
				err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
				if err != nil {
					logger.Error("Error while unmarshalling vote", "error", err)
				}
				sigs := helper.GetSigs(votes)
				fmt.Println("sigs", hex.EncodeToString(sigs))
				fmt.Println("vote", hex.EncodeToString(helper.GetVoteBytes(votes, ctx)))
				// flush span cache
				app.borKeeper.FlushSpanCache(ctx)
			}
		}
	}
	// convert updates from map to array
	var tmValUpdates []abci.ValidatorUpdate
	for _, v := range valUpdates {
		tmValUpdates = append(tmValUpdates, v)
	}

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		ValidatorUpdates: tmValUpdates,
	}
}

// initialize store from a genesis state
func (app *HeimdallApp) initFromGenesisState(ctx sdk.Context, genesisState GenesisState) []abci.ValidatorUpdate {
	genesisState.Sanitize()

	// accounts
	// Load the genesis accounts
	for _, genacc := range genesisState.Accounts {
		acc := app.accountKeeper.NewAccountWithAddress(ctx, genacc.Address[:])
		acc.SetCoins(genacc.Coins)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// check if genesis is actually a genesis
	var isGenesis bool
	if len(genesisState.StakingData.CurrentValSet.Validators) == 0 {
		isGenesis = true
	} else {
		isGenesis = false
	}

	// validator set
	// valSet, valUpdates := app.GetValidatorsFromGenesis(ctx, &genesisState, genesisState.AckCount)
	// if len(valSet.Validators) == 0 {
	// 	panic(errors.New("no valid validators found"))
	// }

	// var currentValSet hmTypes.ValidatorSet
	// if isGenesis {
	// 	currentValSet = valSet
	// } else {
	// 	currentValSet = genesisState.CurrentValSet
	// }

	// // TODO match valSet and genesisState.CurrentValSet for difference in accum
	// // update validator set in store
	// err = app.stakingKeeper.UpdateValidatorSetInStore(ctx, currentValSet)
	// if err != nil {
	// 	logger.Error("Unable to marshall validator set while adding in store", "Error", err)
	// 	panic(err)
	// }

	// set span duration from genesis
	// app.borKeeper.SetSpanDuration(ctx, genesisState.SpanDuration)

	// Set initial ack count
	// app.stakingKeeper.UpdateACKCountWithValue(ctx, genesisState.AckCount)

	// Add checkpoint in buffer
	// app.checkpointKeeper.SetCheckpointBuffer(ctx, genesisState.BufferedCheckpoint)

	// Set Caches
	// app.SetCaches(ctx, &genesisState)

	// Set last no-ack
	// app.checkpointKeeper.SetLastNoAck(ctx, genesisState.LastNoACK)

	// Add all headers
	// app.InsertHeaders(ctx, &genesisState)

	//
	// InitGenesis
	//
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	bor.InitGenesis(ctx, app.borKeeper, genesisState.BorData)
	checkpoint.InitGenesis(ctx, app.checkpointKeeper, genesisState.CheckpointData)
	staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)

	// validate genesis state
	if err := ValidateGenesisState(genesisState); err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	// increment accumulator if starting from genesis
	if isGenesis {
		app.stakingKeeper.IncreamentAccum(ctx, 1)
	}

	//
	// get val updates
	//

	var valUpdates []abci.ValidatorUpdate

	// check if validator is current validator
	// add to val updates else skip
	for _, validator := range genesisState.StakingData.Validators {
		if validator.IsCurrentValidator(genesisState.CheckpointData.AckCount) {
			// convert to Validator Update
			updateVal := abci.ValidatorUpdate{
				Power:  int64(validator.Power),
				PubKey: validator.PubKey.ABCIPubKey(),
			}
			// Add validator to validator updated to be processed below
			valUpdates = append(valUpdates, updateVal)
		}
	}
	return valUpdates
}

func (app *HeimdallApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	logger.Info("Loading validators from genesis and setting defaults")
	common.InitStakingLogger(&ctx)

	// get genesis state
	var genesisState GenesisState
	err := json.Unmarshal(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	// init state from genesis state
	valUpdates := app.initFromGenesisState(ctx, genesisState)

	//
	// draft reponse init chain
	//

	// TODO make sure old validtors dont go in validator updates ie deactivated validators have to be removed
	// udpate validators
	return abci.ResponseInitChain{
		// validator updates
		Validators: valUpdates,

		// consensus params
		ConsensusParams: &abci.ConsensusParams{
			Block: &abci.BlockParams{
				MaxBytes: maxBytesPerBlock,
				MaxGas:   maxGasPerBlock,
			},
			Evidence:  &abci.EvidenceParams{},
			Validator: &abci.ValidatorParams{PubKeyTypes: []string{ABCIPubKeyTypeSecp256k1}},
		},
	}
}

// // returns validator genesis/existing from genesis state
// // TODO add check from main chain if genesis information is right, else people can create invalid genesis and distribute
// func (app *HeimdallApp) GetValidatorsFromGenesis(ctx sdk.Context, genesisState *GenesisState, ackCount uint64) (newValSet hmTypes.ValidatorSet, valUpdates []abci.ValidatorUpdate) {
// 	if len(genesisState.GenValidators) > 0 {
// 		logger.Debug("Loading genesis validators")
// 		for _, validator := range genesisState.GenValidators {
// 			hmValidator := validator.HeimdallValidator()
// 			logger.Debug("gen validator", "gen", validator, "hmVal", hmValidator)
// 			if ok := hmValidator.ValidateBasic(); !ok {
// 				logger.Error("Invalid validator properties", "validator", hmValidator)
// 				return
// 			}
// 			if !hmValidator.IsCurrentValidator(genesisState.AckCount) {
// 				logger.Error("Genesis validators should be current validators", "FaultyValidator", hmValidator)
// 				return
// 			}
// 			if ok := newValSet.Add(&hmValidator); !ok {
// 				panic(errors.New("error while adding genesis validator"))
// 			} else {
// 				// Add individual validator to state
// 				app.stakingKeeper.AddValidator(ctx, hmValidator)
// 				// convert to Validator Update
// 				updateVal := abci.ValidatorUpdate{
// 					Power:  int64(validator.Power),
// 					PubKey: validator.PubKey.ABCIPubKey(),
// 				}

// 				// Add validator to validator updated to be processed below
// 				valUpdates = append(valUpdates, updateVal)
// 			}
// 		}
// 		logger.Debug("Adding validators to state", "ValidatorSet", newValSet, "ValUpdates", valUpdates)
// 		return
// 	}

// 	// read validators
// 	logger.Debug("Loading validators from state-dump")
// 	for _, validator := range genesisState.Validators {
// 		if !validator.ValidateBasic() {
// 			logger.Error("Invalid validator properties", "validator", validator)
// 			return
// 		}
// 		if ok := newValSet.Add(&validator); !ok {
// 			panic(errors.New("Error while addings new validator"))
// 		} else {
// 			// Add individual validator to state
// 			app.stakingKeeper.AddValidator(ctx, validator)

// 			// check if validator is current validator
// 			// add to val updates else skip
// 			if validator.IsCurrentValidator(ackCount) {
// 				// convert to Validator Update
// 				updateVal := abci.ValidatorUpdate{
// 					Power:  int64(validator.Power),
// 					PubKey: validator.PubKey.ABCIPubKey(),
// 				}
// 				// Add validator to validator updated to be processed below
// 				valUpdates = append(valUpdates, updateVal)
// 			}
// 		}
// 	}
// 	logger.Debug("Adding validators to state", "ValidatorSet", newValSet, "ValUpdates", valUpdates)
// 	return
// }

// // Set caches like checkpoint and checkpointACK cache
// // Incase user needs to retry sending last checkpoint or sending ACK
// func (app *HeimdallApp) SetCaches(ctx sdk.Context, genesisState *GenesisState) {
// 	if genesisState.CheckpointCache {
// 		logger.Debug("Found checkpoint cache", "CheckpointCache", genesisState.CheckpointCache)
// 		app.checkpointKeeper.SetCheckpointCache(ctx, checkpoint.DefaultValue)
// 		return
// 	}
// 	if genesisState.CheckpointACKCache {
// 		logger.Debug("Found checkpoint ACK cache", "CheckpointACKCache", genesisState.CheckpointACKCache)
// 		app.checkpointKeeper.SetCheckpointAckCache(ctx, checkpoint.DefaultValue)
// 		return
// 	}
// }

// // Insert headers into state
// func (app *HeimdallApp) InsertHeaders(ctx sdk.Context, genesisState *GenesisState) {
// 	if len(genesisState.Headers) != 0 {
// 		logger.Debug("Trying to add successfull checkpoints", "NoOfHeaders", len(genesisState.Headers))
// 		if int(genesisState.AckCount) != len(genesisState.Headers) {
// 			logger.Error("Number of headers and ack count do not match", "HeaderCount", len(genesisState.Headers), "AckCount", genesisState.AckCount)
// 			panic(errors.New("Incorrect state in state-dump , Please Check "))
// 		}
// 		for i, header := range genesisState.Headers {
// 			checkpointHeaderIndex := helper.GetConfig().ChildBlockInterval * (uint64(i) + 1)
// 			app.checkpointKeeper.AddCheckpoint(ctx, checkpointHeaderIndex, header)
// 		}
// 	}
// 	return
// }

// GetExtraData get extra data for checkpoint
func GetExtraData(ctx sdk.Context, _checkpoint hmTypes.CheckpointBlockHeader) []byte {
	logger.Debug("Creating extra data", "startBlock", _checkpoint.StartBlock, "endBlock", _checkpoint.EndBlock, "roothash", _checkpoint.RootHash, "timestamp", _checkpoint.TimeStamp)

	// craft a message
	// msg := checkpoint.NewMsgCheckpointBlock(
	// 	_checkpoint.Proposer,
	// 	_checkpoint.StartBlock,
	// 	_checkpoint.EndBlock,
	// 	_checkpoint.RootHash,
	// 	_checkpoint.TimeStamp,
	// )

	// helper.GetSignedTxBytes(ct)
	return nil
	// return txBytes[authTypes.PulpHashLength:]
}

// PrepareAndSendCheckpoint prepares all the data required for sending checkpoint and sends tx to rootchain
func PrepareAndSendCheckpoint(ctx sdk.Context, ck checkpoint.Keeper, sk staking.Keeper, caller helper.ContractCaller) {
	// fetch votes from block header
	var votes []tmTypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		logger.Error("Error while unmarshalling vote", "error", err)
	}

	// get sigs from votes
	sigs := helper.GetSigs(votes)

	// Getting latest checkpoint data from store using height as key and unmarshall
	_checkpoint, err := ck.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to unmarshall checkpoint from buffer while preparing checkpoint tx", "error", err, "height", ctx.BlockHeight())
		return
	}

	// Get extra data
	extraData := GetExtraData(ctx, _checkpoint)

	//fetch current child block from rootchain contract
	lastblock, err := caller.CurrentChildBlock()
	if err != nil {
		logger.Error("Could not fetch last block from mainchain", "error", err)
		panic(err)
	}

	// get validator address
	validatorAddress := ethCommon.BytesToAddress(helper.GetPubKey().Address().Bytes())

	// check if we are proposer
	if bytes.Equal(sk.GetCurrentProposer(ctx).Signer.Bytes(), validatorAddress.Bytes()) {
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
		logger.Info("We are not proposer", "proposer", sk.GetCurrentProposer(ctx), "validator", validatorAddress.String())
	}
}
