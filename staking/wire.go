package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// TODO we most likely dont need to register to amino as we are using RLP to encode

func RegisterWire(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgValidatorJoin{}, "staking/MsgValidatorJoin", nil)
	cdc.RegisterConcrete(MsgSignerUpdate{}, "staking/MsgSignerUpdate", nil)
	cdc.RegisterConcrete(MsgValidatorExit{}, "staking/MsgValidatorExit", nil)
}

func RegisterPulp(pulp *authTypes.Pulp) {
	pulp.RegisterConcrete(MsgValidatorJoin{})
	pulp.RegisterConcrete(MsgSignerUpdate{})
	pulp.RegisterConcrete(MsgValidatorExit{})
}

var cdcEmpty = codec.New()

func init() {
	RegisterWire(cdcEmpty)
	codec.RegisterCrypto(cdcEmpty)
}
