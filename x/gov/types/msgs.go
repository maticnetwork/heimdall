package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
func (msg MsgSubmitProposal) ValidateBasic() error {
	// TODO - Check this
	// if msg.Depositor.Empty() {
	// 	return sdk.ErrInvalidAddress(msg.Depositor.String())
	// }
	// if !msg.Amount.IsValid() {
	// 	return sdk.ErrInvalidCoins(msg.Amount.String())
	// }
	// if msg.Amount.IsAnyNegative() {
	// 	return sdk.ErrInvalidCoins(msg.Amount.String())
	// }
	// if msg.Validator == 0 {
	// 	return hmCommon.ErrInvalidMsg(DefaultCodespace, "Invalid validator id")
	// }

	return nil
}

// Implements Msg.
func (msg MsgVote) ValidateBasic() error {
	// TODO - Check this
	// if msg.Depositor.Empty() {
	// 	return sdk.ErrInvalidAddress(msg.Depositor.String())
	// }
	// if !msg.Amount.IsValid() {
	// 	return sdk.ErrInvalidCoins(msg.Amount.String())
	// }
	// if msg.Amount.IsAnyNegative() {
	// 	return sdk.ErrInvalidCoins(msg.Amount.String())
	// }
	// if msg.Validator == 0 {
	// 	return hmCommon.ErrInvalidMsg(DefaultCodespace, "Invalid validator id")
	// }

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
	// TODO - Check this
	// if msg.Depositor.Empty() {
	// 	return sdk.ErrInvalidAddress(msg.Depositor.String())
	// }
	// if !msg.Amount.IsValid() {
	// 	return sdk.ErrInvalidCoins(msg.Amount.String())
	// }
	// if msg.Amount.IsAnyNegative() {
	// 	return sdk.ErrInvalidCoins(msg.Amount.String())
	// }
	// if msg.Validator == 0 {
	// 	return hmCommon.ErrInvalidMsg(DefaultCodespace, "Invalid validator id")
	// }

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

func (m *MsgSubmitProposal) GetInitialDeposit() Coins { return m.InitialDeposit }
