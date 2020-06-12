package clerk_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/clerk"
	"github.com/maticnetwork/heimdall/clerk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
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
	suite.app = setupClerkGenesis()
	suite.ctx = suite.app.BaseApp.NewContext(true, abci.Header{})
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

	recordSequences := make([]string, 5)
	eventRecords := make([]*types.EventRecord, 1)

	for i := range recordSequences {
		recordSequences[i] = strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	}

	for i := range eventRecords {
		hAddr := hmTypes.BytesToHeimdallAddress([]byte(strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))))
		hHash := hmTypes.BytesToHeimdallHash([]byte(strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))))
		testEventRecord := types.NewEventRecord(hHash, uint64(i), uint64(i), hAddr, make([]byte, 0), strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000)), time.Now())
		eventRecords[i] = &testEventRecord
	}
	genesisState := types.GenesisState{
		EventRecords:    eventRecords,
		RecordSequences: recordSequences,
	}
	clerk.InitGenesis(ctx, app.ClerkKeeper, genesisState)

	actualParams := clerk.ExportGenesis(ctx, app.ClerkKeeper)

	require.Equal(t, len(recordSequences), len(actualParams.RecordSequences))
	require.Equal(t, len(eventRecords), len(actualParams.EventRecords))
}
