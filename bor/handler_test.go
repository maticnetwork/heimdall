package bor_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	borCommon "github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/common"
	hmTypes "github.com/maticnetwork/heimdall/types"
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
		spanDuration uint64
		msgID        uint64
		proposer     hmTypes.HeimdallAddress
		startBlock   uint64
		endBlock     uint64
		chainID      string
		seed         borCommon.Hash
		span         *hmTypes.Span
		out          sdk.Result
		msg          string
	}{
		{
			out:          sdk.Result{Events: sdk.EmptyEvents()},
			spanDuration: 2,
			span:         &hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, ChainID: "15001"},
			chainID:      "15001", // default chain id
			startBlock:   2,
			endBlock:     3,
			msgID:        2,
			msg:          "happy flow",
		},
		{
			out: common.ErrInvalidBorChainID("1").Result(),
			msg: "error invalid chain id",
		},
		{
			chainID: "15001", // default chain id
			out:     common.ErrSpanNotFound("1").Result(),
			msg:     "error span not found",
		},
		{
			span:       &hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, ChainID: "15001"},
			chainID:    "15001", // default chain id
			startBlock: 1,
			endBlock:   0,
			msgID:      2,
			out:        common.ErrSpanNotInCountinuity("1").Result(),
			msg:        "error span not in continuity",
		},
		{
			span:       &hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, ChainID: "15001"},
			chainID:    "15001", // default chain id
			startBlock: 0,
			endBlock:   1,
			msgID:      2,
			out:        common.ErrSpanNotInCountinuity("1").Result(),
			msg:        "error span not in continuity",
		},
		{
			span:       &hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, ChainID: "15001"},
			chainID:    "15001", // default chain id
			startBlock: 1,
			endBlock:   1,
			msgID:      0,
			out:        common.ErrSpanNotInCountinuity("1").Result(),
			msg:        "error span not in continuity",
		},
		{
			span:       &hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, ChainID: "15001"},
			chainID:    "15001", // default chain id
			startBlock: 1,
			endBlock:   1,
			msgID:      2,
			out:        common.ErrInvalidSpanDuration("1").Result(),
			msg:        "error span invalid duration",
		},
	}
	for i, c := range tc {
		suite.SetupTest()
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		if c.spanDuration != 0 {
			suite.app.BorKeeper.SetParams(suite.ctx, borTypes.Params{SpanDuration: c.spanDuration})
		}
		if c.span != nil {
			suite.app.BorKeeper.AddNewSpan(suite.ctx, *c.span)
		}

		// cSpan is used to check if span data remains constant post handler execution
		cSpan := suite.app.BorKeeper.GetAllSpans(suite.ctx)

		out := bor.HandleMsgProposeSpan(suite.ctx, borTypes.MsgProposeSpan{ID: c.msgID, Proposer: c.proposer, StartBlock: c.startBlock, EndBlock: c.endBlock, ChainID: c.chainID, Seed: c.seed}, suite.app.BorKeeper)
		suite.Equal(c.out, out, c.msg)

		// pSpan is used to check if span data remains constant post handler execution
		pSpan := suite.app.BorKeeper.GetAllSpans(suite.ctx)
		suite.Equal(cSpan, pSpan, "Invalid: handler should not update span "+c.msg)

	}
}
