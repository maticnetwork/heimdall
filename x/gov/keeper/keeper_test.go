package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
	"github.com/maticnetwork/heimdall/app"
)

// KeeperTestSuite integrate test suite context object
type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite)TestIncrementProposalNumber() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	tp := TestProposal
	_, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposal6, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	require.Equal(t, uint64(6), proposal6.ProposalId)
}

func (suite *KeeperTestSuite)TestProposalQueues() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	// create test proposals
	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	inactiveIterator := app.GovKeeper.InactiveProposalQueueIterator(ctx, proposal.DepositEndTime)
	require.True(t, inactiveIterator.Valid())

	proposalID := types.GetProposalIDFromBytes(inactiveIterator.Value())
	require.Equal(t, proposalID, proposal.ProposalId)
	inactiveIterator.Close()

	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposal.ProposalId)
	require.True(t, ok)

	activeIterator := app.GovKeeper.ActiveProposalQueueIterator(ctx, proposal.VotingEndTime)
	require.True(t, activeIterator.Valid())

	proposalID, _ = types.SplitActiveProposalQueueKey(activeIterator.Key())
	require.Equal(t, proposalID, proposal.ProposalId)

	activeIterator.Close()
}
