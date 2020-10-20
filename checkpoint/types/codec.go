package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCheckpoint{}, "checkpoint/MsgCheckpoint", nil)
	cdc.RegisterConcrete(MsgCheckpointAck{}, "checkpoint/MsgCheckpointACK", nil)
	cdc.RegisterConcrete(MsgCheckpointNoAck{}, "checkpoint/MsgCheckpointNoACK", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.NewLegacyAmino()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
