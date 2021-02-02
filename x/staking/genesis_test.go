package staking_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/maticnetwork/heimdall/x/staking/test_helper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/staking"
	"github.com/maticnetwork/heimdall/x/staking/types"

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
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(true)
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

// TestInitExportGenesis test import and export genesis state
func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5

	stakingSequence := make([]string, n)
	accounts := simulation.RandomAccounts(r1, n)

	for i := range stakingSequence {
		stakingSequence[i] = strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	}

	validators := make([]*hmTypes.Validator, n)
	for i := 0; i < len(validators); i++ {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			uint64(i),
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmCommonTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)
	}

	// validator set
	validatorSet := hmTypes.NewValidatorSet(validators)

	genesisState := types.NewGenesisState(validators, validatorSet, stakingSequence)
	staking.InitGenesis(ctx, initApp.StakingKeeper, *genesisState)

	actualParams := staking.ExportGenesis(ctx, initApp.StakingKeeper)
	require.NotNil(t, actualParams)
	require.LessOrEqual(t, 5, len(actualParams.Validators))
}
