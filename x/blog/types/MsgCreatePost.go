package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

var _ sdk.Msg = &MsgPost{}

func NewMsgPost(creator sdk.AccAddress, title string, body string) *MsgPost {
  return &MsgPost{
    Id: uuid.New().String(),
		Creator: creator,
    Title: title,
    Body: body,
	}
}

func (msg *MsgPost) Route() string {
  return RouterKey
}

func (msg *MsgPost) Type() string {
  return "CreatePost"
}

func (msg *MsgPost) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg *MsgPost) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgPost) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}
