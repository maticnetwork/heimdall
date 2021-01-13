package types

import (
	"fmt"
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// Governance message types and routes
const (
	TypeMsgDeposit        = "deposit"
	TypeMsgVote           = "vote"
	TypeMsgSubmitProposal = "submit_proposal"
)

var cdc = codec.NewLegacyAmino()

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

func (msg MsgVote) Route() string { return RouterKey }
func (msg MsgVote) Type() string  { return TypeMsgVote }

func (msg MsgSubmitProposal) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

func (msg MsgSubmitProposal) Route() string { return RouterKey }
func (msg MsgSubmitProposal) Type() string  { return TypeMsgSubmitProposal }

// Implements Msg.
func (m MsgSubmitProposal) ValidateBasic() error {
	if m.Proposer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Proposer.String())
	}
	if !m.InitialDeposit.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.InitialDeposit.String())
	}
	if m.InitialDeposit.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.InitialDeposit.String())
	}

	content := m.GetContent()
	if content == nil {
		return sdkerrors.Wrap(ErrInvalidProposalContent, "missing content")
	}
	if !IsValidProposalType(content.ProposalType()) {
		return sdkerrors.Wrap(ErrInvalidProposalType, content.ProposalType())
	}
	if err := content.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// Implements Msg.
func (msg MsgVote) ValidateBasic() error {
	if msg.Voter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Voter.String())
	}
	if !ValidVoteOption(msg.Option) {
		return sdkerrors.Wrap(ErrInvalidVote, msg.Option.String())
	}

	return nil
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

func (msg MsgDeposit) Route() string { return RouterKey }
func (msg MsgDeposit) Type() string  { return TypeMsgDeposit }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() error {
	if msg.Depositor.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

func (m *MsgSubmitProposal) GetProposer() sdk.AccAddress {
	return m.Proposer
}

func (m *MsgSubmitProposal) GetContent() Content {
	content, ok := m.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

func (m *MsgSubmitProposal) GetInitialDeposit() sdk.Coins { return m.InitialDeposit }

// NewMsgSubmitProposal creates new submit proposal
func NewMsgSubmitProposal(content Content, initialDeposit sdk.Coins, proposer sdk.AccAddress, validator hmTypes.ValidatorID) (*MsgSubmitProposal, error) {
	m := &MsgSubmitProposal{
		InitialDeposit: initialDeposit,
		Proposer:       proposer,
		Validator:      validator,
	}
	err := m.SetContent(content)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func NewMsgDeposit(depositor sdk.AccAddress, proposalID uint64, amount sdk.Coins, validator hmTypes.ValidatorID) MsgDeposit {
	return MsgDeposit{proposalID, depositor, amount, validator}
}

// NewMsgVote new msg vote
func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option VoteOption, validator hmTypes.ValidatorID) MsgVote {
	return MsgVote{proposalID, voter, option, validator}
}

func (m *MsgSubmitProposal) SetContent(content Content) error {
	msg, ok := content.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Content = any
	return nil
}
