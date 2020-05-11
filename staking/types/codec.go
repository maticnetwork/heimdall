package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// TODO we most likely dont need to register to amino as we are using RLP to encode

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgValidatorJoin{}, "staking/MsgValidatorJoin", nil)
	cdc.RegisterConcrete(MsgSignerUpdate{}, "staking/MsgSignerUpdate", nil)
	cdc.RegisterConcrete(MsgValidatorExit{}, "staking/MsgValidatorExit", nil)
	cdc.RegisterConcrete(MsgStakeUpdate{}, "staking/MsgStakeUpdate", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
