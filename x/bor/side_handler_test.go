package bor_test

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"testing"
	"time"

	ethereum "github.com/maticnetwork/bor"

	hmCommon "github.com/maticnetwork/heimdall/common"

	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"

	abci "github.com/tendermint/tendermint/abci/types"

	ethTypes "github.com/maticnetwork/bor/core/types"

	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/x/bor"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/bor/test_helper"
	borTypes "github.com/maticnetwork/heimdall/x/bor/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SideHandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	sideHandler    hmTypes.SideTxHandler
	postHandler    hmTypes.PostTxHandler
	contractCaller mocks.IContractCaller
	r              *rand.Rand
}

func (suite *SideHandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.sideHandler = bor.NewSideTxHandler(suite.app.BorKeeper, &suite.contractCaller)
	suite.postHandler = bor.NewPostTxHandler(suite.app.BorKeeper)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestSideHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SideHandlerTestSuite))
}

// Test Cases

func (suite *SideHandlerTestSuite) TestSideHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.sideHandler(ctx, nil)
	require.Equal(t, result.Code, sdkerrors.ErrUnknownRequest.ABCICode())
}

// callerMethod is to be used to mock the IContractCaller
type callerMethod struct {
	name string
	args []interface{}
	ret  []interface{}
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgProposeSpan() {
	t, ctx := suite.T(), suite.ctx
	var bi *big.Int
	ethBlockData := `{"parentHash":"0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","miner":"0x0000000000000000000000000000000000000000","stateRoot":"0x5d6cded585e73c4e322c30c2f782a336316f17dd85a4863b9d838d2d4b8b3008","transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","difficulty":"0x2","number":"0x1","gasLimit":"0x9fd801","gasUsed":"0x0","timestamp":"0x5c530ffd","extraData":"0x506172697479205465636820417574686f7269747900000000000000000000002bbf886181970654ed46e3fae0ded41ee53fec702c47431988a7ae80e6576f3552684f069af80ba11d36327aaf846d470526e4a1c461601b2fd4ebdcdc2b734a01","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","hash":"0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a"}`
	ethBlockHash := `0xc3bd2d00745c03048a5616146a96f5ff78e54efb9e5b04af208cdaff6f3830ee`

	var ethHeader ethTypes.Header
	suite.Nil(json.Unmarshal([]byte(ethBlockData), &ethHeader))

	tc := []struct {
		out       abci.ResponseDeliverSideTx
		msg       string
		codeSpace string
		result    tmprototypes.SideTxResultType
		cm        []callerMethod
		seed      string
		span      hmTypes.Span
		error     bool
		code      uint32
	}{

		{
			code:   hmCommon.ErrInvalidMsg.ABCICode(),
			msg:    "error mainchain error",
			result: tmprototypes.SideTxResultType_SKIP,
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{nil, ethereum.NotFound},
				},
			},
			error: true,
		},
		{
			msg:    "error msg seed bytes failure",
			code:   hmCommon.ErrInvalidMsg.ABCICode(),
			result: tmprototypes.SideTxResultType_SKIP,
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
			},
			error: true,
		},
		{
			msg:    "error failed to GetMaticChainBlock",
			code:   hmCommon.ErrInvalidMsg.ABCICode(),
			seed:   hmCommonTypes.HexToHeimdallHash(ethBlockHash).String(),
			result: tmprototypes.SideTxResultType_SKIP,
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
				{
					name: "GetMaticChainBlock",
					args: []interface{}{bi},
					ret:  []interface{}{&ethHeader, ethereum.NotFound},
				},
			},
			error: true,
		},
		{
			msg:    "error failed to lastSpan",
			code:   hmCommon.ErrInvalidMsg.ABCICode(),
			result: tmprototypes.SideTxResultType_SKIP,
			seed:   hmCommonTypes.HexToHeimdallHash(ethBlockHash).String(),
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
				{
					name: "GetMaticChainBlock",
					args: []interface{}{bi},
					ret:  []interface{}{&ethHeader, nil},
				},
			},
			error: true,
		},
		{
			msg:    "error failed to lastSpan validation",
			code:   hmCommon.ErrInvalidMsg.ABCICode(),
			result: tmprototypes.SideTxResultType_SKIP,
			seed:   hmCommonTypes.HexToHeimdallHash(ethBlockHash).String(),
			span:   hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, BorChainId: "15001"},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
				{
					name: "GetMaticChainBlock",
					args: []interface{}{bi},
					ret:  []interface{}{&ethHeader, nil},
				},
			},
			error: true,
		},
		{
			seed:   hmCommonTypes.HexToHeimdallHash(ethBlockHash).String(),
			result: tmprototypes.SideTxResultType_YES,
			span:   hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 1, BorChainId: "15001"},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
				{
					name: "GetMaticChainBlock",
					args: []interface{}{bi},
					ret:  []interface{}{&ethHeader, nil},
				},
			},
			error: false,
			code:  uint32(0),
		},
	}

	for _, c := range tc {
		if c.cm != nil {
			for _, m := range c.cm {
				suite.contractCaller.On(m.name, m.args...).Return(m.ret...)
			}
		}

		if c.code == uint32(0) {
			err := suite.app.BorKeeper.AddNewSpan(suite.ctx, c.span)
			require.NoError(t, err)
		}

		msg := borTypes.MsgProposeSpan{Seed: c.seed}
		result := suite.sideHandler(ctx, &msg)
		if c.error {
			require.Equal(t, c.code, result.Code, "Side tx handler should Fail")
			require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
			require.Equal(t, hmCommon.ErrInvalidMsg.ABCICode(), result.Code)
		} else {
			require.Equal(t, c.code, result.Code, "Side tx handler should be success")
			require.Equal(t, tmprototypes.SideTxResultType_YES, result.Result, "Result should be `yes`")
		}
		suite.contractCaller = mocks.IContractCaller{}
	}
}

// NewPostTxHandler

func (suite *SideHandlerTestSuite) TestPostTxHandler() {
	t, ctx := suite.T(), suite.ctx

	// post tx handler
	result, err := suite.postHandler(ctx, nil, tmprototypes.SideTxResultType_YES)
	require.Nil(t, result)
	require.Error(t, err)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgEventSpan() {
	t, ctx := suite.T(), suite.ctx
	tc := []struct {
		msg            string
		proposeSpanMsg borTypes.MsgProposeSpan
		result         tmprototypes.SideTxResultType
		span           hmTypes.Span
		event          string
		checkSpan      bool
		error          bool
	}{
		{
			msg:    "External call majority validation failed",
			error:  true,
			result: tmprototypes.SideTxResultType_NO,
		},
		{
			msg:            "Old txhash not allowed",
			error:          true,
			span:           hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, BorChainId: "15001"},
			result:         tmprototypes.SideTxResultType_YES,
			proposeSpanMsg: borTypes.MsgProposeSpan{SpanId: 1, StartBlock: 1, EndBlock: 0, BorChainId: "15001"},
		},
		{
			msg:            "success",
			error:          false,
			span:           hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, BorChainId: "15001"},
			result:         tmprototypes.SideTxResultType_YES,
			proposeSpanMsg: borTypes.MsgProposeSpan{SpanId: 0, StartBlock: 1, EndBlock: 0, BorChainId: "15001"},
		},
	}
	for _, c := range tc {
		if c.result == tmprototypes.SideTxResultType_YES {
			err := suite.app.BorKeeper.AddNewSpan(suite.ctx, c.span)
			require.NoError(t, err)
		}
		// cSpan is used to check if span data remains constant post handler execution
		cSpan, err := suite.app.BorKeeper.GetAllSpans(ctx)
		require.NoError(t, err)
		result, err := suite.postHandler(ctx, &c.proposeSpanMsg, c.result)
		if c.error {
			require.Equal(t, err.Error(), c.msg)
			require.Error(t, err)
			require.Nil(t, result)
		} else {
			pSpan, err := suite.app.BorKeeper.GetAllSpans(suite.ctx)
			require.NoError(t, err)
			suite.NotEqual(cSpan, pSpan)
			require.NotNil(t, result)
			require.NoError(t, err)
		}
	}
}
