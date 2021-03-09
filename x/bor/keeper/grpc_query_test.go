package keeper_test

import (
	"math/big"

	"github.com/maticnetwork/bor/common"

	"github.com/maticnetwork/heimdall/helper/mocks"

	chSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"

	ethTypes "github.com/maticnetwork/bor/core/types"

	ethereum "github.com/maticnetwork/bor"

	hmTypes "github.com/maticnetwork/heimdall/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/bor/keeper"
	"github.com/maticnetwork/heimdall/x/bor/types"
	borTypes "github.com/maticnetwork/heimdall/x/bor/types"
)

// callerMethod is to be used to mock the IContractCaller
type callerMethod struct {
	name string
	args []interface{}
	ret  []interface{}
}

func (suite *KeeperTestSuite) TestQueryParams() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	grpcQuery := keeper.Querier{
		Keeper: app.BorKeeper,
	}

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QueryParamsRequest
	}{
		{
			status: "success",
			error:  false,
			msg:    &borTypes.QueryParamsRequest{},
		},
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
	}

	for _, c := range tc {
		resp, err := grpcQuery.Params(sdk.WrapSDKContext(ctx), c.msg)
		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NotNil(t, resp)
			require.NoError(t, err)
			require.Equal(t, resp.SpanDuration, borTypes.DefaultParams().SpanDuration)
			require.Equal(t, resp.Sprint, borTypes.DefaultParams().SprintDuration)
			require.Equal(t, resp.ProducerCount, borTypes.DefaultParams().ProducerCount)
		}
	}
}

func (suite *KeeperTestSuite) TestQueryParam() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	grpcQuery := keeper.Querier{
		Keeper: app.BorKeeper,
	}

	defaultParams := borTypes.DefaultParams()

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QueryParamRequest
		resp   *borTypes.QueryParamResponse
	}{
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
		{
			status: "success",
			error:  false,
			msg: &borTypes.QueryParamRequest{
				ParamsType: keeper.ParamSpan,
			},
			resp: &borTypes.QueryParamResponse{Params: &borTypes.QueryParamResponse_SpanDuration{
				SpanDuration: defaultParams.SpanDuration,
			}},
		},
		{
			status: "success",
			error:  false,
			msg: &borTypes.QueryParamRequest{
				ParamsType: keeper.ParamSprint,
			},
			resp: &borTypes.QueryParamResponse{Params: &borTypes.QueryParamResponse_Sprint{
				Sprint: defaultParams.SprintDuration,
			}},
		},
		{
			status: "success",
			error:  false,
			msg: &borTypes.QueryParamRequest{
				ParamsType: keeper.ParamProducerCount,
			},
			resp: &borTypes.QueryParamResponse{Params: &borTypes.QueryParamResponse_ProducerCount{
				ProducerCount: defaultParams.ProducerCount,
			}},
		},
	}

	for _, c := range tc {
		resp, err := grpcQuery.Param(sdk.WrapSDKContext(ctx), c.msg)
		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NotNil(t, resp)
			require.NoError(t, err)
			require.Equal(t, resp, c.resp)
		}
	}
}

func (suite *KeeperTestSuite) TestQuerySpan() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	spanId := uint64(1)

	grpcQuery := keeper.Querier{
		Keeper: app.BorKeeper,
	}

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QuerySpanRequest
		resp   *borTypes.QuerySpanResponse
		span   hmTypes.Span
	}{
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
		{
			status: "span not found for id",
			error:  true,
			msg:    &borTypes.QuerySpanRequest{SpanId: spanId},
		},
		{
			status: "success",
			error:  false,
			msg:    &borTypes.QuerySpanRequest{SpanId: spanId},
			resp: &borTypes.QuerySpanResponse{Span: &hmTypes.Span{
				ID:         spanId,
				StartBlock: 1,
				EndBlock:   3,
				ChainId:    "15001",
			}},
			span: hmTypes.Span{
				ID:         spanId,
				StartBlock: 1,
				EndBlock:   3,
				ChainId:    "15001",
			},
		},
	}

	for _, c := range tc {
		if !c.error {
			err := app.BorKeeper.AddNewSpan(ctx, c.span)
			require.NoError(t, err)
		}
		resp, err := grpcQuery.Span(sdk.WrapSDKContext(ctx), c.msg)
		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, resp, c.resp)
		}
	}
}

func (suite *KeeperTestSuite) TestQuerySpanList() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	grpcQuery := keeper.Querier{
		Keeper: app.BorKeeper,
	}

	spanId := uint64(1)

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QuerySpanListRequest
		resp   *borTypes.QuerySpanListResponse
		span   hmTypes.Span
	}{
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
		{
			status: "no span list",
			error:  true,
			msg: &borTypes.QuerySpanListRequest{Pagination: &hmTypes.QueryPaginationParams{
				Page:  1,
				Limit: 10,
			}},
		},
		{
			status: "success",
			error:  false,
			msg: &borTypes.QuerySpanListRequest{Pagination: &hmTypes.QueryPaginationParams{
				Page:  1,
				Limit: 10,
			}},
			resp: &borTypes.QuerySpanListResponse{Spans: []*hmTypes.Span{{
				ID:         spanId,
				StartBlock: 1,
				EndBlock:   3,
				ChainId:    "15001",
			}}},
			span: hmTypes.Span{
				ID:         spanId,
				StartBlock: 1,
				EndBlock:   3,
				ChainId:    "15001",
			},
		},
	}

	for _, c := range tc {
		if !c.error {
			err := app.BorKeeper.AddNewSpan(ctx, c.span)
			require.NoError(t, err)
		}
		resp, err := grpcQuery.SpanList(sdk.WrapSDKContext(ctx), c.msg)
		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, resp, c.resp)
		}
	}
}

func (suite *KeeperTestSuite) TestQueryLatestSpan() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	grpcQuery := keeper.Querier{
		Keeper: app.BorKeeper,
	}

	spanId := uint64(1)

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QueryLatestSpanRequest
		resp   *borTypes.QueryLatestSpanResponse
		span   hmTypes.Span
	}{
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
		{
			status: "success",
			error:  false,
			msg:    &borTypes.QueryLatestSpanRequest{},
			resp: &borTypes.QueryLatestSpanResponse{Span: &hmTypes.Span{
				ID:         spanId,
				StartBlock: 1,
				EndBlock:   3,
				ChainId:    "15001",
			}},
			span: hmTypes.Span{
				ID:         spanId,
				StartBlock: 1,
				EndBlock:   3,
				ChainId:    "15001",
			},
		},
	}

	for _, c := range tc {
		if !c.error {
			err := app.BorKeeper.AddNewSpan(ctx, c.span)
			require.NoError(t, err)
		}
		resp, err := grpcQuery.LatestSpan(sdk.WrapSDKContext(ctx), c.msg)
		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, resp, c.resp)
		}
	}
}

func (suite *KeeperTestSuite) TestQueryNextProducers() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx

	grpcQuery := keeper.NewQueryServerImpl(initApp.BorKeeper, &suite.contractCaller)

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QueryNextProducersRequest
		resp   *borTypes.QueryNextProducersResponse
		span   hmTypes.Span
		cm     []callerMethod
	}{
		{
			status: "success",
			error:  false,
			msg:    &borTypes.QueryNextProducersRequest{},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
			},
			resp: &borTypes.QueryNextProducersResponse{NextProducers: []hmTypes.Validator{}},
		},
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
		{
			status: "not found",
			error:  true,
			msg:    &borTypes.QueryNextProducersRequest{},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{nil, ethereum.NotFound},
				},
			},
		},
	}

	for _, c := range tc {
		if c.cm != nil {
			for _, m := range c.cm {
				suite.contractCaller.On(m.name, m.args...).Return(m.ret...)
			}
		}
		if !c.error {
			params := types.DefaultParams()
			initApp.BorKeeper.SetParams(ctx, &params)
			_ = chSim.LoadValidatorSet(4, t, initApp.StakingKeeper, ctx, false, 10)
			initApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 1)
		}
		resp, err := grpcQuery.NextProducers(sdk.WrapSDKContext(ctx), c.msg)

		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
		}

		// reset the contractCaller
		suite.contractCaller = mocks.IContractCaller{}
	}
}

func (suite *KeeperTestSuite) TestQueryNextSpanSeed() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx

	ethHeader := &ethTypes.Header{}
	grpcQuery := keeper.NewQueryServerImpl(initApp.BorKeeper, &suite.contractCaller)

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QueryNextSpanSeedRequest
		resp   *borTypes.QueryNextSpanSeedResponse
		span   hmTypes.Span
		cm     []callerMethod
	}{
		{
			status: "success",
			error:  false,
			msg:    &borTypes.QueryNextSpanSeedRequest{},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{ethHeader, nil},
				},
			},
			resp: &borTypes.QueryNextSpanSeedResponse{
				NextSpanSeed: ethHeader.Hash().String(),
			},
		},
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
		{
			status: "not found",
			error:  true,
			msg:    &borTypes.QueryNextSpanSeedRequest{},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{nil, ethereum.NotFound},
				},
			},
		},
	}

	for _, c := range tc {
		if c.cm != nil {
			for _, m := range c.cm {
				suite.contractCaller.On(m.name, m.args...).Return(m.ret...)
			}
		}
		if !c.error {
			params := types.DefaultParams()
			initApp.BorKeeper.SetParams(ctx, &params)
			_ = chSim.LoadValidatorSet(4, t, initApp.StakingKeeper, ctx, false, 10)
			initApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 1)
		}
		resp, err := grpcQuery.NextSpanSeed(sdk.WrapSDKContext(ctx), c.msg)

		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, resp, c.resp)
		}

		// reset the contractCaller
		suite.contractCaller = mocks.IContractCaller{}
	}
}

func (suite *KeeperTestSuite) TestQueryPrepareNextSpan() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx

	ethHeader := &ethTypes.Header{}
	grpcQuery := keeper.NewQueryServerImpl(initApp.BorKeeper, &suite.contractCaller)

	tc := []struct {
		status string
		error  bool
		msg    *borTypes.QueryPrepareNextSpanRequest
		resp   *borTypes.QueryPrepareNextSpanResponse
		span   hmTypes.Span
		cm     []callerMethod
	}{
		{
			status: "success",
			error:  false,
			msg: &borTypes.QueryPrepareNextSpanRequest{
				SpanId:     1,
				ChainId:    "15001",
				StartBlock: 2,
				Proposer:   common.HexToAddress("1231231").String(),
			},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{ethHeader, nil},
				},
			},
			resp: &borTypes.QueryPrepareNextSpanResponse{},
		},
		{
			status: "invalid request",
			error:  true,
			msg:    nil,
		},
		{
			status: "not found",
			error:  true,
			msg:    &borTypes.QueryPrepareNextSpanRequest{},
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{nil, ethereum.NotFound},
				},
			},
		},
	}

	for _, c := range tc {
		var validatorSet hmTypes.ValidatorSet
		if c.cm != nil {
			for _, m := range c.cm {
				suite.contractCaller.On(m.name, m.args...).Return(m.ret...)
			}
		}
		if !c.error {
			params := types.DefaultParams()
			initApp.BorKeeper.SetParams(ctx, &params)
			validatorSet = chSim.LoadValidatorSet(4, t, initApp.StakingKeeper, ctx, false, 10)
			initApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 1)
		}
		resp, err := grpcQuery.PrepareNextSpan(sdk.WrapSDKContext(ctx), c.msg)

		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, resp.Span.StartBlock, c.msg.StartBlock)
			require.Equal(t, resp.Span.ChainId, c.msg.ChainId)
			require.Equal(t, resp.Span.ID, c.msg.SpanId)
			require.Equal(t, resp.Span.ValidatorSet, validatorSet)
		}

		// reset the contractCaller
		suite.contractCaller = mocks.IContractCaller{}
	}
}
