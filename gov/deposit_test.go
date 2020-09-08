package gov_test

import (
	// "fmt"
	"math/rand"
	"time"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types/simulation"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
// Test suite
//

// DepositTestSuite integrate test suite context object
type DepositTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *DepositTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestDepositTestSuite(t *testing.T) {
	suite.Run(t, new(DepositTestSuite))
}

func (suite *DepositTestSuite) TestAddGetDeposit() {
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

	app.BankKeeper.AddCoins(ctx, accounts[0].Address, sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(int64(4*10)))))
	tp := testProposal()
	proposal, _ := app.GovKeeper.SubmitProposal(ctx, tp)
	proposalID := proposal.ProposalID

	err, _ := app.GovKeeper.AddDeposit(ctx, proposalID, accounts[0].Address, sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(int64(4)))), validators[0].ID)

	require.Nil(t, err)

	_, found := app.GovKeeper.GetDeposit(ctx, proposalID, validators[0].ID)
	require.True(t, found)

	deposits := app.GovKeeper.GetAllDeposits(ctx)
	require.Len(t, deposits, 1)

	deposits = app.GovKeeper.GetDeposits(ctx, proposalID)
	require.Len(t, deposits, 1)

	app.GovKeeper.RefundDeposits(ctx, proposalID)

	deposits = app.GovKeeper.GetDeposits(ctx, proposalID)
	require.Len(t, deposits, 0)
}