package milestone_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/milestone"
	"github.com/maticnetwork/heimdall/milestone/types"
	hmTypes "github.com/maticnetwork/heimdall/types"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

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
	t.Parallel()
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmTypes.HexToHeimdallHash("123")

	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"
	milestoneID := "0000"

	milestoneMock := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)

	params := types.DefaultParams()
	genesisState := types.NewGenesisState(
		params,
		&milestoneMock,
	)

	milestone.InitGenesis(ctx, app.MilestoneKeeper, genesisState)

	actualParams := milestone.ExportGenesis(ctx, app.MilestoneKeeper)

	require.Equal(t, genesisState.Milestone, actualParams.Milestone)
	require.Equal(t, genesisState.Params, actualParams.Params)
}

func (suite *GenesisTestSuite) TestInitExportGenesisWithNilMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	params := types.DefaultParams()
	genesisState := types.NewGenesisState(
		params,
		nil,
	)

	milestone.InitGenesis(ctx, app.MilestoneKeeper, genesisState)

	actualParams := milestone.ExportGenesis(ctx, app.MilestoneKeeper)

	require.Nil(t, actualParams.Milestone)
	require.Equal(t, genesisState.Params, actualParams.Params)
}
