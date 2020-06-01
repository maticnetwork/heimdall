package topup_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/topup"
	"github.com/maticnetwork/heimdall/topup/types"
	"github.com/maticnetwork/heimdall/types/simulation"
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
	suite.app, suite.ctx, _ = createTestApp(true)
}

// TestGenesisTestSuite
func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

//TestInitExportGenesis test import and export genesis state
func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	topupSequences := make([]string, 5)

	for i := range topupSequences {
		topupSequences[i] = strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	}
	genesisState := types.GenesisState{
		TopupSequences: topupSequences,
	}
	topup.InitGenesis(ctx, app.TopupKeeper, genesisState)

	actualParams := topup.ExportGenesis(ctx, app.TopupKeeper)

	require.LessOrEqual(t, len(topupSequences), len(actualParams.TopupSequences))
}
