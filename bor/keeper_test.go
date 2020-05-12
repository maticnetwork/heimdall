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
	"github.com/maticnetwork/heimdall/merr"
	"github.com/maticnetwork/heimdall/test"
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
	}{
		{
			id:  1,
			msg: "testing happy flow",
		},
		{
			id:         1,
			startBlock: 3,
			msg:        "validation error missing id",
		},
		{
			expErr: merr.ValErr{Field: "id", Module: bortypes.ModuleName},
			msg:    "validation error missing id",
		},
	}

	for i, c := range tc {
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		err := suite.app.BorKeeper.FreezeSet(suite.ctx, c.id, c.startBlock, c.endBlock, c.borChainID, c.seed)
		suite.Equal(c.expErr, err, cMsg)
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
		vals := test.LoadValidatorSet(4, suite.T(), suite.app.StakingKeeper, suite.ctx, false, 0).Validators // load a random number of validators
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
	}{
		{
			producerCount: 4,
			msg:           "happy flow",
		},
		{
			producerCount: 40,
			msg:           "happy flow",
		},
	}

	for i, c := range tc {
		suite.app.BorKeeper.SetParams(suite.ctx, bortypes.Params{SprintDuration: 1, SpanDuration: 1, ProducerCount: c.producerCount})
		test.LoadValidatorSet(4, suite.T(), suite.app.StakingKeeper, suite.ctx, true, 0) // load a random number of validators
		cMsg := fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		out, err := suite.app.BorKeeper.SelectNextProducers(suite.ctx, c.seed)
		suite.NotNil(out, cMsg)
		suite.Equal(c.expErr, err, cMsg)
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
