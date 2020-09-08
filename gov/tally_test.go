package gov_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/gov"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/maticnetwork/heimdall/gov/types"
	"github.com/maticnetwork/heimdall/helper/mocks"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TallyTestSuite integrate test suite context object
type TallyTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
	contractCaller mocks.IContractCaller
}

// SetupTest setup necessary things for genesis test
func (suite *TallyTestSuite) SetupTest() {
	suite.app = setupGovGenesis()
	suite.ctx = suite.app.BaseApp.NewContext(true, abci.Header{})
	suite.contractCaller = mocks.IContractCaller{}
}

// TestTallyTestSuite
func TestTallyTestSuite(t *testing.T) {
	suite.Run(t, new(TallyTestSuite))
}

func (suite *TallyTestSuite) TestTallyNoOneVotes() {

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

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := gov.Tally(ctx, app.GovKeeper, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
}

func (suite *TallyTestSuite) TestTallyNoQuorum() {

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

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := gov.Tally(ctx, app.GovKeeper, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
}