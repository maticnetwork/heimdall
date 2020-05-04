package bor

import (
	"fmt"
	"testing"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank"
	chainmanager "github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/clerk"
	"github.com/maticnetwork/heimdall/params"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/supply"
	"github.com/maticnetwork/heimdall/topup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	ModuleBasics = module.NewBasicManager(
		params.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		supply.AppModuleBasic{},
		chainmanager.AppModuleBasic{},
		staking.AppModuleBasic{},
		checkpoint.AppModuleBasic{},
		AppModuleBasic{},
		clerk.AppModuleBasic{},
		topup.AppModuleBasic{},
		//gov.NewAppModuleBasic(paramsClient.ProposalHandler),
	)
)

type KeeperTestSuite struct {
	suite.Suite
	//	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	ModuleBasics.RegisterCodec(cdc)

	cdc.Seal()

	// app = app.Setup(isCheckTx)
	db := dbm.NewMemDB()

	// set prefix
	config := sdk.GetConfig()
	config.Seal()

	logger := log.NewNopLogger()

	pulp := authTypes.GetPulpInstance()
	bApp := bam.NewBaseApp("bor_keeper_test", logger, db, authTypes.RLPTxDecoder(cdc, pulp))

	isCheckTx := false
	suite.ctx = bApp.NewContext(isCheckTx, abci.Header{})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestFreeze() {
	tc := []struct {
		ctx        sdk.Context
		id         uint64
		startBlock uint64
		borChainID string
		expErr     error
		msg        string
	}{
		{
			ctx:        suite.ctx,
			id:         1,
			startBlock: 1,
			borChainID: "keeper test chain",
			msg:        "testing bor chain",
		},
	}

	// func (k *Keeper) FreezeSet(ctx sdk.Context, id uint64, startBlock uint64, borChainID string) error {

	for i, c := range tc {
		tMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		k := Keeper{}
		err := k.FreezeSet(c.ctx, c.id, c.startBlock, c.borChainID)
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
