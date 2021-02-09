package gov_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/gov"
	"github.com/maticnetwork/heimdall/x/gov/test_helper"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// GenesisTestSuite integrate test suite context object
type GenesisTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *GenesisTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
}

// TestGenesisTestSuite
func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

//TestIsEmptyGenesis test empty genesis state
func (suite *GenesisTestSuite) TestIsEmptyGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()
	emptyGenesisState := types.GenesisState{}

	require.True(t, emptyGenesisState.IsEmpty())
	require.False(t, defaultGenesisState.IsEmpty())
}

//TestNewGenesis test new genesis state
func (suite *GenesisTestSuite) TestNewGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()
	newGenesisState := types.NewGenesisState(defaultGenesisState.StartingProposalId, defaultGenesisState.DepositParams, defaultGenesisState.VotingParams, defaultGenesisState.TallyParams)

	require.True(t, defaultGenesisState.Equal(newGenesisState))
}

//TestEqualGenesis test equal genesis state
func (suite *GenesisTestSuite) TestEqualGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState1 := types.DefaultGenesis()
	defaultGenesisState2 := types.DefaultGenesis()
	emptyGenesisState := types.GenesisState{}

	require.True(t, defaultGenesisState1.Equal(*defaultGenesisState2))
	require.False(t, defaultGenesisState1.Equal(emptyGenesisState))
}

//TestValidate test valid genesis state
func (suite *GenesisTestSuite) TestValidate() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()
	require.Nil(t, defaultGenesisState.Validate())

	defaultGenesisState.TallyParams.Threshold = sdk.NewDecWithPrec(-1, 0)
	require.NotNil(t, defaultGenesisState.Validate())

	defaultGenesisState = types.DefaultGenesis()
	defaultGenesisState.TallyParams.Veto = sdk.NewDecWithPrec(-1, 0)
	require.NotNil(t, defaultGenesisState.Validate())
}

//TestInitExportGenesis test import and export genesis state
func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()

	gov.InitGenesis(ctx, app.AccountKeeper, app.BankKeeper, app.GovKeeper, *defaultGenesisState)
	returnedGenesisState := gov.ExportGenesis(ctx, app.GovKeeper)

	require.True(t, defaultGenesisState.Equal(*returnedGenesisState))
}
