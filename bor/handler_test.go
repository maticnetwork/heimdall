package bor_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type handlerSuite struct {
	suite.Suite
	app *app.HeimdallApp
	ctx sdk.Context
}

func TestBorHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

func (suite *handlerSuite) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
}

// func (suite *handlerSuite) TestNewHandler() {
// 	tc := []struct {
// 		k          bor.Keeper
// 		outHandler sdk.Handler
// 		msg        string
// 	}{
// 		{
// 			k:          suite.app.BorKeeper,
// 			outHandler: bor.NewHandler(suite.app.BorKeeper),
// 			msg:        "happy flow",
// 		},
// 	}
// 	for i, c := range tc {
// 		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
// 		out := bor.NewHandler(c.k)
// 		suite.IsType(sdk.Handler(suite.ctx, &suite.app.GetCaller()), out, c.msg)
// 		// suite.Equal(c.outHandler, out, c.msg)
// 	}
// }

func (suite handlerSuite) TestHandleMsgProposeSpan() {
	tc := []struct {
		out sdk.Result
		msg string
	}{
		{
			out: sdk.Result{},
			msg: "error invalid chain id",
		},
	}
	for i, c := range tc {
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out := bor.HandleMsgProposeSpan(suite.ctx, borTypes.MsgProposeSpan{ID: 1}, suite.app.BorKeeper)
		suite.Equal(c.out, out, c.msg)
	}
}
