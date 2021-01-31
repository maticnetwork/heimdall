package chainmanager_test

import (
	"testing"

	"github.com/maticnetwork/heimdall/x/chainmanager/test_helper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/chainmanager"
	"github.com/maticnetwork/heimdall/x/chainmanager/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// GenesisTestSuite integrate test suite context object
type GenesisTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *GenesisTestSuite) SetupTest() {
	suite.app, suite.ctx = test_helper.CreateTestApp(true)
}

// TestGenesisTestSuite
func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

//TestInitExportGenesis test import and export genesis state
func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	params := types.DefaultParams()

	genesisState := types.GenesisState{
		Params: params,
	}
	chainmanager.InitGenesis(ctx, app.ChainKeeper, genesisState)

	actualParams := chainmanager.ExportGenesis(ctx, app.ChainKeeper)
	require.Equal(t, &genesisState, actualParams)

}
