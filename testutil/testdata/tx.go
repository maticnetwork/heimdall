package testdata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.MsgRequest = &SideMsgCreateDog{}

func (msg *SideMsgCreateDog) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{} }
func (msg *SideMsgCreateDog) ValidateBasic() error         { return nil }
func (msg *SideMsgCreateDog) GetSideSignBytes() []byte     { return nil }

func NewServiceSideMsgCreateDog(msg *SideMsgCreateDog) sdk.Msg {
	return sdk.ServiceMsg{
		MethodName: "/testutil.testdata.Msg/CreateDog",
		Request:    msg,
	}
}
