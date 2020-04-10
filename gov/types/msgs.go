package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// Governance message types and routes
const (
	TypeMsgDeposit        = "deposit"
	TypeMsgVote           = "vote"
	TypeMsgSubmitProposal = "submit_proposal"
)

var _, _, _ sdk.Msg = MsgSubmitProposal{}, MsgDeposit{}, MsgVote{}

// MsgSubmitProposal represents submit proposal message
type MsgSubmitProposal struct {
	Content        Content                 `json:"content" yaml:"content"`
	InitialDeposit sdk.Coins               `json:"initial_deposit" yaml:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive
	Proposer       hmTypes.HeimdallAddress `json:"proposer" yaml:"proposer"`               //  Address of the proposer
	Validator      hmTypes.ValidatorID     `json:"validator" yaml:"validator"`             //  Validator id
}

// NewMsgSubmitProposal creates new submit proposal
func NewMsgSubmitProposal(content Content, initialDeposit sdk.Coins, proposer hmTypes.HeimdallAddress, validator hmTypes.ValidatorID) MsgSubmitProposal {
	return MsgSubmitProposal{content, initialDeposit, proposer, validator}
}

//nolint
func (msg MsgSubmitProposal) Route() string { return RouterKey }
func (msg MsgSubmitProposal) Type() string  { return TypeMsgSubmitProposal }

// Implements Msg.
func (msg MsgSubmitProposal) ValidateBasic() sdk.Error {
	if msg.Content == nil {
		return ErrInvalidProposalContent(DefaultCodespace, "missing content")
	}

	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if msg.InitialDeposit.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if !IsValidProposalType(msg.Content.ProposalType()) {
		return ErrInvalidProposalType(DefaultCodespace, msg.Content.ProposalType())
	}
	if msg.Validator == 0 {
		return hmCommon.ErrInvalidMsg(DefaultCodespace, "Invalid validator id", "proposalType", msg.Content.ProposalType())
	}

	return msg.Content.ValidateBasic()
}

func (msg MsgSubmitProposal) String() string {
	return fmt.Sprintf(`Submit Proposal Message:
  Content:         %s
  Initial Deposit: %s
  Validator: %s
`, msg.Content.String(), msg.InitialDeposit, msg.Validator.String())
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.Proposer)}
}

// MsgDeposit represents deposit message
type MsgDeposit struct {
	ProposalID uint64                  `json:"proposal_id" yaml:"proposal_id"` // ID of the proposal
	Depositor  hmTypes.HeimdallAddress `json:"depositor" yaml:"depositor"`     // Address of the depositor
	Amount     sdk.Coins               `json:"amount" yaml:"amount"`           // Coins to add to the proposal's deposit
	Validator  hmTypes.ValidatorID     `json:"validator" yaml:"validator"`     //  Validator id
}

func NewMsgDeposit(depositor hmTypes.HeimdallAddress, proposalID uint64, amount sdk.Coins, validator hmTypes.ValidatorID) MsgDeposit {
	return MsgDeposit{proposalID, depositor, amount, validator}
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Route() string { return RouterKey }
func (msg MsgDeposit) Type() string  { return TypeMsgDeposit }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.Validator == 0 {
		return hmCommon.ErrInvalidMsg(DefaultCodespace, "Invalid validator id")
	}

	return nil
}

func (msg MsgDeposit) String() string {
	return fmt.Sprintf(`Deposit Message:
  Depositer:   %s
  Proposal ID: %d
  Amount:      %s
  Validator:      %s
`, msg.Depositor, msg.ProposalID, msg.Amount, msg.Validator.String())
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.Depositor)}
}

// MsgVote
type MsgVote struct {
	ProposalID uint64                  `json:"proposal_id" yaml:"proposal_id"` // ID of the proposal
	Voter      hmTypes.HeimdallAddress `json:"voter" yaml:"voter"`             //  address of the voter
	Option     VoteOption              `json:"option" yaml:"option"`           //  option from OptionSet chosen by the voter
	Validator  hmTypes.ValidatorID     `json:"validator" yaml:"validator"`     //  validator id of the voter
}

// NewMsgVote new msg vote
func NewMsgVote(voter hmTypes.HeimdallAddress, proposalID uint64, option VoteOption, validator hmTypes.ValidatorID) MsgVote {
	return MsgVote{proposalID, voter, option, validator}
}

// Implements Msg.
func (msg MsgVote) Route() string { return RouterKey }
func (msg MsgVote) Type() string  { return TypeMsgVote }

// Implements Msg.
func (msg MsgVote) ValidateBasic() sdk.Error {
	if msg.Voter.Empty() {
		return sdk.ErrInvalidAddress(msg.Voter.String())
	}
	if !ValidVoteOption(msg.Option) {
		return ErrInvalidVote(DefaultCodespace, msg.Option)
	}
	if msg.Validator == 0 {
		return hmCommon.ErrInvalidMsg(DefaultCodespace, "Invalid validator id")
	}

	return nil
}

func (msg MsgVote) String() string {
	return fmt.Sprintf(`Vote Message:
  Proposal ID: %d
  Option:      %s
  Validator:   %s
`, msg.ProposalID, msg.Option, msg.Validator.String())
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.Voter)}
}
