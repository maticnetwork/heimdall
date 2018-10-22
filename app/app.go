package app

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/ethereum/go-ethereum/rlp"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/checkpoint"
	txHelper "github.com/maticnetwork/heimdall/contracts"
	"github.com/maticnetwork/heimdall/staker"
)

const (
	appName = "HeimdallApp"
)

// HeimdallApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
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
	stakerKeeper     staker.Keeper
}

// NewHeimdallApp returns a reference to a new HeimdallApp given a logger and
// database. Internally, a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
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

	// define and attach the mappers and keepers
	//app.accountMapper = auth.NewAccountMapper(
	//	cdc,
	//	app.keyAccount, // target store
	//
	//	func() auth.Account {
	//		return &types.AppAccount{}
	//	},
	//)
	//TODO change to its own codespace
	app.checkpointKeeper = checkpoint.NewKeeper(app.cdc, app.keyCheckpoint, app.RegisterCodespace(checkpoint.DefaultCodespace))
	app.stakerKeeper = staker.NewKeeper(app.cdc, app.keyStaker, app.RegisterCodespace(checkpoint.DefaultCodespace))
	// register message routes
	app.Router().
		AddRoute("checkpoint", checkpoint.NewHandler(app.checkpointKeeper))
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

// MakeCodec creates a new wire codec and registers all the necessary types
// with the codec.
func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	checkpoint.RegisterWire(cdc)
	// register custom type

	cdc.Seal()

	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *HeimdallApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *HeimdallApp) EndBlocker(ctx sdk.Context, x abci.RequestEndBlock) abci.ResponseEndBlock {
	//logger := ctx.Logger().With("module", "x/baseapp")

	validatorSet := staker.EndBlocker(ctx, app.stakerKeeper)

	//logger.Info("New Validator Set : %v", validatorSet)

	var votes []tmtypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		fmt.Printf("error %v", err)
	}

	var sigs []byte
	sigs = GetSigs(votes)
	// TODO move this check to below check and validate checkpoint proposer
	if bytes.Equal(ctx.BlockHeader().Proposer.Address, txHelper.GetProposer().Bytes()) {
		fmt.Printf("Current Proposer and Block Proposer Matched ! ")
	} else {
		fmt.Printf("Current Proposer :%v , BlockProposer:  %v", txHelper.GetProposer().String(), ctx.BlockHeader().Proposer)
	}

	if ctx.BlockHeader().NumTxs == 1 {
		// Getting latest checkpoint data from store using height as key
		fmt.Printf(" Get Vote Bytes %v", hex.EncodeToString(GetVoteBytes(votes, ctx)))
		var _checkpoint checkpoint.CheckpointBlockHeader
		json.Unmarshal(app.checkpointKeeper.GetCheckpoint(ctx, ctx.BlockHeight()), &_checkpoint)
		extraData := GetExtraData(_checkpoint)
		txHelper.SubmitProof(GetVoteBytes(votes, ctx), sigs, extraData, _checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash)
		//txHelper.SendCheckpoint(int(_checkpoint.StartBlock), int(_checkpoint.EndBlock), sigs)
	}
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorSet,
	}
}

func GetSigs(votes []tmtypes.Vote) (sigs []byte) {
	for _, vote := range votes {
		sigs = append(sigs[:], vote.Signature[:]...)
	}
	return
}
func GetVoteBytes(votes []tmtypes.Vote, ctx sdk.Context) []byte {
	return votes[0].SignBytes(ctx.ChainID())
}
func GetExtraData(_checkpoint checkpoint.CheckpointBlockHeader) []byte {
	msg := checkpoint.NewMsgCheckpointBlock(_checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash, _checkpoint.Proposer.String())

	tx := checkpoint.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		fmt.Printf("Error generating TXBYtes %v", err)
	}
	return txBytes
}

// RLP decodes the txBytes to a BaseTx
func (app *HeimdallApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = checkpoint.BaseTx{}
	fmt.Printf("Decoding Transaction from app.go")
	err := rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode(err.Error())
	}
	return tx, nil
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *HeimdallApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	//genesisState := new(types.GenesisState)
	fmt.Printf("app state bytes is %v", string(stateJSON))

	//err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	//if err != nil {
	//	// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
	//
	//	panic(err)
	//}

	//for _, gacc := range genesisState.Accounts {
	//	acc, err := gacc.ToAppAccount()
	//	if err != nil {
	//		// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
	//		panic(err)
	//	}
	//
	//	acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
	//	app.accountMapper.SetAccount(ctx, acc)
	//}

	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *HeimdallApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	return appState, validators, err
}
