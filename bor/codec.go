package bor

import (
	"github.com/cosmos/cosmos-sdk/codec"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgProposeSpan{}, "bor/MsgProposeSpan", nil)
}

func RegisterPulp(pulp *authTypes.Pulp) {
	pulp.RegisterConcrete(MsgProposeSpan{})
}

var cdcEmpty = codec.New()

func init() {
	RegisterCodec(cdcEmpty)
	codec.RegisterCrypto(cdcEmpty)
}
