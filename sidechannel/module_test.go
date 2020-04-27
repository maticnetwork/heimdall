package sidechannel_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
	sidechannel "github.com/maticnetwork/heimdall/sidechannel"
	simulation "github.com/maticnetwork/heimdall/sidechannel/simulation"
	sidechannelTypes "github.com/maticnetwork/heimdall/sidechannel/types"
)

//
// Test suite
//

// ModuleTestSuite integrate test suite context object
type ModuleTestSuite struct {
	suite.Suite

	app    *app.HeimdallApp
	ctx    sdk.Context
	module sidechannel.AppModule
	r      *rand.Rand
}

func (suite *ModuleTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.module = sidechannel.NewAppModule(suite.app.SidechannelKeeper)

	// get random seed from time as source
	suite.r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestModuleTestSuite(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}

//
// Tests
//

func (suite *ModuleTestSuite) TestGetTxCmd() {
	t, module := suite.T(), suite.module

	require.Nil(t, module.GetTxCmd(sidechannelTypes.ModuleCdc))
}

func (suite *ModuleTestSuite) GetQueryCmd() {
	t, module := suite.T(), suite.module

	require.Nil(t, module.GetQueryCmd(sidechannelTypes.ModuleCdc))
}

func (suite *ModuleTestSuite) TestInitGenesis() {
	t, ctx, module := suite.T(), suite.ctx, suite.module

	data := sidechannelTypes.NewGenesisState([]sidechannelTypes.PastCommit{{Height: 23}})
	genesisState := sidechannelTypes.ModuleCdc.MustMarshalJSON(data)

	// init genesis
	require.NotPanics(t, func() {
		module.InitGenesis(ctx, genesisState)
	}, "Init genesis should not panic")

	data = sidechannelTypes.NewGenesisState([]sidechannelTypes.PastCommit{{Height: 122, Txs: []tmTypes.Tx{[]byte("test-tx122")}}})
	genesisState = sidechannelTypes.ModuleCdc.MustMarshalJSON(data)

	// init genesis
	require.NotPanics(t, func() {
		module.InitGenesis(ctx, genesisState)
	}, "Init genesis should not panic")
}

func (suite *ModuleTestSuite) TestExportGenesisWithoutPastCommit() {
	t, ctx, module := suite.T(), suite.ctx, suite.module

	gs := sidechannelTypes.DefaultGenesisState()
	genesisState := sidechannelTypes.ModuleCdc.MustMarshalJSON(gs)

	// init/export genesis
	module.InitGenesis(ctx, genesisState)
	actualParams := module.ExportGenesis(ctx)

	require.Equal(t, json.RawMessage(genesisState), actualParams, "Default export should be default genesis state")

	// genesis state with past commits
	gs1 := sidechannelTypes.NewGenesisState([]sidechannelTypes.PastCommit{{Height: 23}})
	genesisState1 := sidechannelTypes.ModuleCdc.MustMarshalJSON(gs1)

	// init/export genesis
	module.InitGenesis(ctx, genesisState1)
	actualParams = module.ExportGenesis(ctx)

	// check with empty genesis state
	require.Equal(t, json.RawMessage(genesisState), actualParams, "Default export should be valid genesis state with empty state")
}

func (suite *ModuleTestSuite) TestExportGenesis() {
	t, ctx, module := suite.T(), suite.ctx, suite.module

	// genesis state with past commits
	gs := sidechannelTypes.NewGenesisState(simulation.RandomPastCommits(suite.r, 2, 5, 5))
	genesisState := sidechannelTypes.ModuleCdc.MustMarshalJSON(gs)

	// init/export genesis
	module.InitGenesis(ctx, genesisState)

	// get genesis data
	var gs3 sidechannelTypes.GenesisState
	actualParams := module.ExportGenesis(ctx)
	err := sidechannelTypes.ModuleCdc.UnmarshalJSON(actualParams, &gs3)
	require.Nil(t, err)

	// check with empty genesis state
	require.Equal(t, len(gs.PastCommits), len(gs3.PastCommits))
	require.Equal(t, 2, len(gs3.PastCommits))
}

func (suite *ModuleTestSuite) TestBeginEndBlock() {
	t, r, app, ctx, module := suite.T(), suite.r, suite.app, suite.ctx, suite.module

	{
		// test start block

		var height int64 = 2
		ctx = ctx.WithBlockHeight(height)

		validators := app.SidechannelKeeper.GetValidators(ctx, height)
		require.Equal(t, 0, len(validators), fmt.Sprintf("It should have no validators before height %d", height))
		module.BeginBlock(ctx, abci.RequestBeginBlock{
			Header: abci.Header{
				Height: height,
			},
			LastCommitInfo: simulation.RandomLastCommitInfo(r, 10),
		})

		validators = app.SidechannelKeeper.GetValidators(ctx, height)
		require.Equal(t, 10, len(validators), fmt.Sprintf("It should have valid validators at height %d", height))
	}

	{
		// test start block

		var height int64 = 20
		ctx = ctx.WithBlockHeight(height)

		validators := app.SidechannelKeeper.GetValidators(ctx, height)
		require.Equal(t, 0, len(validators), fmt.Sprintf("It should have no validators before height %d", height))
		module.BeginBlock(ctx, abci.RequestBeginBlock{
			Header: abci.Header{
				Height: height,
			},
			LastCommitInfo: simulation.RandomLastCommitInfo(r, 10),
		})

		validators = app.SidechannelKeeper.GetValidators(ctx, height)
		require.Equal(t, 10, len(validators), fmt.Sprintf("It should have valid validators after begin block at height %d ", height))

		// test end block

		height = 20
		ctx = ctx.WithBlockHeight(height)

		module.EndBlock(ctx, abci.RequestEndBlock{
			Height: height,
		})

		validators = app.SidechannelKeeper.GetValidators(ctx, height)
		require.Equal(t, 0, len(validators), fmt.Sprintf("Validator should be empty after end block at height %d", height))
	}
}
