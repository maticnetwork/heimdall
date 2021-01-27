package gov_test

import (
	// "encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	// abci "github.com/tendermint/tendermint/abci/types"
	// "github.com/tendermint/tendermint/libs/log"
	// tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	// dbm "github.com/tendermint/tm-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	// "github.com/cosmos/cosmos-sdk/x/auth"
	// authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/gov"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// GenesisTestSuite integrate test suite context object
type GenesisTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *GenesisTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
}

// TestGenesisTestSuite
func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

// func (suite *GenesisTestSuite)TestImportExportQueues() {
// 	t, app, ctx := suite.T(), suite.app, suite.ctx
// 	addrs := simapp.AddTestAddrs(app, ctx, 2, valTokens)

// 	SortAddresses(addrs)

// 	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
// 	app.BeginBlock(abci.RequestBeginBlock{Header: header})

// 	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

// 	// Create two proposals, put the second into the voting period
// 	proposal := TestProposal
// 	proposal1, err := app.GovKeeper.SubmitProposal(ctx, proposal)
// 	require.NoError(t, err)
// 	proposalID1 := proposal1.ProposalId

// 	proposal2, err := app.GovKeeper.SubmitProposal(ctx, proposal)
// 	require.NoError(t, err)
// 	proposalID2 := proposal2.ProposalId

// 	votingStarted, err := app.GovKeeper.AddDeposit(ctx, proposalID2, addrs[0], app.GovKeeper.GetDepositParams(ctx).MinDeposit)
// 	require.NoError(t, err)
// 	require.True(t, votingStarted)

// 	proposal1, ok := app.GovKeeper.GetProposal(ctx, proposalID1)
// 	require.True(t, ok)
// 	proposal2, ok = app.GovKeeper.GetProposal(ctx, proposalID2)
// 	require.True(t, ok)
// 	require.True(t, proposal1.Status == types.StatusDepositPeriod)
// 	require.True(t, proposal2.Status == types.StatusVotingPeriod)

// 	authGenState := auth.ExportGenesis(ctx, app.AccountKeeper)
// 	bankGenState := app.BankKeeper.ExportGenesis(ctx)

// 	// export the state and import it into a new app
// 	govGenState := gov.ExportGenesis(ctx, app.GovKeeper)
// 	genesisState := types.DefaultGenesis()

// 	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenState)
// 	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenState)
// 	genesisState[types.ModuleName] = app.AppCodec().MustMarshalJSON(govGenState)

// 	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
// 	if err != nil {
// 		panic(err)
// 	}

// 	db := dbm.NewMemDB()
// 	app2 := simapp.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{})

// 	app2.InitChain(
// 		abci.RequestInitChain{
// 			Validators:      []abci.ValidatorUpdate{},
// 			ConsensusParams: simapp.DefaultConsensusParams,
// 			AppStateBytes:   stateBytes,
// 		},
// 	)

// 	app2.Commit()
// 	app2.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app2.LastBlockHeight() + 1}})

// 	header = tmproto.Header{Height: app2.LastBlockHeight() + 1}
// 	app2.BeginBlock(abci.RequestBeginBlock{Header: header})

// 	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})

// 	// Jump the time forward past the DepositPeriod and VotingPeriod
// 	ctx2 = ctx2.WithBlockTime(ctx2.BlockHeader().Time.Add(app2.GovKeeper.GetDepositParams(ctx2).MaxDepositPeriod).Add(app2.GovKeeper.GetVotingParams(ctx2).VotingPeriod))

// 	// Make sure that they are still in the DepositPeriod and VotingPeriod respectively
// 	proposal1, ok = app2.GovKeeper.GetProposal(ctx2, proposalID1)
// 	require.True(t, ok)
// 	proposal2, ok = app2.GovKeeper.GetProposal(ctx2, proposalID2)
// 	require.True(t, ok)
// 	require.True(t, proposal1.Status == types.StatusDepositPeriod)
// 	require.True(t, proposal2.Status == types.StatusVotingPeriod)

// 	macc := app2.GovKeeper.GetGovernanceAccount(ctx2)
// 	require.Equal(t, app2.GovKeeper.GetDepositParams(ctx2).MinDeposit, app2.BankKeeper.GetAllBalances(ctx2, macc.GetAddress()))

// 	// Run the endblocker. Check to make sure that proposal1 is removed from state, and proposal2 is finished VotingPeriod.
// 	gov.EndBlocker(ctx2, app2.GovKeeper)

// 	proposal1, ok = app2.GovKeeper.GetProposal(ctx2, proposalID1)
// 	require.False(t, ok)

// 	proposal2, ok = app2.GovKeeper.GetProposal(ctx2, proposalID2)
// 	require.True(t, ok)
// 	require.True(t, proposal2.Status == types.StatusRejected)
// }

// func(suite *GenesisTestSuite) TestEqualProposals() {
// 	app := simapp.Setup(false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
// 	addrs := simapp.AddTestAddrs(app, ctx, 2, valTokens)

// 	SortAddresses(addrs)

// 	// Submit two proposals
// 	proposal := TestProposal
// 	proposal1, err := app.GovKeeper.SubmitProposal(ctx, proposal)
// 	require.NoError(t, err)

// 	proposal2, err := app.GovKeeper.SubmitProposal(ctx, proposal)
// 	require.NoError(t, err)

// 	// They are similar but their IDs should be different
// 	require.NotEqual(t, proposal1, proposal2)
// 	require.NotEqual(t, proposal1, proposal2)

// 	// Now create two genesis blocks
// 	state1 := types.GenesisState{Proposals: []types.Proposal{proposal1}}
// 	state2 := types.GenesisState{Proposals: []types.Proposal{proposal2}}
// 	require.NotEqual(t, state1, state2)
// 	require.False(t, state1.Equal(state2))

// 	// Now make proposals identical by setting both IDs to 55
// 	proposal1.ProposalId = 55
// 	proposal2.ProposalId = 55
// 	require.Equal(t, proposal1, proposal1)
// 	require.Equal(t, proposal1, proposal2)

// 	// Reassign proposals into state
// 	state1.Proposals[0] = proposal1
// 	state2.Proposals[0] = proposal2

// 	// State should be identical now..
// 	require.Equal(t, state1, state2)
// 	require.True(t, state1.Equal(state2))
// }

//TestIsEmptyGenesis test empty genesis state
func (suite *GenesisTestSuite) TestIsEmptyGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()
	emptyGenesisState := types.GenesisState{}

	require.True(t, emptyGenesisState.IsEmpty())
	require.False(t, defaultGenesisState.IsEmpty())
}

//TestNewGenesis test new genesis state
func (suite *GenesisTestSuite) TestNewGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()
	newGenesisState := types.NewGenesisState(defaultGenesisState.StartingProposalId, defaultGenesisState.DepositParams, defaultGenesisState.VotingParams, defaultGenesisState.TallyParams)

	require.True(t, defaultGenesisState.Equal(newGenesisState))
}

//TestEqualGenesis test equal genesis state
func (suite *GenesisTestSuite) TestEqualGenesis() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState1 := types.DefaultGenesis()
	defaultGenesisState2 := types.DefaultGenesis()
	emptyGenesisState := types.GenesisState{}

	require.True(t, defaultGenesisState1.Equal(*defaultGenesisState2))
	require.False(t, defaultGenesisState1.Equal(emptyGenesisState))
}

//TestValidate test valid genesis state
func (suite *GenesisTestSuite) TestValidate() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()
	require.Nil(t, defaultGenesisState.Validate())

	defaultGenesisState.TallyParams.Threshold = sdk.NewDecWithPrec(-1, 0)
	require.NotNil(t, defaultGenesisState.Validate())

	defaultGenesisState = types.DefaultGenesis()
	defaultGenesisState.TallyParams.Veto = sdk.NewDecWithPrec(-1, 0)
	require.NotNil(t, defaultGenesisState.Validate())
}

//TestInitExportGenesis test import and export genesis state
func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	defaultGenesisState := types.DefaultGenesis()

	gov.InitGenesis(ctx, app.AccountKeeper, app.BankKeeper, app.GovKeeper, *defaultGenesisState)
	returnedGenesisState := gov.ExportGenesis(ctx, app.GovKeeper)

	require.True(t, defaultGenesisState.Equal(*returnedGenesisState))
}
