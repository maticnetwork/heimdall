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
	"github.com/maticnetwork/heimdall/checkpoint/simulation"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type keeperTest struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func TestBorKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(keeperTest))
}

func (suite *keeperTest) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
}

func (suite *keeperTest) TestFreeze() {
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
		cSpan := suite.app.BorKeeper.GetAllSpans(suite.ctx)
		cEthBlock := suite.app.BorKeeper.GetLastEthBlock(suite.ctx)

		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		err := suite.app.BorKeeper.FreezeSet(suite.ctx, c.id, c.startBlock, c.endBlock, c.borChainID, c.seed)
		suite.Equal(c.expErr, err, cMsg)
		if c.checkSpan {
			// pSpan is used to check if span data remains constant post handler execution
			pSpan := suite.app.BorKeeper.GetAllSpans(suite.ctx)
			suite.NotEqual(cSpan, pSpan, "Invalid: handler should update span "+c.msg)

			pEthBlock := suite.app.BorKeeper.GetLastEthBlock(suite.ctx)
			suite.NotEqual(cEthBlock, pEthBlock, "Invalid: handler should update span "+c.msg)

		}

	}
}

func (suite *keeperTest) TestSelectNextProducers() {
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
		vals := simulation.LoadValidatorSet(4, suite.T(), suite.app.StakingKeeper, suite.ctx, false, 0).Validators // load a random number of validators
		for _, v := range vals {
			c.spanEligibleVals = append(c.spanEligibleVals, *v)
		}
		out, err := bor.SelectNextProducers(c.hash, c.spanEligibleVals, c.producerCount)
		if c.expOut {
			suite.NotNil(out, c.msg)
		} else {
			suite.Equal([]uint64{}, out, c.msg)
		}
		suite.Equal(c.expErr, err, c.msg)
	}
}

func (suite *keeperTest) TestBorKeeperSelectNextProducers() {

	tc := []struct {
		msg           string
		producerCount uint64
		seed          common.Hash
		expErr        error
		out           hmTypes.ValidatorSet
	}{
		{
			producerCount: 10,
			seed:          common.HexToHash("testSeed"),
			msg:           "happy flow with very large validator set",
			out:           simulation.LoadValidatorSet(40, suite.T(), suite.app.StakingKeeper, suite.ctx, true, 0), // load a random number of validators
		},
		{
			producerCount: 1,
			seed:          common.HexToHash("testSeed"),
			msg:           "happy flow with greater validators then allowed",
			out:           simulation.LoadValidatorSet(4, suite.T(), suite.app.StakingKeeper, suite.ctx, true, 0), // load a random number of validators
		},
		{
			producerCount: 10,
			seed:          common.HexToHash("testSeed"),
			msg:           "happy flow with fewer validators than allowed",
			out:           simulation.LoadValidatorSet(4, suite.T(), suite.app.StakingKeeper, suite.ctx, true, 0), // load a random number of validators
		},
	}

	for i, c := range tc {
		suite.SetupTest()
		// cVals is used to check if validators are being modified during execution
		cVals := suite.app.StakingKeeper.GetValidatorSet(suite.ctx)
		suite.app.BorKeeper.SetParams(suite.ctx, bortypes.Params{SprintDuration: 1, SpanDuration: 1, ProducerCount: c.producerCount})
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out, err := suite.app.BorKeeper.SelectNextProducers(suite.ctx, c.seed)

		// pVals is used to check if validators are being modified during execution
		pVals := suite.app.StakingKeeper.GetValidatorSet(suite.ctx)
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

func (suite *keeperTest) TestGetAllSpans() {
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
		err := suite.app.BorKeeper.AddNewSpan(suite.ctx, *c.span)
		suite.Nil(err, c.msg)
		out := suite.app.BorKeeper.GetAllSpans(suite.ctx)
		suite.Equal([]*hmTypes.Span{c.span}, out, c.msg)
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
		suite.app.BorKeeper.SetLastEthBlock(suite.ctx, c.currentBlockHeight)
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out := suite.app.BorKeeper.GetLastEthBlock(suite.ctx)
		suite.Equal(c.expOut, out, cMsg)
	}
}
