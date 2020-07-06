package gov_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/gov"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisTestSuite integrate test suite context object
type GenesisTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *GenesisTestSuite) SetupTest() {
	suite.app = setupGovGenesis()
	suite.ctx = suite.app.BaseApp.NewContext(true, abci.Header{})
}

// TestGenesisTestSuite
func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

//TestIsEmptyGenesis test empty genesis state
func (suite *GenesisTestSuite) TestIsEmptyGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := gov.DefaultGenesisState()
	emptyGenesisState := gov.GenesisState{}

	require.True(t, emptyGenesisState.IsEmpty())
	require.False(t, defaultGenesisState.IsEmpty())
}

//TestEqualGenesis test equal genesis state
func (suite *GenesisTestSuite) TestEqualGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState1 := gov.DefaultGenesisState()
	defaultGenesisState2 := gov.DefaultGenesisState()
	emptyGenesisState := gov.GenesisState{}

	require.True(t, defaultGenesisState1.Equal(defaultGenesisState2))
	require.False(t, defaultGenesisState1.Equal(emptyGenesisState))
}

//TestValidateGenesis test valid genesis state
func (suite *GenesisTestSuite) TestValidateGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := gov.DefaultGenesisState()
	require.Nil(t, gov.ValidateGenesis(defaultGenesisState))

	defaultGenesisState.TallyParams.Threshold = sdk.NewDecWithPrec(-1, 0)
	require.NotNil(t, gov.ValidateGenesis(defaultGenesisState))

	defaultGenesisState.TallyParams.Veto = sdk.NewDecWithPrec(-1, 0)
	require.NotNil(t, gov.ValidateGenesis(defaultGenesisState))

}

//TestInitExportGenesis test import and export genesis state
func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	defaultGenesisState := gov.DefaultGenesisState()

	gov.InitGenesis(ctx, app.GovKeeper, app.SupplyKeeper, defaultGenesisState)
	returnedGenesisState := gov.ExportGenesis(ctx, app.GovKeeper)

	require.True(t, defaultGenesisState.Equal(returnedGenesisState))
}
