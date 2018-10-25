package app

import (
	"os"
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/ethereum/go-ethereum/rlp"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	heimlib "github.com/maticnetwork/heimdall/libs"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	appName = "HeimdallApp"
)

type HeimdallApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the multistore
	keyMain       *sdk.KVStoreKey
	keyCheckpoint *sdk.KVStoreKey
	keyStake      *sdk.KVStoreKey

	keyStaker *sdk.KVStoreKey
	// manage getting and setting accounts
	checkpointKeeper checkpoint.Keeper
	stakerKeeper     staking.Keeper
}

var logger = heimlib.NewMainLogger(log.NewSyncWriter(os.Stdout)).With("module", "app")

func NewHeimdallApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *HeimdallApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create your application type
	var app = &HeimdallApp{
		cdc:           cdc,
		BaseApp:       bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		keyMain:       sdk.NewKVStoreKey("main"),
		keyCheckpoint: sdk.NewKVStoreKey("checkpoint"),
		keyStaker:     sdk.NewKVStoreKey("staker"),
	}

	//TODO change to its own codespace
	app.checkpointKeeper = checkpoint.NewKeeper(app.cdc, app.keyCheckpoint, app.RegisterCodespace(checkpoint.DefaultCodespace))
	app.stakerKeeper = staking.NewKeeper(app.cdc, app.keyStaker, app.RegisterCodespace(checkpoint.DefaultCodespace))
	// register message routes
	app.Router().AddRoute("checkpoint", checkpoint.NewHandler(app.checkpointKeeper))
	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	//app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.SetTxDecoder(app.txDecoder)
	//TODO check if correct
	app.BaseApp.SetTxDecoder(app.txDecoder)
	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyCheckpoint, app.keyStaker)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()

	return app
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	checkpoint.RegisterWire(cdc)
	// register custom type

	cdc.Seal()
	return cdc
}

func (app *HeimdallApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	// todo add flushValidatorSet here
	return abci.ResponseBeginBlock{}
}

func (app *HeimdallApp) EndBlocker(ctx sdk.Context, x abci.RequestEndBlock) abci.ResponseEndBlock {

	//validatorSet := staking.EndBlocker(ctx, app.stakerKeeper)
	//logger.Info("New Validator Set : %v", validatorSet)

	// unmarshall votes from header
	var votes []tmtypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		logger.Error("Unmarshalling Vote Errored", err, "Error")
	}

	// get sigs from votes
	var sigs []byte
	sigs = GetSigs(votes)

	if ctx.BlockHeader().NumTxs == 1 {

		// Getting latest checkpoint data from store using height as key and unmarshall
		var _checkpoint checkpoint.CheckpointBlockHeader
		json.Unmarshal(app.checkpointKeeper.GetCheckpoint(ctx, ctx.BlockHeight()), &_checkpoint)

		extraData := GetExtraData(_checkpoint, ctx)

		helper.SubmitProof(GetVoteBytes(votes, ctx), sigs, extraData, _checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash)
	}

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		//ValidatorUpdates: validatorSet,
	}
}

func GetSigs(votes []tmtypes.Vote) (sigs []byte) {
	// loop votes and append to sig to sigs
	for _, vote := range votes {
		sigs = append(sigs[:], vote.Signature[:]...)
	}

	return
}

func GetVoteBytes(votes []tmtypes.Vote, ctx sdk.Context) []byte {
	// sign bytes for vote
	return votes[0].SignBytes(ctx.ChainID())
}

func GetExtraData(_checkpoint checkpoint.CheckpointBlockHeader, ctx sdk.Context) []byte {
	msg := checkpoint.NewMsgCheckpointBlock(_checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash, _checkpoint.Proposer.String())

	tx := checkpoint.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		logger.Error("Error decoding transaction data ", err, "Error")
	}

	return txBytes
}

// RLP decodes the txBytes to a BaseTx
func (app *HeimdallApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {

	var tx = checkpoint.BaseTx{}
	err := rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		//todo create own error
		return nil, sdk.ErrTxDecode(err.Error())
	}

	return tx, nil
}

func (app *HeimdallApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	return abci.ResponseInitChain{}
}
func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	return appState, validators, err
}
