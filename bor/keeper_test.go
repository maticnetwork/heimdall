package bor

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	chainmanager "github.com/maticnetwork/heimdall/chainmanager"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/params"
	"github.com/maticnetwork/heimdall/params/subspace"
	paramsTypes "github.com/maticnetwork/heimdall/params/types"
	"github.com/maticnetwork/heimdall/staking"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	moduleBasics = module.NewBasicManager(
		AppModuleBasic{},
		staking.AppModuleBasic{},
		chainmanager.AppModuleBasic{},
		bank.AppModuleBasic{},
	)
)

type KeeperTestSuite struct {
	suite.Suite
	Keeper

	// dependencies
	mm               *module.Manager
	bApp             *baseapp.BaseApp
	ctx              sdk.Context
	cdc              *codec.Codec
	accountKeeper    auth.AccountKeeper
	chanKeeper       chainmanager.Keeper
	stakingKeeper    staking.Keeper
	bankKeeper       bank.Keeper
	checkpointKeeper checkpoint.Keeper
	keys             map[string]*sdk.KVStoreKey
	tkeys            map[string]*sdk.TransientStoreKey
	subspaces        map[string]subspace.Subspace
	caller           helper.ContractCaller
}

// GetACKCount returns ack count
func (k KeeperTestSuite) GetACKCount(ctx sdk.Context) uint64 {
	return k.checkpointKeeper.GetACKCount(ctx)
}

// IsCurrentValidatorByAddress check if validator is current validator
func (k KeeperTestSuite) IsCurrentValidatorByAddress(ctx sdk.Context, address []byte) bool {
	return k.stakingKeeper.IsCurrentValidatorByAddress(ctx, address)
}

// AddFeeToDividendAccount add fee to dividend account
func (k KeeperTestSuite) AddFeeToDividendAccount(ctx sdk.Context, valID types.ValidatorID, fee *big.Int) sdk.Error {
	return k.stakingKeeper.AddFeeToDividendAccount(ctx, valID, fee)
}

// GetValidatorFromValID get validator from validator id
func (k KeeperTestSuite) GetValidatorFromValID(ctx sdk.Context, valID types.ValidatorID) (validator types.Validator, ok bool) {
	return k.stakingKeeper.GetValidatorFromValID(ctx, valID)
}

// SetCoins sets coins
func (k KeeperTestSuite) SetCoins(ctx sdk.Context, addr types.HeimdallAddress, amt sdk.Coins) sdk.Error {
	return k.bankKeeper.SetCoins(ctx, addr, amt)
}

// GetCoins gets coins
func (k KeeperTestSuite) GetCoins(ctx sdk.Context, addr types.HeimdallAddress) sdk.Coins {
	return k.bankKeeper.GetCoins(ctx, addr)
}

// SendCoins transfers coins
func (k KeeperTestSuite) SendCoins(ctx sdk.Context, fromAddr types.HeimdallAddress, toAddr types.HeimdallAddress, amt sdk.Coins) sdk.Error {
	return k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.cdc = codec.New()

	codec.RegisterCrypto(suite.cdc)
	sdk.RegisterCodec(suite.cdc)
	moduleBasics.RegisterCodec(suite.cdc)

	suite.cdc.Seal()

	// app = app.Setup(isCheckTx)
	db := dbm.NewMemDB()

	// set prefix
	config := sdk.GetConfig()
	config.Seal()

	logger := log.NewNopLogger()

	pulp := authTypes.GetPulpInstance()
	suite.bApp = bam.NewBaseApp("bor_keeper_test", logger, db, authTypes.RLPTxDecoder(suite.cdc, pulp))
	suite.bApp.SetCommitMultiStoreTracer(nil)
	suite.bApp.SetAppVersion(version.Version)

	// keys
	suite.keys = sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		bankTypes.StoreKey,
		borTypes.StoreKey,
		authTypes.StoreKey,
		checkpointTypes.StoreKey,
		paramsTypes.StoreKey,
	)

	suite.tkeys = sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)
	paramsKeeper := params.NewKeeper(suite.cdc, suite.keys[paramsTypes.StoreKey], suite.tkeys[paramsTypes.TStoreKey], paramsTypes.DefaultCodespace)
	suite.subspaces = map[string]subspace.Subspace{}
	suite.subspaces[authTypes.ModuleName] = paramsKeeper.Subspace(authTypes.DefaultParamspace)
	suite.subspaces[bankTypes.ModuleName] = paramsKeeper.Subspace(bankTypes.DefaultParamspace)
	suite.subspaces[borTypes.ModuleName] = paramsKeeper.Subspace(borTypes.DefaultParamspace)
	suite.subspaces[checkpointTypes.ModuleName] = paramsKeeper.Subspace(checkpointTypes.DefaultParamspace)
	suite.subspaces[chainmanagerTypes.ModuleName] = paramsKeeper.Subspace(chainmanagerTypes.DefaultParamspace)
	suite.subspaces[stakingTypes.ModuleName] = paramsKeeper.Subspace(stakingTypes.DefaultParamspace)

	isCheckTx := false // NOTE using this as a placeholder, it is used while generating the context and would like to emulate that behaviour
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := moduleBasics.DefaultGenesis()
		stateBytes, err := codec.MarshalJSONIndent(suite.cdc, genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		suite.bApp.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}
	var err error
	suite.caller, err = helper.NewContractCaller()
	if err != nil {
		panic(err)
	}

	// create chain keeper
	// account keeper
	suite.accountKeeper = auth.NewAccountKeeper(
		suite.cdc,
		suite.keys[authTypes.StoreKey],
		suite.subspaces[authTypes.ModuleName],
		authTypes.ProtoBaseAccount,
	)
	suite.chainKeeper = chainmanager.NewKeeper(
		suite.cdc,
		suite.keys[chainmanagerTypes.StoreKey],
		suite.subspaces[chainmanagerTypes.ModuleName],
		common.DefaultCodespace,
		suite.caller,
	)
	suite.bankKeeper = bank.NewKeeper(
		suite.cdc,
		suite.keys[bankTypes.StoreKey],
		suite.subspaces[bankTypes.ModuleName],
		bankTypes.DefaultCodespace,
		suite.accountKeeper,
		suite,
	)
	suite.stakingKeeper = staking.NewKeeper(
		suite.cdc,
		suite.keys[stakingTypes.StoreKey],
		suite.subspaces[stakingTypes.ModuleName],
		common.DefaultCodespace,
		suite.chainKeeper,
		suite,
	)
	suite.Keeper = NewKeeper(
		suite.cdc,
		suite.keys[borTypes.StoreKey],
		suite.subspaces[borTypes.ModuleName],
		common.DefaultCodespace,
		suite.chainKeeper,
		suite.stakingKeeper,
		suite.caller,
	)

	suite.ctx = suite.bApp.NewContext(isCheckTx, abci.Header{})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestFreezeSet() {
	tc := []struct {
		ctx          sdk.Context
		id           uint64
		startBlock   uint64
		borChainID   string
		keeperParams borTypes.Params
		expErr       error
		msg          string
	}{
		{
			ctx:          suite.ctx,
			id:           1,
			startBlock:   0,
			borChainID:   "keeper test chain",
			keeperParams: borTypes.Params{SprintDuration: 100, SpanDuration: 10, ProducerCount: 10},
			msg:          "testing bor chain",
		},
	}

	// func (k *Keeper) FreezeSet(ctx sdk.Context, id uint64, startBlock uint64, borChainID string) error {

	for i, c := range tc {
		tMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)

		// set keeper parameters
		suite.Keeper.SetParams(c.ctx, c.keeperParams)
		fmt.Println(suite.Keeper.GetParams(c.ctx))

		err := suite.Keeper.FreezeSet(c.ctx, c.id, c.startBlock, c.borChainID)
		assert.Equal(suite.T(), c.expErr, err, tMsg)
	}

}

// Generic Tests

func TestCodeSpace(t *testing.T) {
	var testCodeSpaceType sdk.CodespaceType = "test codespace"
	k := Keeper{codespace: testCodeSpaceType}
	assert.Equal(t, testCodeSpaceType, k.Codespace(), "testing for codespacetype")
}

// func TestAddNewSpan(t *testing.T) {
// 	valSet := []*types.Validator{}
// 	tc := []struct {
// 		ctx    sdk.Context
// 		span   hmTypes.Span
// 		expErr error
// 		msg    string
// 	}{
// 		{
// 			ctx:  sdk.Context{},
// 			span: hmTypes.NewSpan(1, 1, 1, *hmTypes.NewValidatorSet(valSet), []types.Validator{}, "test chain"),
// 			msg:  "nil ctx failure",
// 		},
// 	}
//
// 	for i, c := range tc {
// 		k := Keeper{}
// 		assert.Equal(t, c.expErr, k.AddNewSpan(c.ctx, c.span), fmt.Sprintf("i: %v, msg: %v", i, c.msg))
// 	}
// }

func TestGetSpanKey(t *testing.T) {
	tc := []struct {
		in  uint64
		out []byte
		msg string
	}{
		{
			in:  1,
			out: []byte{0x36, 0x31},
			msg: "testing for small uint64",
		},
		{
			in:  1234567890123457890,
			out: []byte{0x36, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x37, 0x38, 0x39, 0x30},

			msg: "testing for large uint64",
		},
	}
	for i, c := range tc {
		b := GetSpanKey(c.in)
		assert.Equal(t, c.out, b, fmt.Sprintf("i: %v, msg: %v", i, c.msg))
	}
}
