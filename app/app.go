package app

import (
	"encoding/json"
	"fmt"

	"bytes"
	"encoding/hex"
	"github.com/basecoin/checkpoint"
	txHelper "github.com/basecoin/contracts"
	"github.com/basecoin/staker"
	"github.com/basecoin/staking"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/ethereum/go-ethereum/rlp"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName = "BasecoinApp"
)

// BasecoinApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type BasecoinApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the multistore
	keyMain       *sdk.KVStoreKey
	keyAccount    *sdk.KVStoreKey
	keyIBC        *sdk.KVStoreKey
	keySideBlock  *sdk.KVStoreKey
	keyCheckpoint *sdk.KVStoreKey
	keyStake      *sdk.KVStoreKey

	keyStaker *sdk.KVStoreKey
	// manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	checkpointKeeper    checkpoint.Keeper
	stakeKeeper         stake.Keeper
	stakerKeeper        staker.Keeper
}

// NewBasecoinApp returns a reference to a new BasecoinApp given a logger and
// database. Internally, a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewBasecoinApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *BasecoinApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create your application type
	var app = &BasecoinApp{
		cdc:           cdc,
		BaseApp:       bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		keyMain:       sdk.NewKVStoreKey("main"),
		keyAccount:    sdk.NewKVStoreKey("acc"),
		keyIBC:        sdk.NewKVStoreKey("ibc"),
		keySideBlock:  sdk.NewKVStoreKey("sideBlock"),
		keyCheckpoint: sdk.NewKVStoreKey("checkpoint"),
		keyStake:      sdk.NewKVStoreKey("stake"),
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

	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	//TODO change to its own codespace
	app.checkpointKeeper = checkpoint.NewKeeper(app.cdc, app.keyCheckpoint, app.RegisterCodespace(stake.DefaultCodespace))
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.stakerKeeper = staker.NewKeeper(app.cdc, app.keyStaker, app.RegisterCodespace(stake.DefaultCodespace))
	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("checkpoint", checkpoint.NewHandler(app.checkpointKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("staker", staker.NewHandler(app.stakerKeeper))
	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	//app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.SetTxDecoder(app.txDecoder)
	//TODO check if correct
	app.BaseApp.SetTxDecoder(app.txDecoder)
	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keySideBlock, app.keyCheckpoint, app.keyStake, app.keyStaker)
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
	bank.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	checkpoint.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	staker.RegisterWire(cdc)
	// register custom type

	//cdc.RegisterConcrete(&types.AppAccount{}, "basecoin/Account", nil)

	cdc.Seal()

	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *BasecoinApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *BasecoinApp) EndBlocker(ctx sdk.Context, x abci.RequestEndBlock) abci.ResponseEndBlock {
	//logger := ctx.Logger().With("module", "x/baseapp")

	//validatorSet := staker.EndBlocker(ctx, app.stakerKeeper)

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
		//ValidatorUpdates: validatorSet,
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
func (app *BasecoinApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
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
func (app *BasecoinApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
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
func (app *BasecoinApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	//ctx := app.NewContext(true, abci.Header{})
	//accounts := []*types.GenesisAccount{}
	//
	//appendAccountsFn := func(acc auth.Account) bool {
	//	account := &types.GenesisAccount{
	//		Address: acc.GetAddress(),
	//		Coins:   acc.GetCoins(),
	//	}
	//
	//	accounts = append(accounts, account)
	//	return false
	//}
	//
	//app.accountMapper.IterateAccounts(ctx, appendAccountsFn)
	//
	//genState := types.GenesisState{Accounts: accounts}
	//appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	//if err != nil {
	//	return nil, nil, err
	//}

	return appState, validators, err
}
