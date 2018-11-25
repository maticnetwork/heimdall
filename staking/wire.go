package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// TODO we most likely dont need to register to amino as we are using RLP to encode

func RegisterWire(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgValidatorExit{}, "staking/MsgValidatorExit", nil)
	cdc.RegisterConcrete(MsgValidatorJoin{}, "staking/MsgValidatorJoin", nil)
	cdc.RegisterConcrete(MsgSignerUpdate{}, "staking/MsgSignerUpdate", nil)
}

var cdcEmpty = codec.New()

func init() {
	RegisterWire(cdcEmpty)
	codec.RegisterCrypto(cdcEmpty)
}
