package bor_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/bor/types"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"
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

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	spancount := 5
	params := types.DefaultParams()
	spans := []*hmTypes.Span{}
	chainID := "15001"
	start := uint64(0)
	end := uint64(0)
	chSim.LoadValidatorSet(4, t, app.StakingKeeper, ctx, false, 10)
	app.BorKeeper.SetParams(ctx, types.DefaultParams())
	valSet := app.StakingKeeper.GetValidatorSet(ctx)
	app.CheckpointKeeper.UpdateACKCountWithValue(ctx, 1)
	producers, _ := app.BorKeeper.SelectNextProducers(ctx, hmTypes.ZeroHeimdallHash.EthHash())
	for i := 0; i < spancount; i++ {
		start = end + 1
		end = end + 10
		span := hmTypes.NewSpan(
			uint64(i+1),
			start,
			end,
			valSet,
			producers,
			chainID,
		)
		spans = append(spans, &span)
	}
	genesisState := types.NewGenesisState(
		params,
		spans,
	)
	bor.InitGenesis(ctx, app.BorKeeper, genesisState)
	actualParams := bor.ExportGenesis(ctx, app.BorKeeper)
	require.NotNil(t, actualParams)
	require.LessOrEqual(t, spancount, len(actualParams.Spans))
}
