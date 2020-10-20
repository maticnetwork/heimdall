package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

var _ sdk.Msg = &MsgComment{}

func NewMsgComment(creator sdk.AccAddress, postID string, body string) *MsgComment {
  return &MsgComment{
    Id: uuid.New().String(),
		Creator: creator,
    PostID: postID,
    Body: body,
	}
}

func (msg *MsgComment) Route() string {
  return RouterKey
}

func (msg *MsgComment) Type() string {
  return "CreateComment"
}

func (msg *MsgComment) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg *MsgComment) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgComment) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}
