package slashing_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// Tests

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestMsgTick() {
	t, _, _ := suite.T(), suite.app, suite.ctx

	// create msg Tick message
	msg := slashingTypes.NewMsgTick(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		hmTypes.ZeroHeimdallHash,
	)
	t.Log(hmTypes.BytesToHeimdallAddress(helper.GetAddress()))
	t.Log(hmTypes.ZeroHeimdallHash)

	t.Log(msg.Proposer)
	t.Log(msg.SlashingInfoHash)

	t.Log(msg.Proposer.String())
	t.Log(msg.SlashingInfoHash.String())

	t.Log(hex.EncodeToString(msg.GetSideSignBytes()))
}
