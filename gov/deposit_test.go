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
	// authTypes "github.com/maticnetwork/heimdall/auth/types"
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

func (suite *DepositTestSuite) TestAddDeposit() {
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
	// fourStake := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(4)))

	_, found := app.GovKeeper.GetDeposit(ctx, proposalID, validators[0].ID)
	require.False(t, found)

	// app.BankKeeper.AddCoins(ctx, accounts[0].Address, sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(int64(4*10)))))

	// err, _ = app.GovKeeper.AddDeposit(ctx, proposalID, accounts[0].Address, fourStake, validators[0].ID)
	// require.Nil(t, err)
}

// func TestDeposits(t *testing.T) {
// 	input := getMockApp(t, 2, GenesisState{}, nil)

// 	SortAddresses(input.addrs)

// 	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
// 	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

// 	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

// 	tp := testProposal()
// 	proposal, err := input.keeper.SubmitProposal(ctx, tp)
// 	require.NoError(t, err)
// 	proposalID := proposal.ProposalID

// 	fourStake := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(4)))
// 	fiveStake := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(5)))

// 	addr0Initial := input.mApp.AccountKeeper.GetAccount(ctx, input.addrs[0]).GetCoins()
// 	addr1Initial := input.mApp.AccountKeeper.GetAccount(ctx, input.addrs[1]).GetCoins()

// 	expTokens := sdk.TokensFromConsensusPower(42)
// 	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, expTokens)), addr0Initial)
// 	require.True(t, proposal.TotalDeposit.IsEqual(sdk.NewCoins()))

// 	// Check no deposits at beginning
// 	deposit, found := input.keeper.GetDeposit(ctx, proposalID, input.addrs[1])
// 	require.False(t, found)
// 	proposal, ok := input.keeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	require.True(t, proposal.VotingStartTime.Equal(time.Time{}))

// 	// Check first deposit
// 	err, votingStarted := input.keeper.AddDeposit(ctx, proposalID, input.addrs[0], fourStake)
// 	require.Nil(t, err)
// 	require.False(t, votingStarted)
// 	deposit, found = input.keeper.GetDeposit(ctx, proposalID, input.addrs[0])
// 	require.True(t, found)
// 	require.Equal(t, fourStake, deposit.Amount)
// 	require.Equal(t, input.addrs[0], deposit.Depositor)
// 	proposal, ok = input.keeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	require.Equal(t, fourStake, proposal.TotalDeposit)
// 	require.Equal(t, addr0Initial.Sub(fourStake), input.mApp.AccountKeeper.GetAccount(ctx, input.addrs[0]).GetCoins())

// 	// Check a second deposit from same address
// 	err, votingStarted = input.keeper.AddDeposit(ctx, proposalID, input.addrs[0], fiveStake)
// 	require.Nil(t, err)
// 	require.False(t, votingStarted)
// 	deposit, found = input.keeper.GetDeposit(ctx, proposalID, input.addrs[0])
// 	require.True(t, found)
// 	require.Equal(t, fourStake.Add(fiveStake), deposit.Amount)
// 	require.Equal(t, input.addrs[0], deposit.Depositor)
// 	proposal, ok = input.keeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	require.Equal(t, fourStake.Add(fiveStake), proposal.TotalDeposit)
// 	require.Equal(t, addr0Initial.Sub(fourStake).Sub(fiveStake), input.mApp.AccountKeeper.GetAccount(ctx, input.addrs[0]).GetCoins())

// 	// Check third deposit from a new address
// 	err, votingStarted = input.keeper.AddDeposit(ctx, proposalID, input.addrs[1], fourStake)
// 	require.Nil(t, err)
// 	require.True(t, votingStarted)
// 	deposit, found = input.keeper.GetDeposit(ctx, proposalID, input.addrs[1])
// 	require.True(t, found)
// 	require.Equal(t, input.addrs[1], deposit.Depositor)
// 	require.Equal(t, fourStake, deposit.Amount)
// 	proposal, ok = input.keeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	require.Equal(t, fourStake.Add(fiveStake).Add(fourStake), proposal.TotalDeposit)
// 	require.Equal(t, addr1Initial.Sub(fourStake), input.mApp.AccountKeeper.GetAccount(ctx, input.addrs[1]).GetCoins())

// 	// Check that proposal moved to voting period
// 	proposal, ok = input.keeper.GetProposal(ctx, proposalID)
// 	require.True(t, ok)
// 	require.True(t, proposal.VotingStartTime.Equal(ctx.BlockHeader().Time))

// 	// Test deposit iterator
// 	depositsIterator := input.keeper.GetDepositsIterator(ctx, proposalID)
// 	require.True(t, depositsIterator.Valid())
// 	input.keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), &deposit)
// 	require.Equal(t, input.addrs[0], deposit.Depositor)
// 	require.Equal(t, fourStake.Add(fiveStake), deposit.Amount)
// 	depositsIterator.Next()
// 	input.keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), &deposit)
// 	require.Equal(t, input.addrs[1], deposit.Depositor)
// 	require.Equal(t, fourStake, deposit.Amount)
// 	depositsIterator.Next()
// 	require.False(t, depositsIterator.Valid())
// 	depositsIterator.Close()

// 	// Test Refund Deposits
// 	deposit, found = input.keeper.GetDeposit(ctx, proposalID, input.addrs[1])
// 	require.True(t, found)
// 	require.Equal(t, fourStake, deposit.Amount)
// 	input.keeper.RefundDeposits(ctx, proposalID)
// 	deposit, found = input.keeper.GetDeposit(ctx, proposalID, input.addrs[1])
// 	require.False(t, found)
// 	require.Equal(t, addr0Initial, input.mApp.AccountKeeper.GetAccount(ctx, input.addrs[0]).GetCoins())
// 	require.Equal(t, addr1Initial, input.mApp.AccountKeeper.GetAccount(ctx, input.addrs[1]).GetCoins())

// }