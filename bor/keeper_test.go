package bor_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	bortypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/merr"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type keeperTest struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
	k   bor.Keeper
}

func TestBorKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(keeperTest))
}

func (suite *keeperTest) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.k = suite.app.BorKeeper
}

func (suite *keeperTest) TestFreeze() {
	tc := []struct {
		id, startBlock uint64
		borChainID     string
		expErr         error
		msg            string
	}{
		{
			expErr: merr.ValErr{Field: "id", Module: bortypes.ModuleName},
			msg:    "validation error missing id",
		},
		{
			id:         1,
			startBlock: 3,
			msg:        "validation error missing id",
		},
		{
			id:  1,
			msg: "testing happy flow",
		},
	}

	for i, c := range tc {
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		err := suite.k.FreezeSet(suite.ctx, c.id, c.startBlock, c.borChainID)
		assert.Equal(suite.T(), c.expErr, err, cMsg)
	}
}

func (suite *keeperTest) TestSelectNextProducers() {
	// func (k *Keeper) SelectNextProducers(ctx sdk.Context) (vals []hmTypes.Validator, err error) {
	// return sdk.NewContext(app.deliverState.ms, header, false, app.logger)

}
