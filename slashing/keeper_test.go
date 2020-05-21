package slashing_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Test suite
//

// KeeperTestSuite integrate test suite context object
type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

//
// Tests
//

func (suite *KeeperTestSuite) TestMissedBlockBitArray() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	slashingKeeper := app.SlashingKeeper
	valID := hmTypes.NewValidatorID(uint64(1))
	index := int64(0)
	missed := false
	slashingKeeper.SetValidatorMissedBlockBitArray(ctx, valID, index, missed)
	response := slashingKeeper.GetValidatorMissedBlockBitArray(ctx, valID, index)
	require.Equal(t, missed, response)

	missed = true
	index = int64(1)
	slashingKeeper.SetValidatorMissedBlockBitArray(ctx, valID, index, missed)
	response = slashingKeeper.GetValidatorMissedBlockBitArray(ctx, valID, index)
	require.Equal(t, missed, response)
}

func (suite *KeeperTestSuite) TestValSigningInfo() {

	t, app, ctx := suite.T(), suite.app, suite.ctx
	slashingKeeper := app.SlashingKeeper

}
