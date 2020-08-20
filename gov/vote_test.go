package gov_test

import (
	"testing"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	// "github.com/maticnetwork/heimdall/gov"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/maticnetwork/heimdall/gov/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

// VoteTestSuite integrate test suite context object
type VoteTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *VoteTestSuite) SetupTest() {
	suite.app = setupGovGenesis()
	suite.ctx = suite.app.BaseApp.NewContext(true, abci.Header{})
}

// TestVoteTestSuite
func TestVoteTestSuite(t *testing.T) {
	suite.Run(t, new(VoteTestSuite))
}

func (suite *TallyTestSuite) TestAddVote() {

	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5

	validators := make([]*hmTypes.Validator, n)
	accounts := simulation.RandomAccounts(r1, n)

	for i := range validators {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			1,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)

		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	tp := testProposal()
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalID
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	tp = testProposal()
	proposal, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID = proposal.ProposalID
	proposal.Status = types.StatusNil
	app.GovKeeper.SetProposal(ctx, proposal)

	err = app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID)
	require.Error(t, err)

}

func (suite *TallyTestSuite) TestGetVotesAllFunctions() {

	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5

	validators := make([]*hmTypes.Validator, n)
	accounts := simulation.RandomAccounts(r1, n)

	for i := range validators {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			1,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)

		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	tp := testProposal()
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalID
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	err = app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID)
	require.Nil(t, err)

	err = app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionYes, validators[1].ID)
	require.Nil(t, err)

	votes := app.GovKeeper.GetAllVotes(ctx)
	require.Len(t, votes, 2)

	votes = app.GovKeeper.GetVotes(ctx, proposalID)
	require.Len(t, votes, 2)

	_, found := app.GovKeeper.GetVote(ctx, proposalID, validators[0].ID)
	require.True(t, found)

	_, found = app.GovKeeper.GetVote(ctx, proposalID, validators[2].ID)
	require.False(t, found)

}