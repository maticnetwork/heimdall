package slashing_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	_, app, ctx := suite.T(), suite.app, suite.ctx
	slashingKeeper := app.SlashingKeeper
	// params := slashingKeeper.GetParams(ctx)
	valID := hmTypes.NewValidatorID(uint64(1))
	index := int64(0)
	missed := false
	slashingKeeper.SetValidatorMissedBlockBitArray(ctx, valID, index, missed)
	response := slashingKeeper.GetValidatorMissedBlockBitArray(ctx, valID, index)

	fmt.Println("response -", response)
}
