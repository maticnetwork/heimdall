package bor_test

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	bortypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/merr"
	"github.com/maticnetwork/heimdall/test"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type keeperTest struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
	k   *bor.Keeper
}

func TestBorKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(keeperTest))
}

func (suite *keeperTest) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.k = &suite.app.BorKeeper
}

func (suite *keeperTest) TestFreeze() {
	tc := []struct {
		id, startBlock, endBlock uint64
		seed                     common.Hash
		borChainID               string
		expErr                   error
		msg                      string
	}{
		{
			id:  1,
			msg: "testing happy flow",
		},
		{
			id:         1,
			startBlock: 3,
			msg:        "validation error missing id",
		},
		{
			expErr: merr.ValErr{Field: "id", Module: bortypes.ModuleName},
			msg:    "validation error missing id",
		},
	}

	for i, c := range tc {
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		err := suite.k.FreezeSet(suite.ctx, c.id, c.startBlock, c.endBlock, c.borChainID, c.seed)
		suite.Equal(c.expErr, err, cMsg)
	}
}

func (suite *keeperTest) TestSelectNextProducers() {
	validator := hmTypes.NewValidator(1, 1, 0, 100, 100, hmTypes.NewPubKey([]byte("testKey")), [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
	// hmTypes.NewDividendAccount(
	// 	hmTypes.NewDividendAccountID(uint64(validator.ID)),
	// 	big.NewInt(0).String(),
	// 	big.NewInt(0).String(),
	// )
	suite.k.SetParams(suite.ctx, bortypes.Params{SprintDuration: 1, SpanDuration: 1})
	suite.app.StakingKeeper.AddValidator(suite.ctx, *validator)

	tc := []struct {
		msg string
		// inSpan []hmtypes.Span
		seed   common.Hash
		expOut []hmTypes.Validator
		expErr error
	}{
		{
			msg:    "happy flow",
			expOut: []hmTypes.Validator{},
		},
	}

	for i, c := range tc {
		test.LoadValidatorSet(2, suite.T(), suite.app.StakingKeeper, suite.ctx, false, 10)
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out, err := suite.k.SelectNextProducers(suite.ctx, c.seed)
		suite.Equal(c.expOut, out, cMsg)
		suite.Equal(c.expErr, err, cMsg)
	}
}

func (suite *keeperTest) TestGetLastEthBlock() {
	tc := []struct {
		msg                string
		currentBlockHeight *big.Int
		expOut             *big.Int
	}{
		{
			msg:                "Happy flow",
			currentBlockHeight: big.NewInt(1),
			expOut:             big.NewInt(1),
		},
	}

	for i, c := range tc {
		suite.k.SetLastEthBlock(suite.ctx, c.currentBlockHeight)
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out := suite.k.GetLastEthBlock(suite.ctx)
		suite.Equal(c.expOut, out, cMsg)
	}
}
