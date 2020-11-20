package gov_test

import (
	// "fmt"
	"time"
	"testing"

	"github.com/maticnetwork/heimdall/app"
	paramTypes "github.com/maticnetwork/heimdall/params/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

//
// Test suite
//

// ProposalTestSuite integrate test suite context object
type ProposalTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *ProposalTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestProposalTestSuite(t *testing.T) {
	suite.Run(t, new(ProposalTestSuite))
}

func testProposal(changes ...paramTypes.ParamChange) paramTypes.ParameterChangeProposal {
	return paramTypes.NewParameterChangeProposal(
		"Test",
		"description",
		changes,
	)
}

func (suite *ProposalTestSuite) TestGetSetProposal() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	tp := testProposal()
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalID

	_, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
}

func (suite *ProposalTestSuite) TestDeleteProposal() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	tp := testProposal()
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalID

	app.GovKeeper.DeleteProposal(ctx, proposalID)
	_, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.False(t, ok)
}

func (suite *ProposalTestSuite) TestGetProposals() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	tp := testProposal()
	_, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	proposals := app.GovKeeper.GetProposals(ctx)
	require.Len(t, proposals, 2)
}

func (suite *ProposalTestSuite) TestGetProposalID() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	tp := testProposal()
	_, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	// TODO - Ask is highest proposal id 1 more than total number of proposals

	proposalID, _ := app.GovKeeper.GetProposalID(ctx)
	require.Equal(t, proposalID, uint64(0x3))
}

func (suite *ProposalTestSuite) TestActivateVotingPeriod() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	tp := testProposal()
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	require.True(t, proposal.VotingStartTime.Equal(time.Time{}))

	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)

	require.True(t, proposal.VotingStartTime.Equal(ctx.BlockHeader().Time))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposal.ProposalID)
	require.True(t, ok)

	activeIterator := app.GovKeeper.ActiveProposalQueueIterator(ctx, proposal.VotingEndTime)
	require.True(t, activeIterator.Valid())
	activeIterator.Close()
}