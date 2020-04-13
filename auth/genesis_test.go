package auth_test

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/auth/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types/simulation"
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
	r := rand.New(rand.NewSource(42))            // seed = 42
	accounts := simulation.RandomAccounts(r, 10) // create 10 accounts

	// genesis accounts
	var genesisAccs authTypes.GenesisAccounts
	for _, acc := range accounts {
		bacc := types.NewBaseAccountWithAddress(acc.Address)
		gacc, _ := types.NewGenesisAccountI(&bacc)
		genesisAccs = append(genesisAccs, gacc)
	}

	// app and context
	suite.app = app.SetupWithGenesisAccounts(genesisAccs)
	suite.ctx = suite.app.BaseApp.NewContext(true, abci.Header{})
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

//
// Tests
//

func (suite *GenesisTestSuite) TestInitGenesis() {
	t, happ, ctx := suite.T(), suite.app, suite.ctx

	accounts := happ.AccountKeeper.GetAllAccounts(ctx)
	require.LessOrEqual(t, 10, len(accounts))
}

func (suite *GenesisTestSuite) TestExportGenesis() {
	t, happ, ctx := suite.T(), suite.app, suite.ctx

	genesisState := auth.ExportGenesis(ctx, happ.AccountKeeper)
	require.LessOrEqual(t, 10, len(genesisState.Accounts))
}
