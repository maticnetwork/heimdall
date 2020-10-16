package gov_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/gov"
	"github.com/maticnetwork/heimdall/gov/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
)

//
// Test suite
//

type HandlerTestSuite struct {
	suite.Suite

	app     *app.HeimdallApp
	ctx     sdk.Context
	cliCtx  client.Context
	chainID string
	handler sdk.Handler
	r       *rand.Rand
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.handler = gov.NewHandler(suite.app.GovKeeper)

	// fetch chain id
	suite.chainID = suite.app.ChainKeeper.GetParams(suite.ctx).ChainParams.BorChainID

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

//
// Test cases
//

type validProposal struct{}

func (validProposal) GetTitle() string         { return "title" }
func (validProposal) GetDescription() string   { return "description" }
func (validProposal) ProposalRoute() string    { return types.RouterKey }
func (validProposal) ProposalType() string     { return types.ProposalTypeText }
func (validProposal) String() string           { return "" }
func (validProposal) ValidateBasic() sdk.Error { return nil }

func (suite *HandlerTestSuite) TestHandleMsgSubmitProposal() {
	t, app, ctx, _, _ := suite.T(), suite.app, suite.ctx, suite.chainID, suite.r

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

	msg := types.NewMsgSubmitProposal(
		testProposal(),
		sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(int64(4)))),
		accounts[0].Address,
		validators[0].ID,
	)

	t.Run("Success", func(t *testing.T) {
		result := suite.handler(ctx, msg)
		require.True(t, result.IsOK(), "expected submit proposal to be ok, got %v", result)
	})
}

func (suite *HandlerTestSuite) TestHandleMsgDeposit() {
	t, app, ctx, _, _ := suite.T(), suite.app, suite.ctx, suite.chainID, suite.r

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

	msg := types.NewMsgDeposit(
		accounts[0].Address,
		proposalID,
		sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(int64(4)))),
		validators[0].ID,
	)

	t.Run("Success", func(t *testing.T) {
		result := suite.handler(ctx, msg)
		require.True(t, result.IsOK(), "expected submit proposal to be ok, got %v", result)
	})
}

func (suite *HandlerTestSuite) TestHandleMsgVote() {
	t, app, ctx, _, _ := suite.T(), suite.app, suite.ctx, suite.chainID, suite.r

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
	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)
	proposalID := proposal.ProposalID

	msg := types.NewMsgVote(
		accounts[0].Address,
		proposalID,
		types.OptionYes,
		validators[0].ID,
	)

	t.Run("Success", func(t *testing.T) {
		result := suite.handler(ctx, msg)
		require.True(t, result.IsOK(), "expected add vote to be ok, got %v", result)
	})
}
