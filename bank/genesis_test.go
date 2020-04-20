package bank_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bank"
)

//
// Test suite
//

// GenesisTestSuite integrate test suite context object
type GenesisTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

//
// Tests
//

func (suite *GenesisTestSuite) TestInitGenesis() {
	t, happ, ctx := suite.T(), suite.app, suite.ctx

	require.Equal(t, true, happ.BankKeeper.GetSendEnabled(ctx))
}

func (suite *GenesisTestSuite) TestExportGenesis() {
	t, happ, ctx := suite.T(), suite.app, suite.ctx

	genesisState := bank.ExportGenesis(ctx, happ.BankKeeper)
	require.Equal(t, true, genesisState.SendEnabled)

	happ.BankKeeper.SetSendEnabled(ctx, false)
	genesisState = bank.ExportGenesis(ctx, happ.BankKeeper)
	require.Equal(t, false, genesisState.SendEnabled)
}
