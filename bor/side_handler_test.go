package bor_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/maticnetwork/bor"
	borCommon "github.com/maticnetwork/bor/common"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	bortypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmCommon "github.com/tendermint/tendermint/libs/common"
)

type sideChHandlerSuite struct {
	suite.Suite
	app        *app.HeimdallApp
	ctx        sdk.Context
	mockCaller mocks.IContractCaller
}

func TestBorSideChHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(sideChHandlerSuite))
}

func (suite *sideChHandlerSuite) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.mockCaller = mocks.IContractCaller{}
}

func (suite *sideChHandlerSuite) TestSideHandleMsgSpan() {
	var bi *big.Int

	ethBlockData := `{"parentHash":"0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","miner":"0x0000000000000000000000000000000000000000","stateRoot":"0x5d6cded585e73c4e322c30c2f782a336316f17dd85a4863b9d838d2d4b8b3008","transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","difficulty":"0x2","number":"0x1","gasLimit":"0x9fd801","gasUsed":"0x0","timestamp":"0x5c530ffd","extraData":"0x506172697479205465636820417574686f7269747900000000000000000000002bbf886181970654ed46e3fae0ded41ee53fec702c47431988a7ae80e6576f3552684f069af80ba11d36327aaf846d470526e4a1c461601b2fd4ebdcdc2b734a01","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","hash":"0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a"}`
	// ethBlockHash := `0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a`
	ethBlockHash := `0xc3bd2d00745c03048a5616146a96f5ff78e54efb9e5b04af208cdaff6f3830ee`

	var ethHeader ethTypes.Header
	suite.Nil(json.Unmarshal([]byte(ethBlockData), &ethHeader))

	type callerMethod struct {
		name string
		args []interface{}
		ret  []interface{}
	}
	tc := []struct {
		out       abci.ResponseDeliverSideTx
		msg       string
		codespace string
		code      common.CodeType
		result    abci.SideTxResultType
		cm        []callerMethod
		seed      borCommon.Hash
		span      *hmTypes.Span
	}{
		{
			codespace: "1",
			code:      common.CodeInvalidMsg,
			msg:       "error mainchain error",
			result:    abci.SideTxResultType_Skip,
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, ethereum.NotFound},
				},
			},
		},
		{
			msg:       "error msg seed bytes failure",
			codespace: "1",
			code:      common.CodeInvalidMsg,
			result:    abci.SideTxResultType_Skip,
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
			},
		},
		{
			msg:       "error failed to GetMaticChainBlock",
			codespace: "1",
			code:      common.CodeInvalidMsg,
			seed:      borCommon.HexToHash(ethBlockHash),
			result:    abci.SideTxResultType_Skip,
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
		},
		{
			msg:       "error failed to lastSpan",
			codespace: "1",
			code:      common.CodeInvalidMsg,
			result:    abci.SideTxResultType_Skip,
			seed:      borCommon.HexToHash(ethBlockHash),
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
		},
		{
			msg:       "error failed to lastSpan validation",
			codespace: "1",
			code:      common.CodeInvalidMsg,
			result:    abci.SideTxResultType_Skip,
			seed:      borCommon.HexToHash(ethBlockHash),
			span:      &hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 0, ChainID: "15001"},
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
		},
		{
			msg:    "happy flow",
			seed:   borCommon.HexToHash(ethBlockHash),
			result: abci.SideTxResultType_Yes,
			span:   &hmTypes.Span{ID: 1, StartBlock: 0, EndBlock: 1, ChainID: "15001"},
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
		},
	}

	for i, c := range tc {
		suite.SetupTest()
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		if c.cm != nil {
			for _, m := range c.cm {
				suite.mockCaller.On(m.name, m.args...).Return(m.ret...)
			}
		}
		if c.span != nil {
			suite.app.BorKeeper.AddNewSpan(suite.ctx, *c.span)
		}

		// cSpan is used to check if span data remains constant post handler execution
		cSpan := suite.app.BorKeeper.GetAllSpans(suite.ctx)

		out := bor.SideHandleMsgSpan(suite.ctx, suite.app.BorKeeper, borTypes.MsgProposeSpan{Seed: c.seed}, &suite.mockCaller)
		// construct output
		c.out = abci.ResponseDeliverSideTx{Code: uint32(c.code), Codespace: c.codespace, Result: c.result}
		suite.Equal(c.out, out, c.msg)

		// pSpan is used to check if span data remains constant post handler execution
		pSpan := suite.app.BorKeeper.GetAllSpans(suite.ctx)
		suite.Equal(cSpan, pSpan, "Invalid: handler should not update span "+c.msg)
	}
}

func (suite *sideChHandlerSuite) TestPostHandleMsgEventSpan() {
	tc := []struct {
		msg         string
		spanMsg     borTypes.MsgProposeSpan
		result      abci.SideTxResultType
		span        *hmTypes.Span
		out         sdk.Result
		producerErr bool
	}{
		{
			msg: "error result check",
			out: common.ErrSideTxValidation(suite.app.BorKeeper.Codespace()).Result(),
		},
		{
			msg:     "error span already exists",
			spanMsg: borTypes.MsgProposeSpan{ID: 1},
			span:    &hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, ChainID: "15001"},
			result:  abci.SideTxResultType_Yes,
			out:     hmCommon.ErrOldTx(suite.app.BorKeeper.Codespace()).Result(),
		},
		{
			msg:         "error unable to freeze val",
			spanMsg:     borTypes.MsgProposeSpan{ID: 0, StartBlock: 0, EndBlock: 0, ChainID: "15001"},
			producerErr: true,
			result:      abci.SideTxResultType_Yes,
			out:         common.ErrUnableToFreezeValSet(suite.app.BorKeeper.Codespace()).Result(),
		},
		{
			msg:     "happy flow",
			spanMsg: borTypes.MsgProposeSpan{ID: 0, StartBlock: 1, EndBlock: 0, ChainID: "15001"},
			result:  abci.SideTxResultType_Yes,
			out:     sdk.Result{Events: sdk.Events{sdk.Event{Type: "propose-span", Attributes: []tmCommon.KVPair{tmCommon.KVPair{Key: []uint8{0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e}, Value: []uint8{0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x2d, 0x73, 0x70, 0x61, 0x6e}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}, tmCommon.KVPair{Key: []uint8{0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65}, Value: []uint8{0x62, 0x6f, 0x72}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}, tmCommon.KVPair{Key: []uint8{0x74, 0x78, 0x68, 0x61, 0x73, 0x68}, Value: []uint8{0x30, 0x78, 0x65, 0x33, 0x62, 0x30, 0x63, 0x34, 0x34, 0x32, 0x39, 0x38, 0x66, 0x63, 0x31, 0x63, 0x31, 0x34, 0x39, 0x61, 0x66, 0x62, 0x66, 0x34, 0x63, 0x38, 0x39, 0x39, 0x36, 0x66, 0x62, 0x39, 0x32, 0x34, 0x32, 0x37, 0x61, 0x65, 0x34, 0x31, 0x65, 0x34, 0x36, 0x34, 0x39, 0x62, 0x39, 0x33, 0x34, 0x63, 0x61, 0x34, 0x39, 0x35, 0x39, 0x39, 0x31, 0x62, 0x37, 0x38, 0x35, 0x32, 0x62, 0x38, 0x35, 0x35}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}, tmCommon.KVPair{Key: []uint8{0x73, 0x69, 0x64, 0x65, 0x2d, 0x74, 0x78, 0x2d, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74}, Value: []uint8{0x59, 0x65, 0x73}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}, tmCommon.KVPair{Key: []uint8{0x73, 0x70, 0x61, 0x6e, 0x2d, 0x69, 0x64}, Value: []uint8{0x30}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}, tmCommon.KVPair{Key: []uint8{0x73, 0x74, 0x61, 0x72, 0x74, 0x2d, 0x62, 0x6c, 0x6f, 0x63, 0x6b}, Value: []uint8{0x31}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}, tmCommon.KVPair{Key: []uint8{0x65, 0x6e, 0x64, 0x2d, 0x62, 0x6c, 0x6f, 0x63, 0x6b}, Value: []uint8{0x30}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}}, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: []uint8(nil), XXX_sizecache: 0}}},
		},
	}
	for i, c := range tc {
		suite.SetupTest()
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		if c.span != nil {
			suite.app.BorKeeper.AddNewSpan(suite.ctx, *c.span)
		}
		if c.producerErr {
			suite.app.BorKeeper.SetParams(suite.ctx, bortypes.Params{SprintDuration: 1, SpanDuration: 1, ProducerCount: 0})
		}

		out := bor.PostHandleMsgEventSpan(suite.ctx, suite.app.BorKeeper, c.spanMsg, c.result)
		suite.Equal(c.out, out, c.msg)
	}
}
