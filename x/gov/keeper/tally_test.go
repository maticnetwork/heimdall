package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// "github.com/maticnetwork/heimdall/x/staking"
	// stakingtypes "github.com/maticnetwork/heimdall/x/staking/types"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmTypesCommon "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/maticnetwork/heimdall/x/gov/test_helper"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// TallyTestSuite integrate test suite context object
type TallyTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *TallyTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
}

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
			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
			accounts[i].Address,
		)

		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, _, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.True(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func (suite *TallyTestSuite) TestTallyNoQuorum() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 2

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

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	err = app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID)
	require.Nil(t, err)

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, _, _ := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
}

func (suite *TallyTestSuite) TestTallyOnlyValidatorsAllYes() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 3

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
	tp := test_helper.TestProposal

	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionYes, validators[1].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionYes, validators[2].ID))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	_, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, burnDeposits)
}

func (suite *TallyTestSuite) TestTallyOnlyValidators51No() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 2

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

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionNo, validators[1].ID))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
}

func (suite *TallyTestSuite) TestTallyOnlyValidators51Yes() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 2

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

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionNo, validators[0].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionYes, validators[1].ID))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	_, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, burnDeposits)
	// require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func (suite *TallyTestSuite) TestTallyOnlyValidatorsVetoed() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 3

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

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionYes, validators[1].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionNoWithVeto, validators[2].ID))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, _, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
}

func (suite *TallyTestSuite) TestTallyOnlyValidatorsAbstainPasses() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 4

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

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionAbstain, validators[0].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionNo, validators[1].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionYes, validators[2].ID))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	_, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, burnDeposits)
}

func (suite *TallyTestSuite) TestTallyOnlyValidatorsAbstainFails() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 3

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

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionAbstain, validators[0].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionYes, validators[1].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionNo, validators[2].ID))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
}

func (suite *TallyTestSuite) TestTallyOnlyValidatorsNonVoter() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 2

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
	valAccAddr1, valAccAddr2 := accounts[0].Address, accounts[1].Address

	tp := test_helper.TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddr1, types.OptionYes, validators[0].ID))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddr2, types.OptionNo, validators[1].ID))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	// require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

// func (suite *TallyTestSuite)TestTallyDelgatorOverride() {
// 	t, app, ctx := suite.T(), suite.app, suite.ctx

// 	s1 := rand.NewSource(time.Now().UnixNano())
// 	r1 := rand.New(s1)
// 	n := 2

// 	validators := make([]*hmTypes.Validator, n)
// 	accounts := simulation.RandomAccounts(r1, n)

// 	for i := range validators {
// 		// validator
// 		validators[i] = hmTypes.NewValidator(
// 			hmTypes.NewValidatorID(uint64(int64(i))),
// 			0,
// 			0,
// 			1,
// 			int64(simulation.RandIntBetween(r1, 10, 100)), // power
// 			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
// 			accounts[i].Address,
// 		)

// 		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
// 		if err != nil {
// 			t.Error("Error while adding validator to store", err)
// 		}
// 	}

// 	delTokens := sdk.TokensFromConsensusPower(30)
// 	val1, found := app.StakingKeeper.GetValidator(ctx, accounts[0].Address)
// 	require.True(t, found)

// 	_, err := app.StakingKeeper.Delegate(ctx, accounts[4].Address, delTokens, stakingtypes.Unbonded, val1, true)
// 	require.NoError(t, err)

// 	_ = staking.EndBlocker(ctx, app.StakingKeeper)

// 	tp := test_helper.TestProposal
// 	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
// 	require.NoError(t, err)
// 	proposalID := proposal.ProposalId
// 	proposal.Status = types.StatusVotingPeriod
// 	app.GovKeeper.SetProposal(ctx, proposal)

// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionYes, validators[0].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionYes, validators[1].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[3].Address, types.OptionYes, validators[2].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[4].Address, types.OptionNo, validators[3].ID))

// 	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

// 	require.False(t, passes)
// 	require.False(t, burnDeposits)
// 	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
// }

// func (suite *TallyTestSuite)TestTallyDelgatorInherit() {
// 	t, app, ctx := suite.T(), suite.app, suite.ctx

// 	s1 := rand.NewSource(time.Now().UnixNano())
// 	r1 := rand.New(s1)
// 	n := 2

// 	validators := make([]*hmTypes.Validator, n)
// 	accounts := simulation.RandomAccounts(r1, n)

// 	for i := range validators {
// 		// validator
// 		validators[i] = hmTypes.NewValidator(
// 			hmTypes.NewValidatorID(uint64(int64(i))),
// 			0,
// 			0,
// 			1,
// 			int64(simulation.RandIntBetween(r1, 10, 100)), // power
// 			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
// 			accounts[i].Address,
// 		)

// 		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
// 		if err != nil {
// 			t.Error("Error while adding validator to store", err)
// 		}
// 	}

// 	delTokens := sdk.TokensFromConsensusPower(30)
// 	val3, found := app.StakingKeeper.GetValidator(ctx, vals[2])
// 	require.True(t, found)

// 	_, err := app.StakingKeeper.Delegate(ctx, accounts[3].Address, delTokens, stakingtypes.Unbonded, val3, true)
// 	require.NoError(t, err)

// 	_ = staking.EndBlocker(ctx, app.StakingKeeper)

// 	tp := test_helper.TestProposal
// 	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
// 	require.NoError(t, err)
// 	proposalID := proposal.ProposalId
// 	proposal.Status = types.StatusVotingPeriod
// 	app.GovKeeper.SetProposal(ctx, proposal)

// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionNo, validators[0].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionNo, validators[1].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionYes, validators[2].ID))

// 	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

// 	require.True(t, passes)
// 	require.False(t, burnDeposits)
// 	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
// }

// func (suite *TallyTestSuite)TestTallyDelgatorMultipleOverride() {
// 	t, app, ctx := suite.T(), suite.app, suite.ctx

// 	s1 := rand.NewSource(time.Now().UnixNano())
// 	r1 := rand.New(s1)
// 	n := 2

// 	validators := make([]*hmTypes.Validator, n)
// 	accounts := simulation.RandomAccounts(r1, n)

// 	for i := range validators {
// 		// validator
// 		validators[i] = hmTypes.NewValidator(
// 			hmTypes.NewValidatorID(uint64(int64(i))),
// 			0,
// 			0,
// 			1,
// 			int64(simulation.RandIntBetween(r1, 10, 100)), // power
// 			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
// 			accounts[i].Address,
// 		)

// 		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
// 		if err != nil {
// 			t.Error("Error while adding validator to store", err)
// 		}
// 	}

// 	delTokens := sdk.TokensFromConsensusPower(10)
// 	val1, found := app.StakingKeeper.GetValidator(ctx, vals[0])
// 	require.True(t, found)
// 	val2, found := app.StakingKeeper.GetValidator(ctx, vals[1])
// 	require.True(t, found)

// 	_, err := app.StakingKeeper.Delegate(ctx, accounts[3].Address, delTokens, stakingtypes.Unbonded, val1, true)
// 	require.NoError(t, err)
// 	_, err = app.StakingKeeper.Delegate(ctx, accounts[3].Address, delTokens, stakingtypes.Unbonded, val2, true)
// 	require.NoError(t, err)

// 	_ = staking.EndBlocker(ctx, app.StakingKeeper)

// 	tp := test_helper.TestProposal
// 	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
// 	require.NoError(t, err)
// 	proposalID := proposal.ProposalId
// 	proposal.Status = types.StatusVotingPeriod
// 	app.GovKeeper.SetProposal(ctx, proposal)

// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionYes, validators[1].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionYes, validators[2].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[3].Address, types.OptionNo, validators[3].ID))

// 	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

// 	require.False(t, passes)
// 	require.False(t, burnDeposits)
// 	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
// }

// func(suite *TallyTestSuite) TestTallyDelgatorMultipleInherit() {
// 	t, app, ctx := suite.T(), suite.app, suite.ctx

// 	s1 := rand.NewSource(time.Now().UnixNano())
// 	r1 := rand.New(s1)
// 	n := 2

// 	validators := make([]*hmTypes.Validator, n)
// 	accounts := simulation.RandomAccounts(r1, n)

// 	for i := range validators {
// 		// validator
// 		validators[i] = hmTypes.NewValidator(
// 			hmTypes.NewValidatorID(uint64(int64(i))),
// 			0,
// 			0,
// 			1,
// 			int64(simulation.RandIntBetween(r1, 10, 100)), // power
// 			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
// 			accounts[i].Address,
// 		)

// 		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
// 		if err != nil {
// 			t.Error("Error while adding validator to store", err)
// 		}
// 	}

// 	s1 = rand.NewSource(time.Now().UnixNano())
// 	r1 = rand.New(s1)
// 	n = 1

// 	validators = make([]*hmTypes.Validator, n)
// 	accounts = simulation.RandomAccounts(r1, n)

// 	for i := range validators {
// 		// validator
// 		validators[i] = hmTypes.NewValidator(
// 			hmTypes.NewValidatorID(uint64(int64(i))),
// 			0,
// 			0,
// 			1,
// 			int64(simulation.RandIntBetween(r1, 10, 100)), // power
// 			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
// 			accounts[i].Address,
// 		)

// 		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
// 		if err != nil {
// 			t.Error("Error while adding validator to store", err)
// 		}
// 	}

// 	delTokens := sdk.TokensFromConsensusPower(10)
// 	val2, found := app.StakingKeeper.GetValidator(ctx, vals[1])
// 	require.True(t, found)
// 	val3, found := app.StakingKeeper.GetValidator(ctx, vals[2])
// 	require.True(t, found)

// 	_, err := app.StakingKeeper.Delegate(ctx, accounts[3].Address, delTokens, stakingtypes.Unbonded, val2, true)
// 	require.NoError(t, err)
// 	_, err = app.StakingKeeper.Delegate(ctx, accounts[3].Address, delTokens, stakingtypes.Unbonded, val3, true)
// 	require.NoError(t, err)

// 	_ = staking.EndBlocker(ctx, app.StakingKeeper)

// 	tp := test_helper.TestProposal
// 	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
// 	require.NoError(t, err)
// 	proposalID := proposal.ProposalId
// 	proposal.Status = types.StatusVotingPeriod
// 	app.GovKeeper.SetProposal(ctx, proposal)

// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionNo, validators[1].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionNo, validators[2].ID))

// 	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

// 	require.False(t, passes)
// 	require.False(t, burnDeposits)
// 	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
// }

// func (suite *TallyTestSuite)TestTallyJailedValidator() {
// 	t, app, ctx := suite.T(), suite.app, suite.ctx

// 	s1 := rand.NewSource(time.Now().UnixNano())
// 	r1 := rand.New(s1)
// 	n := 2

// 	validators := make([]*hmTypes.Validator, n)
// 	accounts := simulation.RandomAccounts(r1, n)

// 	for i := range validators {
// 		// validator
// 		validators[i] = hmTypes.NewValidator(
// 			hmTypes.NewValidatorID(uint64(int64(i))),
// 			0,
// 			0,
// 			1,
// 			int64(simulation.RandIntBetween(r1, 10, 100)), // power
// 			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
// 			accounts[i].Address,
// 		)

// 		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
// 		if err != nil {
// 			t.Error("Error while adding validator to store", err)
// 		}
// 	}

// 	delTokens := sdk.TokensFromConsensusPower(10)
// 	val2, found := app.StakingKeeper.GetValidator(ctx, accounts[1].Address)
// 	require.True(t, found)
// 	val3, found := app.StakingKeeper.GetValidator(ctx, accounts[2].Address)
// 	require.True(t, found)

// 	_, err := app.StakingKeeper.Delegate(ctx, accounts[3].Address, delTokens, stakingtypes.Unbonded, val2, true)
// 	require.NoError(t, err)
// 	_, err = app.StakingKeeper.Delegate(ctx, accounts[3].Address, delTokens, stakingtypes.Unbonded, val3, true)
// 	require.NoError(t, err)

// 	_ = staking.EndBlocker(ctx, app.StakingKeeper)

// 	consAddr, err := val2.GetConsAddr()
// 	require.NoError(t, err)
// 	app.StakingKeeper.Jail(ctx, sdk.ConsAddress(consAddr.Bytes()))

// 	tp := test_helper.TestProposal
// 	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
// 	require.NoError(t, err)
// 	proposalID := proposal.ProposalId
// 	proposal.Status = types.StatusVotingPeriod
// 	app.GovKeeper.SetProposal(ctx, proposal)

// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionNo, validators[1].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionNo, validators[2].ID))

// 	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

// 	require.True(t, passes)
// 	require.False(t, burnDeposits)
// 	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
// }

// func (suite *TallyTestSuite)TestTallyValidatorMultipleDelegations() {
// 	t, app, ctx := suite.T(), suite.app, suite.ctx

// 	s1 := rand.NewSource(time.Now().UnixNano())
// 	r1 := rand.New(s1)
// 	n := 2

// 	validators := make([]*hmTypes.Validator, n)
// 	accounts := simulation.RandomAccounts(r1, n)

// 	for i := range validators {
// 		// validator
// 		validators[i] = hmTypes.NewValidator(
// 			hmTypes.NewValidatorID(uint64(int64(i))),
// 			0,
// 			0,
// 			1,
// 			int64(simulation.RandIntBetween(r1, 10, 100)), // power
// 			hmTypesCommon.NewPubKey(accounts[i].Address.Bytes()),
// 			accounts[i].Address,
// 		)

// 		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
// 		if err != nil {
// 			t.Error("Error while adding validator to store", err)
// 		}
// 	}

// 	delTokens := sdk.TokensFromConsensusPower(10)
// 	val2, found := app.StakingKeeper.GetValidator(ctx, accounts[1].Address)
// 	require.True(t, found)

// 	_, err := app.StakingKeeper.Delegate(ctx, accounts[0].Address, delTokens, stakingtypes.Unbonded, val2, true)
// 	require.NoError(t, err)

// 	tp := test_helper.TestProposal
// 	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
// 	require.NoError(t, err)
// 	proposalID := proposal.ProposalId
// 	proposal.Status = types.StatusVotingPeriod
// 	app.GovKeeper.SetProposal(ctx, proposal)

// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[0].Address, types.OptionYes, validators[0].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[1].Address, types.OptionNo, validators[1].ID))
// 	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, accounts[2].Address, types.OptionYes, validators[2].ID))

// 	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

// 	require.True(t, passes)
// 	require.False(t, burnDeposits)

// 	expectedYes := sdk.TokensFromConsensusPower(30)
// 	expectedAbstain := sdk.TokensFromConsensusPower(0)
// 	expectedNo := sdk.TokensFromConsensusPower(10)
// 	expectedNoWithVeto := sdk.TokensFromConsensusPower(0)
// 	expectedTallyResult := types.NewTallyResult(expectedYes, expectedAbstain, expectedNo, expectedNoWithVeto)

// 	require.True(t, tallyResults.Equals(expectedTallyResult))
// }
