package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTopup{}, "topup/MsgTopup", nil)
	cdc.RegisterConcrete(MsgWithdrawFee{}, "bank/MsgWithdrawFee", nil)
}

// RegisterPulp register pulp
func RegisterPulp(pulp *authTypes.Pulp) {
	pulp.RegisterConcrete(MsgTopup{})
	pulp.RegisterConcrete(MsgWithdrawFee{})
}

// ModuleCdc module cdc
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
