package bor_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type querierHandlerSuite struct {
	suite.Suite
	app        *app.HeimdallApp
	ctx        sdk.Context
	mockCaller mocks.IContractCaller
	querier    sdk.Querier
}

func TestBorQuerierHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(querierHandlerSuite))
}

func (suite *querierHandlerSuite) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.mockCaller = mocks.IContractCaller{}
	suite.querier = bor.NewQuerier(suite.app.BorKeeper, &suite.mockCaller)
}

func (suite *querierHandlerSuite) TestNewQueirer() {

	ethBlockData := `{"parentHash":"0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","miner":"0x0000000000000000000000000000000000000000","stateRoot":"0x5d6cded585e73c4e322c30c2f782a336316f17dd85a4863b9d838d2d4b8b3008","transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","difficulty":"0x2","number":"0x1","gasLimit":"0x9fd801","gasUsed":"0x0","timestamp":"0x5c530ffd","extraData":"0x506172697479205465636820417574686f7269747900000000000000000000002bbf886181970654ed46e3fae0ded41ee53fec702c47431988a7ae80e6576f3552684f069af80ba11d36327aaf846d470526e4a1c461601b2fd4ebdcdc2b734a01","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","hash":"0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a"}`
	// ethBlockHash := `0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a`
	//	ethBlockHash := `0xc3bd2d00745c03048a5616146a96f5ff78e54efb9e5b04af208cdaff6f3830ee`

	var ethHeader ethTypes.Header
	suite.Nil(json.Unmarshal([]byte(ethBlockData), &ethHeader))

	recordBytes, err := suite.app.Codec().MarshalJSON(types.QuerySpanParams{RecordID: 1})
	suite.Nil(err)

	querySpanParams, err := suite.app.Codec().MarshalJSON(hmTypes.QueryPaginationParams{Page: 1, Limit: 1})
	suite.Nil(err)

	tc := []struct {
		span                 *hmTypes.Span
		path                 []string
		req                  abci.RequestQuery
		callerMethod         []callerMethod
		expResp              []byte
		expErr               sdk.Error
		msg                  string
		loadVals, ignoreResp bool
	}{
		{
			path:   []string{"unknown_path"},
			expErr: sdk.ErrUnknownRequest("unknown auth query endpoint"),
			msg:    "error unknown request",
		},
		{

			expResp: []byte(`{"sprint_duration":64,"span_duration":6400,"producer_count":4}`),
			path:    []string{types.QueryParams},
			msg:     "happy flow for empty path query",
		},
		{
			expErr: sdk.ErrUnknownRequest("unknown_path is not a valid query request path"),
			path:   []string{types.QueryParams, "unknown_path"},
			msg:    "error invalid path for query",
			req:    abci.RequestQuery{Path: "unknown_path"},
		},
		{
			path:   []string{types.QuerySpan},
			req:    abci.RequestQuery{Data: []byte("invalid bytes")},
			expErr: sdk.ErrInternal("failed to parse params: invalid character 'i' looking for beginning of value"),
			msg:    "error invalid param data",
		},
		{
			path:   []string{types.QuerySpan},
			req:    abci.RequestQuery{Data: recordBytes},
			expErr: sdk.ErrInternal(sdk.AppendMsgToErr("could not get span", "span not found for id")),
			msg:    "error span not found",
		},
		{
			path:   []string{types.QuerySpanList},
			req:    abci.RequestQuery{Data: []byte("invalid bytes")},
			expErr: sdk.ErrInternal("failed to parse params: invalid character 'i' looking for beginning of value"),
			msg:    "error invalid param data query span list",
		},
		{
			path:    []string{types.QuerySpanList},
			span:    &hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, ChainID: "15001"},
			req:     abci.RequestQuery{Data: querySpanParams},
			expResp: []byte(`[{"span_id":1,"start_block":1,"end_block":1,"validator_set":{"validators":null,"proposer":null},"selected_producers":null,"bor_chain_id":"15001"}]`),
			msg:     "happy flow: query span list",
		},
		{
			path:    []string{types.QueryLatestSpan},
			span:    &hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, ChainID: "15001"},
			req:     abci.RequestQuery{Data: querySpanParams},
			expResp: []byte(`{"span_id":1,"start_block":1,"end_block":1,"validator_set":{"validators":null,"proposer":null},"selected_producers":null,"bor_chain_id":"15001"}`),
			msg:     "happy flow: query latest span",
		},
		{
			path:    []string{types.QueryLatestSpan},
			req:     abci.RequestQuery{Data: querySpanParams},
			expResp: []byte(`{"span_id":0,"start_block":0,"end_block":0,"validator_set":{"validators":null,"proposer":null},"selected_producers":null,"bor_chain_id":""}`),
			msg:     "happy flow: latest span default span",
		},
		{
			path:   []string{types.QueryNextProducers},
			req:    abci.RequestQuery{Data: querySpanParams},
			expErr: sdk.ErrInternal((sdk.AppendMsgToErr("cannot fetch next span seed from keeper", "error"))),
			callerMethod: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, fmt.Errorf("error")},
				},
			},
			msg: "error: query Next producers, cannot fetch next span",
		},
		{
			path:     []string{types.QueryNextProducers},
			span:     &hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, ChainID: "15001"},
			req:      abci.RequestQuery{Data: querySpanParams},
			loadVals: true,
			callerMethod: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethHeader, nil},
				},
			},
			ignoreResp: true,
			msg:        "happy flow: query Next producers",
		},
		{
			path:     []string{types.QueryNextSpanSeed},
			span:     &hmTypes.Span{ID: 1, StartBlock: 1, EndBlock: 1, ChainID: "15001"},
			req:      abci.RequestQuery{Data: querySpanParams},
			loadVals: true,
			callerMethod: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethHeader, nil},
				},
			},
			expResp: []byte(`"0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a"`),
			msg:     "happy flow: query Next span seed",
		},
		{
			path:   []string{types.QueryNextSpanSeed},
			expErr: sdk.ErrInternal(sdk.AppendMsgToErr("Error fetching next span seed", "error")),
			callerMethod: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, fmt.Errorf("error")},
				},
			},
			msg: "error: Next span seed not found",
		},
	}
	for i, c := range tc {
		suite.SetupTest()
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		if c.span != nil {
			suite.app.BorKeeper.AddNewSpan(suite.ctx, *c.span)
		}

		if c.callerMethod != nil {
			for _, m := range c.callerMethod {
				suite.mockCaller.On(m.name, m.args...).Return(m.ret...)
			}
		}
		if c.loadVals {
			simulation.LoadValidatorSet(1, suite.T(), suite.app.StakingKeeper, suite.ctx, false, 0)
		}

		out, err := suite.querier(suite.ctx, c.path, c.req)
		if !c.ignoreResp {
			suite.Equal(c.expResp, out, c.msg)
		}
		switch c.expErr {
		case nil:
			suite.Equal(c.expErr, err, c.msg)
		default:
			suite.Require().NotNil(err)
			suite.Equal(c.expErr.Error(), err.Error(), c.msg)
		}
	}
}
