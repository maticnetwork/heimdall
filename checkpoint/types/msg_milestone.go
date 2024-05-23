package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

var _ sdk.Msg = &MsgMilestone{}

// MsgMilestone represents milestone
type MsgMilestone struct {
	Proposer    types.HeimdallAddress `json:"proposer"`
	StartBlock  uint64                `json:"start_block"`
	EndBlock    uint64                `json:"end_block"`
	Hash        types.HeimdallHash    `json:"hash"`
	BorChainID  string                `json:"bor_chain_id"`
	MilestoneID string                `json:"milestone_id"`
}

// NewMsgMilestoneBlock creates new milestone message using mentioned arguments
func NewMsgMilestoneBlock(
	proposer types.HeimdallAddress,
	startBlock uint64,
	endBlock uint64,
	hash types.HeimdallHash,
	borChainID string,
	milestoneID string,
) MsgMilestone {
	return MsgMilestone{
		Proposer:    proposer,
		StartBlock:  startBlock,
		EndBlock:    endBlock,
		Hash:        hash,
		BorChainID:  borChainID,
		MilestoneID: milestoneID,
	}
}

// Type returns message type
func (msg MsgMilestone) Type() string {
	return EventTypeMilestone
}

func (msg MsgMilestone) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgMilestone) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.Proposer)}
}

func (msg MsgMilestone) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgMilestone) ValidateBasic() sdk.Error {
	if bytes.Equal(msg.Hash.Bytes(), helper.ZeroHash.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid rootHash %v", msg.Hash.String())
	}

	if msg.Proposer.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.Proposer.String())
	}

	if msg.StartBlock >= msg.EndBlock || msg.EndBlock == 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid startBlock %v or/and endBlock %v", msg.StartBlock, msg.EndBlock)
	}

	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgMilestone) GetSideSignBytes() []byte {
	return nil
}

var _ sdk.Msg = &MsgMilestoneTimeout{}

type MsgMilestoneTimeout struct {
	From types.HeimdallAddress `json:"from"`
}

func NewMsgMilestoneTimeout(from types.HeimdallAddress) MsgMilestoneTimeout {
	return MsgMilestoneTimeout{
		From: from,
	}
}

func (msg MsgMilestoneTimeout) Type() string {
	return "milestone-timeout"
}

func (msg MsgMilestoneTimeout) Route() string {
	return RouterKey
}

func (msg MsgMilestoneTimeout) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgMilestoneTimeout) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgMilestoneTimeout) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid from %v", msg.From.String())
	}

	return nil
}
