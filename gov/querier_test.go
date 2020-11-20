package gov_test

import (
	"strings"
	"testing"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/gov"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/stretchr/testify/suite"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/gov/types"
)

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	querier        sdk.Querier
}

// SetupTest setup all necessary things for querier tesing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.querier = gov.NewQuerier(suite.app.GovKeeper)
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

// TestQueryParams checks request query
func (suite *QuerierTestSuite) TestQueryParams() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryParams, types.ParamDeposit}, "/"),
		Data: []byte{},
	}
	_, err := querier(ctx, []string{types.QueryParams, types.ParamDeposit}, req)
	require.Nil(t, err)

	req = abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryParams, types.ParamVoting}, "/"),
		Data: []byte{},
	}
	_, err = querier(ctx, []string{types.QueryParams, types.ParamVoting}, req)
	require.Nil(t, err)

	req = abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryParams, types.ParamTallying}, "/"),
		Data: []byte{},
	}
	_, err = querier(ctx, []string{types.QueryParams, types.ParamTallying}, req)
	require.Nil(t, err)
}

func (suite *QuerierTestSuite) TestQueryProposal() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryProposal}, "/"),
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}
	bz, err := querier(ctx, []string{types.QueryProposal}, req)
	require.Nil(t, err)
	require.NotNil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryProposals() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryProposals}, "/"),
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}
	bz, err := querier(ctx, []string{types.QueryProposals}, req)
	require.Nil(t, err)
	require.NotNil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryDeposits() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryDeposits}, "/"),
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}
	bz, err := querier(ctx, []string{types.QueryDeposits}, req)
	require.Nil(t, err)
	require.NotNil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryDeposit() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryDeposit}, "/"),
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}
	bz, err := querier(ctx, []string{types.QueryDeposit}, req)
	require.Nil(t, err)
	require.NotNil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryVote() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryVote}, "/"),
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}
	bz, err := querier(ctx, []string{types.QueryVote}, req)
	require.Nil(t, err)
	require.NotNil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryVotes() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryVotes}, "/"),
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}
	bz, err := querier(ctx, []string{types.QueryVotes}, req)
	require.Nil(t, err)
	require.NotNil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryTally() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	req := abci.RequestQuery{
		Path: strings.Join([]string{"custom", types.QuerierRoute, types.QueryTally}, "/"),
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}
	bz, err := querier(ctx, []string{types.QueryTally}, req)
	require.Nil(t, err)
	require.NotNil(t, bz)
}