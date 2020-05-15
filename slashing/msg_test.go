package slashing_test

import (
	"encoding/hex"
	"testing"

	"github.com/maticnetwork/heimdall/helper"
	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func TestMsgTick(t *testing.T) {
	// create msg Tick message
	msg := slashingTypes.NewMsgTick(
		uint64(2),
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		hmTypes.HexToHexBytes("0xox"),
	)
	t.Log(hmTypes.BytesToHeimdallAddress(helper.GetAddress()))
	t.Log(hmTypes.HexToHexBytes("0xox"))

	t.Log(msg.Proposer)
	t.Log(msg.SlashingInfoBytes)

	t.Log(msg.Proposer.String())
	t.Log(msg.SlashingInfoBytes.String())

	t.Log(hex.EncodeToString(msg.GetSideSignBytes()))
}
