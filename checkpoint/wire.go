package checkpoint

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// TODO we most likely dont need to register to amino as we are using RLP to encode

func RegisterWire(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCheckpoint{}, "checkpoint/MsgCheckpoint", nil)
	cdc.RegisterConcrete(MsgCheckpointAck{}, "checkpoint/MsgCheckpointACK", nil)
}

var cdcEmpty = codec.New()

func init() {
	RegisterWire(cdcEmpty)
	codec.RegisterCrypto(cdcEmpty)
}
