package bor_test

import (
	"testing"

	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"

	"github.com/maticnetwork/heimdall/x/bor"
	chSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"
	"github.com/stretchr/testify/require"

	"github.com/maticnetwork/heimdall/x/bor/test_helper"
	"github.com/maticnetwork/heimdall/x/bor/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *GenesisTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(true)
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	spanCount := 5
	params := types.DefaultParams()
	var spans []*hmTypes.Span
	chainID := "15001"
	start := uint64(0)
	end := uint64(0)
	valSet := chSim.LoadValidatorSet(4, t, initApp.StakingKeeper, ctx, false, 10)
	initApp.BorKeeper.SetParams(ctx, &params)
	initApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 1)
	producers, _ := initApp.BorKeeper.SelectNextProducers(ctx, hmCommonTypes.ZeroHeimdallHash.EthHash())
	for i := 0; i < spanCount; i++ {
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
	bor.InitGenesis(ctx, initApp.BorKeeper, types.GenesisState{
		Params: genesisState.Params,
		Spans:  genesisState.Spans,
	})
	actualParams := bor.ExportGenesis(ctx, initApp.BorKeeper)
	require.NotNil(t, actualParams)
	require.LessOrEqual(t, spanCount, len(actualParams.Spans))
}
