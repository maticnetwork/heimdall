package bor_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/bor"
	borTypes "github.com/maticnetwork/heimdall/x/bor/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/bor/test_helper"
	"github.com/stretchr/testify/suite"
)

type handlerSuite struct {
	suite.Suite
	app     *app.HeimdallApp
	ctx     sdk.Context
	handler sdk.Handler
	r       *rand.Rand
}

func TestBorHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

func (suite *handlerSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(true)
	suite.handler = bor.NewHandler(suite.app.BorKeeper)
	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func (suite *handlerSuite) TestHandler() {
	t, ctx := suite.T(), suite.ctx
	result, err := suite.handler(ctx, nil)
	require.Nil(t, result)
	require.Error(t, err)
}

func (suite *handlerSuite) TestHandleMsgProposeSpan() {
	t, ctx := suite.T(), suite.ctx
	_, _, addr1 := testdata.KeyTestPubAddr()
	tc := []struct {
		msg          string
		span         hmTypes.Span
		spanDuration uint64
		spanId       uint64
		proposer     string
		startBlock   uint64
		endBlock     uint64
		chainID      string
		seed         string
		event        string
		error        bool
	}{
		{
			msg: "success",
			span: hmTypes.Span{
				ID:         2,
				StartBlock: 1,
				EndBlock:   2,
				BorChainId: "15001",
			},
			spanId:       3,
			proposer:     addr1.String(),
			event:        "propose-span",
			startBlock:   3,
			endBlock:     5,
			chainID:      "15001",
			spanDuration: 3,
			seed:         hmCommonTypes.HexToHeimdallHash("123123123").String(),
			error:        false,
		},
		{
			msg:   "Invalid Bor chain id",
			error: true,
		},
		{
			chainID: "15001",
			msg:     "Span not continuous",
			error:   true,
		},
		{
			span:       hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, BorChainId: "15001"},
			chainID:    "15001", // default chain id
			startBlock: 1,
			endBlock:   0,
			spanId:     2,
			msg:        "Span not continuous",
			error:      true,
		},
		{
			span:       hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, BorChainId: "15001"},
			chainID:    "15001", // default chain id
			startBlock: 0,
			endBlock:   1,
			spanId:     2,
			msg:        "Span not continuous",
			error:      true,
		},
		{
			span:       hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, BorChainId: "15001"},
			chainID:    "15001",
			startBlock: 1,
			endBlock:   1,
			spanId:     0,
			msg:        "Span not continuous",
			error:      true,
		},
		{
			msg: "Wrong span duration",
			span: hmTypes.Span{
				ID:         2,
				StartBlock: 1,
				EndBlock:   2,
				BorChainId: "15001",
			},
			spanId:       3,
			proposer:     addr1.String(),
			event:        "propose-span",
			startBlock:   3,
			endBlock:     5,
			chainID:      "15001",
			spanDuration: 5,
			seed:         hmCommonTypes.HexToHeimdallHash("123123123").String(),
			error:        true,
		},
	}

	for _, c := range tc {
		if c.spanDuration != 0 {
			suite.app.BorKeeper.SetParams(ctx, &borTypes.Params{SpanDuration: c.spanDuration, SprintDuration: 1, ProducerCount: 2})
		}
		if len(c.span.BorChainId) != 0 {
			err := suite.app.BorKeeper.AddNewSpan(ctx, c.span)
			require.NoError(t, err)
		}
		msg := borTypes.NewMsgProposeSpan(
			c.spanId,
			c.proposer,
			c.spanId,
			c.endBlock,
			c.chainID,
			c.seed,
		)
		result, err := suite.handler(ctx, &msg)
		if c.error {
			require.Error(t, err)
			require.Equal(t, err.Error(), c.msg)
			require.Nil(t, result)
		} else {
			require.NotNil(t, result)
			require.Equal(t, c.event, result.GetEvents()[0].Type)
			require.NoError(t, err)
		}
	}
}
