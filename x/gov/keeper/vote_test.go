package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmTypesCommon "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// VoteTestSuite integrate test suite context object
type VoteTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *VoteTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
}

func TestVoteTestSuite(t *testing.T) {
	suite.Run(t, new(VoteTestSuite))
}

func (suite *VoteTestSuite) TestVotes() {
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
			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
			accounts[i].Address,
		)

		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId

	var invalidOption types.VoteOption = 0x10

	require.Error(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID), "proposal not on voting period")
	require.Error(t, app.GovKeeper.AddVote(ctx, 10, accounts[0].Address, types.OptionYes, validators[0].ID), "invalid proposal ID")

	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.Error(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, invalidOption, validators[0].ID), "invalid option")

	// Test first vote
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionAbstain, validators[0].ID))
	vote, found := app.GovKeeper.GetVote(ctx, proposalID, validators[0].ID)
	require.True(t, found)
	require.Equal(t, proposalID, vote.ProposalId)
	require.Equal(t, types.OptionAbstain, vote.Option)

	// Test change of vote
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
	vote, found = app.GovKeeper.GetVote(ctx, proposalID, validators[0].ID)
	require.True(t, found)
	//require.Equal(t, accounts[0].Address.String(), vote.Voter)
	require.Equal(t, proposalID, vote.ProposalId)
	require.Equal(t, types.OptionYes, vote.Option)

	// Test second vote
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionNoWithVeto, validators[1].ID))
	vote, found = app.GovKeeper.GetVote(ctx, proposalID, validators[1].ID)
	require.True(t, found)
	//require.Equal(t, accounts[1].Address.String(), vote.Voter)
	require.Equal(t, proposalID, vote.ProposalId)
	require.Equal(t, types.OptionNoWithVeto, vote.Option)

	// Test vote iterator
	// NOTE order of deposits is determined by the addresses
	votes := app.GovKeeper.GetAllVotes(ctx)
	require.Len(t, votes, 2)
	require.Equal(t, votes, app.GovKeeper.GetVotes(ctx, proposalID))
	//require.Equal(t, accounts[0].Address.String(), votes[0].Voter)
	require.Equal(t, proposalID, votes[0].ProposalId)
	require.Equal(t, types.OptionYes, votes[0].Option)
	//require.Equal(t, accounts[1].Address.String(), votes[1].Voter)
	require.Equal(t, proposalID, votes[1].ProposalId)
	require.Equal(t, types.OptionNoWithVeto, votes[1].Option)
}
