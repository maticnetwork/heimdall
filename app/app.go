package app

import (
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
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
		BaseApp:       bam.NewBaseApp(AppName, logger, db, RLPTxDecoder(), baseAppOptions...),
		keyMain:       sdk.NewKVStoreKey("main"),
		keyCheckpoint: sdk.NewKVStoreKey("checkpoint"),
		keyStaker:     sdk.NewKVStoreKey("staker"),
	}

	app.checkpointKeeper = checkpoint.NewKeeper(app.cdc, app.keyCheckpoint, app.RegisterCodespace(checkpoint.DefaultCodespace))
	app.stakerKeeper = staking.NewKeeper(app.cdc, app.keyStaker, app.RegisterCodespace(checkpoint.DefaultCodespace))
	// register message routes
	app.Router().AddRoute("checkpoint", checkpoint.NewHandler(app.checkpointKeeper))
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
	// unmarshall votes from header
	var votes []tmtypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		logger.Error("Unmarshalling vote errored", "error", err)
	}

	// get sigs from votes
	sigs := getSigs(votes)
	if ctx.BlockHeader().NumTxs == 1 {
		// Getting latest checkpoint data from store using height as key and unmarshall
		var _checkpoint checkpoint.CheckpointBlockHeader
		err := json.Unmarshal(app.checkpointKeeper.GetCheckpoint(ctx, ctx.BlockHeight()), &_checkpoint)
		if err != nil {
			logger.Error("Unable to unmarshall checkpoint", "error", err)
		}

		// Get extra data3
		extraData := getExtraData(_checkpoint, ctx)

		logger.Debug("Validating last block from main chain", "lastBlock", helper.GetLastBlock(), "startBlock", _checkpoint.StartBlock)
		if helper.GetLastBlock() == _checkpoint.StartBlock {
			logger.Info("Valid checkpoint")
			helper.SendCheckpoint(GetVoteBytes(votes, ctx), sigs, extraData)
		} else {
			logger.Error("Start block does not match", "lastBlock", helper.GetLastBlock(), "startBlock", _checkpoint.StartBlock)
			// TODO panic ?
		}
		validators = staking.EndBlocker(ctx, app.stakerKeeper)
	}
	// TODO move this to above ie execute when checkpoint
	//if ctx.BlockHeight()%10 ==0 {
	//	logger.Error("Changing Validator set","Height",ctx.BlockHeight())
	//	validators = staking.EndBlocker(ctx,app.stakerKeeper)
	//}

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		ValidatorUpdates: validators,
	}
}

func getSigs(votes []tmtypes.Vote) (sigs []byte) {
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

func getExtraData(_checkpoint checkpoint.CheckpointBlockHeader, ctx sdk.Context) []byte {
	logger.Debug("Creating extra data", "startBlock", _checkpoint.StartBlock, "endBlock", _checkpoint.EndBlock, "roothash", _checkpoint.RootHash)
	msg := checkpoint.NewMsgCheckpointBlock(_checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash)

	tx := checkpoint.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		logger.Error("Error decoding transaction data", "error", err)
	}

	return txBytes
}

func (app *HeimdallApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	return abci.ResponseInitChain{}
}

func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	return appState, validators, err
}

// RLPTxDecoder decodes the txBytes to a BaseTx
func RLPTxDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = checkpoint.BaseTx{}
		err := rlp.DecodeBytes(txBytes, &tx)
		if err != nil {
			//todo create own error
			return nil, sdk.ErrTxDecode(err.Error())
		}

		return tx, nil

	}
}
