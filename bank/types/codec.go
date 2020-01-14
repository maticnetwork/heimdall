package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "bank/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "bank/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgTopup{}, "bank/MsgTopup", nil)
	cdc.RegisterConcrete(MsgWithdrawFee{}, "bank/MsgWithdrawFee", nil)
}

// RegisterPulp register pulp
func RegisterPulp(pulp *authTypes.Pulp) {
	pulp.RegisterConcrete(MsgSend{})
	pulp.RegisterConcrete(MsgMultiSend{})
	pulp.RegisterConcrete(MsgTopup{})
	pulp.RegisterConcrete(MsgWithdrawFee{})

}

// ModuleCdc module cdc
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
