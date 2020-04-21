package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	coinsPos         = sdk.NewCoins(hmTypes.NewInt64Coin(authTypes.FeeToken, 1000))
	coinsZero        = sdk.NewCoins()
	coinsPosNotMatic = sdk.NewCoins(hmTypes.NewInt64Coin("foo", 10000))
	coinsMulti       = sdk.NewCoins(hmTypes.NewInt64Coin(authTypes.FeeToken, 1000), hmTypes.NewInt64Coin("foo", 10000))
	addrs            = []hmTypes.HeimdallAddress{
		hmTypes.SampleHeimdallAddress("test1"),
		hmTypes.SampleHeimdallAddress("test2"),
	}
)

func init() {
	coinsMulti.Sort()
}

// test ValidateBasic for MsgCreateValidator
func TestMsgSubmitProposal(t *testing.T) {
	tests := []struct {
		title, description string
		proposalType       string
		proposerAddr       hmTypes.HeimdallAddress
		initialDeposit     sdk.Coins
		expectPass         bool
	}{
		{"Test Proposal", "the purpose of this proposal is to test", ProposalTypeText, addrs[0], coinsPos, true},
		{"", "the purpose of this proposal is to test", ProposalTypeText, addrs[0], coinsPos, false},
		{"Test Proposal", "", ProposalTypeText, addrs[0], coinsPos, false},
		{"Test Proposal", "the purpose of this proposal is to test", ProposalTypeSoftwareUpgrade, addrs[0], coinsPos, false},
		{"Test Proposal", "the purpose of this proposal is to test", ProposalTypeText, hmTypes.HeimdallAddress{}, coinsPos, false},
		{"Test Proposal", "the purpose of this proposal is to test", ProposalTypeText, addrs[0], coinsZero, true},
		{"Test Proposal", "the purpose of this proposal is to test", ProposalTypeText, addrs[0], coinsMulti, true},
		{strings.Repeat("#", MaxTitleLength*2), "the purpose of this proposal is to test", ProposalTypeText, addrs[0], coinsMulti, false},
		{"Test Proposal", strings.Repeat("#", MaxDescriptionLength*2), ProposalTypeText, addrs[0], coinsMulti, false},
	}

	for i, tc := range tests {
		msg := NewMsgSubmitProposal(
			ContentFromProposalType(tc.title, tc.description, tc.proposalType),
			tc.initialDeposit,
			tc.proposerAddr,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgDepositGetSignBytes(t *testing.T) {
	addr := hmTypes.SampleHeimdallAddress("addr1")
	msg := NewMsgDeposit(addr, 0, coinsPos)
	res := msg.GetSignBytes()

	expected := `{"type":"heimdall/MsgDeposit","value":{"amount":[{"amount":"1000","denom":"matic"}],"depositor":"0x0000000000000000000000000000006164647231","proposal_id":"0"}}`
	require.Equal(t, expected, string(res))
}

// test ValidateBasic for MsgDeposit
func TestMsgDeposit(t *testing.T) {
	tests := []struct {
		proposalID    uint64
		depositorAddr hmTypes.HeimdallAddress
		depositAmount sdk.Coins
		expectPass    bool
	}{
		{0, addrs[0], coinsPos, true},
		{1, hmTypes.HeimdallAddress{}, coinsPos, false},
		{1, addrs[0], coinsZero, true},
		{1, addrs[0], coinsMulti, true},
	}

	for i, tc := range tests {
		msg := NewMsgDeposit(tc.depositorAddr, tc.proposalID, tc.depositAmount)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

// test ValidateBasic for MsgDeposit
func TestMsgVote(t *testing.T) {
	tests := []struct {
		proposalID uint64
		voterAddr  hmTypes.HeimdallAddress
		option     VoteOption
		expectPass bool
	}{
		{0, addrs[0], OptionYes, true},
		{0, hmTypes.HeimdallAddress{}, OptionYes, false},
		{0, addrs[0], OptionNo, true},
		{0, addrs[0], OptionNoWithVeto, true},
		{0, addrs[0], OptionAbstain, true},
		{0, addrs[0], VoteOption(0x13), false},
	}

	for i, tc := range tests {
		msg := NewMsgVote(tc.voterAddr, tc.proposalID, tc.option)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
