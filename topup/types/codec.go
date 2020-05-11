package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTopup{}, "topup/MsgTopup", nil)
	cdc.RegisterConcrete(MsgWithdrawFee{}, "topup/MsgWithdrawFee", nil)
}

// ModuleCdc module cdc
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
