package app

import (
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	"github.com/maticnetwork/heimdall/supply"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
	"github.com/maticnetwork/heimdall/types"
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

var (
	// module account permissions
	maccPerms = map[string][]string{
		authTypes.FeeCollectorName: nil,
		// mint.ModuleName:           {supply.Minter},
		// staking.BondedPoolName:    {supply.Burner, supply.Staking},
		// staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		// gov.ModuleName:            {supply.Burner},
	}
)

// HeimdallApp main heimdall app
type HeimdallApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the multistore
	keyAccount    *sdk.KVStoreKey
	keyBank       *sdk.KVStoreKey
	keySupply     *sdk.KVStoreKey
	keyGov        *sdk.KVStoreKey
	keyCheckpoint *sdk.KVStoreKey
	keyStaking    *sdk.KVStoreKey
	keyBor        *sdk.KVStoreKey
	keyMain       *sdk.KVStoreKey
	keyParams     *sdk.KVStoreKey
	tKeyParams    *sdk.TransientStoreKey

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	supplyKeeper  supply.Keeper
	govKeeper     gov.Keeper
	paramsKeeper  params.Keeper

	checkpointKeeper checkpoint.Keeper
	stakingKeeper    staking.Keeper
	borKeeper        bor.Keeper

	// masterKeeper common.Keeper
	caller helper.ContractCaller

	//  total coins supply
	TotalCoinsSupply types.Coins
}

var logger = helper.Logger.With("module", "app")

// AckRetriever retriever
type AckRetriever struct {
	App *HeimdallApp
}

func (d AckRetriever) GetACKCount(ctx sdk.Context) uint64 {
	return d.App.checkpointKeeper.GetACKCount(ctx)
}

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
		keyBank:    sdk.NewKVStoreKey(bankTypes.StoreKey),
		keySupply:  sdk.NewKVStoreKey(supplyTypes.StoreKey),
		// keyGov:        sdk.NewKVStoreKey(gov.StoreKey),
		keyCheckpoint: sdk.NewKVStoreKey(checkpointTypes.StoreKey),
		keyStaking:    sdk.NewKVStoreKey(stakingTypes.StoreKey),
		keyBor:        sdk.NewKVStoreKey(borTypes.StoreKey),
		keyParams:     sdk.NewKVStoreKey(subspace.StoreKey),
		tKeyParams:    sdk.NewTransientStoreKey(subspace.TStoreKey),
	}

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.caller = contractCallerObj

	// define param keeper
	app.paramsKeeper = params.NewKeeper(cdc, app.keyParams, app.tKeyParams)

	// account keeper
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount, // target store
		app.paramsKeeper.Subspace(authTypes.DefaultParamspace),
		authTypes.ProtoBaseAccount, // prototype
	)

	// bank keeper
	app.bankKeeper = bank.NewBaseKeeper(
		app.cdc,
		app.keyBank, // target store
		app.paramsKeeper.Subspace(bankTypes.DefaultParamspace),
		bankTypes.DefaultCodespace,
		app.accountKeeper,
	)

	// bank keeper
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		app.keyBank, // target store
		app.paramsKeeper.Subspace(supplyTypes.DefaultParamspace),
		maccPerms,
		app.accountKeeper,
		app.bankKeeper,
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
		AckRetriever{App: app},
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
		app.caller,
	)

	// register message routes
	app.Router().
		AddRoute(bankTypes.RouterKey, bank.NewHandler(app.bankKeeper)).
		AddRoute(checkpointTypes.RouterKey, checkpoint.NewHandler(app.checkpointKeeper, &app.caller)).
		AddRoute(stakingTypes.RouterKey, staking.NewHandler(app.stakingKeeper, &app.caller)).
		AddRoute(borTypes.RouterKey, bor.NewHandler(app.borKeeper))

	// query routes
	app.QueryRouter().
		AddRoute(authTypes.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(bankTypes.QuerierRoute, bank.NewQuerier(app.bankKeeper)).
		AddRoute(supplyTypes.QuerierRoute, supply.NewQuerier(app.supplyKeeper)).
		AddRoute(checkpointTypes.QuerierRoute, checkpoint.NewQuerier(app.checkpointKeeper)).
		AddRoute(borTypes.QuerierRoute, bor.NewQuerier(app.borKeeper))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.beginBlocker)
	app.SetEndBlocker(app.endBlocker)
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	// mount the multistore and load the latest state
	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyBank,
		app.keySupply,
		app.keyCheckpoint,
		app.keyStaking,
		app.keyBor,
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

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)

	authTypes.RegisterCodec(cdc)
	bankTypes.RegisterCodec(cdc)
	supplyTypes.RegisterCodec(cdc)

	checkpoint.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	bor.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}

// MakePulp creates pulp codec and registers custom types for decoder
func MakePulp() *authTypes.Pulp {
	pulp := authTypes.GetPulpInstance()

	// register custom type
	bankTypes.RegisterPulp(pulp)
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
		ackCount := app.checkpointKeeper.GetACKCount(ctx)

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
				// PrepareAndSendCheckpoint(ctx, app.checkpointKeeper, app.stakingKeeper, app.caller)
				// clear Checkpoint cache
				app.checkpointKeeper.FlushCheckpointCache(ctx)
			}

			if app.borKeeper.GetSpanCache(ctx) {
				logger.Info("Propose Span processed in block", "ProposeSpanProcesses", true)
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
	// Load the genesis accounts
	for _, genacc := range genesisState.Accounts {
		acc := app.accountKeeper.NewAccountWithAddress(ctx, types.BytesToHeimdallAddress(genacc.Address.Bytes()))
		acc.SetCoins(genacc.Coins)
		acc.SetSequence(genacc.Sequence)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// check if genesis is actually a genesis
	var isGenesis bool
	if len(genesisState.StakingData.CurrentValSet.Validators) == 0 {
		isGenesis = true
	} else {
		isGenesis = false
	}

	//
	// InitGenesis
	//
	auth.InitGenesis(ctx, app.accountKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	supply.InitGenesis(ctx, app.supplyKeeper, app.accountKeeper, genesisState.SupplyData)
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
