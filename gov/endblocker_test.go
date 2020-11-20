package gov_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/gov"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
// Test suite
//

// EndBlockerTestSuite integrate test suite context object
type EndBlockerTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *EndBlockerTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestEndBlockerTestSuite(t *testing.T) {
	suite.Run(t, new(EndBlockerTestSuite))
}

func (suite *EndBlockerTestSuite) TestEndBlocker() {
	_, app, ctx := suite.T(), suite.app, suite.ctx
	gov.EndBlocker(ctx, app.GovKeeper)
}