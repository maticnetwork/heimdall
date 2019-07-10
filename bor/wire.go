package bor

package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

func RegisterWire(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgProposeSpan{}, "bor/MsgProposeSpan", nil)
}

func RegisterPulp(pulp *hmTypes.Pulp) {
	pulp.RegisterConcrete(MsgProposeSpan{})
}

var cdcEmpty = codec.New()

func init() {
	RegisterWire(cdcEmpty)
	codec.RegisterCrypto(cdcEmpty)
}
