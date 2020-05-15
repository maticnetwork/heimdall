package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "slashing/MsgUnjail", nil)
	cdc.RegisterConcrete(MsgTick{}, "slashing/MsgTick", nil)
	cdc.RegisterConcrete(MsgTickAck{}, "slashing/MsgTickAck", nil)

}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
