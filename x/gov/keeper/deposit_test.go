package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmTypesCommon "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
// Test suite
//

// DepositTestSuite integrate test suite context object
type DepositTestSuite struct {
	suite.Suite

	sdk.Fee
	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *DepositTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
}

func TestDepositTestSuite(t *testing.T) {
	suite.Run(t, new(DepositTestSuite))
}

func (suite *DepositTestSuite) TestDeposits() {
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

	app.BankKeeper.AddCoins(ctx, accounts[0].Address, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(20))))
	app.BankKeeper.AddCoins(ctx, accounts[1].Address, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(20))))

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId

	fourStake := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(4)))
	fiveStake := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(5)))

	addr0Initial := app.BankKeeper.GetAllBalances(ctx, accounts[0].Address)
	addr1Initial := app.BankKeeper.GetAllBalances(ctx, accounts[1].Address)

	require.True(t, proposal.TotalDeposit.IsEqual(sdk.NewCoins()))

	// Check no deposits at beginning
	deposit, found := app.GovKeeper.GetDeposit(ctx, proposalID, validators[1].ID)
	require.False(t, found)
	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	require.True(t, proposal.VotingStartTime.Equal(time.Time{}))

	// Check first deposit
	err, votingStarted := app.GovKeeper.AddDeposit(ctx, proposalID, accounts[0].Address, fourStake, validators[0].ID)
	require.NoError(t, err)
	require.False(t, votingStarted)
	deposit, found = app.GovKeeper.GetDeposit(ctx, proposalID, validators[0].ID)
	require.True(t, found)
	require.Equal(t, fourStake, deposit.Amount)
	require.Equal(t, validators[0].ID, deposit.Depositor)
	proposal, ok = app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	require.Equal(t, fourStake, proposal.TotalDeposit)
	require.Equal(t, addr0Initial.Sub(fourStake), app.BankKeeper.GetAllBalances(ctx, accounts[0].Address))

	// Check a second deposit from same address
	err, votingStarted = app.GovKeeper.AddDeposit(ctx, proposalID, accounts[0].Address, fiveStake, validators[0].ID)
	require.NoError(t, err)
	require.False(t, votingStarted)
	deposit, found = app.GovKeeper.GetDeposit(ctx, proposalID, validators[0].ID)
	require.True(t, found)
	require.Equal(t, fourStake.Add(fiveStake...), deposit.Amount)
	require.Equal(t, validators[0].ID, deposit.Depositor)
	proposal, ok = app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	require.Equal(t, fourStake.Add(fiveStake...), proposal.TotalDeposit)
	require.Equal(t, addr0Initial.Sub(fourStake).Sub(fiveStake), app.BankKeeper.GetAllBalances(ctx, accounts[0].Address))

	// Check third deposit from a new address
	err, votingStarted = app.GovKeeper.AddDeposit(ctx, proposalID, accounts[1].Address, fourStake, validators[1].ID)
	require.NoError(t, err)
	require.True(t, votingStarted)
	deposit, found = app.GovKeeper.GetDeposit(ctx, proposalID, validators[1].ID)
	require.True(t, found)
	require.Equal(t, validators[1].ID, deposit.Depositor)
	require.Equal(t, fourStake, deposit.Amount)
	proposal, ok = app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	require.Equal(t, fourStake.Add(fiveStake...).Add(fourStake...), proposal.TotalDeposit)
	require.Equal(t, addr1Initial.Sub(fourStake), app.BankKeeper.GetAllBalances(ctx, accounts[1].Address))

	// Check that proposal moved to voting period
	proposal, ok = app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	require.True(t, proposal.VotingStartTime.Equal(ctx.BlockHeader().Time))

	// Test deposit iterator
	// NOTE order of deposits is determined by the addresses
	deposits := app.GovKeeper.GetAllDeposits(ctx)
	require.Len(t, deposits, 2)
	require.Equal(t, deposits, app.GovKeeper.GetDeposits(ctx, proposalID))
	require.Equal(t, validators[0].ID, deposits[0].Depositor)
	require.Equal(t, fourStake.Add(fiveStake...), deposits[0].Amount)
	require.Equal(t, validators[1].ID, deposits[1].Depositor)
	require.Equal(t, fourStake, deposits[1].Amount)

	// TODO - Check this
	// // Test Refund Deposits
	// deposit, found = app.GovKeeper.GetDeposit(ctx, proposalID, validators[1].ID)
	// require.True(t, found)
	// require.Equal(t, fourStake, deposit.Amount)
	// app.GovKeeper.RefundDeposits(ctx, proposalID)
	// deposit, found = app.GovKeeper.GetDeposit(ctx, proposalID, validators[1].ID)
	// require.False(t, found)
	// require.Equal(t, addr0Initial, app.BankKeeper.GetAllBalances(ctx, accounts[0].Address))
	// require.Equal(t, addr1Initial, app.BankKeeper.GetAllBalances(ctx, accounts[1].Address))
}
