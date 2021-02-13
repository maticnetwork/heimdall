package keeper_test

import (
	"testing"

	"github.com/maticnetwork/heimdall/helper/mocks"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/bor/test_helper"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context

	contractCaller mocks.IContractCaller
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
