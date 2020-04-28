package staking_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         context.CLIContext
	querier        sdk.Querier
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for querier tesing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.querier = staking.NewQuerier(suite.app.StakingKeeper)
}

// TestQuerierTestSuite
func TestQuerierTestSuite(t *testing.T) {
	suite.Run(t, new(QuerierTestSuite))
}

// TestInvalidQuery checks request query
func (suite *QuerierTestSuite) TestInvalidQuery() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, bz)

	bz, err = querier(ctx, []string{types.QuerierRoute}, req)
	require.Error(t, err)
	require.Nil(t, bz)
}

func (suite *QuerierTestSuite) TestHandleQueryCurrentValidatorSet() {

}

func (suite *QuerierTestSuite) TesthandleQuerySigner() {

}

func (suite *QuerierTestSuite) TesthandleQueryValidator() {

}

func (suite *QuerierTestSuite) TestHandleQueryValidatorStatus() {

}

func (suite *QuerierTestSuite) TestHandleQueryProposer() {

}

func (suite *QuerierTestSuite) TestHandleQueryCurrentProposer() {

}

func (suite *QuerierTestSuite) TestHandleQueryDividendAccount() {

}

func (suite *QuerierTestSuite) TestHandleDividendAccountRoot() {

}

func (suite *QuerierTestSuite) TestHandleQueryAccountProof() {

}

func (suite *QuerierTestSuite) TestHandleQueryVerifyAccountProof() {

}

func (suite *QuerierTestSuite) TestHandleQueryStakingSequence() {

}

func (suite *QuerierTestSuite) Test() {

}
