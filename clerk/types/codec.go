package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgEventRecord{}, "cosmos-sdk/MsgEventRecord", nil)
}

// RegisterPulp register pulp
func RegisterPulp(pulp *authTypes.Pulp) {
	pulp.RegisterConcrete(MsgEventRecord{})
}

// ModuleCdc module cdc
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
