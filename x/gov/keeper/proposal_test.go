package keeper_test

import (
	"errors"
	// "fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmTypesCommon "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/maticnetwork/heimdall/x/gov/test_helper"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// ProposalTestSuite integrate test suite context object
type ProposalTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *ProposalTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
}

func TestProposalTestSuite(t *testing.T) {
	suite.Run(t, new(ProposalTestSuite))
}

func (suite *ProposalTestSuite) TestGetSetProposal() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	app.GovKeeper.SetProposal(ctx, proposal)

	gotProposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	require.True(t, proposal.Equal(gotProposal))
}

func (suite *ProposalTestSuite) TestActivateVotingPeriod() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	require.True(t, proposal.VotingStartTime.Equal(time.Time{}))

	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)

	require.True(t, proposal.VotingStartTime.Equal(ctx.BlockHeader().Time))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposal.ProposalId)
	require.True(t, ok)

	activeIterator := app.GovKeeper.ActiveProposalQueueIterator(ctx, proposal.VotingEndTime)
	require.True(t, activeIterator.Valid())

	proposalID := types.GetProposalIDFromBytes(activeIterator.Value())
	require.Equal(t, proposalID, proposal.ProposalId)
	activeIterator.Close()
}

type invalidProposalRoute struct{ types.TextProposal }

func (invalidProposalRoute) ProposalRoute() string { return "nonexistingroute" }

func (suite *ProposalTestSuite) TestSubmitProposal() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	testCases := []struct {
		content     types.Content
		expectedErr error
	}{
		{&types.TextProposal{Title: "title", Description: "description"}, nil},
		// Keeper does not check the validity of title and description, no error
		{&types.TextProposal{Title: "", Description: "description"}, nil},
		{&types.TextProposal{Title: strings.Repeat("1234567890", 100), Description: "description"}, nil},
		{&types.TextProposal{Title: "title", Description: ""}, nil},
		{&types.TextProposal{Title: "title", Description: strings.Repeat("1234567890", 1000)}, nil},
		// error only when invalid route
		{&invalidProposalRoute{}, types.ErrNoProposalHandlerExists},
	}

	for i, tc := range testCases {
		_, err := app.GovKeeper.SubmitProposal(ctx, tc.content)
		require.True(t, errors.Is(tc.expectedErr, err), "tc #%d; got: %v, expected: %v", i, err, tc.expectedErr)
	}
}

func (suite *ProposalTestSuite) TestGetProposalsFiltered() {
	proposalID := uint64(1)
	t, app, ctx := suite.T(), suite.app, suite.ctx

	status := []types.ProposalStatus{types.StatusDepositPeriod, types.StatusVotingPeriod}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 1

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

	for _, s := range status {
		for i := 0; i < 50; i++ {
			p, err := types.NewProposal(test_helper.TestProposal, proposalID, time.Now(), time.Now())
			require.NoError(t, err)

			p.Status = s

			if i%2 == 0 {
				d := types.NewDeposit(proposalID, nil, validators[0].ID)
				v := types.NewVote(proposalID, validators[0].ID, types.OptionYes)
				app.GovKeeper.SetDeposit(ctx, proposalID, validators[0].ID, d)
				app.GovKeeper.SetVote(ctx, proposalID, validators[0].ID, v)
			}

			app.GovKeeper.SetProposal(ctx, p)
			proposalID++
		}
	}

	// testCases := []struct {
	// 	params             types.QueryProposalsParams
	// 	expectedNumResults int
	// }{
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusNil, nil, nil), 50},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, nil, nil), 50},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusVotingPeriod, nil, nil), 50},
	// 	{types.NewQueryProposalsParams(1, 25, types.StatusNil, nil, nil), 25},
	// 	{types.NewQueryProposalsParams(2, 25, types.StatusNil, nil, nil), 25},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusRejected, nil, nil), 0},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusNil, accounts[0].Address, nil), 50},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusNil, nil, accounts[0].Address), 50},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusNil, accounts[0].Address, accounts[0].Address), 50},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, accounts[0].Address, accounts[0].Address), 25},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, nil, nil), 50},
	// 	{types.NewQueryProposalsParams(1, 50, types.StatusVotingPeriod, nil, nil), 50},
	// }

	// for i, tc := range testCases {
	// 	t.Run(fmt.Sprintf("Test Case %d", i), func(t *testing.T) {
	// 		proposals := app.GovKeeper.GetProposalsFiltered(ctx, tc.params)
	// 		require.Len(t, proposals, tc.expectedNumResults)

	// 		for _, p := range proposals {
	// 			if types.ValidProposalStatus(tc.params.ProposalStatus) {
	// 				require.Equal(t, tc.params.ProposalStatus, p.Status)
	// 			}
	// 		}
	// 	})
	// }
}
