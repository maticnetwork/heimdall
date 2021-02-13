package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/bor/keeper"
	borTypes "github.com/maticnetwork/heimdall/x/bor/types"
	"github.com/maticnetwork/heimdall/x/checkpoint/simulation"

	"github.com/maticnetwork/bor/common"

	"github.com/maticnetwork/heimdall/helper/mocks"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/bor/test_helper"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context

	contractCaller mocks.IContractCaller
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestFreeze() {
	initApp, ctx := suite.app, suite.ctx

	tc := []struct {
		id, startBlock, endBlock uint64
		seed                     common.Hash
		borChainID               string
		expErr                   error
		msg                      string
		checkSpan                bool
	}{
		{
			id:        1,
			msg:       "testing happy flow",
			checkSpan: true,
		},
		{
			id:         1,
			startBlock: 3,
			msg:        "validation error missing id",
			checkSpan:  true,
		},
	}

	for i, c := range tc {
		// cSpan is used to check if span data remains constant post handler execution
		cSpan := initApp.BorKeeper.GetAllSpans(ctx)
		cEthBlock := initApp.BorKeeper.GetLastEthBlock(ctx)

		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		err := initApp.BorKeeper.FreezeSet(ctx, c.id, c.startBlock, c.endBlock, c.borChainID, c.seed)
		suite.Equal(c.expErr, err, cMsg)
		if c.checkSpan {
			// pSpan is used to check if span data remains constant post handler execution
			pSpan := initApp.BorKeeper.GetAllSpans(ctx)
			suite.NotEqual(cSpan, pSpan, "Invalid: handler should update span "+c.msg)

			pEthBlock := initApp.BorKeeper.GetLastEthBlock(ctx)
			suite.NotEqual(cEthBlock, pEthBlock, "Invalid: handler should update span "+c.msg)
		}
	}
}

func (suite *KeeperTestSuite) TestSelectNextProducers() {
	initApp, ctx := suite.app, suite.ctx

	tc := []struct {
		hash             common.Hash
		spanEligibleVals []hmTypes.Validator
		producerCount    uint64
		expOut           bool
		expErr           error
		msg              string
	}{
		{
			hash:          common.Hash([32]byte{1, 2, 3, 4}),
			msg:           "happy flow",
			producerCount: 4,
			expOut:        true,
		},
		{
			hash:          common.Hash([32]byte{1, 2, 3, 4}),
			msg:           "happy flow",
			producerCount: 1,
			expOut:        true,
		},
		{
			hash:          common.Hash([32]byte{1, 2, 3, 4}),
			msg:           "nil output check",
			producerCount: 0,
		},
	}

	for i, c := range tc {
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		vals := simulation.LoadValidatorSet(4, suite.T(), initApp.StakingKeeper, ctx, false, 0).Validators // load a random number of validators
		for _, v := range vals {
			c.spanEligibleVals = append(c.spanEligibleVals, *v)
		}
		out, err := keeper.SelectNextProducers(c.hash, c.spanEligibleVals, c.producerCount)
		if c.expOut {
			suite.NotNil(out, c.msg)
		} else {
			suite.Equal([]uint64{}, out, c.msg)
		}
		suite.Equal(c.expErr, err, c.msg)
	}
}

func (suite *KeeperTestSuite) TestBorKeeperSelectNextProducers() {
	initApp, ctx := suite.app, suite.ctx

	tc := []struct {
		msg           string
		producerCount uint64
		seed          common.Hash
		expErr        error
		valCount      int
		out           hmTypes.ValidatorSet
	}{
		{
			producerCount: 10,
			seed:          common.HexToHash("testSeed"),
			msg:           "happy flow with very large validator set",
			valCount:      40,
		},
		{
			producerCount: 1,
			seed:          common.HexToHash("testSeed"),
			msg:           "happy flow with greater validators then allowed",
			valCount:      4,
		},
		{
			producerCount: 10,
			seed:          common.HexToHash("testSeed"),
			msg:           "happy flow with fewer validators than allowed",
			valCount:      4,
		},
	}

	for i, c := range tc {
		//suite.SetupTest()
		c.out = simulation.LoadValidatorSet(c.valCount, suite.T(), initApp.StakingKeeper, ctx, true, 0) // load a random number of validators
		// cVals is used to check if validators are being modified during execution
		cVals := initApp.StakingKeeper.GetValidatorSet(ctx)
		initApp.BorKeeper.SetParams(ctx, &borTypes.Params{SprintDuration: 1, SpanDuration: 1, ProducerCount: c.producerCount})
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out, err := initApp.BorKeeper.SelectNextProducers(ctx, c.seed)

		// pVals is used to check if validators are being modified during execution
		pVals := initApp.StakingKeeper.GetValidatorSet(ctx)
		suite.Equal(cVals, pVals, "Checking validators do not get modified while selecting them "+c.msg)

		for _, v := range out {
			for _, va := range c.out.Validators {
				if va.ID == v.ID {
					va.VotingPower, va.ProposerPriority = v.VotingPower, v.ProposerPriority // modifying values in order to compare
					suite.Equal(*va, v, "invalid validator found: validator does not exist in validator set")
				}
			}
		}
		suite.GreaterOrEqual(int(c.producerCount), len(out), c.msg)
		suite.Equal(c.expErr, err, cMsg)
	}
}

func (suite *KeeperTestSuite) TestGetAllSpans() {
	initApp, ctx := suite.app, suite.ctx

	tc := []struct {
		span *hmTypes.Span
		msg  string
	}{
		{
			span: &hmTypes.Span{ID: 666, StartBlock: 1, EndBlock: 1},
			msg:  "happy flow",
		},
	}

	for i, c := range tc {
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		err := initApp.BorKeeper.AddNewSpan(ctx, *c.span)
		suite.Nil(err, c.msg)
		out := initApp.BorKeeper.GetAllSpans(ctx)
		suite.Equal([]*hmTypes.Span{c.span}, out, c.msg)
	}
}

func (suite *KeeperTestSuite) TestGetLastEthBlock() {
	initApp, ctx := suite.app, suite.ctx

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
		initApp.BorKeeper.SetLastEthBlock(ctx, c.currentBlockHeight)
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out := initApp.BorKeeper.GetLastEthBlock(ctx)
		suite.Equal(c.expOut, out, cMsg)
	}
}
